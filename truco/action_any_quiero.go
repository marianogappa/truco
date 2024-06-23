package truco

type ActionSayEnvidoNoQuiero struct{ act }
type ActionSayEnvidoQuiero struct {
	act
	Score int `json:"score"`
}
type ActionSayTrucoQuiero struct{ act }
type ActionSayTrucoNoQuiero struct{ act }

func (a ActionSayEnvidoNoQuiero) IsPossible(g GameState) bool {
	if g.EnvidoFinished {
		return false
	}
	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (a ActionSayEnvidoQuiero) IsPossible(g GameState) bool {
	if g.EnvidoFinished {
		return false
	}
	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (a ActionSayTrucoQuiero) IsPossible(g GameState) bool {
	return g.TrucoSequence.CanAddStep(a.GetName())
}

func (a ActionSayTrucoNoQuiero) IsPossible(g GameState) bool {
	return g.TrucoSequence.CanAddStep(a.GetName())
}

func (a ActionSayEnvidoNoQuiero) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.EnvidoSequence.AddStep(a.GetName())
	g.EnvidoFinished = true
	cost, err := g.EnvidoSequence.Cost(g.CurrentPlayerScore(), g.OpponentPlayerScore())
	if err != nil {
		return err
	}
	g.Scores[g.OpponentPlayerID()] += cost
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
	g.RoundFinished = true
	cost, err := g.TrucoSequence.Cost()
	if err != nil {
		return err
	}
	g.CurrentRoundResult.TrucoPoints = cost
	g.CurrentRoundResult.TrucoWinnerPlayerID = g.OpponentPlayerID()
	g.Scores[g.OpponentPlayerID()] += cost
	return nil
}
