package game

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

	// ActionStatePartitionNames are the partition names that should have
	// action state values in their params set by the Python websocket
	// server after each step
	ActionStatePartitionNames []string

	// VisualizationConfig holds game-specific visualization settings
	VisualizationConfig *VisualizationConfig

	// ImplementationConfig holds simulation implementation settings
	ImplementationConfig *ImplementationConfig

	// SimulationGenerator builds a simulator.ConfigGenerator independently
	// of the game builder; the framework will call this and wire output
	// callback and output condition.
	SimulationGenerator func() *simulator.ConfigGenerator

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

// VisualizationBuilder provides a fluent API for building complex visualizations
type VisualizationBuilder struct {
	config *VisualizationConfig
}

// NewVisualizationBuilder creates a new VisualizationBuilder
func NewVisualizationBuilder() *VisualizationBuilder {
	return &VisualizationBuilder{
		config: &VisualizationConfig{
			CanvasWidth:      400,
			CanvasHeight:     200,
			BackgroundColor:  "#2a2a2a",
			UpdateIntervalMs: 100,
			Renderers:        make([]RendererConfig, 0),
		},
	}
}

// WithCanvas sets the canvas dimensions
func (vb *VisualizationBuilder) WithCanvas(width, height int) *VisualizationBuilder {
	vb.config.CanvasWidth = width
	vb.config.CanvasHeight = height
	return vb
}

// WithBackground sets the background color
func (vb *VisualizationBuilder) WithBackground(color string) *VisualizationBuilder {
	vb.config.BackgroundColor = color
	return vb
}

// WithUpdateInterval sets the update interval in milliseconds
func (vb *VisualizationBuilder) WithUpdateInterval(ms int) *VisualizationBuilder {
	vb.config.UpdateIntervalMs = ms
	return vb
}

// AddText adds a text renderer
func (vb *VisualizationBuilder) AddText(partitionName string, text string, x, y int, options *TextOptions) *VisualizationBuilder {
	props := map[string]interface{}{
		"text": text,
		"x":    x,
		"y":    y,
	}

	if options != nil {
		if options.FontSize != 0 {
			props["fontSize"] = options.FontSize
		}
		if options.Color != "" {
			props["color"] = options.Color
		}
		if options.FontFamily != "" {
			props["fontFamily"] = options.FontFamily
		}
		if options.TextAlign != "" {
			props["textAlign"] = options.TextAlign
		}
	}

	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "text",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

// AddCircle adds a circle renderer
func (vb *VisualizationBuilder) AddCircle(partitionName string, x, y, radius int, options *ShapeOptions) *VisualizationBuilder {
	props := map[string]interface{}{
		"x":      x,
		"y":      y,
		"radius": radius,
	}

	if options != nil {
		if options.Color != "" {
			props["color"] = options.Color
		}
		if options.FillColor != "" {
			props["fillColor"] = options.FillColor
		}
		if options.StrokeColor != "" {
			props["strokeColor"] = options.StrokeColor
		}
		if options.StrokeWidth != 0 {
			props["strokeWidth"] = options.StrokeWidth
		}
	}

	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "circle",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

// AddRectangle adds a rectangle renderer
func (vb *VisualizationBuilder) AddRectangle(partitionName string, x, y, width, height int, options *ShapeOptions) *VisualizationBuilder {
	props := map[string]interface{}{
		"x":      x,
		"y":      y,
		"width":  width,
		"height": height,
	}

	if options != nil {
		if options.Color != "" {
			props["color"] = options.Color
		}
		if options.FillColor != "" {
			props["fillColor"] = options.FillColor
		}
		if options.StrokeColor != "" {
			props["strokeColor"] = options.StrokeColor
		}
		if options.StrokeWidth != 0 {
			props["strokeWidth"] = options.StrokeWidth
		}
	}

	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "rectangle",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

// AddLine adds a line renderer
func (vb *VisualizationBuilder) AddLine(partitionName string, x1, y1, x2, y2 int, options *LineOptions) *VisualizationBuilder {
	props := map[string]interface{}{
		"x1": x1,
		"y1": y1,
		"x2": x2,
		"y2": y2,
	}

	if options != nil {
		if options.Color != "" {
			props["color"] = options.Color
		}
		if options.Width != 0 {
			props["width"] = options.Width
		}
		if options.DashPattern != nil {
			props["dashPattern"] = options.DashPattern
		}
	}

	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "line",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

// AddBarChart adds a bar chart renderer
func (vb *VisualizationBuilder) AddBarChart(partitionName string, x, y, width, height int, options *ChartOptions) *VisualizationBuilder {
	props := map[string]interface{}{
		"x":      x,
		"y":      y,
		"width":  width,
		"height": height,
	}

	if options != nil {
		if options.Color != "" {
			props["color"] = options.Color
		}
		if options.MaxValue != 0 {
			props["maxValue"] = options.MaxValue
		}
		if options.ShowLabels {
			props["showLabels"] = options.ShowLabels
		}
		if options.LabelFormat != "" {
			props["labelFormat"] = options.LabelFormat
		}
	}

	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "barChart",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

// AddLineChart adds a line chart renderer
func (vb *VisualizationBuilder) AddLineChart(partitionName string, x, y, width, height int, options *ChartOptions) *VisualizationBuilder {
	props := map[string]interface{}{
		"x":      x,
		"y":      y,
		"width":  width,
		"height": height,
	}

	if options != nil {
		if options.Color != "" {
			props["color"] = options.Color
		}
		if options.MaxValue != 0 {
			props["maxValue"] = options.MaxValue
		}
		if options.ShowLabels {
			props["showLabels"] = options.ShowLabels
		}
		if options.LabelFormat != "" {
			props["labelFormat"] = options.LabelFormat
		}
		if options.LineWidth != 0 {
			props["lineWidth"] = options.LineWidth
		}
	}

	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "lineChart",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

// AddProgressBar adds a progress bar renderer
func (vb *VisualizationBuilder) AddProgressBar(partitionName string, x, y, width, height int, options *ProgressBarOptions) *VisualizationBuilder {
	props := map[string]interface{}{
		"x":      x,
		"y":      y,
		"width":  width,
		"height": height,
	}

	if options != nil {
		if options.BackgroundColor != "" {
			props["backgroundColor"] = options.BackgroundColor
		}
		if options.ForegroundColor != "" {
			props["foregroundColor"] = options.ForegroundColor
		}
		if options.BorderColor != "" {
			props["borderColor"] = options.BorderColor
		}
		if options.BorderWidth != 0 {
			props["borderWidth"] = options.BorderWidth
		}
		if options.ShowLabel {
			props["showLabel"] = options.ShowLabel
		}
		if options.LabelFormat != "" {
			props["labelFormat"] = options.LabelFormat
		}
		if options.MaxValue != 0 {
			props["maxValue"] = options.MaxValue
		}
	}

	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "progressBar",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

// AddImage adds an image/sprite renderer
func (vb *VisualizationBuilder) AddImage(partitionName, imagePath string, x, y int, options *ImageOptions) *VisualizationBuilder {
	props := map[string]interface{}{
		"imagePath": imagePath,
		"x":         x,
		"y":         y,
	}

	if options != nil {
		if options.Width != 0 {
			props["width"] = options.Width
		}
		if options.Height != 0 {
			props["height"] = options.Height
		}
		if options.Rotation != 0 {
			props["rotation"] = options.Rotation
		}
		if options.Opacity != 0 {
			props["opacity"] = options.Opacity
		}
		if options.SpriteSheetX != 0 {
			props["spriteSheetX"] = options.SpriteSheetX
		}
		if options.SpriteSheetY != 0 {
			props["spriteSheetY"] = options.SpriteSheetY
		}
		if options.CenterX {
			props["centerX"] = options.CenterX
		}
		if options.CenterY {
			props["centerY"] = options.CenterY
		}
	}

	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "image",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

// Build creates the final VisualizationConfig
func (vb *VisualizationBuilder) Build() *VisualizationConfig {
	return vb.config
}

// Options structs for different renderer types

// TextOptions provides options for text rendering
type TextOptions struct {
	FontSize   int
	Color      string
	FontFamily string
	TextAlign  string
}

// ShapeOptions provides options for shape rendering
type ShapeOptions struct {
	Color       string
	FillColor   string
	StrokeColor string
	StrokeWidth int
}

// LineOptions provides options for line rendering
type LineOptions struct {
	Color       string
	Width       int
	DashPattern []int
}

// ChartOptions provides options for chart rendering
type ChartOptions struct {
	Color       string
	MaxValue    float64
	ShowLabels  bool
	LabelFormat string
	LineWidth   int
}

// ProgressBarOptions provides options for progress bar rendering
type ProgressBarOptions struct {
	BackgroundColor string
	ForegroundColor string
	BorderColor     string
	BorderWidth     int
	ShowLabel       bool
	LabelFormat     string
	MaxValue        float64
}

// ImageOptions provides options for image/sprite rendering
type ImageOptions struct {
	Width        int
	Height       int
	Rotation     float64
	Opacity      float64
	SpriteSheetX int
	SpriteSheetY int
	CenterX      bool
	CenterY      bool
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

// WithServerPartition adds a partition that should be sent to the Python websocket server
func (gb *GameBuilder) WithServerPartition(partitionName string) *GameBuilder {
	gb.config.ServerPartitionNames = append(gb.config.ServerPartitionNames, partitionName)
	return gb
}

// WithActionStatePartition adds a partition that should have action state values in its params set by
// the Python websocket server after each step
func (gb *GameBuilder) WithActionStatePartition(partitionName string) *GameBuilder {
	gb.config.ActionStatePartitionNames = append(gb.config.ActionStatePartitionNames, partitionName)
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

// WithSimulation sets the simulation generator used to produce simulator configs
func (gb *GameBuilder) WithSimulation(simGen func() *simulator.ConfigGenerator) *GameBuilder {
	gb.config.SimulationGenerator = simGen
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

// GetRenderer returns the visualization renderer
func (g *GenericGame) GetRenderer() GameRenderer {
	if g.config.VisualizationConfig == nil {
		return &DefaultRenderer{}
	}
	return &GenericRenderer{Config: g.config.VisualizationConfig}
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
	Config *VisualizationConfig
}

func (r *GenericRenderer) GetVisualizationConfig() *VisualizationConfig {
	return r.Config
}

func (r *GenericRenderer) GetJavaScriptCode() string {
	return `
// Enhanced Generic renderer JavaScript with support for all renderer types
class GenericRenderer {
    constructor(canvas, config) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.config = config;
        this.state = {};
        this.history = {}; // For charts
    }
    
    update(partitionState) {
        this.state[partitionState.partitionName] = partitionState.state.values;
        
        // Store history for charts
        if (!this.history[partitionState.partitionName]) {
            this.history[partitionState.partitionName] = [];
        }
        this.history[partitionState.partitionName].push({
            value: partitionState.state.values[0] || 0,
            time: partitionState.cumulativeTimesteps || 0
        });
        
        // Keep only last 100 points for performance
        if (this.history[partitionState.partitionName].length > 100) {
            this.history[partitionState.partitionName].shift();
        }
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
        if (!state && renderer.partitionName !== '') return;
        
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
            case 'line':
                this.renderLine(renderer, state);
                break;
            case 'barChart':
                this.renderBarChart(renderer, state);
                break;
            case 'lineChart':
                this.renderLineChart(renderer, state);
                break;
            case 'progressBar':
                this.renderProgressBar(renderer, state);
                break;
            case 'image':
                this.renderImage(renderer, state);
                break;
        }
    }
    
    renderText(renderer, state) {
        this.ctx.fillStyle = '#ffffff';
        this.ctx.font = '16px Arial';
        this.ctx.textAlign = 'center';
        
        let text = renderer.properties.text || '{value}';
        text = text.replace('{value}', Math.floor(state[0] || 0));
        
        this.ctx.fillText(text, 
                          renderer.properties.x || this.canvas.width / 2,
                          renderer.properties.y || this.canvas.height / 2);
    }
    
    renderCircle(renderer, state) {
        const x = renderer.properties.x || this.canvas.width / 2;
        const y = renderer.properties.y || this.canvas.height / 2;
        const radius = renderer.properties.radius || 10;
        
        this.ctx.beginPath();
        this.ctx.arc(x, y, radius, 0, 2 * Math.PI);
        
        if (renderer.properties.fillColor) {
            this.ctx.fillStyle = renderer.properties.fillColor;
            this.ctx.fill();
        }
        
        if (renderer.properties.strokeColor) {
            this.ctx.strokeStyle = renderer.properties.strokeColor;
            this.ctx.lineWidth = renderer.properties.strokeWidth || 1;
            this.ctx.stroke();
        }
        
        if (!renderer.properties.fillColor && !renderer.properties.strokeColor) {
            this.ctx.fillStyle = renderer.properties.color || '#ffffff';
            this.ctx.fill();
        }
    }
    
    renderRectangle(renderer, state) {
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        const width = renderer.properties.width || 50;
        const height = renderer.properties.height || 50;
        
        // For static rectangles, always render
        if (renderer.properties.fillColor) {
            this.ctx.fillStyle = renderer.properties.fillColor;
            this.ctx.fillRect(x, y, width, height);
        }
        
        if (renderer.properties.strokeColor) {
            this.ctx.strokeStyle = renderer.properties.strokeColor;
            this.ctx.lineWidth = renderer.properties.strokeWidth || 1;
            this.ctx.strokeRect(x, y, width, height);
        }
        
        if (!renderer.properties.fillColor && !renderer.properties.strokeColor) {
            this.ctx.fillStyle = renderer.properties.color || '#ffffff';
            this.ctx.fillRect(x, y, width, height);
        }
    }
    
    renderLine(renderer, state) {
        const x1 = renderer.properties.x1 || 0;
        const y1 = renderer.properties.y1 || 0;
        const x2 = renderer.properties.x2 || 50;
        const y2 = renderer.properties.y2 || 50;
        
        // For static lines, always render
        this.ctx.beginPath();
        this.ctx.moveTo(x1, y1);
        this.ctx.lineTo(x2, y2);
        this.ctx.strokeStyle = renderer.properties.color || '#ffffff';
        this.ctx.lineWidth = renderer.properties.width || 1;
        this.ctx.stroke();
    }
    
    renderBarChart(renderer, state) {
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        const width = renderer.properties.width || 50;
        const height = renderer.properties.height || 50;
        const maxValue = renderer.properties.maxValue || 100;
        const value = state[0] || 0;
        const normalizedValue = Math.min(value / maxValue, 1.0);
        
        // Draw background
        this.ctx.fillStyle = renderer.properties.color || 'rgba(255,255,255,0.3)';
        this.ctx.fillRect(x, y, width, height);
        
        // Draw bar
        this.ctx.fillStyle = renderer.properties.color || '#4CAF50';
        this.ctx.fillRect(x, y + height * (1 - normalizedValue), width, height * normalizedValue);
        
        // Draw label if requested
        if (renderer.properties.showLabels) {
            this.ctx.fillStyle = '#ffffff';
            this.ctx.font = '12px Arial';
            this.ctx.textAlign = 'center';
            this.ctx.fillText(Math.floor(value), x + width / 2, y + height / 2);
        }
    }
    
    renderLineChart(renderer, state) {
        const history = this.history[renderer.partitionName];
        if (!history || history.length < 2) return;
        
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        const width = renderer.properties.width || 50;
        const height = renderer.properties.height || 50;
        const maxValue = renderer.properties.maxValue || 100;
        
        // Find min/max for scaling
        let minVal = Infinity, maxVal = -Infinity;
        history.forEach(point => {
            minVal = Math.min(minVal, point.value);
            maxVal = Math.max(maxVal, point.value);
        });
        const range = Math.max(maxVal - minVal, 0.1);
        
        this.ctx.strokeStyle = renderer.properties.color || '#4CAF50';
        this.ctx.lineWidth = renderer.properties.lineWidth || 2;
        this.ctx.beginPath();
        
        history.forEach((point, i) => {
            const px = x + (i / (history.length - 1)) * width;
            const py = y + height - ((point.value - minVal) / range) * height;
            
            if (i === 0) {
                this.ctx.moveTo(px, py);
            } else {
                this.ctx.lineTo(px, py);
            }
        });
        
        this.ctx.stroke();
    }
    
    renderProgressBar(renderer, state) {
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        const width = renderer.properties.width || 100;
        const height = renderer.properties.height || 20;
        const maxValue = renderer.properties.maxValue || 100;
        const value = Math.max(0, Math.min(state[0] || 0, maxValue));
        const normalizedValue = value / maxValue;
        
        // Draw background
        this.ctx.fillStyle = renderer.properties.backgroundColor || 'rgba(255,255,255,0.3)';
        this.ctx.fillRect(x, y, width, height);
        
        // Draw progress
        this.ctx.fillStyle = renderer.properties.foregroundColor || '#4CAF50';
        this.ctx.fillRect(x, y, width * normalizedValue, height);
        
        // Draw border if specified
        if (renderer.properties.borderColor) {
            this.ctx.strokeStyle = renderer.properties.borderColor;
            this.ctx.lineWidth = renderer.properties.borderWidth || 1;
            this.ctx.strokeRect(x, y, width, height);
        }
        
        // Draw label if requested
        if (renderer.properties.showLabel) {
            this.ctx.fillStyle = '#ffffff';
            this.ctx.font = '12px Arial';
            this.ctx.textAlign = 'center';
            this.ctx.fillText(Math.floor(value) + '%', x + width / 2, y + height / 2 + 4);
        }
    }
    
    renderImage(renderer, state) {
        const imagePath = renderer.properties.imagePath;
        if (!imagePath) return;
        
        // For now, we'll implement basic rendering
        // In a full implementation, you'd load and cache images
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        
        // Draw placeholder rectangle for now
        this.ctx.fillStyle = 'rgba(255,255,255,0.5)';
        this.ctx.fillRect(x, y, 
            renderer.properties.width || 32, 
            renderer.properties.height || 32);
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
    background-color: ` + r.Config.BackgroundColor + `;
    border: 2px solid #444;
    border-radius: 8px;
}

.generic-game canvas {
    display: block;
    margin: 0 auto;
}
`
}
