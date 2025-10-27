package games

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// BuilderExampleGame demonstrates the GameBuilder pattern
// This creates a more complex game with multiple partitions and custom visualizations
type BuilderExampleGame struct {
	config *GameConfig
}

// NewBuilderExampleGame creates a new game using the GameBuilder pattern
func NewBuilderExampleGame() *BuilderExampleGame {
	// Create a complex game using the fluent API
	config := NewGameBuilder("builder_example").
		WithDescription("A complex game demonstrating the GameBuilder pattern").
		WithPartition("counter", "counter_state", &CounterIteration{}).
		WithPartition("timer", "timer_state", &TimerIteration{}).
		WithServerPartition("counter_state").
		WithServerPartition("timer_state").
		WithParameter("counter_init", []float64{0.0}).
		WithParameter("counter_params", map[string][]float64{"increment": {1.0}}).
		WithParameter("timer_init", []float64{0.0}).
		WithParameter("timer_params", map[string][]float64{"speed": {0.1}}).
		WithMaxTime(60.0).
		WithTimestep(1.0).
		WithVisualization(&VisualizationConfig{
			CanvasWidth:      600,
			CanvasHeight:     400,
			BackgroundColor:  "#1a1a1a",
			UpdateIntervalMs: 50,
			Renderers: []RendererConfig{
				{
					Type:          "text",
					PartitionName: "counter_state",
					Properties: map[string]interface{}{
						"fontSize": 24,
						"color":    "#00ff00",
						"x":        150,
						"y":        100,
						"text":     "Counter: {value}",
					},
				},
				{
					Type:          "text",
					PartitionName: "timer_state",
					Properties: map[string]interface{}{
						"fontSize": 24,
						"color":    "#ff0000",
						"x":        450,
						"y":        100,
						"text":     "Timer: {value}",
					},
				},
				{
					Type:          "circle",
					PartitionName: "counter_state",
					Properties: map[string]interface{}{
						"color":  "#00ff00",
						"x":      150,
						"y":      200,
						"radius": 20,
					},
				},
				{
					Type:          "rectangle",
					PartitionName: "timer_state",
					Properties: map[string]interface{}{
						"color":  "#ff0000",
						"x":      400,
						"y":      180,
						"width":  100,
						"height": 40,
					},
				},
			},
		}).
		Build()

	return &BuilderExampleGame{config: config}
}

// GetName returns the game name
func (b *BuilderExampleGame) GetName() string {
	return b.config.Name
}

// GetDescription returns the game description
func (b *BuilderExampleGame) GetDescription() string {
	return b.config.Description
}

// GetConfig returns the game configuration
func (b *BuilderExampleGame) GetConfig() *GameConfig {
	return b.config
}

// GetConfigGenerator returns a configured ConfigGenerator
func (b *BuilderExampleGame) GetConfigGenerator() *simulator.ConfigGenerator {
	configGen := simulator.NewConfigGenerator()
	configGen.SetGlobalSeed(42)

	// Add counter partition
	counterPartition := &simulator.PartitionConfig{
		Name:            "counter_state",
		Params:          simulator.NewParams(map[string][]float64{"increment": {1.0}}),
		InitStateValues: []float64{0.0},
	}
	configGen.SetPartition(counterPartition)

	// Add timer partition
	timerPartition := &simulator.PartitionConfig{
		Name:            "timer_state",
		Params:          simulator.NewParams(map[string][]float64{"speed": {0.1}}),
		InitStateValues: []float64{0.0},
	}
	configGen.SetPartition(timerPartition)

	// Configure simulation-level settings
	simulationConfig := &simulator.SimulationConfig{
		InitTimeValue: 0.0,
	}
	configGen.SetSimulation(simulationConfig)

	return configGen
}

// GetRenderer returns the visualization renderer
func (b *BuilderExampleGame) GetRenderer() GameRenderer {
	return &GenericRenderer{config: b.config.VisualizationConfig}
}

// CounterIteration implements a simple counter
type CounterIteration struct{}

func (c *CounterIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	// No configuration needed
}

func (c *CounterIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	outputState := stateHistories[partitionIndex].Values.RawRowView(0)
	increment := params.Get("increment")[0]
	outputState[0] += increment
	return outputState
}

// TimerIteration implements a timer that counts up
type TimerIteration struct{}

func (t *TimerIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	// No configuration needed
}

func (t *TimerIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	outputState := stateHistories[partitionIndex].Values.RawRowView(0)
	speed := params.Get("speed")[0]
	outputState[0] += speed
	return outputState
}
