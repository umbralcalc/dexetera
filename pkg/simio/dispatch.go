package simio

import (
	"fmt"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// ActionState and ActionValues are stochadex's action message
// (simulator.ActionState, generated from stochadex's
// cmd/messages/action_state.proto). The wire format is unchanged, so
// existing clients and runtime/action_state_pb.js keep working.
type (
	ActionState  = simulator.ActionState
	ActionValues = simulator.ActionValues
)

// ActionParam is the params key an action partition reads its actions from.
const ActionParam = "action_state_values"

// ActionDispatcher writes incoming action values onto a running
// coordinator's action partitions, between steps, through stochadex's
// ParamsInjector: each partition is found, and its action_state_values
// checked, once, when the dispatcher is made, so an action is an in-place
// copy with no allocation.
//
// When an ActionState's Partitions is non-empty, the per-partition named
// path is used: each entry sets the partition whose name matches the key.
// Names that are not action partitions are skipped. When it is empty, the
// broadcast path applies: Values is set on every action partition. This
// keeps compatibility with action sources that don't emit named partitions
// (e.g. existing dexact Python clients).
type ActionDispatcher struct {
	ordered []*simulator.ParamsInjector // broadcast path, declaration order
	byName  map[string]*simulator.ParamsInjector
}

// NewActionDispatcher resolves each named action partition of coordinator.
// A partition that does not exist, or does not declare action_state_values,
// is an error.
func NewActionDispatcher(
	coordinator *simulator.PartitionCoordinator,
	partitionNames []string,
) (*ActionDispatcher, error) {
	dispatcher := &ActionDispatcher{
		ordered: make([]*simulator.ParamsInjector, 0, len(partitionNames)),
		byName:  make(map[string]*simulator.ParamsInjector, len(partitionNames)),
	}
	for _, name := range partitionNames {
		injector, err := coordinator.NewParamsInjector(name, ActionParam)
		if err != nil {
			return nil, fmt.Errorf("simio: action partition %q: %w", name, err)
		}
		dispatcher.ordered = append(dispatcher.ordered, injector)
		dispatcher.byName[name] = injector
	}
	return dispatcher, nil
}

// Apply sets the actions an ActionState carries, from the next step on. Call
// it only between steps. A vector of the wrong width is an error, and leaves
// that partition's actions as they were.
func (d *ActionDispatcher) Apply(actionState *ActionState) error {
	if actionState == nil {
		return nil
	}
	if len(actionState.Partitions) > 0 {
		for name, values := range actionState.Partitions {
			injector, ok := d.byName[name]
			if !ok {
				continue
			}
			if err := injector.Inject(values.GetValues()); err != nil {
				return err
			}
		}
		return nil
	}
	for _, injector := range d.ordered {
		if err := injector.Inject(actionState.Values); err != nil {
			return err
		}
	}
	return nil
}
