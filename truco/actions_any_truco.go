package truco

type ActionSayTruco struct {
	act
	NoQuieroCost int `json:"noQuieroCost"`
	QuieroCost   int `json:"quieroCost"`
}
type ActionSayQuieroRetruco struct {
	act
	NoQuieroCost int `json:"noQuieroCost"`
	QuieroCost   int `json:"quieroCost"`
}
type ActionSayQuieroValeCuatro struct {
	act
	NoQuieroCost int `json:"noQuieroCost"`
	QuieroCost   int `json:"quieroCost"`
}

func (a ActionSayTruco) IsPossible(g GameState) bool         { return g.AnyTrucoActionIsPossible(&a) }
func (a ActionSayQuieroRetruco) IsPossible(g GameState) bool { return g.AnyTrucoActionIsPossible(&a) }
func (a ActionSayQuieroValeCuatro) IsPossible(g GameState) bool {
	return g.AnyTrucoActionIsPossible(&a)
}

func (a ActionSayTruco) Run(g *GameState) error            { return g.AnyTrucoActionRunAction(&a) }
func (a ActionSayQuieroRetruco) Run(g *GameState) error    { return g.AnyTrucoActionRunAction(&a) }
func (a ActionSayQuieroValeCuatro) Run(g *GameState) error { return g.AnyTrucoActionRunAction(&a) }

func (a *ActionSayTruco) Enrich(g GameState)            { g.AnyTrucoActionTypeEnrich(a) }
func (a *ActionSayQuieroRetruco) Enrich(g GameState)    { g.AnyTrucoActionTypeEnrich(a) }
func (a *ActionSayQuieroValeCuatro) Enrich(g GameState) { g.AnyTrucoActionTypeEnrich(a) }

func (g GameState) AnyTrucoActionIsPossible(a Action) bool {
	if g.IsRoundFinished {
		return false
	}
	if !g.EnvidoSequence.IsEmpty() && !g.IsEnvidoFinished {
		return false
	}
	// Only the player who said "quiero" last can raise the stakes, unless quiero hasn't been said yet
	if (a.GetName() == SAY_QUIERO_RETRUCO || a.GetName() == SAY_QUIERO_VALE_CUATRO) &&
		g.TrucoSequence.QuieroOwnerPlayerID != a.GetPlayerID() &&
		g.TrucoSequence.QuieroOwnerPlayerID != -1 {
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

	// Possible actions are "truco", "quiero retruco" and "quiero vale cuatro", not "quiero"/"no quiero".
	// If this is the first action in a sub-sequence (subsequences are delimited by "quiero" actions),
	// Store the player ID that started the sub-sequence, so that turn can be yielded correctly after
	// a "quiero" action.
	if g.TrucoSequence.IsSubsequenceStart() {
		g.TrucoSequence.StartingPlayerID = g.TurnPlayerID
	}

	return nil
}

func (g GameState) AnyTrucoActionTypeEnrich(a Action) {
	if !a.IsPossible(g) {
		return
	}
	var (
		quieroSeq, _   = g.TrucoSequence.WithStep(SAY_TRUCO_QUIERO)
		quieroCost     = quieroSeq.Cost()
		noQuieroSeq, _ = g.TrucoSequence.WithStep(SAY_TRUCO_NO_QUIERO)
		noQuieroCost   = noQuieroSeq.Cost()
	)

	switch a.GetName() {
	case SAY_TRUCO:
		a.(*ActionSayTruco).QuieroCost = quieroCost
		a.(*ActionSayTruco).NoQuieroCost = noQuieroCost
	case SAY_QUIERO_RETRUCO:
		a.(*ActionSayQuieroRetruco).QuieroCost = quieroCost
		a.(*ActionSayQuieroRetruco).NoQuieroCost = noQuieroCost
	case SAY_QUIERO_VALE_CUATRO:
		a.(*ActionSayQuieroValeCuatro).QuieroCost = quieroCost
		a.(*ActionSayQuieroValeCuatro).NoQuieroCost = noQuieroCost
	}
}
