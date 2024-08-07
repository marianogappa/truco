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
		expectedCustomValidationAfterAction  func(*GameState)
	}

	tests := []struct {
		name                   string
		hands                  []Hand
		changeInitialGameState func(*GameState)
		steps                  []testStep
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
		{
			name: "falta envido with son mejores with 15 points",
			changeInitialGameState: func(g *GameState) {
				g.RuleMaxPoints = 15
			},
			hands: []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: COPA}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}}, // 25
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}}, // 31
			},
			steps: []testStep{
				{
					action: NewActionSayFaltaEnvido(0),
				},
				{
					action: NewActionSayEnvidoQuiero(1),
				},
				{
					action: NewActionSayEnvidoScore(0),
				},
				{
					action: NewActionSaySonMejores(1),
				},
				{
					action: NewActionRevealEnvidoScore(1),
					expectedCustomValidationAfterAction: func(g *GameState) {
						require.Equal(t, 15, g.Players[1].Score)
						require.Equal(t, 0, g.Players[0].Score)
					},
				},
			},
		},
		{
			name: "falta envido with son mejores with 30 points (when starting)",
			changeInitialGameState: func(g *GameState) {
				g.RuleMaxPoints = 30
			},
			hands: []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: COPA}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}}, // 25
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}}, // 31
			},
			steps: []testStep{
				{
					action: NewActionSayFaltaEnvido(0),
				},
				{
					action: NewActionSayEnvidoQuiero(1),
				},
				{
					action: NewActionSayEnvidoScore(0),
				},
				{
					action: NewActionSaySonMejores(1),
				},
				{
					action: NewActionSayMeVoyAlMazo(0),
				},
				{
					action: NewActionRevealEnvidoScore(1),
					expectedCustomValidationAfterAction: func(g *GameState) {
						require.Equal(t, 16, g.Players[1].Score) // 15 + 1 for implicitly winning truco
						require.Equal(t, 0, g.Players[0].Score)
					},
				},
			},
		},
		{
			name: "falta envido with son buenas with 15 points",
			changeInitialGameState: func(g *GameState) {
				g.RuleMaxPoints = 15
			},
			hands: []Hand{
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}}, // 31
				{Unrevealed: []Card{{Number: 1, Suit: COPA}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}}, // 25
			},
			steps: []testStep{
				{
					action: NewActionSayFaltaEnvido(0),
				},
				{
					action: NewActionSayEnvidoQuiero(1),
				},
				{
					action: NewActionSayEnvidoScore(0),
				},
				{
					action: NewActionSaySonBuenas(1),
				},
				{
					action: NewActionRevealEnvidoScore(0),
					expectedCustomValidationAfterAction: func(g *GameState) {
						require.Equal(t, 15, g.Players[0].Score)
						require.Equal(t, 0, g.Players[1].Score)
					},
				},
			},
		},
		{
			name: "falta envido with son buenas with 30 points (when starting)",
			changeInitialGameState: func(g *GameState) {
				g.RuleMaxPoints = 30
			},
			hands: []Hand{
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 5, Suit: ORO}, {Number: 6, Suit: ORO}}}, // 31
				{Unrevealed: []Card{{Number: 1, Suit: COPA}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}}, // 25
			},
			steps: []testStep{
				{
					action: NewActionSayFaltaEnvido(0),
				},
				{
					action: NewActionSayEnvidoQuiero(1),
				},
				{
					action: NewActionSayEnvidoScore(0),
				},
				{
					action: NewActionSaySonBuenas(1),
				},
				{
					action: NewActionSayMeVoyAlMazo(0),
				},
				{
					action: NewActionRevealEnvidoScore(0),
					expectedCustomValidationAfterAction: func(g *GameState) {
						require.Equal(t, 15, g.Players[0].Score)
						require.Equal(t, 1, g.Players[1].Score) // +1 for implicitly winning truco
					},
				},
			},
		},
		{
			name: "falta envido with son mejores with 15 points, but winner is about to lose in 1 point",
			changeInitialGameState: func(g *GameState) {
				g.RuleMaxPoints = 15
				g.Players[0].Score = 14
				g.Players[1].Score = 0
			},
			hands: []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: COPA}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}},       // 25
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 3, Suit: ESPADA}, {Number: 7, Suit: ESPADA}}}, // 30
			},
			steps: []testStep{
				{
					action: NewActionSayFaltaEnvido(0),
				},
				{
					action: NewActionSayEnvidoQuiero(1),
				},
				{
					action: NewActionSayEnvidoScore(0),
				},
				{
					action: NewActionSaySonMejores(1),
				},
				{
					action: NewActionRevealCard(Card{Number: 1, Suit: COPA}, 0),
				},
				{
					action: NewActionRevealCard(Card{Number: 3, Suit: ESPADA}, 1),
				},
				{
					action: NewActionRevealCard(Card{Number: 7, Suit: ESPADA}, 1),
					expectedCustomValidationAfterAction: func(g *GameState) {
						require.Equal(t, 14, g.Players[0].Score)
						require.Equal(t, 1, g.Players[1].Score) // Only wins 1 point for falta envido
					},
				},
			},
		},
		{
			name: "falta envido with son mejores with 30 points, loser has 14 points, but winner still wins 15 points because the game goes to 30",
			changeInitialGameState: func(g *GameState) {
				g.RuleMaxPoints = 30
				g.Players[0].Score = 14
				g.Players[1].Score = 0
			},
			hands: []Hand{
				{Unrevealed: []Card{{Number: 1, Suit: COPA}, {Number: 2, Suit: ORO}, {Number: 3, Suit: ORO}}},       // 25
				{Unrevealed: []Card{{Number: 4, Suit: COPA}, {Number: 3, Suit: ESPADA}, {Number: 7, Suit: ESPADA}}}, // 30
			},
			steps: []testStep{
				{
					action: NewActionSayFaltaEnvido(0),
				},
				{
					action: NewActionSayEnvidoQuiero(1),
				},
				{
					action: NewActionSayEnvidoScore(0),
				},
				{
					action: NewActionSaySonMejores(1),
				},
				{
					action: NewActionRevealCard(Card{Number: 1, Suit: COPA}, 0),
				},
				{
					action: NewActionRevealCard(Card{Number: 3, Suit: ESPADA}, 1),
				},
				{
					action: NewActionRevealCard(Card{Number: 7, Suit: ESPADA}, 1),
					expectedCustomValidationAfterAction: func(g *GameState) {
						require.Equal(t, 14, g.Players[0].Score)
						require.Equal(t, 15, g.Players[1].Score) // Wins 15 points because the game goes to 30
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
			gameState := New(withDeck(newTestDeck(tt.hands)), WithFlorEnabled(false))

			if tt.changeInitialGameState != nil {
				tt.changeInitialGameState(gameState)
			}

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

				if step.expectedCustomValidationAfterAction != nil {
					step.expectedCustomValidationAfterAction(gameState)
				}

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
