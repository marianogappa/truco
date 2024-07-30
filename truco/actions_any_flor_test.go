package truco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlor(t *testing.T) {
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
			name: "in a starting hand with flor, player can say flor or leave",
			steps: []testStep{
				{
					expectedPossibleActionNamesBefore: []string{SAY_FLOR, SAY_ME_VOY_AL_MAZO},
					action:                            NewActionSayFlor(0),
				},
			},
		},
		{
			name: "if both players have flor, opponent can accept, decline, raise or leave, but nothing else",
			steps: []testStep{
				{
					expectedPossibleActionNamesBefore: []string{SAY_FLOR, SAY_ME_VOY_AL_MAZO},
					action:                            NewActionSayFlor(0),
				},
				{
					expectedPossibleActionNamesBefore: []string{
						SAY_CONTRAFLOR,
						SAY_CONTRAFLOR_AL_RESTO,
						SAY_CON_FLOR_ME_ACHICO,
						SAY_CON_FLOR_QUIERO,
					},
				},
			},
		},
		{
			name: "if first player says flor, envido is finished",
			hands: []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: ORO}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}},  // 26
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}}, // no flor
			},
			steps: []testStep{
				{
					expectedPossibleActionNamesBefore: []string{SAY_FLOR, SAY_ME_VOY_AL_MAZO},
					action:                            NewActionSayFlor(0),
				},
				{
					expectedCustomValidationBeforeAction: func(g *GameState) {
						require.True(t, g.IsEnvidoFinished)
					},
				},
			},
		},
		{
			name: "if first player says flor and opponent does'nt have flor, actions are different, and 3 points are not won yet!",
			hands: []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: ORO}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}},  // 26
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}}, // no flor
			},
			steps: []testStep{
				{
					expectedPossibleActionNamesBefore: []string{SAY_FLOR, SAY_ME_VOY_AL_MAZO},
					action:                            NewActionSayFlor(0),
				},
				{
					expectedPossibleActionNamesBefore: []string{
						REVEAL_CARD,
						REVEAL_CARD,
						REVEAL_CARD,
						SAY_TRUCO,
						SAY_ME_VOY_AL_MAZO,
					},
					expectedCustomValidationBeforeAction: func(g *GameState) {
						require.Equal(t, g.Players[0].Score, 0)
					},
				},
			},
		},
		{
			name: "3 flor points are won when cards are revealed",
			hands: []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: ORO}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}},  // 26
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}}, // no flor
			},
			steps: []testStep{
				{
					action: NewActionSayFlor(0),
				},
				{
					action: NewActionRevealCard(Card{Number: 1, Suit: ORO}, 0),
				},
				{
					action: NewActionRevealCard(Card{Number: 4, Suit: COPA}, 1),
				},
				{
					action: NewActionSayMeVoyAlMazo(0),
				},
				{
					action: NewActionRevealFlorScore(0),
				},
				{
					expectedCustomValidationBeforeAction: func(g *GameState) {
						require.Equal(t, 3, g.Players[0].Score)
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
				{Unrevealed: []Card{{Number: 1, Suit: ORO}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}}, // 26
				{Unrevealed: []Card{{Number: 4, Suit: ORO}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}}, // 35
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

func _p[T any](v T) *T {
	return &v
}
