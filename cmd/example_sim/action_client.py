import asyncio

from app.python.client import ActionTaker, launch_websocket_client


class ExampleActionTaker(ActionTaker):
    @property
    def number_of_input_messages() -> int:
        return 3

    def take_next_action(
        self, 
        time: float, 
        states: dict[int, list[float]]
    ) -> list[float]:
        return [s + 0.1 for s in states[0]]


if __name__ == "__main__":
    uri = "ws://localhost:2112/simio"
    action_taker = ExampleActionTaker()
    asyncio.run(launch_websocket_client(uri, action_taker))