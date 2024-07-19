package truco

import "fmt"

type ActionSayEnvidoNoQuiero struct{ act }
type ActionSayEnvidoQuiero struct{ act }
type ActionSayEnvidoScore struct {
	act
	Score int `json:"score"`
}
type ActionRevealEnvidoScore struct {
	act
	Score int `json:"score"`
}
type ActionSayTrucoQuiero struct{ act }
type ActionSayTrucoNoQuiero struct{ act }

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
	if !g.IsRoundFinished {
		return false
	}
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
	revealedHand := Hand{Revealed: g.Players[a.PlayerID].Hand.Revealed}
	return revealedHand.EnvidoScore() != g.Players[a.PlayerID].Hand.EnvidoScore()
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
		isSonMejoresPossible     = NewActionSaySonMejores(0, me).IsPossible(g)
		isSayEnvidoScorePossible = NewActionSayEnvidoScore(0, me).IsPossible(g)
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
		isSonMejoresPossible     = NewActionSaySonMejores(0, me).IsPossible(g)
		isSayEnvidoScorePossible = NewActionSayEnvidoScore(0, me).IsPossible(g)
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
	cost, err := g.EnvidoSequence.Cost(g.Players[g.TurnPlayerID].Score, g.Players[g.TurnOpponentPlayerID].Score)
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
		// and reveal the additional cards of this combination
		for i := range is {
			candidateHand.Revealed = append(candidateHand.Revealed, curPlayersHand.Unrevealed[i])
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
			if !g.tryAwardEnvidoPoints() {
				return fmt.Errorf("couldn't award envido score after running reveal envido score action due to a bug, this code should be unreachable")
			}
			return nil
		}
	}
	// we tried all possible reveal combinations, so it should be impossible that we didn't find the right combination!
	return fmt.Errorf("couldn't reveal envido score due to a bug, this code should be unreachable")
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
	cost, err := g.TrucoSequence.Cost()
	if err != nil {
		return err
	}
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

func (a ActionRevealEnvidoScore) YieldsTurn(g GameState) bool {
	// this action doesn't change turn because the round is finished at this point
	// and the current player must confirm round finished right after this action
	return false
}
