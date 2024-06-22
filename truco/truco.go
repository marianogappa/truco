package truco

import (
	"encoding/json"
	"errors"
	"fmt"
)

type GameState struct {
	RoundTurnPlayerID int           `json:"roundTurnPlayerID"`
	RoundNumber       int           `json:"roundNumber"`
	TurnPlayerID      int           `json:"turnPlayerID"`
	Hands             map[int]*Hand `json:"hands"`
	Scores            map[int]int   `json:"scores"`

	PossibleActions      []string            `json:"possibleActionTypes"`
	EnvidoSequence       *EnvidoSequence     `json:"envidoSequence"`
	TrucoSequence        *TrucoSequence      `json:"trucoSequence"`
	CardRevealSequence   *CardRevealSequence `json:"cardRevealSequence"`
	EnvidoFinished       bool                `json:"envidoFinished"`
	EnvidoWinnerPlayerID int                 `json:"envidoWinnerPlayerID"`
	ValidSonBuenas       bool                `json:"validSonBuenas"`
	ValidSonMejores      bool                `json:"validSonMejores"`
	RoundFinished        bool                `json:"roundFinished"`
	IsEnded              bool                `json:"isEnded"`
	WinnerPlayerID       int                 `json:"winnerPlayerID"`
	CurrentRoundResult   RoundResult         `json:"currentRoundResult"`
	RoundJustStarted     bool                `json:"roundJustStarted"`

	Actions              []json.RawMessage `json:"actions"`
	HandsDealt           []map[int]*Hand   `json:"handsDealt"`
	RoundResults         []RoundResult     `json:"roundResults"`
	ActionOwnerPlayerIDs []int             `json:"actionOwnerPlayerIDs"`
}

func New() *GameState {
	// TODO: support taking player ids, ser/de, ...
	gs := &GameState{
		RoundTurnPlayerID: 1,
		RoundNumber:       0,
		Scores:            map[int]int{0: 0, 1: 0},
		Hands:             map[int]*Hand{0: nil, 1: nil},
		IsEnded:           false,
		WinnerPlayerID:    -1,
		Actions:           []json.RawMessage{},
	}

	gs.StartNewRound()

	return gs
}

func (g *GameState) StartNewRound() {
	deck := NewDeck()
	g.CurrentRoundResult = RoundResult{
		EnvidoWinnerPlayerID: -1,
		EnvidoPoints:         0,
		TrucoWinnerPlayerID:  -1,
		TrucoPoints:          0,
		LastAction:           nil,
	}

	g.RoundJustStarted = true
	g.RoundTurnPlayerID = g.OpponentOf(g.RoundTurnPlayerID)
	g.RoundNumber++
	g.TurnPlayerID = g.RoundTurnPlayerID

	// By default, deal new hands, but if GameState has hands saved, use those
	// TODO: this changes if players are not 0 & 1 (or more than 2 players)
	handPlayer0 := deck.DealHand()
	handPlayer1 := deck.DealHand()
	if len(g.HandsDealt) >= g.RoundNumber {
		handPlayer0 = g.HandsDealt[g.RoundNumber-1][0]
		handPlayer1 = g.HandsDealt[g.RoundNumber-1][1]
	} else {
		g.HandsDealt = append(g.HandsDealt, map[int]*Hand{
			0: handPlayer0,
			1: handPlayer1,
		})
	}

	g.Hands = map[int]*Hand{
		0: handPlayer0,
		1: handPlayer1,
	}
	g.EnvidoWinnerPlayerID = -1
	g.EnvidoSequence = &EnvidoSequence{StartingPlayerID: -1}
	g.TrucoSequence = &TrucoSequence{}
	g.CardRevealSequence = &CardRevealSequence{}
	g.EnvidoFinished = false
	g.ValidSonBuenas = true
	g.ValidSonMejores = true
	g.RoundFinished = false
	g.PossibleActions = g.CalculatePossibleActions()
}

func (g *GameState) RunAction(action Action) error {
	if g.IsEnded {
		return errGameIsEnded
	}

	if !action.IsPossible(*g) {
		return errActionNotPossible
	}
	err := action.Run(g)
	if err != nil {
		return err
	}
	g.RoundJustStarted = false
	bs := SerializeAction(action)
	g.Actions = append(g.Actions, bs)
	g.ActionOwnerPlayerIDs = append(g.ActionOwnerPlayerIDs, g.CurrentPlayerID())

	// Start new round if current round is finished
	if !g.IsEnded && g.RoundFinished {
		g.HandleInvalidEnvidoDeclarations()
		// TODO: here we need to handle the case where players lie about envido score
		g.RoundResults = append(g.RoundResults, g.CurrentRoundResult)
		g.StartNewRound()
		return nil
	}

	// Switch player turn within current round (unless current action doesn't yield turn)
	if !g.IsEnded && !g.RoundFinished && action.YieldsTurn(*g) {
		g.TurnPlayerID = g.OpponentOf(g.TurnPlayerID)
	}
	g.PossibleActions = g.CalculatePossibleActions()

	// Handle end of game due to score
	// TODO: this changes if players are not 0 & 1 (or more than 2 players)
	if g.Scores[0] >= 30 || g.Scores[1] >= 30 {
		g.IsEnded = true
		g.Scores[0] = min(30, g.Scores[0])
		g.Scores[1] = min(30, g.Scores[1])
		g.WinnerPlayerID = 0
		if g.Scores[1] > g.Scores[0] {
			g.WinnerPlayerID = 1
		}
	}

	return nil
}

func (g *GameState) HandleInvalidEnvidoDeclarations() {
	// TODO this is really tricky actually
}

func (g GameState) CurrentPlayerID() int {
	return g.TurnPlayerID
}

func (g GameState) CurrentPlayerScore() int {
	return g.Scores[g.TurnPlayerID]
}

func (g GameState) OpponentPlayerID() int {
	return g.OpponentOf(g.CurrentPlayerID())
}

func (g GameState) OpponentPlayerScore() int {
	return g.Scores[g.OpponentPlayerID()]
}

func (g GameState) RoundTurnOpponentPlayerID() int {
	return g.OpponentOf(g.RoundTurnPlayerID)
}

func (g GameState) OpponentOf(playerID int) int {
	// N.B. this function doesn't check if playerID is valid
	players := []int{}
	for playerID := range g.Hands {
		players = append(players, playerID)
	}
	// TODO: This still assumes 2 players, but at least it doesn't hardcode the player IDs
	for _, p := range players {
		if p != playerID {
			return p
		}
	}
	return -1 // Unreachable
}

func (g GameState) Serialize() ([]byte, error) {
	return json.Marshal(g)
}

func (g *GameState) PrettyPrint() (string, error) {
	var prettyJSON []byte
	prettyJSON, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

type Action interface {
	IsPossible(g GameState) bool
	Run(g *GameState) error
	GetName() string
	YieldsTurn(g GameState) bool
}

var (
	// errUnknownActionType = errors.New("unknown action type")
	errActionNotPossible = errors.New("action not possible")
	errEnvidoFinished    = errors.New("envido finished")
	errGameIsEnded       = errors.New("game is ended")
)

func (g GameState) CalculatePossibleActions() []string {
	allActions := []Action{
		ActionRevealCard{act: act{Name: "reveal_card"}},
		ActionSayEnvido{act: act{Name: "say_envido"}},
		ActionSayRealEnvido{act: act{Name: "say_real_envido"}},
		ActionSayFaltaEnvido{act: act{Name: "say_falta_envido"}},
		ActionSayEnvidoQuiero{act: act{Name: "say_envido_quiero"}},
		ActionSayEnvidoNoQuiero{act: act{Name: "say_envido_no_quiero"}},
		ActionSayTruco{act: act{Name: "say_truco"}},
		ActionSayTrucoQuiero{act: act{Name: "say_truco_quiero"}},
		ActionSayTrucoNoQuiero{act: act{Name: "say_truco_no_quiero"}},
		ActionSayQuieroRetruco{act: act{Name: "say_quiero_retruco"}},
		ActionSayQuieroValeCuatro{act: act{Name: "say_quiero_vale_cuatro"}},
		ActionSaySonBuenas{act: act{Name: "say_son_buenas"}},
		ActionSaySonMejores{act: act{Name: "say_son_mejores"}},
		ActionSayMeVoyAlMazo{act: act{Name: "say_me_voy_al_mazo"}},
	}
	actions := []string{}
	for _, action := range allActions {
		if action.IsPossible(g) {
			actions = append(actions, action.GetName())
		}
	}
	return actions
}

func SerializeAction(action Action) []byte {
	bs, _ := json.Marshal(action)
	return bs
}

func DeserializeAction(bs []byte) (Action, error) {
	var actionName struct {
		Name string `json:"name"`
	}

	err := json.Unmarshal(bs, &actionName)
	if err != nil {
		return nil, err
	}

	var action Action
	switch actionName.Name {
	case "reveal_card":
		action = &ActionRevealCard{}
	case "say_envido":
		action = &ActionSayEnvido{}
	case "say_real_envido":
		action = &ActionSayRealEnvido{}
	case "say_falta_envido":
		action = &ActionSayFaltaEnvido{}
	case "say_envido_quiero":
		action = &ActionSayEnvidoQuiero{}
	case "say_envido_no_quiero":
		action = &ActionSayEnvidoNoQuiero{}
	case "say_truco":
		action = &ActionSayTruco{}
	case "say_truco_quiero":
		action = &ActionSayTrucoQuiero{}
	case "say_truco_no_quiero":
		action = &ActionSayTrucoNoQuiero{}
	case "say_quiero_retruco":
		action = &ActionSayQuieroRetruco{}
	case "say_quiero_vale_cuatro":
		action = &ActionSayQuieroValeCuatro{}
	case "say_son_buenas":
		action = &ActionSaySonBuenas{}
	case "say_son_mejores":
		action = &ActionSaySonMejores{}
	case "say_me_voy_al_mazo":
		action = &ActionSayMeVoyAlMazo{}
	default:
		return nil, fmt.Errorf("unknown action type %v", actionName.Name)
	}

	err = json.Unmarshal(bs, action)
	if err != nil {
		return nil, err
	}

	return action, nil
}

type RoundResult struct {
	EnvidoWinnerPlayerID int             `json:"envidoWinnerPlayerID"`
	EnvidoPoints         int             `json:"envidoPoints"`
	TrucoWinnerPlayerID  int             `json:"trucoWinnerPlayerID"`
	TrucoPoints          int             `json:"trucoPoints"`
	LastAction           json.RawMessage `json:"lastAction"`
}
