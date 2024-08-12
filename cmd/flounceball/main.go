//go:build js && wasm

package main

import (
	"math"
	"strconv"

	"golang.org/x/exp/rand"

	"github.com/umbralcalc/dexetera/pkg/examples"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func main() {
	seeds := make([]uint64, 0)
	stateWidths := make([]int, 0)
	stateHistoryDepths := make([]int, 0)
	initStateValues := make([][]float64, 0)
	params := make([]simulator.Params, 0)
	partitions := make([]simulator.Partition, 0)
	matchStateParamsFromUpstreamPartition := make(map[string]int, 0)
	seeds = append(seeds, 0)
	stateWidths = append(stateWidths, 3)
	stateHistoryDepths = append(stateHistoryDepths, 1)
	initStateValues = append(initStateValues, []float64{0.0, 0.0, 0.0})
	params = append(params, simulator.Params{"param_values": {0.0, 0.0, 0.0}})
	partitions = append(
		partitions,
		simulator.Partition{Iteration: &simulator.ParamValuesIteration{}},
	)
	yourPlayerTemplateParams := simulator.Params{
		"match_state_partition_index":         {21},
		"opposition_player_partition_indices": {11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		"player_space_finding_talent":         {7},
		"player_ball_interaction_speed":       {0.5},
		"player_ball_interaction_inaccuracy":  {0.1},
		"team_possession_state_value":         {0},
		"team_attacking_distance_threshold":   {10.0},
		"team_defensive_distance_threshold":   {10.0},
		"player_movement_speed":               {0.1},
	}
	index := len(partitions)
	radius := 0.0
	angle := 0.0
	for i := 1; i < 11; i++ {
		seeds = append(seeds, uint64(rand.Intn(10000)))
		stateWidths = append(stateWidths, 4)
		stateHistoryDepths = append(stateHistoryDepths, 1)
		initStateValues = append(initStateValues, []float64{radius, angle, 0.0, 0.0})
		if i%2 == 0 {
			radius += examples.PitchRadiusMetres / 5.0
		}
		angle += math.Pi
		copyParams := yourPlayerTemplateParams
		params = append(params, copyParams)
		yourPlayerIteration := &examples.FlounceballPlayerStateIteration{}
		partitions = append(
			partitions,
			simulator.Partition{Iteration: yourPlayerIteration},
		)
		matchStateParamsFromUpstreamPartition["your_player_"+
			strconv.Itoa(i)+"_state"] = index
		index += 1
	}
	otherPlayerTemplateParams := simulator.Params{
		"match_state_partition_index":         {21},
		"opposition_player_partition_indices": {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"player_space_finding_talent":         {7},
		"player_ball_interaction_speed":       {0.5},
		"player_ball_interaction_inaccuracy":  {0.1},
		"team_possession_state_value":         {1},
		"team_attacking_distance_threshold":   {10.0},
		"team_defensive_distance_threshold":   {10.0},
		"player_movement_speed":               {0.1},
	}
	radius = 0.0
	angle = math.Pi
	for i := 1; i < 11; i++ {
		seeds = append(seeds, uint64(rand.Intn(10000)))
		stateWidths = append(stateWidths, 4)
		stateHistoryDepths = append(stateHistoryDepths, 1)
		initStateValues = append(initStateValues, []float64{radius, angle, 0.0, 0.0})
		if i%2 == 0 {
			radius += examples.PitchRadiusMetres / 5.0
		}
		angle += math.Pi
		copyParams := otherPlayerTemplateParams
		params = append(params, copyParams)
		otherPlayerIteration := &examples.FlounceballPlayerStateIteration{}
		partitions = append(
			partitions,
			simulator.Partition{Iteration: otherPlayerIteration},
		)
		matchStateParamsFromUpstreamPartition["other_player_"+
			strconv.Itoa(i)+"_state"] = index
		index += 1
	}
	seeds = append(seeds, uint64(rand.Intn(10000)))
	stateWidths = append(stateWidths, 11)
	stateHistoryDepths = append(stateHistoryDepths, 1)
	initStateValues = append(
		initStateValues,
		[]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
	)
	matchTemplateParams := simulator.Params{
		"max_ball_falling_time": {5.0},
	}
	params = append(params, matchTemplateParams)
	matchIteration := &examples.FlounceballMatchStateIteration{}
	partitions = append(
		partitions,
		simulator.Partition{
			Iteration:                   matchIteration,
			ParamsFromUpstreamPartition: matchStateParamsFromUpstreamPartition,
		},
	)
	settings := &simulator.Settings{
		Params:                params,
		InitStateValues:       initStateValues,
		InitTimeValue:         0.0,
		Seeds:                 seeds,
		StateWidths:           stateWidths,
		StateHistoryDepths:    stateHistoryDepths,
		TimestepsHistoryDepth: 1,
	}
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
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
