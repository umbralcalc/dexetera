package games

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// GameTemplateGenerator generates complete frontend files from game configurations
type GameTemplateGenerator struct {
	game Game
}

// NewGameTemplateGenerator creates a new template generator for a game
func NewGameTemplateGenerator(game Game) *GameTemplateGenerator {
	return &GameTemplateGenerator{game: game}
}

// GenerateFrontend generates all frontend files for the game
func (gtg *GameTemplateGenerator) GenerateFrontend(outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate HTML file
	if err := gtg.generateHTML(outputDir); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Generate CSS file
	if err := gtg.generateCSS(outputDir); err != nil {
		return fmt.Errorf("failed to generate CSS: %w", err)
	}

	// Generate JavaScript file
	if err := gtg.generateJavaScript(outputDir); err != nil {
		return fmt.Errorf("failed to generate JavaScript: %w", err)
	}

	// Generate WebAssembly build script
	if err := gtg.generateBuildScript(outputDir); err != nil {
		return fmt.Errorf("failed to generate build script: %w", err)
	}

	// Copy necessary JavaScript dependencies
	if err := gtg.copyDependencies(outputDir); err != nil {
		return fmt.Errorf("failed to copy dependencies: %w", err)
	}

	return nil
}

// generateHTML generates the main HTML file
func (gtg *GameTemplateGenerator) generateHTML(outputDir string) error {
	config := gtg.game.GetConfig()
	renderer := gtg.game.GetRenderer()

	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>dexetera... {{.GameName}}</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
    <div class="container">
        <h1>🎮 {{.GameName}}</h1>
        <p class="description">{{.Description}}</p>
        
        <div class="game-area">
            <div class="game-container">
                <canvas id="gameCanvas" width="{{.CanvasWidth}}" height="{{.CanvasHeight}}"></canvas>
            </div>
            
            <div class="status" id="status">Waiting for Python server to connect...</div>
        </div>
        
        <div class="info">
            <h3>🔧 Game Configuration</h3>
            <div class="config-info">
                <p><strong>Partitions:</strong> {{range $key, $value := .PartitionNames}}{{$key}} ({{$value}}){{end}}</p>
                <p><strong>Server Partitions:</strong> {{range .ServerPartitionNames}}{{.}} {{end}}</p>
                <p><strong>Update Interval:</strong> {{.UpdateIntervalMs}}ms</p>
            </div>
        </div>
    </div>
    
    <div class="overlay" id="overlay">
        <div class="overlay-content">
            <h2>Loading Game...</h2>
            <p>Please wait while the WebAssembly module loads.</p>
        </div>
    </div>

    <script src="wasm_exec.js"></script>
    <script src="google-protobuf.js"></script>
    <script src="partition_state_pb.js"></script>
    <script src="game.js"></script>
</body>
</html>`

	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	visConfig := renderer.GetVisualizationConfig()
	data := struct {
		GameName             string
		Description          string
		CanvasWidth          int
		CanvasHeight         int
		PartitionNames       map[string]string
		ServerPartitionNames []string
		UpdateIntervalMs     int
	}{
		GameName:             config.Name,
		Description:          config.Description,
		CanvasWidth:          visConfig.CanvasWidth,
		CanvasHeight:         visConfig.CanvasHeight,
		PartitionNames:       config.PartitionNames,
		ServerPartitionNames: config.ServerPartitionNames,
		UpdateIntervalMs:     visConfig.UpdateIntervalMs,
	}

	file, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateCSS generates the CSS file
func (gtg *GameTemplateGenerator) generateCSS(outputDir string) error {
	renderer := gtg.game.GetRenderer()
	visConfig := renderer.GetVisualizationConfig()

	cssTemplate := `/* Generated CSS for {{.GameName}} */
body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    margin: 0;
    padding: 20px;
    background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
    color: #ffffff;
    min-height: 100vh;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

h1 {
    text-align: center;
    margin-bottom: 10px;
    font-size: 2.5em;
    text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
}

.description {
    text-align: center;
    font-size: 1.2em;
    margin-bottom: 30px;
    opacity: 0.9;
}

.game-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
}

.game-container {
    background-color: {{.BackgroundColor}};
    border: 2px solid #444;
    border-radius: 8px;
    padding: 10px;
    box-shadow: 0 4px 8px rgba(0,0,0,0.3);
}

.game-container canvas {
    display: block;
    margin: 0 auto;
    border-radius: 4px;
}

.status {
    font-size: 1.1em;
    padding: 10px 20px;
    background: rgba(255,255,255,0.1);
    border-radius: 20px;
    text-align: center;
    min-width: 300px;
}

.info {
    margin-top: 30px;
    padding: 20px;
    background: rgba(255,255,255,0.1);
    border-radius: 8px;
    backdrop-filter: blur(10px);
}

.info h3 {
    margin-top: 0;
    color: #4CAF50;
}

.config-info p {
    margin: 8px 0;
    font-family: 'Courier New', monospace;
    background: rgba(0,0,0,0.3);
    padding: 8px 12px;
    border-radius: 4px;
}

.overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.8);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
}

.overlay-content {
    text-align: center;
    color: white;
}

.overlay-content h2 {
    margin-bottom: 20px;
    font-size: 2em;
}

.overlay-content p {
    font-size: 1.2em;
    opacity: 0.9;
}

/* Responsive design */
@media (max-width: 768px) {
    .container {
        padding: 10px;
    }
    
    h1 {
        font-size: 2em;
    }
    
    .game-container canvas {
        max-width: 100%;
        height: auto;
    }
}`

	tmpl, err := template.New("css").Parse(cssTemplate)
	if err != nil {
		return err
	}

	data := struct {
		GameName        string
		BackgroundColor string
	}{
		GameName:        gtg.game.GetName(),
		BackgroundColor: visConfig.BackgroundColor,
	}

	file, err := os.Create(filepath.Join(outputDir, "styles.css"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateJavaScript generates the JavaScript file
func (gtg *GameTemplateGenerator) generateJavaScript(outputDir string) error {
	config := gtg.game.GetConfig()
	renderer := gtg.game.GetRenderer()

	jsTemplate := `// Generated JavaScript for {{.GameName}}
// Game configuration
const gameConfig = {
    name: "{{.GameName}}",
    description: "{{.Description}}",
    partitionNames: {
        {{range $key, $value := .PartitionNames}}
        "{{$key}}": "{{$value}}",{{end}}
    },
    serverPartitionNames: [{{range .ServerPartitionNames}}"{{.}}", {{end}}],
    visualization: {
        canvasWidth: {{.CanvasWidth}},
        canvasHeight: {{.CanvasHeight}},
        backgroundColor: "{{.BackgroundColor}}",
        updateIntervalMs: {{.UpdateIntervalMs}},
        renderers: [
            {{range .Renderers}}
            {
                type: "{{.Type}}",
                partitionName: "{{.PartitionName}}",
                properties: {
                    {{range $key, $value := .Properties}}
                    "{{$key}}": {{if eq $key "text"}}"{{$value}}"{{else if eq $key "color"}}"{{$value}}"{{else if eq $key "fontSize"}}{{$value}}{{else if eq $key "x"}}{{$value}}{{else if eq $key "y"}}{{$value}}{{else if eq $key "width"}}{{$value}}{{else if eq $key "height"}}{{$value}}{{else if eq $key "radius"}}{{$value}}{{else if eq $key "fontFamily"}}"{{$value}}"{{else if eq $key "textAlign"}}"{{$value}}"{{else if eq $key "fillColor"}}"{{$value}}"{{else if eq $key "strokeColor"}}"{{$value}}"{{else if eq $key "strokeWidth"}}{{$value}}{{else if eq $key "maxValue"}}{{$value}}{{else if eq $key "showLabels"}}{{$value}}{{else if eq $key "labelFormat"}}"{{$value}}"{{else if eq $key "lineWidth"}}{{$value}}{{else}}{{$value}}{{end}},{{end}}
                }
            },{{end}}
        ]
    }
};

// Binding resolver: supports { bind: { type: 'state'|'param', partition: 'name', index: 0, key?: 'paramKey' } }
const stateCache = {}; // partitionName -> values array
const paramCache = {}; // partitionName -> { key -> [] }

function resolveBinding(val) {
    if (!val || typeof val !== 'object') return val;
    const b = val.bind;
    if (!b || typeof b !== 'object') return val;
    // state binding
    if (b.type === 'state') {
        const part = b.partition;
        const idx = b.index || 0;
        const arr = stateCache[part];
        if (arr && typeof arr[idx] === 'number') return arr[idx];
        return 0;
    }
    // param binding: lookup from paramCache if available
    if (b.type === 'param') {
        const part = b.partition;
        const key = b.key || 'param_values';
        const idx = b.index || 0;
        const obj = paramCache[part];
        if (obj && Array.isArray(obj[key]) && typeof obj[key][idx] === 'number') return obj[key][idx];
        // Fallback: if iteration mirrors param_values into state, reuse state
        const arr = stateCache[part];
        if (arr && typeof arr[idx] === 'number') return arr[idx];
        return 0;
    }
    return val;
}

function prop(renderer, name, fallback) {
    const v = renderer.properties ? renderer.properties[name] : undefined;
    const r = resolveBinding(v);
    return r === undefined ? fallback : r;
}

// Generic renderer with binding-aware properties
class GenericRenderer {
    constructor(canvas, config) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.config = config;
        this.state = {};
    }
    
    update(partitionState) {
        // cache raw array for binding resolver
        stateCache[partitionState.partitionName] = partitionState.state.values;
        this.state[partitionState.partitionName] = partitionState.state.values;
    }
    
    render() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        // Render each configured renderer
        this.config.renderers.forEach(renderer => {
            this.renderElement(renderer);
        });
    }
    
    renderElement(renderer) {
        const state = this.state[renderer.partitionName];
        if (!state) return;
        
        switch (renderer.type) {
            case 'text':
                this.renderText(renderer, state);
                break;
            case 'circle':
                this.renderCircle(renderer, state);
                break;
            case 'rectangle':
                this.renderRectangle(renderer, state);
                break;
        }
    }
    
    renderText(renderer, state) {
        this.ctx.fillStyle = prop(renderer, 'color', '#ffffff');
        const fontSize = prop(renderer, 'fontSize', 16);
        const fontFamily = prop(renderer, 'fontFamily', 'Arial');
        this.ctx.font = String(fontSize) + 'px ' + String(fontFamily);
        this.ctx.textAlign = prop(renderer, 'textAlign', 'center');
        
        let text = prop(renderer, 'text', '{value}');
        if (typeof text === 'string') {
            const v = Math.floor(state[0] || 0);
            text = text.replace('{value}', v);
        }
        
        this.ctx.fillText(
            text,
            prop(renderer, 'x', this.canvas.width / 2),
            prop(renderer, 'y', this.canvas.height / 2),
        );
    }
    
    renderCircle(renderer, state) {
        const x = prop(renderer, 'x', this.canvas.width / 2);
        const y = prop(renderer, 'y', this.canvas.height / 2);
        const radius = prop(renderer, 'radius', 10);
        
        this.ctx.beginPath();
        this.ctx.arc(x, y, radius, 0, 2 * Math.PI);
        
        const fillColor = prop(renderer, 'fillColor', null);
        const strokeColor = prop(renderer, 'strokeColor', null);
        const strokeWidth = prop(renderer, 'strokeWidth', 1);
        
        if (fillColor) {
            this.ctx.fillStyle = fillColor;
            this.ctx.fill();
        }
        if (strokeColor) {
            this.ctx.strokeStyle = strokeColor;
            this.ctx.lineWidth = strokeWidth;
            this.ctx.stroke();
        }
        if (!fillColor && !strokeColor) {
            this.ctx.fillStyle = prop(renderer, 'color', '#ffffff');
            this.ctx.fill();
        }
    }
    
    renderRectangle(renderer, state) {
        const x = prop(renderer, 'x', 0);
        const y = prop(renderer, 'y', 0);
        const width = prop(renderer, 'width', 50);
        const height = prop(renderer, 'height', 50);
        
        const fillColor = prop(renderer, 'fillColor', null);
        const strokeColor = prop(renderer, 'strokeColor', null);
        const strokeWidth = prop(renderer, 'strokeWidth', 1);
        
        if (fillColor) {
            this.ctx.fillStyle = fillColor;
            this.ctx.fillRect(x, y, width, height);
        }
        if (strokeColor) {
            this.ctx.strokeStyle = strokeColor;
            this.ctx.lineWidth = strokeWidth;
            this.ctx.strokeRect(x, y, width, height);
        }
        if (!fillColor && !strokeColor) {
            this.ctx.fillStyle = prop(renderer, 'color', '#ffffff');
            this.ctx.fillRect(x, y, width, height);
        }
    }
}

// Global variables
let gameRenderer = null;
let worker = null;

// Initialize the game
function initializeGame() {
    const canvas = document.getElementById('gameCanvas');
    gameRenderer = new GenericRenderer(canvas, gameConfig.visualization);
    
    // Set up worker for WebAssembly
    worker = new Worker('worker.js');
    
    worker.onmessage = function(e) {
        const { type, data } = e.data;
        
        if (type === 'partitionState') {
            updateVisualization(data);
            // Update status to show we're receiving data
            const partitionName = data.partitionName;
            const value = Math.floor(data.state.values[0] || 0);
            document.getElementById('status').textContent = ` + "`${partitionName}: ${value} (Python-controlled)`" + `;
        } else if (type === 'error') {
            console.error('Worker error:', data);
            document.getElementById('status').textContent = 'Error: ' + data;
        }
    };
    
    worker.onerror = function(error) {
        console.error('Worker error:', error);
        document.getElementById('status').textContent = 'Worker error: ' + error.message;
    };
    
    // Hide loading overlay
    document.getElementById('overlay').style.display = 'none';
    document.getElementById('status').textContent = 'Ready - waiting for Python server...';
    
    // Start the WebAssembly module
    worker.postMessage({ 
        action: 'start', 
        wasmBinary: 'src/main.wasm',
        serverPartitionNames: gameConfig.serverPartitionNames,
        stopAtSimTime: 30.05,
        debugMode: false
    });
}

// Update visualization with new data
function updateVisualization(partitionState) {
    if (gameRenderer) {
        gameRenderer.update(partitionState);
        gameRenderer.render();
    }
}

// Initialize when page loads
window.addEventListener('load', initializeGame);`

	tmpl, err := template.New("js").Parse(jsTemplate)
	if err != nil {
		return err
	}

	visConfig := renderer.GetVisualizationConfig()
	data := struct {
		GameName             string
		Description          string
		CanvasWidth          int
		CanvasHeight         int
		BackgroundColor      string
		UpdateIntervalMs     int
		PartitionNames       map[string]string
		ServerPartitionNames []string
		Renderers            []RendererConfig
	}{
		GameName:             config.Name,
		Description:          config.Description,
		CanvasWidth:          visConfig.CanvasWidth,
		CanvasHeight:         visConfig.CanvasHeight,
		BackgroundColor:      visConfig.BackgroundColor,
		UpdateIntervalMs:     visConfig.UpdateIntervalMs,
		PartitionNames:       config.PartitionNames,
		ServerPartitionNames: config.ServerPartitionNames,
		Renderers:            visConfig.Renderers,
	}

	file, err := os.Create(filepath.Join(outputDir, "game.js"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the template first
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	// Write the renderer JavaScript directly to avoid HTML encoding
	rendererJS := renderer.GetJavaScriptCode()
	_, err = file.WriteString(rendererJS)
	return err
}

// generateBuildScript generates a build script for the WebAssembly module
func (gtg *GameTemplateGenerator) generateBuildScript(outputDir string) error {
	config := gtg.game.GetConfig()

	buildScript := `#!/bin/bash
# Generated build script for {{.GameName}}

echo "Building {{.GameName}} WebAssembly module..."

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Create src directory if it doesn't exist
mkdir -p src

# Build the WebAssembly module from the project root
cd "$PROJECT_ROOT"
GOOS=js GOARCH=wasm go build -o "$SCRIPT_DIR/src/main.wasm" ./cmd/{{.GameName}}

if [ $? -eq 0 ]; then
    echo "✅ WebAssembly module built successfully!"
    echo "📁 Output: $SCRIPT_DIR/src/main.wasm"
else
    echo "❌ Build failed!"
    exit 1
fi

echo "🎮 {{.GameName}} is ready to run!"
echo "📝 Start your Python websocket server and open index.html in a browser"
`

	tmpl, err := template.New("build").Parse(buildScript)
	if err != nil {
		return err
	}

	data := struct {
		GameName string
	}{
		GameName: config.Name,
	}

	file, err := os.Create(filepath.Join(outputDir, "build.sh"))
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	// Make the script executable
	return os.Chmod(filepath.Join(outputDir, "build.sh"), 0755)
}

// GenerateGamePackage generates a complete game package ready for deployment
func GenerateGamePackage(game Game, outputDir string) error {
	generator := NewGameTemplateGenerator(game)
	return generator.GenerateFrontend(outputDir)
}

// copyDependencies copies necessary JavaScript dependencies
func (gtg *GameTemplateGenerator) copyDependencies(outputDir string) error {
	// List of files to copy from the app directory
	dependencies := []string{
		"wasm_exec.js",
		"google-protobuf.js",
		"partition_state_pb.js",
		"worker.js",
	}

	// Find the project root by looking for go.mod
	projectRoot, err := findProjectRoot()
	if err != nil {
		fmt.Printf("Warning: Could not find project root: %v\n", err)
		return nil // Don't fail the entire generation
	}

	for _, dep := range dependencies {
		srcPath := filepath.Join(projectRoot, "app", "src", dep)
		dstPath := filepath.Join(outputDir, dep)

		if err := copyFile(srcPath, dstPath); err != nil {
			// If file doesn't exist, create a placeholder or skip
			fmt.Printf("Warning: Could not copy %s: %v\n", dep, err)
			continue
		}
	}

	return nil
}

// findProjectRoot finds the project root by looking for go.mod
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find project root")
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
