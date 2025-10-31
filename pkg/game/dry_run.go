package game

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// FullDryRun tests the behaviour of the simulation stepper in a full run of the
// game using the action state values as provided.
func FullDryRun(cfg *GameConfig, actionStateValues []float64) error {
	// Use the simulation generator from the game config
	var gen *simulator.ConfigGenerator = cfg.SimulationGenerator()

	settings, implementations := gen.GenerateConfigs()
	implementations.OutputCondition = &simulator.NilOutputCondition{}

	// Resolve websocket partition indices by name
	for _, name := range cfg.ActionStatePartitionNames {
		for _, iteration := range settings.Iterations {
			if iteration.Name == name {
				iteration.Params.Set("action_state_values", actionStateValues)
			}
		}
	}

	// Run first normally
	coordinator := simulator.NewPartitionCoordinator(
		settings,
		implementations,
	)
	coordinator.Run()

	// Then run with test harnesses
	return simulator.RunWithHarnesses(settings, implementations)
}
