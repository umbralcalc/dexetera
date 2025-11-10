//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/network_control"
	"github.com/umbralcalc/dexetera/pkg/simio"
)

func main() {
	simio.RegisterStep(network_control.NewNetworkControlGame().GetConfig(), "", ":2112")
}
