//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/umbralcalc/dexetera/pkg/games"
	"github.com/umbralcalc/dexetera/pkg/simio"
)

func main() {
	// Add debugging
	js.Global().Get("console").Call("log", "Visualization example main function called")

	// Create the visualization example game using the VisualizationBuilder
	game := games.NewVisualizationExampleGame()
	js.Global().Get("console").Call("log", "Visualization example game created")

	// Get the config generator and generate settings
	configGen := game.GetConfigGenerator()
	settings, _ := configGen.GenerateConfigs() // We'll use our own implementations

	// Get implementations from our game configuration
	implementations := game.GetConfig().ImplementationConfig.ToImplementations()
	js.Global().Get("console").Call("log", "Settings and implementations generated from VisualizationBuilder")

	// Register the simulation step function
	// Note: websocketPartitionIndex is 0 since counter_state is the first partition
	js.Global().Get("console").Call("log", "Calling RegisterStep")
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
