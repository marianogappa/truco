package truco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScoreAfterMeVoyAlMazo(t *testing.T) {
	type testStep struct {
		action                           Action
		expectedScorePlayer0AfterRunning int
		expectedScorePlayer1AfterRunning int
	}

	tests := []struct {
		name  string
		hands []Hand
		steps []testStep
	}{
		{
			name: "quiero_vale_cuatro_accepted_leads_to_5_points",
			steps: []testStep{
				{
					action:                           NewActionSayTruco(0),
					expectedScorePlayer0AfterRunning: 0,
					expectedScorePlayer1AfterRunning: 0,
				},
				{
					action:                           NewActionSayQuieroRetruco(1),
					expectedScorePlayer0AfterRunning: 0,
					expectedScorePlayer1AfterRunning: 0,
				},
				{
					action:                           NewActionSayQuieroValeCuatro(0),
					expectedScorePlayer0AfterRunning: 0,
					expectedScorePlayer1AfterRunning: 0,
				},
				{
					action:                           NewActionSayTrucoQuiero(1),
					expectedScorePlayer0AfterRunning: 0,
					expectedScorePlayer1AfterRunning: 0,
				},
				{
					action:                           NewActionSayMeVoyAlMazo(0),
					expectedScorePlayer0AfterRunning: 0,
					expectedScorePlayer1AfterRunning: 5, // 1 for envido, 4 for vale cuatro
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
				err := gameState.RunAction(step.action)
				require.NoError(t, err)
				assert.Equal(t, step.expectedScorePlayer0AfterRunning, gameState.Players[0].Score, "at step %v expected player 0's score to be %v but got %v", i, step.expectedScorePlayer0AfterRunning, gameState.Players[0].Score)
				assert.Equal(t, step.expectedScorePlayer1AfterRunning, gameState.Players[1].Score, "at step %v expected player 1's score to be %v but got %v", i, step.expectedScorePlayer1AfterRunning, gameState.Players[1].Score)
			}
		})
	}
}
