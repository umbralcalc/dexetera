package wasm

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/umbralcalc/dexetera/pkg/io"
	"github.com/umbralcalc/stochadex/pkg/simulator"
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
				io.NewWebsocketIOIteration(connection)
			coordinator := simulator.NewPartitionCoordinator(
				settings,
				implementations,
			)
			coordinator.Run()
		},
	)
	log.Fatal(http.ListenAndServe(address, nil))
}
