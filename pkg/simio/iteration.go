package simio

import "github.com/umbralcalc/stochadex/pkg/simulator"

// ActionParamsIteration implements an iteration in the stochadex
// which directly outputs the configured "action" params.
type ActionParamsIteration struct{}

func (a *ActionParamsIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
}

func (a *ActionParamsIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	return params.FloatParams["action"]
}
