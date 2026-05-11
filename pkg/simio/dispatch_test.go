package simio_test

import (
	"testing"

	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// noopIteration is a stochadex Iteration that just returns the current
// state unchanged. Used here only as a vehicle for ApplyActionState to
// have something to write `action_state_values` onto.
type noopIteration struct{}

func (*noopIteration) Configure(int, *simulator.Settings) {}

func (*noopIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	return stateHistories[partitionIndex].CopyStateRow(0)
}

// buildCoordinator builds a self-contained two-partition coordinator with
// action partitions named "alpha" and "beta". Returns it together with the
// slice and lookup-map of action-partition indices that ApplyActionState
// expects. Kept small and dependency-free so the dispatch tests don't ride
// on any specific example simulation.
func buildCoordinator(t *testing.T) (
	*simulator.PartitionCoordinator,
	[]int,
	map[string]int,
) {
	t.Helper()
	gen := simulator.NewConfigGenerator()
	actionNames := []string{"alpha", "beta"}
	for i, name := range actionNames {
		gen.SetPartition(&simulator.PartitionConfig{
			Name:      name,
			Iteration: &noopIteration{},
			Params: simulator.NewParams(map[string][]float64{
				"action_state_values": {0.0},
			}),
			InitStateValues:   []float64{0.0},
			StateHistoryDepth: 1,
			Seed:              uint64(101 + i),
		})
	}
	gen.SetSimulation(&simulator.SimulationConfig{
		OutputCondition:      &simulator.NilOutputCondition{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{MaxTimeElapsed: 1.0},
		TimestepFunction:     &simulator.ConstantTimestepFunction{Stepsize: 1.0},
		InitTimeValue:        0.0,
	})

	settings, implementations := gen.GenerateConfigs()
	implementations.OutputCondition = &simulator.NilOutputCondition{}
	implementations.OutputFunction = &simulator.NilOutputFunction{}

	indices := make([]int, 0, len(actionNames))
	byName := make(map[string]int, len(actionNames))
	for _, name := range actionNames {
		for i, it := range settings.Iterations {
			if it.Name == name {
				indices = append(indices, i)
				byName[name] = i
			}
		}
	}
	return simulator.NewPartitionCoordinator(settings, implementations), indices, byName
}

func TestApplyActionState_BroadcastPath(t *testing.T) {
	coord, indices, byName := buildCoordinator(t)

	values := []float64{1.0}
	simio.ApplyActionState(coord, indices, byName, &simio.ActionState{Values: values})

	for _, idx := range indices {
		got := coord.Iterators[idx].Params.Get("action_state_values")
		if len(got) != 1 || got[0] != 1.0 {
			t.Errorf("partition %d: expected broadcast values [1.0], got %v", idx, got)
		}
	}
}

func TestApplyActionState_NamedPath(t *testing.T) {
	coord, indices, byName := buildCoordinator(t)

	// Distinct values per named partition to verify they don't collide.
	state := &simio.ActionState{
		Partitions: map[string]*simio.ActionValues{
			"alpha":       {Values: []float64{1.0}},
			"beta":        {Values: []float64{0.0}},
			"nonexistent": {Values: []float64{42.0}}, // must be silently skipped
		},
	}
	simio.ApplyActionState(coord, indices, byName, state)

	if got := coord.Iterators[byName["alpha"]].Params.Get("action_state_values"); len(got) != 1 || got[0] != 1.0 {
		t.Errorf("alpha: expected [1.0], got %v", got)
	}
	if got := coord.Iterators[byName["beta"]].Params.Get("action_state_values"); len(got) != 1 || got[0] != 0.0 {
		t.Errorf("beta: expected [0.0], got %v", got)
	}
}

func TestApplyActionState_NamedPathTakesPrecedenceOverValues(t *testing.T) {
	coord, indices, byName := buildCoordinator(t)

	// Both Values and Partitions present. Partitions must win for matched
	// names; partitions not in the named map must NOT receive the
	// broadcast Values (precedence semantics, not a fallback per-partition).
	state := &simio.ActionState{
		Values: []float64{99.0},
		Partitions: map[string]*simio.ActionValues{
			"alpha": {Values: []float64{1.0}},
		},
	}
	simio.ApplyActionState(coord, indices, byName, state)

	if got := coord.Iterators[byName["alpha"]].Params.Get("action_state_values"); len(got) != 1 || got[0] != 1.0 {
		t.Errorf("alpha: expected [1.0] from named path, got %v", got)
	}
	if got := coord.Iterators[byName["beta"]].Params.Get("action_state_values"); len(got) == 1 && got[0] == 99.0 {
		t.Errorf("beta: legacy broadcast value leaked through named-path call (got %v)", got)
	}
}

func TestApplyActionState_NilIsNoop(t *testing.T) {
	coord, indices, byName := buildCoordinator(t)
	// Should not panic; no observable effect required.
	simio.ApplyActionState(coord, indices, byName, nil)
}
