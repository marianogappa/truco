package newbot

import (
	"github.com/marianogappa/truco/truco"
)

var (
	ruleRespondToTruco = rule{
		name:         "ruleRespondToTruco",
		description:  "Responds to a Truco action",
		isApplicable: ruleRespondToTrucoIsApplicable,
		dependsOn:    []rule{ruleInitState},
		run:          ruleRespondToTrucoRun,
	}
)

func init() {
	registerRule(ruleRespondToTruco)
}

func ruleRespondToTrucoIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_TRUCO_NO_QUIERO, truco.SAY_TRUCO_QUIERO, truco.SAY_QUIERO_RETRUCO)
}

func ruleRespondToTrucoRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	var (
		pointsToLose      = pointsToLose(st)
		costOfNoQuiero    = getAction(st, truco.SAY_TRUCO_NO_QUIERO).(*truco.ActionSayTrucoNoQuiero).Cost
		noQuieroLosesGame = costOfNoQuiero >= pointsToLose
	)

	result := analyzeTruco(st, gs, noQuieroLosesGame)
	switch {
	case result.shouldRaise:
		return ruleResult{
			action:            getAction(st, truco.SAY_QUIERO_RETRUCO),
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
