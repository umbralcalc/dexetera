package game

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// prepareSimulationRun builds settings and implementations from cfg, applies
// action state values to the named action partitions, and uses nil output
// wiring suitable for dry runs and harness checks.
func prepareSimulationRun(
	cfg *GameConfig,
	actionStateValues []float64,
) (*simulator.Settings, *simulator.Implementations) {
	gen := cfg.SimulationGenerator()
	settings, implementations := gen.GenerateConfigs()
	implementations.OutputCondition = &simulator.NilOutputCondition{}
	implementations.OutputFunction = &simulator.NilOutputFunction{}

	for _, name := range cfg.ActionStatePartitionNames {
		for i := range settings.Iterations {
			if settings.Iterations[i].Name == name {
				settings.Iterations[i].Params.Set("action_state_values", actionStateValues)
			}
		}
	}
	return settings, implementations
}

// FullDryRun runs a full coordinator pass with the given action state values.
// Production-style simulation code should use this path (coordinator only),
// not RunWithHarnesses.
func FullDryRun(cfg *GameConfig, actionStateValues []float64) {
	settings, implementations := prepareSimulationRun(cfg, actionStateValues)
	coordinator := simulator.NewPartitionCoordinator(settings, implementations)
	coordinator.Run()
}

// VerifyIterationHarness runs stochadex iteration harness checks on the same
// graph as FullDryRun. Intended for *_test.go only.
func VerifyIterationHarness(cfg *GameConfig, actionStateValues []float64) error {
	settings, implementations := prepareSimulationRun(cfg, actionStateValues)
	return simulator.RunWithHarnesses(settings, implementations)
}
