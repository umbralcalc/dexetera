//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/phenomena"
	"github.com/umbralcalc/stochadex/pkg/simulator"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Rules of Flounceball:
// 1. The objective of the game is for each team to keep the ball in the
// air for as long as possible.
// 2. Possession alternates between each team, where the objective of the
// defensive side is to stop the ball from staying in the air by disrupting
// the play of the attacking side.
// 3. Defensive players are not allowed to touch the ball, otherwise the
// time the ball has remained in the air is subtracted from their team total.
// 4. Defensive players may tackle any attacking player within 10m radius of
// the lat/long ball location.
// 5. Player substitutions may occur continuously throughout the match and
// there is no limit to them.
// 6. The game takes place on a circular field and if the ball leaves this
// area at any time, then the possession of the attacking team is ended.
// 7. A match is 40 minutes of continuous play, without stoppages.

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
