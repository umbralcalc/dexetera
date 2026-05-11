package simio

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// ApplyActionState writes incoming action values onto the running
// coordinator's iterator params.
//
// When actionState.Partitions is non-empty, the per-partition named path
// is used: each entry sets `action_state_values` only on the partition
// whose name matches the map key (looked up via actionPartitionIndexByName).
// Names not present in the map are silently skipped.
//
// When actionState.Partitions is empty, the legacy broadcast path applies:
// the same actionState.Values slice is set on every partition listed in
// actionPartitionIndices. This preserves compatibility with action sources
// that don't yet emit named partitions (e.g. existing dexact Python clients).
func ApplyActionState(
	coordinator *simulator.PartitionCoordinator,
	actionPartitionIndices []int,
	actionPartitionIndexByName map[string]int,
	actionState *ActionState,
) {
	if actionState == nil {
		return
	}
	if len(actionState.Partitions) > 0 {
		for name, av := range actionState.Partitions {
			index, ok := actionPartitionIndexByName[name]
			if !ok {
				continue
			}
			coordinator.Iterators[index].
				Params.Set("action_state_values", av.GetValues())
		}
		return
	}
	for _, index := range actionPartitionIndices {
		coordinator.Iterators[index].
			Params.Set("action_state_values", actionState.Values)
	}
}
