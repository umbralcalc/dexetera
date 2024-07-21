//go:build js && wasm

package main

import (
	"strconv"

	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// PitchRadiusMetres is the radius of the circular pitch.
const PitchRadiusMetres = 100.0

// PossessionValueMap is a mapping to check which team is in possession
// based on the value of the possession state index.
var PossessionValueMap = map[int]string{0: "Your Team", 1: "Other Team"}

// MatchStateValueIndices is a mapping which helps with describing the
// meaning of the values for each match state index.
var MatchStateValueIndices = map[string]int{
	"Possession State":             0,
	"Your Team Total Air Time":     1,
	"Other Team Total Air Time":    2,
	"Ball Possession Air Time":     3,
	"Ball Radial Position State":   4,
	"Ball Angular Position State":  5,
	"Ball Vertical Position State": 6,
}

// PlayerStateValueIndices is a mapping which helps with describing the
// meaning of the values for each player state index.
var PlayerStateValueIndices = map[string]int{
	"Radial Position State":  0,
	"Angular Position State": 1,
}

// PlayerStateIteration describes the iteration of an individual player
// state in a Flounceball match.
type PlayerStateIteration struct {
}

func (p *PlayerStateIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
}

func (p *PlayerStateIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	return make([]float64, 0)
}

// MatchStateIteration describes the iteration of a Flounceball match
// state in response to player positions and manager decisions.
type MatchStateIteration struct {
}

func (m *MatchStateIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
}

func (m *MatchStateIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	ballRadius := stateHistories[partitionIndex].Values.At(
		0,
		MatchStateValueIndices["Ball Radial Position State"],
	)
	ballAngle := stateHistories[partitionIndex].Values.At(
		0,
		MatchStateValueIndices["Ball Angular Position State"],
	)
	for i := 1; i < 11; i++ {
		radiusAngle := params.FloatParams["your_player_"+strconv.Itoa(i)+"_radius_angle"]
	}
	return make([]float64, 0)
}

func main() {
	settings := &simulator.Settings{
		OtherParams: []*simulator.OtherParams{
			{
				FloatParams: map[string][]float64{
					"param_values": {1.0, 1.0, 1.0},
				},
				IntParams: map[string][]int64{},
			},
			{
				FloatParams: map[string][]float64{
					"rates":        {0.5, 1.0, 0.8, 1.0, 1.1},
					"gamma_alphas": {1.0, 2.5, 3.0, 1.8, 1.0},
					"gamma_betas":  {2.0, 1.0, 4.1, 2.0, 1.2},
				},
				IntParams: map[string][]int64{},
			},
			{
				FloatParams: map[string][]float64{
					"rates":        {1.5, 0.2, 0.6},
					"gamma_alphas": {2.3, 5.1, 2.0},
					"gamma_betas":  {2.0, 1.5, 1.1},
				},
				IntParams: map[string][]int64{},
			},
		},
		InitStateValues: [][]float64{
			{1.0, 1.0, 1.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0},
		},
		InitTimeValue:         0.0,
		Seeds:                 []uint64{0, 563, 8312},
		StateWidths:           []int{3, 5, 3},
		StateHistoryDepths:    []int{2, 2, 2},
		TimestepsHistoryDepth: 2,
	}
	partitions := []simulator.Partition{
		{Iteration: &simulator.ParamValuesIteration{}},
		{Iteration: &simulator.ConstantValuesIteration{}},
		{
			Iteration: &simulator.ConstantValuesIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"rates": 0,
			},
		},
	}
	for index, partition := range partitions {
		partition.Iteration.Configure(index, settings)
	}
	implementations := &simulator.Implementations{
		Partitions:      partitions,
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 100,
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 1.0},
	}
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
