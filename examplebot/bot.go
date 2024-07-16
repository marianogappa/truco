package examplebot

import (
	"encoding/json"
	"log"

	"math/rand"

	"github.com/marianogappa/truco/truco"
)

type Bot struct{}

func New() Bot {
	return Bot{}
}

func _deserializeActions(as []json.RawMessage) []truco.Action {
	_as := []truco.Action{}
	for _, a := range as {
		_a, _ := truco.DeserializeAction(a)
		_as = append(_as, _a)
	}
	return _as
}

func possibleActionsMap(gs truco.ClientGameState) map[string]truco.Action {
	possibleActions := make(map[string]truco.Action)
	for _, action := range _deserializeActions(gs.PossibleActions) {
		possibleActions[action.GetName()] = action
	}
	return possibleActions
}

func filter(possibleActions map[string]truco.Action, candidateActions ...truco.Action) []truco.Action {
	filteredActions := []truco.Action{}
	for _, action := range candidateActions {
		if possibleAction, ok := possibleActions[action.GetName()]; ok {
			filteredActions = append(filteredActions, possibleAction)
		}
	}
	return filteredActions
}

func calculateAggresiveness(gs truco.ClientGameState) string {
	aggresiveness := "normal"
	if gs.YourScore-gs.TheirScore >= 5 {
		aggresiveness = "low"
	}
	if gs.YourScore-gs.TheirScore <= -5 {
		aggresiveness = "high"
	}
	return aggresiveness
}

func calculateEnvidoScore(gs truco.ClientGameState) int {
	return truco.Hand{Revealed: gs.YourRevealedCards, Unrevealed: gs.YourUnrevealedCards}.EnvidoScore()
}

func calculateCardStrength(gs truco.Card) int {
	specialValues := map[truco.Card]int{
		{Suit: truco.ESPADA, Number: 1}: 19,
		{Suit: truco.BASTO, Number: 1}:  18,
		{Suit: truco.ESPADA, Number: 7}: 17,
		{Suit: truco.ORO, Number: 7}:    16,
	}
	if _, ok := specialValues[gs]; ok {
		return specialValues[gs]
	}
	if gs.Number <= 3 {
		return gs.Number + 12
	}
	return gs.Number
}

func faceoffResults(gs truco.ClientGameState) []int {
	results := []int{}
	for i := 0; i < min(len(gs.YourRevealedCards), len(gs.TheirRevealedCards)); i++ {
		results = append(results, gs.YourRevealedCards[i].CompareTrucoScore(gs.TheirRevealedCards[i]))
	}
	return results
}

func calculateTrucoHandChance(cards []truco.Card) float64 {
	base := float64(len(cards) * 4)
	sum := -base
	for _, card := range cards {
		sum += float64(calculateCardStrength(card))
	}
	return sum / (19 + 18 + 17 - base)
}

func canAnyEnvido(actions map[string]truco.Action) bool {
	return len(filter(actions,
		truco.NewActionSayEnvido(1),
		truco.NewActionSayRealEnvido(1),
		truco.NewActionSayFaltaEnvido(1),
		truco.NewActionSayEnvidoQuiero(1),
		truco.NewActionSayEnvidoNoQuiero(1),
	)) > 0
}

func possibleEnvidoActionsMap(gs truco.ClientGameState) map[string]truco.Action {
	possible := possibleActionsMap(gs)

	filter := map[string]struct{}{
		truco.SAY_ENVIDO:        {},
		truco.SAY_REAL_ENVIDO:   {},
		truco.SAY_FALTA_ENVIDO:  {},
		truco.SAY_ENVIDO_QUIERO: {},
	}

	possibleEnvidoActions := make(map[string]truco.Action)
	for name, action := range possible {
		if _, ok := filter[name]; ok {
			possibleEnvidoActions[name] = action
		}
	}

	return possibleEnvidoActions
}

func possibleTrucoActionsMap(gs truco.ClientGameState) map[string]truco.Action {
	possible := possibleActionsMap(gs)

	filter := map[string]struct{}{
		truco.SAY_TRUCO_QUIERO:       {},
		truco.SAY_TRUCO:              {},
		truco.SAY_QUIERO_RETRUCO:     {},
		truco.SAY_QUIERO_VALE_CUATRO: {},
	}

	possibleTrucoActions := make(map[string]truco.Action)
	for name, action := range possible {
		if _, ok := filter[name]; ok {
			possibleTrucoActions[name] = action
		}
	}

	return possibleTrucoActions
}

func sortPossibleEnvidoActions(gs truco.ClientGameState) []truco.Action {
	possible := possibleEnvidoActionsMap(gs)
	filter := []string{
		truco.SAY_ENVIDO_QUIERO,
		truco.SAY_ENVIDO,
		truco.SAY_REAL_ENVIDO,
		truco.SAY_FALTA_ENVIDO,
	}

	actions := []truco.Action{}
	for _, name := range filter {
		if action, ok := possible[name]; ok {
			actions = append(actions, action)
		}
	}
	return actions
}

func shouldAnyEnvido(gs truco.ClientGameState, aggresiveness string) bool {
	shouldMap := map[string]int{
		"low":    29,
		"normal": 27,
		"high":   24,
	}
	score := calculateEnvidoScore(gs)

	log.Println("Bot: envido score is", score, "and aggresiveness is", aggresiveness)

	return score >= shouldMap[aggresiveness]
}

func chooseEnvidoAction(gs truco.ClientGameState, aggresiveness string) truco.Action {
	possibleActions := sortPossibleEnvidoActions(gs)
	score := calculateEnvidoScore(gs)

	minScore := map[string]int{
		"low":    29,
		"normal": 27,
		"high":   24,
	}[aggresiveness]
	maxScore := 33

	span := maxScore - minScore
	numActions := len(possibleActions)

	// Calculate bucket width
	bucketWidth := float64(span) / float64(numActions)

	// Determine the bucket for the score
	bucket := int(float64(score-minScore) / bucketWidth)

	// Handle edge cases
	if bucket < 0 {
		bucket = 0
	} else if bucket >= numActions {
		bucket = numActions - 1
	}

	return possibleActions[bucket]
}

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

func cardsWithoutLowest(cards []truco.Card) []truco.Card {
	lowest := cards[0]
	for _, card := range cards {
		if card.CompareTrucoScore(lowest) == -1 {
			lowest = card
		}
	}

	unrevealed := []truco.Card{}
	for _, card := range cards {
		if card != lowest {
			unrevealed = append(unrevealed, card)
		}
	}
	return unrevealed
}

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

func cardsWithoutLowestCardThatBeats(card truco.Card, cards []truco.Card) []truco.Card {
	return cardsWithout(cards, lowestCardThatBeats(card, cards))
}

func cardsWithoutCardThatTies(card truco.Card, cards []truco.Card) []truco.Card {
	return cardsWithout(cards, cardThatTies(card, cards))
}

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

func cardsChance(cards []truco.Card) float64 {
	base := float64(len(cards) * 4)
	divisor := float64(19.0 - base)
	if len(cards) == 2 {
		divisor = 19.0 + 18.0 - base
	}
	if len(cards) == 3 {
		divisor = 19.0 + 18.0 + 17.0 - base
	}
	sum := -base
	for _, card := range cards {
		sum += float64(calculateCardStrength(card))
	}
	return sum / divisor
}

// No cards => Hand strength
//
// They 1
//
//	If I can beat their card:
//		remaining cards strength after beating with lowest beating
//	If I can tie their card:
//		highest card's strength
//	If I can't beat their card:
//		remaining cards strength after throwing lowest card * 0.66
//
// Both 1, my turn
// Both 2, my turn
//
//	In these two cases, we're tied or I'm winning (cause wouldn't be my turn otherwise). Therefore:
//		return Highest unrevealed card's strenth
//
// They 2, me 1
//
//	if first faceoff is a tie:
//		If I can't beat their last card, 0%
//		If I can beat their last card, 100%
//		If I can tie: remaining card's strength after beating with lowest beating
//
//	if first faceoff is their win:
//		If I can't beat or I tie their last card, 0%
//		If I can beat it: remaining card's strength after beating with lowest beating
//
// They 3, me 2 =>
//
//	if I tie or lose against their last card: 0%
//	otherwise, 100%
func chanceOfWinningTruco(gs truco.ClientGameState) float64 {
	if len(gs.YourRevealedCards) == 0 && len(gs.TheirRevealedCards) == 0 {
		return calculateTrucoHandChance(gs.YourUnrevealedCards)
	}

	if len(gs.TheirRevealedCards) == 1 && len(gs.YourRevealedCards) == 0 {
		if canBeatCard(gs.TheirRevealedCards[0], gs.YourUnrevealedCards) {
			return cardsChance(cardsWithoutLowestCardThatBeats(gs.TheirRevealedCards[0], gs.YourUnrevealedCards))
		}
		if canTieCard(gs.TheirRevealedCards[0], gs.YourUnrevealedCards) {
			return cardsChance([]truco.Card{highestOf(gs.YourUnrevealedCards)})
		}
		return cardsChance(cardsWithoutLowest(gs.YourUnrevealedCards)) * 0.66
	}

	// If it's the bot's turn, it means that the faceoff was a tie or the bot is winning
	// Either way, return the highest card's chance
	if len(gs.TheirRevealedCards) == len(gs.YourRevealedCards) { // either 1,1 or 2,2
		return cardsChance([]truco.Card{highestOf(gs.YourUnrevealedCards)})
	}

	if len(gs.TheirRevealedCards) == 2 && len(gs.YourRevealedCards) == 1 {
		results := faceoffResults(gs)
		if results[0] == 0 {
			if canBeatCard(gs.TheirRevealedCards[1], gs.YourUnrevealedCards) {
				return 1.0
			}
			if canTieCard(gs.TheirRevealedCards[1], gs.YourUnrevealedCards) {
				// Note that this will be a single card anyway
				return cardsChance(cardsWithoutCardThatTies(gs.TheirRevealedCards[1], gs.YourUnrevealedCards))
			}
			return 0.0
		}
		if results[0] == -1 {
			if canBeatCard(gs.TheirRevealedCards[1], gs.YourUnrevealedCards) {
				return cardsChance(cardsWithoutLowestCardThatBeats(gs.TheirRevealedCards[1], gs.YourUnrevealedCards))
			}
			return 0.0
		}
	}

	if len(gs.TheirRevealedCards) == 3 && len(gs.YourRevealedCards) == 2 {
		if canBeatCard(gs.TheirRevealedCards[2], gs.YourUnrevealedCards) {
			return 1.0
		}
		return 0.0
	}

	// This should be unreachable, but in this case return 0.0
	return 0.0
}

func sortPossibleTrucoActions(gs truco.ClientGameState) []truco.Action {
	possible := possibleTrucoActionsMap(gs)
	filter := []string{
		truco.SAY_TRUCO_QUIERO,
		truco.SAY_TRUCO,
		truco.SAY_QUIERO_RETRUCO,
		truco.SAY_QUIERO_VALE_CUATRO,
	}

	actions := []truco.Action{}
	for _, name := range filter {
		if action, ok := possible[name]; ok {
			actions = append(actions, action)
		}
	}
	return actions
}

func chooseTrucoAction(gs truco.ClientGameState, aggresiveness string) truco.Action {
	possibleActions := sortPossibleTrucoActions(gs)
	chance := chanceOfWinningTruco(gs)
	log.Println("Bot: chanceOfWinningTruco: ", chance)

	minChance := map[string]float64{
		"low":    0.55,
		"normal": 0.5,
		"high":   0.461, // This is the average hand chance
	}[aggresiveness]
	maxChance := 1.0

	span := maxChance - minChance
	numActions := len(possibleActions)

	// Calculate bucket width
	bucketWidth := float64(span) / float64(numActions)

	// Determine the bucket for the score
	bucket := int(float64(chance-minChance) / bucketWidth)

	// Handle edge cases
	if bucket < 0 {
		bucket = 0
	} else if bucket >= numActions {
		bucket = numActions - 1
	}

	return possibleActions[bucket]
}

func shouldAcceptTruco(gs truco.ClientGameState, aggresiveness string) bool {
	shouldMap := map[string]float64{
		"low":    0.55,
		"normal": 0.5,
		"high":   0.461, // This is the average hand chance
	}
	chance := chanceOfWinningTruco(gs)
	log.Println("Bot: chanceOfWinningTruco: ", chance)

	log.Println("Bot: truco chance is", chance, "and aggresiveness is", aggresiveness)

	return chance >= shouldMap[aggresiveness]
}

func chooseCardToThrow(gs truco.ClientGameState) truco.Action {
	// If there's only one card left, throw it
	if len(gs.YourUnrevealedCards) == 1 {
		return truco.NewActionRevealCard(gs.YourUnrevealedCards[0], 1)
	}

	// If they have no revealed cards, throw the weakest card
	if len(gs.TheirRevealedCards) == 0 {
		weakestCard := gs.YourUnrevealedCards[0]
		for _, card := range gs.YourUnrevealedCards {
			if card.CompareTrucoScore(weakestCard) == -1 {
				weakestCard = card
			}
		}
		return truco.NewActionRevealCard(weakestCard, 1)
	}

	// If they have one more revealed card then me, throw the lowest card that beats their last card
	if len(gs.TheirRevealedCards) == len(gs.YourRevealedCards)+1 {
		lowestCardThatBeats := lowestCardThatBeats(gs.TheirRevealedCards[len(gs.YourRevealedCards)], gs.YourUnrevealedCards)
		if lowestCardThatBeats.Number != 0 {
			return truco.NewActionRevealCard(lowestCardThatBeats, 1)
		}
		// Otherwise throw the lowest card
		return truco.NewActionRevealCard(lowestOf(gs.YourUnrevealedCards), 1)
	}

	// If we have the same amount of revealed cards, and the last faceoff was won by me, throw the lowest card
	results := faceoffResults(gs)
	if results[len(results)-1] == 1 {
		return truco.NewActionRevealCard(lowestOf(gs.YourUnrevealedCards), 1)
	}

	// If they have the same amount of revealed cards as me, throw the highest card left
	return truco.NewActionRevealCard(highestOf(gs.YourUnrevealedCards), 1)
}

func getRandomAction(actions []truco.Action) truco.Action {
	index := rand.Intn(len(actions))
	return actions[index]
}

func sonBuenas() truco.Action {
	return truco.NewActionSaySonBuenas(1)
}
func sonMejores() truco.Action {
	return truco.NewActionSaySonMejores(0, 1)
}
func envidoNoQuiero() truco.Action {
	return truco.NewActionSayEnvidoNoQuiero(1)
}
func envidoQuiero() truco.Action {
	return truco.NewActionSayEnvidoQuiero(1)
}
func trucoQuiero() truco.Action {
	return truco.NewActionSayTrucoQuiero(1)
}
func _truco() truco.Action {
	return truco.NewActionSayTruco(1)
}
func revealCard() truco.Action {
	return truco.NewActionRevealCard(truco.Card{}, 1)
}

func (m Bot) ChooseAction(gs truco.ClientGameState) truco.Action {
	actions := possibleActionsMap(gs)
	log.Println("Bot: possible actions are", actions)

	if len(gs.PossibleActions) == 0 {
		log.Println("Bot: there are no actions left.")
		return nil
	}

	// If there's only a say_son_buenas, say_son_mejores or a single action, choose it
	sonBuenasActions := filter(actions, sonBuenas())
	if len(sonBuenasActions) > 0 {
		log.Println("Bot: I have to say son buenas.")
		return sonBuenasActions[0]
	}
	sonMejoresActions := filter(actions, sonMejores())
	if len(sonMejoresActions) > 0 {
		log.Println("Bot: I have to say son mejores.")
		return sonMejoresActions[0]
	}
	if len(gs.PossibleActions) == 1 {
		log.Println("Bot: there was only one action: ", string(gs.PossibleActions[0]))
		return _deserializeActions(gs.PossibleActions)[0]
	}

	aggresiveness := calculateAggresiveness(gs)

	// Handle envido responses or actions
	if canAnyEnvido(actions) {
		log.Println("Bot: Envido actions are on the table.")
		should := shouldAnyEnvido(gs, aggresiveness)
		log.Println("Bot: should envido?", should)

		if !should && len(filter(actions, envidoNoQuiero())) > 0 {
			log.Println("Bot: I said no quiero to envido due to considering I shouldn't based on my aggresiveness, which is", aggresiveness, "and my envido score is", calculateEnvidoScore(gs))
			return truco.NewActionSayEnvidoNoQuiero(1)
		}
		if should && len(filter(actions, envidoQuiero())) > 0 {
			log.Println("Bot: I chose an envido action due to considering I should based on my aggresiveness, which is", aggresiveness, "and my envido score is", calculateEnvidoScore(gs))
			return chooseEnvidoAction(gs, aggresiveness)
		}
		if should {
			// This is the case where the bot initiates the envido
			// Sometimes (<50%), a human player would hide their envido by not initiating, and hoping the other says it first
			// TODO: should this chance based on aggresiveness?
			if rand.Float64() < 0.67 {
				return chooseEnvidoAction(gs, aggresiveness)
			}
		}
	}

	shouldTruco := shouldAcceptTruco(gs, aggresiveness)
	log.Println("Bot: should truco?", shouldTruco)

	// Handle truco responses
	if len(filter(actions, trucoQuiero())) > 0 {
		if shouldTruco {
			log.Println("Bot: I chose a accept truco with some truco action due to considering I should based on my aggresiveness, which is", aggresiveness)
			return chooseTrucoAction(gs, aggresiveness)
		}
		log.Println("Bot: I chose to say no quiero to truco due to considering I should based on my aggresiveness, which is", aggresiveness)
		return truco.NewActionSayTrucoNoQuiero(1)
	}

	// Handle say truco
	if len(filter(actions, _truco())) > 0 && shouldTruco {
		log.Println("Bot: I chose to say truco due to considering I should based on my aggresiveness, which is", aggresiveness)
		return chooseTrucoAction(gs, aggresiveness)
	}

	// Only throw card left
	if len(filter(actions, revealCard())) > 0 {
		log.Println("Bot: I chose to reveal a card due to being the last action left.")
		return chooseCardToThrow(gs)
	}

	// This should be unreachable, but in this case choose random action
	log.Println("Bot: I shouldn't have arrived here, so I'm choosing a random action.")
	return getRandomAction(_deserializeActions(gs.PossibleActions))
}
