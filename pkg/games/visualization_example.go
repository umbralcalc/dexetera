package games

// VisualizationExampleGame demonstrates the powerful VisualizationBuilder
// This creates a complex game with multiple visualization types
type VisualizationExampleGame struct {
	config *GameConfig
}

// NewVisualizationExampleGame creates a new game using the VisualizationBuilder
func NewVisualizationExampleGame() *BuilderExampleGame {
	// Create a complex visualization using the VisualizationBuilder
	visConfig := NewVisualizationBuilder().
		WithCanvas(800, 600).
		WithBackground("#0a0a0a").
		WithUpdateInterval(50).
		AddText("counter_state", "Counter: {value}", 100, 50, &TextOptions{
			FontSize:   24,
			Color:      "#00ff00",
			FontFamily: "Arial",
			TextAlign:  "left",
		}).
		AddText("timer_state", "Timer: {value}", 100, 100, &TextOptions{
			FontSize:   24,
			Color:      "#ff0000",
			FontFamily: "Arial",
			TextAlign:  "left",
		}).
		AddCircle("counter_state", 200, 200, 30, &ShapeOptions{
			FillColor:   "#00ff00",
			StrokeColor: "#ffffff",
			StrokeWidth: 2,
		}).
		AddRectangle("timer_state", 300, 170, 100, 60, &ShapeOptions{
			FillColor:   "#ff0000",
			StrokeColor: "#ffffff",
			StrokeWidth: 2,
		}).
		AddLine("counter_state", 0, 400, 800, 400, &LineOptions{
			Color: "#444444",
			Width: 2,
		}).
		AddBarChart("counter_state", 50, 450, 200, 100, &ChartOptions{
			Color:       "#00ff00",
			MaxValue:    100,
			ShowLabels:  true,
			LabelFormat: "Count: {value}",
		}).
		AddLineChart("timer_state", 300, 450, 200, 100, &ChartOptions{
			Color:       "#ff0000",
			MaxValue:    50,
			ShowLabels:  true,
			LabelFormat: "Time: {value}s",
			LineWidth:   3,
		}).
		Build()

	// Create the game using GameBuilder with the complex visualization
	config := NewGameBuilder("visualization_example").
		WithDescription("A complex game demonstrating advanced visualization capabilities").
		WithPartition("counter", "counter_state", &CounterIteration{}).
		WithPartition("timer", "timer_state", &TimerIteration{}).
		WithServerPartition("counter_state").
		WithServerPartition("timer_state").
		WithParameter("counter_init", []float64{0.0}).
		WithParameter("counter_params", map[string][]float64{"increment": {2.0}}).
		WithParameter("timer_init", []float64{0.0}).
		WithParameter("timer_params", map[string][]float64{"speed": {0.5}}).
		WithMaxTime(120.0).
		WithTimestep(0.5).
		WithVisualization(visConfig).
		Build()

	return &BuilderExampleGame{config: config}
}
