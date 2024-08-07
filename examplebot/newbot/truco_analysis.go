package newbot

import (
	"fmt"
	"strings"

	"math/rand"

	"github.com/marianogappa/truco/truco"
)

type trucoResult struct {
	shouldQuiero   bool
	shouldInitiate bool
	shouldRaise    bool
	shouldLeave    bool
	revealCard     truco.Card
	description    string
}

func t_raise(c truco.Card, description string) trucoResult {
	return trucoResult{shouldQuiero: true, shouldInitiate: true, shouldRaise: true, revealCard: c, description: description}
}

func t_quiero(c truco.Card, description string) trucoResult {
	return trucoResult{shouldQuiero: true, shouldInitiate: true, shouldRaise: false, revealCard: c, description: description}
}

func t_leave(c truco.Card, description string) trucoResult {
	return trucoResult{shouldQuiero: false, shouldInitiate: false, shouldRaise: false, shouldLeave: true, revealCard: c, description: description}
}

func t_noquiero(c truco.Card, description string) trucoResult {
	return trucoResult{shouldQuiero: false, shouldInitiate: false, shouldRaise: false, revealCard: c, description: description}
}

func analyzeTruco(st state, gs truco.ClientGameState, noQuieroLosesGame bool) trucoResult {
	result := _analyzeTruco(st, gs)
	if noQuieroLosesGame && (!result.shouldQuiero || result.shouldLeave) {
		result.shouldQuiero = true
		result.shouldLeave = false
		result.shouldRaise = true
	}
	// Exception: bot should never raise (but still accept) if points to win is 1
	if pointsToWin(gs) == 1 {
		result.shouldRaise = false
	}
	return result
}

func _analyzeTruco(st state, gs truco.ClientGameState) trucoResult {
	agg := aggresiveness(st)
	faceoffResults := calculateFaceoffResults(gs)

	revealedCardPairs := [2]int{len(gs.YourRevealedCards), len(gs.TheirRevealedCards)}
	switch revealedCardPairs {
	case [2]int{0, 0}, [2]int{1, 0}: // No cards on table, or only our card on table
		var revealCard truco.Card
		if len(gs.YourUnrevealedCards) == 1 {
			revealCard = gs.YourRevealedCards[1]
		} else {
			revealCard = randomCardExceptSpecial(gs.YourUnrevealedCards)
		}
		power := newPower(gs.YourUnrevealedCards)
		switch agg {
		case "low":
			switch {
			case power.gte(TWO_GOOD_ONE_MEDIUM):
				return t_raise(revealCard, "no cards on table (or mine only), low agg, two good cards and one medium or better, good to raise")
			case (power.gte(ONE_GOOD_TWO_MEDIUM) && youMano(gs)) || (power.gte(TWO_GOOD)):
				return t_quiero(revealCard, "no cards on table (or mine only), low agg, one good card & too medium as mano (or two good) or better, good to accept")
			default:
				return t_noquiero(revealCard, "no cards on table (or mine only), low agg, not even two good cards, should not accept")
			}
		case "normal":
			switch {
			case power.gte(TWO_GOOD):
				return t_raise(revealCard, "no cards on table (or mine only), normal agg, two good cards or better, good to raise")
			case power.gte(THREE_MEDIUM):
				return t_quiero(revealCard, "no cards on table (or mine only), normal agg, three medium cards or better, good to accept")
			default:
				return t_noquiero(revealCard, "no cards on table (or mine only), normal agg, not even three medium cards, should not accept")
			}
		case "high":
			switch {
			case power.gte(ONE_GOOD_ONE_MEDIUM):
				return t_raise(revealCard, "no cards on table (or mine only), high agg, one good card and one medium or better, good to raise")
			case power.gte(TWO_MEDIUM):
				return t_quiero(revealCard, "no cards on table (or mine only), high agg, two medium cards or better, good to accept")
			default:
				return t_noquiero(revealCard, "no cards on table (or mine only), high agg, not even two medium cards, should not accept")
			}
		}
		panic("unreachable")
	case [2]int{0, 1}: // Only their card on table
		canBeat := canBeatCard(gs.TheirRevealedCards[0], gs.YourUnrevealedCards)
		canTie := canTieCard(gs.TheirRevealedCards[0], gs.YourUnrevealedCards)
		switch {
		case canBeat:
			revealCard := lowestCardThatBeats(gs.TheirRevealedCards[0], gs.YourUnrevealedCards)
			power := newPower(cardsWithout(gs.YourUnrevealedCards, revealCard))
			switch agg {
			case "low":
				if power.gte(ONE_GOOD) {
					return t_raise(revealCard, "only their card on table, can beat it, low agg, one good card or better, good to raise")
				}
				if power.lte(TWO_BAD) {
					return t_noquiero(revealCard, "only their card on table, can beat it, low agg, two bad cards, should not accept")
				}
				return t_quiero(revealCard, "only their card on table, can beat it, low agg, at least one medium card, should accept")
			case "normal":
				if power.gte(TWO_MEDIUM) {
					return t_raise(revealCard, "only their card on table, can beat it, normal agg, two medium cards or better, good to raise")
				}
				if power.lte(TWO_BAD) {
					return t_noquiero(revealCard, "only their card on table, can beat it, normal agg, two bad cards, should not accept")
				}
				return t_quiero(revealCard, "only their card on table, can beat it, normal agg, at least one medium card, should accept")
			case "high":
				if power.gte(TWO_MEDIUM) {
					return t_raise(revealCard, "only their card on table, can beat it, high agg, two medium cards or better, good to raise")
				}
				return t_quiero(revealCard, "only their card on table, can beat it, high agg, less than two medium cards, should accept")
			}
			panic("unreachable")
		case canTie:
			revealCard := cardThatTies(gs.TheirRevealedCards[0], gs.YourUnrevealedCards)
			power := newPower(cardsWithout(gs.YourUnrevealedCards, revealCard))
			switch agg {
			case "low":
				if power.gte(ONE_SPECIAL) {
					return t_raise(revealCard, "only their card on table, can tie it, low agg, one special card or better, good to raise")
				}
				if power.lt(ONE_GOOD) {
					return t_noquiero(revealCard, "only their card on table, can tie it, low agg, less than one good card, should not accept")
				}
				return t_quiero(revealCard, "only their card on table, can tie it, low agg, at least one good card, should accept")

			case "normal":
				if power.gte(TWO_GOOD) {
					return t_raise(revealCard, "only their card on table, can tie it, normal agg, two good cards or better, good to raise")
				}
				if power.lte(TWO_MEDIUM) {
					return t_noquiero(revealCard, "only their card on table, can tie it, normal agg, up to two medium cards, should not accept")
				}
				return t_quiero(revealCard, "only their card on table, can tie it, normal agg, at least one good card, should accept")
			case "high":
				if power.gte(TWO_MEDIUM) {
					return t_raise(revealCard, "only their card on table, can tie it, high agg, two medium cards or better, good to raise")
				}
				if power.gte(ONE_MEDIUM) {
					return t_quiero(revealCard, "only their card on table, can tie it, high agg, at least one medium card, should accept")
				}
				return t_noquiero(revealCard, "only their card on table, can tie it, high agg, less than one medium card, should not accept")
			}
		default: // will lose faceoff
			revealCard := lowestOf(gs.YourUnrevealedCards)
			power := newPower(cardsWithout(gs.YourUnrevealedCards, revealCard))
			switch agg {
			case "low":
				if power.gte(ONE_SPECIAL_ONE_GOOD) {
					return t_raise(revealCard, "only their card on table, will lose faceoff, low agg, one special and one good card or better, good to raise")
				}
				if power.lte(TWO_GOOD) {
					return t_noquiero(revealCard, "only their card on table, will lose faceoff, low agg, up to two good cards, should not accept")
				}
				return t_quiero(revealCard, "only their card on table, will lose faceoff, low agg, at least one medium card, should accept")

			case "normal":
				if power.gte(TWO_GOOD) {
					return t_raise(revealCard, "only their card on table, will lose faceoff, normal agg, two good cards or better, good to raise")
				}
				if power.lte(ONE_GOOD_ONE_MEDIUM) {
					return t_noquiero(revealCard, "only their card on table, will lose faceoff, normal agg, up to one good and one medium card, should not accept")
				}
				return t_quiero(revealCard, "only their card on table, will lose faceoff, normal agg, at least one good card, should accept")
			case "high":
				if power.gte(ONE_GOOD_ONE_MEDIUM) {
					return t_raise(revealCard, "only their card on table, will lose faceoff, high agg, one good and one medium card or better, good to raise")
				}
				if power.gte(ONE_GOOD) {
					return t_quiero(revealCard, "only their card on table, will lose faceoff, high agg, at least one good card, should accept")
				}
				return t_noquiero(revealCard, "only their card on table, will lose faceoff, high agg, less than one good card, should not accept")
			}
		}
		panic("unreachable")
	case [2]int{1, 1}, [2]int{2, 1}: // One faceoff done
		var revealCard truco.Card
		if len(gs.YourUnrevealedCards) == 1 {
			revealCard = gs.YourRevealedCards[1]
		}
		switch faceoffResults[0] {
		case FACEOFF_WIN:
			if revealCard == (truco.Card{}) {
				revealCard = lowestOf(gs.YourUnrevealedCards)
			}
			power := newPower(cardsWithout(gs.YourUnrevealedCards, revealCard))
			switch agg {
			case "low":
				if power.gte(TWO_MEDIUM) {
					return t_raise(revealCard, "won faceoff, low agg, two medium cards or better, good to raise")
				}
				if power.lt(ONE_MEDIUM) {
					return t_noquiero(revealCard, "won faceoff, low agg, less than one medium card, should not accept")
				}
				return t_quiero(revealCard, "won faceoff, low agg, one medium card, should accept")

			case "normal":
				if power.gte(ONE_MEDIUM) {
					return t_raise(revealCard, "won faceoff, normal agg, one medium card or better, good to raise")
				}
				if power.lte(TWO_BAD) {
					return t_noquiero(revealCard, "won faceoff, normal agg, two bad cards, should not accept")
				}
				return t_quiero(revealCard, "won faceoff, normal agg, at least one good card, should accept")
			case "high":
				if power.gte(ONE_MEDIUM) {
					return t_raise(revealCard, "won faceoff, high agg, one medium card or better, good to raise")
				}
				return t_quiero(revealCard, "won faceoff, high agg, less than one medium card, should accept")
			}
		case FACEOFF_TIE:
			if revealCard == (truco.Card{}) {
				revealCard = highestOf(gs.YourUnrevealedCards)
			}
			power := newPower(cardsWithout(gs.YourUnrevealedCards, revealCard))
			switch agg {
			case "low":
				if power.gte(ONE_SPECIAL) {
					return t_raise(revealCard, "tied faceoff, low agg, one special card or better, good to raise")
				}
				if power.lt(ONE_GOOD) {
					return t_noquiero(revealCard, "tied faceoff, low agg, less than one good card, should not accept")
				}
				return t_quiero(revealCard, "tied faceoff, low agg, at least one good card, should accept")

			case "normal":
				if power.gte(ONE_GOOD_ONE_MEDIUM) {
					return t_raise(revealCard, "tied faceoff, normal agg, one good and one medium card or better, good to raise")
				}
				if power.lte(TWO_MEDIUM) {
					return t_noquiero(revealCard, "tied faceoff, normal agg, up to two medium cards, should not accept")
				}
				return t_quiero(revealCard, "tied faceoff, normal agg, at least one good card, should accept")
			case "high":
				if power.gte(ONE_MEDIUM) {
					return t_raise(revealCard, "tied faceoff, high agg, one medium card or better, good to raise")
				}
				return t_quiero(revealCard, "tied faceoff, high agg, less than one medium card, should accept")
			}
			panic("unreachable")
		case FACEOFF_LOSS:
			if revealCard == (truco.Card{}) {
				revealCard = highestOf(gs.YourUnrevealedCards)
			}
			power := newPower(cardsWithout(gs.YourUnrevealedCards, revealCard))
			switch agg {
			case "low":
				if power.gte(TWO_SPECIALS) {
					return t_raise(revealCard, "lost faceoff, low agg, two special cards or better, good to raise")
				}
				if power.lt(TWO_GOOD) {
					return t_noquiero(revealCard, "lost faceoff, low agg, less than two good cards, should not accept")
				}
				return t_quiero(revealCard, "lost faceoff, low agg, at least two good cards, should accept")

			case "normal":
				if power.gte(ONE_SPECIAL_ONE_GOOD) {
					return t_raise(revealCard, "lost faceoff, normal agg, one special and one good card or better, good to raise")
				}
				if power.lte(ONE_GOOD_ONE_MEDIUM) {
					return t_noquiero(revealCard, "lost faceoff, normal agg, up to one good and one medium card, should not accept")
				}
				return t_quiero(revealCard, "lost faceoff, normal agg, two good cards, should accept")
			case "high":
				if power.gte(ONE_GOOD) {
					return t_raise(revealCard, "lost faceoff, high agg, one good card or better, good to raise")
				}
				if power.lte(TWO_MEDIUM) {
					return t_noquiero(revealCard, "lost faceoff, high agg, up to two medium cards, should not accept")
				}
				return t_quiero(revealCard, "lost faceoff, high agg, at least one medium card, should accept")
			}
		}
		panic("unreachable")
	case [2]int{1, 2}: // One faceoff done, their card on the table
		canBeat := canBeatCard(gs.TheirRevealedCards[1], gs.YourUnrevealedCards)
		canTie := canTieCard(gs.TheirRevealedCards[1], gs.YourUnrevealedCards)
		switch faceoffResults[0] {
		case FACEOFF_WIN:
			switch {
			case canBeat:
				revealCard := lowestCardThatBeats(gs.TheirRevealedCards[1], gs.YourUnrevealedCards)
				return t_raise(revealCard, "won faceoff, their card on table, can beat it, good to raise")
			case canTie:
				revealCard := cardThatTies(gs.TheirRevealedCards[1], gs.YourUnrevealedCards)
				return t_raise(revealCard, "won faceoff, their card on table, can tie it, good to raise")
			default: // can't beat their card
				revealCard := lowestOf(gs.YourUnrevealedCards)
				power := newPower(cardsWithout(gs.YourUnrevealedCards, revealCard))
				switch agg {
				case "low":
					if power.gte(ONE_GOOD) {
						return t_raise(revealCard, "won faceoff, their card on table, can't beat it, low agg, one good card or better, good to raise")
					}
					return t_noquiero(revealCard, "won faceoff, their card on table, can't beat it, low agg, less than one good card, should not accept")
				case "normal":
					if power.gte(ONE_GOOD) {
						return t_raise(revealCard, "won faceoff, their card on table, can't beat it, normal agg, one good card or better, good to raise")
					}
					if power.lte(ONE_BAD) {
						return t_noquiero(revealCard, "won faceoff, their card on table, can't beat it, normal agg, up to one bad card, should not accept")
					}
					return t_quiero(revealCard, "won faceoff, their card on table, can't beat it, normal agg, at least one medium card, should accept")
				case "high":
					return t_raise(revealCard, "won faceoff, their card on table, can't beat it, high agg, good to raise")
				}
				panic("unreachable")
			}
		case FACEOFF_TIE:
			switch {
			case canBeat:
				revealCard := lowestCardThatBeats(gs.TheirRevealedCards[1], gs.YourUnrevealedCards)
				return t_raise(revealCard, "tied faceoff, their card on table, can beat it, good to raise")
			case canTie:
				revealCard := cardThatTies(gs.TheirRevealedCards[1], gs.YourUnrevealedCards)
				return t_raise(revealCard, "tied faceoff, their card on table, can tie it, good to raise")
			default: // can't beat their card, so about to lose hand
				revealCard := lowestOf(gs.YourUnrevealedCards)
				return t_leave(revealCard, "tied faceoff, their card on table, can't beat it, should leave")
			}
		case FACEOFF_LOSS:
			switch {
			case canBeat:
				revealCard := lowestCardThatBeats(gs.TheirRevealedCards[1], gs.YourUnrevealedCards)
				power := newPower(cardsWithout(gs.YourUnrevealedCards, revealCard))
				switch agg {
				case "low":
					if power.gte(ONE_SPECIAL_ONE_GOOD) {
						return t_raise(revealCard, "lost faceoff, their card on table, can beat it, low agg, one special and one good card or better, good to raise")
					}
					if power.lt(TWO_GOOD) {
						return t_noquiero(revealCard, "lost faceoff, their card on table, can beat it, low agg, less than two good cards, should not accept")
					}
					return t_quiero(revealCard, "lost faceoff, their card on table, can beat it, low agg, at least two good cards, should accept")
				case "normal":
					if power.gte(TWO_GOOD) {
						return t_raise(revealCard, "lost faceoff, their card on table, can beat it, normal agg, two good cards or better, good to raise")
					}
					if power.lt(ONE_GOOD_ONE_MEDIUM) {
						return t_noquiero(revealCard, "lost faceoff, their card on table, can beat it, normal agg, less than one good and one medium card, should not accept")
					}
					return t_quiero(revealCard, "lost faceoff, their card on table, can beat it, normal agg, at least one good card and one medium card, should accept")
				case "high":
					if power.gte(ONE_GOOD_ONE_MEDIUM) {
						return t_raise(revealCard, "lost faceoff, their card on table, can beat it, high agg, one good and one medium card or better, good to raise")
					}
					if power.lt(ONE_GOOD) {
						return t_noquiero(revealCard, "lost faceoff, their card on table, can beat it, high agg, less than one good card, should not accept")
					}
					return t_quiero(revealCard, "lost faceoff, their card on table, can beat it, high agg, at least one good card, should accept")
				}
			default: // tie or lose, either way it's a loss
				revealCard := lowestOf(gs.YourUnrevealedCards)
				return t_leave(revealCard, "lost faceoff, their card on table, can't beat it, should leave")
			}
		}
		panic("unreachable")
	case [2]int{2, 2}, [2]int{3, 2}: // Two faceoffs done
		// TODO: at this point envido or flor info could tell bot what human has
		var revealCard truco.Card
		if len(gs.YourUnrevealedCards) == 0 {
			revealCard = gs.YourRevealedCards[2]
		} else {
			revealCard = gs.YourUnrevealedCards[0]
		}
		power := newPower([]truco.Card{revealCard})
		switch agg {
		case "low":
			if power.gte(ONE_GOOD) {
				return t_raise(revealCard, "two faceoffs done, low agg, one good card or better, good to raise")
			}
			if power.lt(ONE_MEDIUM) {
				return t_noquiero(revealCard, "two faceoffs done, low agg, less than one medium card, should not accept")
			}
			return t_quiero(revealCard, "two faceoffs done, low agg, at least one medium card, should accept")
		case "normal":
			if power.gte(ONE_MEDIUM) {
				return t_raise(revealCard, "two faceoffs done, normal agg, one medium card or better, good to raise")
			}
			if power.lte(ONE_BAD) {
				return t_noquiero(revealCard, "two faceoffs done, normal agg, one bad card, should not accept")
			}
		case "high":
			return t_quiero(revealCard, "two faceoffs done, high agg, should accept")
		}
		panic("unreachable")
	case [2]int{2, 3}: // Two faceoffs done, their card on the table
		canBeat := canBeatCard(gs.TheirRevealedCards[2], gs.YourUnrevealedCards)
		canTie := canTieCard(gs.TheirRevealedCards[2], gs.YourUnrevealedCards)
		switch {
		case canBeat, canTie && youMano(gs): // sure win
			return t_raise(gs.YourUnrevealedCards[0], "two faceoffs done, their card on table, can beat it (or tie as mano), good to raise")
		default: // sure loss
			return t_leave(gs.YourUnrevealedCards[0], "two faceoffs done, their card on table, can't beat it, should leave")
		}
	default:
		panic(fmt.Sprintf("Unexpected number of revealed card pairs: %v", revealedCardPairs))
	}
}

const (
	POWER_SPECIAL = 4
	POWER_GOOD    = 3
	POWER_MEDIUM  = 2
	POWER_BAD     = 1
)

var (
	TWO_SPECIALS         = newPower([]truco.Card{{Suit: truco.ESPADA, Number: 1}, {Suit: truco.BASTO, Number: 1}}, "two_specials cards")
	ONE_SPECIAL_ONE_GOOD = newPower([]truco.Card{{Suit: truco.ESPADA, Number: 1}, {Suit: truco.ORO, Number: 3}}, "one special card, one good card")
	TWO_GOOD             = newPower([]truco.Card{{Suit: truco.ORO, Number: 3}, {Suit: truco.COPA, Number: 3}}, "two good cards")
	TWO_GOOD_ONE_MEDIUM  = newPower([]truco.Card{{Suit: truco.ORO, Number: 3}, {Suit: truco.COPA, Number: 3}, {Suit: truco.BASTO, Number: 11}}, "two good cards, one medium card")
	ONE_GOOD_TWO_MEDIUM  = newPower([]truco.Card{{Suit: truco.ORO, Number: 3}, {Suit: truco.BASTO, Number: 11}, {Suit: truco.COPA, Number: 11}}, "one good card, two medium cards")
	ONE_GOOD_ONE_MEDIUM  = newPower([]truco.Card{{Suit: truco.ORO, Number: 3}, {Suit: truco.BASTO, Number: 11}}, "one good card, one medium card")
	THREE_MEDIUM         = newPower([]truco.Card{{Suit: truco.BASTO, Number: 11}, {Suit: truco.COPA, Number: 11}, {Suit: truco.ORO, Number: 11}}, "three medium cards")
	TWO_MEDIUM           = newPower([]truco.Card{{Suit: truco.BASTO, Number: 11}, {Suit: truco.COPA, Number: 11}}, "two medium cards")
	TWO_BAD              = newPower([]truco.Card{{Suit: truco.ORO, Number: 4}, {Suit: truco.COPA, Number: 4}}, "two bad cards")
	ONE_GOOD             = newPower([]truco.Card{{Suit: truco.ORO, Number: 3}}, "one good card")
	ONE_SPECIAL          = newPower([]truco.Card{{Suit: truco.ESPADA, Number: 1}}, "one special card")
	ONE_MEDIUM           = newPower([]truco.Card{{Suit: truco.BASTO, Number: 11}}, "one medium card")
	ONE_BAD              = newPower([]truco.Card{{Suit: truco.ORO, Number: 4}}, "one bad card")
)

type trucoPower struct {
	power       map[int]int
	totalPower  int
	count       int
	description string
}

func newPower(cards []truco.Card, description ...string) trucoPower {
	return trucoPower{
		power:       cardsToPowers(cards),
		totalPower:  _sumPowers(cardsToPowers(cards)),
		count:       len(cards),
		description: strings.Join(description, " "),
	}
}

func (tp trucoPower) _cmp(tp2 trucoPower) int {
	// If the count is a mismatch, fill the smaller one with BADs (i.e. +1)
	fixTP, fixTP2 := 0, 0
	if tp.count > tp2.count {
		fixTP2 = tp.count - tp2.count
	} else if tp.count < tp2.count {
		fixTP = tp2.count - tp.count
	}

	if tp.totalPower+fixTP > tp2.totalPower+fixTP2 {
		return 1
	}
	if tp.totalPower+fixTP < tp2.totalPower+fixTP2 {
		return -1
	}
	return 0
}

func (tp trucoPower) lt(tp2 trucoPower) bool {
	return tp._cmp(tp2) == -1
}

func (tp trucoPower) gte(tp2 trucoPower) bool {
	return tp._cmp(tp2) >= 0
}

func (tp trucoPower) lte(tp2 trucoPower) bool {
	return tp._cmp(tp2) <= 0
}

func _sumPowers(powers map[int]int) int {
	sum := 0
	for power, count := range powers {
		sum += power * count
	}
	return sum
}

func cardsToPowers(cards []truco.Card) map[int]int {
	powers := map[int]int{}
	for _, c := range cards {
		powers[cardToPower(c)]++
	}
	return powers
}

func cardToPower(c truco.Card) int {
	if isCardSpecial(c) {
		return POWER_SPECIAL
	}
	if c.Number == 3 || c.Number == 2 {
		return POWER_GOOD
	}
	if c.Number == 1 || c.Number == 12 || c.Number == 11 {
		return POWER_MEDIUM
	}
	return POWER_BAD
}

// func calculateCardStrength(gs truco.Card) int {
// 	specialValues := map[truco.Card]int{
// 		{Suit: truco.ESPADA, Number: 1}: 15,
// 		{Suit: truco.BASTO, Number: 1}:  14,
// 		{Suit: truco.ESPADA, Number: 7}: 13,
// 		{Suit: truco.ORO, Number: 7}:    12,
// 	}
// 	if _, ok := specialValues[gs]; ok {
// 		return specialValues[gs]
// 	}
// 	if gs.Number <= 3 {
// 		return gs.Number + 12 - 4
// 	}
// 	return gs.Number - 4
// }

// func sortCardsByValue(cards []truco.Card) []truco.Card {
// 	sort.Slice(cards, func(i, j int) bool {
// 		return calculateCardStrength(cards[i]) > calculateCardStrength(cards[j])
// 	})
// 	return cards
// }

func isCardSpecial(card truco.Card) bool {
	return card == (truco.Card{Suit: truco.ESPADA, Number: 1}) || card == (truco.Card{Suit: truco.BASTO, Number: 1}) || card == (truco.Card{Suit: truco.ESPADA, Number: 7}) || card == (truco.Card{Suit: truco.ORO, Number: 7})
}

func forAllCards(cards []truco.Card, f func(truco.Card) bool) bool {
	for _, c := range cards {
		if !f(c) {
			return false
		}
	}
	return true
}

func randomCardExceptSpecial(cards []truco.Card) truco.Card {
	if forAllCards(cards, isCardSpecial) {
		return cards[rand.Intn(len(cards))]
	}
	for {
		card := cards[rand.Intn(len(cards))]
		if !isCardSpecial(card) {
			return card
		}
	}
}

// func randomCard(cards []truco.Card) truco.Card {
// 	return cards[rand.Intn(len(cards))]
// }

func canBeatCard(card truco.Card, cards []truco.Card) bool {
	for _, c := range cards {
		if c.CompareTrucoScore(card) == 1 {
			return true
		}
	}
	return false
}

func canTieCard(card truco.Card, cards []truco.Card) bool {
	for _, c := range cards {
		if c.CompareTrucoScore(card) == 0 {
			return true
		}
	}
	return false
}

// func cardsWithoutLowest(cards []truco.Card) []truco.Card {
// 	lowest := cards[0]
// 	for _, card := range cards {
// 		if card.CompareTrucoScore(lowest) == -1 {
// 			lowest = card
// 		}
// 	}

// 	unrevealed := []truco.Card{}
// 	for _, card := range cards {
// 		if card != lowest {
// 			unrevealed = append(unrevealed, card)
// 		}
// 	}
// 	return unrevealed
// }

func lowestOf(cards []truco.Card) truco.Card {
	lowest := cards[0]
	for _, card := range cards {
		if card.CompareTrucoScore(lowest) == -1 {
			lowest = card
		}
	}
	return lowest
}

func highestOf(cards []truco.Card) truco.Card {
	highest := cards[0]
	for _, card := range cards {
		if card.CompareTrucoScore(highest) == 1 {
			highest = card
		}
	}
	return highest
}

func cardsWithout(cards []truco.Card, without truco.Card) []truco.Card {
	filtered := []truco.Card{}
	for _, card := range cards {
		if card != without {
			filtered = append(filtered, card)
		}
	}
	return filtered
}

// func cardsWithoutLowestCardThatBeats(card truco.Card, cards []truco.Card) []truco.Card {
// 	return cardsWithout(cards, lowestCardThatBeats(card, cards))
// }

// func cardsWithoutCardThatTies(card truco.Card, cards []truco.Card) []truco.Card {
// 	return cardsWithout(cards, cardThatTies(card, cards))
// }

func cardThatTies(card truco.Card, cards []truco.Card) truco.Card {
	for _, c := range cards {
		if c.CompareTrucoScore(card) == 0 {
			return c
		}
	}
	return truco.Card{} // This should be unreachable
}

func lowestCardThatBeats(card truco.Card, cards []truco.Card) truco.Card {
	cardsThatBeatCard := []truco.Card{}
	for _, c := range cards {
		if c.CompareTrucoScore(card) == 1 {
			cardsThatBeatCard = append(cardsThatBeatCard, c)
		}
	}
	if len(cardsThatBeatCard) == 0 {
		return truco.Card{}
	}
	return lowestOf(cardsThatBeatCard)
}
