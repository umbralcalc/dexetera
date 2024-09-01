//go:build js && wasm

package main

import (
	"math"
	"strconv"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"

	"github.com/umbralcalc/dexetera/pkg/examples"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func main() {
	totalMatchSeconds := 300.0
	timeStepsizeSeconds := 0.01
	seeds := make([]uint64, 0)
	stateWidths := make([]int, 0)
	stateHistoryDepths := make([]int, 0)
	initStateValues := make([][]float64, 0)
	params := make([]simulator.Params, 0)
	partitions := make([]simulator.Partition, 0)
	matchStateParamsFromUpstreamPartition := make(map[string]int, 0)
	seeds = append(seeds, 0)
	stateWidths = append(stateWidths, 3)
	stateHistoryDepths = append(stateHistoryDepths, 1)
	initStateValues = append(initStateValues, []float64{0.0, 0.0, 0.0})
	params = append(params, simulator.Params{"param_values": {0.0, 0.0, 0.0}})
	partitions = append(
		partitions,
		simulator.Partition{Iteration: &simulator.ParamValuesIteration{}},
	)
	uniformDist := &distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(uint64(rand.Intn(10000))),
	}
	index := len(partitions)
	for i := 1; i < 11; i++ {
		seeds = append(seeds, uint64(rand.Intn(10000)))
		stateWidths = append(stateWidths, 2)
		stateHistoryDepths = append(stateHistoryDepths, 1)
		initStateValues = append(
			initStateValues,
			[]float64{
				examples.PitchRadiusMetres * uniformDist.Rand(),
				2.0 * math.Pi * uniformDist.Rand(),
			},
		)
		params = append(params, simulator.Params{})
		yourPlayerIteration := &examples.FlounceballPlayerStateIteration{}
		partitions = append(
			partitions,
			simulator.Partition{Iteration: yourPlayerIteration},
		)
		matchStateParamsFromUpstreamPartition["your_player_"+
			strconv.Itoa(i)+"_state"] = index
		index += 1
	}
	for i := 1; i < 11; i++ {
		seeds = append(seeds, uint64(rand.Intn(10000)))
		stateWidths = append(stateWidths, 2)
		stateHistoryDepths = append(stateHistoryDepths, 1)
		initStateValues = append(
			initStateValues,
			[]float64{
				examples.PitchRadiusMetres * uniformDist.Rand(),
				2.0 * math.Pi * uniformDist.Rand(),
			},
		)
		params = append(params, simulator.Params{})
		otherPlayerIteration := &examples.FlounceballPlayerStateIteration{}
		partitions = append(
			partitions,
			simulator.Partition{Iteration: otherPlayerIteration},
		)
		matchStateParamsFromUpstreamPartition["other_player_"+
			strconv.Itoa(i)+"_state"] = index
		index += 1
	}
	seeds = append(seeds, uint64(rand.Intn(10000)))
	stateWidths = append(stateWidths, 2)
	stateHistoryDepths = append(stateHistoryDepths, 1)
	initStateValues = append(initStateValues, []float64{0.0, 0.0})
	params = append(params, simulator.Params{})
	matchIteration := &examples.FlounceballMatchStateIteration{}
	partitions = append(
		partitions,
		simulator.Partition{
			Iteration:                   matchIteration,
			ParamsFromUpstreamPartition: matchStateParamsFromUpstreamPartition,
		},
	)
	settings := &simulator.Settings{
		Params:                params,
		InitStateValues:       initStateValues,
		InitTimeValue:         0.0,
		Seeds:                 seeds,
		StateWidths:           stateWidths,
		StateHistoryDepths:    stateHistoryDepths,
		TimestepsHistoryDepth: 1,
	}
	for i, partition := range partitions {
		partition.Iteration.Configure(i, settings)
	}
	implementations := &simulator.Implementations{
		Partitions:      partitions,
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: int(totalMatchSeconds / timeStepsizeSeconds),
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: timeStepsizeSeconds},
	}
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
