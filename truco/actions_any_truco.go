package truco

type ActionSayTruco struct{ act }
type ActionSayQuieroRetruco struct{ act }
type ActionSayQuieroValeCuatro struct{ act }

func (a ActionSayTruco) IsPossible(g GameState) bool            { return g.AnyTrucoActionIsPossible(a) }
func (a ActionSayQuieroRetruco) IsPossible(g GameState) bool    { return g.AnyTrucoActionIsPossible(a) }
func (a ActionSayQuieroValeCuatro) IsPossible(g GameState) bool { return g.AnyTrucoActionIsPossible(a) }

func (a ActionSayTruco) Run(g *GameState) error            { return g.AnyTrucoActionRunAction(a) }
func (a ActionSayQuieroRetruco) Run(g *GameState) error    { return g.AnyTrucoActionRunAction(a) }
func (a ActionSayQuieroValeCuatro) Run(g *GameState) error { return g.AnyTrucoActionRunAction(a) }

func (g GameState) AnyTrucoActionIsPossible(a Action) bool {
	if !g.EnvidoSequence.IsEmpty() && !g.EnvidoFinished {
		return false
	}
	return g.TrucoSequence.CanAddStep(a.GetName())
}

func (g *GameState) AnyTrucoActionRunAction(at Action) error {
	if !g.AnyTrucoActionIsPossible(at) {
		return errActionNotPossible
	}
	ok := g.TrucoSequence.AddStep(at.GetName())
	if !ok {
		return errActionNotPossible
	}
	return nil
}
