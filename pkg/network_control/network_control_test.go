package network_control

import (
	"testing"

	"github.com/umbralcalc/dexetera/pkg/game"
)

func TestNetworkControl(t *testing.T) {
	t.Run(
		"test that the network control game runs",
		func(t *testing.T) {
			cfg := NewNetworkControlGame().GetConfig()
			game.FullDryRun(cfg, []float64{0.0, 0.0})
			if err := game.VerifyIterationHarness(cfg, []float64{0.0, 0.0}); err != nil {
				t.Errorf("iteration harness: %v", err)
			}
		},
	)

	t.Run(
		"test that control actions work",
		func(t *testing.T) {
			cfg := NewNetworkControlGame().GetConfig()
			game.FullDryRun(cfg, []float64{1.0, 1.0})
			if err := game.VerifyIterationHarness(cfg, []float64{1.0, 1.0}); err != nil {
				t.Errorf("iteration harness: %v", err)
			}
		},
	)
}
