# Hyperspace Traffic Controller Cheatsheet

## List value meanings for `states["latest_controller_actions"]`

0. is the categorical control input for the left jump point (choose to serve queues from: [4, 5, 6] or anything else to serve none of them)
1. is the categorical control input for the upper jump point (choose to serve queues from: [8, 9] or anything else to serve none of them)
2. is the categorical control input for the right jump point (choose to serve queues from: [11, 12] or anything else to serve none of them)

## List value meanings for `states["latest_left_jump_point_queues"]`

0. is the size of the middle outside triangle queue (served with ID 4)
1. is the size of the upper outside triangle queue (served with ID 5)
2. is the size of the lower outside triangle queue (served with ID 6)

## List value meanings for `states["latest_upper_jump_point_queues"]`

0. is the size of the outside triangle queue (served with ID 8)
1. is the size of the inside triangle queue (served with ID 9)


## List value meanings for `states["latest_right_jump_point_queues"]`

0. is the size of the lower inside triangle queue (served with ID 11)
1. is the size of the upper inside triangle queue (served with ID 12)
