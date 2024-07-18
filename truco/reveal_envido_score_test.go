package truco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRevealEnvidoScore(t *testing.T) {
	type testStep struct {
		action                         Action
		expectedIsValid                bool
		expectedPlayerTurnAfterRunning int
		ignoreAction                   bool
	}

	tests := []struct {
		name  string
		hands []Hand
		steps []testStep
	}{
		{
			name: "test reveal envido score",
			steps: []testStep{
				{
					action:       NewActionSayEnvido(0),
					ignoreAction: true,
				},
				{
					action:       NewActionSayEnvidoQuiero(1),
					ignoreAction: true,
				},
				{
					action:       NewActionSayEnvidoScore(25, 0),
					ignoreAction: true,
				},
				{
					action:       NewActionSaySonMejores(31, 1),
					ignoreAction: true,
				},
				{
					action:       NewActionSayTruco(0),
					ignoreAction: true,
				},
				{
					action:       NewActionSayTrucoNoQuiero(1),
					ignoreAction: true,
				},
				// At this point the round is finished, so it's valid for player1 to reveal the envido score, but not to confirm round end
				{
					action:          NewActionConfirmRoundFinished(1),
					expectedIsValid: false,
				},
				// Revealing the envido score is valid
				{
					action:                         NewActionRevealEnvidoScore(1, 31),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 1, // doesn't yield turn
				},
				// Now that the envido score is revealed, it's valid to confirm the round end
				{
					action:                         NewActionConfirmRoundFinished(1),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 0, // yields turn
				},
				// Revealing the envido score again is invalid
				{
					action:          NewActionRevealEnvidoScore(1, 31),
					expectedIsValid: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultHands := []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: ORO}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}}, // 25
				{Unrevealed: []Card{{Number: 4, Suit: ORO}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}}, // 31
			}
			if len(tt.hands) == 0 {
				tt.hands = defaultHands
			}
			gameState := New(withDeck(newTestDeck(tt.hands)))

			require.Equal(t, 0, gameState.TurnPlayerID)

			for i, step := range tt.steps {
				if !step.ignoreAction {
					actualIsValid := step.action.IsPossible(*gameState)
					require.Equal(t, step.expectedIsValid, actualIsValid, "at step %v expected isValid to be %v but wasn't", i, step.expectedIsValid)
					if !step.expectedIsValid {
						continue
					}
				}

				err := gameState.RunAction(step.action)
				require.NoError(t, err)

				if step.ignoreAction {
					continue
				}

				assert.Equal(t, step.expectedPlayerTurnAfterRunning, gameState.TurnPlayerID, "at step %v expected player turn %v but got %v", i, step.expectedPlayerTurnAfterRunning, gameState.TurnPlayerID)
			}
		})
	}
}
