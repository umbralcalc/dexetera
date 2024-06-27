protoc -I=. \
    --go_out=./.. \
    --js_out=. \
    --python_out=./python \
    partition_state.proto