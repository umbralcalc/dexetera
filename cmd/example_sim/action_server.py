from app.python.server import ActionTaker, launch_websocket_server


class ExampleActionTaker(ActionTaker):
    @property
    def state_map(self) -> dict[int, str]:
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
        return [s + 0.1 for s in states["actions"]]


if __name__ == "__main__":
    launch_websocket_server(ExampleActionTaker())