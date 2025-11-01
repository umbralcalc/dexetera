# dexetera...

Decision-making games for the Python programmer built using the [stochadex](https://github.com/umbralcalc/stochadex) simulation framework. Games are Go simulations compiled to WebAssembly and run in the browser with real-time visualizations.

## Quick Start

### Build minimal_example
```bash
# Generate complete game package
go run ./cmd/minimal_example/generate_game

# Build and run
cd minimal_example
./build.sh                    # Builds WebAssembly
python -m http.server 8000   # Start server
# Open http://localhost:8000
```

### Build team_sport
```bash
# Generate complete game package
go run ./cmd/team_sport/generate_game

# Build and run
cd team_sport
./build.sh                    # Builds WebAssembly
python -m http.server 8000   # Start server
# Open http://localhost:8000
```

## Available Games

- `minimal_example` - Simple counter game (demonstrates the workflow)
- `team_sport` - Team sport manager game where you make strategic substitutions to win!

## Python Integration

The Python server drives the simulation - the frontend only iterates the world state in reponse to the server action states and displays the results. You can control games via action states with Python using the [dexAct](https://pypi.org/project/dexact/) package:

```bash
pip install dexact
python cmd/minimal_example/action_server.py
```

**Note**: Since the Python action server controls the simulation timing and execution, you can add sleeps into the `ActionTaker.take_next_action` call to slow the simulation frontend updates down to make them more human-friendly.
