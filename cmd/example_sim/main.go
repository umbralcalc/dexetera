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
	settings := simulator.LoadSettingsFromYaml("settings.yaml")
	partitions := make([]simulator.Partition, 0)
	partitions = append(
		partitions,
		simulator.Partition{Iteration: &simulator.ConstantValuesIteration{}},
	)
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
	partitions = append(
		partitions,
		simulator.Partition{Iteration: iteration1},
	)
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
	partitions = append(
		partitions,
		simulator.Partition{
			Iteration: iteration2,
			ParamsFromUpstreamPartition: map[string]int{
				"rates": 0,
			},
		},
	)
	implementations := &simulator.Implementations{
		Partitions:      partitions,
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 100,
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 1.0},
	}
	websocketPartitionIndex := 0
	handle := "/simio"
	address := ":2112"
	simio.RegisterRun(
		settings,
		implementations,
		websocketPartitionIndex,
		handle,
		address,
	)
}
