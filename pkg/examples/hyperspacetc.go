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
	"Upstream Entry Count":                   0,
	"Downstream Exit Count":                  1,
	"Downstream Queue Size":                  2,
	"Min Upstream Entry Time Index In Queue": 3,
	"Time Since Last Exit":                   4,
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
	stateHistory *simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) {
	minEntryTimeIndex := int(stateHistory.Values.At(
		0, LaneCountStateValueIndices["Min Upstream Entry Time Index In Queue"]))
	for i := minEntryTimeIndex - 1; i >= 1; i-- {
		if stateHistory.Values.At(
			i, LaneCountStateValueIndices["Upstream Entry Count"]) > 0.0 {
			queueSize := stateHistory.Values.At(
				0, LaneCountStateValueIndices["Downstream Queue Size"])
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
				outputState[LaneCountStateValueIndices["Downstream Queue Size"]] = queueSize + 1
				outputState[LaneCountStateValueIndices["Min Upstream Entry Time Index In Queue"]] =
					float64(i)
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
	stateHistory := stateHistories[partitionIndex]
	outputState := stateHistory.Values.RawRowView(0)

	// Update the upstream entries into the lane from a lane connector
	outputState[LaneCountStateValueIndices["Upstream Entry Count"]] =
		stateHistories[int(params["lane_connector_partition"][0])].Values.At(
			0, int(params["lane_connector_value_index"][0]),
		)

	// Update the upstream arrivals into the queue, assuming no overtaking
	// is allowed in this simple model
	s.arrivals(outputState, params, stateHistory, timestepsHistory)

	// Update the downstream departures from the queue into a lane
	// connector which should be conditional on the lane having allowed flow
	outputState[LaneCountStateValueIndices["Downstream Exit Count"]] = 0.0
	if params["flow_allowed"][0] > 0.0 && stateHistory.Values.At(
		0, LaneCountStateValueIndices["Downstream Queue Size"]) > 0.0 {
		if stateHistory.Values.At(
			0,
			LaneCountStateValueIndices["Time Since Last Exit"],
		) > params["time_to_exit"][0] {
			outputState[LaneCountStateValueIndices["Downstream Exit Count"]] = 1.0
			outputState[LaneCountStateValueIndices["Downstream Queue Size"]] -= 1.0
			outputState[LaneCountStateValueIndices["Time Since Last Exit"]] = 0.0
		} else {
			outputState[LaneCountStateValueIndices["Time Since Last Exit"]] +=
				timestepsHistory.NextIncrement
		}
	} else {
		outputState[LaneCountStateValueIndices["Time Since Last Exit"]] = 0.0
	}

	return outputState
}

// SpacecraftLaneConnectorIteration
type SpacecraftLaneConnectorIteration struct {
	categoricalDist distuv.Categorical
}

func (s *SpacecraftLaneConnectorIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	weights := make([]float64, 0)
	for i := 0; i < settings.StateWidths[partitionIndex]; i++ {
		weights = append(weights, 1.0)
	}
	s.categoricalDist = distuv.NewCategorical(
		weights,
		rand.NewSource(settings.Seeds[partitionIndex]),
	)
}

func (s *SpacecraftLaneConnectorIteration) Iterate(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	stateHistory := stateHistories[partitionIndex]
	outputState := make([]float64, stateHistory.StateWidth)
	for i := 0; i < stateHistory.StateWidth; i++ {
		outputState = append(outputState, 0.0)
	}
	for i, count := range params["connected_partition_input_counts"] {
		if count > 0.0 {
			s.categoricalDist.Reweight(i, 0.0) // don't exit by same lane
			outputState[int(s.categoricalDist.Rand())] = count
			s.categoricalDist.Reweight(i, 1.0)
		}
	}
	return outputState
}
