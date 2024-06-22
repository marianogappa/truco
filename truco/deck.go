package truco

import (
	"errors"
	"fmt"
	"math/rand"
)

const (
	ORO    = "oro"
	COPA   = "copa"
	ESPADA = "espada"
	BASTO  = "basto"
)

type Card struct {
	Suit   string `json:"suit"`
	Number int    `json:"number"`
}

func (c Card) String() string {
	return fmt.Sprintf("%d de %s", c.Number, c.Suit)
}

type Deck struct {
	cards []Card
}

type Hand struct {
	Unrevealed []Card `json:"unrevealed"`
	Revealed   []Card `json:"revealed"`
}

func (h Hand) HasUnrevealedCard(c Card) bool {
	for _, card := range h.Unrevealed {
		if card == c {
			return true
		}
	}
	return false
}

func (h *Hand) RevealCard(card Card) error {
	for _, c := range h.Revealed {
		if c == card {
			return errCardAlreadyRevealed
		}
	}
	for i, c := range h.Unrevealed {
		if c == card {
			h.Revealed = append(h.Revealed, c)
			h.Unrevealed = append(h.Unrevealed[:i], h.Unrevealed[i+1:]...)
			return nil
		}
	}
	return errCardNotInHand
}

var (
	errCardNotInHand       = errors.New("card not in hand")
	errCardAlreadyRevealed = errors.New("card already revealed")
)

func makeSpanishCards() []Card {
	cards := []Card{}
	suits := []string{ORO, COPA, ESPADA, BASTO}
	for _, suit := range suits {
		for i := 1; i <= 12; i++ {
			if i == 8 || i == 9 {
				continue
			}
			cards = append(cards, Card{Suit: suit, Number: i})
		}
	}

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

func NewDeck() *Deck {
	d := Deck{cards: makeSpanishCards()}
	return &d
}

func (d *Deck) DealHand() *Hand {
	hand := &Hand{}
	for i := 0; i < 3; i++ {
		hand.Unrevealed = append(hand.Unrevealed, d.cards[i])
	}
	d.cards = d.cards[3:]
	return hand
}

// EnvidoScore returns the score of the hand according to the Envido rules.
//
// The score is an integer between 0 and 33.
//
// For the purpose of calculating the score, all 3 cards are considered, regardless of being revealed.
// Therefore, use all cards from the hand, not just the revealed ones.
//
// The score is calculated as follows:
// - The value of each cards is its number, but 10, 11 and 12 are worth 0.
// - If all cards are different suit, the score is the highest value card's value.
// - If two cards are the same suit, the score is the sum of the two cards' values plus 20.
// - If all cards are the same suit, the score is the maximum score of the three possible pairs of cards.
func (h Hand) EnvidoScore() int {
	var (
		maxScore     = 0
		suitToValues = make(map[string][]int)
	)
	for _, card := range append(h.Unrevealed, h.Revealed...) {
		suitToValues[card.Suit] = append(suitToValues[card.Suit], card.Number)
		if card.Number >= 10 {
			suitToValues[card.Suit][len(suitToValues[card.Suit])-1] = 0
		}
	}
	for suit, values := range suitToValues {
		switch len(values) {
		case 3:
			maxScore = max(
				maxScore,
				suitToValues[suit][0]+suitToValues[suit][1]+20,
				suitToValues[suit][0]+suitToValues[suit][2]+20,
				suitToValues[suit][1]+suitToValues[suit][2]+20,
			)
		case 2:
			maxScore = max(maxScore, suitToValues[suit][0]+suitToValues[suit][1]+20)
		case 1:
			maxScore = max(maxScore, suitToValues[suit][0])
		}
	}
	return maxScore
}

// CompareTrucoScore returns:
// -  1 if the receiver card has a higher Truco score than the other card
// - -1 if it has a lower score
// -  0 if they have the same score
//
// So essentially it's a sort order for the cards, according to Truco rules.
//
// The sort is calculated as follows:
//
// - There are 4 special cards, which are the highest in this order:
//
//  1. 1 of Espada
//
//  2. 1 of Basto
//
//  3. 7 of Espada
//
//  4. 7 of Oro
//
//     - The rest of the cards are ordered by their number, but the numerical order is:
//     3, 2, 1, 12, 11, 10, 7, 6, 5, 4
func (c Card) CompareTrucoScore(other Card) int {
	specialValues := map[Card]int{
		{Suit: ESPADA, Number: 1}: 19,
		{Suit: BASTO, Number: 1}:  18,
		{Suit: ESPADA, Number: 7}: 17,
		{Suit: ORO, Number: 7}:    16,
	}
	a := c.Number
	b := other.Number
	if specialValue, ok := specialValues[c]; ok {
		a = specialValue
	}
	if specialValue, ok := specialValues[other]; ok {
		b = specialValue
	}
	if a == b {
		return 0
	}
	if a <= 3 {
		a += 12
	}
	if b <= 3 {
		b += 12
	}
	return sign(a - b)
}

func sign(i int) int {
	if i < 0 {
		return -1
	}
	return 1
}
