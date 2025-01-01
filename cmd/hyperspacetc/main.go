//go:build js && wasm

package main

import (
	"math/rand"

	"github.com/umbralcalc/dexetera/pkg/games"
	"github.com/umbralcalc/stochadex/pkg/discrete"
	"github.com/umbralcalc/stochadex/pkg/general"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func main() {
	seeds := make([]uint64, 0)
	for i := 0; i < 15; i++ {
		seeds = append(seeds, uint64(rand.Intn(10000)))
	}
	settings := &simulator.Settings{
		Params: []simulator.Params{
			{
				// action taker
				"param_values": {-1, -1, -1},
			},
			{
				// left-node queue counts
				"partition_indices":      {4, 5, 6},
				"partition_state_values": {2, 2, 2},
			},
			{
				// upper-node queue counts
				"partition_indices":      {8, 9},
				"partition_state_values": {2, 2},
			},
			{
				// right-node queue counts
				"partition_indices":      {11, 12},
				"partition_state_values": {2, 2},
			},
			{
				// queue-left-node-middle-outside-triangle
				"upstream_partition":        {14},
				"upstream_value_index":      {0},
				"line_length":               {322.509871659},
				"spacecraft_length":         {2.0},
				"spacecraft_speed":          {50.0},
				"spacecraft_speed_variance": {5.0},
				"time_to_exit":              {10.0},
			},
			{
				// queue-left-node-upper-outside-triangle
				"upstream_partition":        {14},
				"upstream_value_index":      {1},
				"line_length":               {354.675108482},
				"spacecraft_length":         {2.0},
				"spacecraft_speed":          {50.0},
				"spacecraft_speed_variance": {5.0},
				"time_to_exit":              {10.0},
			},
			{
				// queue-left-node-lower-outside-triangle
				"upstream_partition":        {14},
				"upstream_value_index":      {2},
				"line_length":               {312.773010188},
				"spacecraft_length":         {2.0},
				"spacecraft_speed":          {50.0},
				"spacecraft_speed_variance": {5.0},
				"time_to_exit":              {10.0},
			},
			{
				// left-node
				"connected_incoming_partitions": {4, 5, 6},
			},
			{
				// queue-upper-node-outside-triangle
				"upstream_partition":        {14},
				"upstream_value_index":      {3},
				"line_length":               {469.583908959},
				"spacecraft_length":         {2.0},
				"spacecraft_speed":          {50.0},
				"spacecraft_speed_variance": {5.0},
				"time_to_exit":              {10.0},
			},
			{
				// queue-upper-node-inside-triangle
				"upstream_partition":        {7},
				"upstream_value_index":      {0},
				"line_length":               {253.799143859},
				"spacecraft_length":         {2.0},
				"spacecraft_speed":          {50.0},
				"spacecraft_speed_variance": {5.0},
				"time_to_exit":              {10.0},
			},
			{
				// upper-node
				"connected_incoming_partitions": {8, 9},
			},
			{
				// queue-right-node-lower-inside-triangle
				"upstream_partition":        {10},
				"upstream_value_index":      {0},
				"line_length":               {292.836844271},
				"spacecraft_length":         {2.0},
				"spacecraft_speed":          {50.0},
				"spacecraft_speed_variance": {5.0},
				"time_to_exit":              {10.0},
			},
			{
				// queue-right-node-upper-inside-triangle
				"upstream_partition":        {10},
				"upstream_value_index":      {1},
				"line_length":               {207.953361475},
				"spacecraft_length":         {2.0},
				"spacecraft_speed":          {50.0},
				"spacecraft_speed_variance": {5.0},
				"time_to_exit":              {10.0},
			},
			{
				// right-node
				"connected_incoming_partitions": {11, 12},
			},
			{
				// incoming counts
				"observed_values":                 {1, 1, 1, 1},
				"state_value_observation_indices": {0, 1, 2, 3},
				"state_value_observation_probs":   {0.9, 0.9, 0.9, 0.9},
			},
		},
		InitStateValues: [][]float64{
			{-1.0, -1.0, -1.0},
			{0.0, 0.0, 0.0},
			{0.0, 0.0},
			{0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 0.0},
			{0.0, 0.0, 0.0, 0.0},
		},
		InitTimeValue:         0.0,
		Seeds:                 seeds,
		StateWidths:           []int{3, 3, 2, 2, 5, 5, 5, 2, 5, 5, 3, 5, 5, 2, 4},
		StateHistoryDepths:    []int{1, 1, 1, 1, 500, 500, 500, 1, 500, 500, 1, 500, 500, 1, 1},
		TimestepsHistoryDepth: 500,
	}
	partitions := []simulator.Partition{
		{
			// action taker
			Iteration: &general.ParamValuesIteration{},
		},
		{
			// left-node queue counts
			Iteration: &general.CopyValuesIteration{},
		},
		{
			// upper-node queue counts
			Iteration: &general.CopyValuesIteration{},
		},
		{
			// right-node queue counts
			Iteration: &general.CopyValuesIteration{},
		},
		{
			// queue-left-node-middle-outside-triangle
			Iteration: &games.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {0},
			},
		},
		{
			// queue-left-node-upper-outside-triangle
			Iteration: &games.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {0},
			},
		},
		{
			// queue-left-node-lower-outside-triangle
			Iteration: &games.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {0},
			},
		},
		{
			// left-node
			Iteration: &games.SpacecraftLineConnectorIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_4_input_count": 4,
				"partition_5_input_count": 5,
				"partition_6_input_count": 6,
			},
			ParamsFromIndices: map[string][]int{
				"partition_4_input_count": {1},
				"partition_5_input_count": {1},
				"partition_6_input_count": {1},
			},
		},
		{
			// queue-upper-node-outside-triangle
			Iteration: &games.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {1},
			},
		},
		{
			// queue-upper-node-inside-triangle
			Iteration: &games.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {1},
			},
		},
		{
			// upper-node
			Iteration: &games.SpacecraftLineConnectorIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_8_input_count": 8,
				"partition_9_input_count": 9,
			},
			ParamsFromIndices: map[string][]int{
				"partition_8_input_count": {1},
				"partition_9_input_count": {1},
			},
		},
		{
			// queue-right-node-lower-inside-triangle
			Iteration: &games.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {2},
			},
		},
		{
			// queue-right-node-upper-inside-triangle
			Iteration: &games.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {2},
			},
		},
		{
			// right-node
			Iteration: &games.SpacecraftLineConnectorIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_11_input_count": 11,
				"partition_12_input_count": 12,
			},
			ParamsFromIndices: map[string][]int{
				"partition_11_input_count": {1},
				"partition_12_input_count": {1},
			},
		},
		{
			// incoming counts
			Iteration: &discrete.BinomialObservationProcessIteration{},
		},
	}
	for index, partition := range partitions {
		partition.Iteration.Configure(index, settings)
	}
	implementations := &simulator.Implementations{
		Partitions:      partitions,
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.StdoutOutputFunction{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{
			MaxTimeElapsed: 3652.5,
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 1.0},
	}
	simulator.NewPartitionCoordinator(settings, implementations).Run()
	// simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
