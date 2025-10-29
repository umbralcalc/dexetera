//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/umbralcalc/dexetera/pkg/games"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// onlyNamesCond filters outputs to only the given partition names
type onlyNamesCond struct{ allow map[string]struct{} }

func (o *onlyNamesCond) IsOutputStep(partitionName string, state []float64, cumulativeTimesteps float64) bool {
	_, ok := o.allow[partitionName]
	return ok
}

func makeOnlyNames(names []string) simulator.OutputCondition {
	m := make(map[string]struct{}, len(names))
	for _, n := range names {
		m[n] = struct{}{}
	}
	return &onlyNamesCond{allow: m}
}

func main() {
	js.Global().Get("console").Call("log", "Minimal example main function called")

	game := games.NewMinimalExampleGame()
	js.Global().Get("console").Call("log", "MinimalExampleGame created")

	// Use the simulation generator from the game config
	cfg := game.GetConfig()
	var gen *simulator.ConfigGenerator
	if cfg.SimulationGenerator != nil {
		gen = cfg.SimulationGenerator()
	} else {
		// Fallback to existing method if not provided
		gen = game.GetConfigGenerator()
	}

	settings, implementations := gen.GenerateConfigs()
	js.Global().Get("console").Call("log", "Settings and implementations generated from SimulationGenerator")

	// Overwrite output condition to only output given partitions configured by the user
	if len(cfg.ServerPartitionNames) > 0 {
		implementations.OutputCondition = makeOnlyNames(cfg.ServerPartitionNames)
	}

	// Resolve websocket partition index by name (fallback to 0)
	websocketPartitionIndex := 0
	if len(cfg.ServerPartitionNames) > 0 {
		// TODO: derive by name from settings when API exposes names in settings ordering
	}

	// Register the simulation step function
	js.Global().Get("console").Call("log", "Calling RegisterStep")
	simio.RegisterStep(settings, implementations, websocketPartitionIndex, "", ":2112")
}
