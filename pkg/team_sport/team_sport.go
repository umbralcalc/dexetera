package team_sport

import (
	"math"
	"math/rand"

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
		WithUpdateInterval(0).
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
			BackgroundColor: "rgba(255,255,255,0.3)",
			ForegroundColor: "#4CAF50",
			BorderColor:     "#ffffff",
			BorderWidth:     2,
			ShowLabel:       true,
			MaxValue:        100,
		}).
		// Team B stamina bar
		AddText("team_b_stamina", "Team B Stamina", 450, 70, &game.TextOptions{
			FontSize:   14,
			Color:      "#ffffff",
			FontFamily: "Arial",
		}).
		AddProgressBar("team_b_stamina", 450, 90, 200, 30, &game.ProgressBarOptions{
			BackgroundColor: "rgba(255,255,255,0.3)",
			ForegroundColor: "#f44336",
			BorderColor:     "#ffffff",
			BorderWidth:     2,
			ShowLabel:       true,
			MaxValue:        100,
		}).
		// Score display
		AddText("score", "Score: {value}", 400, 50, &game.TextOptions{
			FontSize:   24,
			Color:      "#ffffff",
			FontFamily: "Arial",
		}).
		// Team A substitutions remaining
		AddText("team_a_substitutions", "Subs: {value}", 150, 130, &game.TextOptions{
			FontSize:   12,
			Color:      "#ffffff",
			FontFamily: "Arial",
		}).
		// Team B substitutions remaining
		AddText("team_b_substitutions", "Subs: {value}", 450, 130, &game.TextOptions{
			FontSize:   12,
			Color:      "#ffffff",
			FontFamily: "Arial",
		}).
		// Player markers
		AddPointSet("team_a_players", &game.PointSetOptions{
			FillColor:   "#4CAF50",
			StrokeColor: "#1B5E20",
			StrokeWidth: 2,
			Radius:      8,
		}).
		AddPointSet("team_b_players", &game.PointSetOptions{
			FillColor:   "#f44336",
			StrokeColor: "#B71C1C",
			StrokeWidth: 2,
			Radius:      8,
		}).
		Build()

	// Create the game using the fluent GameBuilder API
	config := game.NewGameBuilder("team_sport").
		WithDescription("Manage your team - make substitutions to win!").
		WithServerPartition("score").
		WithServerPartition("team_a_stamina").
		WithServerPartition("team_b_stamina").
		WithServerPartition("team_a_substitutions").
		WithServerPartition("team_b_substitutions").
		WithServerPartition("team_a_players").
		WithServerPartition("team_b_players").
		WithActionStatePartition("team_a_stamina").
		WithActionStatePartition("team_a_substitutions").
		WithVisualization(visConfig).
		WithSimulation(BuildTeamSportSimulation).
		Build()

	return &TeamSportGame{config: config}
}

// BuildTeamSportSimulation produces the simulation config generator
func BuildTeamSportSimulation() *simulator.ConfigGenerator {
	gen := simulator.NewConfigGenerator()

	playerCount := 11
	teamASpawn := createTeamSpawnPositions(playerCount, 200.0, 220.0, 28.0)
	teamBSpawn := createTeamSpawnPositions(playerCount, 600.0, 220.0, 28.0)

	scorePartitionIndex := 6
	teamAStaminaIndex := 4
	teamBStaminaIndex := 5

	teamAPlayers := &simulator.PartitionConfig{
		Name: "team_a_players",
		Iteration: NewTeamPlayersIteration(
			"team_a",
			teamASpawn,
			1.0,
			200.0,
			610.0,
			200.0,
			500.0,
			6.0,
			3.0,
			101,
			scorePartitionIndex,
			teamAStaminaIndex,
		),
		InitStateValues:   teamASpawn,
		StateHistoryDepth: 1,
		Seed:              101,
	}
	gen.SetPartition(teamAPlayers)

	teamBPlayers := &simulator.PartitionConfig{
		Name: "team_b_players",
		Iteration: NewTeamPlayersIteration(
			"team_b",
			teamBSpawn,
			-1.0,
			190.0,
			600.0,
			200.0,
			500.0,
			6.0,
			3.0,
			102,
			scorePartitionIndex,
			teamBStaminaIndex,
		),
		InitStateValues:   teamBSpawn,
		StateHistoryDepth: 1,
		Seed:              102,
	}
	gen.SetPartition(teamBPlayers)

	// Team A substitutions remaining partition
	teamASubstitutions := &simulator.PartitionConfig{
		Name:      "team_a_substitutions",
		Iteration: &SubstitutionCountIteration{},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {0.0}, // substitution action
			"max_substitutions":   {3.0}, // 3 substitutions allowed
		}),
		InitStateValues:   []float64{3.0}, // Start with 3 substitutions
		StateHistoryDepth: 1,
		Seed:              103,
	}
	gen.SetPartition(teamASubstitutions)

	// Team B substitutions remaining partition
	teamBSubstitutions := &simulator.PartitionConfig{
		Name:      "team_b_substitutions",
		Iteration: &SubstitutionCountIteration{},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {0.0}, // substitution action (not used for team B)
			"max_substitutions":   {3.0}, // 3 substitutions allowed
		}),
		InitStateValues:   []float64{3.0}, // Start with 3 substitutions
		StateHistoryDepth: 1,
		Seed:              104,
	}
	gen.SetPartition(teamBSubstitutions)

	// Team A stamina partition - average stamina of all players on team A
	teamAStamina := &simulator.PartitionConfig{
		Name:      "team_a_stamina",
		Iteration: &TeamStaminaIteration{SubstitutionPartitionIndex: 2},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {0.0}, // substitution action
			"base_stamina":        {85.0},
			"stamina_decay":       {0.25},
		}),
		InitStateValues:   []float64{85.0}, // Start at 80% stamina
		StateHistoryDepth: 1,
		Seed:              105,
	}
	gen.SetPartition(teamAStamina)

	// Team B stamina partition - average stamina of all players on team B
	teamBStamina := &simulator.PartitionConfig{
		Name:      "team_b_stamina",
		Iteration: &TeamStaminaIteration{SubstitutionPartitionIndex: 3},
		Params: simulator.NewParams(map[string][]float64{
			"action_state_values": {0.0}, // substitution action
			"base_stamina":        {85.0},
			"stamina_decay":       {0.25},
		}),
		InitStateValues:   []float64{85.0}, // Start at 80% stamina
		StateHistoryDepth: 1,
		Seed:              106,
	}
	gen.SetPartition(teamBStamina)

	// Score partition - tracks the current score and broadcasting goal events
	score := &simulator.PartitionConfig{
		Name: "score",
		Iteration: &ScoreIteration{
			TeamAPlayersIndex: 0,
			TeamBPlayersIndex: 1,
			TeamAGoalX:        600.0,
			TeamBGoalX:        210.0,
		},
		InitStateValues:   []float64{0.0, 0.0}, // [score_diff, goal_flag]
		StateHistoryDepth: 1,
		Seed:              107,
	}
	gen.SetPartition(score)

	sim := &simulator.SimulationConfig{
		OutputCondition:      &simulator.EveryStepOutputCondition{},
		TerminationCondition: &simulator.TimeElapsedTerminationCondition{MaxTimeElapsed: 10000.0},
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

// createTeamSpawnPositions generates evenly spaced spawn positions for players.
func createTeamSpawnPositions(playerCount int, startX, startY, spacing float64) []float64 {
	positions := make([]float64, playerCount*2)
	for i := 0; i < playerCount; i++ {
		positions[2*i] = startX
		positions[2*i+1] = startY + float64(i)*spacing
	}
	return positions
}

// TeamPlayersIteration moves player markers across the field based on stamina.
type TeamPlayersIteration struct {
	teamName              string
	direction             float64
	minX                  float64
	maxX                  float64
	minY                  float64
	maxY                  float64
	baseSpeed             float64
	jitter                float64
	spawnPositions        []float64
	seed                  int64
	scorePartitionIndex   int
	staminaPartitionIndex int
	rng                   *rand.Rand
}

// NewTeamPlayersIteration constructs a new TeamPlayersIteration with the provided configuration.
func NewTeamPlayersIteration(
	teamName string,
	spawn []float64,
	direction float64,
	minX, maxX, minY, maxY float64,
	baseSpeed float64,
	jitter float64,
	seed int64,
	scorePartitionIndex int,
	staminaPartitionIndex int,
) *TeamPlayersIteration {
	return &TeamPlayersIteration{
		teamName:              teamName,
		direction:             direction,
		minX:                  math.Min(minX, maxX),
		maxX:                  math.Max(minX, maxX),
		minY:                  minY,
		maxY:                  maxY,
		baseSpeed:             baseSpeed,
		jitter:                jitter,
		spawnPositions:        append([]float64(nil), spawn...),
		seed:                  seed,
		scorePartitionIndex:   scorePartitionIndex,
		staminaPartitionIndex: staminaPartitionIndex,
	}
}

func (t *TeamPlayersIteration) Configure(partitionIndex int, settings *simulator.Settings) {
	// Seed the RNG for deterministic jitter per partition
	t.rng = rand.New(rand.NewSource(t.seed))
}

func (t *TeamPlayersIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	state := stateHistories[partitionIndex].CopyStateRow(0)

	// Reset to spawn positions after a goal
	scoreState := stateHistories[t.scorePartitionIndex].CopyStateRow(0)
	goalFlag := scoreState[1]
	if (t.teamName == "team_a" && goalFlag > 0.5) || (t.teamName == "team_b" && goalFlag < -0.5) {
		copy(state, t.spawnPositions)
	}

	stamina := stateHistories[t.staminaPartitionIndex].CopyStateRow(0)[0]
	speed := t.baseSpeed * (stamina / 100.0)
	if speed < 0.3 {
		speed = 0.3
	}

	for i := 0; i < len(state); i += 2 {
		x := state[i]
		y := state[i+1]

		currentSpeed := speed
		if t.direction > 0 {
			if x > t.maxX-60.0 {
				currentSpeed *= 1.5
			}
		} else {
			if x < t.minX+60.0 {
				currentSpeed *= 1.5
			}
		}

		x += t.direction * currentSpeed

		if x < t.minX {
			x = t.minX
		}
		if x > t.maxX {
			x = t.maxX
		}

		if t.jitter > 0 {
			y += (t.rng.Float64() - 0.5) * t.jitter
		}

		if y < t.minY {
			y = t.minY
		}
		if y > t.maxY {
			y = t.maxY
		}

		state[i] = x
		state[i+1] = y
	}

	return state
}

// TeamStaminaIteration simulates stamina changes for a team
type TeamStaminaIteration struct {
	SubstitutionPartitionIndex int
}

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

	substitutionPartitionIndex := t.SubstitutionPartitionIndex

	if substitutionAction > 0.5 && substitutionPartitionIndex >= 0 {
		substitutionsRemaining := stateHistories[substitutionPartitionIndex].CopyStateRow(0)[0]
		if substitutionsRemaining > 0 {
			currentStamina[0] = baseStamina
		} else {
			currentStamina[0] = currentStamina[0] - staminaDecay
		}
	} else {
		currentStamina[0] = currentStamina[0] - staminaDecay
	}

	if currentStamina[0] < 0.0 {
		currentStamina[0] = 0.0
	}
	if currentStamina[0] > 100.0 {
		currentStamina[0] = 100.0
	}

	return currentStamina
}

// ScoreIteration simulates scoring based on player positions and goal events
type ScoreIteration struct {
	TeamAPlayersIndex int
	TeamBPlayersIndex int
	TeamAGoalX        float64
	TeamBGoalX        float64
}

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
	if len(currentScore) < 2 {
		currentScore = append(currentScore, 0.0)
	}

	goalFlag := 0.0

	teamAPlayers := stateHistories[s.TeamAPlayersIndex].CopyStateRow(0)
	teamBPlayers := stateHistories[s.TeamBPlayersIndex].CopyStateRow(0)

	teamAGoal := detectGoal(teamAPlayers, s.TeamAGoalX, 1)
	teamBGoal := detectGoal(teamBPlayers, s.TeamBGoalX, -1)

	switch {
	case teamAGoal && !teamBGoal:
		currentScore[0] += 1.0
		goalFlag = 1.0
	case teamBGoal && !teamAGoal:
		currentScore[0] -= 1.0
		goalFlag = -1.0
	default:
		goalFlag = 0.0
	}

	currentScore[1] = goalFlag
	return currentScore
}

func detectGoal(playerStates []float64, goalX float64, direction int) bool {
	for i := 0; i < len(playerStates); i += 2 {
		x := playerStates[i]
		if direction > 0 {
			if x >= goalX {
				return true
			}
		} else {
			if x <= goalX {
				return true
			}
		}
	}
	return false
}

type SubstitutionCountIteration struct{}

func (s *SubstitutionCountIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	// No configuration needed
}

func (s *SubstitutionCountIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	currentSubs := stateHistories[partitionIndex].CopyStateRow(0)
	substitutionAction := params.Get("action_state_values")[0]
	maxSubstitutions := params.Get("max_substitutions")[0]

	if substitutionAction > 0.5 && currentSubs[0] > 0 {
		currentSubs[0] = currentSubs[0] - 1.0
	}

	if currentSubs[0] < 0.0 {
		currentSubs[0] = 0.0
	}
	if currentSubs[0] > maxSubstitutions {
		currentSubs[0] = maxSubstitutions
	}

	return currentSubs
}
