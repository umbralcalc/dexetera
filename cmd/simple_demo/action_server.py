from dexact.server import ActionTaker, launch_websocket_server


class SimpleDemoActionTaker(ActionTaker):
    def take_next_action(
        self, 
        time: float,
        states: dict[str, list[float]]
    ) -> list[float]:
        """
        Simple demo action taker - moves a particle around
        """
        print(f"Simple demo - time: {time}, states: {states}")
        if "particle_state" in states:
            # Simple movement: increment x and y coordinates
            x, y = states["particle_state"][:2]
            new_x = x + 0.1
            new_y = y + 0.05
            result = [new_x, new_y]
            print(f"Moving particle to: {result}")
            return result
        else:
            print("No particle_state found, returning default position")
            return [0.0, 0.0]  # Default position


if __name__ == "__main__":
    launch_websocket_server(SimpleDemoActionTaker(), num_state_keys=2)
