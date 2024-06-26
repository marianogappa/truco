package truco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitialOptions(t *testing.T) {
	gameState := New()

	expectedActions := []Action{
		newActionRevealCard(gameState.Players[gameState.TurnPlayerID].Hand.Unrevealed[0]),
		newActionRevealCard(gameState.Players[gameState.TurnPlayerID].Hand.Unrevealed[1]),
		newActionRevealCard(gameState.Players[gameState.TurnPlayerID].Hand.Unrevealed[2]),
		newActionSayEnvido(),
		newActionSayRealEnvido(),
		newActionSayFaltaEnvido(),
		newActionSayTruco(),
		newActionSayMeVoyAlMazo(),
	}

	require.Equal(
		t,
		_serializeActions(expectedActions),
		gameState.PossibleActions,
	)
}

func TestAfterRealEnvidoOptions(t *testing.T) {
	gameState := New()

	expectedActions := []Action{
		newActionSayFaltaEnvido(),
		newActionSayEnvidoQuiero(gameState.Players[gameState.TurnOpponentPlayerID].Hand.EnvidoScore()),
		newActionSayEnvidoNoQuiero(),
	}

	err := gameState.RunAction(newActionSayRealEnvido())
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(
		t,
		_serializeActions(expectedActions),
		gameState.PossibleActions,
	)
}
