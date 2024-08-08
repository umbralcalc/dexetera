package examples

import (
	"math"
	"strconv"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// LaneCountSateValueIndices is a mapping which helps with describing
// the meaning of the values for each spacecraft lane count state index.
var LaneCountStateValueIndices = map[string]int{
	"Upstream Entry Count":                            0,
	"Downstream Exit Count":                           1,
	"Downstream Queue Size":                           2,
	"Min Upstream Entry Time Index In Queue":          3,
	"Time Since Last Exit":                            4,
	"Opposing Upstream Entry Count":                   5,
	"Opposing Downstream Exit Count":                  6,
	"Opposing Downstream Queue Size":                  7,
	"Opposing Min Upstream Entry Time Index In Queue": 8,
	"Opposing Time Since Last Exit":                   9,
}

// laneIndexFromBaseKeyName returns the state value index for the lane
// count using the base key name and whether or not it is opposing.
func laneIndexFromBaseKeyName(baseKeyName string, opposing bool) int {
	if opposing {
		return LaneCountStateValueIndices["Opposing "+baseKeyName]
	} else {
		return LaneCountStateValueIndices[baseKeyName]
	}
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

// SpacecraftLaneCountIteration iterates the state of a hyperspace lane in the
// Hyperspace Traffic Control example.
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
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	stateHistory := stateHistories[partitionIndex]
	outputState := stateHistory.Values.RawRowView(0)

	// handle both lane directions in turn
	for _, opposing := range []bool{true, false} {
		oppStr := ""
		if opposing {
			oppStr = "opposing_"
		}
		// Update the upstream entries into the lane from a lane connector if it exists
		if connectorParts, ok := params[oppStr+
			"upstream_lane_connector_partition"]; ok {
			outputState[laneIndexFromBaseKeyName(
				"Upstream Entry Count", opposing)] =
				stateHistories[int(connectorParts[0])].Values.At(
					0, int(params[oppStr+
						"upstream_lane_connector_value_index"][0]),
				)
		}

		// Update the upstream arrivals into the queue, assuming no overtaking
		// is allowed in this simple model
		minEntryTimeIndex := int(
			stateHistory.Values.At(0, laneIndexFromBaseKeyName(
				"Min Upstream Entry Time Index In Queue", opposing)))
		for i := minEntryTimeIndex - 1; i >= 1; i-- {
			if stateHistory.Values.At(
				i, laneIndexFromBaseKeyName("Upstream Entry Count", opposing)) > 0.0 {
				queueSize := stateHistory.Values.At(
					0, laneIndexFromBaseKeyName("Downstream Queue Size", opposing))
				effectiveLaneLength := params["lane_length"][0] -
					(params["spacecraft_length"][0] * queueSize)
				timeSinceEntry := timestepsHistory.NextIncrement +
					timestepsHistory.Values.AtVec(0) - timestepsHistory.Values.AtVec(i)
				// This probabilistic sample draw answers the question: has the craft
				// reached the back of the queue?
				if s.uniformDist.Rand() < inverseGaussianCdf(
					timeSinceEntry,
					muFromParams(effectiveLaneLength, params["spacecraft_speed"][0]),
					lambdaFromParams(
						effectiveLaneLength, params["spacecraft_speed_variance"][0]),
				) {
					// If it has, then update the state values accordingly
					outputState[laneIndexFromBaseKeyName(
						"Downstream Queue Size", opposing)] = queueSize + 1
					outputState[laneIndexFromBaseKeyName(
						"Min Upstream Entry Time Index In Queue", opposing)] =
						float64(i)
					break
				}
			}
			i += 1
		}

		// Update the downstream departures from the queue into a lane
		// connector which should be conditional on the lane having allowed flow
		outputState[laneIndexFromBaseKeyName(
			"Downstream Exit Count", opposing)] = 0.0
		if params["flow_allowed"][0] > 0.0 && stateHistory.Values.At(
			0, laneIndexFromBaseKeyName("Downstream Queue Size", opposing)) > 0.0 {
			if stateHistory.Values.At(
				0,
				laneIndexFromBaseKeyName("Time Since Last Exit", opposing),
			) > params["time_to_exit"][0] {
				outputState[laneIndexFromBaseKeyName(
					"Downstream Exit Count", opposing)] = 1.0
				outputState[laneIndexFromBaseKeyName(
					"Downstream Queue Size", opposing)] -= 1.0
				outputState[laneIndexFromBaseKeyName(
					"Time Since Last Exit", opposing)] = 0.0
			} else {
				outputState[laneIndexFromBaseKeyName(
					"Time Since Last Exit", opposing)] +=
					timestepsHistory.NextIncrement
			}
		} else {
			outputState[laneIndexFromBaseKeyName(
				"Time Since Last Exit", opposing)] = 0.0
		}
	}

	return outputState
}

// SpacecraftLaneConnectorIteration iterates the state of a connection between
// hyperspace lanes in the Hyperspace Traffic Control example.
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
		outputState[i] = 0.0
	}
	for i, index := range params["connected_partitions"] {
		count := params["partition_"+strconv.Itoa(int(index))+"_input_count"][0]
		if count > 0.0 {
			s.categoricalDist.Reweight(i, 0.0) // don't exit by same lane
			outputState[int(s.categoricalDist.Rand())] = count
			s.categoricalDist.Reweight(i, 1.0)
		}
	}
	return outputState
}
