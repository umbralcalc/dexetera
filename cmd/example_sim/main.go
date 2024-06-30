//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/phenomena"
	"github.com/umbralcalc/stochadex/pkg/simulator"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

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
				FloatParams: map[string][]float64{},
				IntParams: map[string][]int64{
					"action": {1.0, 1.0, 1.0},
				},
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
				Src: rand.NewSource(
					settings.Seeds[1],
				),
			},
		},
	}
	iteration1.Configure(1, settings)
	iteration2 := &phenomena.CompoundPoissonProcessIteration{
		JumpDist: &gammaJumpDistribution{
			dist: &distuv.Gamma{
				Alpha: 1.0,
				Beta:  1.0,
				Src: rand.NewSource(
					settings.Seeds[2],
				),
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
