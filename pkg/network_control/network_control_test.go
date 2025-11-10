package network_control

import (
	"testing"

	"github.com/umbralcalc/dexetera/pkg/game"
)

func TestTeamSport(t *testing.T) {
	t.Run(
		"test that the team sport game runs",
		func(t *testing.T) {
			if err := game.FullDryRun(
				NewNetworkControlGame().GetConfig(),
				[]float64{0.0, 0.0},
			); err != nil {
				t.Errorf("test harness failed: %v", err)
			}
		},
	)

	t.Run(
		"test that substitutions work",
		func(t *testing.T) {
			if err := game.FullDryRun(
				NewNetworkControlGame().GetConfig(),
				[]float64{1.0, 1.0},
			); err != nil {
				t.Errorf("test harness failed: %v", err)
			}
		},
	)
}
