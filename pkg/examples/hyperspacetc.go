package examples

import (
	"math"
	"strconv"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// LineCountSateValueIndices is a mapping which helps with describing
// the meaning of the values for each spacecraft line count state index.
var LineCountStateValueIndices = map[string]int{
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

// muFromParams returns the mu value derived from line length and spacecraft speed.
func muFromParams(lineLength float64, speed float64) float64 {
	return lineLength / speed
}

// lambdaFromParams returns the lambda value derived from line length and spacecraft
// speed variance over their journey.
func lambdaFromParams(lineLength float64, speedVariance float64) float64 {
	return lineLength * lineLength / speedVariance
}

// SpacecraftLineCountIteration iterates the state of a hyperspace line in the
// Hyperspace Traffic Control example.
type SpacecraftLineCountIteration struct {
	uniformDist *distuv.Uniform
}

func (s *SpacecraftLineCountIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	s.uniformDist = &distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(settings.Seeds[partitionIndex]),
	}
}

func (s *SpacecraftLineCountIteration) Iterate(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	stateHistory := stateHistories[partitionIndex]
	outputState := stateHistory.Values.RawRowView(0)

	// Update the upstream entries into the line from a line connector if it exists
	if connectorParts, ok := params["upstream_partition"]; ok {
		outputState[LineCountStateValueIndices["Upstream Entry Count"]] =
			stateHistories[int(connectorParts[0])].Values.At(
				0, int(params["upstream_value_index"][0]),
			)
	}

	// Update the upstream arrivals into the queue, assuming no overtaking
	// is allowed in this simple model
	minEntryTimeIndex := int(stateHistory.Values.At(0,
		LineCountStateValueIndices["Min Upstream Entry Time Index In Queue"]))
	for i := minEntryTimeIndex - 1; i >= 1; i-- {
		if stateHistory.Values.At(
			i, LineCountStateValueIndices["Upstream Entry Count"]) > 0.0 {
			queueSize := stateHistory.Values.At(
				0, LineCountStateValueIndices["Downstream Queue Size"])
			effectiveLineLength := params["line_length"][0] -
				(params["spacecraft_length"][0] * queueSize)
			timeSinceEntry := timestepsHistory.NextIncrement +
				timestepsHistory.Values.AtVec(0) - timestepsHistory.Values.AtVec(i)
			// This probabilistic sample draw answers the question: has the craft
			// reached the back of the queue?
			if s.uniformDist.Rand() < inverseGaussianCdf(
				timeSinceEntry,
				muFromParams(effectiveLineLength, params["spacecraft_speed"][0]),
				lambdaFromParams(
					effectiveLineLength, params["spacecraft_speed_variance"][0]),
			) {
				// If it has, then update the state values accordingly
				outputState[LineCountStateValueIndices["Downstream Queue Size"]] +=
					queueSize + 1
				outputState[LineCountStateValueIndices["Min Upstream Entry Time Index In Queue"]] =
					float64(i)
				break
			}
		}
		i += 1
	}

	// Update the downstream departures from the queue into a line
	// connector which should be conditional on the line having allowed flow
	outputState[LineCountStateValueIndices["Downstream Exit Count"]] = 0.0
	if params["flow_allowed"][0] > 0.0 && stateHistory.Values.At(
		0, LineCountStateValueIndices["Downstream Queue Size"]) > 0.0 {
		if stateHistory.Values.At(
			0, LineCountStateValueIndices["Time Since Last Exit"],
		) > params["time_to_exit"][0] {
			outputState[LineCountStateValueIndices["Downstream Exit Count"]] = 1.0
			outputState[LineCountStateValueIndices["Downstream Queue Size"]] -= 1.0
			outputState[LineCountStateValueIndices["Time Since Last Exit"]] = 0.0
		} else {
			outputState[LineCountStateValueIndices["Time Since Last Exit"]] +=
				timestepsHistory.NextIncrement
		}
	} else {
		outputState[LineCountStateValueIndices["Time Since Last Exit"]] = 0.0
	}

	return outputState
}

// SpacecraftLineConnectorIteration iterates the state of a connection between
// hyperspace lines in the Hyperspace Traffic Control example.
type SpacecraftLineConnectorIteration struct {
	categoricalDist distuv.Categorical
}

func (s *SpacecraftLineConnectorIteration) Configure(
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

func (s *SpacecraftLineConnectorIteration) Iterate(
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
	for _, index := range params["connected_incoming_partitions"] {
		count := params["partition_"+strconv.Itoa(int(index))+"_input_count"][0]
		if count > 0.0 {
			outputState[int(s.categoricalDist.Rand())] = count
		}
	}
	return outputState
}
