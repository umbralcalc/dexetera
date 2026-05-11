//go:build js && wasm

// register_step is the growth example compiled as a WebAssembly module.
// Same shape as the other examples' register_step mains: it registers
// `stepSimulation` on the JS global, then blocks forever so the Go runtime
// stays alive to service per-step calls from runtime/worker.js.
//
// Build with growth/build.sh, or directly:
//
//	GOOS=js GOARCH=wasm go build -o growth/src/main.wasm ./cmd/growth/register_step
package main

import (
	"github.com/umbralcalc/dexetera/pkg/growth"
	"github.com/umbralcalc/dexetera/pkg/simio"
)

func main() {
	simio.RegisterStep(growth.NewConfig())
}
