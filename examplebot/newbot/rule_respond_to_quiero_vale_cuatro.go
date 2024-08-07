package newbot

import (
	"github.com/marianogappa/truco/truco"
)

var (
	ruleRespondToQuieroValeCuatro = rule{
		name:         "ruleRespondToQuieroValeCuatro",
		description:  "Responds to a Quiero Vale Cuatro action",
		isApplicable: ruleRespondToQuieroValeCuatroIsApplicable,
		dependsOn:    []rule{ruleRespondToQuieroRetruco},
		run:          ruleRespondToQuieroValeCuatroRun,
	}
)

func init() {
	registerRule(ruleRespondToQuieroValeCuatro)
}

func ruleRespondToQuieroValeCuatroIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_TRUCO_NO_QUIERO, truco.SAY_TRUCO_QUIERO)
}

func ruleRespondToQuieroValeCuatroRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	var (
		pointsToLose      = pointsToLose(st)
		costOfNoQuiero    = getAction(st, truco.SAY_TRUCO_NO_QUIERO).(*truco.ActionSayTrucoNoQuiero).Cost
		noQuieroLosesGame = costOfNoQuiero >= pointsToLose
	)

	result := analyzeTruco(st, gs, noQuieroLosesGame)
	switch {
	case result.shouldRaise, result.shouldQuiero:
		return ruleResult{
			action:            getAction(st, truco.SAY_TRUCO_QUIERO),
			stateChanges:      []stateChange{},
			resultDescription: result.description,
		}, nil
	default:
		return ruleResult{
			action:            getAction(st, truco.SAY_TRUCO_NO_QUIERO),
			stateChanges:      []stateChange{},
			resultDescription: result.description,
		}, nil
	}
}
