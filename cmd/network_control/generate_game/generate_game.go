package main

import (
	"github.com/umbralcalc/dexetera/pkg/game"
	"github.com/umbralcalc/dexetera/pkg/network_control"
)

func main() {
	game.GenerateGamePackage(network_control.NewNetworkControlGame(), "network_control/")
}
