package newbot

import (
	"errors"
	"fmt"

	"github.com/marianogappa/truco/truco"
)

var (
	ruleRespondToEnvido = rule{
		name:         "ruleRespondToEnvido",
		description:  "Responds to an Envido action",
		isApplicable: ruleRespondToEnvidoIsApplicable,
		dependsOn:    []rule{ruleInitState},
		run:          ruleRespondToEnvidoRun,
	}
)

func init() {
	registerRule(ruleRespondToEnvido)
}

// TODO: replace envido action with falta_envido when with either bot wins, but if they lose they lose less points.

func ruleRespondToEnvidoIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_ENVIDO_NO_QUIERO, truco.SAY_ENVIDO_QUIERO, truco.SAY_REAL_ENVIDO, truco.SAY_FALTA_ENVIDO)
}

func ruleRespondToEnvidoRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	agg := aggresiveness(st)
	envidoScore := envidoScore(st)
	pointsToLose := pointsToLose(st)
	pointsToWin := pointsToWin(gs)
	costOfNoQuiero := getAction(st, truco.SAY_ENVIDO_NO_QUIERO).(*truco.ActionSayEnvidoNoQuiero).Cost

	if costOfNoQuiero >= pointsToLose {
		return ruleResult{
			action:            getAction(st, truco.SAY_ENVIDO_QUIERO),
			stateChanges:      []stateChange{},
			resultDescription: fmt.Sprintf("Responded to Envido with quiero because no quiero costs %v and I will lose in %v points.", costOfNoQuiero, pointsToLose),
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
			truco.SAY_ENVIDO_NO_QUIERO: [2]int{0, 25},
			truco.SAY_ENVIDO_QUIERO:    [2]int{26, 28},
			truco.SAY_REAL_ENVIDO:      [2]int{29, 30},
			truco.SAY_FALTA_ENVIDO:     [2]int{31, 33},
		},
		"normal": {
			truco.SAY_ENVIDO_NO_QUIERO: [2]int{0, 24},
			truco.SAY_ENVIDO_QUIERO:    [2]int{25, 27},
			truco.SAY_REAL_ENVIDO:      [2]int{28, 29},
			truco.SAY_FALTA_ENVIDO:     [2]int{30, 33},
		},
		"high": {
			truco.SAY_ENVIDO_NO_QUIERO: [2]int{0, 22},
			truco.SAY_ENVIDO_QUIERO:    [2]int{23, 26},
			truco.SAY_REAL_ENVIDO:      [2]int{27, 28},
			truco.SAY_FALTA_ENVIDO:     [2]int{29, 33},
		},
	}

	decisionTreeForAgg := decisionTree[agg]

	for actionName, scoreRange := range decisionTreeForAgg {
		if envidoScore >= scoreRange[0] && envidoScore <= scoreRange[1] {

			// Exception: if pointsToWin == 1, SAY_FALTA_ENVIDO is the same as any other action.
			if pointsToLose == 1 {
				actionName = truco.SAY_ENVIDO_QUIERO
			}

			return ruleResult{
				action:            getAction(st, actionName),
				stateChanges:      []stateChange{},
				resultDescription: fmt.Sprintf("Responded to Envido with %v given decision tree for %v aggressiveness and envido score of %v.", actionName, agg, envidoScore),
			}, nil
		}
	}

	return ruleResult{}, errors.New("Couldn't find a suitable action to respond to Envido.")
}
