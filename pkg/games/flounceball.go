package games

import (
	"math"

	"math/rand/v2"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/stat/distuv"
)

// PitchRadiusMetres is the radius of the circular pitch.
const PitchRadiusMetres = 100.0

// MatchStateValueIndices is a mapping which helps with describing the
// meaning of the values for each match state index.
var MatchStateValueIndices = map[string]int{
	"Your Team Total Points":  0,
	"Other Team Total Points": 1,
}

// PlayerStateValueIndices is a mapping which helps with describing the
// meaning of the values for each player state index.
var PlayerStateValueIndices = map[string]int{
	"Radial Position State":  0,
	"Angular Position State": 1,
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
		Src: rand.NewPCG(
			settings.Iterations[partitionIndex].Seed,
			settings.Iterations[partitionIndex].Seed,
		),
	}
}

func (f *FlounceballPlayerStateIteration) Iterate(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	outputState := stateHistories[partitionIndex].Values.RawRowView(0)
	if p, ok := params.GetOk("manager_directed_coordinates"); ok {
		outputState[PlayerStateValueIndices["Radial Position State"]] = p[0]
		outputState[PlayerStateValueIndices["Angular Position State"]] = p[1]
	} else {
		// Opposition manager randomly positions players
		outputState[PlayerStateValueIndices["Radial Position State"]] =
			f.uniformDist.Rand() * PitchRadiusMetres
		outputState[PlayerStateValueIndices["Angular Position State"]] =
			f.uniformDist.Rand() * 2.0 * math.Pi
	}
	return outputState
}

// proximity returns the proximity (in terms of absolute distance) of the
// input coordinates to the entity.
func proximity(radial1, angular1, radial2, angular2 float64) float64 {
	diffX := (radial1 * math.Cos(angular1)) - (radial2 * math.Cos(angular2))
	diffY := (radial1 * math.Sin(angular1)) - (radial2 * math.Sin(angular2))
	return math.Sqrt((diffX * diffX) + (diffY * diffY))
}

// FlounceballMatchStateValuesFunction describes the iteration of a
// Flounceball match state in response to player positions.
func FlounceballMatchStateValuesFunction(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	outputState := stateHistories[partitionIndex].Values.RawRowView(0)
	isYourWin := params["team_proximity_minima"][0] < params["team_proximity_minima"][1]
	if isYourWin {
		outputState[MatchStateValueIndices["Your Team Total Points"]] += 1
	} else {
		outputState[MatchStateValueIndices["Other Team Total Points"]] += 1
	}
	return outputState
}
