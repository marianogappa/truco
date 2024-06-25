package truco

type act struct {
	Name string `json:"name"`
}

func (a act) GetName() string {
	return a.Name
}

func (a act) YieldsTurn(g GameState) bool {
	return true
}

func newActionSayEnvido() Action {
	return ActionSayEnvido{act: act{Name: SAY_ENVIDO}}
}

func newActionSayEnvidoNoQuiero() Action {
	return ActionSayEnvidoNoQuiero{act: act{Name: SAY_ENVIDO_NO_QUIERO}}
}

func newActionSayEnvidoQuiero(score int) Action {
	return ActionSayEnvidoQuiero{act: act{Name: SAY_ENVIDO_QUIERO}, Score: score}
}

func newActionSayTrucoQuiero() Action {
	return ActionSayTrucoQuiero{act: act{Name: SAY_TRUCO_QUIERO}}
}

func newActionSayTrucoNoQuiero() Action {
	return ActionSayTrucoNoQuiero{act: act{Name: SAY_TRUCO_NO_QUIERO}}
}

func newActionSayTruco() Action {
	return ActionSayTruco{act: act{Name: SAY_TRUCO}}
}

func newActionSayQuieroRetruco() Action {
	return ActionSayQuieroRetruco{act: act{Name: SAY_QUIERO_RETRUCO}}
}

func newActionSayQuieroValeCuatro() Action {
	return ActionSayQuieroValeCuatro{act: act{Name: SAY_QUIERO_VALE_CUATRO}}
}

func newActionSaySonBuenas() Action {
	return ActionSaySonBuenas{act: act{Name: SAY_SON_BUENAS}}
}

func newActionSaySonMejores(score int) Action {
	return ActionSaySonMejores{act: act{Name: SAY_SON_MEJORES}, Score: score}
}

func newActionRevealCard(card Card) Action {
	return ActionRevealCard{act: act{Name: REVEAL_CARD}, Card: card}
}
