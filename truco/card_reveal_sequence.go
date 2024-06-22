package truco

import (
	"fmt"
)

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
	case 2: // If there are two steps, the third step must be from the first faceoff winner
		return step.playerID == crs.BistepWinners[0]
	case 3: // If there are 3 steps, the 4th step must be from the first faceoff winner's opponent
		return step.playerID == g.OpponentOf(crs.BistepWinners[0])
	case 4: // If there are 4 steps, the 5th step must be from the second faceoff winner
		return step.playerID == crs.BistepWinners[1]
	case 5: // If there are 5 steps, the 6th step must be from the second faceoff winner's opponent
		return step.playerID == g.OpponentOf(crs.BistepWinners[1])
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

	if len(crs.Steps)%2 == 0 && len(crs.Steps) > 0 {
		previousStep := crs.Steps[len(crs.Steps)-2]
		comparisonResult := step.card.CompareTrucoScore(previousStep.card)
		if comparisonResult == 1 || (comparisonResult == 0 && step.playerID == g.RoundTurnPlayerID) {
			fmt.Printf("Player %d won the faceoff, because their card %v beats %d's card %v\n", step.playerID, step.card, previousStep.playerID, previousStep.card)
			crs.BistepWinners = append(crs.BistepWinners, step.playerID)
		} else {
			fmt.Printf("Player %d won the faceoff, because their card %v beats %d's card %v\n", previousStep.playerID, previousStep.card, step.playerID, step.card)
			crs.BistepWinners = append(crs.BistepWinners, previousStep.playerID)
		}
	}

	return true
}

func (crs CardRevealSequence) IsFinished() bool {
	if len(crs.BistepWinners) < 2 {
		return false
	}
	if len(crs.BistepWinners) == 2 && crs.BistepWinners[0] != crs.BistepWinners[1] {
		return false
	}
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
		winsByPlayer[winner]++
	}
	for playerID, wins := range winsByPlayer {
		if wins >= 2 {
			return playerID
		}
	}
	return -1 // Unreachable
}
