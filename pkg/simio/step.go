//go:build js && wasm

// Package simio is the wasm-side runtime that hosts a stochadex simulation
// inside the browser. The browser worker (runtime/worker.js) loads the
// compiled module, calls RegisterStep at startup, then drives the
// simulation one step at a time by invoking the registered global JS
// function `stepSimulation(callback, actionStateBytes-or-null)`.
//
// Output flows in the opposite direction: each step the wasm module calls
// `callback(uint8Array)` once per output partition with a marshalled
// PartitionState protobuf message. The JS side decodes those messages and
// either renders them or forwards them to an external action source.
package simio

import (
	"sync"
	"syscall/js"

	"github.com/umbralcalc/dexetera/pkg/dashboard"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"google.golang.org/protobuf/proto"
)

// JsCallbackOutputFunction is a stochadex OutputFunction that delivers
// each output step to the surrounding JavaScript by invoking the most
// recently registered callback. The callback is set on every step (the
// first argument to stepSimulation), which is what lets the worker swap
// callbacks if it ever needs to.
type JsCallbackOutputFunction struct {
	callback *js.Value
}

func (j *JsCallbackOutputFunction) Configure(*simulator.Settings) {}

func (j *JsCallbackOutputFunction) Output(
	partitionName string,
	state []float64,
	cumulativeTimesteps float64,
) {
	if j.callback == nil || j.callback.Type() != js.TypeFunction {
		return
	}
	sendBytes, err := proto.Marshal(
		&simulator.PartitionState{
			CumulativeTimesteps: cumulativeTimesteps,
			PartitionName:       partitionName,
			State:               state,
		},
	)
	if err != nil {
		panic(err)
	}
	uint8Array := js.Global().Get("Uint8Array").New(len(sendBytes))
	js.CopyBytesToJS(uint8Array, sendBytes)
	callback := *j.callback
	callback.Invoke(uint8Array)
}

// OnlyNamesCondition is a stochadex OutputCondition that gates output to
// just the partitions whose names appear in `allow`. Used by RegisterStep
// so that only the partitions the GameConfig explicitly declares as
// "server partitions" are ever marshalled across the wasm/JS boundary.
type OnlyNamesCondition struct {
	allow map[string]struct{}
}

func (o *OnlyNamesCondition) IsOutputStep(
	partitionName string,
	state []float64,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) bool {
	_, ok := o.allow[partitionName]
	return ok
}

func NewOnlyNamesCondition(names []string) *OnlyNamesCondition {
	m := make(map[string]struct{}, len(names))
	for _, n := range names {
		m[n] = struct{}{}
	}
	return &OnlyNamesCondition{allow: m}
}

// GenerateStepClosure builds the JS-side step entrypoint.
//
// The returned function is registered as `stepSimulation` on the JS global
// scope. It expects two arguments on every call:
//
//	args[0]  the output callback to invoke for each emitted PartitionState
//	         this step. Re-set every step so the caller can swap it.
//	args[1]  either null (no new action input) or a Uint8Array of bytes
//	         encoding an ActionState protobuf. When present, the bytes are
//	         decoded and routed through ApplyActionState, which updates
//	         the relevant partitions' `action_state_values` params before
//	         the step runs.
//
// The closure then advances the coordinator by one step and returns nil.
func GenerateStepClosure(
	wg *sync.WaitGroup,
	callback *js.Value,
	coordinator *simulator.PartitionCoordinator,
	actionPartitionIndices []int,
	actionPartitionIndexByName map[string]int,
) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		*callback = args[0]
		if !args[1].IsNull() {
			var actionState ActionState
			stateBytes := make([]byte, args[1].Get("length").Int())
			js.CopyBytesToGo(stateBytes, args[1])
			if err := proto.Unmarshal(stateBytes, &actionState); err != nil {
				panic(err)
			}
			ApplyActionState(
				coordinator,
				actionPartitionIndices,
				actionPartitionIndexByName,
				&actionState,
			)
		}
		coordinator.Step(wg)
		return nil
	}
}

// RegisterStep is the wasm `main` for an example: it builds the stochadex
// coordinator from cfg, wires the JS output callback in, and registers a
// `stepSimulation` global on `js.Global()`. It then blocks forever
// (`select {}`) so the Go runtime stays alive to service further calls.
//
// The two index structures it builds — actionPartitionIndices (slice,
// declaration order) and actionPartitionIndexByName (map) — exist so that
// ApplyActionState can serve both action-delivery paths efficiently:
//   - Broadcast (legacy ActionState.Values): iterate the slice.
//   - Per-partition named (ActionState.Partitions): look up by name.
func RegisterStep(cfg *dashboard.Config) {
	settings, implementations := cfg.SimulationGenerator().GenerateConfigs()

	// Restrict output to the partitions the Config declares as "server"
	// partitions, so neither the renderer nor any external action source
	// receives partitions that weren't explicitly opted in.
	if len(cfg.ServerPartitionNames) > 0 {
		implementations.OutputCondition = NewOnlyNamesCondition(cfg.ServerPartitionNames)
	}

	actionPartitionIndices := make([]int, 0, len(cfg.ActionStatePartitionNames))
	actionPartitionIndexByName := make(map[string]int, len(cfg.ActionStatePartitionNames))
	for _, name := range cfg.ActionStatePartitionNames {
		for index, iteration := range settings.Iterations {
			if iteration.Name == name {
				actionPartitionIndices = append(actionPartitionIndices, index)
				actionPartitionIndexByName[name] = index
			}
		}
	}

	var wg sync.WaitGroup
	var callback js.Value
	implementations.OutputFunction = &JsCallbackOutputFunction{callback: &callback}

	coordinator := simulator.NewPartitionCoordinator(settings, implementations)
	step := GenerateStepClosure(
		&wg,
		&callback,
		coordinator,
		actionPartitionIndices,
		actionPartitionIndexByName,
	)

	js.Global().Set("stepSimulation", js.FuncOf(step))
	select {}
}
