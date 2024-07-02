import asyncio
import websockets

from typing import Protocol
from websockets.server import WebSocketServerProtocol
from .partition_state_pb2 import PartitionState, State


class ActionTaker(Protocol):
    @property
    def input_state_map(self) -> dict[int, str]:
        ...
    def take_next_action(
        self, 
        time: float, 
        states: dict[str, list[float]]
    ) -> list[float]:
        ...


async def _launch_websocket_server(action_taker: ActionTaker):
    received_messages: dict[int, list[float]] = {}

    async def _handle(websocket: WebSocketServerProtocol, path: str):
        async for binary_message in websocket:
            message = PartitionState()
            message.ParseFromString(binary_message)
            received_messages[
                action_taker.input_state_map[message.partition_index]
            ] = message.state.values
            if len(received_messages) == len(action_taker.input_state_map):
                action_state = action_taker.take_next_action(
                    message.cumulative_timesteps, 
                    received_messages,
                )
                await websocket.send(State(values=action_state).SerializeToString())
                received_messages.clear()

    async with websockets.serve(_handle, "localhost", 2112):
        print("WebSocket server started on ws://localhost:2112")
        await asyncio.Future()


def launch_websocket_server(action_taker: ActionTaker):
    asyncio.run(_launch_websocket_server(action_taker))