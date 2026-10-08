package dashboard

import (
	"strings"
	"testing"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// widthSettings is one action partition declaring declared action values.
func widthSettings(declared int) *simulator.Settings {
	return &simulator.Settings{Iterations: []simulator.IterationSettings{{
		Name: "policy", Params: simulator.NewParams(map[string][]float64{
			"action_state_values": make([]float64, declared)}),
	}}}
}

func TestCheckActionWidths(t *testing.T) {
	config := func(indices ...int) *Config {
		builder := NewConfigBuilder("test").WithActionStatePartition("policy")
		for i, index := range indices {
			builder.WithSlider(Slider{Name: string(rune('a' + i)), Partition: "policy", ValueIndex: index})
		}
		return builder.Build()
	}
	t.Run("sliders covering the declared width pass", func(t *testing.T) {
		if err := CheckActionWidths(config(0, 2, 1), widthSettings(3)); err != nil {
			t.Error(err)
		}
	})
	t.Run("an action partition with no sliders passes", func(t *testing.T) {
		if err := CheckActionWidths(config(), widthSettings(8)); err != nil {
			t.Error(err)
		}
	})
	t.Run("sliders sending fewer values than declared fail, naming both widths", func(t *testing.T) {
		err := CheckActionWidths(config(0, 5), widthSettings(8))
		if err == nil || !strings.Contains(err.Error(),
			`action partition "policy" declares 8 action_state_values but its sliders send 6 (indices 0..5)`) {
			t.Errorf("expected a width error, got %v", err)
		}
	})
	t.Run("sliders sending more values than declared fail", func(t *testing.T) {
		if err := CheckActionWidths(config(0, 3), widthSettings(2)); err == nil {
			t.Error("expected a width error")
		}
	})
}
