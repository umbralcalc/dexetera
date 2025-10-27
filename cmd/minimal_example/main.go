//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/umbralcalc/dexetera/pkg/games"
	"github.com/umbralcalc/dexetera/pkg/simio"
)

func main() {
	js.Global().Get("console").Call("log", "Minimal example main function called")

	game := games.NewMinimalExampleGame()
	js.Global().Get("console").Call("log", "MinimalExampleGame created")

	configGen := game.GetConfigGenerator()
	settings, _ := configGen.GenerateConfigs()

	implementations := game.GetConfig().ImplementationConfig.ToImplementations()
	js.Global().Get("console").Call("log", "Settings and implementations generated")

	// Note: websocketPartitionIndex is 0 since counter_state is the first partition
	js.Global().Get("console").Call("log", "Calling RegisterStep")
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
