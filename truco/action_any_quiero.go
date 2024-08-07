package truco

import (
	"slices"
)

type ActionSayEnvidoNoQuiero struct {
	act
	Cost int `json:"cost"`
}
type ActionSayEnvidoQuiero struct {
	act
	Cost int `json:"cost"`
}
type ActionSayTrucoQuiero struct {
	act
	Cost int `json:"cost"`
	// RequiresReminder is true if a player ran say_truco and the other player
	// initiated an envido sequence. This action might seem out of context.
	RequiresReminder bool `json:"requires_reminder"`
}
type ActionSayTrucoNoQuiero struct {
	act
	Cost int `json:"cost"`
	// RequiresReminder is true if a player ran say_truco and the other player
	// initiated an envido sequence. This action might seem out of context.
	RequiresReminder bool `json:"requires_reminder"`
}

func (a ActionSayEnvidoNoQuiero) IsPossible(g GameState) bool {
	if g.IsRoundFinished {
		return false
	}
	if g.IsEnvidoFinished {
		return false
	}
	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (a ActionSayEnvidoQuiero) IsPossible(g GameState) bool {
	if g.IsRoundFinished {
		return false
	}
	if g.IsEnvidoFinished {
		return false
	}
	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (a ActionSayTrucoQuiero) IsPossible(g GameState) bool {
	if g.IsRoundFinished {
		return false
	}
	// Edge case: Truco -> Envido -> ???
	// In this case, until envido is resolved, truco cannot continue
	var (
		me                       = a.PlayerID
		isEnvidoQuieroPossible   = NewActionSayEnvidoQuiero(me).IsPossible(g)
		isSonBuenasPossible      = NewActionSaySonBuenas(me).IsPossible(g)
		isSonMejoresPossible     = NewActionSaySonMejores(me).IsPossible(g)
		isSayEnvidoScorePossible = NewActionSayEnvidoScore(me).IsPossible(g)
	)
	if isEnvidoQuieroPossible || isSonBuenasPossible || isSonMejoresPossible || isSayEnvidoScorePossible {
		return false
	}

	return g.TrucoSequence.CanAddStep(a.GetName())
}

func (a ActionSayTrucoNoQuiero) IsPossible(g GameState) bool {
	if g.IsRoundFinished {
		return false
	}
	// Edge case: Truco -> Envido -> ???
	// In this case, until envido is resolved, truco cannot continue
	var (
		me                       = a.PlayerID
		isEnvidoQuieroPossible   = NewActionSayEnvidoQuiero(me).IsPossible(g)
		isSonBuenasPossible      = NewActionSaySonBuenas(me).IsPossible(g)
		isSonMejoresPossible     = NewActionSaySonMejores(me).IsPossible(g)
		isSayEnvidoScorePossible = NewActionSayEnvidoScore(me).IsPossible(g)
	)
	if isEnvidoQuieroPossible || isSonBuenasPossible || isSonMejoresPossible || isSayEnvidoScorePossible {
		return false
	}

	return g.TrucoSequence.CanAddStep(a.GetName())
}

func (a ActionSayEnvidoNoQuiero) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.EnvidoSequence.AddStep(a.GetName())
	g.IsEnvidoFinished = true
	cost, err := g.EnvidoSequence.Cost(g.RuleMaxPoints, g.Players[g.TurnPlayerID].Score, g.Players[g.TurnOpponentPlayerID].Score, false)
	if err != nil {
		return err
	}
	g.RoundsLog[g.RoundNumber].EnvidoPoints = cost
	g.RoundsLog[g.RoundNumber].EnvidoWinnerPlayerID = g.TurnOpponentPlayerID
	g.Players[g.TurnOpponentPlayerID].Score += cost
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
	g.TrucoSequence.QuieroOwnerPlayerID = g.TurnPlayerID
	return nil
}

func (a ActionSayTrucoNoQuiero) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.TrucoSequence.AddStep(a.GetName())
	g.IsRoundFinished = true
	cost := g.TrucoSequence.Cost()
	g.RoundsLog[g.RoundNumber].TrucoPoints = cost
	g.RoundsLog[g.RoundNumber].TrucoWinnerPlayerID = g.TurnOpponentPlayerID
	g.Players[g.TurnOpponentPlayerID].Score += cost
	return nil
}

func (a ActionSayTrucoQuiero) YieldsTurn(g GameState) bool {
	// Next turn belongs to the player who started the truco
	// "sub-sequence". Thus, yield turn if the current player
	// is not the one who started the sub-sequence.
	return g.TurnPlayerID != g.TrucoSequence.StartingPlayerID
}

func (a ActionSayEnvidoNoQuiero) YieldsTurn(g GameState) bool {
	// In son_buenas/son_mejores/no_quiero, the turn should go to whoever started the sequence
	return g.TurnPlayerID != g.EnvidoSequence.StartingPlayerID
}

func (a ActionSayEnvidoQuiero) YieldsTurn(g GameState) bool {
	// In envido_quiero, the next turn should go to whoever has to reveal the score.
	// This should always be the "mano" player.
	return g.TurnPlayerID != g.RoundTurnPlayerID
}

func (a *ActionSayTrucoQuiero) Enrich(g GameState) {
	a.RequiresReminder = _doesTrucoActionRequireReminder(g)
	quieroSeq, _ := g.TrucoSequence.WithStep(SAY_TRUCO_QUIERO)
	quieroCost := quieroSeq.Cost()
	a.Cost = quieroCost
}

func (a *ActionSayTrucoNoQuiero) Enrich(g GameState) {
	a.RequiresReminder = _doesTrucoActionRequireReminder(g)
	noQuieroSeq, _ := g.TrucoSequence.WithStep(SAY_TRUCO_NO_QUIERO)
	quieroCost := noQuieroSeq.Cost()
	a.Cost = quieroCost
}

func (a *ActionSayEnvidoQuiero) Enrich(g GameState) {
	if !a.IsPossible(g) {
		return
	}
	var (
		youScore         = g.Players[a.GetPlayerID()].Score
		theirScore       = g.Players[g.OpponentOf(a.GetPlayerID())].Score
		quieroSeq, err   = g.EnvidoSequence.WithStep(SAY_ENVIDO_QUIERO)
		quieroCost, err2 = quieroSeq.Cost(g.RuleMaxPoints, youScore, theirScore, true)
	)
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
	a.Cost = quieroCost
}

func (a *ActionSayEnvidoNoQuiero) Enrich(g GameState) {
	if !a.IsPossible(g) {
		return
	}
	var (
		youScore           = g.Players[a.GetPlayerID()].Score
		theirScore         = g.Players[g.OpponentOf(a.GetPlayerID())].Score
		noQuieroSeq, err   = g.EnvidoSequence.WithStep(SAY_ENVIDO_NO_QUIERO)
		noQuieroCost, err2 = noQuieroSeq.Cost(g.RuleMaxPoints, youScore, theirScore, true)
	)
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
	a.Cost = noQuieroCost
}

func _doesTrucoActionRequireReminder(g GameState) bool {
	if len(g.RoundsLog[g.RoundNumber].ActionsLog) == 0 {
		return false
	}
	lastAction := _deserializeCurrentRoundLastAction(g)
	// If the last action wasn't a truco action, then an envido sequence
	// got in the middle of the truco sequence. A reminder is needed.
	return !slices.Contains([]string{SAY_TRUCO, SAY_QUIERO_RETRUCO, SAY_QUIERO_VALE_CUATRO}, lastAction.GetName())
}
