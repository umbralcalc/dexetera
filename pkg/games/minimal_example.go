package games

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// MinimalExampleGame is the simplest possible game that demonstrates the framework.
// It just increments a counter and shows how easy it is to create new games.
type MinimalExampleGame struct {
	config *GameConfig
}

// NewMinimalExampleGame creates a new minimal example game
func NewMinimalExampleGame() *MinimalExampleGame {
	return &MinimalExampleGame{
		config: &GameConfig{
			Name:        "minimal_example",
			Description: "The simplest possible game - just a counter",
			PartitionNames: map[string]string{
				"counter": "counter_state",
			},
			ServerPartitionNames: []string{"counter_state"},
			VisualizationConfig: &VisualizationConfig{
				CanvasWidth:      400,
				CanvasHeight:     200,
				BackgroundColor:  "#2a2a2a",
				UpdateIntervalMs: 100,
				Renderers: []RendererConfig{
					{
						Type:          "text",
						PartitionName: "counter_state",
						Properties: map[string]interface{}{
							"fontSize": 24,
							"color":    "#ffffff",
							"x":        200,
							"y":        100,
							"text":     "Count: {value}",
						},
					},
				},
			},
			ImplementationConfig: &ImplementationConfig{
				Iterations: map[string]simulator.Iteration{
					"counter_state": &MinimalCounterIteration{},
				},
				OutputCondition: &simulator.EveryStepOutputCondition{},
				OutputFunction:  &simulator.StdoutOutputFunction{},
				TerminationCondition: &simulator.TimeElapsedTerminationCondition{
					MaxTimeElapsed: 30.0, // 30 seconds
				},
				TimestepFunction: &simulator.ConstantTimestepFunction{
					Stepsize: 1.0, // 1 second per step
				},
			},
			Parameters: map[string]interface{}{
				"initial_value": 0,
				"increment":     1,
			},
		},
	}
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
	return &MinimalExampleRenderer{config: m.config.VisualizationConfig}
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

// MinimalExampleRenderer handles the visualization
type MinimalExampleRenderer struct {
	config *VisualizationConfig
}

func (r *MinimalExampleRenderer) GetVisualizationConfig() *VisualizationConfig {
	return r.config
}

func (r *MinimalExampleRenderer) GetJavaScriptCode() string {
	return `
// Minimal example visualization JavaScript
class MinimalExampleRenderer {
    constructor(canvas, config) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.config = config;
        this.counterValue = 0;
    }
    
    update(partitionState) {
        if (partitionState.partitionName === 'counter_state') {
            this.counterValue = partitionState.state.values[0];
        }
    }
    
    render() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        // Draw counter text
        this.ctx.fillStyle = '#ffffff';
        this.ctx.font = '24px Arial';
        this.ctx.textAlign = 'center';
        this.ctx.fillText('Count: ' + Math.floor(this.counterValue), 
                          this.canvas.width / 2, this.canvas.height / 2);
    }
}

// Global renderer instance
let minimalRenderer = null;

function initializeRenderer(canvas, config) {
    minimalRenderer = new MinimalExampleRenderer(canvas, config);
}

function updateVisualization(partitionState) {
    if (minimalRenderer) {
        minimalRenderer.update(partitionState);
        minimalRenderer.render();
    }
}
`
}

func (r *MinimalExampleRenderer) GetCSSCode() string {
	return `
.minimal-example {
    background-color: #2a2a2a;
    border: 2px solid #444;
    border-radius: 8px;
}

.minimal-example canvas {
    display: block;
    margin: 0 auto;
}
`
}
