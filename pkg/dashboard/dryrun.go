package dashboard

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// prepareSimulationRun materialises a stochadex Settings + Implementations
// pair from cfg with output wiring stubbed out, then seeds every action
// partition with actionStateValues so subsequent steps see those values as
// if they had arrived from a live action source. The two return values are
// the same pair that PartitionCoordinator and RunWithHarnesses both want.
//
// Output is set to the nil pair on purpose: dry runs and harness checks
// don't want noise on stdout, they just want the iteration graph to
// execute.
func prepareSimulationRun(
	cfg *Config,
	actionStateValues []float64,
) (*simulator.Settings, *simulator.Implementations) {
	settings, implementations := cfg.SimulationGenerator().GenerateConfigs()
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

// FullDryRun executes one complete coordinator run against the simulation
// described by cfg, with every action partition pre-seeded to
// actionStateValues. Production code paths use this entry point (a real
// PartitionCoordinator, not the harness machinery); test code should
// prefer VerifyIterationHarness for stronger checks on each iteration.
func FullDryRun(cfg *Config, actionStateValues []float64) {
	settings, implementations := prepareSimulationRun(cfg, actionStateValues)
	coordinator := simulator.NewPartitionCoordinator(settings, implementations)
	coordinator.Run()
}

// VerifyIterationHarness runs stochadex's iteration-harness validation
// over the same partition graph FullDryRun would execute. The harness
// confirms each partition's Iterate implementation respects the contracts
// stochadex relies on (no out-of-bounds history reads, no state-width
// surprises, etc.). Intended for *_test.go usage.
func VerifyIterationHarness(cfg *Config, actionStateValues []float64) error {
	settings, implementations := prepareSimulationRun(cfg, actionStateValues)
	return simulator.RunWithHarnesses(settings, implementations)
}
