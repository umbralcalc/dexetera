#!/bin/bash
# Generated build script for minimal_example

echo "Building minimal_example WebAssembly module..."

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Create src directory if it doesn't exist
mkdir -p src

# Build the WebAssembly module from the project root
cd "$PROJECT_ROOT"
GOOS=js GOARCH=wasm go build -o "$SCRIPT_DIR/src/main.wasm" ./cmd/minimal_example/register_step

if [ $? -eq 0 ]; then
    echo "✅ WebAssembly module built successfully!"
    echo "📁 Output: $SCRIPT_DIR/src/main.wasm"
else
    echo "❌ Build failed!"
    exit 1
fi

echo "🎮 minimal_example is ready to run!"
echo "📝 Start your Python websocket server and open index.html in a browser"
