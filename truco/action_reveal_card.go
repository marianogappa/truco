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
	if !g.EnvidoFinished && !g.EnvidoSequence.IsEmpty() && !g.EnvidoSequence.IsFinished() {
		return false
	}

	// If truco was said and it hasn't been accepted or rejected, then the card can't be revealed
	if !g.TrucoSequence.IsEmpty() && !g.TrucoSequence.IsFinished() {
		return false
	}

	step := CardRevealSequenceStep{
		card:     a.Card,
		playerID: g.CurrentPlayerID(),
	}

	// Note that CalculatePossibleActions will call this without arguments
	// in this case, let's try all unrevealed cards
	if a.Card == (Card{}) {
		result := false
		for _, card := range g.Hands[g.CurrentPlayerID()].Unrevealed {
			step.card = card
			result = result || g.CardRevealSequence.CanAddStep(step, g)
		}
		return result
	}

	return g.CardRevealSequence.CanAddStep(step, g)
}

func (a ActionRevealCard) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	step := CardRevealSequenceStep{
		card:     a.Card,
		playerID: g.CurrentPlayerID(),
	}
	g.CardRevealSequence.AddStep(step, *g)
	err := g.Hands[g.CurrentPlayerID()].RevealCard(a.Card)
	if err != nil {
		return err
	}
	if g.CardRevealSequence.IsFinished() {
		g.RoundFinished = true

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

		g.Scores[g.CardRevealSequence.WinnerPlayerID()] += score
		g.CurrentRoundResult.TrucoPoints = score
		g.CurrentRoundResult.TrucoWinnerPlayerID = g.CardRevealSequence.WinnerPlayerID()
	}
	// If both players have revealed a card, then envido cannot be played anymore
	if !g.EnvidoFinished && len(g.Hands[g.TurnPlayerID].Revealed) >= 1 && len(g.Hands[g.OpponentOf(g.TurnPlayerID)].Revealed) >= 1 {
		g.EnvidoFinished = true
	}
	return nil
}

func (a ActionRevealCard) YieldsTurn(g GameState) bool {
	return g.CardRevealSequence.YieldsTurn(g)
}
