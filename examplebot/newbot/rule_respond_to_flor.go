package newbot

import (
	"errors"
	"fmt"

	"github.com/marianogappa/truco/truco"
)

var (
	ruleRespondToFlor = rule{
		name:         "ruleRespondToFlor",
		description:  "Responds to a Flor action",
		isApplicable: ruleRespondToFlorIsApplicable,
		dependsOn:    []rule{ruleInitState},
		run:          ruleRespondToFlorRun,
	}
)

func init() {
	registerRule(ruleRespondToFlor)
}

// TODO: rule for when saying no quiero makes you lose the game
// TODO: rule for when less actions are possible (e.g. if opponent said truco.SAY_CONTRAFLOR)

func ruleRespondToFlorIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_CON_FLOR_ME_ACHICO, truco.SAY_CON_FLOR_QUIERO, truco.SAY_CONTRAFLOR, truco.SAY_CONTRAFLOR_AL_RESTO)
}

func ruleRespondToFlorRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	agg := aggresiveness(st)
	florScore := florScore(st)

	decisionTree := map[string]map[string][2]int{
		"low": {
			truco.SAY_CON_FLOR_ME_ACHICO:  [2]int{20, 26},
			truco.SAY_CON_FLOR_QUIERO:     [2]int{27, 29},
			truco.SAY_CONTRAFLOR:          [2]int{30, 33},
			truco.SAY_CONTRAFLOR_AL_RESTO: [2]int{34, 38},
		},
		"normal": {
			truco.SAY_CON_FLOR_ME_ACHICO:  [2]int{20, 24},
			truco.SAY_CON_FLOR_QUIERO:     [2]int{25, 27},
			truco.SAY_CONTRAFLOR:          [2]int{28, 30},
			truco.SAY_CONTRAFLOR_AL_RESTO: [2]int{31, 38},
		},
		"high": {
			truco.SAY_CON_FLOR_ME_ACHICO:  [2]int{20, 23},
			truco.SAY_CON_FLOR_QUIERO:     [2]int{24, 25},
			truco.SAY_CONTRAFLOR:          [2]int{26, 28},
			truco.SAY_CONTRAFLOR_AL_RESTO: [2]int{29, 38},
		},
	}

	decisionTreeForAgg := decisionTree[agg]

	for actionName, scoreRange := range decisionTreeForAgg {
		if florScore >= scoreRange[0] && florScore <= scoreRange[1] {
			return ruleResult{
				action:            getAction(st, actionName),
				stateChanges:      []stateChange{},
				resultDescription: fmt.Sprintf("Responded to Flor with %v given decision tree for %v aggressiveness and flor score of %v.", actionName, agg, florScore),
			}, nil
		}
	}

	return ruleResult{}, errors.New("Couldn't find a suitable action to respond to Flor.")
}
