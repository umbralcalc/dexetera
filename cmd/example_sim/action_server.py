from app.python.server import ActionTaker, launch_websocket_server


class ExampleActionTaker(ActionTaker):
    @property
    def number_of_input_messages(self) -> int:
        return 3

    def take_next_action(
        self, 
        time: float, 
        states: dict[int, list[float]]
    ) -> list[float]:
        return [s + 0.1 for s in states[0]]


if __name__ == "__main__":
    launch_websocket_server(ExampleActionTaker())