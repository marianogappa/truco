package truco

type ActionRevealCard struct {
	act
	Card   Card `json:"card"`
	EnMesa bool `json:"en_mesa"`
	Score  int  `json:"score"`
}

func (a ActionRevealCard) IsPossible(g GameState) bool {
	if g.IsRoundFinished {
		return false
	}
	// If envido was said and it hasn't finished, then the card can't be revealed
	if !g.IsEnvidoFinished && !g.EnvidoSequence.IsEmpty() && !g.EnvidoSequence.IsFinished() {
		return false
	}

	// If truco was said and it hasn't been accepted or rejected, then the card can't be revealed
	if !g.TrucoSequence.IsEmpty() && !g.TrucoSequence.IsFinished() {
		return false
	}

	step := CardRevealSequenceStep{
		card:     a.Card,
		playerID: g.TurnPlayerID,
	}

	return g.CardRevealSequence.CanAddStep(step, g)
}

func (a *ActionRevealCard) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	step := CardRevealSequenceStep{
		card:     a.Card,
		playerID: g.TurnPlayerID,
	}
	g.CardRevealSequence.AddStep(step, *g)
	err := g.Players[g.TurnPlayerID].Hand.RevealCard(a.Card)
	if err != nil {
		return err
	}
	if g.CardRevealSequence.IsFinished() {
		g.IsRoundFinished = true

		var score int

		// Calculate scores. If there was no truco sequence, 1. Else, calculate the cost.
		if g.TrucoSequence.IsEmpty() {
			score = 1
		} else {
			score = g.TrucoSequence.Cost()
		}

		g.Players[g.CardRevealSequence.WinnerPlayerID()].Score += score
		g.RoundsLog[g.RoundNumber].TrucoPoints = score
		g.RoundsLog[g.RoundNumber].TrucoWinnerPlayerID = g.CardRevealSequence.WinnerPlayerID()
	}
	// If both players have revealed a card, then envido cannot be played anymore
	if !g.IsEnvidoFinished && len(g.Players[g.TurnPlayerID].Hand.Revealed) >= 1 && len(g.Players[g.TurnOpponentPlayerID].Hand.Revealed) >= 1 {
		g.IsEnvidoFinished = true
	}
	// Revealing a card may cause the envido score to be revealed
	if g.tryAwardEnvidoPoints(a.PlayerID) {
		a.EnMesa = true
		a.Score = g.Players[a.PlayerID].Hand.EnvidoScore() // it must be the action's player
	}
	// Revealing a card may cause the flor score to be revealed
	if ok, _ := g.tryAwardFlorPoints(); ok {
		a.EnMesa = true
		a.Score = g.Players[a.PlayerID].Hand.FlorScore() // it must be the action's player
	}
	return nil
}

func (a ActionRevealCard) YieldsTurn(g GameState) bool {
	return g.CardRevealSequence.YieldsTurn(g)
}

func (a *ActionRevealCard) Enrich(g GameState) {
	if g.canAwardEnvidoPoints(Hand{Revealed: append(g.Players[g.TurnPlayerID].Hand.Revealed, a.Card)}) {
		a.EnMesa = true
		a.Score = g.Players[g.TurnPlayerID].Hand.EnvidoScore()
	}
}
