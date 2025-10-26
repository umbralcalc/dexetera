protoc -I=. \
    --go_out=./.. \
    --js_out=library=./src/action_state_pb,binary:. \
    action_state.proto;
protoc --python_out=../../dexact/dexact action_state.proto;