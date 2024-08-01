package truco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvido(t *testing.T) {
	type testStep struct {
		action                               Action
		expectedPlayerTurnAfterRunning       *int
		expectedIsFinishedAfterRunning       *bool
		expectedPossibleActionNamesBefore    []string
		expectedPossibleActionNamesAfter     []string
		expectedCustomValidationBeforeAction func(*GameState)
	}

	tests := []struct {
		name  string
		hands []Hand
		steps []testStep
	}{
		{
			name: "it is still possible to say envido after opponent's first card is revealed",
			steps: []testStep{
				{
					action: NewActionRevealCard(Card{Number: 1, Suit: COPA}, 0),
				},
				{
					expectedPossibleActionNamesBefore: []string{
						REVEAL_CARD,
						REVEAL_CARD,
						REVEAL_CARD,
						SAY_ENVIDO,
						SAY_REAL_ENVIDO,
						SAY_FALTA_ENVIDO,
						SAY_TRUCO,
						SAY_ME_VOY_AL_MAZO,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultHands := []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: COPA}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}},
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}},
			}
			if len(tt.hands) == 0 {
				tt.hands = defaultHands
			}
			gameState := New(withDeck(newTestDeck(tt.hands)), WithFlorEnabled(true))

			require.Equal(t, 0, gameState.TurnPlayerID)

			for i, step := range tt.steps {

				if step.expectedPossibleActionNamesBefore != nil {
					actualAvailableActionNamesBefore := []string{}
					for _, a := range gameState.CalculatePossibleActions() {
						actualAvailableActionNamesBefore = append(actualAvailableActionNamesBefore, a.GetName())
					}
					assert.ElementsMatch(t, step.expectedPossibleActionNamesBefore, actualAvailableActionNamesBefore, "at step %v", i)
				}

				if step.expectedCustomValidationBeforeAction != nil {
					step.expectedCustomValidationBeforeAction(gameState)
				}

				if step.action == nil {
					continue
				}

				step.action.Enrich(*gameState)
				err := gameState.RunAction(step.action)
				require.NoError(t, err, "at step %v", i)

				if step.expectedPossibleActionNamesAfter != nil {
					actualAvailableActionNamesAfter := []string{}
					for _, a := range gameState.CalculatePossibleActions() {
						actualAvailableActionNamesAfter = append(actualAvailableActionNamesAfter, a.GetName())
					}
					assert.ElementsMatch(t, step.expectedPossibleActionNamesAfter, actualAvailableActionNamesAfter, "at step %v", i)
				}

				if step.expectedPlayerTurnAfterRunning != nil {
					assert.Equal(t, *step.expectedPlayerTurnAfterRunning, gameState.TurnPlayerID, "at step %v expected player turn %v but got %v", i, step.expectedPlayerTurnAfterRunning, gameState.TurnPlayerID)
				}

				if step.expectedIsFinishedAfterRunning != nil {
					assert.Equal(t, *step.expectedIsFinishedAfterRunning, gameState.EnvidoSequence.IsFinished(), "at step %v expected isFinished to be %v but wasn't", i, step.expectedIsFinishedAfterRunning)
				}
			}
		})
	}
}
