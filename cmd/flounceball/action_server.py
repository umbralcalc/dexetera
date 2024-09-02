import math, random, itertools # libs only needed for example
from dexact.server import ActionTaker, launch_websocket_server


class FlounceballActionTaker(ActionTaker):
    @property
    def state_map(self) -> dict[int, str]:
        """You can ignore this config property."""
        return {
            0: "latest_manager_actions",
            22: "latest_match_state",
        }

    def take_next_action(
        self, 
        time: float,
        states: dict[str, list[float]]
    ) -> list[float]:
        """
        Modify this method to play!
        Also check out the manager cheatsheet if needed:
        umbralcalc.github.io/dexetera/cmd/flounceball/cheatsheet.md
        """
        return list(itertools.chain(*[
            (random.uniform(0, 100), random.uniform(0, 2*math.pi)) 
            for _ in range(0, 10)
        ]))


if __name__ == "__main__":
    launch_websocket_server(FlounceballActionTaker())