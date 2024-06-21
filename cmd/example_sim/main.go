//go:build js && wasm

package main

import "github.com/umbralcalc/dexetera/pkg/wasm"

func main() {
	wasm.RunLoop()
}
