package newbot

import (
	"errors"
	"fmt"

	"github.com/marianogappa/truco/truco"
)

var (
	ruleRespondToContraflorAlResto = rule{
		name:         "ruleRespondToContraflorAlResto",
		description:  "Responds to a Contraflor Al Resto action",
		isApplicable: ruleRespondToContraflorAlRestoIsApplicable,
		dependsOn:    []rule{ruleRespondToContraflor},
		run:          ruleRespondToContraflorAlRestoRun,
	}
)

func init() {
	registerRule(ruleRespondToContraflorAlResto)
}

// TODO: rule for when saying no quiero makes you lose the game
// TODO: rule for when less actions are possible (e.g. if opponent said truco.SAY_CONTRAFLOR)

func ruleRespondToContraflorAlRestoIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_CON_FLOR_ME_ACHICO, truco.SAY_CON_FLOR_QUIERO)
}

func ruleRespondToContraflorAlRestoRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	agg := aggresiveness(st)
	florScore := florScore(st)

	decisionTree := map[string]map[string][2]int{
		"low": {
			truco.SAY_CON_FLOR_ME_ACHICO: [2]int{20, 33},
			truco.SAY_CON_FLOR_QUIERO:    [2]int{34, 38},
		},
		"normal": {
			truco.SAY_CON_FLOR_ME_ACHICO: [2]int{20, 30},
			truco.SAY_CON_FLOR_QUIERO:    [2]int{30, 38},
		},
		"high": {
			truco.SAY_CON_FLOR_ME_ACHICO: [2]int{20, 26},
			truco.SAY_CON_FLOR_QUIERO:    [2]int{27, 38},
		},
	}

	decisionTreeForAgg := decisionTree[agg]

	for actionName, scoreRange := range decisionTreeForAgg {
		if florScore >= scoreRange[0] && florScore <= scoreRange[1] {
			return ruleResult{
				action:            getAction(st, actionName),
				stateChanges:      []stateChange{},
				resultDescription: fmt.Sprintf("Responded to Contraflor with %v given decision tree for %v aggressiveness and flor score of %v.", actionName, agg, florScore),
			}, nil
		}
	}

	return ruleResult{}, errors.New("Couldn't find a suitable action to respond to Contraflor.")
}
