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

	// Get the settings and implementations from the game
	settings := game.GetSettings()
	implementations := game.GetImplementations()
	js.Global().Get("console").Call("log", "Settings and implementations obtained")

	// Register the simulation step function
	// Note: websocketPartitionIndex is 0 since we only have one partition
	js.Global().Get("console").Call("log", "Calling RegisterStep")
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
