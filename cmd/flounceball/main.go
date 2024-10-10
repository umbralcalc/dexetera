//go:build js && wasm

package main

import (
	"math"
	"strconv"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"

	"github.com/umbralcalc/dexetera/pkg/examples"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/general"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

type MatchPartitionConfig struct {
	Partitions                            []simulator.Partition
	TeamMinParamsFromUpstreamPartition    map[string]int
	MatchStateParamsFromUpstreamPartition map[string]int
}

func initEmptySettings() *simulator.Settings {
	return &simulator.Settings{
		Params:                make([]simulator.Params, 0),
		InitStateValues:       make([][]float64, 0),
		InitTimeValue:         0.0,
		Seeds:                 make([]uint64, 0),
		StateWidths:           make([]int, 0),
		StateHistoryDepths:    make([]int, 0),
		TimestepsHistoryDepth: 1,
	}
}

func addActionTaker(
	settings *simulator.Settings,
	partitionConfig *MatchPartitionConfig,
) {
	settings.Seeds = append(settings.Seeds, 0)
	settings.StateWidths = append(settings.StateWidths, 20)
	settings.StateHistoryDepths = append(settings.StateHistoryDepths, 1)
	uniformDist := &distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(uint64(rand.Intn(10000))),
	}
	initialActionValues := make([]float64, 0)
	for i := 1; i < 11; i++ {
		initialActionValues = append(
			initialActionValues,
			examples.PitchRadiusMetres*uniformDist.Rand(),
		)
		initialActionValues = append(
			initialActionValues,
			2.0*math.Pi*uniformDist.Rand(),
		)
	}
	settings.InitStateValues = append(
		settings.InitStateValues,
		initialActionValues,
	)
	paramValues := initialActionValues
	settings.Params = append(
		settings.Params,
		simulator.Params{"param_values": paramValues},
	)
	partitionConfig.Partitions = append(
		partitionConfig.Partitions,
		simulator.Partition{Iteration: &general.ParamValuesIteration{}},
	)
}

func addYourPlayers(
	settings *simulator.Settings,
	partitionConfig *MatchPartitionConfig,
) {
	index := len(partitionConfig.Partitions)
	uniformDist := &distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(uint64(rand.Intn(10000))),
	}
	for i := 1; i < 11; i++ {
		settings.Seeds = append(settings.Seeds, uint64(rand.Intn(10000)))
		settings.StateWidths = append(settings.StateWidths, 2)
		settings.StateHistoryDepths = append(settings.StateHistoryDepths, 1)
		settings.InitStateValues = append(
			settings.InitStateValues,
			[]float64{
				examples.PitchRadiusMetres * uniformDist.Rand(),
				2.0 * math.Pi * uniformDist.Rand(),
			},
		)
		settings.Params = append(settings.Params, simulator.Params{})
		yourPlayerIteration := &examples.FlounceballPlayerStateIteration{}
		partitionConfig.Partitions = append(
			partitionConfig.Partitions,
			simulator.Partition{
				Iteration: yourPlayerIteration,
				ParamsFromUpstreamPartition: map[string]int{
					"manager_directed_coordinates": 0,
				},
				ParamsFromIndices: map[string][]int{
					"manager_directed_coordinates": {2 * (i - 1), 2*(i-1) + 1},
				},
			},
		)
		partitionConfig.TeamMinParamsFromUpstreamPartition["your_player_"+
			strconv.Itoa(i)+"_state"] = index
		index += 1
	}
}

func addOtherPlayers(
	settings *simulator.Settings,
	partitionConfig *MatchPartitionConfig,
) {
	index := len(partitionConfig.Partitions)
	uniformDist := &distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(uint64(rand.Intn(10000))),
	}
	for i := 1; i < 11; i++ {
		settings.Seeds = append(settings.Seeds, uint64(rand.Intn(10000)))
		settings.StateWidths = append(settings.StateWidths, 2)
		settings.StateHistoryDepths = append(settings.StateHistoryDepths, 1)
		settings.InitStateValues = append(
			settings.InitStateValues,
			[]float64{
				examples.PitchRadiusMetres * uniformDist.Rand(),
				2.0 * math.Pi * uniformDist.Rand(),
			},
		)
		settings.Params = append(settings.Params, simulator.Params{})
		otherPlayerIteration := &examples.FlounceballPlayerStateIteration{}
		partitionConfig.Partitions = append(
			partitionConfig.Partitions,
			simulator.Partition{Iteration: otherPlayerIteration},
		)
		partitionConfig.TeamMinParamsFromUpstreamPartition["other_player_"+
			strconv.Itoa(i)+"_state"] = index
		index += 1
	}
}

func addChosenCoordsGenerator(
	settings *simulator.Settings,
	partitionConfig *MatchPartitionConfig,
) {
	index := len(partitionConfig.Partitions)
	settings.Seeds = append(settings.Seeds, 0)
	settings.StateWidths = append(settings.StateWidths, 2)
	settings.StateHistoryDepths = append(settings.StateHistoryDepths, 1)
	settings.InitStateValues = append(
		settings.InitStateValues,
		[]float64{0.0, 0.0},
	)
	settings.Params = append(
		settings.Params,
		simulator.Params{"param_values": {0.0, 0.0}},
	)
	partitionConfig.Partitions = append(
		partitionConfig.Partitions,
		simulator.Partition{Iteration: &general.ParamValuesIteration{}},
	)
	partitionConfig.TeamMinParamsFromUpstreamPartition["chosen_coordinates"] = index
}

func addTeamMinCalculator(
	settings *simulator.Settings,
	partitionConfig *MatchPartitionConfig,
) {
	index := len(partitionConfig.Partitions)
	settings.Seeds = append(settings.Seeds, 0)
	settings.StateWidths = append(settings.StateWidths, 2)
	settings.StateHistoryDepths = append(settings.StateHistoryDepths, 1)
	settings.InitStateValues = append(settings.InitStateValues, []float64{1000.0, 1000.0})
	settings.Params = append(
		settings.Params,
		simulator.Params{
			"accepted_value_groups": {0, 1},
			"default_values":        {1000.0, 1000.0},
		},
	)
	minGroupingIteration := &general.ValuesGroupedAggregationIteration{
		ValuesFunction: examples.PlayerProximityValuesFunction,
		AggFunction:    general.MinAggFunction,
	}
	partitionConfig.Partitions = append(
		partitionConfig.Partitions,
		simulator.Partition{
			Iteration:                   minGroupingIteration,
			ParamsFromUpstreamPartition: partitionConfig.TeamMinParamsFromUpstreamPartition,
		},
	)
	partitionConfig.MatchStateParamsFromUpstreamPartition["team_proximity_minima"] = index
}

func addMatchState(
	settings *simulator.Settings,
	partitionConfig *MatchPartitionConfig,
) {
	settings.Seeds = append(settings.Seeds, uint64(rand.Intn(10000)))
	settings.StateWidths = append(settings.StateWidths, 2)
	settings.StateHistoryDepths = append(settings.StateHistoryDepths, 1)
	settings.InitStateValues = append(settings.InitStateValues, []float64{0.0, 0.0})
	settings.Params = append(settings.Params, simulator.Params{})
	matchIteration := &general.ValuesFunctionIteration{
		Function: examples.FlounceballMatchStateValuesFunction,
	}
	partitionConfig.Partitions = append(
		partitionConfig.Partitions,
		simulator.Partition{
			Iteration:                   matchIteration,
			ParamsFromUpstreamPartition: partitionConfig.MatchStateParamsFromUpstreamPartition,
		},
	)
}

func main() {
	totalMatchSeconds := 300.0
	timeStepsizeSeconds := 0.05
	settings := initEmptySettings()
	partitionConfig := &MatchPartitionConfig{
		Partitions:                            make([]simulator.Partition, 0),
		MatchStateParamsFromUpstreamPartition: make(map[string]int, 0),
		TeamMinParamsFromUpstreamPartition:    make(map[string]int, 0),
	}
	addActionTaker(settings, partitionConfig)
	addYourPlayers(settings, partitionConfig)
	addOtherPlayers(settings, partitionConfig)
	addChosenCoordsGenerator(settings, partitionConfig)
	addTeamMinCalculator(settings, partitionConfig)
	addMatchState(settings, partitionConfig)
	for i, partition := range partitionConfig.Partitions {
		partition.Iteration.Configure(i, settings)
	}
	implementations := &simulator.Implementations{
		Partitions:      partitionConfig.Partitions,
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.StdoutOutputFunction{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{
			MaxTimeElapsed: totalMatchSeconds,
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: timeStepsizeSeconds},
	}
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
