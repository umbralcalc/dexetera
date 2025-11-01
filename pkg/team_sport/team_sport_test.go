package team_sport

import (
	"testing"

	"github.com/umbralcalc/dexetera/pkg/game"
)

func TestTeamSport(t *testing.T) {
	t.Run(
		"test that the team sport game runs",
		func(t *testing.T) {
			if err := game.FullDryRun(
				NewTeamSportGame().GetConfig(),
				[]float64{0.0}, // No substitution initially
			); err != nil {
				t.Errorf("test harness failed: %v", err)
			}
		},
	)

	t.Run(
		"test that substitutions work",
		func(t *testing.T) {
			if err := game.FullDryRun(
				NewTeamSportGame().GetConfig(),
				[]float64{1.0}, // Make a substitution
			); err != nil {
				t.Errorf("test harness failed: %v", err)
			}
		},
	)
}
