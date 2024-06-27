import asyncio
import websockets

from websockets.client import WebSocketClientProtocol
from typing import Any, Protocol
from partition_state_pb2 import PartitionState


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


async def send(websocket: WebSocketClientProtocol, message: Any):
    binary_message = message.SerializeToString()
    await websocket.send(binary_message)


async def receive(websocket: WebSocketClientProtocol) -> Any:
    binary_message = await websocket.recv()
    message = PartitionState()
    message.ParseFromString(binary_message)
    return message


async def websocket_client(uri: str, action_taker: ActionTaker):
    async with websockets.connect(uri) as websocket:
        received_messages: dict[int, list[float]] = {}
        while len(received_messages) < action_taker.number_of_input_messages:
            message = await receive(websocket)
            received_messages[message.partition_index] = message.state
            time = message.cumulative_timesteps

        action_state = action_taker.take_next_action(time, received_messages)

        # maybe don't need the full protobuf message type to send back?
        message_to_send = list(received_messages)[0]
        message_to_send.state = action_state

        await send(websocket, message_to_send)

if __name__ == "__main__":
    uri = "ws://localhost:2112/simio"
    asyncio.run(websocket_client(uri))
