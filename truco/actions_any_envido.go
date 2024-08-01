package truco

import "fmt"

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
type ActionSayEnvidoScore struct {
	act
	Score int `json:"score"`
}
type ActionRevealEnvidoScore struct {
	act
	Score int `json:"score"`
}

func (a ActionSayEnvido) IsPossible(g GameState) bool { return g.AnyEnvidoActionTypeIsPossible(&a) }
func (a ActionSayFaltaEnvido) IsPossible(g GameState) bool {
	return g.AnyEnvidoActionTypeIsPossible(&a)
}
func (a ActionSayRealEnvido) IsPossible(g GameState) bool { return g.AnyEnvidoActionTypeIsPossible(&a) }

func (a ActionSayEnvidoScore) IsPossible(g GameState) bool {
	if len(g.RoundsLog[g.RoundNumber].ActionsLog) == 0 {
		return false
	}
	lastAction := _deserializeCurrentRoundLastAction(g)
	if lastAction.GetName() != SAY_ENVIDO_QUIERO {
		return false
	}
	return g.EnvidoSequence.CanAddStep(a.GetName())
}

func (a ActionRevealEnvidoScore) IsPossible(g GameState) bool {
	if !g.EnvidoSequence.WasAccepted() {
		return false
	}
	if g.EnvidoSequence.EnvidoPointsAwarded {
		return false
	}
	roundLog := g.RoundsLog[g.RoundNumber]
	if roundLog.EnvidoWinnerPlayerID != a.PlayerID {
		return false
	}
	if !g.IsRoundFinished && g.Players[a.PlayerID].Score+roundLog.EnvidoPoints < g.RuleMaxPoints {
		return false
	}
	revealedHand := Hand{Revealed: g.Players[a.PlayerID].Hand.Revealed}
	return revealedHand.EnvidoScore() != g.Players[a.PlayerID].Hand.EnvidoScore()
}

func (a ActionSayEnvido) Run(g *GameState) error      { return g.AnyEnvidoActionTypeRunAction(&a) }
func (a ActionSayFaltaEnvido) Run(g *GameState) error { return g.AnyEnvidoActionTypeRunAction(&a) }
func (a ActionSayRealEnvido) Run(g *GameState) error  { return g.AnyEnvidoActionTypeRunAction(&a) }

func (a ActionSayEnvidoScore) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.EnvidoSequence.AddStep(a.GetName())
	return nil
}

func (a ActionRevealEnvidoScore) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	// We need to reveal the least amount of cards such that the envido score is revealed.
	// Since we don't know which cards to reveal, let's try all possible reveal combinations.
	//
	// allPossibleReveals is a `map[unrevealed_len][]map[card_index]struct{}{}`
	//
	// Note: len(unrevealed) == 0 must be impossible if this line is reached
	_s := struct{}{}
	allPossibleReveals := map[int][]map[int]struct{}{
		1: {{0: _s}}, // i.e. if there's only one unrevealed card, only option is to reveal that card
		2: {{0: _s}, {1: _s}, {0: _s, 1: _s}},
		3: {{0: _s}, {1: _s}, {2: _s}, {0: _s, 1: _s}, {0: _s, 2: _s}, {1: _s, 2: _s}},
	}
	curPlayersHand := g.Players[a.PlayerID].Hand

	// for each possible combination of card reveals
	for _, is := range allPossibleReveals[len(curPlayersHand.Unrevealed)] {
		// create a candidate hand but only with reveal cards
		candidateHand := Hand{Revealed: append([]Card{}, curPlayersHand.Revealed...)}
		for i := range curPlayersHand.Unrevealed {
			card := curPlayersHand.Unrevealed[i]
			candidateHand.displayUnrevealedCards = append(candidateHand.displayUnrevealedCards, DisplayCard{Number: card.Number, Suit: card.Suit})
		}

		// and reveal the additional cards of this combination
		for i := range is {
			candidateHand.Revealed = append(candidateHand.Revealed, curPlayersHand.Unrevealed[i])
			candidateHand.displayUnrevealedCards[i].IsHole = true
		}
		// if by revealing these cards we reach the expected envido score, this is the right reveal
		// Note: this is only true if the reveal combinations are sorted by reveal count ascending!
		// Note: we didn't add the unrevealed cards to the candidate hand yet, because we need to
		//       reach the expected envido score only with revealed cards! That's the whole point!
		if candidateHand.EnvidoScore() == curPlayersHand.EnvidoScore() {
			// don't forget to add the unrevealed cards to the candidate hand
			for i := range curPlayersHand.Unrevealed {
				// add all unrevealed cards from the players hand, except if we revealed them
				if _, ok := is[i]; !ok {
					candidateHand.Unrevealed = append(candidateHand.Unrevealed, curPlayersHand.Unrevealed[i])
				}
			}
			// replace hand with our satisfactory candidate hand
			g.Players[a.PlayerID].Hand = &candidateHand
			if !g.tryAwardEnvidoPoints(a.PlayerID) {
				panic("couldn't award envido score after running reveal envido score action due to a bug, this code should be unreachable")
			}
			return nil
		}
	}
	// we tried all possible reveal combinations, so it should be impossible that we didn't find the right combination!
	return fmt.Errorf("couldn't reveal envido score due to a bug, this code should be unreachable")
}

func (a *ActionSayEnvido) Enrich(g GameState)      { g.AnyEnvidoActionTypeEnrich(a) }
func (a *ActionSayFaltaEnvido) Enrich(g GameState) { g.AnyEnvidoActionTypeEnrich(a) }
func (a *ActionSayRealEnvido) Enrich(g GameState)  { g.AnyEnvidoActionTypeEnrich(a) }

func (a *ActionSayEnvidoScore) Enrich(g GameState) {
	a.Score = g.Players[a.PlayerID].Hand.EnvidoScore()
}

func (a *ActionRevealEnvidoScore) Enrich(g GameState) {
	a.Score = g.Players[a.PlayerID].Hand.EnvidoScore()
}

func (a ActionRevealEnvidoScore) YieldsTurn(g GameState) bool {
	// this action doesn't change turn because the round is finished at this point
	// and the current player must confirm round finished right after this action
	return false
}

func (a *ActionSayEnvidoScore) GetPriority() int {
	return 1
}

func (a ActionRevealEnvidoScore) GetPriority() int {
	return 2 // Because it's higher than confirming round finished
}

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
		seq, _             = g.EnvidoSequence.WithStep(a.GetName())
		youScore           = g.Players[a.GetPlayerID()].Score
		theirScore         = g.Players[g.OpponentOf(a.GetPlayerID())].Score
		quieroSeq, err     = seq.WithStep(SAY_ENVIDO_QUIERO)
		quieroCost, err2   = quieroSeq.Cost(g.RuleMaxPoints, youScore, theirScore, true)
		noQuieroSeq, err3  = seq.WithStep(SAY_ENVIDO_NO_QUIERO)
		noQuieroCost, err4 = noQuieroSeq.Cost(g.RuleMaxPoints, youScore, theirScore, true)
	)
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
	if err3 != nil {
		panic(err3)
	}
	if err4 != nil {
		panic(err4)
	}

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
