package truco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitialOptions(t *testing.T) {
	gameState := New()

	expectedActions := []Action{
		NewActionRevealCard(gameState.Players[gameState.TurnPlayerID].Hand.Unrevealed[0], 0),
		NewActionRevealCard(gameState.Players[gameState.TurnPlayerID].Hand.Unrevealed[1], 0),
		NewActionRevealCard(gameState.Players[gameState.TurnPlayerID].Hand.Unrevealed[2], 0),
		NewActionSayEnvido(0),
		NewActionSayRealEnvido(0),
		NewActionSayFaltaEnvido(0),
		NewActionSayTruco(0),
		NewActionSayMeVoyAlMazo(0),
	}
	for _, action := range expectedActions {
		action.Enrich(*gameState)
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
		NewActionSayFaltaEnvido(1),
		NewActionSayEnvidoQuiero(1),
		NewActionSayEnvidoNoQuiero(1),
	}

	err := gameState.RunAction(NewActionSayRealEnvido(0))
	if err != nil {
		t.Fatal(err)
	}
	for _, action := range expectedActions {
		action.Enrich(*gameState)
	}
	for i, expectedAction := range _serializeActions(expectedActions) {
		require.Equal(t, string(expectedAction), string(gameState.PossibleActions[i]))
	}
}
