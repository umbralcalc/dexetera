//go:build js && wasm

package simio

import (
	"log"
	"net/http"
	"syscall/js"

	"github.com/gorilla/websocket"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"google.golang.org/protobuf/proto"
)

// RunAndServeWebsocket runs a simulation while serving a websocket with
// the io.WebsocketIOIteration for one of the state partition iterations.
func RunAndServeWebsocket(
	settings *simulator.Settings,
	implementations *simulator.Implementations,
	websocketPartitionIndex int,
	handle string,
	address string,
) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	http.HandleFunc(
		handle,
		func(w http.ResponseWriter, r *http.Request) {
			connection, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println("Error upgrading to WebSocket:", err)
				return
			}
			defer connection.Close()

			implementations.Partitions[websocketPartitionIndex].Iteration =
				NewWebsocketIOIteration(connection)
			coordinator := simulator.NewPartitionCoordinator(
				settings,
				implementations,
			)
			coordinator.Run()
		},
	)
	log.Fatal(http.ListenAndServe(address, nil))
}

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
			State:               state,
		},
	)
	if err != nil {
		panic(err)
	}
	j.callback.Invoke(sendBytes)
}

// GenerateRunClosure creates a function which runs the stochadex
// simulation engine given the provided configured inputs.
func GenerateRunClosure(
	settings *simulator.Settings,
	implementations *simulator.Implementations,
	websocketPartitionIndex int,
	handle string,
	address string,
) func(this js.Value, p []js.Value) interface{} {
	return func(this js.Value, p []js.Value) interface{} {
		implementations.OutputFunction = &JsCallbackOutputFunction{
			callback: p[0],
		}
		RunAndServeWebsocket(
			settings,
			implementations,
			websocketPartitionIndex,
			handle,
			address,
		)
		return nil
	}
}

// RegisterRun registers the simulation run function as a JavaScript
// function which can be called from JavaScript code.
func RegisterRun(
	settings *simulator.Settings,
	implementations *simulator.Implementations,
	websocketPartitionIndex int,
	handle string,
	address string,
) {
	run := GenerateRunClosure(
		settings,
		implementations,
		websocketPartitionIndex,
		handle,
		address,
	)
	js.Global().Set("run", js.FuncOf(run))
	select {}
}
