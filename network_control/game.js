// Generated JavaScript for network_control
// Game configuration
const gameConfig = {
    name: "network_control",
    description: "Control junction phases to maximise traffic throughput.",
    partitionNames: {
        
    },
    serverPartitionNames: ["edge_states", "vehicle_rectangles", "flow_metrics", "junction_a_control", "junction_b_control", ],
    visualization: {
        canvasWidth: 800,
        canvasHeight: 600,
        backgroundColor: "#1a1a1a",
        updateIntervalMs: 0,
        renderers: [
            
            {
                type: "line",
                partitionName: "",
                properties: {
                    
                    "color": "#3d3d3d",
                    "width": 14,
                    "x1": 80,
                    "x2": 300,
                    "y1": 280,
                    "y2": 280,
                }
            },
            {
                type: "line",
                partitionName: "",
                properties: {
                    
                    "color": "#3d3d3d",
                    "width": 14,
                    "x1": 300,
                    "x2": 300,
                    "y1": 520,
                    "y2": 300,
                }
            },
            {
                type: "line",
                partitionName: "",
                properties: {
                    
                    "color": "#3d3d3d",
                    "width": 14,
                    "x1": 300,
                    "x2": 500,
                    "y1": 280,
                    "y2": 280,
                }
            },
            {
                type: "line",
                partitionName: "",
                properties: {
                    
                    "color": "#3d3d3d",
                    "width": 14,
                    "x1": 300,
                    "x2": 300,
                    "y1": 280,
                    "y2": 80,
                }
            },
            {
                type: "line",
                partitionName: "",
                properties: {
                    
                    "color": "#3d3d3d",
                    "width": 14,
                    "x1": 500,
                    "x2": 500,
                    "y1": 80,
                    "y2": 260,
                }
            },
            {
                type: "line",
                partitionName: "",
                properties: {
                    
                    "color": "#3d3d3d",
                    "width": 14,
                    "x1": 500,
                    "x2": 720,
                    "y1": 280,
                    "y2": 280,
                }
            },
            {
                type: "line",
                partitionName: "",
                properties: {
                    
                    "color": "#3d3d3d",
                    "width": 14,
                    "x1": 500,
                    "x2": 500,
                    "y1": 280,
                    "y2": 520,
                }
            },
            {
                type: "rectangleSet",
                partitionName: "vehicle_rectangles",
                properties: {
                    
                    "defaultHeight": 10,
                    "defaultWidth": 18,
                    "fillColor": "#f4a261",
                    "strokeColor": "#e76f51",
                    "strokeWidth": 1,
                }
            },
            {
                type: "text",
                partitionName: "junction_a_control",
                properties: {
                    
                    "color": "#ffffff",
                    "fontFamily": "Arial",
                    "fontSize": 14,
                    "text": "Junction A Phase: {value}",
                    "textAlign": "left",
                    "x": 220,
                    "y": 90,
                }
            },
            {
                type: "text",
                partitionName: "junction_b_control",
                properties: {
                    
                    "color": "#ffffff",
                    "fontFamily": "Arial",
                    "fontSize": 14,
                    "text": "Junction B Phase: {value}",
                    "textAlign": "left",
                    "x": 520,
                    "y": 90,
                }
            },
            {
                type: "text",
                partitionName: "flow_metrics",
                properties: {
                    
                    "color": "#ffffff",
                    "fontFamily": "Arial",
                    "fontSize": 16,
                    "text": "Vehicles exited: {value}",
                    "textAlign": "center",
                    "x": 600,
                    "y": 520,
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
        stopAtSimTime: 10000.00,
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
        if (!state && renderer.partitionName !== '') return;
        
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
            case 'rectangleSet':
                this.renderRectangleSet(renderer, state);
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
            case 'progressBar':
                this.renderProgressBar(renderer, state);
                break;
            case 'image':
                this.renderImage(renderer, state);
                break;
            case 'playerSet':
            case 'pointSet':
                this.renderPointSet(renderer, state);
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
        
        // For static rectangles, always render
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
    
    renderRectangleSet(renderer, state) {
        const defaultWidth = renderer.properties.defaultWidth || 12;
        const defaultHeight = renderer.properties.defaultHeight || 8;
        const fill = renderer.properties.fillColor || renderer.properties.color || '#ffffff';
        const stroke = renderer.properties.strokeColor;
        const strokeWidth = renderer.properties.strokeWidth || 1;

        for (let i = 0; i + 3 < state.length; i += 4) {
            const x = state[i];
            const y = state[i + 1];
            const width = state[i + 2] || defaultWidth;
            const height = state[i + 3] || defaultHeight;

            if (!width || !height) {
                continue;
            }

            const drawWidth = Math.abs(width);
            const drawHeight = Math.abs(height);
            const left = x - drawWidth / 2;
            const top = y - drawHeight / 2;

            this.ctx.fillStyle = fill;
            this.ctx.fillRect(left, top, drawWidth, drawHeight);

            if (stroke) {
                this.ctx.strokeStyle = stroke;
                this.ctx.lineWidth = strokeWidth;
                this.ctx.strokeRect(left, top, drawWidth, drawHeight);
            }
        }
    }
    
    renderLine(renderer, state) {
        const x1 = renderer.properties.x1 || 0;
        const y1 = renderer.properties.y1 || 0;
        const x2 = renderer.properties.x2 || 50;
        const y2 = renderer.properties.y2 || 50;
        
        // For static lines, always render
        this.ctx.beginPath();
        this.ctx.moveTo(x1, y1);
        this.ctx.lineTo(x2, y2);
        this.ctx.strokeStyle = renderer.properties.color || '#ffffff';
        this.ctx.lineWidth = renderer.properties.width || 1;
        this.ctx.stroke();
    }
    
    renderBarChart(renderer, state) {
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        const width = renderer.properties.width || 50;
        const height = renderer.properties.height || 50;
        const maxValue = renderer.properties.maxValue || 100;
        const value = state[0] || 0;
        const normalizedValue = Math.min(value / maxValue, 1.0);
        
        // Draw background
        this.ctx.fillStyle = renderer.properties.color || 'rgba(255,255,255,0.3)';
        this.ctx.fillRect(x, y, width, height);
        
        // Draw bar
        this.ctx.fillStyle = renderer.properties.color || '#4CAF50';
        this.ctx.fillRect(x, y + height * (1 - normalizedValue), width, height * normalizedValue);
        
        // Draw label if requested
        if (renderer.properties.showLabels) {
            this.ctx.fillStyle = '#ffffff';
            this.ctx.font = '12px Arial';
            this.ctx.textAlign = 'center';
            this.ctx.fillText(Math.floor(value), x + width / 2, y + height / 2);
        }
    }
    
    renderLineChart(renderer, state) {
        const history = this.history[renderer.partitionName];
        if (!history || history.length < 2) return;
        
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        const width = renderer.properties.width || 50;
        const height = renderer.properties.height || 50;
        const maxValue = renderer.properties.maxValue || 100;
        
        // Find min/max for scaling
        let minVal = Infinity, maxVal = -Infinity;
        history.forEach(point => {
            minVal = Math.min(minVal, point.value);
            maxVal = Math.max(maxVal, point.value);
        });
        const range = Math.max(maxVal - minVal, 0.1);
        
        this.ctx.strokeStyle = renderer.properties.color || '#4CAF50';
        this.ctx.lineWidth = renderer.properties.lineWidth || 2;
        this.ctx.beginPath();
        
        history.forEach((point, i) => {
            const px = x + (i / (history.length - 1)) * width;
            const py = y + height - ((point.value - minVal) / range) * height;
            
            if (i === 0) {
                this.ctx.moveTo(px, py);
            } else {
                this.ctx.lineTo(px, py);
            }
        });
        
        this.ctx.stroke();
    }
    
    renderProgressBar(renderer, state) {
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        const width = renderer.properties.width || 100;
        const height = renderer.properties.height || 20;
        const maxValue = renderer.properties.maxValue || 100;
        const value = Math.max(0, Math.min(state[0] || 0, maxValue));
        const normalizedValue = value / maxValue;
        
        // Draw background
        this.ctx.fillStyle = renderer.properties.backgroundColor || 'rgba(255,255,255,0.3)';
        this.ctx.fillRect(x, y, width, height);
        
        // Draw progress
        this.ctx.fillStyle = renderer.properties.foregroundColor || '#4CAF50';
        this.ctx.fillRect(x, y, width * normalizedValue, height);
        
        // Draw border if specified
        if (renderer.properties.borderColor) {
            this.ctx.strokeStyle = renderer.properties.borderColor;
            this.ctx.lineWidth = renderer.properties.borderWidth || 1;
            this.ctx.strokeRect(x, y, width, height);
        }
        
        // Draw label if requested
        if (renderer.properties.showLabel) {
            this.ctx.fillStyle = '#ffffff';
            this.ctx.font = '12px Arial';
            this.ctx.textAlign = 'center';
            this.ctx.fillText(Math.floor(value) + '%', x + width / 2, y + height / 2 + 4);
        }
    }
    
    renderImage(renderer, state) {
        const imagePath = renderer.properties.imagePath;
        if (!imagePath) return;
        
        // For now, we'll implement basic rendering
        // In a full implementation, you'd load and cache images
        const x = renderer.properties.x || 0;
        const y = renderer.properties.y || 0;
        
        // Draw placeholder rectangle for now
        this.ctx.fillStyle = 'rgba(255,255,255,0.5)';
        this.ctx.fillRect(x, y, 
            renderer.properties.width || 32, 
            renderer.properties.height || 32);
    }

    renderPointSet(renderer, state) {
        const radius = renderer.properties.radius || 8;
        const fill = renderer.properties.fillColor || renderer.properties.color || '#ffffff';
        const stroke = renderer.properties.strokeColor;
        const strokeWidth = renderer.properties.strokeWidth || 1;

        for (let i = 0; i < state.length; i += 2) {
            const x = state[i];
            const y = state[i + 1];
            if (typeof x !== 'number' || typeof y !== 'number') {
                continue;
            }

            this.ctx.beginPath();
            this.ctx.arc(x, y, radius, 0, 2 * Math.PI);
            this.ctx.fillStyle = fill;
            this.ctx.fill();

            if (stroke) {
                this.ctx.strokeStyle = stroke;
                this.ctx.lineWidth = strokeWidth;
                this.ctx.stroke();
            }
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
