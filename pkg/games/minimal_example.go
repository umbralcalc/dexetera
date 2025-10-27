package games

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// MinimalExampleGame is the simplest possible game that demonstrates the framework.
// It just increments a counter and shows how easy it is to create new games.
type MinimalExampleGame struct {
	config *GameConfig
}

// NewMinimalExampleGame creates a new minimal example game using GameBuilder and VisualizationBuilder
func NewMinimalExampleGame() *MinimalExampleGame {
	// Create visualization using VisualizationBuilder
	visConfig := NewVisualizationBuilder().
		WithCanvas(400, 200).
		WithBackground("#2a2a2a").
		WithUpdateInterval(100).
		AddText("counter_state", "Count: {value}", 200, 100, &TextOptions{
			FontSize: 24,
			Color:    "#ffffff",
		}).
		Build()

	// Create the game using the fluent GameBuilder API
	config := NewGameBuilder("minimal_example").
		WithDescription("The simplest possible game - just a counter").
		WithPartition("counter", "counter_state", &MinimalCounterIteration{}).
		WithServerPartition("counter_state").
		WithParameter("counter_init", []float64{0.0}).
		WithParameter("counter_params", map[string][]float64{"increment": {1.0}}).
		WithMaxTime(30.0).
		WithTimestep(1.0).
		WithVisualization(visConfig).
		Build()

	return &MinimalExampleGame{config: config}
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
func (m *MinimalExampleGame) GetConfig() *GameConfig {
	return m.config
}

// GetConfigGenerator returns a configured ConfigGenerator that builds the simulation
// configuration step-by-step using the fluent API
func (m *MinimalExampleGame) GetConfigGenerator() *simulator.ConfigGenerator {
	initialValue := m.config.Parameters["initial_value"].(int)
	increment := m.config.Parameters["increment"].(int)

	// Create a new ConfigGenerator
	configGen := simulator.NewConfigGenerator()

	// Set global seed
	configGen.SetGlobalSeed(42)

	// Configure the counter partition
	counterPartition := &simulator.PartitionConfig{
		Name:            "counter_state",
		Params:          simulator.NewParams(map[string][]float64{"increment": {float64(increment)}}),
		InitStateValues: []float64{float64(initialValue)},
	}
	configGen.SetPartition(counterPartition)

	// Configure simulation-level settings
	simulationConfig := &simulator.SimulationConfig{
		InitTimeValue: 0.0,
	}
	configGen.SetSimulation(simulationConfig)

	return configGen
}

// GetRenderer returns the visualization renderer
func (m *MinimalExampleGame) GetRenderer() GameRenderer {
	return &GenericRenderer{config: m.config.VisualizationConfig}
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
	outputState := stateHistories[partitionIndex].Values.RawRowView(0)
	increment := params.Get("increment")[0]

	// Simple counter: just increment the value
	outputState[0] += increment

	return outputState
}
