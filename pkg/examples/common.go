package examples

import "github.com/umbralcalc/stochadex/pkg/simulator"

// OutputState is a simple container for referencing states on output.
type OutputState struct {
	Values []float64
}

// GenerateStateValueGetter creates a closure which tidies up the
// code required to retrieve state values from a single time index in
// the history of a partition.
func GenerateStateValueGetter(
	stateValueIndicesMap map[string]int,
	stateHistory *simulator.StateHistory,
) func(key string, timeIndex int) float64 {
	return func(key string, timeIndex int) float64 {
		return stateHistory.Values.At(timeIndex, stateValueIndicesMap[key])
	}
}

// GenerateStateValueSetter creates a closure which tidies up the
// code required to reassign state values.
func GenerateStateValueSetter(
	stateValueIndicesMap map[string]int,
	outputState *OutputState,
) func(key string, value float64) {
	return func(key string, value float64) {
		outputState.Values[stateValueIndicesMap[key]] = value
	}
}

// GenerateMultiStateValuesGetter creates a closure which tidies up
// the code required to retrieve state values from a single time index
// in the history of multiple partitions.
func GenerateMultiStateValuesGetter(
	stateValueIndicesMap map[string]int,
	stateHistories []*simulator.StateHistory,
) func(keys []string, partitionIndices []int64, timeIndex int) [][]float64 {
	return func(keys []string, partitionIndices []int64, timeIndex int) [][]float64 {
		values := make([][]float64, 0)
		for _, index := range partitionIndices {
			valuesVec := make([]float64, 0)
			for _, key := range keys {
				valuesVec = append(
					valuesVec,
					stateHistories[index].Values.At(timeIndex, stateValueIndicesMap[key]),
				)
			}
			values = append(values, valuesVec)
		}
		return values
	}
}
