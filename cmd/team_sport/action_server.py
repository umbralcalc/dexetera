from dexact.server import ActionTaker, launch_websocket_server


class TeamSportActionTaker(ActionTaker):
    def __init__(self):
        self.last_substitution_time = 0.0
        self.substitution_cooldown = 10.0  # Can only substitute every 10 seconds

    def take_next_action(
        self,
        time: float,
        states: dict[str, list[float]]
    ) -> list[float]:
        """
        Make substitution decisions based on team stamina and remaining substitutions.
        Returns substitution action: [0.0] = no substitution, [1.0] = substitute
        """
        team_a_stamina = states["team_a_stamina"][0]
        substitutions_remaining = states.get("team_a_substitutions", [3.0])[0]

        # Strategy: If our stamina drops below 50%, we have substitutions left,
        # and we haven't substituted recently, make a substitution
        if (
            team_a_stamina < 50.0
            and substitutions_remaining > 0
            and (time - self.last_substitution_time) >= self.substitution_cooldown
        ):
            self.last_substitution_time = time
            return [1.0]  # Trigger substitution

        return [0.0]


if __name__ == "__main__":
    launch_websocket_server(TeamSportActionTaker(), num_state_keys=7)

