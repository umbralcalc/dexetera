package examples

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// Planned approach:
// - Essentially a fake car-following model for the underlying dynamics
// of spacecraft.
// - Use the histogram node iteraton when constructing the node
// controller logic.

// SpacecraftFollowingLaneIteration
type SpacecraftFollowingLaneIteration struct {
}

func (s *SpacecraftFollowingLaneIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
}

func (s *SpacecraftFollowingLaneIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	return make([]float64, 0)
}
