//go:build js && wasm

package simio

import (
	"sync"
	"syscall/js"

	"github.com/umbralcalc/dexetera/pkg/game"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"google.golang.org/protobuf/proto"
)

// JsCallbackOutputFunction sets the callback function which passes
// the simulation states to the surrounding JavaScript code.
type JsCallbackOutputFunction struct {
	callback *js.Value
}

func (j *JsCallbackOutputFunction) Output(
	partitionName string,
	state []float64,
	cumulativeTimesteps float64,
) {
	// Check if callback is set before trying to invoke it
	if j.callback == nil || j.callback.Type() != js.TypeFunction {
		return // Skip output if callback is not ready
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
	// Call the callback with the Uint8Array directly
	callback.Invoke(uint8Array)
}

// OnlyNamesCondition filters outputs to only the given partition names.
type OnlyNamesCondition struct{ allow map[string]struct{} }

func (o *OnlyNamesCondition) IsOutputStep(partitionName string, state []float64, cumulativeTimesteps float64) bool {
	_, ok := o.allow[partitionName]
	return ok
}

// NewOnlyNamesCondition creates a new OnlyNamesCondition.
func NewOnlyNamesCondition(names []string) *OnlyNamesCondition {
	m := make(map[string]struct{}, len(names))
	for _, n := range names {
		m[n] = struct{}{}
	}
	return &OnlyNamesCondition{allow: m}
}

// GenerateStepClosure creates a function which steps the stochadex
// simulation engine given the provided configured inputs.
func GenerateStepClosure(
	wg *sync.WaitGroup,
	callback *js.Value,
	coordinator *simulator.PartitionCoordinator,
	websocketPartitionIndices []int,
	handle string,
	address string,
) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		*callback = args[0]
		// Update action state from server if data is received
		if !args[1].IsNull() {
			var actionState ActionState
			stateBytes := make([]byte, args[1].Get("length").Int())
			js.CopyBytesToGo(stateBytes, args[1])
			err := proto.Unmarshal(stateBytes, &actionState)
			if err != nil {
				panic(err)
			}
			for _, index := range websocketPartitionIndices {
				coordinator.Iterators[index].
					Params.Set("action_state_values", actionState.Values)
			}
		}
		coordinator.Step(wg)
		return nil
	}
}

// RegisterStep registers the simulation step function as a JavaScript
// function which can be called from JavaScript code.
func RegisterStep(cfg *game.GameConfig, handle string, address string) {
	js.Global().Get("console").Call("log", "main function called")

	js.Global().Get("console").Call("log", "game created")

	// Use the simulation generator from the game config
	var gen *simulator.ConfigGenerator = cfg.SimulationGenerator()

	settings, implementations := gen.GenerateConfigs()
	js.Global().Get("console").Call("log", "Settings and implementations generated from SimulationGenerator")

	// Overwrite output condition to only output given partitions configured by the user
	if len(cfg.ServerPartitionNames) > 0 {
		implementations.OutputCondition = NewOnlyNamesCondition(cfg.ServerPartitionNames)
	}

	// Resolve websocket partition indices by name
	websocketPartitionIndices := make([]int, 0)
	for _, name := range cfg.ActionStatePartitionNames {
		for index, iteration := range settings.Iterations {
			if iteration.Name == name {
				websocketPartitionIndices = append(websocketPartitionIndices, index)
			}
		}
	}

	// Register the simulation step function
	js.Global().Get("console").Call("log", "Calling RegisterStep")

	// Add debugging
	js.Global().Get("console").Call("log", "RegisterStep called")

	var wg sync.WaitGroup
	var callback js.Value
	implementations.OutputFunction = &JsCallbackOutputFunction{
		callback: &callback,
	}

	js.Global().Get("console").Call("log", "Creating PartitionCoordinator")
	coordinator := simulator.NewPartitionCoordinator(
		settings,
		implementations,
	)

	js.Global().Get("console").Call("log", "Creating step closure")
	step := GenerateStepClosure(
		&wg,
		&callback,
		coordinator,
		websocketPartitionIndices,
		handle,
		address,
	)

	js.Global().Get("console").Call("log", "Registering stepSimulation function")
	js.Global().Set("stepSimulation", js.FuncOf(step))
	js.Global().Get("console").Call("log", "stepSimulation function registered successfully")

	select {}
}
