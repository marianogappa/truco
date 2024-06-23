package truco

type ActionSayMeVoyAlMazo struct {
	act
}

func (a ActionSayMeVoyAlMazo) IsPossible(g GameState) bool {
	if !g.EnvidoSequence.IsEmpty() && !g.EnvidoFinished && !g.EnvidoSequence.IsFinished() {
		return false
	}
	if g.EnvidoFinished && !g.TrucoSequence.IsEmpty() && !g.TrucoSequence.IsFinished() {
		return false
	}
	return true
}

func (a ActionSayMeVoyAlMazo) Run(g *GameState) error {
	var cost int
	if g.TrucoSequence.IsEmpty() && g.EnvidoFinished {
		// Envido is finished, so either the envido cost was updated already, or it's zero
		cost = 1
	}
	if g.EnvidoSequence.IsEmpty() && g.TrucoSequence.IsEmpty() && !g.EnvidoFinished {
		cost = 2
	}
	if g.EnvidoFinished && !g.TrucoSequence.IsEmpty() {
		var err error
		cost, err = g.TrucoSequence.Cost()
		if err != nil {
			return err
		}
	}
	g.CurrentRoundResult.TrucoPoints = cost
	g.CurrentRoundResult.TrucoWinnerPlayerID = g.OpponentPlayerID()
	g.Scores[g.OpponentPlayerID()] += cost
	g.RoundFinished = true
	return nil
}
