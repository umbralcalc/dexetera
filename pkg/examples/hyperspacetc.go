package examples

import (
	"math"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Planned approach:
// - Use the histogram node iteraton when constructing the node
// controller logic.

// LaneCountSateValueIndices is a mapping which helps with describing
// the meaning of the values for each spacecraft lane count state index.
var LaneCountStateValueIndices = map[string]int{
	"Upstream Entry Detection":               0,
	"Downstream Exit Detection":              1,
	"Downstream Queue Size":                  2,
	"Min Upstream Entry Time Index In Queue": 3,
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

// SpacecraftLaneCountIteration
type SpacecraftLaneCountIteration struct {
	GetLaneState func(key string, timeIndex int) float64
	uniformDist  *distuv.Uniform
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
	// Create a state setter for convenience.
	outputState := stateHistories[partitionIndex].Values.RawRowView(0)
	setState := GenerateStateValueSetter(LaneCountStateValueIndices, outputState)

	// TODO: get lane length parameter
	// TODO: get spacecraft lane speed parameter
	timeSinceEntry := 0.0
	laneLength := 0.0
	craftSpeed := 0.0
	craftSpeedVariance := 0.0
	craftLength := 0.0

	// assume no overtaking for this simple model
	minEntryTimeIndex := int(s.GetLaneState("Min Upstream Entry Time Index In Queue", 0))
	for i := minEntryTimeIndex - 1; i >= 1; i-- {
		if s.GetLaneState("Upstream Entry Detection", i) > 0.0 {
			// has the craft reached the end of the queue?
			queueSize := s.GetLaneState("Downstream Queue Size", i)
			effectiveLaneLength := laneLength - (craftLength * queueSize)
			if s.uniformDist.Rand() < inverseGaussianCdf(
				timeSinceEntry,
				muFromParams(effectiveLaneLength, craftSpeed),
				lambdaFromParams(effectiveLaneLength, craftSpeedVariance),
			) {
				// if so, then update the lane count state accordingly
				setState("Downstream Queue Size", queueSize+1)
				break
			}
		}
		i += 1
	}

	return outputState
}
