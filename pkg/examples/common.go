package examples

import "github.com/umbralcalc/stochadex/pkg/simulator"

// OutputState is a simple container for referencing states on output.
type OutputState struct {
	Values []float64
}

// GenerateStateValueGetter creates a closure which tidies up the
// code required to retrieve state values.
func GenerateStateValueGetter(
	stateValueIndicesMap map[string]int,
	stateHistory *simulator.StateHistory,
) func(key string) float64 {
	return func(key string) float64 {
		return stateHistory.Values.At(0, stateValueIndicesMap[key])
	}
}

// GenerateStateValueSetter creates a closure which tidies up the
// code required to reassign state values.
func GenerateStateValueSetter(
	stateValueIndicesMap map[string]int,
	state *OutputState,
) func(key string, value float64) {
	return func(key string, value float64) {
		state.Values[stateValueIndicesMap[key]] = value
	}
}

// GenerateMultiStateValuesGetter creates a closure which tidies up
// the code required to retrieve state values for multiple partitions.
func GenerateMultiStateValuesGetter(
	stateValueIndicesMap map[string]int,
	stateHistories []*simulator.StateHistory,
) func(partitionIndices []int64, keys []string) [][]float64 {
	return func(partitionIndices []int64, keys []string) [][]float64 {
		values := make([][]float64, 0)
		for _, index := range partitionIndices {
			valuesVec := make([]float64, 0)
			for _, key := range keys {
				valuesVec = append(
					valuesVec,
					stateHistories[index].Values.At(0, stateValueIndicesMap[key]),
				)
			}
			values = append(values, valuesVec)
		}
		return values
	}
}
