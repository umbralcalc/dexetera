from dexact.server import ActionTaker, launch_websocket_server
import time


class TeamSportActionTaker(ActionTaker):
    def __init__(self):
        super().__init__()
        self.last_substitution_time = 0.0
        self.substitution_cooldown = 10.0  # Can only substitute every 10 seconds
    
    def take_next_action(
        self, 
        time_val: float,
        states: dict[str, list[float]]
    ) -> list[float]:
        """
        Make substitution decisions based on team stamina and remaining substitutions.
        Returns substitution action: [0.0] = no substitution, [1.0] = substitute
        """
        print(f"Time: {time_val:.1f}s")
        
        if "team_a_stamina" in states and "team_b_stamina" in states:
            team_a_stamina = states["team_a_stamina"][0]
            team_b_stamina = states["team_b_stamina"][0]
            score = states.get("score", [0.0])[0]
            substitutions_remaining = states.get("team_a_substitutions", [3.0])[0]
            
            print(f"  Team A Stamina: {team_a_stamina:.1f}%")
            print(f"  Team B Stamina: {team_b_stamina:.1f}%")
            print(f"  Score: {score:.1f}")
            print(f"  Substitutions Remaining: {substitutions_remaining:.0f}")
            
            # Strategy: If our stamina drops below 50%, we have substitutions left,
            # and we haven't substituted recently, make a substitution
            if (team_a_stamina < 50.0 and 
                substitutions_remaining > 0 and
                (time_val - self.last_substitution_time) >= self.substitution_cooldown):
                print(f"  ⚽ MAKING SUBSTITUTION!")
                self.last_substitution_time = time_val
                return [1.0]  # Trigger substitution
            elif substitutions_remaining <= 0:
                print(f"  ⚠️  No substitutions remaining!")
        
        print(f"  No substitution")
        return [0.0]  # No substitution


if __name__ == "__main__":
    launch_websocket_server(TeamSportActionTaker(), num_state_keys=1)

