//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/examples"
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/stochadex/pkg/observations"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func main() {
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
				"upstream_partition":   {14},
				"upstream_value_index": {0},
				// "line_length":               {15.0},
				// "spacecraft_length":         {1.0},
				// "spacecraft_speed":          {3.5},
				// "spacecraft_speed_variance": {3.5},
				// "partition_flow_allowed":    {0},
				// "time_to_exit":              {0.0},
			},
			{
				// queue-left-node-upper-outside-triangle
				"upstream_partition":   {14},
				"upstream_value_index": {1},
				// "line_length":               {15.0},
				// "spacecraft_length":         {1.0},
				// "spacecraft_speed":          {3.5},
				// "spacecraft_speed_variance": {3.5},
				// "partition_flow_allowed":    {0},
				// "time_to_exit":              {0.0},
			},
			{
				// queue-left-node-lower-outside-triangle
				"upstream_partition":   {14},
				"upstream_value_index": {2},
				// "line_length":               {15.0},
				// "spacecraft_length":         {1.0},
				// "spacecraft_speed":          {3.5},
				// "spacecraft_speed_variance": {3.5},
				// "partition_flow_allowed":    {0},
				// "time_to_exit":              {0.0},
			},
			{
				// left-node
				"connected_incoming_partitions": {4, 5, 6},
			},
			{
				// queue-upper-node-outside-triangle
				"upstream_partition":   {14},
				"upstream_value_index": {3},
				// "line_length":               {15.0},
				// "spacecraft_length":         {1.0},
				// "spacecraft_speed":          {3.5},
				// "spacecraft_speed_variance": {3.5},
				// "partition_flow_allowed":    {0},
				// "time_to_exit":              {0.0},
			},
			{
				// queue-upper-node-inside-triangle
				"upstream_partition":   {7},
				"upstream_value_index": {0},
				// "line_length":               {15.0},
				// "spacecraft_length":         {1.0},
				// "spacecraft_speed":          {3.5},
				// "spacecraft_speed_variance": {3.5},
				// "partition_flow_allowed":    {0},
				// "time_to_exit":              {0.0},
			},
			{
				// upper-node
				"connected_incoming_partitions": {8, 9},
			},
			{
				// queue-right-node-lower-inside-triangle
				"upstream_partition":   {10},
				"upstream_value_index": {0},
				// "line_length":               {15.0},
				// "spacecraft_length":         {1.0},
				// "spacecraft_speed":          {3.5},
				// "spacecraft_speed_variance": {3.5},
				// "partition_flow_allowed":    {0},
				// "time_to_exit":              {0.0},
			},
			{
				// queue-right-node-upper-inside-triangle
				"upstream_partition":   {10},
				"upstream_value_index": {1},
				// "line_length":               {15.0},
				// "spacecraft_length":         {1.0},
				// "spacecraft_speed":          {3.5},
				// "spacecraft_speed_variance": {3.5},
				// "partition_flow_allowed":    {0},
				// "time_to_exit":              {0.0},
			},
			{
				// right-node
				"connected_incoming_partitions": {11, 12},
			},
			{
				// incoming counts
				"observed_values":                 {1, 1, 1, 1},
				"state_value_observation_indices": {0, 1, 2, 3},
				"state_value_observation_probs":   {0.5, 0.5, 0.5, 0.5},
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
		Seeds:                 []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		StateWidths:           []int{3, 3, 2, 2, 5, 5, 5, 2, 5, 5, 3, 5, 5, 2, 4},
		StateHistoryDepths:    []int{1, 1, 1, 1, 150, 150, 150, 1, 150, 150, 1, 150, 150, 1, 1},
		TimestepsHistoryDepth: 150,
	}
	partitions := []simulator.Partition{
		{
			// action taker
			Iteration: &simulator.ParamValuesIteration{},
		},
		{
			// left-node queue counts
			Iteration: &simulator.CopyValuesIteration{},
		},
		{
			// upper-node queue counts
			Iteration: &simulator.CopyValuesIteration{},
		},
		{
			// right-node queue counts
			Iteration: &simulator.CopyValuesIteration{},
		},
		{
			// queue-left-node-middle-outside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {0},
			},
		},
		{
			// queue-left-node-upper-outside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {0},
			},
		},
		{
			// queue-left-node-lower-outside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {0},
			},
		},
		{
			// left-node
			Iteration: &examples.SpacecraftLineConnectorIteration{},
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
			Iteration: &examples.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {1},
			},
		},
		{
			// queue-upper-node-inside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {1},
			},
		},
		{
			// upper-node
			Iteration: &examples.SpacecraftLineConnectorIteration{},
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
			Iteration: &examples.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {2},
			},
		},
		{
			// queue-right-node-upper-inside-triangle
			Iteration: &examples.SpacecraftLineCountIteration{},
			ParamsFromUpstreamPartition: map[string]int{
				"partition_flow_allowed": 0,
			},
			ParamsFromIndices: map[string][]int{
				"partition_flow_allowed": {2},
			},
		},
		{
			// right-node
			Iteration: &examples.SpacecraftLineConnectorIteration{},
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
			Iteration: &observations.BinomialStaticPartialStateObservationIteration{},
		},
	}
	for index, partition := range partitions {
		partition.Iteration.Configure(index, settings)
	}
	implementations := &simulator.Implementations{
		Partitions:      partitions,
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{
			MaxTimeElapsed: 3652.5,
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 10.0},
	}
	simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
