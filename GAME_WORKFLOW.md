# 🎮 Complete Dexetera Game Development Workflow

This document explains the complete end-to-end workflow for creating and deploying Dexetera games.

## 📋 Overview

The Dexetera framework provides a complete workflow from Go game code to deployed web applications:

1. **Write Game Code** - Create game logic in Go using the GameBuilder pattern
2. **Generate Frontend** - Use the template system to create HTML/JS/CSS
3. **Build WebAssembly** - Compile Go code to WASM
4. **Deploy** - Run the complete game in a browser

## 🚀 Step-by-Step Workflow

### Step 1: Create a New Game

Create a new game file in `pkg/games/`:

```go
// pkg/games/my_new_game.go
package games

import (
    "github.com/umbralcalc/stochadex/pkg/simulator"
)

type MyNewGame struct {
    config *GameConfig
}

func NewMyNewGame() *MyNewGame {
    // Create visualization using VisualizationBuilder
    visConfig := NewVisualizationBuilder().
        WithCanvas(600, 400).
        WithBackground("#1a1a1a").
        WithUpdateInterval(50).
        AddText("game_state", "Score: {value}", 300, 200, &TextOptions{
            FontSize: 24,
            Color:    "#ffffff",
        }).
        Build()

    // Create the game using GameBuilder
    config := NewGameBuilder("my_new_game").
        WithDescription("My awesome new game").
        WithPartition("game", "game_state", &MyGameIteration{}).
        WithServerPartition("game_state").
        WithParameter("game_init", []float64{0.0}).
        WithParameter("game_params", map[string][]float64{"speed": {1.0}}).
        WithMaxTime(60.0).
        WithTimestep(1.0).
        WithVisualization(visConfig).
        Build()

    return &MyNewGame{config: config}
}

// Implement the Game interface
func (g *MyNewGame) GetName() string { return g.config.Name }
func (g *MyNewGame) GetDescription() string { return g.config.Description }
func (g *MyNewGame) GetConfig() *GameConfig { return g.config }
func (g *MyNewGame) GetConfigGenerator() *simulator.ConfigGenerator {
    // Implementation here
}
func (g *MyNewGame) GetRenderer() GameRenderer {
    return &GenericRenderer{config: g.config.VisualizationConfig}
}

// Your game iteration logic
type MyGameIteration struct{}

func (m *MyGameIteration) Configure(partitionIndex int, settings *simulator.Settings) {}
func (m *MyGameIteration) Iterate(params *simulator.Params, partitionIndex int, 
    stateHistories []*simulator.StateHistory, timestepsHistory *simulator.CumulativeTimestepsHistory) []float64 {
    // Your game logic here
    outputState := stateHistories[partitionIndex].Values.RawRowView(0)
    speed := params.Get("speed")[0]
    outputState[0] += speed
    return outputState
}
```

### Step 2: Create WebAssembly Command

Create a command file in `cmd/my_new_game/main.go`:

```go
//go:build js && wasm

package main

import (
    "syscall/js"
    "github.com/umbralcalc/dexetera/pkg/games"
    "github.com/umbralcalc/dexetera/pkg/simio"
)

func main() {
    js.Global().Get("console").Call("log", "My new game main function called")

    game := games.NewMyNewGame()
    js.Global().Get("console").Call("log", "MyNewGame created")

    configGen := game.GetConfigGenerator()
    settings, _ := configGen.GenerateConfigs()
    
    implementations := game.GetConfig().ImplementationConfig.ToImplementations()
    js.Global().Get("console").Call("log", "Settings and implementations generated")

    // Register the simulation step function
    websocketPartitionIndex := 0
    if len(game.GetConfig().ServerPartitionNames) > 0 {
        for i, pName := range settings.Iterations {
            if pName.Name == game.GetConfig().ServerPartitionNames[0] {
                websocketPartitionIndex = i
                break
            }
        }
    }

    js.Global().Get("console").Call("log", "Calling RegisterStep")
    simio.RegisterStep(settings, implementations, websocketPartitionIndex, "", ":2112")
}
```

### Step 3: Add to Template Generator

Add your game to the template generator in `cmd/generate_game/main.go`:

```go
// Add to the switch statement
case "my_new_game":
    game = games.NewMyNewGame()
```

### Step 4: Generate Frontend Package

```bash
# Generate the complete frontend package
go run ./cmd/generate_game -game my_new_game -output ./my_new_game_package
```

This creates:
- `index.html` - Complete HTML page
- `styles.css` - Professional styling
- `game.js` - JavaScript with game logic
- `build.sh` - Build script for WebAssembly
- `wasm_exec.js` - Go WebAssembly runtime
- `worker.js` - WebAssembly worker
- `google-protobuf.js` - Protocol buffer support
- `partition_state_pb.js` - Partition state definitions

### Step 5: Build and Run

```bash
cd my_new_game_package
./build.sh                    # Builds WebAssembly module
python3 -m http.server 8000   # Start local server
# Open http://localhost:8000 in browser
```

## 🔧 Alternative: Integrated Development

For development and testing, you can also work directly in the `app/` directory:

### Option A: Use Existing App Structure

1. **Build WebAssembly**:
   ```bash
   GOOS=js GOARCH=wasm go build -o app/src/my_new_game/main.wasm ./cmd/my_new_game
   ```

2. **Create HTML Page**:
   ```bash
   cp app/minimal_example.html app/my_new_game.html
   # Edit the HTML to use your game
   ```

3. **Run**:
   ```bash
   cd app
   python3 -m http.server 8000
   # Open http://localhost:8000/my_new_game.html
   ```

### Option B: Use Template System (Recommended)

The template system is the recommended approach because it:
- ✅ **Generates complete packages** - Everything needed in one command
- ✅ **Handles dependencies** - Copies all required JS files
- ✅ **Creates proper build scripts** - Handles WebAssembly compilation
- ✅ **Professional styling** - Consistent, modern appearance
- ✅ **Easy deployment** - Self-contained packages

## 🎯 Key Benefits

### Template System Advantages:
- **Zero Manual Work** - Complete frontend generated automatically
- **Consistent Styling** - All games have professional appearance
- **Self-Contained** - Each package includes everything needed
- **Easy Deployment** - Single command creates deployable package
- **Extensible** - Generated code can be customized

### Development Workflow:
- **Fast Iteration** - Generate, build, test cycle
- **Type Safety** - All configuration from Go code
- **Version Control** - Only Go code needs to be committed
- **Reproducible** - Same Go code always generates same frontend

## 🚨 Common Issues and Solutions

### Issue 1: Build Script Fails
**Problem**: `./build.sh` fails with "directory not found"
**Solution**: The build script now correctly navigates to the project root

### Issue 2: Missing JavaScript Dependencies
**Problem**: Browser shows errors about missing JS files
**Solution**: Template system now copies all required dependencies

### Issue 3: Game Not Listed in CLI
**Problem**: `./generate_game -list` doesn't show your game
**Solution**: Add your game to the switch statement in `cmd/generate_game/main.go`

### Issue 4: WebAssembly Build Fails
**Problem**: Go build fails for WebAssembly
**Solution**: Ensure your command file has `//go:build js && wasm` directive

## 📁 File Structure

```
dexetera/
├── pkg/games/
│   ├── game.go                    # Core game framework
│   ├── template_generator.go      # Template system
│   ├── minimal_example.go         # Example games
│   ├── builder_example.go
│   └── my_new_game.go            # Your new game
├── cmd/
│   ├── generate_game/            # Template generator CLI
│   ├── minimal_example/          # WebAssembly commands
│   ├── builder_example/
│   └── my_new_game/              # Your WebAssembly command
├── app/                          # Existing app structure
│   ├── src/                      # WebAssembly modules
│   ├── *.html                    # Example pages
│   └── *.js                      # JavaScript dependencies
└── my_new_game_package/          # Generated package
    ├── index.html
    ├── styles.css
    ├── game.js
    ├── build.sh
    ├── src/main.wasm
    └── *.js                      # Dependencies
```

## 🎉 Summary

The Dexetera framework provides a complete, automated workflow:

1. **Write Go game code** using GameBuilder and VisualizationBuilder
2. **Create WebAssembly command** with proper build directives
3. **Generate complete frontend** with one command
4. **Build and deploy** with simple scripts

This eliminates manual frontend work while providing professional, consistent results for all games!
