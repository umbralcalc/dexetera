//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func main() {
	settings := &simulator.Settings{}
	implementations := &simulator.Implementations{}
	websocketPartitionIndex := 0
	handle := ""
	address := ""
	simio.RegisterRun(
		settings,
		implementations,
		websocketPartitionIndex,
		handle,
		address,
	)
}
