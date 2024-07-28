package truco

import (
	"fmt"
	"strings"
)

type act struct {
	Name     string `json:"name"`
	PlayerID int    `json:"playerID"`

	fmt.Stringer `json:"-"`
}

func (a act) GetName() string {
	return a.Name
}

func (a act) GetPlayerID() int {
	return a.PlayerID
}

func (a act) GetPriority() int {
	return 0
}

// By default, actions don't need to be enriched.
func (a act) Enrich(g GameState) {}

func (a act) String() string {
	name := strings.ReplaceAll(strings.TrimPrefix(a.Name, "say_"), "_", " ")
	return fmt.Sprintf("Player %v says %v", a.PlayerID, name)
}

func (a act) YieldsTurn(g GameState) bool {
	return true
}

func NewActionSayEnvido(playerID int) Action {
	return &ActionSayEnvido{act: act{Name: SAY_ENVIDO, PlayerID: playerID}}
}

func NewActionSayRealEnvido(playerID int) Action {
	return &ActionSayRealEnvido{act: act{Name: SAY_REAL_ENVIDO, PlayerID: playerID}}
}

func NewActionSayFaltaEnvido(playerID int) Action {
	return &ActionSayFaltaEnvido{act: act{Name: SAY_FALTA_ENVIDO, PlayerID: playerID}}
}

func NewActionSayEnvidoNoQuiero(playerID int) Action {
	return &ActionSayEnvidoNoQuiero{act: act{Name: SAY_ENVIDO_NO_QUIERO, PlayerID: playerID}}
}

func NewActionSayEnvidoQuiero(playerID int) Action {
	return &ActionSayEnvidoQuiero{act: act{Name: SAY_ENVIDO_QUIERO, PlayerID: playerID}}
}

func NewActionSayEnvidoScore(playerID int) Action {
	return &ActionSayEnvidoScore{act: act{Name: SAY_ENVIDO_SCORE, PlayerID: playerID}}
}

func NewActionSayTrucoQuiero(playerID int) Action {
	return &ActionSayTrucoQuiero{act: act{Name: SAY_TRUCO_QUIERO, PlayerID: playerID}}
}

func NewActionSayTrucoNoQuiero(playerID int) Action {
	return &ActionSayTrucoNoQuiero{act: act{Name: SAY_TRUCO_NO_QUIERO, PlayerID: playerID}}
}

func NewActionSayTruco(playerID int) Action {
	return &ActionSayTruco{act: act{Name: SAY_TRUCO, PlayerID: playerID}}
}

func NewActionSayQuieroRetruco(playerID int) Action {
	return &ActionSayQuieroRetruco{act: act{Name: SAY_QUIERO_RETRUCO, PlayerID: playerID}}
}

func NewActionSayQuieroValeCuatro(playerID int) Action {
	return &ActionSayQuieroValeCuatro{act: act{Name: SAY_QUIERO_VALE_CUATRO, PlayerID: playerID}}
}

func NewActionSaySonBuenas(playerID int) Action {
	return &ActionSaySonBuenas{act: act{Name: SAY_SON_BUENAS, PlayerID: playerID}}
}

func NewActionSaySonMejores(playerID int) Action {
	return &ActionSaySonMejores{act: act{Name: SAY_SON_MEJORES, PlayerID: playerID}}
}

func NewActionRevealCard(card Card, playerID int) Action {
	return &ActionRevealCard{act: act{Name: REVEAL_CARD, PlayerID: playerID}, Card: card}
}

func NewActionsRevealCards(playerID int, gameState GameState) []Action {
	actions := []Action{}
	for _, card := range gameState.Players[playerID].Hand.Unrevealed {
		actions = append(actions, NewActionRevealCard(card, playerID))
	}
	return actions
}

func NewActionSayMeVoyAlMazo(playerID int) Action {
	return &ActionSayMeVoyAlMazo{act: act{Name: SAY_ME_VOY_AL_MAZO, PlayerID: playerID}}
}

func NewActionConfirmRoundFinished(playerID int) Action {
	return &ActionConfirmRoundFinished{act: act{Name: CONFIRM_ROUND_FINISHED, PlayerID: playerID}}
}

func NewActionRevealEnvidoScore(playerID int) Action {
	return &ActionRevealEnvidoScore{act: act{Name: REVEAL_ENVIDO_SCORE, PlayerID: playerID}}
}

func NewActionSayFlor(playerID int) Action {
	return &ActionSayFlor{act: act{Name: SAY_FLOR, PlayerID: playerID}}
}

func NewActionSayConFlorMeAchico(playerID int) Action {
	return &ActionSayConFlorMeAchico{act: act{Name: SAY_CON_FLOR_ME_ACHICO, PlayerID: playerID}}
}

func NewActionSayContraflor(playerID int) Action {
	return &ActionSayContraflor{act: act{Name: SAY_CONTRAFLOR, PlayerID: playerID}}
}

func NewActionSayContraflorAlResto(playerID int) Action {
	return &ActionSayContraflorAlResto{act: act{Name: SAY_CONTRAFLOR_AL_RESTO, PlayerID: playerID}}
}

func NewActionSayConFlorQuiero(playerID int) Action {
	return &ActionSayConFlorQuiero{act: act{Name: SAY_CON_FLOR_QUIERO, PlayerID: playerID}}
}

func NewActionSayFlorScore(playerID int) Action {
	return &ActionSayFlorScore{act: act{Name: SAY_FLOR_SCORE, PlayerID: playerID}}
}

func NewActionSayFlorSonBuenas(playerID int) Action {
	return &ActionSayFlorSonBuenas{act: act{Name: SAY_FLOR_SON_BUENAS, PlayerID: playerID}}
}

func NewActionSayFlorSonMejores(playerID int) Action {
	return &ActionSayFlorSonMejores{act: act{Name: SAY_FLOR_SON_MEJORES, PlayerID: playerID}}
}

func NewActionRevealFlorScore(playerID int) Action {
	return &ActionRevealFlorScore{act: act{Name: REVEAL_FLOR_SCORE, PlayerID: playerID}}
}

func (a ActionSaySonMejores) String() string {
	return fmt.Sprintf("Player %v says %v son mejores", a.PlayerID, a.Score)
}

func (a ActionRevealEnvidoScore) String() string {
	return fmt.Sprintf("Player %v says %v en mesa", a.PlayerID, a.Score)
}

func (a ActionRevealCard) String() string {
	text := fmt.Sprintf("Player %v reveals %v of %v", a.PlayerID, a.Card.Number, a.Card.Suit)
	if a.EnMesa {
		text = fmt.Sprintf("%v (%v en mesa)", text, a.Score)
	}
	return text
}

func (a ActionConfirmRoundFinished) String() string {
	return fmt.Sprintf("Player %v confirms round finished", a.PlayerID)
}
