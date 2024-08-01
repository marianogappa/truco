package truco

type ActionSayMeVoyAlMazo struct {
	act
}

func (a ActionSayMeVoyAlMazo) IsPossible(g GameState) bool {
	if g.IsRoundFinished {
		return false
	}
	if !g.EnvidoSequence.IsEmpty() && !g.IsEnvidoFinished && !g.EnvidoSequence.IsFinished() {
		return false
	}
	if !g.FlorSequence.IsEmpty() && !g.FlorSequence.IsFinished() {
		return false
	}
	if g.IsEnvidoFinished && !g.TrucoSequence.IsEmpty() && !g.TrucoSequence.IsFinished() {
		return false
	}
	if NewActionRevealFlorScore(a.PlayerID).IsPossible(g) || NewActionRevealEnvidoScore(a.PlayerID).IsPossible(g) {
		return false
	}
	return true
}

func (a ActionSayMeVoyAlMazo) Run(g *GameState) error {
	cost := func() int {
		if g.TrucoSequence.IsEmpty() {
			if g.IsEnvidoFinished {
				// In this case:
				// - Envido was played, "no quiero" happened and score was updated already.
				// - Envido was played, "quiero" happened and score either was or will be updated after this.
				// - Envido wasn't played, and it "expired" after first faceoff.
				//
				// In all cases, the cost is 0 for envido, so just 1 point for truco.
				return 1
			}
			// Envido is not finished. This can only happen if it wasn't played but can still be.
			// So the cost is 1 point for envido and 1 point for truco.
			return 2
		}
		// If truco was played, let's start by calculating its cost
		_cost := g.TrucoSequence.Cost()
		// If envido wasn't finished (wasn't played but can still be), we must add 1 point
		if !g.IsEnvidoFinished {
			_cost++
		}
		return _cost
	}()

	g.RoundsLog[g.RoundNumber].TrucoPoints = cost
	g.RoundsLog[g.RoundNumber].TrucoWinnerPlayerID = g.TurnOpponentPlayerID
	g.Players[g.TurnOpponentPlayerID].Score += cost
	g.IsRoundFinished = true
	return nil
}

func (a ActionSayMeVoyAlMazo) GetPriority() int {
	return 2
}

func (a ActionSayMeVoyAlMazo) AllowLowerPriority() bool {
	return true
}
