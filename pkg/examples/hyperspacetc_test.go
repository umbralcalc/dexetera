package examples

import (
	"testing"

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
			iterationByEvent := map[float64]simulator.Iteration{
				0: &general.ConstantValuesIteration{},
				1: &general.ValuesCollectionPushIteration{
					PushFunction: general.ParamValuesPushFunction,
				},
				2: &general.ValuesCollectionPopIteration{
					PopIndexFunction: general.NextNonEmptyPopIndexFunction,
				},
				3: &general.SerialIterationsIteration{
					Iterations: []simulator.Iteration{
						&general.ValuesCollectionPushIteration{
							PushFunction: general.ParamValuesPushFunction,
						},
						&general.ValuesCollectionPopIteration{
							PopIndexFunction: general.NextNonEmptyPopIndexFunction,
						},
					},
				},
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
					Iteration: &general.ValuesChangingEventsIteration{
						EventIteration:   &SpacecraftLineEventIteration{},
						IterationByEvent: iterationByEvent,
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
						EventIteration:   &SpacecraftLineEventIteration{},
						IterationByEvent: iterationByEvent,
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
						EventIteration:   &SpacecraftLineEventIteration{},
						IterationByEvent: iterationByEvent,
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
						"partition_0_input_value": 0,
						"partition_2_input_value": 2,
						"partition_4_input_value": 4,
					},
					ParamsFromIndices: map[string][]int{
						"partition_0_input_value": {0},
						"partition_2_input_value": {0},
						"partition_4_input_value": {0},
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
						EventIteration: &general.ValuesFunctionIteration{
							Function: SpacecraftQueueEventFunction,
						},
						IterationByEvent: downstreamIterationByEvent,
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
						EventIteration: &general.ValuesFunctionIteration{
							Function: SpacecraftQueueEventFunction,
						},
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
