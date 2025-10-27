//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/umbralcalc/dexetera/pkg/games"
	"github.com/umbralcalc/dexetera/pkg/simio"
)

func main() {
	// Add debugging
	js.Global().Get("console").Call("log", "Minimal example main function called")

	// Create the minimal example game
	game := games.NewMinimalExampleGame()
	js.Global().Get("console").Call("log", "Game created")

	// Get the config generator and generate settings and implementations
	configGen := game.GetConfigGenerator()
	settings, implementations := configGen.GenerateConfigs()
	js.Global().Get("console").Call("log", "Settings and implementations generated from ConfigGenerator")

	// Register the simulation step function
	// Note: websocketPartitionIndex is 0 since we only have one partition
	js.Global().Get("console").Call("log", "Calling RegisterStep")
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
