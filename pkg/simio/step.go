//go:build js && wasm

package simio

import (
	"sync"
	"syscall/js"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"google.golang.org/protobuf/proto"
)

// JsCallbackOutputFunction sets the callback function which passes
// the simulation states to the surrounding JavaScript code.
type JsCallbackOutputFunction struct {
	callback *js.Value
}

func (j *JsCallbackOutputFunction) Output(
	partitionIndex int,
	state []float64,
	cumulativeTimesteps float64,
) {
	sendBytes, err := proto.Marshal(
		&PartitionState{
			CumulativeTimesteps: cumulativeTimesteps,
			PartitionIndex:      int64(partitionIndex),
			State:               &State{Values: state},
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

// GenerateStepClosure creates a function which steps the stochadex
// simulation engine given the provided configured inputs.
func GenerateStepClosure(
	wg *sync.WaitGroup,
	callback *js.Value,
	coordinator *simulator.PartitionCoordinator,
	websocketPartitionIndex int,
	handle string,
	address string,
) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		*callback = args[0]
		// Update action state from server if data is received
		if !args[1].IsNull() {
			var stateBytes []byte
			var actionState State
			js.CopyBytesToGo(stateBytes, args[1])
			err := proto.Unmarshal(stateBytes, &actionState)
			if err != nil {
				panic(err)
			}
			coordinator.Iterators[websocketPartitionIndex].
				Params.FloatParams["action"] = actionState.Values
		}
		coordinator.Step(wg)
		return nil
	}
}

// RegisterStep registers the simulation step function as a JavaScript
// function which can be called from JavaScript code.
func RegisterStep(
	settings *simulator.Settings,
	implementations *simulator.Implementations,
	websocketPartitionIndex int,
	handle string,
	address string,
) {
	var wg sync.WaitGroup
	var callback js.Value
	implementations.OutputFunction = &JsCallbackOutputFunction{
		callback: &callback,
	}
	coordinator := simulator.NewPartitionCoordinator(
		settings,
		implementations,
	)
	step := GenerateStepClosure(
		&wg,
		&callback,
		coordinator,
		websocketPartitionIndex,
		handle,
		address,
	)
	js.Global().Set("stepSimulation", js.FuncOf(step))
	select {}
}
