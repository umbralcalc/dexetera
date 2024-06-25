package io

import (
	sync "sync"

	"github.com/gorilla/websocket"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"google.golang.org/protobuf/proto"
)

// WebsocketIOIteration implements an iteration in the stochadex
// based on I/O with a WebSocket connection.
type WebsocketInputIteration struct {
	conn           *websocket.Conn
	mutex          *sync.Mutex
	sendPartitions []int64
}

func (w *WebsocketInputIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	w.sendPartitions =
		settings.OtherParams[partitionIndex].IntParams["send_partitions"]
}

func (w *WebsocketInputIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	// Send the configured states to the client
	w.mutex.Lock()
	for _, index := range w.sendPartitions {
		sendData, err := proto.Marshal(
			&PartitionState{
				CumulativeTimesteps: timestepsHistory.Values.AtVec(0),
				PartitionIndex:      int64(index),
				State:               stateHistories[index].Values.RawRowView(0),
			},
		)
		if err != nil {
			panic(err)
		}
		err = w.conn.WriteMessage(0, sendData)
		if err != nil {
			panic(err)
		}
	}

	// Read message from WebSocket connection
	_, readData, err := w.conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	var data PartitionState
	err = proto.Unmarshal(readData, &data)
	if err != nil {
		panic(err)
	}
	w.mutex.Unlock()

	return data.State
}

// NewWebsocketInputIteration creates a new WebsocketInputIteration
func NewWebsocketInputIteration(
	conn *websocket.Conn,
	mutex *sync.Mutex,
) *WebsocketInputIteration {
	return &WebsocketInputIteration{conn: conn, mutex: mutex}
}
