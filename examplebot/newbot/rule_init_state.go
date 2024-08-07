package newbot

import (
	"fmt"

	"github.com/marianogappa/truco/truco"
)

var (
	ruleInitState = rule{
		name:         "ruleInitState",
		description:  "Initialises the bot's state",
		isApplicable: ruleInitStateIsApplicable,
		dependsOn:    []rule{},
		run:          ruleInitStateRun,
	}
)

func init() {
	registerRule(ruleInitState)
}

func ruleInitStateIsApplicable(state, truco.ClientGameState) bool {
	return true
}

func ruleInitStateRun(_ state, gs truco.ClientGameState) (ruleResult, error) {
	var (
		aggresiveness         = calculateAggresiveness(gs)
		possibleActions       = possibleActionsMap(gs)
		possibleActionNameSet = possibleActionNameSet(possibleActions)
		envidoScore           = calculateEnvidoScore(gs)
		florScore             = calculateFlorScore(gs)
		faceoffResults        = calculateFaceoffResults(gs)
		pointsToLose          = calculatePointsToLose(gs)
	)

	var (
		stateChangeAggresiveness = stateChange{
			fn: func(st *state) {
				(*st)["aggresiveness"] = aggresiveness
			},
			description: fmt.Sprintf("Set aggresiveness to %v", aggresiveness),
		}
		statePossibleActions = stateChange{
			fn: func(st *state) {
				(*st)["possibleActions"] = possibleActions
			},
			description: fmt.Sprintf("Set possibleActions to %v", possibleActions),
		}
		statePossibleActionNameSet = stateChange{
			fn: func(st *state) {
				(*st)["possibleActionNameSet"] = possibleActionNameSet
			},
			description: fmt.Sprintf("Set possibleActionNameSet to %v", possibleActionNameSet),
		}
		stateEnvidoScore = stateChange{
			fn: func(st *state) {
				(*st)["envidoScore"] = envidoScore
			},
			description: fmt.Sprintf("Set envidoScore to %v", envidoScore),
		}
		stateFlorScore = stateChange{
			fn: func(st *state) {
				(*st)["florScore"] = florScore
			},
			description: fmt.Sprintf("Set florScore to %v", florScore),
		}
		stateFaceoffResults = stateChange{
			fn: func(st *state) {
				(*st)["faceoffResults"] = faceoffResults
			},
			description: fmt.Sprintf("Set faceoffResults to %v", faceoffResults),
		}
		statePointsToLose = stateChange{
			fn: func(st *state) {
				(*st)["pointsToLose"] = pointsToLose
			},
			description: fmt.Sprintf("Set pointsToLose to %v", pointsToLose),
		}
	)

	return ruleResult{
		action: nil,
		stateChanges: []stateChange{
			stateChangeAggresiveness,
			statePossibleActions,
			statePossibleActionNameSet,
			stateEnvidoScore,
			stateFlorScore,
			stateFaceoffResults,
			statePointsToLose,
		},
		resultDescription: "Initialised bot's state.",
	}, nil
}

func calculateAggresiveness(gs truco.ClientGameState) string {
	aggresiveness := "normal"
	if gs.YourScore-gs.TheirScore >= 5 {
		aggresiveness = "low"
	}
	if gs.YourScore-gs.TheirScore <= -5 {
		aggresiveness = "high"
	}
	return aggresiveness
}

func possibleActionsMap(gs truco.ClientGameState) map[string]truco.Action {
	possibleActions := make(map[string]truco.Action)
	for _, action := range _deserializeActions(gs.PossibleActions) {
		possibleActions[action.GetName()] = action
	}
	return possibleActions
}

func possibleActionNameSet(mp map[string]truco.Action) map[string]struct{} {
	possibleActionNames := make(map[string]struct{})
	for name := range mp {
		possibleActionNames[name] = struct{}{}
	}
	return possibleActionNames
}

func calculateEnvidoScore(gs truco.ClientGameState) int {
	return truco.Hand{Revealed: gs.YourRevealedCards, Unrevealed: gs.YourUnrevealedCards}.EnvidoScore()
}

func calculateFlorScore(gs truco.ClientGameState) int {
	return truco.Hand{Revealed: gs.YourRevealedCards, Unrevealed: gs.YourUnrevealedCards}.FlorScore()
}

func calculateFaceoffResults(gs truco.ClientGameState) []int {
	results := []int{}
	for i := 0; i < min(len(gs.YourRevealedCards), len(gs.TheirRevealedCards)); i++ {
		results = append(results, gs.YourRevealedCards[i].CompareTrucoScore(gs.TheirRevealedCards[i]))
	}
	return results
}

const (
	FACEOFF_WIN  = 1
	FACEOFF_LOSS = -1
	FACEOFF_TIE  = 0
)

func calculatePointsToLose(gs truco.ClientGameState) int {
	return gs.RuleMaxPoints - gs.TheirScore
}

func pointsToWin(gs truco.ClientGameState) int {
	return gs.RuleMaxPoints - gs.YourScore
}

func youMano(gs truco.ClientGameState) bool {
	return gs.RoundTurnPlayerID == gs.YouPlayerID
}
