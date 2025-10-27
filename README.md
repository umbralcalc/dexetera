# dexetera...

Decision-making games for the Python programmer built using the [stochadex](https://github.com/umbralcalc/stochadex) simulation framework. Games are Go simulations compiled to WebAssembly and run in the browser with real-time visualizations.

## Quick Start

### 1. Create a Game
Write your game logic in Go using the GameBuilder pattern:

```go
// pkg/games/my_game.go
func NewMyGame() *MyGame {
    visConfig := NewVisualizationBuilder().
        WithCanvas(600, 400).
        AddText("game_state", "Score: {value}", 300, 200, &TextOptions{
            FontSize: 24, Color: "#ffffff",
        }).
        Build()

    config := NewGameBuilder("my_game").
        WithDescription("My awesome game").
        WithPartition("game", "game_state", &MyGameIteration{}).
        WithServerPartition("game_state").
        WithMaxTime(60.0).
        WithVisualization(visConfig).
        Build()

    return &MyGame{config: config}
}
```

### 2. Create WebAssembly Command
```go
// cmd/my_game/main.go
//go:build js && wasm

package main

import (
    "syscall/js"
    "github.com/umbralcalc/dexetera/pkg/games"
    "github.com/umbralcalc/dexetera/pkg/simio"
)

func main() {
    game := games.NewMyGame()
    configGen := game.GetConfigGenerator()
    settings, _ := configGen.GenerateConfigs()
    implementations := game.GetConfig().ImplementationConfig.ToImplementations()
    
    simio.RegisterStep(settings, implementations, 0, "", ":2112")
}
```

### 3. Generate Complete Frontend
```bash
# Generate HTML, CSS, JS, and build scripts
go run ./cmd/generate_game -game my_game -output ./my_game_package
```

### 4. Build and Run
```bash
cd my_game_package
./build.sh                    # Builds WebAssembly
python3 -m http.server 8000   # Start server
# Open http://localhost:8000
```

## Available Games

- `minimal_example` - Simple counter game

## Python Integration

Control games with Python using the [dexAct](https://pypi.org/project/dexact/) package:

```bash
pip install dexact
python cmd/game_name/action_server.py
```
