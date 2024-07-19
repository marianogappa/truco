package truco

import "fmt"

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

func (a act) String() string {
	return fmt.Sprintf("%v runs %v", a.PlayerID, a.Name)
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

func NewActionSayEnvidoScore(score int, playerID int) Action {
	return &ActionSayEnvidoScore{act: act{Name: SAY_ENVIDO_SCORE, PlayerID: playerID}, Score: score}
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

func NewActionSaySonMejores(score int, playerID int) Action {
	return &ActionSaySonMejores{act: act{Name: SAY_SON_MEJORES, PlayerID: playerID}, Score: score}
}

func NewActionRevealCard(card Card, playerID int) Action {
	return &ActionRevealCard{act: act{Name: REVEAL_CARD, PlayerID: playerID}, Card: card}
}

func NewActionSayMeVoyAlMazo(playerID int) Action {
	return &ActionSayMeVoyAlMazo{act: act{Name: SAY_ME_VOY_AL_MAZO, PlayerID: playerID}}
}

func NewActionConfirmRoundFinished(playerID int) Action {
	return &ActionConfirmRoundFinished{act: act{Name: CONFIRM_ROUND_FINISHED, PlayerID: playerID}}
}

func NewActionRevealEnvidoScore(playerID int, score int) Action {
	return &ActionRevealEnvidoScore{act: act{Name: REVEAL_ENVIDO_SCORE, PlayerID: playerID}, Score: score}
}
