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
			partitions := make([]simulator.Partition, 0)
			rateIteration := &SpacecraftFollowingLaneIteration{}
			rateIteration.Configure(0, settings)
			partitions = append(partitions, simulator.Partition{Iteration: rateIteration})
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
