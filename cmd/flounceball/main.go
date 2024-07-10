//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/phenomena"
	"github.com/umbralcalc/stochadex/pkg/simulator"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// PitchRadiusMetres is the radius of the circular pitch.
const PitchRadiusMetres = 100.0

// PossessionValueMap is a mapping to check which team is in possession
// based on the value of the possession state index.
var PossessionValueMap = map[int]string{0: "Your Team", 1: "Other Team"}

// MatchStateValueIndices is a mapping which helps with describing the
// meaning of the values for each match state index.
var MatchStateValueIndices = map[string]int{
	"Possession State":                            0,
	"Your Team Total Time Score":                  1,
	"Other Team Total Time Score":                 2,
	"Ball Radial Position State":                  3,
	"Ball Angular Position State":                 4,
	"Your Team Player 1 Radial Position State":    5,
	"Your Team Player 1 Angular Position State":   6,
	"Your Team Player 2 Radial Position State":    7,
	"Your Team Player 2 Angular Position State":   8,
	"Your Team Player 3 Radial Position State":    9,
	"Your Team Player 3 Angular Position State":   10,
	"Your Team Player 4 Radial Position State":    11,
	"Your Team Player 4 Angular Position State":   12,
	"Your Team Player 5 Radial Position State":    13,
	"Your Team Player 5 Angular Position State":   14,
	"Your Team Player 6 Radial Position State":    15,
	"Your Team Player 6 Angular Position State":   16,
	"Your Team Player 7 Radial Position State":    17,
	"Your Team Player 7 Angular Position State":   18,
	"Your Team Player 8 Radial Position State":    19,
	"Your Team Player 8 Angular Position State":   20,
	"Your Team Player 9 Radial Position State":    21,
	"Your Team Player 9 Angular Position State":   22,
	"Your Team Player 10 Radial Position State":   23,
	"Your Team Player 10 Angular Position State":  24,
	"Other Team Player 1 Radial Position State":   25,
	"Other Team Player 1 Angular Position State":  26,
	"Other Team Player 2 Radial Position State":   27,
	"Other Team Player 2 Angular Position State":  28,
	"Other Team Player 3 Radial Position State":   29,
	"Other Team Player 3 Angular Position State":  30,
	"Other Team Player 4 Radial Position State":   31,
	"Other Team Player 4 Angular Position State":  32,
	"Other Team Player 5 Radial Position State":   33,
	"Other Team Player 5 Angular Position State":  34,
	"Other Team Player 6 Radial Position State":   35,
	"Other Team Player 6 Angular Position State":  36,
	"Other Team Player 7 Radial Position State":   37,
	"Other Team Player 7 Angular Position State":  38,
	"Other Team Player 8 Radial Position State":   39,
	"Other Team Player 8 Angular Position State":  40,
	"Other Team Player 9 Radial Position State":   41,
	"Other Team Player 9 Angular Position State":  42,
	"Other Team Player 10 Radial Position State":  43,
	"Other Team Player 10 Angular Position State": 44,
}

// gammaJumpDistribution jumps the compound Poisson process with samples
// drawn from a gamma distribution - this is just for testing.
type gammaJumpDistribution struct {
	dist *distuv.Gamma
}

func (g *gammaJumpDistribution) NewJump(
	params *simulator.OtherParams,
	stateElement int,
) float64 {
	g.dist.Alpha = params.FloatParams["gamma_alphas"][stateElement]
	g.dist.Beta = params.FloatParams["gamma_betas"][stateElement]
	return g.dist.Rand()
}

func main() {
	settings := &simulator.Settings{
		OtherParams: []*simulator.OtherParams{
			{
				FloatParams: map[string][]float64{
					"action": {1.0, 1.0, 1.0},
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
	iteration0 := &simio.ActionParamsIteration{}
	iteration0.Configure(0, settings)
	iteration1 := &phenomena.CompoundPoissonProcessIteration{
		JumpDist: &gammaJumpDistribution{
			dist: &distuv.Gamma{
				Alpha: 1.0,
				Beta:  1.0,
				Src:   rand.NewSource(settings.Seeds[1]),
			},
		},
	}
	iteration1.Configure(1, settings)
	iteration2 := &phenomena.CompoundPoissonProcessIteration{
		JumpDist: &gammaJumpDistribution{
			dist: &distuv.Gamma{
				Alpha: 1.0,
				Beta:  1.0,
				Src:   rand.NewSource(settings.Seeds[2]),
			},
		},
	}
	iteration2.Configure(2, settings)
	partitions := []simulator.Partition{
		{Iteration: iteration0},
		{Iteration: iteration1},
		{
			Iteration: iteration2,
			ParamsFromUpstreamPartition: map[string]int{
				"rates": 0,
			},
		},
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
