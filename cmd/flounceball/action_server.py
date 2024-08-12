from dexact.server import ActionTaker, launch_websocket_server


class ExampleActionTaker(ActionTaker):
    @property
    def state_map(self) -> dict[int, str]:
        return {
            0: "latest_manager_actions",
            21: "latest_match_state",
        }

    def take_next_action(
        self, 
        time: float,
        states: dict[str, list[float]]
    ) -> list[float]:
        return [s + 0.1 for s in states["latest_manager_actions"]]


if __name__ == "__main__":
    launch_websocket_server(ExampleActionTaker())