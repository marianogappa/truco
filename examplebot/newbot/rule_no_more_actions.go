package newbot

import (
	"github.com/marianogappa/truco/truco"
)

var (
	ruleNoMoreActions = rule{
		name:         "ruleNoMoreActions",
		description:  "Blows up because no more actions are possible",
		isApplicable: ruleNoMoreActionsIsApplicable,
		dependsOn:    []rule{ruleRevealCard},
		run:          ruleNoMoreActionsRun,
	}
)

func init() {
	registerRule(ruleNoMoreActions)
}

func ruleNoMoreActionsIsApplicable(st state, _ truco.ClientGameState) bool {
	return true
}

func ruleNoMoreActionsRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	panic("No more actions are possible.")
}
