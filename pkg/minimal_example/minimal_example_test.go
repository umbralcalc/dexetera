package minimal_example

import (
	"testing"

	"github.com/umbralcalc/dexetera/pkg/game"
)

func TestMinimalExample(t *testing.T) {
	t.Run(
		"test that the minimal example game runs",
		func(t *testing.T) {
			if err := game.FullDryRun(
				NewMinimalExampleGame().GetConfig(),
				[]float64{1.0},
			); err != nil {
				t.Errorf("test harness failed: %v", err)
			}
		},
	)
}
