package truco

type CardRevealSequenceStep struct {
	card     Card
	playerID int
}

type CardRevealSequence struct {
	Steps         []CardRevealSequenceStep `json:"steps"`
	BistepWinners []int                    `json:"bistepWinners"`
}

func (crs CardRevealSequence) CanAddStep(step CardRevealSequenceStep, g GameState) bool {
	// Sanity check: the action's player must be the current player
	if g.CurrentPlayerID() != step.playerID {
		return false
	}
	// Sanity check: the card must be in the player's hand, and it must be unrevealed
	if !g.Hands[step.playerID].HasUnrevealedCard(step.card) {
		return false
	}
	// Sanity check: the sequence must not be finished (i.e. neither player must have won)
	if crs.IsFinished() {
		return false
	}
	switch len(crs.Steps) {
	case 0: // Sanity check: if there are no steps, the first step must be from the rounds's first player
		return step.playerID == g.RoundTurnPlayerID
	case 1: // If there is one step, the second step must be from the round's second player
		return step.playerID == g.RoundTurnOpponentPlayerID()
	case 2: // If there are two steps, the third step must be from the first faceoff winner, or round's first player if tied
		if crs.BistepWinners[0] == -1 {
			return step.playerID == g.RoundTurnPlayerID
		}
		return step.playerID == crs.BistepWinners[0]
	case 3: // If there are 3 steps, the 4th step must be from the first faceoff winner's opponent, or round's second player if tied
		if crs.BistepWinners[0] == -1 {
			return step.playerID == g.RoundTurnOpponentPlayerID()
		}
		return step.playerID == g.OpponentOf(crs.BistepWinners[0])
	case 4:
		// This can get tricky so let's outline the cases
		var (
			mano  = g.RoundTurnPlayerID
			other = g.OpponentOf(g.RoundTurnPlayerID)
			tie   = -1
		)

		nextPlayer := map[[2]int]int{
			{tie, tie}:     mano,
			{tie, mano}:    -1, // mano wins
			{tie, other}:   -1, // other wins
			{mano, tie}:    -1, // mano wins
			{mano, mano}:   -1, // mano wins
			{mano, other}:  other,
			{other, tie}:   -1, // other wins
			{other, mano}:  mano,
			{other, other}: -1, // other wins
		}[[2]int{crs.BistepWinners[0], crs.BistepWinners[1]}]

		if nextPlayer == -1 {
			return false
		}

		return step.playerID == nextPlayer
	case 5:
		// The last card can only be thrown by the opponent of whoever played the previous card
		lastStep := crs.Steps[len(crs.Steps)-2]
		return step.playerID != lastStep.playerID
	}

	// if 6 cards were revealed, the sequence is finished (unreachable due to finishing)
	return false
}

func (crs *CardRevealSequence) AddStep(step CardRevealSequenceStep, g GameState) bool {
	if !crs.CanAddStep(step, g) {
		return false
	}
	crs.Steps = append(crs.Steps, step)

	// Edge case: as de espadas may win the round on the 3rd step
	if len(crs.Steps) == 3 && step.card == (Card{Suit: ESPADA, Number: 1}) && crs.BistepWinners[0] == step.playerID {
		crs.BistepWinners = append(crs.BistepWinners, step.playerID)
		return true
	}

	// If there are 2 steps, compare the cards and compute the winner (or tie)
	if len(crs.Steps)%2 == 0 && len(crs.Steps) > 0 {
		previousStep := crs.Steps[len(crs.Steps)-2]
		comparisonResult := step.card.CompareTrucoScore(previousStep.card)
		switch comparisonResult {
		case 1:
			crs.BistepWinners = append(crs.BistepWinners, step.playerID)
		case -1:
			crs.BistepWinners = append(crs.BistepWinners, previousStep.playerID)
		case 0:
			crs.BistepWinners = append(crs.BistepWinners, -1)
		}
	}

	return true
}

func (crs CardRevealSequence) IsFinished() bool {
	if len(crs.BistepWinners) < 2 {
		return false
	}
	// If there are two finished faceoffs
	if len(crs.BistepWinners) == 2 {
		// If each one won one, not finished
		if crs.BistepWinners[0] != crs.BistepWinners[1] && crs.BistepWinners[0] != -1 && crs.BistepWinners[1] != -1 {
			return false
		}
		// If one of them won both, finished
		if crs.BistepWinners[0] == crs.BistepWinners[1] && crs.BistepWinners[0] != -1 {
			return true
		}
		// If both are tied, not finished
		if crs.BistepWinners[0] == crs.BistepWinners[1] && crs.BistepWinners[0] == -1 {
			return false
		}
	}

	// Otherwise, finished
	return true
}

// NOTE: this must be called AFTER AddStep
func (crs CardRevealSequence) YieldsTurn(g GameState) bool {
	// If an even number of cards were revealed, and the last faceoff winner is the current player, the turn is NOT yielded
	// because the winner gets to start the next faceoff
	if len(crs.Steps)%2 == 0 && len(crs.Steps) > 0 && crs.BistepWinners[len(crs.BistepWinners)-1] == crs.Steps[len(crs.Steps)-1].playerID {
		return false
	}
	return true
}

func (crs CardRevealSequence) WinnerPlayerID() int {
	if !crs.IsFinished() {
		// Shouldn't be called if the sequence is not finished
		return -1
	}
	winsByPlayer := map[int]int{}
	for _, winner := range crs.BistepWinners {
		if winner == -1 {
			continue
		}
		winsByPlayer[winner]++
	}
	winningPlayerID := -1
	mostWins := 0
	for playerID, wins := range winsByPlayer {
		if wins > mostWins {
			winningPlayerID = playerID
			mostWins = wins
		}
	}
	return winningPlayerID
}
