package examplebot

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"math/rand"

	"github.com/marianogappa/truco/truco"
)

type Bot struct {
	log func(format string, v ...any)
}

func New() Bot {
	return Bot{log: func(format string, v ...any) { log.Printf(fmt.Sprintf("Bot: %v\n", format), v...) }}
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

func calculateFlorScore(gs truco.ClientGameState) int {
	return truco.Hand{Revealed: gs.YourRevealedCards, Unrevealed: gs.YourUnrevealedCards}.FlorScore()
}

func calculateCardStrength(gs truco.Card) int {
	specialValues := map[truco.Card]int{
		{Suit: truco.ESPADA, Number: 1}: 15,
		{Suit: truco.BASTO, Number: 1}:  14,
		{Suit: truco.ESPADA, Number: 7}: 13,
		{Suit: truco.ORO, Number: 7}:    12,
	}
	if _, ok := specialValues[gs]; ok {
		return specialValues[gs]
	}
	if gs.Number <= 3 {
		return gs.Number + 12 - 4
	}
	return gs.Number - 4
}

func faceoffResults(gs truco.ClientGameState) []int {
	results := []int{}
	for i := 0; i < min(len(gs.YourRevealedCards), len(gs.TheirRevealedCards)); i++ {
		results = append(results, gs.YourRevealedCards[i].CompareTrucoScore(gs.TheirRevealedCards[i]))
	}
	return results
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

func canAnyFlor(actions map[string]truco.Action) bool {
	return len(filter(actions,
		truco.NewActionSayFlor(1),
		truco.NewActionSayContraflor(1),
		truco.NewActionSayContraflorAlResto(1),
		truco.NewActionSayConFlorQuiero(1),
		truco.NewActionSayConFlorMeAchico(1),
	)) > 0
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
	possible := possibleActionsMap(gs)
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

	// Sort actions based on their cost
	// TODO: this is broken at the moment because the cost doesn't work well
	sort.Slice(actions, func(i, j int) bool {
		return _getEnvidoActionQuieroCost(actions[i]) < _getEnvidoActionQuieroCost(actions[j])
	})

	return actions
}

func sortPossibleFlorActions(gs truco.ClientGameState) []truco.Action {
	possible := possibleActionsMap(gs)
	filter := []string{
		truco.SAY_FLOR,
		truco.SAY_CON_FLOR_QUIERO,
		truco.SAY_CONTRAFLOR,
		truco.SAY_CONTRAFLOR_AL_RESTO,
	}

	actions := []truco.Action{}
	for _, name := range filter {
		if action, ok := possible[name]; ok {
			actions = append(actions, action)
		}
	}

	// Sort actions based on their cost
	// TODO: this is broken at the moment because the cost doesn't work well
	sort.Slice(actions, func(i, j int) bool {
		return _getFlorActionQuieroCost(actions[i]) < _getFlorActionQuieroCost(actions[j])
	})

	return actions
}

func _getFlorActionQuieroCost(action truco.Action) int {
	switch a := action.(type) {
	case *truco.ActionSayFlor:
		return a.QuieroCost
	case *truco.ActionSayConFlorQuiero:
		return a.Cost
	case *truco.ActionSayContraflor:
		return a.QuieroCost
	case *truco.ActionSayContraflorAlResto:
		return a.QuieroCost
	default:
		panic("this code should be unreachable! bug in _getFlorActionCost! please report this bug.")
	}
}

func _getEnvidoActionQuieroCost(action truco.Action) int {
	switch a := action.(type) {
	case *truco.ActionSayEnvidoQuiero:
		return a.Cost
	case *truco.ActionSayEnvido:
		return a.QuieroCost
	case *truco.ActionSayRealEnvido:
		return a.QuieroCost
	case *truco.ActionSayFaltaEnvido:
		return a.QuieroCost
	default:
		panic("this code should be unreachable! bug in _getEnvidoActionCost! please report this bug.")
	}
}

func shouldAnyEnvido(gs truco.ClientGameState, aggresiveness string, log func(string, ...any)) bool {
	// if "no quiero" is possible and saying no quiero means losing, return true
	possible := possibleActionsMap(gs)
	noQuieroActions := filter(possible, truco.NewActionSayEnvidoNoQuiero(gs.YouPlayerID))
	if len(noQuieroActions) > 0 {
		cost := noQuieroActions[0].(*truco.ActionSayEnvidoNoQuiero).Cost
		if gs.TheirScore+cost >= gs.RuleMaxPoints {
			return true
		}
	}

	shouldMap := map[string]int{
		"low":    29,
		"normal": 27,
		"high":   24,
	}
	score := calculateEnvidoScore(gs)

	log("shouldAcceptEnvido: should[%v] = %v, score = %v", aggresiveness, shouldMap[aggresiveness], score)

	return score >= shouldMap[aggresiveness]
}

func shouldAnyFlor(gs truco.ClientGameState, aggresiveness string, log func(string, ...any)) bool {
	// If bot doesn't have flor, bot shouldn't say flor
	if calculateFlorScore(gs) == 0 {
		log("I don't have flor, so I'm not going to say flor.")
		return false
	}

	possible := possibleActionsMap(gs)

	// If human doesn't necessarily have flor, and bot has flor then bot should say flor
	if quieroActions := filter(possible, truco.NewActionSayConFlorQuiero(gs.YouPlayerID)); len(quieroActions) == 0 {
		log("Human doesn't necessarily have flor, so I'm going to say flor.")
		return true
	}

	// if "no quiero" is possible and saying no quiero means losing, return true
	noQuieroActions := filter(possible, truco.NewActionSayConFlorMeAchico(gs.YouPlayerID))
	if len(noQuieroActions) > 0 {
		log("Bot can say no quiero to flor, and saying no quiero might mean losing.")
		cost := noQuieroActions[0].(*truco.ActionSayConFlorMeAchico).Cost
		if gs.TheirScore+cost >= gs.RuleMaxPoints {
			log("Bot should say quiero to flor, because bot loses otherwise.")
			return true
		}
	}

	// Both have flor and saying no quiero doesn't lose the game. At this point it depends on the score and the aggresiveness
	log("Bot has flor, and human has flor. Bot's flor score is %v, and aggresiveness is: %v", calculateFlorScore(gs), aggresiveness)
	return calculateFlorScore(gs) >= map[string]int{
		"low":    31,
		"normal": 29,
		"high":   26,
	}[aggresiveness]
}

func chooseFlorAction(gs truco.ClientGameState, aggresiveness string) truco.Action {
	possibleActions := sortPossibleFlorActions(gs)
	score := calculateFlorScore(gs)

	minScore := map[string]int{
		"low":    31,
		"normal": 29,
		"high":   26,
	}[aggresiveness]
	maxScore := 38

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
	divisor := float64([]float64{1, 15.0, 15.0 + 14.0, 15.0 + 14.0 + 13.0}[len(cards)])
	sum := 0.0
	for _, card := range cards {
		sum += float64(calculateCardStrength(card))
	}
	return sum / divisor
}

func cardsChanceTwoAttempts(cards []truco.Card) float64 {
	highestNumber := float64(calculateCardStrength(highestOf(cards)))
	return highestNumber/15.0 + (15.0-highestNumber)/(15.0*15.0)
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
	if len(gs.YourRevealedCards) <= 1 && len(gs.TheirRevealedCards) == 0 {
		return cardsChance(append(gs.YourRevealedCards, gs.YourUnrevealedCards...))
	}

	if len(gs.TheirRevealedCards) == 2 && len(gs.YourRevealedCards) == 3 {
		return cardsChance([]truco.Card{gs.YourRevealedCards[2]})
	}

	if len(gs.TheirRevealedCards) == 1 && len(gs.YourRevealedCards) == 0 {
		if canBeatCard(gs.TheirRevealedCards[0], gs.YourUnrevealedCards) {
			return cardsChance(cardsWithoutLowestCardThatBeats(gs.TheirRevealedCards[0], gs.YourUnrevealedCards))
		}
		if canTieCard(gs.TheirRevealedCards[0], gs.YourUnrevealedCards) {
			return cardsChance([]truco.Card{highestOf(gs.YourUnrevealedCards)})
		}
		// In this case, bot cannot win the first faceoff. Therefore, in order to win, the next two faceoffs have to be won
		chance := cardsChance(cardsWithoutLowest(gs.YourUnrevealedCards))
		return chance * chance
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

	// Bot won first round
	if len(gs.TheirRevealedCards) == 1 && len(gs.YourRevealedCards) == 2 {
		// In this case the bot only has to win one of the next two faceoffs
		chance := cardsChanceTwoAttempts([]truco.Card{gs.YourUnrevealedCards[0], gs.YourRevealedCards[len(gs.YourRevealedCards)-1]})
		return chance
	}

	// This should be unreachable, but in this case return 0.0
	panic("this code should be unreachable! bug in chanceOfWinningTruco! please report this bug.")
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

func shouldAcceptTruco(gs truco.ClientGameState, aggresiveness string, log func(string, ...any)) bool {
	// if "no quiero" is possible and saying no quiero means losing, return true
	possible := possibleActionsMap(gs)
	noQuieroActions := filter(possible, truco.NewActionSayTrucoNoQuiero(gs.YouPlayerID))
	if len(noQuieroActions) > 0 {
		cost := noQuieroActions[0].(*truco.ActionSayTrucoNoQuiero).Cost
		if gs.TheirScore+cost >= gs.RuleMaxPoints {
			return true
		}
	}

	shouldMap := map[string]float64{
		"low":    0.55,
		"normal": 0.5,
		"high":   0.461, // This is the average hand chance
	}
	chance := chanceOfWinningTruco(gs)
	log("shouldAcceptTruco: should[%v] = %v, chance = %v", aggresiveness, shouldMap[aggresiveness], chance)
	return chance >= shouldMap[aggresiveness]
}

func losesHandWithNextCard(gs truco.ClientGameState) bool {
	if len(gs.TheirRevealedCards) < 2 {
		return false // This face off doesn't decide who wins
	}
	if len(gs.TheirRevealedCards) != len(gs.YourRevealedCards)+1 {
		return false // It's not the bot's turn to play a card
	}
	var (
		youMano         = gs.RoundTurnPlayerID == gs.YouPlayerID
		faceoffResults  = faceoffResults(gs)
		theirCard       = gs.TheirRevealedCards[len(gs.YourRevealedCards)]
		yourHighestCard = highestOf(gs.YourUnrevealedCards)
	)
	// The result of the current faceoff between bot & other
	switch yourHighestCard.CompareTrucoScore(theirCard) {
	case 1:
		return false // Bot wins, so it doesn't lose with the next card
	case -1:
		return true // Bot loses
	case 0: // If bot ties, then it depends on previous faceoffs
		switch len(faceoffResults) {
		// There was only one previous faceoff
		case 1:
			switch faceoffResults[0] {
			case 0, 1: // If bot tied or won, a tie doesn't lose the hand
				return false
			case -1: // If bot lost, a tie loses the hand
				return true
			}
		case 2:
			// If bot won any of the previous faceoffs, a tie doesn't lose the hand
			if faceoffResults[0] == 1 || faceoffResults[1] == 1 {
				return false
			}
			// If bot lost any of the previous faceoffs, a tie loses the hand
			if faceoffResults[0] == -1 || faceoffResults[1] == -1 {
				return true
			}
			// If both faceoffs were ties, then it depends on who's mano
			if faceoffResults[0] == 0 && faceoffResults[1] == 0 {
				return !youMano
			}
		}
	}
	panic("this code should be unreachable! bug in losesHandWithNextCard! please report this bug.")
}

func chooseCardToThrow(gs truco.ClientGameState, log func(string, ...any)) truco.Action {
	actions := possibleActionsMap(gs)
	// If me_voy_al_mazo is possible and the card is lower than the other's revealed card, say me_voy_al_mazo
	if len(filter(actions, meVoy(gs))) > 0 && len(gs.TheirRevealedCards) > len(gs.YourRevealedCards) && losesHandWithNextCard(gs) {
		log("I'm losing the hand with the next card, so I'm going to say me_voy_al_mazo")
		return truco.NewActionSayMeVoyAlMazo(gs.YouPlayerID)
	}

	// If there's only one card left, throw it
	if len(gs.YourUnrevealedCards) == 1 {
		return truco.NewActionRevealCard(gs.YourUnrevealedCards[0], gs.YouPlayerID)
	}

	// If they have no revealed cards, throw the weakest card
	if len(gs.TheirRevealedCards) == 0 {
		weakestCard := gs.YourUnrevealedCards[0]
		for _, card := range gs.YourUnrevealedCards {
			if card.CompareTrucoScore(weakestCard) == -1 {
				weakestCard = card
			}
		}
		return truco.NewActionRevealCard(weakestCard, gs.YouPlayerID)
	}

	// If they have one more revealed card then me, throw the lowest card that beats their last card
	if len(gs.TheirRevealedCards) == len(gs.YourRevealedCards)+1 {
		lowestCardThatBeats := lowestCardThatBeats(gs.TheirRevealedCards[len(gs.YourRevealedCards)], gs.YourUnrevealedCards)
		if lowestCardThatBeats.Number != 0 {
			return truco.NewActionRevealCard(lowestCardThatBeats, gs.YouPlayerID)
		}
		// Otherwise throw the lowest card
		return truco.NewActionRevealCard(lowestOf(gs.YourUnrevealedCards), gs.YouPlayerID)
	}

	// If we have the same amount of revealed cards, and the last faceoff was won by me, throw the lowest card
	results := faceoffResults(gs)
	if results[len(results)-1] == 1 {
		return truco.NewActionRevealCard(lowestOf(gs.YourUnrevealedCards), gs.YouPlayerID)
	}

	// If they have the same amount of revealed cards as me, throw the highest card left
	return truco.NewActionRevealCard(highestOf(gs.YourUnrevealedCards), gs.YouPlayerID)
}

func getRandomAction(actions []truco.Action) truco.Action {
	index := rand.Intn(len(actions))
	return actions[index]
}

func sonBuenas(gs truco.ClientGameState) truco.Action {
	return truco.NewActionSaySonBuenas(gs.YouPlayerID)
}
func sonMejores(gs truco.ClientGameState) truco.Action {
	return truco.NewActionSaySonMejores(gs.YouPlayerID)
}
func envidoNoQuiero(gs truco.ClientGameState) truco.Action {
	return truco.NewActionSayEnvidoNoQuiero(gs.YouPlayerID)
}
func envidoQuiero(gs truco.ClientGameState) truco.Action {
	return truco.NewActionSayEnvidoQuiero(gs.YouPlayerID)
}
func florQuiero(gs truco.ClientGameState) truco.Action {
	return truco.NewActionSayConFlorQuiero(gs.YouPlayerID)
}
func florNoQuiero(gs truco.ClientGameState) truco.Action {
	return truco.NewActionSayConFlorMeAchico(gs.YouPlayerID)
}
func trucoQuiero(gs truco.ClientGameState) truco.Action {
	return truco.NewActionSayTrucoQuiero(gs.YouPlayerID)
}
func _truco(gs truco.ClientGameState) truco.Action {
	return truco.NewActionSayTruco(gs.YouPlayerID)
}
func revealCard(gs truco.ClientGameState) truco.Action {
	return truco.NewActionRevealCard(truco.Card{}, gs.YouPlayerID)
}
func meVoy(gs truco.ClientGameState) truco.Action {
	return truco.NewActionSayMeVoyAlMazo(gs.YouPlayerID)
}

func (m Bot) ChooseAction(gs truco.ClientGameState) truco.Action {
	if len(gs.PossibleActions) == 0 {
		m.log("there are no actions left.")
		return nil
	}
	if len(gs.PossibleActions) == 1 {
		m.log("there was only one action: %v", string(gs.PossibleActions[0]))
		return _deserializeActions(gs.PossibleActions)[0]
	}

	// If there's only a say_son_buenas, say_son_mejores or a single action, choose it
	actions := possibleActionsMap(gs)
	for _, action := range actions {
		m.log("possible action: %v", action)
	}
	sonBuenasActions := filter(actions, sonBuenas(gs))
	if len(sonBuenasActions) > 0 {
		m.log("I have to say son buenas.")
		return sonBuenasActions[0]
	}
	sonMejoresActions := filter(actions, sonMejores(gs))
	if len(sonMejoresActions) > 0 {
		m.log("I have to say son mejores.")
		return sonMejoresActions[0]
	}

	var (
		aggresiveness = calculateAggresiveness(gs)
		shouldEnvido  = shouldAnyEnvido(gs, aggresiveness, m.log)
		shouldFlor    = shouldAnyFlor(gs, aggresiveness, m.log)
		shouldTruco   = shouldAcceptTruco(gs, aggresiveness, m.log)
	)

	// Handle flor responses or actions
	if canAnyFlor(actions) {
		m.log("Flor actions are on the table.")

		if shouldFlor && len(filter(actions, florQuiero(gs))) > 0 {
			m.log("I chose an flor action due to considering I should based on my aggresiveness, which is %v and my flor score is %v", aggresiveness, calculateFlorScore(gs))
			return chooseFlorAction(gs, aggresiveness)
		}
		if !shouldFlor && len(filter(actions, florNoQuiero(gs))) > 0 {
			m.log("I said no quiero to flor due to considering I shouldn't based on my aggresiveness, which is %v and my flor score is %v", aggresiveness, calculateFlorScore(gs))
			return truco.NewActionSayConFlorMeAchico(gs.YouPlayerID)
		}
		if shouldFlor {
			// This is the case where the bot initiates the flor
			// Sometimes (<50%), a human player would hide their envido by not initiating, and hoping the other says it first
			// TODO: should this chance based on aggresiveness?
			if rand.Float64() < 0.67 {
				return chooseFlorAction(gs, aggresiveness)
			}
		}
	}

	// Handle envido responses or actions
	if canAnyEnvido(actions) {
		m.log("Envido actions are on the table.")

		if shouldEnvido && len(filter(actions, envidoQuiero(gs))) > 0 {
			m.log("I chose an envido action due to considering I should based on my aggresiveness, which is %v and my envido score is %v", aggresiveness, calculateEnvidoScore(gs))
			return chooseEnvidoAction(gs, aggresiveness)
		}
		if !shouldEnvido && len(filter(actions, envidoNoQuiero(gs))) > 0 {
			m.log("I said no quiero to envido due to considering I shouldn't based on my aggresiveness, which is %v and my envido score is %v", aggresiveness, calculateEnvidoScore(gs))
			return truco.NewActionSayEnvidoNoQuiero(gs.YouPlayerID)
		}
		if shouldEnvido {
			// This is the case where the bot initiates the envido
			// Sometimes (<50%), a human player would hide their envido by not initiating, and hoping the other says it first
			// TODO: should this chance based on aggresiveness?
			if rand.Float64() < 0.67 {
				return chooseEnvidoAction(gs, aggresiveness)
			}
		}
	}

	// Handle truco responses
	if len(filter(actions, trucoQuiero(gs))) > 0 {
		m.log("I have to answer a truco question. My previous analysis is: %v", shouldTruco)
		if shouldTruco {
			m.log("Choosing truco acceptance action")
			return chooseTrucoAction(gs, aggresiveness)
		}
		m.log("Choosing no quiero truco action")
		return truco.NewActionSayTrucoNoQuiero(gs.YouPlayerID)
	}

	// Handle say truco
	if len(filter(actions, _truco(gs))) > 0 && shouldTruco {
		m.log("Even though I haven't been asked, I'm going to say truco due to analysis that I should.")
		return chooseTrucoAction(gs, aggresiveness)
	}

	// Only throw card left
	if len(filter(actions, revealCard(gs))) > 0 {
		m.log("I chose to reveal a card due to being the last action left.")
		return chooseCardToThrow(gs, m.log)
	}

	// This should be unreachable, but in this case choose random action
	m.log("I shouldn't have arrived here, so I'm choosing a random action.")
	return getRandomAction(_deserializeActions(gs.PossibleActions))
}
