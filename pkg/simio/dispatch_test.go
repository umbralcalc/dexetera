package simio_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/umbralcalc/dexetera/pkg/growth"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/floats"
)

// echoIteration's next state is its action_state_values (or, for a
// partition without actions, its "other" param), so each step's state shows
// exactly which actions that step read.
type echoIteration struct{}

func (*echoIteration) Configure(int, *simulator.Settings) {}

func (*echoIteration) Iterate(
	params *simulator.Params, _ int, _ []*simulator.StateHistory,
	_ *simulator.CumulativeTimestepsHistory,
) []float64 {
	if actions, ok := params.GetOk(simio.ActionParam); ok {
		return append([]float64(nil), actions...)
	}
	return append([]float64(nil), params.Get("other")...)
}

// buildCoordinator is two action partitions, "alpha" and "beta", each
// declaring a width-2 action vector, plus "gamma", which declares none.
func buildCoordinator(t *testing.T) *simulator.PartitionCoordinator {
	t.Helper()
	gen := simulator.NewConfigGenerator()
	for i, name := range []string{"alpha", "beta", "gamma"} {
		params := map[string][]float64{simio.ActionParam: {0, 0}}
		if name == "gamma" {
			params = map[string][]float64{"other": {0, 0}}
		}
		gen.SetPartition(&simulator.PartitionConfig{
			Name: name, Iteration: &echoIteration{}, Params: simulator.NewParams(params),
			InitStateValues: []float64{0, 0}, StateHistoryDepth: 1, Seed: uint64(101 + i),
		})
	}
	gen.SetSimulation(&simulator.SimulationConfig{
		OutputCondition: &simulator.NilOutputCondition{}, OutputFunction: &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{MaxNumberOfSteps: 100},
		TimestepFunction:     &simulator.ConstantTimestepFunction{Stepsize: 1.0},
		ExecutionStrategy:    &simulator.InlineExecution{},
	})
	return simulator.NewPartitionCoordinator(gen.GenerateConfigs())
}

// stepAndRead takes one step and returns each partition's state.
func stepAndRead(stepper simulator.Stepper, coordinator *simulator.PartitionCoordinator) map[string][]float64 {
	stepper.Step()
	out := map[string][]float64{}
	for _, iterator := range coordinator.Iterators {
		out[iterator.Partition.Name] = append([]float64(nil),
			coordinator.Shared.StateHistories[iterator.Partition.Index].Values.RawRowView(0)...)
	}
	return out
}

func TestActionDispatcher(t *testing.T) {
	t.Run("the broadcast path sets every action partition", func(t *testing.T) {
		coordinator := buildCoordinator(t)
		dispatcher, err := simio.NewActionDispatcher(coordinator, []string{"alpha", "beta"})
		if err != nil {
			t.Fatal(err)
		}
		if err := dispatcher.Apply(&simio.ActionState{Values: []float64{1, 2}}); err != nil {
			t.Fatal(err)
		}
		states := stepAndRead(coordinator.NewStepper(), coordinator)
		if !floats.Equal(states["alpha"], []float64{1, 2}) || !floats.Equal(states["beta"], []float64{1, 2}) {
			t.Errorf("states = %v, want both [1 2]", states)
		}
	})

	t.Run("the named path sets only the named partition, and takes precedence", func(t *testing.T) {
		coordinator := buildCoordinator(t)
		dispatcher, _ := simio.NewActionDispatcher(coordinator, []string{"alpha", "beta"})
		err := dispatcher.Apply(&simio.ActionState{
			Values:     []float64{9, 9}, // ignored: named entries take precedence
			Partitions: map[string]*simio.ActionValues{"beta": {Values: []float64{3, 4}}, "nobody": {Values: []float64{5}}},
		})
		if err != nil {
			t.Fatal(err)
		}
		states := stepAndRead(coordinator.NewStepper(), coordinator)
		if !floats.Equal(states["alpha"], []float64{0, 0}) || !floats.Equal(states["beta"], []float64{3, 4}) {
			t.Errorf("states = %v, want alpha [0 0] and beta [3 4]", states)
		}
	})

	t.Run("actions hold until the next ones, and a nil ActionState changes nothing", func(t *testing.T) {
		coordinator := buildCoordinator(t)
		dispatcher, _ := simio.NewActionDispatcher(coordinator, []string{"alpha"})
		stepper := coordinator.NewStepper()
		dispatcher.Apply(&simio.ActionState{Values: []float64{7, 8}})
		stepAndRead(stepper, coordinator)
		if err := dispatcher.Apply(nil); err != nil {
			t.Fatal(err)
		}
		if states := stepAndRead(stepper, coordinator); !floats.Equal(states["alpha"], []float64{7, 8}) {
			t.Errorf("alpha = %v, want the held [7 8]", states["alpha"])
		}
	})

	t.Run("a vector of the wrong width is an error and leaves the actions", func(t *testing.T) {
		coordinator := buildCoordinator(t)
		dispatcher, _ := simio.NewActionDispatcher(coordinator, []string{"alpha"})
		err := dispatcher.Apply(&simio.ActionState{Values: []float64{1, 2, 3}})
		if err == nil || !strings.Contains(err.Error(), "has width 2, got 3 values") {
			t.Fatalf("expected a width error, got %v", err)
		}
		if states := stepAndRead(coordinator.NewStepper(), coordinator); !floats.Equal(states["alpha"], []float64{0, 0}) {
			t.Errorf("alpha = %v, want its declared [0 0]", states["alpha"])
		}
	})

	t.Run("a named vector of the wrong width is an error too", func(t *testing.T) {
		coordinator := buildCoordinator(t)
		dispatcher, _ := simio.NewActionDispatcher(coordinator, []string{"alpha"})
		err := dispatcher.Apply(&simio.ActionState{Partitions: map[string]*simio.ActionValues{
			"alpha": {Values: []float64{1}}}})
		if err == nil || !strings.Contains(err.Error(), "has width 2, got 1 values") {
			t.Errorf("expected a width error, got %v", err)
		}
	})

	t.Run("an action partition that cannot take actions is an error at startup", func(t *testing.T) {
		for names, want := range map[string]string{
			"alpha,ghost": `action partition "ghost": simulator: InjectParams: no partition named "ghost"`,
			"gamma":       `action partition "gamma": simulator: InjectParams: partition "gamma" has no params key "action_state_values"`,
		} {
			_, err := simio.NewActionDispatcher(buildCoordinator(t), strings.Split(names, ","))
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s: expected %q, got %v", names, want, err)
			}
		}
	})

	t.Run("applying actions allocates nothing", func(t *testing.T) {
		coordinator := buildCoordinator(t)
		dispatcher, _ := simio.NewActionDispatcher(coordinator, []string{"alpha", "beta"})
		named := &simio.ActionState{Partitions: map[string]*simio.ActionValues{"alpha": {Values: []float64{1, 2}}}}
		broadcast := &simio.ActionState{Values: []float64{3, 4}}
		if allocs := testing.AllocsPerRun(1000, func() {
			dispatcher.Apply(named)
			dispatcher.Apply(broadcast)
		}); allocs != 0 {
			t.Errorf("Apply allocated %v times, want 0", allocs)
		}
	})
}

// TestTheGrowthDashboardStepsAsBefore drives the growth example the new way —
// inline execution, actions through the dispatcher — and the old way — the
// default strategy, actions written with Params.Set as ApplyActionState did —
// with the same slider moves, and checks every step's state is identical.
func TestTheGrowthDashboardStepsAsBefore(t *testing.T) {
	actions := map[int][]float64{0: {0.05, 500}, 40: {0.12, 800}, 90: {0.01, 50}, 150: {0.2, 1000}}
	run := func(newPath bool) [][]float64 {
		cfg := growth.NewConfig()
		settings, implementations := cfg.SimulationGenerator().GenerateConfigs()
		implementations.OutputCondition = &simulator.NilOutputCondition{}
		implementations.OutputFunction = &simulator.NilOutputFunction{}
		if newPath {
			implementations.ExecutionStrategy = &simulator.InlineExecution{}
		}
		coordinator := simulator.NewPartitionCoordinator(settings, implementations)
		dispatcher, err := simio.NewActionDispatcher(coordinator, cfg.ActionStatePartitionNames)
		if err != nil {
			t.Fatal(err)
		}
		stepper := coordinator.NewStepper()
		defer stepper.Close()
		rows := [][]float64{}
		for step := 0; step < 200; step++ {
			if values, ok := actions[step]; ok {
				state := &simio.ActionState{Partitions: map[string]*simio.ActionValues{
					"population": {Values: values}}}
				if newPath {
					if err := dispatcher.Apply(state); err != nil {
						t.Fatal(err)
					}
				} else {
					coordinator.Iterators[0].Params.Set(simio.ActionParam,
						state.Partitions["population"].GetValues())
				}
			}
			stepper.Step()
			rows = append(rows, append([]float64(nil), coordinator.Shared.StateHistories[0].Values.RawRowView(0)...))
		}
		return rows
	}
	newRows, oldRows := run(true), run(false)
	if fmt.Sprint(newRows) != fmt.Sprint(oldRows) {
		t.Fatalf("the new path diverges from the old: step 50 %v vs %v", newRows[50], oldRows[50])
	}
	if newRows[39][0] == newRows[199][0] {
		t.Error("positive control: the slider moves should change the population")
	}
}
