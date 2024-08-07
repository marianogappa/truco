package newbot

import (
	"github.com/marianogappa/truco/truco"
)

var (
	ruleRespondToQuieroRetruco = rule{
		name:         "ruleRespondToQuieroRetruco",
		description:  "Responds to a Quiero Retruco action",
		isApplicable: ruleRespondToQuieroRetrucoIsApplicable,
		dependsOn:    []rule{ruleRespondToTruco},
		run:          ruleRespondToQuieroRetrucoRun,
	}
)

func init() {
	registerRule(ruleRespondToQuieroRetruco)
}

func ruleRespondToQuieroRetrucoIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_TRUCO_NO_QUIERO, truco.SAY_TRUCO_QUIERO, truco.SAY_QUIERO_VALE_CUATRO)
}

func ruleRespondToQuieroRetrucoRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	var (
		pointsToLose      = pointsToLose(st)
		costOfNoQuiero    = getAction(st, truco.SAY_TRUCO_NO_QUIERO).(*truco.ActionSayTrucoNoQuiero).Cost
		noQuieroLosesGame = costOfNoQuiero >= pointsToLose
	)

	result := analyzeTruco(st, gs, noQuieroLosesGame)
	switch {
	case result.shouldRaise:
		return ruleResult{
			action:            getAction(st, truco.SAY_QUIERO_VALE_CUATRO),
			stateChanges:      []stateChange{},
			resultDescription: result.description,
		}, nil
	case result.shouldQuiero:
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
