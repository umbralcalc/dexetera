import json
import asyncio
import websockets

from typing import Protocol
from websockets.server import WebSocketServerProtocol
from .partition_state_pb2 import PartitionState, State


class ActionTaker(Protocol):
    @property
    def number_of_input_messages(self) -> int:
        ...
    def take_next_action(
        self, 
        time: float, 
        states: dict[int, list[float]]
    ) -> list[float]:
        ...


async def _launch_websocket_server(action_taker: ActionTaker):
    async def _handle(websocket: WebSocketServerProtocol, path: str):
        received_messages: dict[int, list[float]] = {}
        async for binary_message in websocket:
            message = PartitionState()
            message.ParseFromString(binary_message.data)
            received_messages[message.partition_index] = message.state.values
            time = message.cumulative_timesteps
            if len(received_messages) == action_taker.number_of_input_messages:
                action_state = action_taker.take_next_action(time, received_messages)
                await websocket.send(
                    json.dumps({"data" : State(values=action_state).SerializeToString()})
                )
                received_messages.clear()

    async with websockets.serve(_handle, "localhost", 2112):
        print("WebSocket server started on ws://localhost:2112")
        await asyncio.Future()


def launch_websocket_server(action_taker: ActionTaker):
    asyncio.run(_launch_websocket_server(action_taker))