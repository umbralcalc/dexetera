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
	gen = cfg.SimulationGenerator()

	settings, implementations := gen.GenerateConfigs()
	js.Global().Get("console").Call("log", "Settings and implementations generated from SimulationGenerator")

	// Overwrite output condition to only output given partitions configured by the user
	if len(cfg.ServerPartitionNames) > 0 {
		implementations.OutputCondition = makeOnlyNames(cfg.ServerPartitionNames)
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
	simio.RegisterStep(settings, implementations, websocketPartitionIndices, "", ":2112")
}
