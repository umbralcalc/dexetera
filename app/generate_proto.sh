protoc -I=. \
    --go_out=./.. \
    --js_out=library=partition_state_pb,binary:. \
    --python_out=./python \
    partition_state.proto