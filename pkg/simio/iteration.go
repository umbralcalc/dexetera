package simio

import (
	"github.com/gorilla/websocket"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"google.golang.org/protobuf/proto"
)

// WebsocketIOIteration implements an iteration in the stochadex
// based on I/O with a WebSocket connection.
type WebsocketIOIteration struct {
	conn           *websocket.Conn
	sendPartitions []int64
}

func (w *WebsocketIOIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	w.sendPartitions =
		settings.OtherParams[partitionIndex].IntParams["send_partitions"]
}

func (w *WebsocketIOIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	// Send the configured states to the client
	for _, index := range w.sendPartitions {
		sendBytes, err := proto.Marshal(
			&PartitionState{
				CumulativeTimesteps: timestepsHistory.Values.AtVec(0),
				PartitionIndex:      index,
				State: &State{
					Values: stateHistories[index].Values.RawRowView(0),
				},
			},
		)
		if err != nil {
			panic(err)
		}
		err = w.conn.WriteMessage(0, sendBytes)
		if err != nil {
			panic(err)
		}
	}

	// Read data from WebSocket connection
	_, readBytes, err := w.conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	var data State
	err = proto.Unmarshal(readBytes, &data)
	if err != nil {
		panic(err)
	}

	return data.Values
}

// NewWebsocketIOIteration creates a new WebsocketIOIteration
func NewWebsocketIOIteration(
	conn *websocket.Conn,
) *WebsocketIOIteration {
	return &WebsocketIOIteration{conn: conn}
}
