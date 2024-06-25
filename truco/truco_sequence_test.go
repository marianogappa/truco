package truco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrucoSequence(t *testing.T) {
	type testStep struct {
		action                         Action
		expectedIsValid                bool
		expectedPlayerTurnAfterRunning int
		expectedIsFinishedAfterRunning bool
		expectedCostAfterRunning       int
		ignoreAction                   bool
	}

	tests := []struct {
		name  string
		hands []Hand
		steps []testStep
	}{
		{
			name: "cannot start with truco_quiero",
			steps: []testStep{
				{
					action:          newActionSayTrucoQuiero(),
					expectedIsValid: false,
				},
			},
		},
		{
			name: "cannot start with truco_no_quiero",
			steps: []testStep{
				{
					action:          newActionSayTrucoNoQuiero(),
					expectedIsValid: false,
				},
			},
		},
		{
			name: "cannot start with quiero retruco",
			steps: []testStep{
				{
					action:          newActionSayQuieroRetruco(),
					expectedIsValid: false,
				},
			},
		},
		{
			name: "cannot start with quiero vale cuatro",
			steps: []testStep{
				{
					action:          newActionSayQuieroValeCuatro(),
					expectedIsValid: false,
				},
			},
		},
		{
			name: "truco_quiero is valid after truco",
			steps: []testStep{
				{
					action:                         newActionSayTruco(),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 1,
				},
				{
					action:                         newActionSayTrucoQuiero(),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 0,
					expectedIsFinishedAfterRunning: true,
					expectedCostAfterRunning:       2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultHands := []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: ORO}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}},
				{Unrevealed: []Card{{Number: 4, Suit: ORO}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}},
			}
			if len(tt.hands) == 0 {
				tt.hands = defaultHands
			}
			gameState := New(withDeck(newTestDeck(tt.hands)))

			require.Equal(t, 0, gameState.TurnPlayerID)

			for i, step := range tt.steps {
				if !step.ignoreAction {
					actualIsValid := gameState.TrucoSequence.CanAddStep(step.action.GetName())
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

				assert.Equal(t, step.expectedIsFinishedAfterRunning, gameState.TrucoSequence.IsFinished(), "at step %v expected isFinished to be %v but wasn't", i, step.expectedIsFinishedAfterRunning)

				if !step.expectedIsFinishedAfterRunning {
					continue
				}

				cost, err := gameState.TrucoSequence.Cost()
				require.NoError(t, err)
				assert.Equal(t, step.expectedCostAfterRunning, cost, "at step %v expected cost %v but got %v", i, step.expectedCostAfterRunning, cost)
			}
		})
	}
}
