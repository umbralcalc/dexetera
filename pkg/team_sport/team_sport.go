package team_sport

import (
	"github.com/umbralcalc/dexetera/pkg/game"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// TeamSportGame simulates a team sport match where players have stamina,
// and the user must make substitutions at the right time to win.
type TeamSportGame struct {
	config *game.GameConfig
}

// NewTeamSportGame creates a new team sport game
func NewTeamSportGame() *TeamSportGame {
	// Create visualization using VisualizationBuilder
	visConfig := game.NewVisualizationBuilder().
		WithCanvas(800, 600).
		WithBackground("#0d7f3e").
		WithUpdateInterval(50).
		// Add field markings
		AddRectangle("", 50, 150, 700, 400, &game.ShapeOptions{
			StrokeColor: "#ffffff",
			StrokeWidth: 2,
		}).
		// Center line
		AddLine("", 400, 150, 400, 550, &game.LineOptions{
			Color: "#ffffff",
			Width: 2,
		}).
		// Team A stamina bar
		AddText("team_a_stamina", "Team A Stamina", 150, 70, &game.TextOptions{
			FontSize:   14,
			Color:      "#ffffff",
			FontFamily: "Arial",
		}).
		AddProgressBar("team_a_stamina", 150, 90, 200, 30, &game.ProgressBarOptions{
			BackgroundColor:  "rgba(255,255,255,0.3)",
			ForegroundColor:  "#4CAF50",
			BorderColor:      "#ffffff",
			BorderWidth:      2,
			ShowLabel:        true,
			MaxValue:         100,
		}).
		// Team B stamina bar
		AddText("team_b_stamina", "Team B Stamina", 450, 70, &game.TextOptions{
			FontSize:   14,
			Color:      "#ffffff",
			FontFamily: "Arial",
		}).
		AddProgressBar("team_b_stamina", 450, 90, 200, 30, &game.ProgressBarOptions{
			BackgroundColor:  "rgba(255,255,255,0.3)",
			ForegroundColor:  "#f44336",
			BorderColor:      "#ffffff",
			BorderWidth:      2,
			ShowLabel:        true,
			MaxValue:         100,
		}).
		// Score display
		AddText("score", "Score: {value}", 400, 50, &game.TextOptions{
			FontSize:   24,
			Color:      "#ffffff",
			FontFamily: "Arial",
		}).
		Build()

	// Create the game using the fluent GameBuilder API
	config := game.NewGameBuilder("team_sport").
		WithDescription("Manage your team - make substitutions to win!").
		WithServerPartition("score").
		WithServerPartition("team_a_stamina").
		WithServerPartition("team_b_stamina").
		WithActionStatePartition("team_a_stamina").
		WithVisualization(visConfig).
		WithSimulation(BuildTeamSportSimulation).
		Build()

	return &TeamSportGame{config: config}
}

// BuildTeamSportSimulation produces the simulation config generator
func BuildTeamSportSimulation() *simulator.ConfigGenerator {
	gen := simulator.NewConfigGenerator()

	// Team A stamina partition - average stamina of all players on team A
	teamAStamina := &simulator.PartitionConfig{
		Name:      "team_a_stamina",
		Iteration: &TeamStaminaIteration{},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {0.0}, // substitution action
			"base_stamina":        {80.0},
			"stamina_decay":       {0.5},
		}),
		InitStateValues:   []float64{80.0}, // Start at 80% stamina
		StateHistoryDepth: 1,
		Seed:              123,
	}
	gen.SetPartition(teamAStamina)

	// Team B stamina partition - average stamina of all players on team B
	teamBStamina := &simulator.PartitionConfig{
		Name:      "team_b_stamina",
		Iteration: &TeamStaminaIteration{},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {0.0}, // substitution action
			"base_stamina":        {75.0},
			"stamina_decay":       {0.5},
		}),
		InitStateValues:   []float64{75.0}, // Start at 75% stamina
		StateHistoryDepth: 1,
		Seed:              124,
	}
	gen.SetPartition(teamBStamina)

	// Score partition - tracks the current score
	score := &simulator.PartitionConfig{
		Name:      "score",
		Iteration: &ScoreIteration{},
		Params: simulator.NewParams(map[string][]float64{
			"scoring_rate": {0.02},
		}),
		InitStateValues:   []float64{0.0}, // Start at 0-0
		StateHistoryDepth: 1,
		Seed:              125,
	}
	gen.SetPartition(score)

	sim := &simulator.SimulationConfig{
		OutputCondition:      &simulator.EveryStepOutputCondition{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{MaxTimeElapsed: 91.0},
		TimestepFunction:     &simulator.ConstantTimestepFunction{Stepsize: 1.0},
		InitTimeValue:        0.0,
	}
	gen.SetSimulation(sim)
	return gen
}

// GetName returns the game name
func (t *TeamSportGame) GetName() string {
	return t.config.Name
}

// GetDescription returns the game description
func (t *TeamSportGame) GetDescription() string {
	return t.config.Description
}

// GetConfig returns the game configuration
func (t *TeamSportGame) GetConfig() *game.GameConfig {
	return t.config
}

// GetRenderer returns the visualization renderer
func (t *TeamSportGame) GetRenderer() game.GameRenderer {
	return &game.GenericRenderer{Config: t.config.VisualizationConfig}
}

// TeamStaminaIteration simulates stamina changes for a team
type TeamStaminaIteration struct{}

func (t *TeamStaminaIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	// No configuration needed
}

func (t *TeamStaminaIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	currentStamina := stateHistories[partitionIndex].CopyStateRow(0)
	baseStamina := params.Get("base_stamina")[0]
	staminaDecay := params.Get("stamina_decay")[0]
	substitutionAction := params.Get("action_state_values")[0]

	// If substitution action is triggered, boost stamina
	if substitutionAction > 0.5 {
		currentStamina[0] = baseStamina // Fresh players on field
	} else {
		// Otherwise, decay stamina over time
		currentStamina[0] = currentStamina[0] - staminaDecay
	}

	// Clamp stamina between 0 and 100
	if currentStamina[0] < 0.0 {
		currentStamina[0] = 0.0
	}
	if currentStamina[0] > 100.0 {
		currentStamina[0] = 100.0
	}

	return currentStamina
}

// ScoreIteration simulates scoring based on stamina differences
type ScoreIteration struct{}

func (s *ScoreIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	// No configuration needed
}

func (s *ScoreIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	currentScore := stateHistories[partitionIndex].CopyStateRow(0)
	scoringRate := params.Get("scoring_rate")[0]

	// Get team stamina values (partition indices: Team A = 0, Team B = 1)
	teamAStamina := stateHistories[0].CopyStateRow(0)[0]
	teamBStamina := stateHistories[1].CopyStateRow(0)[0]

	// More stamina = better chance to score
	// Score changes based on stamina difference
	staminaDiff := teamAStamina - teamBStamina
	scoreChange := staminaDiff * scoringRate

	// Update score
	currentScore[0] += scoreChange

	return currentScore
}

