package newbot

import (
	"errors"
	"fmt"

	"github.com/marianogappa/truco/truco"
)

var (
	ruleRespondToFaltaEnvido = rule{
		name:         "ruleRespondToFaltaEnvido",
		description:  "Responds to a Falta Envido action",
		isApplicable: ruleRespondToFaltaEnvidoIsApplicable,
		dependsOn:    []rule{ruleRespondToRealEnvido},
		run:          ruleRespondToFaltaEnvidoRun,
	}
)

func init() {
	registerRule(ruleRespondToFaltaEnvido)
}

func ruleRespondToFaltaEnvidoIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_ENVIDO_NO_QUIERO, truco.SAY_ENVIDO_QUIERO)
}

func ruleRespondToFaltaEnvidoRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	var (
		agg            = aggresiveness(st)
		envidoScore    = envidoScore(st)
		pointsToLose   = pointsToLose(st)
		pointsToWin    = pointsToWin(gs)
		costOfNoQuiero = getAction(st, truco.SAY_ENVIDO_NO_QUIERO).(*truco.ActionSayEnvidoNoQuiero).Cost
	)

	if costOfNoQuiero >= pointsToLose {
		return ruleResult{
			action:            getAction(st, truco.SAY_ENVIDO_QUIERO),
			stateChanges:      []stateChange{},
			resultDescription: fmt.Sprintf("Responded to Falta Envido with quiero because no quiero costs %v and I will lose in %v points.", costOfNoQuiero, pointsToLose),
		}, nil
	}
	if pointsToWin == 1 {
		return ruleResult{
			action:            getAction(st, truco.SAY_ENVIDO_QUIERO),
			stateChanges:      []stateChange{},
			resultDescription: "When there's only 1 point left to win, always accept Falta Envido.",
		}, nil
	}

	decisionTree := map[string]map[string][2]int{
		"low": {
			truco.SAY_ENVIDO_NO_QUIERO: [2]int{0, 29},
			truco.SAY_ENVIDO_QUIERO:    [2]int{30, 33},
		},
		"normal": {
			truco.SAY_ENVIDO_NO_QUIERO: [2]int{0, 27},
			truco.SAY_ENVIDO_QUIERO:    [2]int{28, 33},
		},
		"high": {
			truco.SAY_ENVIDO_NO_QUIERO: [2]int{0, 26},
			truco.SAY_ENVIDO_QUIERO:    [2]int{27, 33},
		},
	}

	decisionTreeForAgg := decisionTree[agg]

	for actionName, scoreRange := range decisionTreeForAgg {
		if envidoScore >= scoreRange[0] && envidoScore <= scoreRange[1] {
			return ruleResult{
				action:            getAction(st, actionName),
				stateChanges:      []stateChange{},
				resultDescription: fmt.Sprintf("Responded to Falta Envido with %v given decision tree for %v aggressiveness and envido score of %v.", actionName, agg, envidoScore),
			}, nil
		}
	}

	return ruleResult{}, errors.New("Couldn't find a suitable action to respond to Falta Envido.")
}
