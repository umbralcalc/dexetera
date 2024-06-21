package wasm

import "github.com/umbralcalc/stochadex/pkg/simulator"

func Run() {
	settings := &simulator.Settings{}
	implementations := &simulator.Implementations{}
	_ = simulator.NewPartitionCoordinator(settings, implementations)
}
