package truco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrucoRequiresReminder(t *testing.T) {
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
			name: "truco, then flor, so the truco_quiero should have requires_reminder = true",
			steps: []testStep{
				{
					action: NewActionSayTruco(0),
				},
				{
					action:                         NewActionSayFlor(1),
					expectedPlayerTurnAfterRunning: _p(1),
				},
				{
					expectedCustomValidationBeforeAction: func(g *GameState) {
						actions := g.CalculatePossibleActions()
						found := false
						for _, a := range actions {
							if a.GetName() == SAY_TRUCO_QUIERO {
								found = true
								assert.True(t, a.(*ActionSayTrucoQuiero).RequiresReminder, "expected say_truco_quiero to have requires_reminder = true")
							}
						}
						assert.True(t, found, "expected to find truco_quiero with requires_reminder = true")
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if tt.name != "if first player says flor and opponent does'nt have flor, actions are different" {
			// 	t.Skip()
			// }
			defaultHands := []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: COPA}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}}, // no flor
				{Unrevealed: []Card{{Number: 4, Suit: ORO}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}},  // flor
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
					assert.Equal(t, *step.expectedPlayerTurnAfterRunning, gameState.TurnPlayerID, "at step %v expected player turn %v but got %v", i, *step.expectedPlayerTurnAfterRunning, gameState.TurnPlayerID)
				}

				if step.expectedIsFinishedAfterRunning != nil {
					assert.Equal(t, *step.expectedIsFinishedAfterRunning, gameState.EnvidoSequence.IsFinished(), "at step %v expected isFinished to be %v but wasn't", i, *step.expectedIsFinishedAfterRunning)
				}
			}
		})
	}
}
