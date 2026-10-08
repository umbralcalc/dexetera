//go:build js && wasm

package simio_test

import (
	"strings"
	"syscall/js"
	"testing"

	"github.com/umbralcalc/dexetera/pkg/dashboard"
	"github.com/umbralcalc/dexetera/pkg/growth"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"google.golang.org/protobuf/proto"
)

// These run in real WebAssembly under Node (GOOS=js GOARCH=wasm go test
// -exec=$(go env GOROOT)/lib/wasm/go_js_wasm_exec ./pkg/simio/), driving the
// registered step function the way runtime/worker.js does.

// recorder is a JS callback collecting the PartitionStates a step emits.
func recorder(t *testing.T) (js.Value, *[]*simulator.PartitionState) {
	states := &[]*simulator.PartitionState{}
	callback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		bytes := make([]byte, args[0].Get("length").Int())
		js.CopyBytesToGo(bytes, args[0])
		state := &simulator.PartitionState{}
		if err := proto.Unmarshal(bytes, state); err != nil {
			t.Errorf("unmarshal: %v", err)
		}
		*states = append(*states, state)
		return nil
	})
	return callback.Value, states
}

func uint8Array(bytes []byte) js.Value {
	array := js.Global().Get("Uint8Array").New(len(bytes))
	js.CopyBytesToJS(array, bytes)
	return array
}

func TestStepFunctionInWasm(t *testing.T) {
	step, err := simio.NewStepFunc(growth.NewConfig())
	if err != nil {
		t.Fatal(err)
	}
	callback, states := recorder(t)

	// A step with no action input emits the population.
	step(js.Undefined(), []js.Value{callback, js.Null()})
	if len(*states) != 1 || (*states)[0].PartitionName != "population" {
		t.Fatalf("first step emitted %v, want one population state", *states)
	}

	// A carrying-capacity cut below the population (K=5 from N=10) turns
	// growth into decay towards K: the action reached the partition before
	// the step ran.
	before := (*states)[0].State[0]
	cut, _ := proto.Marshal(&simio.ActionState{Partitions: map[string]*simio.ActionValues{
		"population": {Values: []float64{0.1, 5}}}})
	step(js.Undefined(), []js.Value{callback, uint8Array(cut)})
	after := (*states)[1].State[0]
	if !(after < before) {
		t.Errorf("population went from %v to %v; the K=5 action should shrink it", before, after)
	}

	// Bytes that are not an ActionState, or a vector of the wrong width, are
	// dropped (reported on the console) and the step still runs on the held
	// actions.
	wrong, _ := proto.Marshal(&simio.ActionState{Values: []float64{1, 2, 3}})
	for _, bad := range [][]byte{{0xff, 0xff, 0xff}, wrong} {
		step(js.Undefined(), []js.Value{callback, uint8Array(bad)})
	}
	if len(*states) != 4 {
		t.Fatalf("expected 4 steps' output, got %d", len(*states))
	}
	if last, previous := (*states)[3].State[0], (*states)[2].State[0]; !(last < previous && last > 5) {
		t.Errorf("after dropped actions the population went %v -> %v; under the held K=5 it should keep decaying towards 5",
			previous, last)
	}
}

func TestStepFunctionRejectsAMisconfiguredPage(t *testing.T) {
	cfg := growth.NewConfig()
	cfg.Sliders = append(cfg.Sliders, dashboard.Slider{Name: "extra", Partition: "population", ValueIndex: 2})
	if _, err := simio.NewStepFunc(cfg); err == nil ||
		!strings.Contains(err.Error(), "declares 2 action_state_values but its sliders send 3") {
		t.Errorf("expected a width error at startup, got %v", err)
	}
}
