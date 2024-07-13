package examplebot

import (
	"math"
	"testing"

	"github.com/marianogappa/truco/truco"
)

func TestCalculateCardStrength(t *testing.T) {
	testCases := []struct {
		card     truco.Card
		expected int
	}{
		{card: truco.Card{Suit: truco.ESPADA, Number: 1}, expected: 19},
		{card: truco.Card{Suit: truco.ESPADA, Number: 2}, expected: 14},
		{card: truco.Card{Suit: truco.ESPADA, Number: 3}, expected: 15},
		{card: truco.Card{Suit: truco.ESPADA, Number: 4}, expected: 4},
	}

	for _, tc := range testCases {
		actual := calculateCardStrength(tc.card)
		if actual != tc.expected {
			t.Errorf("calculateCardStrength(%v) = %d, expected %d", tc.card, actual, tc.expected)
		}
	}
}

func TestCalculateTrucoHandChance(t *testing.T) {
	testCases := []struct {
		cards    []truco.Card
		expected float64
	}{
		{
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 7},
				{Suit: truco.BASTO, Number: 1},
			},
			expected: 1,
		},
		{
			cards: []truco.Card{
				{Suit: truco.BASTO, Number: 4},
				{Suit: truco.COPA, Number: 4},
				{Suit: truco.ORO, Number: 4},
			},
			expected: 0.0,
		},
		// Add more test cases here...
	}

	for _, tc := range testCases {
		actual := calculateTrucoHandChance(tc.cards)
		if !floatEquals(actual, tc.expected, 0.01) {
			t.Errorf("calculateTrucoHandChance(%v) = %f, expected %f", tc.cards, actual, tc.expected)
		}

	}
}

// Function to compare floating-point numbers within a tolerance.
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestCanBeatCard(t *testing.T) {
	testCases := []struct {
		card     truco.Card
		cards    []truco.Card
		expected bool
	}{
		// Test case 1: card can beat the highest card in the list
		{
			card: truco.Card{Suit: truco.ESPADA, Number: 7},
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 4},
				{Suit: truco.ORO, Number: 7},
			},
			expected: true,
		},
		// Test case 2: card cannot beat any card in the list
		{
			card: truco.Card{Suit: truco.BASTO, Number: 3},
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.COPA, Number: 4},
				{Suit: truco.BASTO, Number: 7},
			},
			expected: false,
		},
		// Test case 3: card can beat some cards in the list
		{
			card: truco.Card{Suit: truco.ORO, Number: 3},
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.COPA, Number: 4},
				{Suit: truco.ORO, Number: 7},
			},
			expected: true,
		},
		// Add more test cases here...
	}
	for _, tc := range testCases {
		actual := canBeatCard(tc.card, tc.cards)
		if actual != tc.expected {
			t.Errorf("canBeatCard(%v, %v) = %t, expected %t", tc.card, tc.cards, actual, tc.expected)
		}
	}
}

func TestCanTieCard(t *testing.T) {
	testCases := []struct {
		card     truco.Card
		cards    []truco.Card
		expected bool
	}{
		// Test case 1: card can tie the highest card in the list
		{
			card: truco.Card{Suit: truco.ESPADA, Number: 7},
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 4},
				{Suit: truco.ORO, Number: 7},
			},
			expected: false,
		},
		// Test case 2: card cannot tie any card in the list
		{
			card: truco.Card{Suit: truco.BASTO, Number: 3},
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.COPA, Number: 4},
				{Suit: truco.BASTO, Number: 7},
			},
			expected: false,
		},
		// Test case 3: card can tie some cards in the list
		{
			card: truco.Card{Suit: truco.ORO, Number: 3},
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.COPA, Number: 3},
				{Suit: truco.ORO, Number: 7},
			},
			expected: true,
		},
		// Test case 4: card is the highest card in the list
		{
			card: truco.Card{Suit: truco.COPA, Number: 5},
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.BASTO, Number: 5},
				{Suit: truco.COPA, Number: 7},
			},
			expected: true,
		},
		// Test case 5: card is the lowest card in the list
		{
			card: truco.Card{Suit: truco.ORO, Number: 1},
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 4},
				{Suit: truco.COPA, Number: 7},
				{Suit: truco.BASTO, Number: 2},
			},
			expected: false,
		},
	}
	for _, tc := range testCases {
		actual := canTieCard(tc.card, tc.cards)
		if actual != tc.expected {
			t.Errorf("canTieCard(%v, %v) = %t, expected %t", tc.card, tc.cards, actual, tc.expected)
		}
	}
}

func TestCardsWithoutLowest(t *testing.T) {
	testCases := []struct {
		cards    []truco.Card
		expected []truco.Card
	}{
		// Test case 1: Remove lowest card from a list of cards
		{
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 2},
				{Suit: truco.ESPADA, Number: 3},
			},
			expected: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 3},
			},
		},
		// Test case 2: Remove lowest card from a list of cards with duplicates
		{
			cards: []truco.Card{
				{Suit: truco.BASTO, Number: 4},
				{Suit: truco.COPA, Number: 5},
				{Suit: truco.ORO, Number: 6},
			},
			expected: []truco.Card{
				{Suit: truco.COPA, Number: 5},
				{Suit: truco.ORO, Number: 6},
			},
		},
		// Test case 3: Remove lowest card from a list of cards with lowest card as the only card
		{
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
			},
			expected: []truco.Card{},
		},
		// Test case 5: Remove lowest card from a list of cards with multiple lowest cards
		{
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.COPA, Number: 2},
			},
			expected: []truco.Card{
				{Suit: truco.COPA, Number: 2},
			},
		},
	}
	for _, tc := range testCases {
		actual := cardsWithoutLowest(tc.cards)
		if !cardsEqual(actual, tc.expected) {
			t.Errorf("cardsWithoutLowest(%v) = %v, expected %v", tc.cards, actual, tc.expected)
		}
	}
}

// Function to compare two slices of cards for equality.
func cardsEqual(a, b []truco.Card) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestLowestOf(t *testing.T) {
	testCases := []struct {
		cards    []truco.Card
		expected truco.Card
	}{
		// Test case 1: Find the lowest card in a list of cards
		{
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 2},
				{Suit: truco.ESPADA, Number: 3},
			},
			expected: truco.Card{Suit: truco.ESPADA, Number: 2},
		},
		// Test case 2: Find the lowest card in a list of cards with duplicates
		{
			cards: []truco.Card{
				{Suit: truco.BASTO, Number: 4},
				{Suit: truco.COPA, Number: 5},
				{Suit: truco.ORO, Number: 6},
			},
			expected: truco.Card{Suit: truco.BASTO, Number: 4},
		},
		// Test case 3: Find the lowest card in a list of cards with lowest card as the only card
		{
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
			},
			expected: truco.Card{Suit: truco.ORO, Number: 1},
		},
		// Test case 4: Find the lowest card in a list of cards with multiple lowest cards
		{
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.COPA, Number: 2},
			},
			expected: truco.Card{Suit: truco.ORO, Number: 1},
		},
	}
	for _, tc := range testCases {
		actual := lowestOf(tc.cards)
		if actual != tc.expected {
			t.Errorf("lowestOf(%v) = %v, expected %v", tc.cards, actual, tc.expected)
		}
	}
}

func TestHighestOf(t *testing.T) {
	testCases := []struct {
		cards    []truco.Card
		expected truco.Card
	}{
		// Test case 1: Find the highest card in a list of cards
		{
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 2},
				{Suit: truco.ESPADA, Number: 3},
			},
			expected: truco.Card{Suit: truco.ESPADA, Number: 1},
		},
		// Test case 2: Find the highest card in a list of cards with duplicates
		{
			cards: []truco.Card{
				{Suit: truco.BASTO, Number: 4},
				{Suit: truco.COPA, Number: 5},
				{Suit: truco.ORO, Number: 6},
			},
			expected: truco.Card{Suit: truco.ORO, Number: 6},
		},
		// Test case 3: Find the highest card in a list of cards with highest card as the only card
		{
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 7},
			},
			expected: truco.Card{Suit: truco.ORO, Number: 7},
		},
		// Test case 4: Find the highest card in a list of cards with multiple highest cards
		{
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 7},
				{Suit: truco.COPA, Number: 6},
			},
			expected: truco.Card{Suit: truco.ORO, Number: 7},
		},
	}
	for _, tc := range testCases {
		actual := highestOf(tc.cards)
		if actual != tc.expected {
			t.Errorf("highestOf(%v) = %v, expected %v", tc.cards, actual, tc.expected)
		}
	}
}

func TestCardsWithout(t *testing.T) {
	testCases := []struct {
		cards    []truco.Card
		card     truco.Card
		expected []truco.Card
	}{
		// Test case 1: Remove a card from a list of cards
		{
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 2},
				{Suit: truco.ESPADA, Number: 3},
			},
			card: truco.Card{Suit: truco.ESPADA, Number: 2},
			expected: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 3},
			},
		},
		// Test case 2: Remove a card from a list of cards with duplicates
		{
			cards: []truco.Card{
				{Suit: truco.BASTO, Number: 4},
				{Suit: truco.COPA, Number: 5},
				{Suit: truco.ORO, Number: 6},
			},
			card: truco.Card{Suit: truco.COPA, Number: 5},
			expected: []truco.Card{
				{Suit: truco.BASTO, Number: 4},
				{Suit: truco.ORO, Number: 6},
			},
		},
		// Test case 3: Remove a card from a list of cards with the card to remove as the only card
		{
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
			},
			card:     truco.Card{Suit: truco.ORO, Number: 1},
			expected: []truco.Card{},
		},
		// Test case 4: Remove a card from a list of cards with multiple cards to remove
		{
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.COPA, Number: 2},
			},
			card:     truco.Card{Suit: truco.COPA, Number: 2},
			expected: []truco.Card{{Suit: truco.ORO, Number: 1}},
		},
	}
	for _, tc := range testCases {
		actual := cardsWithout(tc.cards, tc.card)
		if !cardsEqual(actual, tc.expected) {
			t.Errorf("cardsWithout(%v, %v) = %v, expected %v", tc.cards, tc.card, actual, tc.expected)
		}
	}
}

func TestCardsWithoutLowestCardThatBeats(t *testing.T) {
	testCases := []struct {
		card     truco.Card
		cards    []truco.Card
		expected []truco.Card
	}{
		// Test case 1: Remove lowest card that beats the given card from a list of cards
		{
			card: truco.Card{Suit: truco.ESPADA, Number: 2},
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.COPA, Number: 2},
				{Suit: truco.ESPADA, Number: 3},
			},
			expected: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.COPA, Number: 2},
			},
		},
		// Test case 2: Remove lowest card that beats the given card from a list of cards with duplicates
		{
			card: truco.Card{Suit: truco.COPA, Number: 5},
			cards: []truco.Card{
				{Suit: truco.BASTO, Number: 4},
				{Suit: truco.ESPADA, Number: 5},
				{Suit: truco.ORO, Number: 6},
			},
			expected: []truco.Card{
				{Suit: truco.BASTO, Number: 4},
				{Suit: truco.ESPADA, Number: 5},
			},
		},
		{
			card: truco.Card{Suit: truco.ORO, Number: 3},
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.ORO, Number: 7},
			},
			expected: []truco.Card{
				{Suit: truco.ORO, Number: 1},
			},
		},
	}
	for _, tc := range testCases {
		actual := cardsWithoutLowestCardThatBeats(tc.card, tc.cards)
		if !cardsEqual(actual, tc.expected) {
			t.Errorf("cardsWithoutLowestCardThatBeats(%v, %v) = %v, expected %v", tc.card, tc.cards, actual, tc.expected)
		}
	}
}

func TestLowestCardThatBeats(t *testing.T) {
	testCases := []struct {
		card     truco.Card
		cards    []truco.Card
		expected truco.Card
	}{
		// Test case 1: Find the lowest card that beats the given card
		{
			card: truco.Card{Suit: truco.ESPADA, Number: 2},
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.COPA, Number: 2},
				{Suit: truco.ESPADA, Number: 3},
			},
			expected: truco.Card{Suit: truco.ESPADA, Number: 3},
		},
		{
			card: truco.Card{Suit: truco.COPA, Number: 5},
			cards: []truco.Card{
				{Suit: truco.BASTO, Number: 4},
				{Suit: truco.ESPADA, Number: 5},
				{Suit: truco.ORO, Number: 6},
			},
			expected: truco.Card{Suit: truco.ORO, Number: 6},
		},
		// Test case 3: Find the lowest card that beats the given card from a list of cards with no card that beats it
		{
			card: truco.Card{Suit: truco.ORO, Number: 3},
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.ORO, Number: 7},
			},
			expected: truco.Card{Suit: truco.ORO, Number: 7},
		},
		// Test case 4: Find the lowest card that beats the given card from a list of cards with multiple cards that beat it
		{
			card: truco.Card{Suit: truco.ORO, Number: 4},
			cards: []truco.Card{
				{Suit: truco.ORO, Number: 1},
				{Suit: truco.COPA, Number: 3},
				{Suit: truco.ORO, Number: 7},
			},
			expected: truco.Card{Suit: truco.ORO, Number: 1},
		},
	}
	for _, tc := range testCases {
		actual := lowestCardThatBeats(tc.card, tc.cards)
		if actual != tc.expected {
			t.Errorf("Expected lowest card that beats %v to be %v, but got %v", tc.card, tc.expected, actual)
		}
	}
}

func TestCardsChance(t *testing.T) {
	testCases := []struct {
		cards    []truco.Card
		expected float64
	}{
		// Test case 1: Calculate chance with no cards
		{
			cards:    []truco.Card{},
			expected: 0.0,
		},
		// Test case 2: Calculate chance with one card
		{
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
			},
			expected: 1.0,
		},
		// Test case 3: Calculate chance with multiple cards
		{
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.BASTO, Number: 1},
			},
			expected: 1.0,
		},
		{
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 1},
				{Suit: truco.ESPADA, Number: 7},
				{Suit: truco.BASTO, Number: 1},
			},
			expected: 1.0,
		},
		{
			cards: []truco.Card{
				{Suit: truco.ESPADA, Number: 4},
			},
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		actual := cardsChance(tc.cards)
		if math.Abs(actual-tc.expected) > 0.01 {
			t.Errorf("cardsChance(%v) = %f; expected %f", tc.cards, actual, tc.expected)
		}
	}
}

func TestChanceOfWinningTruco(t *testing.T) {
	testCases := []struct {
		name     string
		gs       truco.ClientGameState
		expected float64
	}{
		{
			name: "Best hand possible, no info from the opponent: 100% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 1},
					{Suit: truco.BASTO, Number: 1},
					{Suit: truco.ESPADA, Number: 7},
				},
				TheirRevealedCards: []truco.Card{},
			},
			expected: 1.0,
		},
		{
			name: "Worst hand possible, no info from the opponent: 0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
					{Suit: truco.BASTO, Number: 4},
					{Suit: truco.ORO, Number: 4},
				},
				TheirRevealedCards: []truco.Card{},
			},
			expected: 0.0,
		},
		{
			name: "Last card to play, will win the faceoff: 100%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.BASTO, Number: 5}, // first faceoff won
					{Suit: truco.ORO, Number: 4},   // second faceoff lost
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 1}, // Will win the third faceoff
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
					{Suit: truco.BASTO, Number: 6},
					{Suit: truco.ORO, Number: 3},
				},
			},
			expected: 1.0,
		},
		{
			name: "Last card to play, will lost the faceoff: 0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.BASTO, Number: 5}, // first faceoff won
					{Suit: truco.ORO, Number: 4},   // second faceoff lost
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 2}, // Will lose the third faceoff
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
					{Suit: truco.BASTO, Number: 6},
					{Suit: truco.ORO, Number: 3},
				},
			},
			expected: 0.0,
		},
		{
			name: "Beats the opponent's revealed card, and has perfect remaining cards: 100%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 5},
					{Suit: truco.ESPADA, Number: 1},
					{Suit: truco.BASTO, Number: 1},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
				},
			},
			expected: 1.0,
		},
		{
			name: "Beats the opponent's revealed card, and has the worst remaining hand: 0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 5},
					{Suit: truco.ESPADA, Number: 4},
					{Suit: truco.BASTO, Number: 4},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
				},
			},
			expected: 0.0,
		},
		{
			name: "Ties the opponent's revealed card, highest remaining card 3: (3+12-4)/(19-4) = 73.3%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
					{Suit: truco.ESPADA, Number: 4},
					{Suit: truco.BASTO, Number: 3},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
				},
			},
			expected: 0.733,
		},
		{
			name: "Can't beat the opponent's revealed card, highest remaining card 3: (10+7-8)/(19+18-8)*0.66 = 20.4%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 10},
					{Suit: truco.ORO, Number: 4},
					{Suit: truco.BASTO, Number: 7},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 7},
				},
			},
			expected: 0.204,
		},
		{
			name: "Tie at first revealed card, highest possible remaining card = 100.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 1},
					{Suit: truco.ORO, Number: 4},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
				},
			},
			expected: 1.0,
		},
		{
			name: "Tie at first revealed card, lowest possible remaining card = 0.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
					{Suit: truco.ORO, Number: 4},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
				},
			},
			expected: 0.0,
		},
		{
			name: "Tie at second revealed card, highest possible remaining card = 100.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
					{Suit: truco.ORO, Number: 3},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 1},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
					{Suit: truco.BASTO, Number: 3},
				},
			},
			expected: 1.0,
		},
		{
			name: "Tie at second revealed card, lowest possible remaining card = 0.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
					{Suit: truco.ORO, Number: 3},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
					{Suit: truco.BASTO, Number: 3},
				},
			},
			expected: 0.0,
		},
		{
			name: "Win at first revealed card, highest possible remaining card = 100.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 1},
					{Suit: truco.ORO, Number: 4},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 2},
				},
			},
			expected: 1.0,
		},
		{ // TODO this is clearly not correct chance; should be better than if tie first faceoff, but how much?
			name: "Win at first revealed card, lowest possible remaining card = 0.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
					{Suit: truco.ORO, Number: 4},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 2},
				},
			},
			expected: 0.0,
		},
		{
			name: "Tie at first revealed card, can beat their revealed card = 100.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
					{Suit: truco.ORO, Number: 6},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
					{Suit: truco.ORO, Number: 5},
				},
			},
			expected: 1.0,
		},
		{
			name: "Tie at first revealed card, can't beat their revealed card = 0.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 4},
					{Suit: truco.ORO, Number: 6},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
					{Suit: truco.ORO, Number: 7},
				},
			},
			expected: 0.0,
		},
		{
			name: "Tie at first revealed card, can tie their revealed card, remaining card 10: (10-4)/(19-4) = 40.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 3},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 2},
					{Suit: truco.ORO, Number: 10},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
					{Suit: truco.ORO, Number: 2},
				},
			},
			expected: 0.4,
		},
		{
			name: "Loss at first revealed card, can tie their revealed card = 0.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 2},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 2},
					{Suit: truco.ORO, Number: 10},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
					{Suit: truco.ORO, Number: 2},
				},
			},
			expected: 0.0,
		},
		{
			name: "Loss at first revealed card, can't beat their revealed card = 0.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 2},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.COPA, Number: 1},
					{Suit: truco.ORO, Number: 10},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
					{Suit: truco.ORO, Number: 2},
				},
			},
			expected: 0.0,
		},
		{
			name: "Loss at first revealed card, can beat their revealed card = (10-4)/(19-4) = 40.0%% chance of winning",
			gs: truco.ClientGameState{
				YourRevealedCards: []truco.Card{
					{Suit: truco.ESPADA, Number: 2},
				},
				YourUnrevealedCards: []truco.Card{
					{Suit: truco.BASTO, Number: 1},
					{Suit: truco.ORO, Number: 10},
				},
				TheirRevealedCards: []truco.Card{
					{Suit: truco.ORO, Number: 3},
					{Suit: truco.ORO, Number: 2},
				},
			},
			expected: 0.4,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := chanceOfWinningTruco(tc.gs)
			if math.Abs(actual-tc.expected) > 0.01 {
				t.Errorf("Expected %f, but got %f", tc.expected, actual)
			}
		})
	}
}
