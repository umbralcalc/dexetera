package examples

import "github.com/umbralcalc/stochadex/pkg/simulator"

// EasyStateAccess is a struct which updates itself and provides easier
// state access to the iterations that hold it.
type EasyStateAccess struct {
	stateValueMaps []map[string]int
	stateHistories []*simulator.StateHistory
	outputState    []float64
	partitionIndex int
}

// Get retrieves state values from a single time index in
// the history of a partition specified by index.
func (e *EasyStateAccess) Get(
	valueName string,
	timeIndex int,
	partitionIndex int,
) float64 {
	return e.stateHistories[partitionIndex].Values.At(
		timeIndex,
		e.stateValueMaps[partitionIndex][valueName],
	)
}

// Set assigns state values at this point in time for the
// encapsulating iteration.
func (e *EasyStateAccess) Set(valueName string, value float64) {
	e.outputState[e.stateValueMaps[e.partitionIndex][valueName]] = value
}

// Update will update the state access struct given the input data.
func (e *EasyStateAccess) Update(
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
) {
	e.outputState = stateHistories[e.partitionIndex].Values.RawRowView(0)
	e.partitionIndex = partitionIndex
	e.stateHistories = stateHistories
}

// Output returns the state values which can be used as the iteration output.
func (e *EasyStateAccess) Output() []float64 {
	return e.outputState
}

// NewEasyStateAccess creates a new EasyStateAccess struct.
func NewEasyStateAccess(stateValueMaps []map[string]int) *EasyStateAccess {
	return &EasyStateAccess{stateValueMaps: stateValueMaps}
}
