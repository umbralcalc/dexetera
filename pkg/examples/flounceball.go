package examples

import (
	"strconv"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// PitchRadiusMetres is the radius of the circular pitch.
const PitchRadiusMetres = 100.0

// PossessionValueMap is a mapping to check which team is in possession
// based on the value of the possession state index.
var PossessionValueMap = map[int]string{0: "Your Team", 1: "Other Team"}

// MatchStateValueIndices is a mapping which helps with describing the
// meaning of the values for each match state index.
var MatchStateValueIndices = map[string]int{
	"Possession State":             0,
	"Your Team Total Air Time":     1,
	"Other Team Total Air Time":    2,
	"Ball Possession Air Time":     3,
	"Ball Radial Position State":   4,
	"Ball Angular Position State":  5,
	"Ball Vertical Position State": 6,
}

// PlayerStateValueIndices is a mapping which helps with describing the
// meaning of the values for each player state index.
var PlayerStateValueIndices = map[string]int{
	"Radial Position State":  0,
	"Angular Position State": 1,
}

// generatePlayerStateValuesGetter creates a closure which reduces the
// amount of code required to retrieve state values for all players.
func generatePlayerStateValuesGetter(
	playerPartitionIndices []int64,
	stateHistories []*simulator.StateHistory,
) func(key string) []float64 {
	return func(key string) []float64 {
		values := make([]float64, 0)
		for _, index := range playerPartitionIndices {
			values = append(
				values,
				stateHistories[index].Values.At(0, PlayerStateValueIndices[key]),
			)
		}
		return values
	}
}

// generatePlayerStateValueSetter creates a closure which reduces the
// amount of code required to reassign state values for a player.
func generatePlayerStateValueSetter(
	stateHistory *simulator.StateHistory,
) func(key string, value float64) {
	return func(key string, value float64) {
		stateHistory.Values.Set(0, PlayerStateValueIndices[key], value)
	}
}

// FlounceballPlayerStateIteration describes the iteration of an individual
// player state in a Flounceball match.
type FlounceballPlayerStateIteration struct {
}

func (f *FlounceballPlayerStateIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
}

func (f *FlounceballPlayerStateIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	getMatchState := generateMatchStateValueGetter(
		stateHistories[params.IntParams["match_partition_index"][0]],
	)
	getYourPlayerStates := generatePlayerStateValuesGetter(
		params.IntParams["your_player_partition_indices"],
		stateHistories,
	)
	getOtherPlayerStates := generatePlayerStateValuesGetter(
		params.IntParams["other_player_partition_indices"],
		stateHistories,
	)
	setYourPlayerState := generatePlayerStateValueSetter(stateHistories[partitionIndex])
	// TODO: Logic for attacking player disruptions from defensive players - limits accuracy
	// TODO: Logic for attacking player attempted trajectory choice
	// TODO: Logic for team positioning tactics when in possession and not in possession
	return make([]float64, 0)
}

// generateMatchStateValueGetter creates a closure which reduces the
// amount of code required to retrieve state values.
func generateMatchStateValueGetter(
	stateHistory *simulator.StateHistory,
) func(key string) float64 {
	return func(key string) float64 {
		return stateHistory.Values.At(0, MatchStateValueIndices[key])
	}
}

// generateMatchStateValueSetter creates a closure which reduces the
// amount of code required to reassign state values.
func generateMatchStateValueSetter(
	stateHistory *simulator.StateHistory,
) func(key string, value float64) {
	return func(key string, value float64) {
		stateHistory.Values.Set(0, MatchStateValueIndices[key], value)
	}
}

// FlounceballMatchStateIteration describes the iteration of a Flounceball
// match state in response to player positions and manager decisions.
type FlounceballMatchStateIteration struct {
}

func (f *FlounceballMatchStateIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
}

func (f *FlounceballMatchStateIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	getMatchState := generateMatchStateValueGetter(stateHistories[partitionIndex])
	setMatchState := generateMatchStateValueSetter(stateHistories[partitionIndex])
	// TODO: Logic for ball trajectories - needs to hit the ground and be motionless
	// TODO: Logic for possession and total air time updates when ball goes out of play or hits ground
	// TODO: Logic for posession air time updates when ball is in play
	ballRadius := getMatchState("Ball Radial Position State")
	ballAngle := getMatchState("Ball Angular Position State")
	ballVert := getMatchState("Ball Vertical Position State")
	for i := 1; i < 11; i++ {
		radiusAngle := params.FloatParams["your_player_"+strconv.Itoa(i)+"_radius_angle"]
	}
	return make([]float64, 0)
}
