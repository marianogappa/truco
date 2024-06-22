package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/marianogappa/truco/truco"
	"github.com/nsf/termbox-go"
)

func play(playerID int, gameState truco.GameState) (truco.Action, error) {
	err := printState(playerID, gameState, true, false)
	if err != nil {
		return nil, err
	}

	if gameState.IsEnded {
		return nil, nil
	}

	var (
		action truco.Action
		input  string
	)
	for {
		input, err = readChar()
		if err != nil {
			return nil, err
		}

		if input == "exit" || input == "quit" {
			break
		}
		num, err := strconv.Atoi(input)
		if err == nil {
			var actionName string
			var err error
			actionName, input, err = numToAction(num, gameState)
			if err != nil {
				continue
			}
			if actionName == truco.SAY_ENVIDO_QUIERO || actionName == truco.SAY_SON_BUENAS || actionName == truco.SAY_SON_MEJORES {
				input = fmt.Sprintf(`{"name":"%v","score":%d}`, actionName, gameState.Hands[gameState.TurnPlayerID].EnvidoScore())
			}
			if actionName == "reveal_card" {
				err := printState(playerID, gameState, false, false)
				if err != nil {
					return nil, err
				}
				input, err = readChar()
				if err != nil {
					return nil, err
				}

				which, err := strconv.Atoi(input)
				if err != nil {
					return nil, err
				}
				card := gameState.Hands[gameState.TurnPlayerID].Unrevealed[which-1]
				jsonCard, _ := json.Marshal(card)
				input = fmt.Sprintf(`{"name":"reveal_card","card":%v}`, string(jsonCard))
			}
		}

		action, err = truco.DeserializeAction([]byte(input))
		if err != nil {
			fmt.Printf("Invalid action:	%v\n", err)
			continue
		}
		break
	}
	return action, nil
}

func readChar() (string, error) {
	event := termbox.PollEvent()

	if event.Type != termbox.EventKey {
		return readChar()
	}

	if event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlD || event.Key == termbox.KeyCtrlZ || event.Ch == 'q' {
		os.Exit(0)
	}

	switch event.Ch {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return string(event.Ch), nil
	}

	return readChar()
}

func readAnyKey() {
	event := termbox.PollEvent()

	if event.Type != termbox.EventKey && event.Key != termbox.KeySpace {
		readAnyKey()
	}
}

func numToAction(num int, state truco.GameState) (string, string, error) {
	actions := state.CalculatePossibleActions()
	if num > len(actions) {
		return "", "", fmt.Errorf("Invalid action")
	}

	return actions[num-1], fmt.Sprintf(`{"name":"%v"}`, actions[num-1]), nil
}

func printState(playerID int, state truco.GameState, initialActions bool, showRoundResult bool) error {
	err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	if err != nil {
		return err
	}

	var (
		mx, my = termbox.Size()
		you    = playerID
		them   = state.OpponentOf(you)
		hand   = *state.Hands[them]
	)

	if showRoundResult {
		hand = *state.HandsDealt[len(state.HandsDealt)-2][them]
	}

	var (
		unrevealed = strings.Repeat("[] ", len(hand.Unrevealed))
		revealed   = getCardsString(hand.Revealed, false)
	)

	printAt(0, 0, unrevealed)
	printAt(0, my/2-3, revealed)

	printUpToAt(mx-1, 0, fmt.Sprintf("Round %d", state.RoundNumber))
	printUpToAt(mx-1, 1, fmt.Sprintf("You %d points", state.Scores[you]))
	printUpToAt(mx-1, 2, fmt.Sprintf("Them %d points", state.Scores[them]))

	hand = *state.Hands[you]

	if showRoundResult {
		hand = *state.HandsDealt[len(state.HandsDealt)-2][you]
	}

	unrevealed = getCardsString(hand.Unrevealed, false)
	revealed = getCardsString(hand.Revealed, false)

	printAt(0, my/2+3, revealed)
	printAt(0, my-4, unrevealed)

	if !showRoundResult {
		lastActionString, err := getLastActionString(you, state)
		if err != nil {
			return err
		}

		printAt(0, my/2, lastActionString)
	} else {
		lastRoundResult := state.RoundResults[len(state.RoundResults)-1]

		envidoPart := "envido was not played"
		if lastRoundResult.EnvidoWinnerPlayerID != -1 {
			envidoWinner := "you"
			if lastRoundResult.EnvidoWinnerPlayerID == them {
				envidoWinner = "they"
			}
			envidoPart = fmt.Sprintf("%v won %v envido points", envidoWinner, lastRoundResult.EnvidoPoints)
		}
		trucoWinner := "you"
		if lastRoundResult.TrucoWinnerPlayerID == them {
			trucoWinner = "they"
		}

		result := fmt.Sprintf(
			"Round ended, %v and %v won %v truco points!",
			envidoPart,
			trucoWinner,
			lastRoundResult.TrucoPoints,
		)
		printAt(0, my/2, result)
	}

	if showRoundResult {
		printAt(0, my-2, "Press any key to continue...")
		termbox.Flush()
		readAnyKey()
		return nil
	}

	if state.TurnPlayerID == playerID {
		if initialActions {
			printAt(0, my-2, "Available Actions: ")
			actionsString := ""
			for i, action := range state.PossibleActions {
				action := strings.TrimPrefix(action, "say_")
				if strings.HasSuffix(action, "no_quiero") {
					action = "no quiero"
				} else if strings.HasSuffix(action, "_quiero") {
					action = "quiero"
				}
				actionsString += fmt.Sprintf("%d. %s, ", i+1, action)
			}
			printAt(0, my-1, actionsString)
		} else {
			printAt(0, my-2, "Which card do you want to reveal: ")
			unrevealed = getCardsString(hand.Unrevealed, true)
			printAt(0, my-1, unrevealed)
		}
	} else {
		_, my := termbox.Size()
		printAt(0, my-2, "Waiting for the other player...")
	}

	termbox.Flush()
	return nil
}

func printAt(x, y int, s string) {
	for i, r := range s {
		termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

// Write so that the output ends at x, y
func printUpToAt(x, y int, s string) {
	for i, r := range s {
		termbox.SetCell(x-len(s)+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func getCardsString(cards []truco.Card, withNumbers bool) string {
	var cs []string
	for i, card := range cards {
		if withNumbers {
			cs = append(cs, fmt.Sprintf("%v. %v", i+1, getCardString(card)))
		} else {
			cs = append(cs, getCardString(card))
		}
	}
	return strings.Join(cs, " ")
}

func getCardString(card truco.Card) string {
	return fmt.Sprintf("[  %v  %v]", card.Number, suitEmoji(card.Suit))
}

func suitEmoji(suit string) string {
	switch suit {
	case truco.ESPADA:
		return "🔪"
	case truco.BASTO:
		return "🌿"
	case truco.ORO:
		return "💰"
	case truco.COPA:
		return "🍷"
	default:
		return "❓"
	}
}

func getLastActionString(playerID int, state truco.GameState) (string, error) {
	if len(state.Actions) == 0 {
		return "Game started!", nil
	}
	if state.RoundJustStarted {
		return "Round started!", nil
	}

	lastActionBs := state.Actions[len(state.Actions)-1]
	lastActionOwnerPlayerID := state.ActionOwnerPlayerIDs[len(state.ActionOwnerPlayerIDs)-1]
	return getActionString(lastActionBs, lastActionOwnerPlayerID, playerID)
}

func getActionString(lastActionBs json.RawMessage, lastActionOwnerPlayerID int, playerID int) (string, error) {
	lastAction, err := truco.DeserializeAction(lastActionBs)
	if err != nil {
		return "", err
	}

	who := "You"
	if playerID != lastActionOwnerPlayerID {
		who = "They"
	}

	var what string
	switch lastAction.GetName() {
	case "reveal_card":
		action := lastAction.(*truco.ActionRevealCard)
		what = fmt.Sprintf("revealed a %v!", getCardString(action.Card))
	case "say_envido":
		what = "said Envido!"
	case "say_real_envido":
		what = "said Real Envido!"
	case "say_falta_envido":
		what = "said Falta Envido!"
	case "say_envido_quiero":
		action := lastAction.(*truco.ActionSayEnvidoQuiero)
		what = fmt.Sprintf("said Quiero with %d!", action.Score)
	case "say_envido_no_quiero":
		what = "said No Quiero!"
	case "say_truco":
		what = "said Truco!"
	case "say_truco_quiero":
		what = "said Quiero!"
	case "say_truco_no_quiero":
		what = "said No Quiero!"
	case "say_quiero_retruco":
		what = "said Quiero Retruco!"
	case "say_quiero_vale_cuatro":
		what = "said Quiero Vale Cuatro!"
	case "say_son_buenas":
		what = "said Son Buenas!"
	case "say_son_mejores":
		action := lastAction.(*truco.ActionSaySonMejores)
		what = fmt.Sprintf("said %d son mejores!", action.Score)
	case "say_me_voy_al_mazo":
		what = "said Me Voy Al Mazo!"
	default:
		what = "???"
	}

	return fmt.Sprintf("%v %v\n", who, what), nil
}
