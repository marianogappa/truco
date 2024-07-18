package truco

type ActionConfirmRoundFinished struct {
	act
}

func (a ActionConfirmRoundFinished) IsPossible(g GameState) bool {
	return g.IsRoundFinished &&
		!NewActionRevealEnvidoScore(a.PlayerID, 0).IsPossible(g) &&
		!g.RoundFinishedConfirmedPlayerIDs[a.PlayerID]
}

func (a ActionConfirmRoundFinished) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.RoundFinishedConfirmedPlayerIDs[a.PlayerID] = true
	return nil
}

func (a ActionConfirmRoundFinished) YieldsTurn(g GameState) bool {
	// The turn should go to the player who is left to confirm the round finished
	return a.PlayerID == g.TurnPlayerID
}
