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
