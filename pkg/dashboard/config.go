// Package game describes a stochadex simulation that's meant to run in
// the browser as a WebAssembly module, together with the visualization
// metadata the in-browser renderer needs to draw its state.
//
// The headline type is Config: it bundles the partition wiring (which
// state names are streamed to action sources, which are driven by them)
// with a VisualizationConfig (the canvas + a list of renderer descriptors)
// and a SimulationGenerator function (which the wasm-side runtime calls
// to obtain a stochadex ConfigGenerator).
//
// The two builders, ConfigBuilder and VisualizationBuilder, exist so that
// example simulations can describe themselves declaratively. Neither
// performs validation; both are just typed convenience wrappers over the
// underlying config structs.
package dashboard

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// Config is the complete description of one browser-runnable simulation:
// what to simulate, what state to stream out, what state can be driven by
// external action input, and how to draw it.
//
// A single value of this type is consumed by two code paths:
//   - the template generator (which uses it to emit the per-example HTML/JS
//     shell), and
//   - the wasm-side runtime in pkg/simio (which uses it to wire up the
//     stochadex coordinator together with the action-state dispatch).
type Config struct {
	// Name identifies the simulation. Used as a folder name by the template
	// generator and as the page title in the generated HTML.
	Name string

	// Description appears beneath the title in the generated HTML.
	Description string

	// ServerPartitionNames lists the stochadex partition names whose state
	// the runtime should publish on every output step. Each named partition
	// is both forwarded to the visualization renderer and (when an external
	// action driver is connected, e.g. a Python websocket server) sent to
	// that driver so it can decide its next action.
	ServerPartitionNames []string

	// ActionStatePartitionNames lists the partitions whose `action_state_values`
	// param the runtime is allowed to overwrite each step from incoming
	// ActionState messages. The two delivery paths use this list differently:
	//
	//   - The legacy broadcast path (ActionState.Values, no named map) writes
	//     the same Values slice to every partition in this list.
	//   - The per-partition named path (ActionState.Partitions) writes only to
	//     partitions whose name appears in the incoming map AND in this list.
	//
	// A simulation that takes no external action input leaves this empty.
	ActionStatePartitionNames []string

	// VisualizationConfig is the renderer-side description: canvas size,
	// background, and the ordered list of shapes/charts to draw using which
	// partition's state.
	VisualizationConfig *VisualizationConfig

	// SimulationGenerator is invoked once at startup to obtain a fresh
	// stochadex ConfigGenerator. The runtime then replaces its OutputCondition
	// and OutputFunction with the wasm-side equivalents and wires action
	// state delivery into the resulting coordinator.
	SimulationGenerator func() *simulator.ConfigGenerator

	// Sliders declare HTML range inputs the codegen should emit into the
	// Live controls panel. Each slider writes to one (Partition, ValueIndex)
	// slot in the inline driver's outgoing action vector.
	Sliders []Slider

	// Readouts declare DOM text elements the codegen should emit into the
	// chart panel(s). Each readout subscribes to one partition's state and
	// formats it via a small template (see Readout.Template).
	Readouts []Readout

	// ShowReset toggles a "Reset simulation" button in the controls panel
	// that re-launches the worker on click. Useful for dashboards where
	// the user wants to restart the simulation without reloading the page.
	ShowReset bool

	// Driver selects which action driver runtime/worker.js loads and what
	// options to pass it. Build() fills in a sensible default if unset.
	Driver DriverSpec
}

// Slider declares a numeric range input that drives one slot of one
// action partition's `action_state_values` vector. Sliders are wired by
// the codegen-emitted game.js to postMessage 'setActions' to the inline
// driver on every input event.
//
// Multiple sliders may share a Partition (with distinct ValueIndex) to
// drive a multi-dimensional action vector. The generated JS groups them
// by partition and emits one entry per partition per publish.
type Slider struct {
	// Name is a unique slug used as the HTML id of both the <input> and
	// its readout element (id="<name>-slider", id="<name>-readout").
	Name string

	// Label is the human-readable string displayed alongside the slider.
	Label string

	// Partition is the ActionStatePartitionName this slider's value lands on.
	Partition string

	// ValueIndex is the position in the partition's action vector that
	// this slider writes. The publish step zero-fills any unused slots.
	ValueIndex int

	Min, Max, Step, Default float64

	// Decimals is the number of fractional digits shown in the on-page
	// readout. Defaults to 3 when zero.
	Decimals int
}

// Readout declares a DOM text element that displays values from one
// partition's most recent state. The Template uses simple tokens that
// the generated JS substitutes at render time:
//
//	{t}        cumulative timesteps (rendered as an integer)
//	{v}, {v0}  state[0] formatted to Decimals decimal places
//	{vN}       state[N] (N a non-negative integer)
//
// One Readout per Partition is typical; multiple Readouts per partition
// are allowed.
type Readout struct {
	Partition string
	Template  string
	// Decimals is the number of fractional digits used to format {v}, {vN}.
	// Defaults to 2 when zero.
	Decimals int
}

// DriverSpec selects which file under runtime/drivers/ the worker loads
// at startup. Kind is the driver name (e.g. "inline", "websocket");
// Options is forwarded verbatim to createDriver as a JS object — its
// expected shape is the relevant driver's documentation.
type DriverSpec struct {
	Kind    string
	Options map[string]interface{}
}

// VisualizationConfig is the static description of a canvas-based view of a
// simulation. The runtime hands one of these to runtime/renderer.js, which
// draws the listed Renderers on each animation frame using the partition
// states it has most recently received from the wasm module.
type VisualizationConfig struct {
	CanvasWidth      int
	CanvasHeight     int
	BackgroundColor  string
	UpdateIntervalMs int

	// Renderers describes how to project named partition states onto the
	// canvas. Rendering happens in list order; later entries draw on top.
	Renderers []RendererConfig
}

// RendererConfig is one drawing element on the canvas, bound to a single
// partition's state.
type RendererConfig struct {
	// Type selects which renderer in runtime/renderer.js handles this entry
	// (e.g. "text", "circle", "rectangleSet", "progressBar"). The set of
	// supported types is defined by the renderer's switch statement.
	Type string

	// PartitionName names the stochadex partition whose state values feed
	// this renderer. The empty string indicates a static element (e.g. a
	// background frame) that draws unconditionally.
	PartitionName string

	// Properties carries the renderer-specific draw parameters (positions,
	// colours, formatting). Marshalled into the generated JS as a plain
	// JS object literal.
	Properties map[string]interface{}
}

// VisualizationBuilder is a small fluent helper for assembling a
// VisualizationConfig one renderer at a time. Each Add* method appends a
// single RendererConfig in declaration order; nothing is validated.
type VisualizationBuilder struct {
	config *VisualizationConfig
}

// NewVisualizationBuilder returns a builder seeded with conservative
// defaults (small dark canvas, 100 ms update interval). The defaults
// matter only if the caller doesn't override them via WithCanvas /
// WithBackground / WithUpdateInterval.
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

func (vb *VisualizationBuilder) WithCanvas(width, height int) *VisualizationBuilder {
	vb.config.CanvasWidth = width
	vb.config.CanvasHeight = height
	return vb
}

func (vb *VisualizationBuilder) WithBackground(color string) *VisualizationBuilder {
	vb.config.BackgroundColor = color
	return vb
}

func (vb *VisualizationBuilder) WithUpdateInterval(ms int) *VisualizationBuilder {
	vb.config.UpdateIntervalMs = ms
	return vb
}

// AddText appends a text label. The optional template token "{value}" inside
// `text` is replaced at render time by the floored first element of the
// bound partition's state (typically used for score readouts).
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

// AddRectangle appends a static axis-aligned rectangle with fixed position
// and size. For position/size driven by simulation state, use AddRectangleSet
// instead.
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

// AddRectangleSet appends a renderer that draws one rectangle for every
// (x, y, width, height) group of four floats found in the bound partition's
// state. Entries whose width or height is zero are skipped, so simulations
// can compact entire slots by emitting all-zeros. By default (x, y) is the
// rectangle's centre; set options.Anchor = "topLeft" to use the top-left
// instead.
func (vb *VisualizationBuilder) AddRectangleSet(partitionName string, width, height int, options *ShapeOptions) *VisualizationBuilder {
	props := map[string]interface{}{
		"defaultWidth":  width,
		"defaultHeight": height,
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
		if options.Anchor != "" {
			props["anchor"] = options.Anchor
		}
	}
	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "rectangleSet",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

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

// AddBarChart appends a vertical bar that fills proportionally to the first
// value of the bound partition's state, clamped against options.MaxValue.
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

// AddLineChart appends a rolling line plot of the bound partition's first
// state value over time. The renderer keeps the most recent 100 samples.
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

// AddProgressBar appends a horizontal fill bar driven by the first value of
// the bound partition's state, scaled by options.MaxValue.
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

// AddImage appends a sprite or static image. The current runtime renderer
// draws a placeholder rectangle in place of the image; full image loading
// has not been implemented yet.
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

// AddPointSet appends a renderer that draws one filled circle for every
// (x, y) pair of floats found in the bound partition's state. Useful for
// drawing dynamic populations of point-shaped objects (players, vehicles,
// agents) without one renderer per object.
func (vb *VisualizationBuilder) AddPointSet(partitionName string, options *PointSetOptions) *VisualizationBuilder {
	props := map[string]interface{}{}
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
		if options.Radius != 0 {
			props["radius"] = options.Radius
		}
	}
	vb.config.Renderers = append(vb.config.Renderers, RendererConfig{
		Type:          "pointSet",
		PartitionName: partitionName,
		Properties:    props,
	})
	return vb
}

func (vb *VisualizationBuilder) Build() *VisualizationConfig {
	return vb.config
}

// Renderer-specific option structs. Fields are passed straight through to
// runtime/renderer.js as JS object properties; semantics are documented on
// the renderer methods that consume them.

type TextOptions struct {
	FontSize   int
	Color      string
	FontFamily string
	TextAlign  string
}

type ShapeOptions struct {
	Color       string
	FillColor   string
	StrokeColor string
	StrokeWidth int
	// Anchor controls how (x, y) in a rectangleSet's state is interpreted.
	// "topLeft" matches Canvas/AddRectangle convention; the default ("" or
	// "center") keeps (x, y) as the rectangle's centre.
	Anchor string
}

type LineOptions struct {
	Color       string
	Width       int
	DashPattern []int
}

type ChartOptions struct {
	Color       string
	MaxValue    float64
	ShowLabels  bool
	LabelFormat string
	LineWidth   int
}

type ProgressBarOptions struct {
	BackgroundColor string
	ForegroundColor string
	BorderColor     string
	BorderWidth     int
	ShowLabel       bool
	LabelFormat     string
	MaxValue        float64
}

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

type PointSetOptions struct {
	Color       string
	FillColor   string
	StrokeColor string
	StrokeWidth int
	Radius      int
}

// ConfigBuilder is a small fluent helper for assembling a Config. Like
// VisualizationBuilder, it performs no validation; the only invariant it
// enforces is that the slice fields start non-nil.
type ConfigBuilder struct {
	config *Config
}

func NewConfigBuilder(name string) *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{
			Name:                 name,
			ServerPartitionNames: make([]string, 0),
		},
	}
}

func (gb *ConfigBuilder) WithDescription(description string) *ConfigBuilder {
	gb.config.Description = description
	return gb
}

// WithServerPartition declares that the named partition's state should be
// streamed out of the wasm module each step — both to the visualization
// renderer and to any connected action driver.
func (gb *ConfigBuilder) WithServerPartition(partitionName string) *ConfigBuilder {
	gb.config.ServerPartitionNames = append(gb.config.ServerPartitionNames, partitionName)
	return gb
}

// WithActionStatePartition declares that the named partition reads its
// `action_state_values` param from incoming ActionState messages. See
// Config.ActionStatePartitionNames for the full dispatch semantics.
func (gb *ConfigBuilder) WithActionStatePartition(partitionName string) *ConfigBuilder {
	gb.config.ActionStatePartitionNames = append(gb.config.ActionStatePartitionNames, partitionName)
	return gb
}

func (gb *ConfigBuilder) WithVisualization(config *VisualizationConfig) *ConfigBuilder {
	gb.config.VisualizationConfig = config
	return gb
}

// WithSimulation registers the per-step simulation builder. The runtime
// calls this once at startup to obtain a fresh stochadex ConfigGenerator
// (with its partitions and simulation already declared); the runtime then
// replaces the generator's output wiring with the wasm-side OutputFunction.
func (gb *ConfigBuilder) WithSimulation(simGen func() *simulator.ConfigGenerator) *ConfigBuilder {
	gb.config.SimulationGenerator = simGen
	return gb
}

// WithSlider appends a slider to the dashboard's Live controls panel.
// The slider drives Partition[ValueIndex] of the inline driver's next
// outgoing action vector. Defaults to 3-decimal readout formatting if
// Slider.Decimals is zero.
func (gb *ConfigBuilder) WithSlider(s Slider) *ConfigBuilder {
	gb.config.Sliders = append(gb.config.Sliders, s)
	return gb
}

// WithReadout appends a DOM readout that displays formatted values from
// one partition's most-recent state. Defaults to 2-decimal value
// formatting if Readout.Decimals is zero.
func (gb *ConfigBuilder) WithReadout(r Readout) *ConfigBuilder {
	gb.config.Readouts = append(gb.config.Readouts, r)
	return gb
}

// WithResetButton enables the "Reset simulation" button in the controls
// panel. The button terminates and re-launches the wasm worker so the
// simulation restarts from its initial state.
func (gb *ConfigBuilder) WithResetButton() *ConfigBuilder {
	gb.config.ShowReset = true
	return gb
}

// WithInlineDriver selects the in-page action driver and sets its tick
// interval. Typical values are 30–100 ms (30 Hz – 10 Hz). Use this for
// dashboards driven by page UI (sliders, buttons, keyboard).
func (gb *ConfigBuilder) WithInlineDriver(intervalMs int) *ConfigBuilder {
	gb.config.Driver = DriverSpec{
		Kind:    "inline",
		Options: map[string]interface{}{"intervalMs": intervalMs},
	}
	return gb
}

// WithWebsocketDriver selects the WebSocket-based action driver. The
// driver connects to `url` (default ws://localhost:2112), forwards every
// streamed partition state to that socket, and treats each inbound
// message as ActionState bytes for the next step.
func (gb *ConfigBuilder) WithWebsocketDriver(url string) *ConfigBuilder {
	opts := map[string]interface{}{}
	if url != "" {
		opts["url"] = url
	}
	// forwardPartitions is filled in at Build() time from ServerPartitionNames
	// (so the caller doesn't have to repeat them).
	gb.config.Driver = DriverSpec{Kind: "websocket", Options: opts}
	return gb
}

// Build finalises and returns the Config. It fills in any defaults that
// depend on prior builder calls (currently: the websocket driver's
// forwardPartitions defaulting to ServerPartitionNames; and the driver
// itself defaulting to "websocket" if WithInlineDriver/WithWebsocketDriver
// wasn't called).
func (gb *ConfigBuilder) Build() *Config {
	if gb.config.Driver.Kind == "" {
		gb.config.Driver = DriverSpec{Kind: "websocket"}
	}
	if gb.config.Driver.Kind == "websocket" {
		if gb.config.Driver.Options == nil {
			gb.config.Driver.Options = map[string]interface{}{}
		}
		if _, set := gb.config.Driver.Options["forwardPartitions"]; !set {
			// Copy to avoid the generated JS aliasing the live slice.
			fp := make([]string, len(gb.config.ServerPartitionNames))
			copy(fp, gb.config.ServerPartitionNames)
			gb.config.Driver.Options["forwardPartitions"] = fp
		}
	}
	return gb.config
}
