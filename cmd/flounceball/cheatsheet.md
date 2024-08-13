# Flounceball Manager Cheatsheet

## List value meanings for `states["latest_manager_actions"]`

## List value meanings for `states["latest_match_state"]`

0. is a [0,1] indicator of the 'restart state'
1. is a [0,1] indicator of the 'possession state', where 0 means your team is in possession
2. is your team's total Air Time score (see the rules for the definition)
3. is the opposition team's total Air Time score (see the rules for the definition)
4. is the cumulative Air Time for the ball in current possession
5. is the (horizontal) speed of the ball
6. is the radial position of the ball on the pitch
7. is the angular position of the ball on the pitch
8. is the projected radial position that the ball is heading towards
9. is the projected angular position that the ball is heading towards
10. is the cumulative time that the ball has spent falling once it reaches the end of its path
