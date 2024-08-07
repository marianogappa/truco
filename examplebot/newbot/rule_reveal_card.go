package newbot

import (
	"github.com/marianogappa/truco/truco"
)

var (
	ruleRevealCard = rule{
		name:         "ruleRevealCard",
		description:  "Decides whether to reveal a card (or leave)",
		isApplicable: ruleRevealCardIsApplicable,
		dependsOn:    []rule{ruleInitiateTruco},
		run:          ruleRevealCardRun,
	}
)

func init() {
	registerRule(ruleRevealCard)
}

func ruleRevealCardIsApplicable(st state, _ truco.ClientGameState) bool {
	return isPossibleAll(st, truco.REVEAL_CARD, truco.SAY_ME_VOY_AL_MAZO)
}

func ruleRevealCardRun(st state, gs truco.ClientGameState) (ruleResult, error) {
	result := analyzeTruco(st, gs, false)

	switch {
	case result.shouldLeave:
		return ruleResult{
			action:            getAction(st, truco.SAY_ME_VOY_AL_MAZO),
			stateChanges:      []stateChange{},
			resultDescription: result.description,
		}, nil
	default:
		act := getAction(st, truco.REVEAL_CARD).(*truco.ActionRevealCard)
		act.Card = result.revealCard
		return ruleResult{
			action:            act,
			stateChanges:      []stateChange{},
			resultDescription: result.description,
		}, nil
	}
}
