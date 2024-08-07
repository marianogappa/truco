package newbot

import (
	"github.com/marianogappa/truco/truco"
)

var (
	ruleInitiateFlor = rule{
		name:         "ruleInitiateFlor",
		description:  "Decides whether to initiate Flor action",
		isApplicable: ruleInitiateFlorIsApplicable,
		dependsOn:    []rule{ruleRespondToQuieroValeCuatro},
		run:          ruleInitiateFlorRun,
	}
)

func init() {
	registerRule(ruleInitiateFlor)
}

func ruleInitiateFlorIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.SAY_FLOR)
}

func ruleInitiateFlorRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	return ruleResult{
		action:            getAction(st, truco.SAY_FLOR),
		stateChanges:      []stateChange{},
		resultDescription: "Always initiate Flor",
	}, nil
}
