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
			// queue-left-node-middle-outside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
		},
		{
			// queue-left-node-upper-outside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
		},
		{
			// queue-left-node-lower-outside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
		},
		{
			// left-node
			Iteration: &examples.SpacecraftLineConnectorIteration{},
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
			// queue-upper-node-outside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
		},
		{
			// queue-upper-node-inside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
		},
		{
			// upper-node
			Iteration: &examples.SpacecraftLineConnectorIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_4_input_count": 4,
				"partition_5_input_count": 5,
			},
			ParamsFromIndices: map[string][]int{
				"partition_4_input_count": {1},
				"partition_5_input_count": {1},
			},
		},
		{
			// queue-right-node-lower-inside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
		},
		{
			// queue-right-node-upper-inside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
		},
		{
			// right-node
			Iteration: &examples.SpacecraftLineConnectorIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_7_input_count": 7,
				"partition_8_input_count": 8,
			},
			ParamsFromIndices: map[string][]int{
				"partition_7_input_count": {1},
				"partition_8_input_count": {1},
			},
		},
	}
	for index, partition := range partitions {
		partition.Iteration.Configure(index, settings)
	}
	implementations := &simulator.Implementations{
		Partitions:      partitions,
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{
			MaxTimeElapsed: 100,
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 1.0},
	}
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
