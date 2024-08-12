//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/examples"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func main() {
	yourPlayerTemplateParams := simulator.Params{
		"match_state_partition_index":         {20},
		"opposition_player_partition_indices": {10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		"player_space_finding_talent":         {7},
		"player_ball_interaction_speed":       {0.5},
		"player_ball_interaction_inaccuracy":  {0.1},
		"team_possession_state_value":         {0},
		"team_attacking_distance_threshold":   {10.0},
		"team_defensive_distance_threshold":   {10.0},
		"player_movement_speed":               {0.1},
	}
	otherPlayerTemplateParams := simulator.Params{
		"match_state_partition_index":         {20},
		"opposition_player_partition_indices": {0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		"player_space_finding_talent":         {7},
		"player_ball_interaction_speed":       {0.5},
		"player_ball_interaction_inaccuracy":  {0.1},
		"team_possession_state_value":         {1},
		"team_attacking_distance_threshold":   {10.0},
		"team_defensive_distance_threshold":   {10.0},
		"player_movement_speed":               {0.1},
	}
	matchTemplateParams := simulator.Params{
		"max_ball_falling_time": {5.0},
	}
	settings := &simulator.Settings{
		Params: []simulator.Params{
			yourPlayerTemplateParams,
			yourPlayerTemplateParams,
			yourPlayerTemplateParams,
			yourPlayerTemplateParams,
			yourPlayerTemplateParams,
			yourPlayerTemplateParams,
			yourPlayerTemplateParams,
			yourPlayerTemplateParams,
			yourPlayerTemplateParams,
			yourPlayerTemplateParams,
			otherPlayerTemplateParams,
			otherPlayerTemplateParams,
			otherPlayerTemplateParams,
			otherPlayerTemplateParams,
			otherPlayerTemplateParams,
			otherPlayerTemplateParams,
			otherPlayerTemplateParams,
			otherPlayerTemplateParams,
			otherPlayerTemplateParams,
			otherPlayerTemplateParams,
			matchTemplateParams,
		},
		InitStateValues: [][]float64{
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		},
		InitTimeValue:         0.0,
		Seeds:                 []uint64{563, 8312, 111, 24253, 55524, 63, 12, 1, 2253, 524, 1563, 822312, 11211, 23, 24, 6, 2, 1000, 3, 4, 8898},
		StateWidths:           []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 11},
		StateHistoryDepths:    []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		TimestepsHistoryDepth: 1,
	}
	partitions := make([]simulator.Partition, 0)
	for i := 0; i < 10; i++ {
		yourPlayerIteration := &examples.FlounceballPlayerStateIteration{}
		partitions = append(
			partitions,
			simulator.Partition{Iteration: yourPlayerIteration},
		)
	}
	for i := 0; i < 10; i++ {
		otherPlayerIteration := &examples.FlounceballPlayerStateIteration{}
		partitions = append(
			partitions,
			simulator.Partition{Iteration: otherPlayerIteration},
		)
	}
	matchIteration := &examples.FlounceballMatchStateIteration{}
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
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
