package newbot

import (
	"errors"
	"fmt"

	"github.com/marianogappa/truco/truco"
)

var (
	ruleRespondToRealEnvido = rule{
		name:         "ruleRespondToRealEnvido",
		description:  "Responds to a Real Envido action",
		isApplicable: ruleRespondToRealEnvidoIsApplicable,
		dependsOn:    []rule{ruleRespondToEnvido},
		run:          ruleRespondToRealEnvidoRun,
	}
)

func init() {
	registerRule(ruleRespondToRealEnvido)
}

// TODO: rule for when saying no quiero makes you lose the game
// TODO: rule for when less actions are possible (e.g. if opponent said truco.SAY_CONTRAFLOR)

func ruleRespondToRealEnvidoIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_ENVIDO_NO_QUIERO, truco.SAY_ENVIDO_QUIERO, truco.SAY_FALTA_ENVIDO)
}

func ruleRespondToRealEnvidoRun(st state, gs truco.ClientGameState) (ruleResult, error) {
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
			resultDescription: fmt.Sprintf("Responded to Real Envido with quiero because no quiero costs %v and I will lose in %v points.", costOfNoQuiero, pointsToLose),
		}, nil
	}
	if pointsToWin == 1 {
		return ruleResult{
			action:            getAction(st, truco.SAY_FALTA_ENVIDO),
			stateChanges:      []stateChange{},
			resultDescription: "When there's only 1 point left to win, always respond with Falta Envido.",
		}, nil
	}

	decisionTree := map[string]map[string][2]int{
		"low": {
			truco.SAY_ENVIDO_NO_QUIERO: [2]int{0, 28},
			truco.SAY_ENVIDO_QUIERO:    [2]int{29, 30},
			truco.SAY_FALTA_ENVIDO:     [2]int{31, 33},
		},
		"normal": {
			truco.SAY_ENVIDO_NO_QUIERO: [2]int{0, 27},
			truco.SAY_ENVIDO_QUIERO:    [2]int{28, 29},
			truco.SAY_FALTA_ENVIDO:     [2]int{30, 33},
		},
		"high": {
			truco.SAY_ENVIDO_NO_QUIERO: [2]int{0, 26},
			truco.SAY_ENVIDO_QUIERO:    [2]int{27, 28},
			truco.SAY_FALTA_ENVIDO:     [2]int{29, 33},
		},
	}

	decisionTreeForAgg := decisionTree[agg]

	for actionName, scoreRange := range decisionTreeForAgg {
		if envidoScore >= scoreRange[0] && envidoScore <= scoreRange[1] {
			return ruleResult{
				action:            getAction(st, actionName),
				stateChanges:      []stateChange{},
				resultDescription: fmt.Sprintf("Responded to Real Envido with %v given decision tree for %v aggressiveness and envido score of %v.", actionName, agg, envidoScore),
			}, nil
		}
	}

	return ruleResult{}, errors.New("Couldn't find a suitable action to respond to Real Envido.")
}
