package truco

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTruco(t *testing.T) {
	gameState := New()
	err := gameState.RunAction(ActionSayRealEnvido{})
	if err != nil {
		t.Fatal(err)
	}
	pretty, err := gameState.PrettyPrint()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("gameState: \n%v\n", pretty)
	t.Fail()
}

func TestInitialOptions(t *testing.T) {
	gameState := New()
	require.Equal(
		t,
		[]string{
			"reveal_card",
			"say_envido",
			"say_real_envido",
			"say_falta_envido",
			"say_truco",
			"me_voy_al_mazo",
		},
		gameState.PossibleActions,
	)
}

func TestAfterRealEnvidoOptions(t *testing.T) {
	gameState := New()
	err := gameState.RunAction(ActionSayRealEnvido{})
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(
		t,
		[]string{
			"say_falta_envido",
			"say_envido_quiero",
			"say_envido_no_quiero",
		},
		gameState.PossibleActions,
	)
}
