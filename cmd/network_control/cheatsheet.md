# network_control Cheatsheet

This cheat sheet summarises the state keys and action values exchanged with the
`network_control` action server so you can interpret the simulation output and drive the junctions.

## State keys streamed to the Python action server

| Key | Description |
| --- | ----------- |
| `edge_states` | Flat list describing vehicle positions along each edge followed by the cumulative exit count. Each slot stores the distance travelled by a vehicle (≥ 0) or `-1` if empty. Edges are ordered as: `west_entry` (8 slots), `south_entry` (7), `junction_a_to_b` (7), `junction_a_to_north_exit` (6), `north_entry` (6), `junction_b_to_east_exit` (8), `junction_b_to_south_exit` (7), and finally the cumulative exit counter. |
| `vehicle_rectangles` | Convenience projection for the frontend renderer. Values are grouped as `[center_x, center_y, width, height]` for each potential vehicle slot. Width or height equal to zero means the slot is empty. |
| `flow_metrics` | Aggregated metrics: `[vehicles_exited, vehicles_on_network, average_occupancy]`. |
| `junction_a_control` | Current phase of junction A (0 = west→east, 1 = south→north). |
| `junction_b_control` | Current phase of junction B (0 = through-east, 1 = north→south). |

## Action values expected from the Python controller

Send back a list of two floats: `[junction_a_phase, junction_b_phase]`. The values are rounded to the nearest integer and clamped to the available phase count for each junction.

* Junction A phases: `0` (release westbound traffic) or `1` (release southbound traffic).
* Junction B phases: `0` (release vehicles continuing east) or `1` (release north-entry vehicles heading south).

> Tip: Maintaining each phase for a few timesteps avoids oscillations and allows queues to dissipate. The baseline controller provided in `action_server.py` switches only after a minimum dwell time and prefers the approaches with the largest queues.

