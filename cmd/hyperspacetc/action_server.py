from dexact.server import ActionTaker, launch_websocket_server


class HyperspaceTCActionTaker(ActionTaker):
    @property
    def state_map(self) -> dict[int, str]:
        """You can ignore this config property."""
        return {
            0: "actions",
            1: "process_1",
            2: "process_2",
        }

    def take_next_action(
        self, 
        time: float,
        states: dict[str, list[float]]
    ) -> list[float]:
        """
        Modify this method to play!
        Also check out the controller cheatsheet if needed:
        umbralcalc.github.io/dexetera/cmd/hyperspacetc/cheatsheet.md
        """
        return [s + 0.1 for s in states["actions"]]


if __name__ == "__main__":
    launch_websocket_server(HyperspaceTCActionTaker())