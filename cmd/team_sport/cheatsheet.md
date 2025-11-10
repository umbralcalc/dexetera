# Team Sport Game Cheatsheet

## Game Overview
This game simulates a team sport match where you manage team substitutions. Your goal is to strategically substitute players to maintain high stamina and outscore the opponent. Each team has 11 players on the field and 3 substitutions available.

## State Partition Meanings

### `states["score"]`
- **Index 0**: Current score differential (Team A - Team B)
  - Positive values mean Team A is winning
  - Negative values mean Team B is winning
- **Index 1**: Goal flag from the previous step
  - `1.0` indicates Team A scored on the last step
  - `-1.0` indicates Team B scored on the last step
  - `0.0` means no recent goal (players keep moving)

### `states["team_a_players"]`
- Contains the `(x, y)` positions for Team A players
- Even indices (0, 2, 4, ...) are x-positions
- Odd indices (1, 3, 5, ...) are y-positions
- Players start near the left of the field and advance toward the right goal

### `states["team_b_players"]`
- Contains the `(x, y)` positions for Team B players
- Even indices are x-positions, odd indices are y-positions
- Players start near the right of the field and advance toward the left goal

### `states["team_a_stamina"]`
- **Index 0**: Average stamina percentage of Team A (0-100)
  - Starts at 80%
  - Decreases by 0.5% per second naturally
  - Can be restored to base level (80%) by making a substitution (if substitutions remain)

### `states["team_b_stamina"]`
- **Index 0**: Average stamina percentage of Team B (0-100)
  - Starts at 80%
  - Decreases by 0.5% per second naturally

### `states["team_a_substitutions"]`
- **Index 0**: Number of substitutions remaining for Team A (0-3)
  - Starts at 3
  - Decrements by 1 each time you make a substitution
  - Substitutions only work if this value is greater than 0

### `states["team_b_substitutions"]`
- **Index 0**: Number of substitutions remaining for Team B (0-3)
  - Starts at 3
  - Team B doesn't make substitutions (controlled by AI)

## Action State

### `action_state_values`
Your action should be a single-element list:
- **`[0.0]`**: No substitution
- **`[1.0]`**: Make a substitution (restores Team A stamina to 80%, decrements substitution count)

## Game Mechanics

- **Players on field**: Each team has 11 players whose positions update each step
- **Movement speed**: Players move faster when stamina is high; low stamina slows them down
- **Substitutions available**: Each team starts with 3 substitutions
- **Stamina decay**: Both teams lose 0.5% stamina per second
- **Substitutions**: Restore Team A stamina to 80%, but only if you have substitutions remaining
- **Scoring**: Team A scores when one of their players reaches the opponent's goal line (x ≈ 710); Team B scores when reaching x ≈ 90
- **Goal flag**: The second value in `score` tells you who scored in the previous step so you can react or display effects
- **Match duration**: The simulation keeps running (long horizon) so your strategy can be tested over many possessions
- **Strategy**: Make timely substitutions to keep stamina high, but use them wisely - you only have 3!

## Tips

1. Monitor your stamina closely - let it drop too low and you'll fall behind
2. Substitutions are powerful but limited - you only have 3, so use them strategically
3. Check `team_a_substitutions` before making a substitution - if it's 0, substitutions won't work
4. If your score is negative, you're losing - consider making a substitution if you have any left
5. The opponent (Team B) doesn't substitute, but their players still move with their stamina baseline
6. Save at least one substitution for the final minutes if possible to counter fatigue

