# Team Sport Game Cheatsheet

## Game Overview
This game simulates a team sport match where you manage team substitutions. Your goal is to strategically substitute players to maintain high stamina and outscore the opponent.

## State Partition Meanings

### `states["score"]`
- **Index 0**: Current score differential (Team A - Team B)
  - Positive values mean Team A is winning
  - Negative values mean Team B is winning
  - Score increases when your team has more stamina than the opponent

### `states["team_a_stamina"]`
- **Index 0**: Average stamina percentage of Team A (0-100)
  - Starts at 80%
  - Decreases by 0.5% per second naturally
  - Can be restored to base level (80%) by making a substitution

### `states["team_b_stamina"]`
- **Index 0**: Average stamina percentage of Team B (0-100)
  - Starts at 75%
  - Decreases by 0.5% per second naturally

## Action State

### `action_state_values`
Your action should be a single-element list:
- **`[0.0]`**: No substitution
- **`[1.0]`**: Make a substitution (restores Team A stamina to 80%)

## Game Mechanics

- **Stamina decay**: Both teams lose 0.5% stamina per second
- **Substitutions**: Restore your team's stamina to 80%
- **Scoring**: The team with higher stamina scores more points
- **Match duration**: 90 seconds
- **Strategy**: Make timely substitutions to keep stamina high and score points!

## Tips

1. Monitor your stamina closely - let it drop too low and you'll fall behind
2. Substitutions are powerful but you need to time them well
3. If your score is negative, you're losing - make a substitution!
4. The opponent (Team B) starts slightly weaker but doesn't substitute

