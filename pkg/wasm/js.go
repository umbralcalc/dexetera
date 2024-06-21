//go:build js && wasm

package wasm

import (
	"syscall/js"
)

func loop(this js.Value, p []js.Value) interface{} {
	callback := p[0] // The JavaScript callback function
	for i := 0; i < 10; i++ {
		data := i * 2         // Example data computation
		callback.Invoke(data) // Pass data to the callback
	}
	return nil
}

func RunLoop() {
	js.Global().Set("loop", js.FuncOf(loop))
	select {} // Prevents the Go program from exiting
}
