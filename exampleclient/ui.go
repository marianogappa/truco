package exampleclient

import (
	"encoding/json"
	"fmt"
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

func calculateRenderState(state truco.ClientGameState, mode renderMode) renderState {
	var (
		viewportWidth, viewportHeight = termbox.Size()
		possibleActions               = _deserializeActions(state.PossibleActions)
		gs                            = state
	)

	if mode == PRINT_MODE_SHOW_ROUND_RESULT {
		// Note that RoundNumber has already been incremented
		// so we need to get the previous round's hands.
		gs.YourRevealedCards = gs.LastRoundLog.HandsDealt[gs.YouPlayerID].Revealed
		gs.YourUnrevealedCards = gs.LastRoundLog.HandsDealt[gs.YouPlayerID].Unrevealed
		gs.TheirRevealedCards = gs.LastRoundLog.HandsDealt[gs.ThemPlayerID].Revealed
	}

	return renderState{
		mode:            mode,
		gs:              gs,
		possibleActions: possibleActions,
		viewportWidth:   viewportWidth,
		viewportHeight:  viewportHeight,
	}
}

func (u *ui) render(state truco.ClientGameState, mode renderMode) error {
	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}

	rs := calculateRenderState(state, mode)

	renderScores(rs)
	renderTheirUnrevealedCards(rs)
	renderTheirRevealedCards(rs)
	renderLastAction(rs)
	renderEndSummary(rs)
	renderYourRevealedCards(rs)
	renderYourUnrevealedCards(rs)
	renderActions(rs)

	termbox.Flush()
	return nil
}

func renderScores(rs renderState) {
	renderUpToAt(rs.viewportWidth-1, 0, fmt.Sprintf("Mano número %d", rs.gs.RoundNumber))

	youMano := ""
	themMano := ""
	if rs.gs.TurnPlayerID == rs.gs.YouPlayerID {
		youMano = " (mano)"
	} else {
		themMano = " (mano)"
	}

	renderUpToAt(rs.viewportWidth-1, 1, fmt.Sprintf("Vos%v %v", youMano, spanishScore(rs.gs.YourScore)))
	renderUpToAt(rs.viewportWidth-1, 2, fmt.Sprintf("Elle%v %v", themMano, spanishScore(rs.gs.TheirScore)))
}

func renderTheirUnrevealedCards(rs renderState) {
	renderAt(0, 0, strings.Repeat("[] ", rs.gs.TheirUnrevealedCardLength))
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
		if rs.gs.LastRoundLog.EnvidoWinnerPlayerID != -1 {
			envidoWinner := "vos"
			won := "ganaste"
			if rs.gs.LastRoundLog.EnvidoWinnerPlayerID == rs.gs.ThemPlayerID {
				envidoWinner = "elle"
				won = "ganó"
			}
			envidoPart = fmt.Sprintf("%v %v %v puntos por el envido", envidoWinner, won, rs.gs.LastRoundLog.EnvidoPoints)
		}
		trucoWinner := "vos"
		won := "ganaste"
		if rs.gs.LastRoundLog.TrucoWinnerPlayerID == rs.gs.ThemPlayerID {
			trucoWinner = "elle"
			won = "ganó"
		}

		renderText = fmt.Sprintf(
			"Terminó la mano, %v y %v %v %v puntos por el truco.",
			envidoPart,
			trucoWinner,
			won,
			rs.gs.LastRoundLog.TrucoPoints,
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
	renderAt(0, rs.viewportHeight-4, getCardsString(rs.gs.YourUnrevealedCards))
}

func renderActions(rs renderState) {
	var renderText string

	switch rs.mode {
	case PRINT_MODE_SHOW_ROUND_RESULT, PRINT_MODE_END:
		renderText = "Presioná cualquier tecla para continuar..."
	default:
		if rs.gs.TurnPlayerID == rs.gs.YouPlayerID {
			actionsString := ""
			for i, action := range rs.possibleActions {
				action := spanishAction(action)
				actionsString += fmt.Sprintf("%d. %s   ", i+1, action)
			}
			renderText = actionsString
		} else {
			renderText = "Esperando al otro jugador..."
		}

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
