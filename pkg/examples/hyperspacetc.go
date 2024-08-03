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
	GetState func(
		key string,
		timeIndex int,
		stateHistory *simulator.StateHistory,
	) float64
	SetState    func(key string, value float64, outputState []float64)
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

func (s *SpacecraftLaneCountIteration) arrivals(
	outputState []float64,
	params simulator.Params,
	laneStateHistory *simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) {
	minEntryTimeIndex := int(s.GetState(
		"Min Upstream Entry Time Index In Queue",
		0,
		laneStateHistory,
	))
	for i := minEntryTimeIndex - 1; i >= 1; i-- {
		if s.GetState("Upstream Entry Detection", i, laneStateHistory) > 0.0 {
			queueSize := s.GetState("Downstream Queue Size", i, laneStateHistory)
			effectiveLaneLength := params["lane_length"][0] -
				(params["spacecraft_length"][0] * queueSize)
			timeSinceEntry := timestepsHistory.NextIncrement +
				timestepsHistory.Values.AtVec(0) - timestepsHistory.Values.AtVec(i)
			// This probabilistic sample draw answers the question: has the craft
			// reached the back of the queue?
			if s.uniformDist.Rand() < inverseGaussianCdf(
				timeSinceEntry,
				muFromParams(effectiveLaneLength, params["spacecraft_speed"][0]),
				lambdaFromParams(effectiveLaneLength, params["spacecraft_speed_variance"][0]),
			) {
				// If it has, then update the state values accordingly
				s.SetState("Downstream Queue Size", queueSize+1, outputState)
				s.SetState("Min Upstream Entry Time Index In Queue", float64(i), outputState)
				break
			}
		}
		i += 1
	}
}

func (s *SpacecraftLaneCountIteration) Iterate(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	// Create a state setter for convenience
	laneStateHistory := stateHistories[partitionIndex]
	outputState := laneStateHistory.Values.RawRowView(0)

	// TODO: Deal with the upstream entries into the lane from a node

	// Deal with upstream arrivals into the queue, assuming no overtaking
	// is allowed in this simple model
	s.arrivals(outputState, params, laneStateHistory, timestepsHistory)

	// TODO: Deal with the downstream departures from the queue into a node

	return outputState
}

// SpacecraftNodeCountIteration
type SpacecraftNodeCountIteration struct {
	GetState func(
		key string,
		timeIndex int,
		stateHistory *simulator.StateHistory,
	) float64
	SetState    func(key string, value float64, outputState []float64)
	uniformDist *distuv.Uniform
}

func (s *SpacecraftNodeCountIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	s.uniformDist = &distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(settings.Seeds[partitionIndex]),
	}
}

func (s *SpacecraftNodeCountIteration) Iterate(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	laneStateHistory := stateHistories[partitionIndex]
	outputState := laneStateHistory.Values.RawRowView(0)

	// TODO: Handle logic for connecting lanes together
	// TODO: Handle logic for moving spacecraft between connected lanes
	return outputState
}
