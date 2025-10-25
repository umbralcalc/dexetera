package games

import (
	"math/rand"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/stat/distuv"
)

// ParticleSystemGame is a simple example game that demonstrates the new framework.
// It shows a collection of particles that move around randomly on a 2D canvas.
type ParticleSystemGame struct {
	config *GameConfig
}

// NewParticleSystemGame creates a new particle system game
func NewParticleSystemGame() *ParticleSystemGame {
	return &ParticleSystemGame{
		config: &GameConfig{
			Name:        "particle_system",
			Description: "A simple particle system with random movement",
			PartitionNames: map[string]string{
				"particles": "particle_positions",
			},
			ServerPartitionNames: []string{"particle_positions"},
			VisualizationConfig: &VisualizationConfig{
				CanvasWidth:      800,
				CanvasHeight:     600,
				BackgroundColor:  "#1a1a1a",
				UpdateIntervalMs: 50,
				Renderers: []RendererConfig{
					{
						Type:          "circle",
						PartitionName: "particle_positions",
						Properties: map[string]interface{}{
							"radius":      5.0,
							"color":       "#00ff88",
							"strokeColor": "#ffffff",
							"strokeWidth": 1.0,
						},
					},
				},
			},
			Parameters: map[string]interface{}{
				"num_particles": 50,
				"canvas_width":  800.0,
				"canvas_height": 600.0,
				"max_speed":     2.0,
			},
		},
	}
}

// GetName returns the game name
func (p *ParticleSystemGame) GetName() string {
	return p.config.Name
}

// GetDescription returns the game description
func (p *ParticleSystemGame) GetDescription() string {
	return p.config.Description
}

// GetConfig returns the game configuration
func (p *ParticleSystemGame) GetConfig() *GameConfig {
	return p.config
}

// GetSettings returns the stochadex settings for this game
func (p *ParticleSystemGame) GetSettings() *simulator.Settings {
	numParticles := p.config.Parameters["num_particles"].(int)
	canvasWidth := p.config.Parameters["canvas_width"].(float64)
	canvasHeight := p.config.Parameters["canvas_height"].(float64)
	maxSpeed := p.config.Parameters["max_speed"].(float64)

	// Initialize particle positions randomly
	initialPositions := make([]float64, numParticles*2)
	for i := 0; i < numParticles; i++ {
		initialPositions[i*2] = rand.Float64() * canvasWidth    // x
		initialPositions[i*2+1] = rand.Float64() * canvasHeight // y
	}

	settings := &simulator.Settings{
		Iterations: []simulator.IterationSettings{
			{
				Name: "particle_positions",
				Params: simulator.NewParams(map[string][]float64{
					"canvas_width":  {canvasWidth},
					"canvas_height": {canvasHeight},
					"max_speed":     {maxSpeed},
				}),
				InitStateValues:   initialPositions,
				Seed:              42,
				StateWidth:        numParticles * 2, // x, y for each particle
				StateHistoryDepth: 1,
			},
		},
		InitTimeValue:         0.0,
		TimestepsHistoryDepth: 1,
	}

	return settings
}

// GetImplementations returns the stochadex implementations
func (p *ParticleSystemGame) GetImplementations() *simulator.Implementations {
	return &simulator.Implementations{
		Iterations: []simulator.Iteration{
			&ParticleSystemIteration{},
		},
		OutputCondition: &simulator.EveryStepOutputCondition{},
		OutputFunction:  &simulator.StdoutOutputFunction{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{
			MaxTimeElapsed: 300.0, // 5 minutes
		},
		TimestepFunction: &simulator.ConstantTimestepFunction{
			Stepsize: 0.016, // ~60 FPS
		},
	}
}

// GetRenderer returns the visualization renderer
func (p *ParticleSystemGame) GetRenderer() GameRenderer {
	return &ParticleSystemRenderer{config: p.config.VisualizationConfig}
}

// ParticleSystemIteration implements the particle movement logic
type ParticleSystemIteration struct {
	uniformDist *distuv.Uniform
}

func (p *ParticleSystemIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	p.uniformDist = &distuv.Uniform{
		Min: -1.0,
		Max: 1.0,
		Src: rand.New(rand.NewSource(int64(settings.Iterations[partitionIndex].Seed))),
	}
}

func (p *ParticleSystemIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	outputState := stateHistories[partitionIndex].Values.RawRowView(0)
	canvasWidth := params.Get("canvas_width")[0]
	canvasHeight := params.Get("canvas_height")[0]
	maxSpeed := params.Get("max_speed")[0]

	// Move each particle randomly
	for i := 0; i < len(outputState); i += 2 {
		// Add random movement
		dx := p.uniformDist.Rand() * maxSpeed
		dy := p.uniformDist.Rand() * maxSpeed

		// Update position
		newX := outputState[i] + dx
		newY := outputState[i+1] + dy

		// Wrap around canvas edges
		if newX < 0 {
			newX = canvasWidth
		} else if newX > canvasWidth {
			newX = 0
		}

		if newY < 0 {
			newY = canvasHeight
		} else if newY > canvasHeight {
			newY = 0
		}

		outputState[i] = newX
		outputState[i+1] = newY
	}

	return outputState
}

// ParticleSystemRenderer handles the visualization of the particle system
type ParticleSystemRenderer struct {
	config *VisualizationConfig
}

func (r *ParticleSystemRenderer) GetVisualizationConfig() *VisualizationConfig {
	return r.config
}

func (r *ParticleSystemRenderer) GetJavaScriptCode() string {
	return `
// Particle system visualization JavaScript
class ParticleSystemRenderer {
    constructor(canvas, config) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.config = config;
        this.particles = [];
    }
    
    update(partitionState) {
        if (partitionState.partitionName === 'particle_positions') {
            this.particles = [];
            const values = partitionState.state.values;
            for (let i = 0; i < values.length; i += 2) {
                this.particles.push({
                    x: values[i],
                    y: values[i + 1]
                });
            }
        }
    }
    
    render() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        this.particles.forEach(particle => {
            this.ctx.beginPath();
            this.ctx.arc(particle.x, particle.y, 5, 0, 2 * Math.PI);
            this.ctx.fillStyle = '#00ff88';
            this.ctx.fill();
            this.ctx.strokeStyle = '#ffffff';
            this.ctx.lineWidth = 1;
            this.ctx.stroke();
        });
    }
}

// Global renderer instance
let particleRenderer = null;

function initializeRenderer(canvas, config) {
    particleRenderer = new ParticleSystemRenderer(canvas, config);
}

function updateVisualization(partitionState) {
    if (particleRenderer) {
        particleRenderer.update(partitionState);
        particleRenderer.render();
    }
}
`
}

func (r *ParticleSystemRenderer) GetCSSCode() string {
	return `
.particle-system {
    background-color: #1a1a1a;
    border: 2px solid #333;
    border-radius: 8px;
}

.particle-system canvas {
    display: block;
    margin: 0 auto;
}
`
}
