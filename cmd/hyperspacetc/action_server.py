from random import choice # lib only needed for example
from dexact.server import ActionTaker, launch_websocket_server


class HyperspaceTCActionTaker(ActionTaker):
    @property
    def state_map(self) -> dict[int, str]:
        """You can ignore this config property."""
        return {
            0: "latest_controller_actions",
            1: "latest_left_jump_point_queues",
            2: "latest_upper_jump_point_queues",
            3: "latest_right_jump_point_queues",
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
        # print(
        #     states["latest_left_jump_point_queues"],
        #     states["latest_upper_jump_point_queues"],
        #     states["latest_right_jump_point_queues"],
        # )
        return [choice([4, 5, 6]), choice([8, 9]), choice([11, 12])]


if __name__ == "__main__":
    launch_websocket_server(HyperspaceTCActionTaker())