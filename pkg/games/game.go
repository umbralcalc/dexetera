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

	// GetSettings returns the stochadex simulator settings for this game
	GetSettings() *simulator.Settings

	// GetImplementations returns the stochadex implementations for this game
	GetImplementations() *simulator.Implementations

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
