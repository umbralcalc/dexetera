from dexact.server import ActionTaker, launch_websocket_server
import random


class ParticleSystemActionTaker(ActionTaker):
    def __init__(self):
        super().__init__()
        self.particle_count = 10
    
    def take_next_action(
        self, 
        time: float,
        states: dict[str, list[float]]
    ) -> list[float]:
        """
        Particle system action taker - moves multiple particles
        """
        print(f"Particle system - time: {time}, states: {states}")
        if "particle_state" in states:
            # Get current particle positions
            current_state = states["particle_state"]
            new_state = []
            
            # Move each particle (x, y, vx, vy)
            for i in range(0, len(current_state), 4):
                x, y, vx, vy = current_state[i:i+4]
                
                # Update position based on velocity
                new_x = x + vx * 0.016  # 60 FPS
                new_y = y + vy * 0.016
                
                # Add some random movement
                new_vx = vx + (random.random() - 0.5) * 0.1
                new_vy = vy + (random.random() - 0.5) * 0.1
                
                # Keep particles in bounds
                if new_x < 0 or new_x > 1:
                    new_vx = -new_vx
                    new_x = max(0, min(1, new_x))
                if new_y < 0 or new_y > 1:
                    new_vy = -new_vy
                    new_y = max(0, min(1, new_y))
                
                new_state.extend([new_x, new_y, new_vx, new_vy])
            
            print(f"Updated {len(new_state)//4} particles")
            return new_state
        else:
            print("No particle_state found, returning default particles")
            # Return default particle positions
            default_state = []
            for i in range(self.particle_count):
                default_state.extend([
                    random.random(),  # x
                    random.random(),  # y
                    (random.random() - 0.5) * 0.1,  # vx
                    (random.random() - 0.5) * 0.1   # vy
                ])
            return default_state


if __name__ == "__main__":
    launch_websocket_server(ParticleSystemActionTaker(), num_state_keys=40)  # 10 particles * 4 values each
