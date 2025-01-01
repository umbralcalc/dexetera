package games

import (
	"testing"

	"github.com/umbralcalc/stochadex/pkg/continuous"
	"github.com/umbralcalc/stochadex/pkg/discrete"
	"github.com/umbralcalc/stochadex/pkg/general"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func TestHyperspacetc(t *testing.T) {
	t.Run(
		"test that the Hyperspace Traffic Control sim runs",
		func(t *testing.T) {
			settings := simulator.LoadSettingsFromYaml(
				"hyperspacetc_settings.yaml",
			)
			sourceIterationByEvent := map[float64]simulator.Iteration{
				0: &general.ParamValuesIteration{},
				1: &continuous.CumulativeTimeIteration{},
			}
			downstreamIterationByEvent := map[float64]simulator.Iteration{
				0: &general.ConstantValuesIteration{},
				1: &general.ValuesCollectionPushIteration{
					PushFunction: EntryTimeFromUpstreamPushFunction,
				},
				2: &general.ValuesCollectionPopIteration{
					PopIndexFunction: general.NextNonEmptyPopIndexFunction,
				},
				3: &general.SerialIterationsIteration{
					Iterations: []simulator.Iteration{
						&general.ValuesCollectionPushIteration{
							PushFunction: EntryTimeFromUpstreamPushFunction,
						},
						&general.ValuesCollectionPopIteration{
							PopIndexFunction: general.NextNonEmptyPopIndexFunction,
						},
					},
				},
			}
			partitions := []simulator.Partition{
				{
					Iteration: &general.ValuesGroupedAggregationIteration{
						ValuesFunction: SpacecraftQueueValuesFunction,
						AggFunction:    general.CountAggFunction,
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration:   &discrete.BernoulliProcessIteration{},
						IterationByEvent: sourceIterationByEvent,
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration:   &SpacecraftLineEventIteration{},
						IterationByEvent: downstreamIterationByEvent,
					},
					ParamsFromUpstreamPartition: map[string]int{
						"queue_size": 0,
					},
					ParamsFromIndices: map[string][]int{
						"queue_size": {0},
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration: &general.ValuesFunctionIteration{
							Function: SpacecraftQueueEventFunction,
						},
						IterationByEvent: downstreamIterationByEvent,
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration:   &discrete.BernoulliProcessIteration{},
						IterationByEvent: sourceIterationByEvent,
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration:   &SpacecraftLineEventIteration{},
						IterationByEvent: downstreamIterationByEvent,
					},
					ParamsFromUpstreamPartition: map[string]int{
						"queue_size": 0,
					},
					ParamsFromIndices: map[string][]int{
						"queue_size": {1},
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration: &general.ValuesFunctionIteration{
							Function: SpacecraftQueueEventFunction,
						},
						IterationByEvent: downstreamIterationByEvent,
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration:   &discrete.BernoulliProcessIteration{},
						IterationByEvent: sourceIterationByEvent,
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration:   &SpacecraftLineEventIteration{},
						IterationByEvent: downstreamIterationByEvent,
					},
					ParamsFromUpstreamPartition: map[string]int{
						"queue_size": 0,
					},
					ParamsFromIndices: map[string][]int{
						"queue_size": {2},
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration: &general.ValuesFunctionIteration{
							Function: SpacecraftQueueEventFunction,
						},
						IterationByEvent: downstreamIterationByEvent,
					},
				},
				{
					Iteration: &SpacecraftLineConnectorIteration{},
					ParamsFromUpstreamPartition: map[string]int{
						"partition_3_input_value": 3,
						"partition_6_input_value": 6,
						"partition_9_input_value": 9,
					},
					ParamsFromIndices: map[string][]int{
						"partition_3_input_value": {0},
						"partition_6_input_value": {0},
						"partition_9_input_value": {0},
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration:   &SpacecraftLineEventIteration{},
						IterationByEvent: downstreamIterationByEvent,
					},
				},
				{
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration:   &SpacecraftLineEventIteration{},
						IterationByEvent: downstreamIterationByEvent,
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
					MaxNumberOfSteps: 5,
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
