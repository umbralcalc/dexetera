from __future__ import annotations

from dataclasses import dataclass
from typing import Dict, List

from dexact.server import ActionTaker, launch_websocket_server


@dataclass
class EdgeLayout:
    name: str
    capacity: int
    start_index: int


EDGE_LAYOUT: List[EdgeLayout] = [
    EdgeLayout("west_entry", capacity=8, start_index=0),
    EdgeLayout("south_entry", capacity=7, start_index=8),
    EdgeLayout("junction_a_to_b", capacity=7, start_index=15),
    EdgeLayout("junction_a_to_north_exit", capacity=6, start_index=22),
    EdgeLayout("north_entry", capacity=6, start_index=28),
    EdgeLayout("junction_b_to_east_exit", capacity=8, start_index=34),
    EdgeLayout("junction_b_to_south_exit", capacity=7, start_index=42),
]


class NetworkControlActionTaker(ActionTaker):
    """
    Baseline controller that cycles traffic light phases while favouring
    the approaches with the longest queues.
    """

    def __init__(self) -> None:
        self.last_switch_time = 0.0
        self.min_phase_duration = 6.0
        self.junction_a_phase = 0
        self.junction_b_phase = 0

    def take_next_action(
        self,
        time: float,
        states: Dict[str, List[float]],
    ) -> List[float]:
        edge_state = states.get("edge_states", [])
        queues = self._estimate_queues(edge_state)

        if queues:
            # Junction A controls west vs south approach (edge 0 vs edge 1)
            if queues.get("south_entry", 0.0) > queues.get("west_entry", 0.0) * 1.2:
                desired_a = 1
            elif queues.get("south_entry", 0.0) * 1.2 < queues.get("west_entry", 0.0):
                desired_a = 0
            else:
                desired_a = self.junction_a_phase

            # Junction B controls through movement vs north entry (edge 2 vs edge 4)
            if queues.get("north_entry", 0.0) > queues.get("junction_a_to_b", 0.0) * 1.2:
                desired_b = 1
            elif queues.get("north_entry", 0.0) * 1.2 < queues.get("junction_a_to_b", 0.0):
                desired_b = 0
            else:
                desired_b = self.junction_b_phase
        else:
            desired_a = self.junction_a_phase
            desired_b = self.junction_b_phase

        if time - self.last_switch_time >= self.min_phase_duration:
            if desired_a != self.junction_a_phase:
                self.junction_a_phase = desired_a
                self.last_switch_time = time
            if desired_b != self.junction_b_phase:
                self.junction_b_phase = desired_b
                self.last_switch_time = time

        return [float(self.junction_a_phase), float(self.junction_b_phase)]

    def _estimate_queues(self, edge_state: List[float]) -> Dict[str, float]:
        queues: Dict[str, float] = {}
        for layout in EDGE_LAYOUT:
            total = 0.0
            for slot in range(layout.capacity):
                idx = layout.start_index + slot
                if idx >= len(edge_state):
                    break
                if edge_state[idx] >= 0.0:
                    total += 1.0
            queues[layout.name] = total
        return queues


if __name__ == "__main__":
    launch_websocket_server(NetworkControlActionTaker(), num_state_keys=5)

