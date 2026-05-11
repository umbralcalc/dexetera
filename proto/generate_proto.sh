#!/usr/bin/env bash
# Regenerate Go and JS stubs from action_state.proto.
#   - Go output is routed by the proto's `option go_package = "./pkg/simio";`
#     and lands in <repo>/pkg/simio/.
#   - JS output lands in <repo>/runtime/ alongside the rest of the runtime.
# Python stubs for the dexact package are not regenerated here; if you maintain
# dexact alongside this repo, regenerate them from there.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$HERE"

protoc -I=. --go_out=./.. action_state.proto
protoc -I=. --js_out=library=action_state_pb,binary:../runtime action_state.proto
