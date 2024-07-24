//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/examples"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func main() {
	settings := &simulator.Settings{
		OtherParams: []*simulator.OtherParams{
			{
				FloatParams: map[string][]float64{
					"param_values": {1.0, 1.0, 1.0},
				},
				IntParams: map[string][]int64{},
			},
			{
				FloatParams: map[string][]float64{
					"rates":        {0.5, 1.0, 0.8, 1.0, 1.1},
					"gamma_alphas": {1.0, 2.5, 3.0, 1.8, 1.0},
					"gamma_betas":  {2.0, 1.0, 4.1, 2.0, 1.2},
				},
				IntParams: map[string][]int64{},
			},
			{
				FloatParams: map[string][]float64{
					"rates":        {1.5, 0.2, 0.6},
					"gamma_alphas": {2.3, 5.1, 2.0},
					"gamma_betas":  {2.0, 1.5, 1.1},
				},
				IntParams: map[string][]int64{},
			},
		},
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
		{Iteration: &simulator.ParamValuesIteration{}},
		{Iteration: &examples.SpacecraftFollowingLaneIteration{}},
		{
			Iteration: &simulator.ConstantValuesIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"rates": 0,
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
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 100,
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 1.0},
	}
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
