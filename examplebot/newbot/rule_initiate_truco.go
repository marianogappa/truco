package newbot

import (
	"github.com/marianogappa/truco/truco"
)

var (
	ruleInitiateTruco = rule{
		name:         "ruleInitiateTruco",
		description:  "Decides whether to initiate a Truco action",
		isApplicable: ruleInitiateTrucoIsApplicable,
		dependsOn:    []rule{ruleInitiateEnvido},
		run:          ruleInitiateTrucoRun,
	}
)

func init() {
	registerRule(ruleInitiateTruco)
}

func ruleInitiateTrucoIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_TRUCO)
}

func ruleInitiateTrucoRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	result := analyzeTruco(st, gs, false)

	switch {
	case result.shouldInitiate:
		return ruleResult{
			action:            getAction(st, truco.SAY_TRUCO),
			stateChanges:      []stateChange{},
			resultDescription: result.description,
		}, nil
	default:
		return ruleResult{
			action:            nil,
			stateChanges:      []stateChange{},
			resultDescription: result.description,
		}, nil
	}
}
