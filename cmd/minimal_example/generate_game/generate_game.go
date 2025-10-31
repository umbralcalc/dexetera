package main

import (
	"github.com/umbralcalc/dexetera/pkg/game"
	"github.com/umbralcalc/dexetera/pkg/minimal_example"
)

func main() {
	game.GenerateGamePackage(minimal_example.NewMinimalExampleGame(), "minimal_example/")
}
