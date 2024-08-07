package newbot

import (
	"fmt"
	"math/rand"

	"github.com/marianogappa/truco/truco"
)

var (
	ruleInitiateEnvido = rule{
		name:         "ruleInitiateEnvido",
		description:  "Decides whether to initiate an Envido action",
		isApplicable: ruleInitiateEnvidoIsApplicable,
		dependsOn:    []rule{ruleRespondToQuieroValeCuatro},
		run:          ruleInitiateEnvidoRun,
	}
)

func init() {
	registerRule(ruleInitiateEnvido)
}

func ruleInitiateEnvidoIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_ENVIDO, truco.SAY_REAL_ENVIDO, truco.SAY_FALTA_ENVIDO)
}

func ruleInitiateEnvidoRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	agg := aggresiveness(st)
	envidoScore := envidoScore(st)

	decisionTree := map[string]map[string][2]int{
		"low": {
			truco.SAY_ENVIDO:       [2]int{26, 28},
			truco.SAY_REAL_ENVIDO:  [2]int{29, 30},
			truco.SAY_FALTA_ENVIDO: [2]int{31, 33},
		},
		"normal": {
			truco.SAY_ENVIDO:       [2]int{25, 27},
			truco.SAY_REAL_ENVIDO:  [2]int{28, 29},
			truco.SAY_FALTA_ENVIDO: [2]int{30, 33},
		},
		"high": {
			truco.SAY_ENVIDO:       [2]int{23, 25},
			truco.SAY_REAL_ENVIDO:  [2]int{26, 28},
			truco.SAY_FALTA_ENVIDO: [2]int{29, 33},
		},
	}

	decisionTreeForAgg := decisionTree[agg]
	lied := false

	for actionName, scoreRange := range decisionTreeForAgg {
		if envidoScore >= scoreRange[0] && envidoScore <= scoreRange[1] {
			// Lie one out of 3 times
			if rand.Intn(3) == 0 {
				lied = true
				break
			}

			// Exception: if pointsToWin == 1, choose SAY_FALTA_ENVIDO
			if pointsToWin(gs) == 1 {
				actionName = truco.SAY_FALTA_ENVIDO
			}

			return ruleResult{
				action:            getAction(st, actionName),
				stateChanges:      []stateChange{},
				resultDescription: fmt.Sprintf("Decided to initiate %v given decision tree for %v aggressiveness and envido score of %v.", actionName, agg, envidoScore),
			}, nil
		}
	}

	// If didn't lie before, one out of 3 times decide to initiate envido as a lie
	if !lied && rand.Intn(3) == 0 {
		return ruleResult{
			action:            getAction(st, truco.SAY_ENVIDO),
			stateChanges:      []stateChange{},
			resultDescription: fmt.Sprintf("Decided to lie (33%% chance) and say envido even though I shouldn't according to rules."),
		}, nil
	}

	return ruleResult{
		action:            nil,
		stateChanges:      []stateChange{},
		resultDescription: fmt.Sprintf("Decided not to initiate an envido action given decision tree for %v aggressiveness and envido score of %v.", agg, envidoScore),
	}, nil
}
