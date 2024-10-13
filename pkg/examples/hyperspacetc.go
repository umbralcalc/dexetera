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

// LabelToEventIndex maps the event label to the relevant index.
var LabelToEventIndex = map[string]int{
	"Nothing":    0,
	"Entry":      1,
	"Exit":       2,
	"Entry&Exit": 3,
}

// EntryExitToEventIndex is a convenience function to map from entry
// and exit event booleans to the appropriate event index.
func EntryExitToEventIndex(entry, exit bool) int {
	type boolPair struct{ entry, exit bool }
	switch (boolPair{entry: entry, exit: exit}) {
	case boolPair{entry: false, exit: false}:
		return LabelToEventIndex["Nothing"]
	case boolPair{entry: true, exit: false}:
		return LabelToEventIndex["Entry"]
	case boolPair{entry: false, exit: true}:
		return LabelToEventIndex["Exit"]
	case boolPair{entry: true, exit: true}:
		return LabelToEventIndex["Entry&Exit"]
	}
	panic("Couldn't find event")
}

// EntryBoolFromUpstreamPartition retrieves the entry event boolean
// based on the data which is in the relevant upstream.
func EntryBoolFromUpstreamPartition(
	params simulator.Params,
	stateHistories []*simulator.StateHistory,
) bool {
	entry := false
	// Get the upstream entries from an upstream if it exists
	if upPart, ok := params["upstream_partition"]; ok {
		entry = stateHistories[int(upPart[0])].Values.At(
			0, int(params["upstream_value_index"][0]),
		) != params["empty_upstream_value"][0]
	}
	return entry
}

// SpacecraftQueueEventFunction returns the index of the latest event
// to happen to the spacecraft queue at the end of the line.
func SpacecraftQueueEventFunction(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	entry := EntryBoolFromUpstreamPartition(params, stateHistories)
	exit := false
	// Get a downstream exit if flow is allowed and collection isn't empty
	if int(params["partition_flow_allowed"][0]) == partitionIndex {
		emptyValue := params["empty_value"][0]
		stateHistory := stateHistories[partitionIndex]
		for i := 1; i < stateHistory.StateWidth; i++ {
			if stateHistory.Values.At(0, i) != emptyValue {
				exit = true
				break
			}
		}
	}
	return []float64{float64(EntryExitToEventIndex(entry, exit))}
}

// EntryTimeFromUpstreamPushFunction retrieves the next values to
// push from the popped values of the configured upstream.
func EntryTimeFromUpstreamPushFunction(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) ([]float64, bool) {
	if params["push_empty"][0] == stateHistories[int(
		params["upstream_partition"][0])].Values.At(
		0, int(params["upstream_state_value_index"][0])) {
		return nil, false
	}
	return []float64{
		timestepsHistory.Values.AtVec(0) + timestepsHistory.NextIncrement,
	}, true
}

// standardNormalCdf returns the CDF value for a standard normal distribution.
func standardNormalCdf(x float64) float64 {
	return 0.5 * (1.0 + math.Erf(x/math.Sqrt(2.0)))
}

// inverseGaussianCdf returns the CDF value for an inverse-Gaussian distribution.
func inverseGaussianCdf(x, mu, lambda float64) float64 {
	return standardNormalCdf(math.Sqrt(lambda/x)*((x/mu)-1.0)) +
		(math.Exp(2.0*lambda/mu) * standardNormalCdf(-math.Sqrt(lambda/x)*((x/mu)+1.0)))
}

// muFromParams returns the mu value derived from line length and spacecraft speed.
func muFromParams(lineLength, speed float64) float64 {
	return lineLength / speed
}

// lambdaFromParams returns the lambda value derived from line length and spacecraft
// speed variance over their journey.
func lambdaFromParams(lineLength, speedVariance float64) float64 {
	return lineLength * lineLength / speedVariance
}

// SpacecraftLineEventIteration returns the index of the latest event
// to happen to the spacecraft line, where exits correspond to downstream
// arrivals into the queue at the end of the line.
type SpacecraftLineEventIteration struct {
	uniformDist *distuv.Uniform
}

func (s *SpacecraftLineEventIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	s.uniformDist = &distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(settings.Seeds[partitionIndex]),
	}
}

func (s *SpacecraftLineEventIteration) Iterate(
	params simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	entry := EntryBoolFromUpstreamPartition(params, stateHistories)
	// Get a downstream arrival depending on a probability distribution of
	// line traversal times
	downstreamArrival := false
	queueSize := params["queue_size"][0]
	emptyValue := params["empty_value"][0]
	stateHistory := stateHistories[partitionIndex]
	for i := 1; i < stateHistory.StateWidth; i++ {
		timeSinceEntry := stateHistory.Values.At(0, i)
		if timeSinceEntry != emptyValue {
			effectiveLineLength := params["line_length"][0] -
				(params["spacecraft_length"][0] * queueSize)
			if s.uniformDist.Rand() < inverseGaussianCdf(
				timeSinceEntry,
				muFromParams(effectiveLineLength, params["spacecraft_speed"][0]),
				lambdaFromParams(
					effectiveLineLength, params["spacecraft_speed_variance"][0]),
			) {
				downstreamArrival = true
				break
			}
		}
	}
	return []float64{float64(EntryExitToEventIndex(entry, downstreamArrival))}
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
		value := params["partition_"+strconv.Itoa(int(index))+"_input_value"][0]
		if value > 0.0 {
			outputState[int(s.categoricalDist.Rand())] = value
		}
	}
	return outputState
}
