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
	if g.IsRoundFinished {
		return false
	}
	if !g.EnvidoSequence.IsEmpty() && !g.IsEnvidoFinished {
		return false
	}
	// Only the player who said "quiero" last can raise the stakes, unless quiero hasn't been said yet,
	// which can only happen if the last action is "truco"
	if !g.IsLastActionOfName(SAY_TRUCO) &&
		(a.GetName() == SAY_QUIERO_RETRUCO || a.GetName() == SAY_QUIERO_VALE_CUATRO) &&
		g.TrucoSequence.QuieroOwnerPlayerID != g.TurnPlayerID {
		return false
	}
	return g.TrucoSequence.CanAddStep(a.GetName())
}

func (g GameState) IsLastActionOfName(name string) bool {
	actionsLog := g.RoundsLog[g.RoundNumber].ActionsLog
	if len(actionsLog) == 0 {
		return false
	}
	lastActionBs := actionsLog[len(actionsLog)-1].Action
	lastAction, err := DeserializeAction(lastActionBs)
	if err != nil {
		return false
	}
	return lastAction.GetName() == name
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
