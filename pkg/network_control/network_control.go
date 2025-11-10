package network_control

import (
	"math"
	"math/rand"

	"github.com/umbralcalc/dexetera/pkg/game"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// NetworkControlGame simulates traffic across a small network where the player controls
// junction phases to maximise throughput.
type NetworkControlGame struct {
	config *game.GameConfig
}

// NewNetworkControlGame constructs a new network control game instance with visualization
// and simulation wired in.
func NewNetworkControlGame() *NetworkControlGame {
	network := defaultNetworkDefinition()

	visBuilder := game.NewVisualizationBuilder().
		WithCanvas(800, 600).
		WithBackground("#1a1a1a").
		WithUpdateInterval(0)

	// Static road layout
	for _, edge := range network.Edges {
		visBuilder.AddLine(
			"",
			int(edge.StartX),
			int(edge.StartY),
			int(edge.EndX),
			int(edge.EndY),
			&game.LineOptions{
				Color: "#3d3d3d",
				Width: 14,
			},
		)
	}

	// Dynamic vehicle rectangles
	visBuilder.AddRectangleSet("vehicle_rectangles", 18, 10, &game.ShapeOptions{
		FillColor:   "#f4a261",
		StrokeColor: "#e76f51",
		StrokeWidth: 1,
	})

	// Junction state readouts
	visBuilder.AddText("junction_a_control", "Junction A Phase: {value}", 220, 90, &game.TextOptions{
		FontSize:   14,
		Color:      "#ffffff",
		FontFamily: "Arial",
		TextAlign:  "left",
	})
	visBuilder.AddText("junction_b_control", "Junction B Phase: {value}", 520, 90, &game.TextOptions{
		FontSize:   14,
		Color:      "#ffffff",
		FontFamily: "Arial",
		TextAlign:  "left",
	})

	// Throughput metrics
	visBuilder.AddText("flow_metrics", "Vehicles exited: {value}", 600, 520, &game.TextOptions{
		FontSize:   16,
		Color:      "#ffffff",
		FontFamily: "Arial",
		TextAlign:  "center",
	})

	config := game.NewGameBuilder("network_control").
		WithDescription("Control junction phases to maximise traffic throughput.").
		WithServerPartition("edge_states").
		WithServerPartition("vehicle_rectangles").
		WithServerPartition("flow_metrics").
		WithServerPartition("junction_a_control").
		WithServerPartition("junction_b_control").
		WithActionStatePartition("junction_a_control").
		WithActionStatePartition("junction_b_control").
		WithVisualization(visBuilder.Build()).
		WithSimulation(BuildNetworkControlSimulation).
		Build()

	return &NetworkControlGame{config: config}
}

// GetName returns the game name.
func (g *NetworkControlGame) GetName() string {
	return g.config.Name
}

// GetDescription returns the game description.
func (g *NetworkControlGame) GetDescription() string {
	return g.config.Description
}

// GetConfig exposes the game configuration.
func (g *NetworkControlGame) GetConfig() *game.GameConfig {
	return g.config
}

// GetRenderer returns a renderer built from the visualization config.
func (g *NetworkControlGame) GetRenderer() game.GameRenderer {
	return &game.GenericRenderer{Config: g.config.VisualizationConfig}
}

// BuildNetworkControlSimulation constructs the simulator configuration generator.
func BuildNetworkControlSimulation() *simulator.ConfigGenerator {
	def := defaultNetworkDefinition()

	gen := simulator.NewConfigGenerator()

	totalSlots := 0
	edgeOffsets := make([]int, len(def.Edges))
	for i, edge := range def.Edges {
		edgeOffsets[i] = totalSlots
		totalSlots += edge.Capacity
	}
	exitCounterIndex := totalSlots

	edgeInit := make([]float64, totalSlots+1)
	for i := 0; i < totalSlots; i++ {
		edgeInit[i] = -1.0
	}

	vehicleRectInit := make([]float64, totalSlots*4)

	// Junction control partitions must be registered first to ensure indices are stable.
	junctionAIndex := 0
	junctionBIndex := 1
	edgeStatesIndex := 2

	junctionA := &simulator.PartitionConfig{
		Name: "junction_a_control",
		Iteration: &JunctionControlIteration{
			PhaseCount:  len(def.Junctions[0].Phases),
			ActionIndex: 0,
		},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {0.0},
		}),
		InitStateValues:   []float64{0.0},
		StateHistoryDepth: 2,
		Seed:              1001,
	}
	gen.SetPartition(junctionA)

	junctionB := &simulator.PartitionConfig{
		Name: "junction_b_control",
		Iteration: &JunctionControlIteration{
			PhaseCount:  len(def.Junctions[1].Phases),
			ActionIndex: 1,
		},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {1.0},
		}),
		InitStateValues:   []float64{1.0},
		StateHistoryDepth: 2,
		Seed:              1002,
	}
	gen.SetPartition(junctionB)

	junctionPartitions := map[int]int{
		def.Junctions[0].ID: junctionAIndex,
		def.Junctions[1].ID: junctionBIndex,
	}

	edgeStates := &simulator.PartitionConfig{
		Name: "edge_states",
		Iteration: NewNetworkFlowIteration(
			def,
			edgeOffsets,
			exitCounterIndex,
			junctionPartitions,
			4321,
		),
		InitStateValues:   edgeInit,
		StateHistoryDepth: 2,
		Seed:              1901,
	}
	gen.SetPartition(edgeStates)

	vehicleRectangles := &simulator.PartitionConfig{
		Name: "vehicle_rectangles",
		Iteration: NewVehicleRectanglesIteration(
			def,
			edgeStatesIndex,
			edgeOffsets,
			totalSlots,
			18.0,
			10.0,
		),
		InitStateValues:   vehicleRectInit,
		StateHistoryDepth: 1,
		Seed:              1902,
	}
	gen.SetPartition(vehicleRectangles)

	flowMetrics := &simulator.PartitionConfig{
		Name: "flow_metrics",
		Iteration: NewFlowMetricsIteration(
			def,
			edgeStatesIndex,
			edgeOffsets,
			exitCounterIndex,
		),
		InitStateValues:   []float64{0.0, 0.0, 0.0},
		StateHistoryDepth: 1,
		Seed:              1903,
	}
	gen.SetPartition(flowMetrics)

	sim := &simulator.SimulationConfig{
		OutputCondition:      &simulator.EveryStepOutputCondition{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{MaxTimeElapsed: 10000.0},
		TimestepFunction:     &simulator.ConstantTimestepFunction{Stepsize: 1.0},
		InitTimeValue:        0.0,
	}
	gen.SetSimulation(sim)

	return gen
}

// Network definition and supporting iteration implementations.

type networkDefinition struct {
	Edges     []EdgeConfig
	Junctions []JunctionConfig
}

type EdgeConfig struct {
	ID            int
	Name          string
	StartX        float64
	StartY        float64
	EndX          float64
	EndY          float64
	Length        float64
	Capacity      int
	Speed         float64
	SpawnRate     float64
	NextEdge      int
	JunctionIndex int
}

type JunctionConfig struct {
	ID        int
	Name      string
	Phases    []JunctionPhase
	PositionX float64
	PositionY float64
}

type JunctionPhase struct {
	Label            string
	AllowedIncoming  map[int]bool
	DefaultSelection bool
}

func defaultNetworkDefinition() networkDefinition {
	edges := []EdgeConfig{
		{
			ID:            0,
			Name:          "west_entry",
			StartX:        80.0,
			StartY:        280.0,
			EndX:          300.0,
			EndY:          280.0,
			Capacity:      8,
			Speed:         14.0,
			SpawnRate:     0.55,
			NextEdge:      2,
			JunctionIndex: 0,
		},
		{
			ID:            1,
			Name:          "south_entry",
			StartX:        300.0,
			StartY:        520.0,
			EndX:          300.0,
			EndY:          300.0,
			Capacity:      7,
			Speed:         12.0,
			SpawnRate:     0.45,
			NextEdge:      3,
			JunctionIndex: 0,
		},
		{
			ID:            2,
			Name:          "junction_a_to_b",
			StartX:        300.0,
			StartY:        280.0,
			EndX:          500.0,
			EndY:          280.0,
			Capacity:      7,
			Speed:         13.0,
			NextEdge:      5,
			JunctionIndex: 1,
		},
		{
			ID:            3,
			Name:          "junction_a_to_north_exit",
			StartX:        300.0,
			StartY:        280.0,
			EndX:          300.0,
			EndY:          80.0,
			Capacity:      6,
			Speed:         12.0,
			NextEdge:      -1,
			JunctionIndex: -1,
		},
		{
			ID:            4,
			Name:          "north_entry",
			StartX:        500.0,
			StartY:        80.0,
			EndX:          500.0,
			EndY:          260.0,
			Capacity:      6,
			Speed:         11.0,
			SpawnRate:     0.35,
			NextEdge:      6,
			JunctionIndex: 1,
		},
		{
			ID:            5,
			Name:          "junction_b_to_east_exit",
			StartX:        500.0,
			StartY:        280.0,
			EndX:          720.0,
			EndY:          280.0,
			Capacity:      8,
			Speed:         14.0,
			NextEdge:      -1,
			JunctionIndex: -1,
		},
		{
			ID:            6,
			Name:          "junction_b_to_south_exit",
			StartX:        500.0,
			StartY:        280.0,
			EndX:          500.0,
			EndY:          520.0,
			Capacity:      7,
			Speed:         12.0,
			NextEdge:      -1,
			JunctionIndex: -1,
		},
	}

	for i := range edges {
		dx := edges[i].EndX - edges[i].StartX
		dy := edges[i].EndY - edges[i].StartY
		edges[i].Length = math.Hypot(dx, dy)
	}

	junctions := []JunctionConfig{
		{
			ID:        0,
			Name:      "junction_a",
			PositionX: 300.0,
			PositionY: 280.0,
			Phases: []JunctionPhase{
				{
					Label:           "West → East",
					AllowedIncoming: map[int]bool{0: true},
				},
				{
					Label:           "South → North",
					AllowedIncoming: map[int]bool{1: true},
				},
			},
		},
		{
			ID:        1,
			Name:      "junction_b",
			PositionX: 500.0,
			PositionY: 280.0,
			Phases: []JunctionPhase{
				{
					Label:           "Through East",
					AllowedIncoming: map[int]bool{2: true},
				},
				{
					Label:           "North → South",
					AllowedIncoming: map[int]bool{4: true},
				},
			},
		},
	}

	return networkDefinition{
		Edges:     edges,
		Junctions: junctions,
	}
}

// JunctionControlIteration simply copies the action state value into the partition state.
type JunctionControlIteration struct {
	PhaseCount  int
	ActionIndex int
}

func (j *JunctionControlIteration) Configure(partitionIndex int, settings *simulator.Settings) {
	// No configuration required.
}

func (j *JunctionControlIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	state := stateHistories[partitionIndex].CopyStateRow(0)
	actionValues := params.Get("action_state_values")
	if len(actionValues) > 0 {
		index := j.ActionIndex
		if index < 0 {
			index = 0
		}
		if index >= len(actionValues) {
			index = len(actionValues) - 1
		}
		if index < 0 || index >= len(actionValues) {
			state[0] = 0.0
			return state
		}

		phase := int(math.Round(actionValues[index]))
		if j.PhaseCount > 0 {
			phase = ((phase % j.PhaseCount) + j.PhaseCount) % j.PhaseCount
		} else {
			phase = 0
		}
		state[0] = float64(phase)
	}
	return state
}

// NetworkFlowIteration handles vehicle spawning, movement, and transitions between edges.
type NetworkFlowIteration struct {
	edges               []EdgeConfig
	edgeOffsets         []int
	edgeCapacities      []int
	edgeLengths         []float64
	edgeSpeeds          []float64
	edgeNext            []int
	edgeSpawnRates      []float64
	exitCounterIndex    int
	junctionAllowed     map[int][]map[int]bool
	junctionPartitions  map[int]int
	junctionPhaseCounts map[int]int
	junctionOrder       []int
	minSpacing          float64
	seed                int64
	rng                 *rand.Rand
}

func NewNetworkFlowIteration(
	def networkDefinition,
	edgeOffsets []int,
	exitCounterIndex int,
	junctionPartitions map[int]int,
	seed int64,
) *NetworkFlowIteration {
	edgeCapacities := make([]int, len(def.Edges))
	edgeLengths := make([]float64, len(def.Edges))
	edgeSpeeds := make([]float64, len(def.Edges))
	edgeNext := make([]int, len(def.Edges))
	edgeSpawnRates := make([]float64, len(def.Edges))

	for i, edge := range def.Edges {
		edgeCapacities[i] = edge.Capacity
		edgeLengths[i] = edge.Length
		edgeSpeeds[i] = edge.Speed
		edgeNext[i] = edge.NextEdge
		edgeSpawnRates[i] = edge.SpawnRate
	}

	junctionAllowed := make(map[int][]map[int]bool, len(def.Junctions))
	junctionPhaseCounts := make(map[int]int, len(def.Junctions))
	junctionOrder := make([]int, 0, len(def.Junctions))

	for _, junction := range def.Junctions {
		phaseMaps := make([]map[int]bool, len(junction.Phases))
		for idx, phase := range junction.Phases {
			allowedCopy := make(map[int]bool, len(phase.AllowedIncoming))
			for incoming, ok := range phase.AllowedIncoming {
				allowedCopy[incoming] = ok
			}
			phaseMaps[idx] = allowedCopy
		}
		junctionAllowed[junction.ID] = phaseMaps
		junctionPhaseCounts[junction.ID] = len(junction.Phases)
		junctionOrder = append(junctionOrder, junction.ID)
	}

	return &NetworkFlowIteration{
		edges:               def.Edges,
		edgeOffsets:         append([]int(nil), edgeOffsets...),
		edgeCapacities:      edgeCapacities,
		edgeLengths:         edgeLengths,
		edgeSpeeds:          edgeSpeeds,
		edgeNext:            edgeNext,
		edgeSpawnRates:      edgeSpawnRates,
		exitCounterIndex:    exitCounterIndex,
		junctionAllowed:     junctionAllowed,
		junctionPartitions:  junctionPartitions,
		junctionPhaseCounts: junctionPhaseCounts,
		junctionOrder:       junctionOrder,
		minSpacing:          6.0,
		seed:                seed,
	}
}

func (n *NetworkFlowIteration) Configure(partitionIndex int, settings *simulator.Settings) {
	n.rng = rand.New(rand.NewSource(n.seed))
}

func (n *NetworkFlowIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	state := stateHistories[partitionIndex].CopyStateRow(0)
	exitCount := state[n.exitCounterIndex]

	// Read junction phases
	junctionPhases := make(map[int]int, len(n.junctionPartitions))
	for _, junctionID := range n.junctionOrder {
		partitionIdx, ok := n.junctionPartitions[junctionID]
		if !ok || partitionIdx >= len(stateHistories) {
			continue
		}
		phaseState := stateHistories[partitionIdx].CopyStateRow(0)
		phase := 0
		if len(phaseState) > 0 {
			phase = int(math.Round(phaseState[0]))
		}
		phaseCount := n.junctionPhaseCounts[junctionID]
		if phaseCount > 0 {
			phase = ((phase % phaseCount) + phaseCount) % phaseCount
		}
		junctionPhases[junctionID] = phase
	}

	// Spawn new vehicles on source edges.
	for edgeIdx, rate := range n.edgeSpawnRates {
		if rate <= 0 {
			continue
		}
		if n.rng.Float64() < rate {
			n.insertVehicle(state, edgeIdx)
		}
	}

	// Move vehicles along each edge.
	for edgeIdx := range n.edges {
		offset := n.edgeOffsets[edgeIdx]
		capacity := n.edgeCapacities[edgeIdx]
		length := n.edgeLengths[edgeIdx]
		speed := n.edgeSpeeds[edgeIdx]

		if capacity == 0 {
			continue
		}

		for slot := 0; slot < capacity; slot++ {
			idx := offset + slot
			position := state[idx]
			if position < 0 {
				continue
			}

			newPosition := position + speed
			if newPosition > length {
				newPosition = length
			}

			if slot > 0 {
				aheadIdx := offset + slot - 1
				aheadPos := state[aheadIdx]
				if aheadPos >= 0 {
					maxPos := aheadPos - n.minSpacing
					if maxPos < 0 {
						maxPos = 0
					}
					if newPosition > maxPos {
						newPosition = math.Max(position, maxPos)
					}
				}
			}

			state[idx] = newPosition
		}
	}

	// Transfer vehicles between edges or to exits.
	for edgeIdx := range n.edges {
		length := n.edgeLengths[edgeIdx]
		nextEdge := n.edgeNext[edgeIdx]
		junctionIndex := n.edges[edgeIdx].JunctionIndex

		for {
			offset := n.edgeOffsets[edgeIdx]
			if n.edgeCapacities[edgeIdx] == 0 {
				break
			}

			front := state[offset]
			if front < 0 || front < length-1e-3 {
				break
			}

			if nextEdge == -1 {
				exitCount++
				n.removeFrontVehicle(state, edgeIdx)
				state[n.exitCounterIndex] = exitCount
				continue
			}

			allowed := true
			if junctionIndex >= 0 {
				phase := junctionPhases[junctionIndex]
				allowedPhases := n.junctionAllowed[junctionIndex]
				if phase >= 0 && phase < len(allowedPhases) {
					if !allowedPhases[phase][edgeIdx] {
						allowed = false
					}
				}
			}

			if !allowed {
				state[offset] = length
				break
			}

			if !n.insertVehicle(state, nextEdge) {
				state[offset] = length
				break
			}

			n.removeFrontVehicle(state, edgeIdx)
		}
	}

	state[n.exitCounterIndex] = exitCount
	return state
}

func (n *NetworkFlowIteration) insertVehicle(state []float64, edgeIdx int) bool {
	offset := n.edgeOffsets[edgeIdx]
	capacity := n.edgeCapacities[edgeIdx]
	for slot := capacity - 1; slot >= 0; slot-- {
		idx := offset + slot
		if state[idx] < 0 {
			state[idx] = 0.0
			return true
		}
	}
	return false
}

func (n *NetworkFlowIteration) removeFrontVehicle(state []float64, edgeIdx int) {
	offset := n.edgeOffsets[edgeIdx]
	capacity := n.edgeCapacities[edgeIdx]
	if capacity == 0 {
		return
	}
	for slot := 0; slot < capacity-1; slot++ {
		state[offset+slot] = state[offset+slot+1]
	}
	state[offset+capacity-1] = -1.0
}

// VehicleRectanglesIteration projects vehicle positions into rectangle shapes for rendering.
type VehicleRectanglesIteration struct {
	edgeStatesIndex int
	edgeOffsets     []int
	edgeLengths     []float64
	edges           []EdgeConfig
	totalSlots      int
	rectangleWidth  float64
	rectangleHeight float64
}

func NewVehicleRectanglesIteration(
	def networkDefinition,
	edgeStatesIndex int,
	edgeOffsets []int,
	totalSlots int,
	rectWidth float64,
	rectHeight float64,
) *VehicleRectanglesIteration {
	edgeLengths := make([]float64, len(def.Edges))
	for i, edge := range def.Edges {
		edgeLengths[i] = edge.Length
	}

	return &VehicleRectanglesIteration{
		edgeStatesIndex: edgeStatesIndex,
		edgeOffsets:     append([]int(nil), edgeOffsets...),
		edgeLengths:     edgeLengths,
		edges:           def.Edges,
		totalSlots:      totalSlots,
		rectangleWidth:  rectWidth,
		rectangleHeight: rectHeight,
	}
}

func (v *VehicleRectanglesIteration) Configure(partitionIndex int, settings *simulator.Settings) {
	// No configuration required.
}

func (v *VehicleRectanglesIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	output := stateHistories[partitionIndex].CopyStateRow(0)
	for i := range output {
		output[i] = 0.0
	}

	if v.edgeStatesIndex >= len(stateHistories) {
		return output
	}

	edgeState := stateHistories[v.edgeStatesIndex].CopyStateRow(0)

	rectIdx := 0
	for edgeIdx, edge := range v.edges {
		offset := v.edgeOffsets[edgeIdx]
		length := v.edgeLengths[edgeIdx]
		if length <= 0 {
			continue
		}

		for slot := 0; slot < edge.Capacity; slot++ {
			if rectIdx >= v.totalSlots {
				break
			}

			stateIndex := offset + slot
			position := edgeState[stateIndex]
			base := rectIdx * 4

			if position < 0 {
				output[base] = 0.0
				output[base+1] = 0.0
				output[base+2] = 0.0
				output[base+3] = 0.0
				rectIdx++
				continue
			}

			t := position / length
			if t < 0 {
				t = 0
			}
			if t > 1 {
				t = 1
			}

			x := edge.StartX + (edge.EndX-edge.StartX)*t
			y := edge.StartY + (edge.EndY-edge.StartY)*t

			output[base] = x
			output[base+1] = y
			output[base+2] = v.rectangleWidth
			output[base+3] = v.rectangleHeight

			rectIdx++
		}
	}

	return output
}

// FlowMetricsIteration computes aggregate metrics for server consumption and UI.
type FlowMetricsIteration struct {
	edgeStatesIndex  int
	edgeOffsets      []int
	edgeCapacities   []int
	exitCounterIndex int
}

func NewFlowMetricsIteration(
	def networkDefinition,
	edgeStatesIndex int,
	edgeOffsets []int,
	exitCounterIndex int,
) *FlowMetricsIteration {
	capacities := make([]int, len(def.Edges))
	for i, edge := range def.Edges {
		capacities[i] = edge.Capacity
	}

	return &FlowMetricsIteration{
		edgeStatesIndex:  edgeStatesIndex,
		edgeOffsets:      append([]int(nil), edgeOffsets...),
		edgeCapacities:   capacities,
		exitCounterIndex: exitCounterIndex,
	}
}

func (f *FlowMetricsIteration) Configure(partitionIndex int, settings *simulator.Settings) {
	// No configuration required.
}

func (f *FlowMetricsIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	state := stateHistories[partitionIndex].CopyStateRow(0)
	for i := range state {
		state[i] = 0.0
	}

	if f.edgeStatesIndex >= len(stateHistories) {
		return state
	}

	edgeState := stateHistories[f.edgeStatesIndex].CopyStateRow(0)

	totalVehicles := 0.0
	totalCapacity := 0.0

	for edgeIdx := range f.edgeOffsets {
		offset := f.edgeOffsets[edgeIdx]
		capacity := f.edgeCapacities[edgeIdx]
		totalCapacity += float64(capacity)
		for slot := 0; slot < capacity; slot++ {
			if edgeState[offset+slot] >= 0 {
				totalVehicles++
			}
		}
	}

	exitCount := 0.0
	if f.exitCounterIndex < len(edgeState) {
		exitCount = edgeState[f.exitCounterIndex]
	}

	averageOccupancy := 0.0
	if totalCapacity > 0 {
		averageOccupancy = totalVehicles / totalCapacity
	}

	state[0] = exitCount
	state[1] = totalVehicles
	state[2] = averageOccupancy

	return state
}
