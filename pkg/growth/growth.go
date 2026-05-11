// Package growth is the minimal dashboard example: a single
// logistic-growth partition whose growth rate r and carrying capacity K
// are driven live by sliders through the inline action driver. Everything
// the page needs — visualization, sliders, readout, reset button, driver
// choice — is declared via the dashboard builder, so the static-site
// shell (index.html, styles.css, game.js, build.sh) is produced by
// `go run ./cmd/growth/generate` rather than hand-written.
package growth

import (
	"github.com/umbralcalc/dexetera/pkg/dashboard"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// NewConfig returns the dashboard.Config for the growth example.
//
// The visualization is two renderers on a small white canvas: a baseline
// stroke and a rolling line chart of the population N over the last ~100
// steps. The live numeric value of N (and the current timestep) is shown
// in a DOM readout beneath the chart — not in the canvas — so it picks
// up the page's CSS styling instead of being rendered into pixels.
func NewConfig() *dashboard.Config {
	visConfig := dashboard.NewVisualizationBuilder().
		WithCanvas(320, 160).
		WithBackground("#ffffff").
		WithUpdateInterval(0).
		// Static baseline at the bottom of the chart area.
		AddLine("", 18, 142, 302, 142, &dashboard.LineOptions{
			Color: "#2c3e50",
			Width: 1,
		}).
		// Rolling line chart of N(t). Renderer auto-scales the y-axis
		// to whatever's in its rolling window.
		AddLineChart("population", 18, 18, 284, 124, &dashboard.ChartOptions{
			Color:     "#3c78d8",
			LineWidth: 2,
		}).
		Build()

	return dashboard.NewConfigBuilder("growth").
		WithDescription("Logistic growth: drag the sliders to set r and K live.").
		WithServerPartition("population").
		WithActionStatePartition("population").
		WithVisualization(visConfig).
		WithSimulation(BuildGrowthSimulation).
		// Two sliders, both writing into the "population" partition's
		// action vector. By convention (see BuildGrowthSimulation below),
		// index 0 is r and index 1 is K.
		WithSlider(dashboard.Slider{
			Name: "r", Label: "r (growth rate)",
			Partition: "population", ValueIndex: 0,
			Min: 0, Max: 0.2, Step: 0.005, Default: 0.05,
			Decimals: 3,
		}).
		WithSlider(dashboard.Slider{
			Name: "K", Label: "K (carrying capacity)",
			Partition: "population", ValueIndex: 1,
			Min: 0, Max: 1000, Step: 10, Default: 500,
			Decimals: 0,
		}).
		WithReadout(dashboard.Readout{
			Partition: "population",
			Template:  "t = {t} · N = {v}",
			Decimals:  2,
		}).
		WithResetButton().
		// 50 ms ≈ 20 Hz; the renderer keeps the most recent 100 samples,
		// so the chart shows roughly the last five seconds of growth.
		WithInlineDriver(50).
		Build()
}

// BuildGrowthSimulation constructs the stochadex generator for the single
// logistic-growth partition. The action_state_values param convention is
// `[r, K]`: index 0 is the per-step growth rate, index 1 is the carrying
// capacity. The defaults below are placeholders; the inline driver
// overwrites them on the first slider movement.
//
// MaxTimeElapsed is 10 000 to match the other examples — at the inline
// driver's 50 ms tick that's roughly eight minutes of in-page run-time.
// The page just stops emitting after that; no automatic restart.
func BuildGrowthSimulation() *simulator.ConfigGenerator {
	gen := simulator.NewConfigGenerator()

	population := &simulator.PartitionConfig{
		Name:      "population",
		Iteration: &LogisticGrowthIteration{},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {0.05, 500.0}, // [r, K] defaults
		}),
		InitStateValues:   []float64{10.0},
		StateHistoryDepth: 1,
		Seed:              7,
	}
	gen.SetPartition(population)

	gen.SetSimulation(&simulator.SimulationConfig{
		OutputCondition:      &simulator.EveryStepOutputCondition{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{MaxTimeElapsed: 10000.0},
		TimestepFunction:     &simulator.ConstantTimestepFunction{Stepsize: 1.0},
		InitTimeValue:        0.0,
	})
	return gen
}

// LogisticGrowthIteration implements the discrete logistic update
//
//	N(t+1) = N(t) + r * N(t) * (1 - N(t)/K)
//
// reading r and K from the partition's action_state_values param. With
// K close to zero the population collapses; with K well above N the
// population grows toward K. The result is clamped to ≥ 0.
type LogisticGrowthIteration struct{}

func (l *LogisticGrowthIteration) Configure(int, *simulator.Settings) {}

func (l *LogisticGrowthIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	state := stateHistories[partitionIndex].CopyStateRow(0)

	r, K := 0.05, 500.0
	if action := params.Get("action_state_values"); len(action) >= 2 {
		r = action[0]
		K = action[1]
	}

	n := state[0]
	if K <= 0 {
		state[0] = 0
		return state
	}
	state[0] = n + r*n*(1.0-n/K)
	if state[0] < 0 {
		state[0] = 0
	}
	return state
}
