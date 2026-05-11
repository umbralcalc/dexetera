package growth

import (
	"testing"

	"github.com/umbralcalc/dexetera/pkg/dashboard"
)

func TestGrowth(t *testing.T) {
	t.Run("steady-state action values run end-to-end", func(t *testing.T) {
		cfg := NewConfig()
		dashboard.FullDryRun(cfg, []float64{0.05, 500.0})
		if err := dashboard.VerifyIterationHarness(cfg, []float64{0.05, 500.0}); err != nil {
			t.Errorf("iteration harness: %v", err)
		}
	})

	t.Run("zero carrying capacity collapses safely", func(t *testing.T) {
		cfg := NewConfig()
		// K = 0 is the edge case the iteration explicitly handles.
		dashboard.FullDryRun(cfg, []float64{0.05, 0.0})
		if err := dashboard.VerifyIterationHarness(cfg, []float64{0.05, 0.0}); err != nil {
			t.Errorf("iteration harness: %v", err)
		}
	})
}
