# dexetera...

Decision-making games for the Python programmer built using the [stochadex](https://github.com/umbralcalc/stochadex) simulation framework. Games are Go simulations compiled to WebAssembly and run in the browser with real-time visualizations.

## Quick Start

### Build
```bash
# Generate complete game package
go run ./cmd/minimal_example/generate_game

# Build and run
cd minimal_example
./build.sh                    # Builds WebAssembly
python -m http.server 8000   # Start server
# Open http://localhost:8000
```

## Available Games

- `minimal_example` - Simple counter game (demonstrates the new workflow)

## Python Integration

The Python server drives the simulation - the frontend only displays the results. Control games with Python using the [dexAct](https://pypi.org/project/dexact/) package:

```bash
pip install dexact
python cmd/game_name/action_server.py
```

**Note**: This Python action server controls the simulation timing and execution. For example, you can add sleeps into the `ActionTaker.take_next_action` call to slow things down.
