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
//	         decoded and applied through dispatcher, which updates the
//	         relevant partitions' `action_state_values` params before the
//	         step runs. Actions that cannot be decoded or applied are
//	         reported on the console and dropped; the step still runs.
//
// The closure then advances the simulation by one step and returns nil.
func GenerateStepClosure(
	callback *js.Value,
	stepper simulator.Stepper,
	dispatcher *ActionDispatcher,
) func(this js.Value, args []js.Value) interface{} {
	console := js.Global().Get("console")
	var actionState ActionState
	return func(this js.Value, args []js.Value) interface{} {
		*callback = args[0]
		if !args[1].IsNull() {
			stateBytes := make([]byte, args[1].Get("length").Int())
			js.CopyBytesToGo(stateBytes, args[1])
			if err := proto.Unmarshal(stateBytes, &actionState); err != nil {
				console.Call("error", "dexetera: decoding an action: "+err.Error())
			} else if err := dispatcher.Apply(&actionState); err != nil {
				console.Call("error", "dexetera: applying an action: "+err.Error())
			}
		}
		stepper.Step()
		return nil
	}
}

// RegisterStep is the wasm `main` for an example: it builds the stochadex
// simulation from cfg, wires the JS output callback in, and registers a
// `stepSimulation` global on `js.Global()`. It then blocks forever
// (`select {}`) so the Go runtime stays alive to service further calls.
//
// Unless the simulation chooses an execution strategy, it is stepped with
// stochadex's inline execution: in WebAssembly there is one thread, so the
// default strategy's per-step goroutines are pure overhead (12-15x slower
// per step, measured under Node). A misconfigured page — an action
// partition that does not exist or declare action_state_values, or sliders
// that send a different width than it declares — fails here, at startup.
func RegisterStep(cfg *dashboard.Config) {
	step, err := NewStepFunc(cfg)
	if err != nil {
		panic(err)
	}
	js.Global().Set("stepSimulation", js.FuncOf(step))
	select {}
}

// NewStepFunc builds the `stepSimulation` function RegisterStep registers:
// the simulation from cfg, stepped as RegisterStep describes, with actions
// applied through an ActionDispatcher. A misconfigured page is an error.
func NewStepFunc(cfg *dashboard.Config) (func(this js.Value, args []js.Value) interface{}, error) {
	settings, implementations := cfg.SimulationGenerator().GenerateConfigs()
	if err := dashboard.CheckActionWidths(cfg, settings); err != nil {
		return nil, err
	}

	// Restrict output to the partitions the Config declares as "server"
	// partitions, so neither the renderer nor any external action source
	// receives partitions that weren't explicitly opted in.
	if len(cfg.ServerPartitionNames) > 0 {
		implementations.OutputCondition = NewOnlyNamesCondition(cfg.ServerPartitionNames)
	}
	if implementations.ExecutionStrategy == nil {
		implementations.ExecutionStrategy = &simulator.InlineExecution{}
	}

	var callback js.Value
	implementations.OutputFunction = &JsCallbackOutputFunction{callback: &callback}

	coordinator := simulator.NewPartitionCoordinator(settings, implementations)
	dispatcher, err := NewActionDispatcher(coordinator, cfg.ActionStatePartitionNames)
	if err != nil {
		return nil, err
	}
	// The page runs until it is closed, so the stepper is never closed: the
	// JS callback output holds nothing to finalize.
	return GenerateStepClosure(&callback, coordinator.NewStepper(), dispatcher), nil
}
