package truco

import (
	"fmt"
	"testing"
)

func TestCard_CompareTrucoScore(t *testing.T) {
	var (
		card_espada_1  = Card{Suit: ESPADA, Number: 1}
		card_espada_7  = Card{Suit: ESPADA, Number: 7}
		card_oro_6     = Card{Suit: ORO, Number: 6}
		card_copa_1    = Card{Suit: COPA, Number: 1}
		card_basto_7   = Card{Suit: BASTO, Number: 7}
		card_espada_4  = Card{Suit: ESPADA, Number: 4}
		card_copa_3    = Card{Suit: COPA, Number: 3}
		card_basto_2   = Card{Suit: BASTO, Number: 2}
		card_espada_5  = Card{Suit: ESPADA, Number: 5}
		card_copa_4    = Card{Suit: COPA, Number: 4}
		card_basto_3   = Card{Suit: BASTO, Number: 3}
		card_espada_6  = Card{Suit: ESPADA, Number: 6}
		card_copa_5    = Card{Suit: COPA, Number: 5}
		card_basto_4   = Card{Suit: BASTO, Number: 4}
		card_copa_6    = Card{Suit: COPA, Number: 6}
		card_basto_5   = Card{Suit: BASTO, Number: 5}
		card_espada_10 = Card{Suit: ESPADA, Number: 10}
		card_copa_7    = Card{Suit: COPA, Number: 7}
		card_basto_6   = Card{Suit: BASTO, Number: 6}
		card_espada_11 = Card{Suit: ESPADA, Number: 11}
		card_oro_1     = Card{Suit: ORO, Number: 1}
	)
	// Define test cases
	tests := []struct {
		card1    Card
		card2    Card
		expected int
	}{
		{card_espada_1, card_espada_7, 1},
		{card_espada_7, card_espada_1, -1},
		{card_oro_6, card_copa_1, -1},
		{card_basto_7, card_espada_4, 1},
		{card_copa_3, card_basto_2, 1},
		{card_copa_4, card_espada_5, -1},
		{card_basto_3, card_espada_6, 1},
		{card_basto_4, card_copa_5, -1},
		{card_copa_6, card_basto_5, 1},
		{card_espada_10, card_copa_7, 1},
		{card_espada_11, card_basto_6, 1},
		{card_espada_1, card_copa_1, 1},
		{card_copa_1, card_espada_1, -1},
		{card_copa_1, card_oro_1, 0},
		{card_espada_7, card_basto_7, 1},
		{card_basto_7, card_espada_7, -1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v vs %v", tt.card1, tt.card2), func(t *testing.T) {
			result := tt.card1.CompareTrucoScore(tt.card2)
			if result != tt.expected {
				t.Errorf("Expected %d for cards %v and %v, got %d", tt.expected, tt.card1, tt.card2, result)
			}
		})
	}
}
