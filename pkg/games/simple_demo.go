package games

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// SimpleDemoGame is a minimal example that demonstrates the new game framework.
// This shows how a Go programmer can easily create a new game with the framework.
type SimpleDemoGame struct {
	config *GameConfig
}

// NewSimpleDemoGame creates a new simple demo game
func NewSimpleDemoGame() *SimpleDemoGame {
	return &SimpleDemoGame{
		config: &GameConfig{
			Name:        "simple_demo",
			Description: "A simple demonstration of the new game framework",
			PartitionNames: map[string]string{
				"counter": "counter_state",
			},
			ServerPartitionNames: []string{"counter_state"},
			VisualizationConfig: &VisualizationConfig{
				CanvasWidth:      400,
				CanvasHeight:     300,
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
							"y":        150,
							"text":     "Counter: {value}",
						},
					},
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
func (s *SimpleDemoGame) GetName() string {
	return s.config.Name
}

// GetDescription returns the game description
func (s *SimpleDemoGame) GetDescription() string {
	return s.config.Description
}

// GetConfig returns the game configuration
func (s *SimpleDemoGame) GetConfig() *GameConfig {
	return s.config
}

// GetSettings returns the stochadex settings for this game
func (s *SimpleDemoGame) GetSettings() *simulator.Settings {
	initialValue := s.config.Parameters["initial_value"].(int)
	increment := s.config.Parameters["increment"].(int)

	settings := &simulator.Settings{
		Iterations: []simulator.IterationSettings{
			{
				Name:              "counter_state",
				Params:            simulator.NewParams(map[string][]float64{"increment": {float64(increment)}}),
				InitStateValues:   []float64{float64(initialValue)},
				Seed:              42,
				StateWidth:        1,
				StateHistoryDepth: 1,
			},
		},
		InitTimeValue:         0.0,
		TimestepsHistoryDepth: 1,
	}

	return settings
}

// GetImplementations returns the stochadex implementations
func (s *SimpleDemoGame) GetImplementations() *simulator.Implementations {
	return &simulator.Implementations{
		Iterations: []simulator.Iteration{
			&SimpleCounterIteration{},
		},
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.StdoutOutputFunction{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{
			MaxTimeElapsed: 60.0, // 1 minute
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{
			Stepsize: 1.0, // 1 second per step
		},
	}
}

// GetRenderer returns the visualization renderer
func (s *SimpleDemoGame) GetRenderer() GameRenderer {
	return &SimpleDemoRenderer{config: s.config.VisualizationConfig}
}

// SimpleCounterIteration implements a simple counter that increments
type SimpleCounterIteration struct{}

func (s *SimpleCounterIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	// No configuration needed for this simple counter
}

func (s *SimpleCounterIteration) Iterate(
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

// SimpleDemoRenderer handles the visualization of the simple demo
type SimpleDemoRenderer struct {
	config *VisualizationConfig
}

func (r *SimpleDemoRenderer) GetVisualizationConfig() *VisualizationConfig {
	return r.config
}

func (r *SimpleDemoRenderer) GetJavaScriptCode() string {
	return `
// Simple demo visualization JavaScript
class SimpleDemoRenderer {
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
        this.ctx.fillText('Counter: ' + Math.floor(this.counterValue), 
                          this.canvas.width / 2, this.canvas.height / 2);
    }
}

// Global renderer instance
let simpleDemoRenderer = null;

function initializeRenderer(canvas, config) {
    simpleDemoRenderer = new SimpleDemoRenderer(canvas, config);
}

function updateVisualization(partitionState) {
    if (simpleDemoRenderer) {
        simpleDemoRenderer.update(partitionState);
        simpleDemoRenderer.render();
    }
}
`
}

func (r *SimpleDemoRenderer) GetCSSCode() string {
	return `
.simple-demo {
    background-color: #2a2a2a;
    border: 2px solid #444;
    border-radius: 8px;
}

.simple-demo canvas {
    display: block;
    margin: 0 auto;
}
`
}
