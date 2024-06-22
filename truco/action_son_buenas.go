package truco

type ActionSaySonBuenas struct {
	act
	Score int `json:"score"`
}

func (a ActionSaySonBuenas) IsPossible(g GameState) bool {
	if g.EnvidoFinished {
		return false
	}
	// TODO: should I allow people to lose voluntarily?
	if g.Hands[g.TurnPlayerID].EnvidoScore() > g.Hands[g.OpponentOf(g.TurnPlayerID)].EnvidoScore() {
		return false
	}
	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (a ActionSaySonBuenas) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.EnvidoSequence.AddStep(a.GetName())
	cost, err := g.EnvidoSequence.Cost(g.CurrentPlayerScore(), g.OpponentPlayerScore())
	if err != nil {
		return err
	}
	curPlayerEnvidoScore := g.Hands[g.CurrentPlayerID()].EnvidoScore()
	oppPlayerEnvidoScore := g.Hands[g.OpponentPlayerID()].EnvidoScore()
	g.ValidSonBuenas = true
	if curPlayerEnvidoScore > oppPlayerEnvidoScore || (curPlayerEnvidoScore == oppPlayerEnvidoScore && g.TurnPlayerID == g.CurrentPlayerID()) {
		g.ValidSonBuenas = false
	}
	g.CurrentRoundResult.EnvidoPoints = cost
	g.CurrentRoundResult.EnvidoWinnerPlayerID = g.OpponentOf(g.CurrentPlayerID())
	g.EnvidoWinnerPlayerID = g.OpponentOf(g.CurrentPlayerID())
	g.EnvidoFinished = true
	g.Scores[g.OpponentPlayerID()] += cost
	return nil
}

func (a ActionSaySonBuenas) YieldsTurn(g GameState) bool {
	// In son_buenas/son_mejores, the turn should go to whoever started the sequence
	return g.TurnPlayerID != g.EnvidoSequence.StartingPlayerID
}
