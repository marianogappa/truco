package truco

type ActionSayEnvido struct{ act }
type ActionSayFaltaEnvido struct{ act }
type ActionSayRealEnvido struct{ act }

func (a ActionSayEnvido) IsPossible(g GameState) bool      { return g.AnyEnvidoActionTypeIsPossible(a) }
func (a ActionSayFaltaEnvido) IsPossible(g GameState) bool { return g.AnyEnvidoActionTypeIsPossible(a) }
func (a ActionSayRealEnvido) IsPossible(g GameState) bool  { return g.AnyEnvidoActionTypeIsPossible(a) }

func (a ActionSayEnvido) Run(g *GameState) error      { return g.AnyEnvidoActionTypeRunAction(a) }
func (a ActionSayFaltaEnvido) Run(g *GameState) error { return g.AnyEnvidoActionTypeRunAction(a) }
func (a ActionSayRealEnvido) Run(g *GameState) error  { return g.AnyEnvidoActionTypeRunAction(a) }

func (g GameState) AnyEnvidoActionTypeIsPossible(a Action) bool {
	if g.IsEnvidoFinished {
		return false
	}
	// If there was a "truco" and an answer to it, regardless when, envido is not possible anymore.
	if len(g.TrucoSequence.Sequence) >= 2 {
		return false
	}
	// If the initial two cards have been revealed, envido is finished
	if len(g.CardRevealSequence.Steps) > 2 {
		return false
	}
	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (g *GameState) AnyEnvidoActionTypeRunAction(a Action) error {
	if g.IsEnvidoFinished {
		return errEnvidoFinished
	}
	if !g.AnyEnvidoActionTypeIsPossible(a) {
		return errActionNotPossible
	}
	if g.EnvidoSequence.IsEmpty() {
		g.EnvidoSequence.StartingPlayerID = g.TurnPlayerID
	}
	ok := g.EnvidoSequence.AddStep(a.GetName())
	if !ok {
		return errActionNotPossible
	}
	return nil
}
