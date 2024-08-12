//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/examples"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func main() {
	settings := &simulator.Settings{
		Params: []simulator.Params{},
		InitStateValues: [][]float64{
			{1.0, 1.0, 1.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0},
		},
		InitTimeValue:         0.0,
		Seeds:                 []uint64{0, 563, 8312},
		StateWidths:           []int{3, 5, 3},
		StateHistoryDepths:    []int{2, 2, 2},
		TimestepsHistoryDepth: 2,
	}
	partitions := []simulator.Partition{
		{
			Iteration: &examples.SpacecraftLaneCountIteration{},
		},
		{
			Iteration: &examples.SpacecraftLaneCountIteration{},
		},
		{
			Iteration: &examples.SpacecraftLaneCountIteration{},
		},
		{
			Iteration: &examples.SpacecraftLaneConnectorIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_0_input_count": 0,
				"partition_1_input_count": 1,
				"partition_2_input_count": 2,
			},
			ParamsFromIndices: map[string][]int{
				"partition_0_input_count": {1},
				"partition_1_input_count": {1},
				"partition_2_input_count": {1},
			},
		},
		{
			Iteration: &examples.SpacecraftLaneCountIteration{},
		},
		{
			Iteration: &examples.SpacecraftLaneCountIteration{},
		},
	}
	for index, partition := range partitions {
		partition.Iteration.Configure(index, settings)
	}
	implementations := &simulator.Implementations{
		Partitions:      partitions,
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 100,
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 1.0},
	}
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
