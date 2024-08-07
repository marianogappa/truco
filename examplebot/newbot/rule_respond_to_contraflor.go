package newbot

import (
	"errors"
	"fmt"

	"github.com/marianogappa/truco/truco"
)

var (
	ruleRespondToContraflor = rule{
		name:         "ruleRespondToContraflor",
		description:  "Responds to a Contraflor action",
		isApplicable: ruleRespondToContraflorIsApplicable,
		dependsOn:    []rule{ruleRespondToFlor},
		run:          ruleRespondToContraflorRun,
	}
)

func init() {
	registerRule(ruleRespondToContraflor)
}

// TODO: rule for when saying no quiero makes you lose the game
// TODO: rule for when less actions are possible (e.g. if opponent said truco.SAY_CONTRAFLOR)

func ruleRespondToContraflorIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_CON_FLOR_ME_ACHICO, truco.SAY_CON_FLOR_QUIERO, truco.SAY_CONTRAFLOR_AL_RESTO)
}

func ruleRespondToContraflorRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	agg := aggresiveness(st)
	florScore := florScore(st)

	decisionTree := map[string]map[string][2]int{
		"low": {
			truco.SAY_CON_FLOR_ME_ACHICO:  [2]int{20, 28},
			truco.SAY_CON_FLOR_QUIERO:     [2]int{29, 33},
			truco.SAY_CONTRAFLOR_AL_RESTO: [2]int{34, 38},
		},
		"normal": {
			truco.SAY_CON_FLOR_ME_ACHICO:  [2]int{20, 26},
			truco.SAY_CON_FLOR_QUIERO:     [2]int{27, 31},
			truco.SAY_CONTRAFLOR_AL_RESTO: [2]int{32, 38},
		},
		"high": {
			truco.SAY_CON_FLOR_ME_ACHICO:  [2]int{20, 23},
			truco.SAY_CON_FLOR_QUIERO:     [2]int{24, 27},
			truco.SAY_CONTRAFLOR_AL_RESTO: [2]int{28, 38},
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
