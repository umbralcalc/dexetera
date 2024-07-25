package examples

import (
	"math"
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
	"Possession State":                      0,
	"Your Team Total Air Time":              1,
	"Other Team Total Air Time":             2,
	"Ball Possession Air Time":              3,
	"Ball Speed State":                      4,
	"Ball Radial Position State":            5,
	"Ball Angular Position State":           6,
	"Ball Projected Radial Position State":  7,
	"Ball Projected Angular Position State": 8,
}

// PlayerStateValueIndices is a mapping which helps with describing the
// meaning of the values for each player state index.
var PlayerStateValueIndices = map[string]int{
	"Radial Position State":  0,
	"Angular Position State": 1,
	"Ball Interaction Value": 2,
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
	// TODO: Logic for attacking player attempted trajectory choice and noise on this
	// TODO: Logic for team positioning tactics when in possession and not in possession
	// TODO: ...which can depend on the ball location, projected ball location and other players
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
	ballRadial := getMatchState("Ball Radial Position State")
	ballAngle := getMatchState("Ball Angular Position State")
	ballProjRadial := getMatchState("Ball Projected Radial Position State")
	ballProjAngle := getMatchState("Ball Projected Angular Position State")

	// TODO: Classify all attackers and defenders near the ball
	// TODO: Update the ball speed state to whatever the nearest attacking players'
	// interaction value is from parameters
	for i := 1; i < 11; i++ {
		radiusAngle := params.FloatParams["your_player_"+strconv.Itoa(i)+"_radius_angle"]
		setMatchState("Ball Speed State", 0.0)
	}
	ballSpeed := getMatchState("Ball Speed State")

	// dx = r_p * cos Q_p - r * cos Q
	// dy = r_p * sin Q_p - r * sin Q
	// dx^2 + dy^2 = r_p^2 + r^2 - 2 * (r_p * r) * cos (Q_p - Q)
	// x' = x + |v| * dt * dx / sqrt( dx^2 + dy^2 )
	// y' = y + |v| * dt * dy / sqrt( dx^2 + dy^2 )
	// r' = sqrt( (x')^2 + (y')^2 )
	// Q' = arctan( y' / x' )
	// Above trig is used to compute the next ball radius and angle
	norm := math.Sqrt((ballProjRadial * ballProjRadial) + (ballRadial * ballRadial) -
		(2.0 * (ballProjRadial * ballRadial) * math.Cos(ballProjAngle-ballAngle)))
	newX := (ballRadial * math.Cos(ballAngle)) + (ballSpeed * timestepsHistory.NextIncrement *
		((ballProjRadial * math.Cos(ballProjAngle)) - (ballRadial * math.Cos(ballAngle))) / norm)
	newY := (ballRadial * math.Sin(ballAngle)) + (ballSpeed * timestepsHistory.NextIncrement *
		((ballProjRadial * math.Sin(ballProjAngle)) - (ballRadial * math.Sin(ballAngle))) / norm)
	setMatchState("Ball Radial Position State", math.Sqrt((newX*newX)+(newY*newY)))
	setMatchState("Ball Angular Position State", math.Atan(newY/newX))

	// TODO: Logic for possession and total air time updates when ball goes out of play or hits ground
	// TODO: Logic for posession air time updates when ball is in play

	return make([]float64, 0)
}
