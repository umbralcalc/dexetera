package examples

import (
	"testing"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func TestFlounceball(t *testing.T) {
	t.Run(
		"test that the Flounceball match runs",
		func(t *testing.T) {
			settings := simulator.LoadSettingsFromYaml(
				"flounceball_settings.yaml",
			)
			partitions := make([]simulator.Partition, 0)
			rateIteration := &FlounceballPlayerStateIteration{}
			rateIteration.Configure(0, settings)
			partitions = append(partitions, simulator.Partition{Iteration: rateIteration})
			coxIteration := &FlounceballMatchStateIteration{}
			coxIteration.Configure(1, settings)
			partitions = append(
				partitions,
				simulator.Partition{
					Iteration:                   coxIteration,
					ParamsFromUpstreamPartition: map[string]int{"rates": 0},
				},
			)
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
