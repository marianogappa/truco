package truco

type ActionSaySonBuenas struct {
	act
	Score int `json:"score"`
}

func (a ActionSaySonBuenas) IsPossible(g GameState) bool {
	if g.IsEnvidoFinished {
		return false
	}

	var (
		mano       = g.RoundTurnPlayerID
		me         = g.TurnPlayerID
		other      = g.TurnOpponentPlayerID
		meScore    = g.Players[me].Hand.EnvidoScore()
		otherScore = g.Players[other].Hand.EnvidoScore()
	)

	// TODO: should I allow people to lose voluntarily?

	if meScore > otherScore {
		return false
	}
	if meScore == otherScore && mano == me {
		return false
	}

	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (a ActionSaySonBuenas) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.EnvidoSequence.AddStep(a.GetName())
	cost, err := g.EnvidoSequence.Cost(g.Players[g.TurnPlayerID].Score, g.Players[g.TurnOpponentPlayerID].Score)
	if err != nil {
		return err
	}
	g.RoundsLog[g.RoundNumber].EnvidoPoints = cost
	g.RoundsLog[g.RoundNumber].EnvidoWinnerPlayerID = g.TurnOpponentPlayerID
	g.IsEnvidoFinished = true
	g.Players[g.TurnOpponentPlayerID].Score += cost
	return nil
}

func (a ActionSaySonBuenas) YieldsTurn(g GameState) bool {
	// In son_buenas/son_mejores/no_quiero, the turn should go to whoever started the sequence
	return g.TurnPlayerID != g.EnvidoSequence.StartingPlayerID
}
