package examples

import (
	"math"
	"strconv"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// PitchRadiusMetres is the radius of the circular pitch.
const PitchRadiusMetres = 100.0

// InteractionRadiusMetres is the maximum radius between two entities
// allowed for them to interact.
const InteractionRadiusMetres = 1.0

// MaxBallFallingLifetime is the maximum time a ball is assumed to be in
// still airborne when it reaches within the InteractionRadiusMetres of
// its projected location and no player has yet interacted with it.
const MaxBallFallingLifetime = 1.0

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
	"Ball Cumulative Falling Time":          9,
}

// PlayerStateValueIndices is a mapping which helps with describing the
// meaning of the values for each player state index.
var PlayerStateValueIndices = map[string]int{
	"Radial Position State":       0,
	"Angular Position State":      1,
	"Ball Interaction Speed":      2,
	"Ball Interaction Inaccuracy": 3,
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

// ApplyShift applies a cartesian coordinate shift to the radial and
// angular position.
func (c *Coordinates) ApplyShift(xDiff float64, yDiff float64) {
	newX := (c.Radial * math.Cos(c.Angular)) + xDiff
	newY := (c.Radial * math.Sin(c.Angular)) + yDiff
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

// MinProximity finds the minimum proximity value to this Coordinates struct
// among the slice of other Coordinates structs provided.
func (c *Coordinates) MinProximity(otherCoords []*Coordinates) (int, float64) {
	outputIndex := 0
	outputProx := c.Proximity(otherCoords[0])
	for i, coords := range otherCoords {
		if i == 0 {
			continue
		}
		prox := c.Proximity(coords)
		if prox < outputProx {
			outputProx = prox
			outputIndex = i
		}
	}
	return outputIndex, outputProx
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

// FlounceballPlayerStateIteration describes the iteration of an individual
// player state in a Flounceball match.
type FlounceballPlayerStateIteration struct {
	uniformDist *distuv.Uniform
}

func (f *FlounceballPlayerStateIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	f.uniformDist = &distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(settings.Seeds[partitionIndex]),
	}
}

func (f *FlounceballPlayerStateIteration) Iterate(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	// Reorganise state data
	stateHistory := stateHistories[partitionIndex]
	outputState := stateHistory.Values.RawRowView(0)
	playerCoords := &Coordinates{
		Radial: stateHistory.Values.At(
			0, PlayerStateValueIndices["Radial Position State"]),
		Angular: stateHistory.Values.At(
			0, PlayerStateValueIndices["Angular Position State"]),
	}
	oppositionPlayerCoords := make([]*Coordinates, 0)
	for _, index := range params["opposition_player_partition_indices"] {
		oppositionPlayerCoords = append(
			oppositionPlayerCoords,
			&Coordinates{
				Radial: stateHistories[int(index)].Values.At(
					0, PlayerStateValueIndices["Radial Position State"]),
				Angular: stateHistories[int(index)].Values.At(
					0, PlayerStateValueIndices["Angular Position State"]),
			},
		)
	}

	// Logic for player substitutions which will change these values
	spaceFindingTalent := int(params["player_space_finding_talent"][0])
	outputState[PlayerStateValueIndices["Ball Interaction Speed"]] =
		params["player_ball_interaction_speed"][0]
	outputState[PlayerStateValueIndices["Ball Interaction Inaccuracy"]] =
		params["player_ball_interaction_inaccuracy"][0]

	// TODO: Below is movement for attack but need to handle movement for defence

	// Logic for player movement and positioning when in possession
	_, bestProx := playerCoords.MinProximity(oppositionPlayerCoords)
	coordsOption := playerCoords
	plannedPlayerCoords := playerCoords
	for i := 0; i < spaceFindingTalent; i++ {
		coordsOption.Radial = f.uniformDist.Rand()
		coordsOption.Angular = f.uniformDist.Rand()
		_, prox := coordsOption.MinProximity(oppositionPlayerCoords)
		if prox > bestProx {
			plannedPlayerCoords.Radial = coordsOption.Radial
			plannedPlayerCoords.Angular = coordsOption.Angular
			bestProx = prox
		}
	}
	playerCoords.Update(
		plannedPlayerCoords,
		params["player_movement_speed"][0],
		timestepsHistory.NextIncrement,
	)
	outputState[PlayerStateValueIndices["Radial Position State"]] =
		playerCoords.Radial
	outputState[PlayerStateValueIndices["Angular Position State"]] =
		playerCoords.Angular

	return outputState
}

// FlounceballMatchStateIteration describes the iteration of a Flounceball
// match state in response to player positions and manager decisions.
type FlounceballMatchStateIteration struct {
	normDist *distuv.Normal
}

func (f *FlounceballMatchStateIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	f.normDist = &distuv.Normal{
		Mu:    0.0,
		Sigma: 1.0,
		Src:   rand.NewSource(settings.Seeds[partitionIndex]),
	}
}

func (f *FlounceballMatchStateIteration) Iterate(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	// Reorganise state data
	stateHistory := stateHistories[partitionIndex]
	outputState := stateHistory.Values.RawRowView(0)
	ballCoords := &Coordinates{
		Radial: stateHistory.Values.At(
			0, MatchStateValueIndices["Ball Radial Position State"]),
		Angular: stateHistory.Values.At(
			0, MatchStateValueIndices["Ball Angular Position State"]),
	}
	ballProjCoords := &Coordinates{
		Radial: stateHistory.Values.At(
			0, MatchStateValueIndices["Ball Projected Radial Position State"]),
		Angular: stateHistory.Values.At(
			0, MatchStateValueIndices["Ball Projected Angular Position State"]),
	}
	posMatcher := &PossessionNameMatcher{
		possession: int(stateHistory.Values.At(
			0, MatchStateValueIndices["Possession State"])),
	}

	// TODO: Use the attacking team tactics and basic heuristics to set the
	// initial intended ball projected position states.
	// Heuristics can be built from:
	// - the average proximity of defending players to the projected location
	// - the average proximity of attacking players to the projected location
	// - the total distance to the projected location

	// Apply the player ball interaction logic once the ball is within the interaction
	// radius of the projected location
	if ballCoords.Proximity(ballProjCoords) <= InteractionRadiusMetres {
		playerCoords := &Coordinates{}
		for i := 1; i < 11; i++ {
			// For the attacking players
			attackState := params[posMatcher.Attacking(
				"your_player_"+strconv.Itoa(i)+"_state",
				"other_player_"+strconv.Itoa(i)+"_state",
			)]
			playerCoords.Radial =
				attackState[PlayerStateValueIndices["Radial Position State"]]
			playerCoords.Angular =
				attackState[PlayerStateValueIndices["Angular Position State"]]
			if ballCoords.Proximity(playerCoords) <= InteractionRadiusMetres {
				// Add noise to the projected ball location based on player inaccuracy
				// - good attackers have lower inaccuracy
				f.normDist.Sigma =
					attackState[PlayerStateValueIndices["Ball Interaction Inaccuracy"]]
				ballProjCoords.ApplyShift(f.normDist.Rand(), f.normDist.Rand())
				// Add speed to the ball based on player ball interaction speed
				outputState[MatchStateValueIndices["Ball Speed State"]] =
					attackState[PlayerStateValueIndices["Ball Interaction Speed"]]
			}
			// For the defending players
			defendState := params[posMatcher.Defending(
				"your_player_"+strconv.Itoa(i)+"_state",
				"other_player_"+strconv.Itoa(i)+"_state",
			)]
			playerCoords.Radial =
				defendState[PlayerStateValueIndices["Radial Position State"]]
			playerCoords.Angular =
				defendState[PlayerStateValueIndices["Angular Position State"]]
			if ballCoords.Proximity(playerCoords) <= InteractionRadiusMetres {
				// Add noise to the projected ball location based on player accuracy
				// - good defenders have higher inaccuracy
				f.normDist.Sigma =
					defendState[PlayerStateValueIndices["Ball Interaction Inaccuracy"]]
				ballProjCoords.ApplyShift(f.normDist.Rand(), f.normDist.Rand())
			}
		}
	}

	// Compute the next ball radius and angle
	ballCoords.Update(
		ballProjCoords,
		stateHistory.Values.At(0, MatchStateValueIndices["Ball Speed State"]),
		timestepsHistory.NextIncrement,
	)
	outputState[MatchStateValueIndices["Ball Radial Position State"]] = ballCoords.Radial
	outputState[MatchStateValueIndices["Ball Angular Position State"]] = ballCoords.Angular

	// TODO: Logic for possession and total air time updates when ball goes out of play or hits ground
	// - out of play is easy logic but hitting the ground could be hard so simple logic is to say the
	// ball has a 'MaxBallFallingLifetime' when it reaches the interaction radius of the projected
	// location which gets reset when a player interacts with it

	// TODO: Logic for posession air time updates when ball is in play

	return outputState
}
