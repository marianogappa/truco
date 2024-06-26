package truco

type ActionSayMeVoyAlMazo struct {
	act
}

func (a ActionSayMeVoyAlMazo) IsPossible(g GameState) bool {
	if !g.EnvidoSequence.IsEmpty() && !g.IsEnvidoFinished && !g.EnvidoSequence.IsFinished() {
		return false
	}
	if g.IsEnvidoFinished && !g.TrucoSequence.IsEmpty() && !g.TrucoSequence.IsFinished() {
		return false
	}
	return true
}

func (a ActionSayMeVoyAlMazo) Run(g *GameState) error {
	var cost int
	if g.TrucoSequence.IsEmpty() && g.IsEnvidoFinished {
		// Envido is finished, so either the envido cost was updated already, or it's zero
		cost = 1
	}
	if g.EnvidoSequence.IsEmpty() && g.TrucoSequence.IsEmpty() && !g.IsEnvidoFinished {
		cost = 2
	}
	if g.IsEnvidoFinished && !g.TrucoSequence.IsEmpty() {
		var err error
		cost, err = g.TrucoSequence.Cost()
		if err != nil {
			return err
		}
	}
	g.RoundsLog[g.RoundNumber].TrucoPoints = cost
	g.RoundsLog[g.RoundNumber].TrucoWinnerPlayerID = g.TurnOpponentPlayerID
	g.Players[g.TurnOpponentPlayerID].Score += cost
	g.IsRoundFinished = true
	return nil
}
