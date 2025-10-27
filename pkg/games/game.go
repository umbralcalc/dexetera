package games

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// Game represents a complete game implementation that can be compiled to WebAssembly
// and run in the browser with stochadex simulations.
type Game interface {
	// GetName returns the unique name of the game
	GetName() string

	// GetDescription returns a human-readable description of the game
	GetDescription() string

	// GetConfig returns the game-specific configuration
	GetConfig() *GameConfig

	// GetConfigGenerator returns a configured ConfigGenerator that builds the simulation
	// configuration step-by-step using the fluent API
	GetConfigGenerator() *simulator.ConfigGenerator

	// GetRenderer returns the visualization renderer for this game
	GetRenderer() GameRenderer
}

// GameConfig holds configuration for a game including partition names,
// visualization settings, and other game-specific parameters.
type GameConfig struct {
	// Name is the unique identifier for the game
	Name string

	// Description is a human-readable description
	Description string

	// PartitionNames maps logical names to stochadex partition names
	// This allows the frontend to reference partitions by meaningful names
	// instead of indices
	PartitionNames map[string]string

	// ServerPartitionNames are the partition names that should be sent
	// to the user's Python websocket server
	ServerPartitionNames []string

	// VisualizationConfig holds game-specific visualization settings
	VisualizationConfig *VisualizationConfig

	// ImplementationConfig holds simulation implementation settings
	ImplementationConfig *ImplementationConfig

	// Game-specific parameters
	Parameters map[string]interface{}
}

// VisualizationConfig holds configuration for how the game should be rendered
// in the browser. This allows Go code to specify visualization details.
type VisualizationConfig struct {
	// CanvasWidth and CanvasHeight define the rendering area
	CanvasWidth  int
	CanvasHeight int

	// BackgroundColor is the CSS color for the background
	BackgroundColor string

	// Renderers defines what should be rendered and how
	Renderers []RendererConfig

	// UpdateIntervalMs is how often to update the visualization (in milliseconds)
	UpdateIntervalMs int
}

// RendererConfig defines how a specific game element should be rendered
type RendererConfig struct {
	// Type defines the renderer type (e.g., "circle", "rectangle", "line", "text")
	Type string

	// PartitionName is the stochadex partition this renderer displays
	PartitionName string

	// Properties holds renderer-specific configuration
	Properties map[string]interface{}
}

// ImplementationConfig holds configuration for simulation implementations
type ImplementationConfig struct {
	// Iterations defines the iteration implementations for each partition
	Iterations map[string]simulator.Iteration

	// OutputCondition defines when to output simulation state
	OutputCondition simulator.OutputCondition

	// OutputFunction defines how to output simulation state
	OutputFunction simulator.OutputFunction

	// TerminationCondition defines when the simulation should end
	TerminationCondition simulator.TerminationCondition

	// TimestepFunction defines how to compute time increments
	TimestepFunction simulator.TimestepFunction
}

// GameRenderer defines how a game should be visualized in the browser.
// This interface allows Go code to specify visualization details that get
// converted to JavaScript configuration.
type GameRenderer interface {
	// GetVisualizationConfig returns the complete visualization configuration
	GetVisualizationConfig() *VisualizationConfig

	// GetJavaScriptCode returns any custom JavaScript code needed for rendering
	GetJavaScriptCode() string

	// GetCSSCode returns any custom CSS code needed for styling
	GetCSSCode() string
}

// Helper functions for creating common implementation configurations

// NewDefaultImplementationConfig creates a default implementation configuration
// suitable for most games
func NewDefaultImplementationConfig() *ImplementationConfig {
	return &ImplementationConfig{
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.StdoutOutputFunction{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{
			MaxTimeElapsed: 30.0, // 30 seconds default
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{
			Stepsize: 1.0, // 1 second per step default
		},
		Iterations: make(map[string]simulator.Iteration),
	}
}

// NewWebSocketImplementationConfig creates an implementation configuration
// optimized for WebSocket-based games
func NewWebSocketImplementationConfig() *ImplementationConfig {
	return &ImplementationConfig{
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.StdoutOutputFunction{}, // Will be replaced with JsCallbackOutputFunction
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{
			MaxTimeElapsed: 60.0, // 60 seconds for WebSocket games
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{
			Stepsize: 1.0, // 1 second per step
		},
		Iterations: make(map[string]simulator.Iteration),
	}
}

// ToImplementations converts an ImplementationConfig to a simulator.Implementations
func (ic *ImplementationConfig) ToImplementations() *simulator.Implementations {
	iterations := make([]simulator.Iteration, 0, len(ic.Iterations))
	for _, iteration := range ic.Iterations {
		iterations = append(iterations, iteration)
	}

	return &simulator.Implementations{
		Iterations:           iterations,
		OutputCondition:      ic.OutputCondition,
		OutputFunction:       ic.OutputFunction,
		TerminationCondition: ic.TerminationCondition,
		TimestepFunction:     ic.TimestepFunction,
	}
}

// GameBuilder provides a fluent API for building games
type GameBuilder struct {
	config *GameConfig
}

// NewGameBuilder creates a new GameBuilder with the given game name
func NewGameBuilder(name string) *GameBuilder {
	return &GameBuilder{
		config: &GameConfig{
			Name:                 name,
			PartitionNames:       make(map[string]string),
			ServerPartitionNames: make([]string, 0),
			Parameters:           make(map[string]interface{}),
			ImplementationConfig: NewDefaultImplementationConfig(),
		},
	}
}

// WithDescription sets the game description
func (gb *GameBuilder) WithDescription(description string) *GameBuilder {
	gb.config.Description = description
	return gb
}

// WithPartition adds a partition to the game
func (gb *GameBuilder) WithPartition(logicalName, partitionName string, iteration simulator.Iteration) *GameBuilder {
	gb.config.PartitionNames[logicalName] = partitionName
	gb.config.ImplementationConfig.Iterations[partitionName] = iteration
	return gb
}

// WithServerPartition adds a partition that should be sent to the Python websocket server
func (gb *GameBuilder) WithServerPartition(partitionName string) *GameBuilder {
	gb.config.ServerPartitionNames = append(gb.config.ServerPartitionNames, partitionName)
	return gb
}

// WithParameter adds a game-specific parameter
func (gb *GameBuilder) WithParameter(key string, value interface{}) *GameBuilder {
	gb.config.Parameters[key] = value
	return gb
}

// WithVisualization sets the visualization configuration
func (gb *GameBuilder) WithVisualization(config *VisualizationConfig) *GameBuilder {
	gb.config.VisualizationConfig = config
	return gb
}

// WithOutputCondition sets the output condition
func (gb *GameBuilder) WithOutputCondition(condition simulator.OutputCondition) *GameBuilder {
	gb.config.ImplementationConfig.OutputCondition = condition
	return gb
}

// WithOutputFunction sets the output function
func (gb *GameBuilder) WithOutputFunction(function simulator.OutputFunction) *GameBuilder {
	gb.config.ImplementationConfig.OutputFunction = function
	return gb
}

// WithTerminationCondition sets the termination condition
func (gb *GameBuilder) WithTerminationCondition(condition simulator.TerminationCondition) *GameBuilder {
	gb.config.ImplementationConfig.TerminationCondition = condition
	return gb
}

// WithTimestepFunction sets the timestep function
func (gb *GameBuilder) WithTimestepFunction(function simulator.TimestepFunction) *GameBuilder {
	gb.config.ImplementationConfig.TimestepFunction = function
	return gb
}

// WithMaxTime sets the maximum simulation time
func (gb *GameBuilder) WithMaxTime(maxTime float64) *GameBuilder {
	gb.config.ImplementationConfig.TerminationCondition = &simulator.TimeElapsedTerminationCondition{
		MaxTimeElapsed: maxTime,
	}
	return gb
}

// WithTimestep sets the timestep size
func (gb *GameBuilder) WithTimestep(stepsize float64) *GameBuilder {
	gb.config.ImplementationConfig.TimestepFunction = &simulator.ConstantTimestepFunction{
		Stepsize: stepsize,
	}
	return gb
}

// Build creates the final GameConfig
func (gb *GameBuilder) Build() *GameConfig {
	return gb.config
}

// GenericGame is a generic game implementation that works with GameBuilder
type GenericGame struct {
	config *GameConfig
}

// NewGenericGame creates a new GenericGame from a GameConfig
func NewGenericGame(config *GameConfig) *GenericGame {
	return &GenericGame{config: config}
}

// GetName returns the game name
func (g *GenericGame) GetName() string {
	return g.config.Name
}

// GetDescription returns the game description
func (g *GenericGame) GetDescription() string {
	return g.config.Description
}

// GetConfig returns the game configuration
func (g *GenericGame) GetConfig() *GameConfig {
	return g.config
}

// GetConfigGenerator returns a configured ConfigGenerator that builds the simulation
// configuration step-by-step using the fluent API
func (g *GenericGame) GetConfigGenerator() *simulator.ConfigGenerator {
	configGen := simulator.NewConfigGenerator()
	configGen.SetGlobalSeed(42) // Default seed, can be made configurable

	// Add partitions from the game configuration
	for logicalName, partitionName := range g.config.PartitionNames {
		// Get parameters for this partition
		var params map[string][]float64
		if paramValue, exists := g.config.Parameters[logicalName+"_params"]; exists {
			if paramMap, ok := paramValue.(map[string][]float64); ok {
				params = paramMap
			}
		}
		if params == nil {
			params = make(map[string][]float64)
		}

		// Get initial state values
		var initValues []float64
		if initValue, exists := g.config.Parameters[logicalName+"_init"]; exists {
			if values, ok := initValue.([]float64); ok {
				initValues = values
			}
		}
		if initValues == nil {
			initValues = []float64{0.0} // Default initial value
		}

		partition := &simulator.PartitionConfig{
			Name:            partitionName,
			Params:          simulator.NewParams(params),
			InitStateValues: initValues,
		}
		configGen.SetPartition(partition)
	}

	// Configure simulation-level settings
	simulationConfig := &simulator.SimulationConfig{
		InitTimeValue: 0.0, // Default, can be made configurable
	}
	configGen.SetSimulation(simulationConfig)

	return configGen
}

// GetRenderer returns the visualization renderer
func (g *GenericGame) GetRenderer() GameRenderer {
	if g.config.VisualizationConfig == nil {
		return &DefaultRenderer{}
	}
	return &GenericRenderer{config: g.config.VisualizationConfig}
}

// DefaultRenderer provides a basic default renderer
type DefaultRenderer struct{}

func (r *DefaultRenderer) GetVisualizationConfig() *VisualizationConfig {
	return &VisualizationConfig{
		CanvasWidth:      400,
		CanvasHeight:     200,
		BackgroundColor:  "#2a2a2a",
		UpdateIntervalMs: 100,
		Renderers:        []RendererConfig{},
	}
}

func (r *DefaultRenderer) GetJavaScriptCode() string {
	return `
// Default renderer JavaScript
class DefaultRenderer {
    constructor(canvas, config) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.config = config;
    }
    
    update(partitionState) {
        // Default: do nothing
    }
    
    render() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.fillStyle = '#ffffff';
        this.ctx.font = '16px Arial';
        this.ctx.textAlign = 'center';
        this.ctx.fillText('Default Game Renderer', 
                          this.canvas.width / 2, this.canvas.height / 2);
    }
}

// Global renderer instance
let defaultRenderer = null;

function initializeRenderer(canvas, config) {
    defaultRenderer = new DefaultRenderer(canvas, config);
}

function updateVisualization(partitionState) {
    if (defaultRenderer) {
        defaultRenderer.update(partitionState);
        defaultRenderer.render();
    }
}
`
}

func (r *DefaultRenderer) GetCSSCode() string {
	return `
.default-game {
    background-color: #2a2a2a;
    border: 2px solid #444;
    border-radius: 8px;
}

.default-game canvas {
    display: block;
    margin: 0 auto;
}
`
}

// GenericRenderer provides a generic renderer based on configuration
type GenericRenderer struct {
	config *VisualizationConfig
}

func (r *GenericRenderer) GetVisualizationConfig() *VisualizationConfig {
	return r.config
}

func (r *GenericRenderer) GetJavaScriptCode() string {
	return `
// Generic renderer JavaScript
class GenericRenderer {
    constructor(canvas, config) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.config = config;
        this.state = {};
    }
    
    update(partitionState) {
        this.state[partitionState.partitionName] = partitionState.state.values;
    }
    
    render() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        // Render each configured renderer
        this.config.renderers.forEach(renderer => {
            this.renderElement(renderer);
        });
    }
    
    renderElement(renderer) {
        const state = this.state[renderer.partitionName];
        if (!state) return;
        
        switch (renderer.type) {
            case 'text':
                this.renderText(renderer, state);
                break;
            case 'circle':
                this.renderCircle(renderer, state);
                break;
            case 'rectangle':
                this.renderRectangle(renderer, state);
                break;
        }
    }
    
    renderText(renderer, state) {
        this.ctx.fillStyle = renderer.properties.color || '#ffffff';
        this.ctx.font = (renderer.properties.fontSize || 16) + 'px Arial';
        this.ctx.textAlign = 'center';
        
        let text = renderer.properties.text || '{value}';
        text = text.replace('{value}', Math.floor(state[0] || 0));
        
        this.ctx.fillText(text, 
                          renderer.properties.x || this.canvas.width / 2,
                          renderer.properties.y || this.canvas.height / 2);
    }
    
    renderCircle(renderer, state) {
        this.ctx.fillStyle = renderer.properties.color || '#ffffff';
        this.ctx.beginPath();
        this.ctx.arc(renderer.properties.x || this.canvas.width / 2,
                     renderer.properties.y || this.canvas.height / 2,
                     renderer.properties.radius || 10,
                     0, 2 * Math.PI);
        this.ctx.fill();
    }
    
    renderRectangle(renderer, state) {
        this.ctx.fillStyle = renderer.properties.color || '#ffffff';
        this.ctx.fillRect(renderer.properties.x || 0,
                         renderer.properties.y || 0,
                         renderer.properties.width || 50,
                         renderer.properties.height || 50);
    }
}

// Global renderer instance
let genericRenderer = null;

function initializeRenderer(canvas, config) {
    genericRenderer = new GenericRenderer(canvas, config);
}

function updateVisualization(partitionState) {
    if (genericRenderer) {
        genericRenderer.update(partitionState);
        genericRenderer.render();
    }
}
`
}

func (r *GenericRenderer) GetCSSCode() string {
	return `
.generic-game {
    background-color: ` + r.config.BackgroundColor + `;
    border: 2px solid #444;
    border-radius: 8px;
}

.generic-game canvas {
    display: block;
    margin: 0 auto;
}
`
}
