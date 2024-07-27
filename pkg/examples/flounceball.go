package examples

import (
	"math"
	"strconv"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// PitchRadiusMetres is the radius of the circular pitch.
const PitchRadiusMetres = 100.0

// InteractionRadiusMetres is the maximum radius between two entities
// allowed for them to interact.
const InteractionRadiusMetres = 1.0

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

// OutputState is a simple container for referencing states on output.
type OutputState struct {
	Values []float64
}

// Coordinates is a convenient struct to hold coordinates on the field
// and operates on them.
type Coordinates struct {
	Radial  float64
	Angular float64
}

// Update updates the radial and angular coordinates of the entity
// given the current and projected coordinates as well as its speed.
func (c *Coordinates) Update(
	projCoords *Coordinates,
	speed float64,
	timestep float64,
) {
	// dx = r_p * cos Q_p - r * cos Q
	// dy = r_p * sin Q_p - r * sin Q
	// dx^2 + dy^2 = r_p^2 + r^2 - 2 * (r_p * r) * cos (Q_p - Q)
	// x' = x + |v| * dt * dx / sqrt( dx^2 + dy^2 )
	// y' = y + |v| * dt * dy / sqrt( dx^2 + dy^2 )
	// r' = sqrt( (x')^2 + (y')^2 )
	// Q' = arctan( y' / x' )
	// The trig above is used to compute the next radius and angle
	projRadial := projCoords.Radial
	projAngle := projCoords.Angular
	norm := math.Sqrt((projRadial * projRadial) + (c.Radial * c.Radial) -
		(2.0 * (projRadial * c.Radial) * math.Cos(projAngle-c.Angular)))
	newX := (c.Radial * math.Cos(c.Angular)) + (speed * timestep *
		((projRadial * math.Cos(projAngle)) - (c.Radial * math.Cos(c.Angular))) / norm)
	newY := (c.Radial * math.Sin(c.Angular)) + (speed * timestep *
		((projRadial * math.Sin(projAngle)) - (c.Radial * math.Sin(c.Angular))) / norm)
	c.Radial = math.Sqrt((newX * newX) + (newY * newY))
	c.Angular = math.Atan(newY / newX)
}

// Proximity returns the proximity (in terms of absolute distance) of the
// input coordinates to the entity.
func (c *Coordinates) Proximity(otherCoords *Coordinates) float64 {
	diffX := (c.Radial * math.Cos(c.Angular)) -
		(otherCoords.Radial * math.Cos(otherCoords.Angular))
	diffY := (c.Radial * math.Sin(c.Angular)) -
		(otherCoords.Radial * math.Sin(otherCoords.Angular))
	return math.Sqrt((diffX * diffX) + (diffY * diffY))
}

// PossessionNameMatcher helps to ensure that the right team names are used
// when referring to data about attack or defence.
type PossessionNameMatcher struct {
	possession int
}

// Attacking finds the right name to refer to the attacking side, depending
// on possession.
func (p *PossessionNameMatcher) Attacking(
	yourName string,
	otherName string,
) string {
	switch PossessionValueMap[p.possession] {
	case "Your Team":
		return yourName
	case "Other Team":
		return otherName
	default:
		panic("possession value was invalid: " + strconv.Itoa(p.possession))
	}
}

// Defending finds the right name to refer to the defending side, depending
// on possession.
func (p *PossessionNameMatcher) Defending(
	yourName string,
	otherName string,
) string {
	switch PossessionValueMap[p.possession] {
	case "Your Team":
		return otherName
	case "Other Team":
		return yourName
	default:
		panic("possession value was invalid: " + strconv.Itoa(p.possession))
	}
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
	state *OutputState,
) func(key string, value float64) {
	return func(key string, value float64) {
		state.Values[PlayerStateValueIndices[key]] = value
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
	// Reorganise state data and specify getters and setters for convenience.
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
	outputState := &OutputState{Values: make([]float64, len(PlayerStateValueIndices))}
	setThisPlayerState := generatePlayerStateValueSetter(outputState)

	// TODO: Logic for attacking player disruptions from defensive players - limits accuracy
	// TODO: Logic for attacking player attempted trajectory choice and noise on this
	// TODO: Logic for team positioning tactics when in possession and not in possession
	// TODO: ...which can depend on the ball location, projected ball location and other players
	return outputState.Values
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
	state *OutputState,
) func(key string, value float64) {
	return func(key string, value float64) {
		state.Values[MatchStateValueIndices[key]] = value
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
	// Reorganise state data and specify getters and setters for convenience.
	getMatchState := generateMatchStateValueGetter(stateHistories[partitionIndex])
	outputState := &OutputState{Values: make([]float64, len(MatchStateValueIndices))}
	setMatchState := generateMatchStateValueSetter(outputState)
	ballCoords := &Coordinates{
		Radial:  getMatchState("Ball Radial Position State"),
		Angular: getMatchState("Ball Angular Position State"),
	}
	ballProjCoords := &Coordinates{
		Radial:  getMatchState("Ball Projected Radial Position State"),
		Angular: getMatchState("Ball Projected Angular Position State"),
	}
	posMatcher := &PossessionNameMatcher{
		possession: int(getMatchState("Possession State")),
	}

	// TODO: Classify all attackers and defenders near the ball
	// TODO: Update the ball speed state to whatever the nearest attacking players'
	// interaction value is from parameters
	attackPlayerCoords := &Coordinates{}
	for i := 1; i < 11; i++ {
		radiusAngle := params.FloatParams[posMatcher.Attacking(
			"your_player_"+strconv.Itoa(i)+"_radius_angle",
			"other_player_"+strconv.Itoa(i)+"_radius_angle",
		)]
		attackPlayerCoords.Radial = radiusAngle[0]
		attackPlayerCoords.Angular = radiusAngle[1]
		if ballCoords.Proximity(attackPlayerCoords) <= InteractionRadiusMetres {
			// FIXME: What happens when there is more than one???
			setMatchState("Ball Speed State", 0.0)
		}
	}
	ballSpeed := getMatchState("Ball Speed State")

	// Compute the next ball radius and angle
	ballCoords.Update(
		ballProjCoords,
		ballSpeed,
		timestepsHistory.NextIncrement,
	)
	setMatchState("Ball Radial Position State", ballCoords.Radial)
	setMatchState("Ball Angular Position State", ballCoords.Angular)

	// TODO: Logic for possession and total air time updates when ball goes out of play or hits ground
	// TODO: Logic for posession air time updates when ball is in play

	return outputState.Values
}
