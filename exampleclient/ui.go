package exampleclient

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/marianogappa/truco/truco"
	"github.com/nsf/termbox-go"
)

type ui struct {
	wantKeyPressCh chan struct{}
	sendKeyPressCh chan rune
}

func NewUI() *ui {
	ui := &ui{
		wantKeyPressCh: make(chan struct{}),
		sendKeyPressCh: make(chan rune),
	}
	ui.startKeyEventLoop()
	return ui
}

func (u *ui) play(playerID int, gameState truco.GameState) (truco.Action, error) {
	err := u.printState(playerID, gameState, PRINT_MODE_NORMAL)
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
		num := u.pressAnyNumber()
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
			err := u.printState(playerID, gameState, PRINT_MODE_WHICH_CARD_REVEAL)
			if err != nil {
				return nil, err
			}
			var card truco.Card
			for {
				which := u.pressAnyNumber()
				if which > len(gameState.Hands[gameState.TurnPlayerID].Unrevealed) {
					continue
				}
				if which == 0 {
					return u.play(playerID, gameState)
				}
				card = gameState.Hands[gameState.TurnPlayerID].Unrevealed[which-1]
				break
			}
			jsonCard, _ := json.Marshal(card)
			input = fmt.Sprintf(`{"name":"reveal_card","card":%v}`, string(jsonCard))
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

func numToAction(num int, state truco.GameState) (string, string, error) {
	actions := state.CalculatePossibleActions()
	if num > len(actions) {
		return "", "", fmt.Errorf("Invalid action")
	}

	return actions[num-1], fmt.Sprintf(`{"name":"%v"}`, actions[num-1]), nil
}

type printMode int

const (
	PRINT_MODE_NORMAL printMode = iota
	PRINT_MODE_WHICH_CARD_REVEAL
	PRINT_MODE_SHOW_ROUND_RESULT
	PRINT_MODE_END
)

func (u *ui) printState(playerID int, state truco.GameState, mode printMode) error {
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

	if mode == PRINT_MODE_SHOW_ROUND_RESULT {
		hand = *state.HandsDealt[len(state.HandsDealt)-2][them]
	}

	var (
		unrevealed = strings.Repeat("[] ", len(hand.Unrevealed))
		revealed   = getCardsString(hand.Revealed, false, false)
	)

	printAt(0, 0, unrevealed)
	printAt(0, my/2-3, revealed)

	printUpToAt(mx-1, 0, fmt.Sprintf("Mano número %d", state.RoundNumber))

	youMano := ""
	themMano := ""
	if state.TurnPlayerID == you {
		youMano = " (mano)"
	} else {
		themMano = " (mano)"
	}

	printUpToAt(mx-1, 1, fmt.Sprintf("Vos%v %v", youMano, spanishScore(state.Scores[you])))
	printUpToAt(mx-1, 2, fmt.Sprintf("Elle%v %v", themMano, spanishScore(state.Scores[them])))

	hand = *state.Hands[you]

	if mode == PRINT_MODE_SHOW_ROUND_RESULT {
		hand = *state.HandsDealt[len(state.HandsDealt)-2][you]
	}

	unrevealed = getCardsString(hand.Unrevealed, false, false)
	revealed = getCardsString(hand.Revealed, false, false)

	printAt(0, my/2+3, revealed)
	printAt(0, my-4, unrevealed)

	switch mode {
	case PRINT_MODE_NORMAL, PRINT_MODE_WHICH_CARD_REVEAL:
		lastActionString, err := getLastActionString(you, state)
		if err != nil {
			return err
		}

		printAt(0, my/2, lastActionString)
	case PRINT_MODE_SHOW_ROUND_RESULT:
		lastActionString, err := getActionString(state.Actions[len(state.Actions)-1], state.ActionOwnerPlayerIDs[len(state.ActionOwnerPlayerIDs)-1], you)
		if err != nil {
			return err
		}

		printAt(0, my/2, lastActionString)
		lastRoundResult := state.RoundResults[len(state.RoundResults)-1]

		envidoPart := "el envido no se jugó"
		if lastRoundResult.EnvidoWinnerPlayerID != -1 {
			envidoWinner := "vos"
			won := "ganaste"
			if lastRoundResult.EnvidoWinnerPlayerID == them {
				envidoWinner = "elle"
				won = "ganó"
			}
			envidoPart = fmt.Sprintf("%v %v %v puntos por el envido", envidoWinner, won, lastRoundResult.EnvidoPoints)
		}
		trucoWinner := "vos"
		won := "ganaste"
		if lastRoundResult.TrucoWinnerPlayerID == them {
			trucoWinner = "elle"
			won = "ganó"
		}

		result := fmt.Sprintf(
			"Terminó la mano, %v y %v %v %v puntos por el truco.",
			envidoPart,
			trucoWinner,
			won,
			lastRoundResult.TrucoPoints,
		)
		printAt(0, my/2+1, result)
	case PRINT_MODE_END:
		lastActionString, err := getActionString(state.Actions[len(state.Actions)-1], state.ActionOwnerPlayerIDs[len(state.ActionOwnerPlayerIDs)-1], you)
		if err != nil {
			return err
		}

		if playerID == state.WinnerPlayerID {
			printAt(0, my/2, fmt.Sprintf("%v Ganaste 🥰!", lastActionString))
		} else {
			printAt(0, my/2, fmt.Sprintf("%v Perdiste 😭!", lastActionString))
		}
	}

	if mode == PRINT_MODE_SHOW_ROUND_RESULT || mode == PRINT_MODE_END {
		printAt(0, my-2, "Presioná cualquier tecla para continuar...")
		termbox.Flush()
		u.pressAnyKey()
		return nil
	}

	if state.TurnPlayerID == playerID {
		if mode == PRINT_MODE_NORMAL {
			actionsString := ""
			for i, action := range state.PossibleActions {
				action := spanishAction(action, state)
				actionsString += fmt.Sprintf("%d. %s   ", i+1, action)
			}
			printAt(0, my-2, actionsString)
		} else if mode == PRINT_MODE_WHICH_CARD_REVEAL {
			printAt(0, my-2, "¿Cuál carta querés tirar?")
			unrevealed = getCardsString(hand.Unrevealed, true, true)
			printAt(0, my-1, unrevealed)
		}
	} else {
		_, my := termbox.Size()
		printAt(0, my-2, "Esperando al otro jugador...")
	}

	termbox.Flush()
	return nil
}

func printAt(x, y int, s string) {
	_s := []rune(s)
	for i, r := range _s {
		termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

// Write so that the output ends at x, y
func printUpToAt(x, y int, s string) {
	_s := []rune(s)
	for i, r := range _s {
		termbox.SetCell(x-len(_s)+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func getCardsString(cards []truco.Card, withNumbers bool, withBack bool) string {
	var cs []string
	for i, card := range cards {
		if withNumbers {
			cs = append(cs, fmt.Sprintf("%v. %v", i+1, getCardString(card)))
		} else {
			cs = append(cs, getCardString(card))
		}
	}
	if withBack {
		cs = append(cs, "0. Volver")
	}
	return strings.Join(cs, "  ")
}

func getCardString(card truco.Card) string {
	return fmt.Sprintf("[%v%v ]", card.Number, suitEmoji(card.Suit))
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
		return "¡Empezó el juego!", nil
	}
	if state.RoundJustStarted {
		return "¡Empezó la mano!", nil
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

	said := "dijiste"
	revealed := "tiraste"
	who := "Vos"
	if playerID != lastActionOwnerPlayerID {
		who = "Elle"
		said = "dijo"
		revealed = "tiró"
	}

	var what string
	switch lastAction.GetName() {
	case truco.REVEAL_CARD:
		action := lastAction.(*truco.ActionRevealCard)
		what = fmt.Sprintf("%v la carta %v", revealed, getCardString(action.Card))
	case truco.SAY_ENVIDO:
		what = fmt.Sprintf("%v envido", said)
	case truco.SAY_REAL_ENVIDO:
		what = fmt.Sprintf("%v real envido", said)
	case truco.SAY_FALTA_ENVIDO:
		what = fmt.Sprintf("%v falta envido!", said)
	case truco.SAY_ENVIDO_QUIERO:
		action := lastAction.(*truco.ActionSayEnvidoQuiero)
		what = fmt.Sprintf("%v quiero con %d", said, action.Score)
	case truco.SAY_ENVIDO_NO_QUIERO:
		what = fmt.Sprintf("%v no quiero", said)
	case truco.SAY_TRUCO:
		what = fmt.Sprintf("%v truco", said)
	case truco.SAY_TRUCO_QUIERO:
		what = fmt.Sprintf("%v quiero", said)
	case truco.SAY_TRUCO_NO_QUIERO:
		what = fmt.Sprintf("%v no quiero", said)
	case truco.SAY_QUIERO_RETRUCO:
		what = fmt.Sprintf("%v quiero retruco", said)
	case truco.SAY_QUIERO_VALE_CUATRO:
		what = fmt.Sprintf("%v quiero vale cuatro", said)
	case truco.SAY_SON_BUENAS:
		what = fmt.Sprintf("%v son buenas", said)
	case truco.SAY_SON_MEJORES:
		action := lastAction.(*truco.ActionSaySonMejores)
		what = fmt.Sprintf("%v %d son mejores", said, action.Score)
	case truco.SAY_ME_VOY_AL_MAZO:
		what = fmt.Sprintf("%v me voy al mazo", said)
	default:
		what = "???"
	}

	return fmt.Sprintf("%v %v\n", who, what), nil
}

func (u *ui) startKeyEventLoop() {
	keyPressesCh := make(chan termbox.Event)
	go func() {
		for {
			event := termbox.PollEvent()
			if event.Type != termbox.EventKey {
				continue
			}
			if event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlD || event.Key == termbox.KeyCtrlZ || event.Ch == 'q' {
				termbox.Close()
				log.Println("Chau!")
				os.Exit(0)
			}
			keyPressesCh <- event
		}
	}()

	go func() {
		for {
			select {
			case <-keyPressesCh:
			case <-u.wantKeyPressCh:
				event := <-keyPressesCh
				u.sendKeyPressCh <- event.Ch
			}
		}
	}()
}

func (u *ui) pressAnyKey() {
	u.wantKeyPressCh <- struct{}{}
	<-u.sendKeyPressCh
}

func (u *ui) pressAnyNumber() int {
	u.wantKeyPressCh <- struct{}{}
	r := <-u.sendKeyPressCh
	num, err := strconv.Atoi(string(r))
	if err != nil {
		return u.pressAnyNumber()
	}
	return num
}

func spanishScore(score int) string {
	if score == 1 {
		return "1 mala"
	}
	if score < 15 {
		return fmt.Sprintf("%d malas", score)
	}
	if score == 15 {
		return "entraste"
	}
	return fmt.Sprintf("%d buenas", score-14)
}

func spanishAction(action string, state truco.GameState) string {
	switch action {
	case truco.REVEAL_CARD:
		return "tirar carta"
	case truco.SAY_ENVIDO:
		return "envido"
	case truco.SAY_REAL_ENVIDO:
		return "real envido"
	case truco.SAY_FALTA_ENVIDO:
		return "falta envido"
	case truco.SAY_ENVIDO_QUIERO:
		return "quiero"
	case truco.SAY_ENVIDO_NO_QUIERO:
		return "no quiero"
	case truco.SAY_TRUCO:
		return "truco"
	case truco.SAY_TRUCO_QUIERO:
		return "quiero"
	case truco.SAY_TRUCO_NO_QUIERO:
		return "no quiero"
	case truco.SAY_QUIERO_RETRUCO:
		return "quiero retruco"
	case truco.SAY_QUIERO_VALE_CUATRO:
		return "quiero vale cuatro"
	case truco.SAY_SON_BUENAS:
		return "son buenas"
	case truco.SAY_SON_MEJORES:
		score := state.Hands[state.TurnPlayerID].EnvidoScore()
		return fmt.Sprintf("%v son mejores", score)
	case truco.SAY_ME_VOY_AL_MAZO:
		return "me voy al mazo"
	default:
		return "???"
	}
}
