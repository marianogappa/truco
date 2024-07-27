package truco

type ActionSayEnvido struct {
	act
	NoQuieroCost int `json:"noQuieroCost"`
	QuieroCost   int `json:"quieroCost"`
}
type ActionSayFaltaEnvido struct {
	act
	NoQuieroCost int `json:"noQuieroCost"`
	QuieroCost   int `json:"quieroCost"`
}
type ActionSayRealEnvido struct {
	act
	NoQuieroCost int `json:"noQuieroCost"`
	QuieroCost   int `json:"quieroCost"`
}

func (a ActionSayEnvido) IsPossible(g GameState) bool { return g.AnyEnvidoActionTypeIsPossible(&a) }
func (a ActionSayFaltaEnvido) IsPossible(g GameState) bool {
	return g.AnyEnvidoActionTypeIsPossible(&a)
}
func (a ActionSayRealEnvido) IsPossible(g GameState) bool { return g.AnyEnvidoActionTypeIsPossible(&a) }

func (a ActionSayEnvido) Run(g *GameState) error      { return g.AnyEnvidoActionTypeRunAction(&a) }
func (a ActionSayFaltaEnvido) Run(g *GameState) error { return g.AnyEnvidoActionTypeRunAction(&a) }
func (a ActionSayRealEnvido) Run(g *GameState) error  { return g.AnyEnvidoActionTypeRunAction(&a) }

func (a *ActionSayEnvido) Enrich(g GameState)      { g.AnyEnvidoActionTypeEnrich(a) }
func (a *ActionSayFaltaEnvido) Enrich(g GameState) { g.AnyEnvidoActionTypeEnrich(a) }
func (a *ActionSayRealEnvido) Enrich(g GameState)  { g.AnyEnvidoActionTypeEnrich(a) }

func (g GameState) AnyEnvidoActionTypeIsPossible(a Action) bool {
	if g.IsRoundFinished {
		return false
	}
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

func (g GameState) AnyEnvidoActionTypeEnrich(a Action) {
	if !a.IsPossible(g) {
		return
	}
	var (
		youScore        = g.Players[a.GetPlayerID()].Score
		theirScore      = g.Players[g.OpponentOf(a.GetPlayerID())].Score
		quieroSeq, _    = g.EnvidoSequence.WithStep(SAY_ENVIDO_QUIERO)
		quieroCost, _   = quieroSeq.Cost(g.RuleMaxPoints, youScore, theirScore)
		noQuieroSeq, _  = g.EnvidoSequence.WithStep(SAY_ENVIDO_NO_QUIERO)
		noQuieroCost, _ = noQuieroSeq.Cost(g.RuleMaxPoints, youScore, theirScore)
	)

	switch a.GetName() {
	case SAY_ENVIDO:
		a.(*ActionSayEnvido).QuieroCost = quieroCost
		a.(*ActionSayEnvido).NoQuieroCost = noQuieroCost
	case SAY_FALTA_ENVIDO:
		a.(*ActionSayFaltaEnvido).QuieroCost = quieroCost
		a.(*ActionSayFaltaEnvido).NoQuieroCost = noQuieroCost
	case SAY_REAL_ENVIDO:
		a.(*ActionSayRealEnvido).QuieroCost = quieroCost
		a.(*ActionSayRealEnvido).NoQuieroCost = noQuieroCost
	}
}
