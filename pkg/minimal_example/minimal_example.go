package minimal_example

import (
	"github.com/umbralcalc/dexetera/pkg/game"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// MinimalExampleGame is the simplest possible game that demonstrates the framework.
// It just increments a counter and shows how easy it is to create new games.
type MinimalExampleGame struct {
	config *game.GameConfig
}

// NewMinimalExampleGame creates a new minimal example game using GameBuilder and VisualizationBuilder
func NewMinimalExampleGame() *MinimalExampleGame {
	// Create visualization using VisualizationBuilder
	visConfig := game.NewVisualizationBuilder().
		WithCanvas(400, 200).
		WithBackground("#2a2a2a").
		WithUpdateInterval(100).
		AddText("counter_state", "Count: {value}", 200, 100, &game.TextOptions{
			FontSize: 24,
			Color:    "#ffffff",
		}).
		Build()

	// Create the game using the fluent GameBuilder API
	config := game.NewGameBuilder("minimal_example").
		WithDescription("The simplest possible game - just a counter").
		WithServerPartition("counter_state").
		WithActionStatePartition("counter_state").
		WithVisualization(visConfig).
		WithSimulation(BuildMinimalSimulation).
		Build()

	return &MinimalExampleGame{config: config}
}

// BuildMinimalSimulation produces the simulation config generator used by the framework
func BuildMinimalSimulation() *simulator.ConfigGenerator {
	gen := simulator.NewConfigGenerator()

	counter := &simulator.PartitionConfig{
		Name:      "counter_state",
		Iteration: &MinimalCounterIteration{},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {1.0},
		}),
		InitStateValues:   []float64{0.0},
		StateHistoryDepth: 1,
		Seed:              123,
	}
	gen.SetPartition(counter)

	sim := &simulator.SimulationConfig{
		// Output callback wired by the framework (JS callback OutputFunction)
		OutputCondition:      &simulator.EveryStepOutputCondition{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{MaxTimeElapsed: 31.0},
		TimestepFunction:     &simulator.ConstantTimestepFunction{Stepsize: 1.0},
		InitTimeValue:        0.0,
	}
	gen.SetSimulation(sim)
	return gen
}

// GetName returns the game name
func (m *MinimalExampleGame) GetName() string {
	return m.config.Name
}

// GetDescription returns the game description
func (m *MinimalExampleGame) GetDescription() string {
	return m.config.Description
}

// GetConfig returns the game configuration
func (m *MinimalExampleGame) GetConfig() *game.GameConfig {
	return m.config
}

// GetRenderer returns the visualization renderer
func (m *MinimalExampleGame) GetRenderer() game.GameRenderer {
	return &game.GenericRenderer{Config: m.config.VisualizationConfig}
}

// MinimalCounterIteration implements the simplest possible counter
type MinimalCounterIteration struct{}

func (m *MinimalCounterIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	// No configuration needed for this simple counter
}

func (m *MinimalCounterIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	outputState := stateHistories[partitionIndex].CopyStateRow(0)

	// Simple counter: get the output "action_state_values" from the ActionState
	// sent by the Python server
	outputState[0] = params.Get("action_state_values")[0]

	return outputState
}
