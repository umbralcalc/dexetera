//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/minimal_example"
	"github.com/umbralcalc/dexetera/pkg/simio"
)

func main() {
	simio.RegisterStep(minimal_example.NewMinimalExampleGame().GetConfig(), "", ":2112")
}
