//go:build js && wasm

package main

import (
	"github.com/umbralcalc/dexetera/pkg/simio"
	"github.com/umbralcalc/dexetera/pkg/team_sport"
)

func main() {
	simio.RegisterStep(team_sport.NewTeamSportGame().GetConfig(), "", ":2112")
}

