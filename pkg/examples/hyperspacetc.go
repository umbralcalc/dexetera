package examples

import (
	"math"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Planned approach:
// - Essentially a fake car-following model for the underlying dynamics
// of spacecraft.
// - Use the histogram node iteraton when constructing the node
// controller logic.

// LaneCountSateValueIndices is a mapping which helps with describing
// the meaning of the values for each spacecraft lane count state index.
var LaneCountStateValueIndices = map[string]int{
	"Upstream Entry Detection":            0,
	"Downstream Exit Detection":           1,
	"Downstream Queue Size":               2,
	"Latest Upstream Entry Time In Queue": 3,
}

// standardNormalCdf returns the CDF value for a standard normal distribution.
func standardNormalCdf(x float64) float64 {
	return 0.5 * (1.0 + math.Erf(x/math.Sqrt(2.0)))
}

// inverseGaussianCdf returns the CDF value for an inverse-Gaussian distribution.
func inverseGaussianCdf(x float64, mu float64, lambda float64) float64 {
	return standardNormalCdf(math.Sqrt(lambda/x)*((x/mu)-1.0)) +
		(math.Exp(2.0*lambda/mu) * standardNormalCdf(-math.Sqrt(lambda/x)*((x/mu)+1.0)))
}

// muFromParams returns the mu value derived from lane length and spacecraft speed.
func muFromParams(laneLength float64, speed float64) float64 {
	return laneLength / speed
}

// lambdaFromParams returns the lambda value derived from lane length and spacecraft
// speed variance over their journey.
func lambdaFromParams(laneLength float64, speedVariance float64) float64 {
	return laneLength * laneLength / speedVariance
}

// generateLaneCountStateValueGetter creates a closure which reduces the
// amount of code required to retrieve state values.
func generateLaneCountStateValueGetter(
	stateHistory *simulator.StateHistory,
) func(key string) float64 {
	return func(key string) float64 {
		return stateHistory.Values.At(0, LaneCountStateValueIndices[key])
	}
}

// generateLaneCountStateValueSetter creates a closure which reduces the
// amount of code required to reassign state values.
func generateLaneCountStateValueSetter(
	state *OutputState,
) func(key string, value float64) {
	return func(key string, value float64) {
		state.Values[LaneCountStateValueIndices[key]] = value
	}
}

// SpacecraftLaneCountIteration
type SpacecraftLaneCountIteration struct {
	uniformDist *distuv.Uniform
}

func (s *SpacecraftLaneCountIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	s.uniformDist = &distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(settings.Seeds[partitionIndex]),
	}
}

func (s *SpacecraftLaneCountIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	// Reorganise state data and specify getters and setters for convenience.
	getState := generateLaneCountStateValueGetter(stateHistories[partitionIndex])
	outputState := &OutputState{Values: make([]float64, len(MatchStateValueIndices))}
	setState := generateLaneCountStateValueSetter(outputState)

	// TODO: get lane length parameter
	// TODO: get spacecraft lane speed parameter
	timeSinceEntry := 0.0
	laneLength := 0.0
	craftSpeed := 0.0
	craftSpeedVariance := 0.0
	craftLength := 0.0

	// has the craft reached the end of the queue?
	queueSize := getState("Downstream Queue Size")
	effectiveLaneLength := laneLength - (craftLength * queueSize)
	if s.uniformDist.Rand() < inverseGaussianCdf(
		timeSinceEntry,
		muFromParams(effectiveLaneLength, craftSpeed),
		lambdaFromParams(effectiveLaneLength, craftSpeedVariance),
	) {
		// if so, then update the lane count state accordingly
		setState("Downstream Queue Size", queueSize+1)
	}
	return outputState.Values
}
