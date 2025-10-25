//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/games"
	"github.com/umbralcalc/dexetera/pkg/simio"
)

func main() {
	// Create the simple demo game
	game := games.NewSimpleDemoGame()

	// Get the settings and implementations from the game
	settings := game.GetSettings()
	implementations := game.GetImplementations()

	// Register the simulation step function
	// Note: websocketPartitionIndex is 0 since we only have one partition
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
