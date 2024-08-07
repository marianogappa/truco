//go:build !tinygo
// +build !tinygo

package exampleclient

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/marianogappa/truco/truco"
	"github.com/nsf/termbox-go"
)

type ui struct {
	keyCh chan rune
}

func NewUI() *ui {
	ui := &ui{}
	ui.keyCh = ui.startKeyEventLoop()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	return ui
}

func (u *ui) Close() {
	termbox.Close()
}

type renderMode int

const (
	PRINT_MODE_NORMAL renderMode = iota
	PRINT_MODE_SHOW_ROUND_RESULT
	PRINT_MODE_END
)

type renderState struct {
	mode            renderMode
	viewportWidth   int
	viewportHeight  int
	gs              truco.ClientGameState
	possibleActions []truco.Action
}

func calculateRenderState(state truco.ClientGameState) renderState {
	var (
		viewportWidth, viewportHeight = termbox.Size()
		possibleActions               = _deserializeActions(state.PossibleActions)
		gs                            = state
		mode                          = PRINT_MODE_NORMAL
	)
	if state.IsRoundFinished {
		mode = PRINT_MODE_SHOW_ROUND_RESULT
	}
	if state.IsGameEnded {
		mode = PRINT_MODE_END
	}

	return renderState{
		mode:            mode,
		gs:              gs,
		possibleActions: possibleActions,
		viewportWidth:   viewportWidth,
		viewportHeight:  viewportHeight,
	}
}

func (u *ui) render(state truco.ClientGameState) error {
	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}

	rs := calculateRenderState(state)

	renderScores(rs)
	renderTheirUnrevealedCards(rs)
	renderTheirRevealedCards(rs)
	renderLastAction(rs)
	renderEndSummary(rs)
	renderYourRevealedCards(rs)
	renderYourUnrevealedCards(rs)
	renderActions(rs)

	termbox.Flush()
	// This is an artificial delay to make the game more human-like.
	time.Sleep(1 * time.Second)

	return nil
}

func renderScores(rs renderState) {
	renderUpToAt(rs.viewportWidth-1, 0, fmt.Sprintf("Mano número %d", rs.gs.RoundNumber))

	youMano := ""
	themMano := ""
	if rs.gs.RoundTurnPlayerID == rs.gs.YouPlayerID {
		youMano = " (mano)"
	} else {
		themMano = " (mano)"
	}

	renderUpToAt(rs.viewportWidth-1, 1, fmt.Sprintf("Vos%v %v", youMano, spanishScore(rs.gs.YourScore)))
	renderUpToAt(rs.viewportWidth-1, 2, fmt.Sprintf("Elle%v %v", themMano, spanishScore(rs.gs.TheirScore)))
}

func renderTheirUnrevealedCards(rs renderState) {
	displayText := ""
	for _, card := range rs.gs.TheirDisplayUnrevealedCards {
		if card.IsHole {
			displayText += "  "
		} else {
			displayText += "[]"
		}
		displayText += " "
	}

	renderAt(0, 0, displayText)
}

func renderTheirRevealedCards(rs renderState) {
	renderAt(0, rs.viewportHeight/2-3, getCardsString(rs.gs.TheirRevealedCards))
}

func renderLastAction(rs renderState) {
	renderAt(0, rs.viewportHeight/2, getLastActionString(rs))
}

func renderEndSummary(rs renderState) {
	var renderText string

	switch rs.mode {
	case PRINT_MODE_SHOW_ROUND_RESULT:
		envidoPart := "el envido no se jugó"
		if rs.gs.EnvidoWinnerPlayerID != -1 {
			envidoWinner := "vos"
			won := "ganaste"
			if rs.gs.EnvidoWinnerPlayerID == rs.gs.ThemPlayerID {
				envidoWinner = "elle"
				won = "ganó"
			}
			envidoPart = fmt.Sprintf("%v %v %v puntos por el envido", envidoWinner, won, rs.gs.EnvidoPoints)
		}
		trucoWinner := "vos"
		won := "ganaste"
		if rs.gs.TrucoWinnerPlayerID == rs.gs.ThemPlayerID {
			trucoWinner = "elle"
			won = "ganó"
		}

		renderText = fmt.Sprintf(
			"Terminó la mano, %v y %v %v %v puntos por el truco.",
			envidoPart,
			trucoWinner,
			won,
			rs.gs.TrucoPoints,
		)
	case PRINT_MODE_END:
		var resultText string
		if rs.gs.YouPlayerID == rs.gs.WinnerPlayerID {
			resultText = "Ganaste 🥰"
		} else {
			resultText = "Perdiste 😭"
		}
		renderText = fmt.Sprintf("%v %v!", getLastActionString(rs), resultText)
	}

	renderAt(0, rs.viewportHeight/2, renderText)
}

func renderYourRevealedCards(rs renderState) {
	renderAt(0, rs.viewportHeight/2+3, getCardsString(rs.gs.YourRevealedCards))
}

func renderYourUnrevealedCards(rs renderState) {
	displayText := ""
	for _, card := range rs.gs.YourDisplayUnrevealedCards {
		if !card.IsHole {
			displayText += getDisplayCardString(card)
		} else {
			displayText += "     "
		}
		displayText += "  "
	}
	renderAt(0, rs.viewportHeight-4, displayText)
}

func renderActions(rs renderState) {
	var renderText string

	actionsString := ""
	for i, action := range rs.possibleActions {
		action := spanishAction(action)
		actionsString += fmt.Sprintf("%d. %s   ", i+1, action)
	}
	renderText = actionsString

	if len(rs.possibleActions) == 0 {
		renderText = "Esperando al otro jugador..."
	}

	if rs.mode == PRINT_MODE_END {
		renderText = "Presioná cualquier tecla para continuar..."
	}

	renderAt(0, rs.viewportHeight-2, renderText)
}

func renderAt(x, y int, s string) {
	_s := []rune(s)
	for i, r := range _s {
		termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

// Write so that the output ends at x, y
func renderUpToAt(x, y int, s string) {
	_s := []rune(s)
	for i, r := range _s {
		termbox.SetCell(x-len(_s)+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func getCardsString(cards []truco.Card) string {
	var cs []string
	for _, card := range cards {
		cs = append(cs, getCardString(card))
	}
	return strings.Join(cs, "  ")
}

func getCardString(card truco.Card) string {
	return fmt.Sprintf("[%v%v ]", card.Number, suitEmoji(card.Suit))
}
func getDisplayCardString(card truco.DisplayCard) string {
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

func _deserializeActions(as []json.RawMessage) []truco.Action {
	_as := []truco.Action{}
	for _, a := range as {
		_a, _ := truco.DeserializeAction(a)
		_as = append(_as, _a)
	}
	return _as
}
