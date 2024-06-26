package truco

type ActionSayEnvidoNoQuiero struct{ act }
type ActionSayEnvidoQuiero struct {
	act
	Score int `json:"score"`
}
type ActionSayTrucoQuiero struct{ act }
type ActionSayTrucoNoQuiero struct{ act }

func (a ActionSayEnvidoNoQuiero) IsPossible(g GameState) bool {
	if g.IsEnvidoFinished {
		return false
	}
	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (a ActionSayEnvidoQuiero) IsPossible(g GameState) bool {
	if g.IsEnvidoFinished {
		return false
	}
	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (a ActionSayTrucoQuiero) IsPossible(g GameState) bool {
	// Edge case: Truco -> Envido -> ???
	// In this case, until envido is resolved, truco cannot continue
	actionEnvidoQuiero := ActionSayEnvidoQuiero{act: act{Name: SAY_ENVIDO_QUIERO}}
	actionSonBuenas := ActionSaySonBuenas{act: act{Name: SAY_SON_BUENAS}}
	actionSonMejores := ActionSaySonMejores{act: act{Name: SAY_SON_MEJORES}}
	if actionEnvidoQuiero.IsPossible(g) || actionSonBuenas.IsPossible(g) || actionSonMejores.IsPossible(g) {
		return false
	}

	return g.TrucoSequence.CanAddStep(a.GetName())
}

func (a ActionSayTrucoNoQuiero) IsPossible(g GameState) bool {
	// Edge case: Truco -> Envido -> ???
	// In this case, until envido is resolved, truco cannot continue
	actionEnvidoQuiero := ActionSayEnvidoQuiero{act: act{Name: SAY_ENVIDO_QUIERO}}
	actionSonBuenas := ActionSaySonBuenas{act: act{Name: SAY_SON_BUENAS}}
	actionSonMejores := ActionSaySonMejores{act: act{Name: SAY_SON_MEJORES}}
	if actionEnvidoQuiero.IsPossible(g) || actionSonBuenas.IsPossible(g) || actionSonMejores.IsPossible(g) {
		return false
	}

	return g.TrucoSequence.CanAddStep(a.GetName())
}

func (a ActionSayEnvidoNoQuiero) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.EnvidoSequence.AddStep(a.GetName())
	g.IsEnvidoFinished = true
	cost, err := g.EnvidoSequence.Cost(g.Players[g.TurnPlayerID].Score, g.Players[g.TurnOpponentPlayerID].Score)
	if err != nil {
		return err
	}
	g.Players[g.TurnOpponentPlayerID].Score += cost
	return nil
}

func (a ActionSayEnvidoQuiero) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.EnvidoSequence.AddStep(a.GetName())
	return nil
}

func (a ActionSayTrucoQuiero) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.TrucoSequence.AddStep(a.GetName())
	g.TrucoQuieroOwnerPlayerId = g.TurnPlayerID
	return nil
}

func (a ActionSayTrucoNoQuiero) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.TrucoSequence.AddStep(a.GetName())
	g.IsRoundFinished = true
	cost, err := g.TrucoSequence.Cost()
	if err != nil {
		return err
	}
	g.CurrentRoundResult.TrucoPoints = cost
	g.CurrentRoundResult.TrucoWinnerPlayerID = g.TurnOpponentPlayerID
	g.Players[g.TurnOpponentPlayerID].Score += cost
	return nil
}

func (a ActionSayTrucoQuiero) YieldsTurn(g GameState) bool {
	// Next turn belongs to the player who started the truco
	// "sub-sequence". Thus, yield turn if the current player
	// is not the one who started the sub-sequence.
	return g.TurnPlayerID != g.TrucoSequence.StartingPlayerID
}

func (a ActionSayEnvidoNoQuiero) YieldsTurn(g GameState) bool {
	// In son_buenas/son_mejores/no_quiero, the turn should go to whoever started the sequence
	return g.TurnPlayerID != g.EnvidoSequence.StartingPlayerID
}
