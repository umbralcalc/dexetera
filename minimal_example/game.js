// Generated JavaScript for minimal_example
// Game configuration
const gameConfig = {
    name: "minimal_example",
    description: "The simplest possible game - just a counter",
    partitionNames: {
        
        "counter": "counter_state",
    },
    serverPartitionNames: ["counter_state", ],
    visualization: {
        canvasWidth: 400,
        canvasHeight: 200,
        backgroundColor: "#2a2a2a",
        updateIntervalMs: 100,
        renderers: [
            
            {
                type: "text",
                partitionName: "counter_state",
                properties: {
                    
                    "color": "#ffffff",
                    "fontSize": 24,
                    "text": "Count: {value}",
                    "x": 200,
                    "y": 100,
                }
            },
        ]
    }
};

// Global variables
let worker = null;

// Initialize the game
function initializeGame() {
    const canvas = document.getElementById('gameCanvas');
    // Defer renderer creation to renderer-provided code
    if (typeof initializeRenderer === 'function') {
        initializeRenderer(canvas, gameConfig.visualization);
    }
    
    // Set up worker for WebAssembly
    worker = new Worker('worker.js');
    
    worker.onmessage = function(e) {
        const { type, data } = e.data;
        
        if (type === 'partitionState') {
            if (typeof updateVisualization === 'function') {
                updateVisualization(data);
            }
            // Update status to show we're receiving data
            const partitionName = data.partitionName;
            const value = Math.floor(data.state.values[0] || 0);
            document.getElementById('status').textContent = partitionName + ': ' + value + ' (Python-controlled)';
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

// Initialize when page loads
window.addEventListener('load', initializeGame);
// Enhanced Generic renderer JavaScript with support for all renderer types
class GenericRenderer {
    constructor(canvas, config) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.config = config;
        this.state = {};
        this.history = {}; // For charts
    }
    
    update(partitionState) {
        this.state[partitionState.partitionName] = partitionState.state.values;
        
        // Store history for charts
        if (!this.history[partitionState.partitionName]) {
            this.history[partitionState.partitionName] = [];
        }
        this.history[partitionState.partitionName].push({
            value: partitionState.state.values[0] || 0,
            time: partitionState.cumulativeTimesteps || 0
        });
        
        // Keep only last 100 points for performance
        if (this.history[partitionState.partitionName].length > 100) {
            this.history[partitionState.partitionName].shift();
        }
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
            case 'line':
                this.renderLine(renderer, state);
                break;
            case 'barChart':
                this.renderBarChart(renderer, state);
                break;
            case 'lineChart':
                this.renderLineChart(renderer, state);
                break;
        }
    }
    
    renderText(renderer, state) {
        this.ctx.fillStyle = '#ffffff';
        this.ctx.font = '16px Arial';
        this.ctx.textAlign = 'center';
        
        let text = renderer.properties.text || '{value}';
        text = text.replace('{value}', Math.floor(state[0] || 0));
        
        this.ctx.fillText(text, 
                          renderer.properties.x || this.canvas.width / 2,
                          renderer.properties.y || this.canvas.height / 2);
    }
    
    renderCircle(renderer, state) {
        const x = renderer.properties.x || this.canvas.width / 2;
        const y = renderer.properties.y || this.canvas.height / 2;
        const radius = renderer.properties.radius || 10;
        
        this.ctx.beginPath();
        this.ctx.arc(x, y, radius, 0, 2 * Math.PI);
        
        if (renderer.properties.fillColor) {
            this.ctx.fillStyle = renderer.properties.fillColor;
            this.ctx.fill();
        }
        
        if (renderer.properties.strokeColor) {
            this.ctx.strokeStyle = renderer.properties.strokeColor;
            this.ctx.lineWidth = renderer.properties.strokeWidth || 1;
            this.ctx.stroke();
        }
        
        if (!renderer.properties.fillColor && !renderer.properties.strokeColor) {
            this.ctx.fillStyle = renderer.properties.color || '#ffffff';
            this.ctx.fill();
        }
    }
    
    renderRectangle(renderer, state) {
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        const width = renderer.properties.width || 50;
        const height = renderer.properties.height || 50;
        
        if (renderer.properties.fillColor) {
            this.ctx.fillStyle = renderer.properties.fillColor;
            this.ctx.fillRect(x, y, width, height);
        }
        
        if (renderer.properties.strokeColor) {
            this.ctx.strokeStyle = renderer.properties.strokeColor;
            this.ctx.lineWidth = renderer.properties.strokeWidth || 1;
            this.ctx.strokeRect(x, y, width, height);
        }
        
        if (!renderer.properties.fillColor && !renderer.properties.strokeColor) {
            this.ctx.fillStyle = renderer.properties.color || '#ffffff';
            this.ctx.fillRect(x, y, width, height);
        }
    }
}

// Global renderer instance
let genericRenderer = null;

function initializeRenderer(canvas, config) {
    genericRenderer = new GenericRenderer(canvas, config);
}

function updateVisualization(partitionState) {
    if (genericRenderer) {
        genericRenderer.update(partitionState);
        genericRenderer.render();
    }
}
