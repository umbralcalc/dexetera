package team_sport

import (
	"testing"

	"github.com/umbralcalc/dexetera/pkg/game"
)

func TestTeamSport(t *testing.T) {
	t.Run(
		"test that the team sport game runs",
		func(t *testing.T) {
			cfg := NewTeamSportGame().GetConfig()
			game.FullDryRun(cfg, []float64{0.0}) // No substitution initially
			if err := game.VerifyIterationHarness(cfg, []float64{0.0}); err != nil {
				t.Errorf("iteration harness: %v", err)
			}
		},
	)

	t.Run(
		"test that substitutions work",
		func(t *testing.T) {
			cfg := NewTeamSportGame().GetConfig()
			game.FullDryRun(cfg, []float64{1.0}) // Make a substitution
			if err := game.VerifyIterationHarness(cfg, []float64{1.0}); err != nil {
				t.Errorf("iteration harness: %v", err)
			}
		},
	)
}
