protoc -I=. \
    --go_out=$(pwd) \
    --js_out=./app \
    ./app/partition_state.proto