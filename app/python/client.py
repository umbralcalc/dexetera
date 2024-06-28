import websockets

from typing import Protocol
from websockets.client import WebSocketClientProtocol
from .partition_state_pb2 import PartitionState, State


class ActionTaker(Protocol):
    @property
    def number_of_input_messages() -> int:
        ...
    def take_next_action(
        self, 
        time: float, 
        states: dict[int, list[float]]
    ) -> list[float]:
        ...


async def send(websocket: WebSocketClientProtocol, message: PartitionState):
    binary_message = message.SerializeToString()
    await websocket.send(binary_message)


async def receive(websocket: WebSocketClientProtocol) -> PartitionState:
    binary_message = await websocket.recv()
    message = PartitionState()
    message.ParseFromString(binary_message)
    return message


async def launch_websocket_client(uri: str, action_taker: ActionTaker):
    async with websockets.connect(uri) as websocket:
        received_messages: dict[int, list[float]] = {}
        while len(received_messages) < action_taker.number_of_input_messages:
            message = await receive(websocket)
            received_messages[message.partition_index] = message.state.values
            time = message.cumulative_timesteps

        action_state = action_taker.take_next_action(time, received_messages)

        await send(websocket, State(values=action_state))
