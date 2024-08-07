package truco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvidoSequence(t *testing.T) {
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
			name: "cannot start with envido_quiero",
			steps: []testStep{
				{
					action:          NewActionSayEnvidoQuiero(0),
					expectedIsValid: false,
				},
			},
		},
		{
			name: "cannot start with envido_no_quiero",
			steps: []testStep{
				{
					action:          NewActionSayEnvidoNoQuiero(0),
					expectedIsValid: false,
				},
			},
		},
		{
			name: "cannot start with son_buenas",
			steps: []testStep{
				{
					action:          NewActionSaySonBuenas(0),
					expectedIsValid: false,
				},
			},
		},
		{
			name: "cannot start with son_mejores",
			steps: []testStep{
				{
					action:          NewActionSaySonMejores(0),
					expectedIsValid: false,
				},
			},
		},
		{
			name: "envido_quiero is valid after envido",
			steps: []testStep{
				{
					action:                         NewActionSayEnvido(0),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 1,
				},
				{
					action:                         NewActionSayEnvidoQuiero(1),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 0,
				},
			},
		},
		{
			name: "basic envido finished sequence with son mejores",
			steps: []testStep{
				{
					action:                         NewActionSayEnvido(0),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 1,
				},
				{
					action:                         NewActionSayEnvidoQuiero(1),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 0,
				},
				{
					action:                         NewActionSayEnvidoScore(0),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 1,
				},
				{
					action:                         NewActionSaySonMejores(1),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 0, // doesn't yield turn because envido is over, so player who started gets to play
					expectedIsFinishedAfterRunning: true,
					expectedCostAfterRunning:       2,
				},
			},
		},
		{
			name: "basic envido finished sequence with son buenas, but this time player 1 starts",
			steps: []testStep{
				{
					action:       NewActionRevealCard(Card{Number: 1, Suit: ORO}, 0),
					ignoreAction: true,
				},
				{
					action:                         NewActionSayEnvido(1),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 0,
				},
				{
					action:                         NewActionSayEnvidoQuiero(0),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 0,
				},
				{
					action:                         NewActionSayEnvidoScore(0),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 1,
				},
				{
					action:                         NewActionSaySonMejores(1),
					expectedIsValid:                true,
					expectedPlayerTurnAfterRunning: 1, // doesn't yield turn because envido is over, so player who started gets to play
					expectedIsFinishedAfterRunning: true,
					expectedCostAfterRunning:       2,
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
					actualIsValid := gameState.EnvidoSequence.CanAddStep(step.action.GetName())
					require.Equal(t, step.expectedIsValid, actualIsValid, "at step %v expected isValid to be %v but wasn't", i, step.expectedIsValid)
					if !step.expectedIsValid {
						continue
					}
				}

				step.action.Enrich(*gameState)
				err := gameState.RunAction(step.action)
				require.NoError(t, err, "at step %v", i)

				if step.ignoreAction {
					continue
				}

				assert.Equal(t, step.expectedPlayerTurnAfterRunning, gameState.TurnPlayerID, "at step %v expected player turn %v but got %v", i, step.expectedPlayerTurnAfterRunning, gameState.TurnPlayerID)

				assert.Equal(t, step.expectedIsFinishedAfterRunning, gameState.EnvidoSequence.IsFinished(), "at step %v expected isFinished to be %v but wasn't", i, step.expectedIsFinishedAfterRunning)

				if !step.expectedIsFinishedAfterRunning {
					continue
				}

				cost, err := gameState.EnvidoSequence.Cost(gameState.RuleMaxPoints, gameState.Players[gameState.TurnPlayerID].Score, gameState.Players[gameState.TurnOpponentPlayerID].Score, true)
				require.NoError(t, err)
				assert.Equal(t, step.expectedCostAfterRunning, cost, "at step %v expected cost %v but got %v", i, step.expectedCostAfterRunning, cost)
			}
		})
	}
}
