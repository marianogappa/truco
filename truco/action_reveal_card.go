package truco

type ActionRevealCard struct {
	act
	Card Card `json:"card"`
}

func NewActionRevealCard(card Card) ActionRevealCard {
	return ActionRevealCard{act: act{Name: "reveal_card"}, Card: card}
}

func (a ActionRevealCard) IsPossible(g GameState) bool {
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

func (a ActionRevealCard) Run(g *GameState) error {
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
			cost, err := g.TrucoSequence.Cost()
			if err != nil {
				return err
			}
			score = cost
		}

		g.Players[g.CardRevealSequence.WinnerPlayerID()].Score += score
		g.CurrentRoundResult.TrucoPoints = score
		g.CurrentRoundResult.TrucoWinnerPlayerID = g.CardRevealSequence.WinnerPlayerID()
	}
	// If both players have revealed a card, then envido cannot be played anymore
	if !g.IsEnvidoFinished && len(g.Players[g.TurnPlayerID].Hand.Revealed) >= 1 && len(g.Players[g.TurnOpponentPlayerID].Hand.Revealed) >= 1 {
		g.IsEnvidoFinished = true
	}
	return nil
}

func (a ActionRevealCard) YieldsTurn(g GameState) bool {
	return g.CardRevealSequence.YieldsTurn(g)
}
