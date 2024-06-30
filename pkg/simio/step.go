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
	callback js.Value
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
	j.callback.Invoke(uint8Array)
}

// GenerateStepClosure creates a function which steps the stochadex
// simulation engine given the provided configured inputs.
func GenerateStepClosure(
	settings *simulator.Settings,
	implementations *simulator.Implementations,
	websocketPartitionIndex int,
	handle string,
	address string,
) func(this js.Value, p []js.Value) interface{} {
	var wg sync.WaitGroup
	var coordinator *simulator.PartitionCoordinator
	return func(this js.Value, p []js.Value) interface{} {
		if coordinator != nil {
			implementations.OutputFunction = &JsCallbackOutputFunction{
				callback: p[0],
			}
			coordinator = simulator.NewPartitionCoordinator(
				settings,
				implementations,
			)
		}
		// Update action state from server if data is received
		if !p[1].IsNull() {
			var stateBytes []byte
			var actionState State
			js.CopyBytesToGo(stateBytes, p[1])
			proto.Unmarshal(stateBytes, &actionState)
			coordinator.Iterators[websocketPartitionIndex].
				Params.FloatParams["action"] = actionState.Values
		}
		coordinator.Step(&wg)
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
	step := GenerateStepClosure(
		settings,
		implementations,
		websocketPartitionIndex,
		handle,
		address,
	)
	js.Global().Set("stepSimulation", js.FuncOf(step))
	select {}
}
