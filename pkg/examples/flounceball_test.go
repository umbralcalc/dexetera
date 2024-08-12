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
			for i := 0; i < 10; i++ {
				yourPlayerIteration := &FlounceballPlayerStateIteration{}
				partitions = append(
					partitions,
					simulator.Partition{Iteration: yourPlayerIteration},
				)
			}
			for i := 0; i < 10; i++ {
				otherPlayerIteration := &FlounceballPlayerStateIteration{}
				partitions = append(
					partitions,
					simulator.Partition{Iteration: otherPlayerIteration},
				)
			}
			matchIteration := &FlounceballMatchStateIteration{}
			partitions = append(
				partitions,
				simulator.Partition{
					Iteration: matchIteration,
					ParamsFromUpstreamPartition: map[string]int{
						"your_player_1_state":   0,
						"your_player_2_state":   1,
						"your_player_3_state":   2,
						"your_player_4_state":   3,
						"your_player_5_state":   4,
						"your_player_6_state":   5,
						"your_player_7_state":   6,
						"your_player_8_state":   7,
						"your_player_9_state":   8,
						"your_player_10_state":  9,
						"other_player_1_state":  10,
						"other_player_2_state":  11,
						"other_player_3_state":  12,
						"other_player_4_state":  13,
						"other_player_5_state":  14,
						"other_player_6_state":  15,
						"other_player_7_state":  16,
						"other_player_8_state":  17,
						"other_player_9_state":  18,
						"other_player_10_state": 19,
					},
				},
			)
			for i, partition := range partitions {
				partition.Iteration.Configure(i, settings)
			}
			implementations := &simulator.Implementations{
				Partitions:      partitions,
				OutputCondition: &simulator.EveryStepOutputCondition{},
				OutputFunction:  &simulator.NilOutputFunction{},
				TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
					MaxNumberOfSteps: 100,
				},
				TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 0.05},
			}
			coordinator := simulator.NewPartitionCoordinator(
				settings,
				implementations,
			)
			coordinator.Run()
		},
	)
}
