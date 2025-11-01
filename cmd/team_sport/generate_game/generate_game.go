package main

import (
	"github.com/umbralcalc/dexetera/pkg/game"
	"github.com/umbralcalc/dexetera/pkg/team_sport"
)

func main() {
	game.GenerateGamePackage(team_sport.NewTeamSportGame(), "team_sport/")
}

