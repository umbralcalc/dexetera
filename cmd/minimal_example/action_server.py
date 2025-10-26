from dexact.server import ActionTaker, launch_websocket_server


class MinimalExampleActionTaker(ActionTaker):
    def take_next_action(
        self, 
        time: float,
        states: dict[str, list[float]]
    ) -> list[float]:
        """
        Modify this method to play!
        Also check out the controller cheatsheet if needed:
        umbralcalc.github.io/dexetera/cmd/minimal_example/cheatsheet.md
        """
        print(f"take_next_action called with time: {time}, states: {states}")
        if "counter_state" in states:
            result = [states["counter_state"][0] + 1]
            print(f"Returning: {result}")
            return result
        else:
            print("No counter_state found, returning [1]")
            return [1]  # Default increment


if __name__ == "__main__":
    launch_websocket_server(MinimalExampleActionTaker(), num_state_keys=1)