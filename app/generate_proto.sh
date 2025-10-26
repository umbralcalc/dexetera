protoc -I=. \
    --go_out=./.. \
    --js_out=library=./src/partition_state_pb,binary:. \
    partition_state.proto;
# protoc --python_out=../../dexact/dexact partition_state.proto;