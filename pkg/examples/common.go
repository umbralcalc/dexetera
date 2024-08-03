package examples

import "github.com/umbralcalc/stochadex/pkg/simulator"

// GenerateStateValueGetter creates a closure which tidies up the
// code required to retrieve state values from a single time index in
// the history of a partition.
func GenerateStateValueGetter(
	stateValueIndicesMap map[string]int,
) func(key string, timeIndex int, stateHistory *simulator.StateHistory) float64 {
	return func(key string, timeIndex int, stateHistory *simulator.StateHistory) float64 {
		return stateHistory.Values.At(
			timeIndex,
			stateValueIndicesMap[key],
		)
	}
}

// GenerateStateValueSetter creates a closure which tidies up the
// code required to reassign state values.
func GenerateStateValueSetter(
	stateValueIndicesMap map[string]int,
) func(key string, value float64, outputState []float64) {
	return func(key string, value float64, outputState []float64) {
		outputState[stateValueIndicesMap[key]] = value
	}
}
