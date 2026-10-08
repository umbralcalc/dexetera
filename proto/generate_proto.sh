#!/usr/bin/env bash
# Regenerate the JS stub for ActionState, the action message the browser
# runtime sends into a simulation.
#
# The message is stochadex's: its definition lives in stochadex's
# cmd/messages/action_state.proto (wire-compatible with the one dexetera used
# to keep here), and its Go bindings are stochadex's simulator.ActionState,
# which pkg/simio aliases. This script reads the .proto from the stochadex
# version in go.mod, so the JS always matches the Go the wasm module is built
# with. Output lands in <repo>/runtime/ alongside the rest of the runtime.
#
# Python stubs for the dexact package are not regenerated here; if you maintain
# dexact alongside this repo, regenerate them from stochadex's
# cmd/messages/action_state_pb2.py.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STOCHADEX="$(cd "$HERE/.." && go list -m -f '{{.Dir}}' github.com/umbralcalc/stochadex)"
cd "$STOCHADEX/cmd/messages"

protoc -I=. --js_out=library=action_state_pb,binary:"$HERE/../runtime" action_state.proto
