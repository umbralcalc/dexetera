package examples

import (
	"testing"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func TestHyperspacetc(t *testing.T) {
	t.Run(
		"test that the Hyperspace Traffic Control sim runs",
		func(t *testing.T) {
			settings := simulator.LoadSettingsFromYaml(
				"hyperspacetc_settings.yaml",
			)
			partitions := []simulator.Partition{
				{
					Iteration: &SpacecraftLineCountIteration{},
				},
				{
					Iteration: &SpacecraftLineCountIteration{},
				},
				{
					Iteration: &SpacecraftLineCountIteration{},
				},
				{
					Iteration: &SpacecraftLineConnectorIteration{},
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
					Iteration: &SpacecraftLineCountIteration{},
				},
				{
					Iteration: &SpacecraftLineCountIteration{},
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
			coordinator := simulator.NewPartitionCoordinator(
				settings,
				implementations,
			)
			coordinator.Run()
		},
	)
}
