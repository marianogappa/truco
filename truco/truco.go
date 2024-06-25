package truco

import (
	"encoding/json"
	"errors"
	"fmt"
)

// GameState represents the state of a Truco game.
//
// It is returned by the server on every single call, so if you want to implement a client,
// you need to be very familiar with this struct.
type GameState struct {
	// RoundTurnPlayerID is the player ID of the player who starts the round, or "mano".
	RoundTurnPlayerID int `json:"roundTurnPlayerID"`

	// RoundNumber is the number of the current round, starting from 1.
	RoundNumber int `json:"roundNumber"`

	// TurnPlayerID is the player ID of the player whose turn it is to play an action.
	// This is different from RoundTurnPlayerID, which is the player who starts the round.
	// They are the same at the beginning of the round.
	TurnPlayerID int `json:"turnPlayerID"`

	// Hands is a map of player IDs to their respective hands.
	Hands map[int]*Hand `json:"hands"`

	// Scores is a map of player IDs to their respective scores.
	// Scores go from 0 to 30.
	Scores map[int]int `json:"scores"`

	// PossibleActions is a list of possible actions that the current player can take.
	// Possible actions are calculated based on game state at the beginnin of the round and after
	// each action is run (i.e. GameState.RunAction).
	// The actions are strings, which are the names of the actions. In the case of REVEAL_CARD,
	// the card is not specified.
	PossibleActions []string `json:"possibleActionTypes"`

	// EnvidoSequence is the sequence of envido actions that have been taken in the current round.
	// Example sequence is: [SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO]
	// The player who started the sequence is saved too, so that certain "YieldsTurn" methods can work.
	EnvidoSequence *EnvidoSequence `json:"envidoSequence"`

	// TrucoSequence is the sequence of truco actions that have been taken in the current round.
	// Example sequence is: [SAY_TRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_RETRUCO, SAY_TRUCO_NO_QUIERO]
	TrucoSequence *TrucoSequence `json:"trucoSequence"`

	// CardRevealSequence is the sequence of card reveal actions that have been taken in the current round.
	// Each step is each card that was revealed (by both players).
	// `BistepWinners` (TODO: bad name) stores the result of the faceoff between each pair of cards.
	// A faceoff result will have the playerID of the winner, or -1 if it was a tie.
	CardRevealSequence *CardRevealSequence `json:"cardRevealSequence"`

	// EnvidoFinished is true if the envido sequence is finished, or can no longer be continued.
	// TODO: can we remove this? Looks redundant to other state. But need tests first.
	EnvidoFinished bool `json:"envidoFinished"`

	// EnvidoWinnerPlayerID is the player ID of the player who won the envido sequence.
	// TODO: looks like this is assigned to but never used. Can we remove?
	EnvidoWinnerPlayerID int `json:"envidoWinnerPlayerID"`

	// RoundFinished is true if the current round is finished. Each action's `Run()` method is responsible
	// for setting this. During `GameState.RunAction()`, If the action's `Run()` method sets this to true,
	// then `GameState.startNewRound()` will be called.
	//
	// Clients are not really notified of a round change, so they should keep track of the "last round
	// number" to see if it changes.
	RoundFinished bool `json:"roundFinished"`

	// IsEnded is true if the whole game is ended, rather than an individual round. This happens when
	// a player reaches 30 points.
	IsEnded bool `json:"isEnded"`

	// WinnerPlayerID is the player ID of the player who won the game. This is only set when `IsEnded` is
	// `true`. Otherwise, it's -1.
	WinnerPlayerID int `json:"winnerPlayerID"`

	// CurrentRoundResult contains the live results of the ongoing round for envido/truco winners & points.
	// It is set when actions are run, and is reset at the beginning of each round. Be careful when using
	// this in the client, because if the last action caused the round to finish, it will be reset before
	// you can use it to summarise what happened. You should use `RoundResults` in this case.
	CurrentRoundResult RoundResult `json:"currentRoundResult"`

	// RoundJustStarted is true if a round has just started (i.e. no actions have been run on this round).
	// This isn't used at all by the engine. It's strictly for clients to know, since there's no way to
	// relate actions to rounds. TODO: that's probably a problem?
	RoundJustStarted bool `json:"roundJustStarted"`

	// TrucoQuieroOwnerPlayerId is the player ID of the player who said "quiero" last in the truco
	// sequence. This is used to determine who can raise the stakes in the truco sequence.
	//
	// TODO: this should probably be inside TrucoSequence?
	TrucoQuieroOwnerPlayerId int `json:"trucoQuieroOwnerPlayerId"`

	// Actions is the list of actions that have been run in the game.
	// Each element is a JSON-serialized action. This is because `Action` is an interface, and we can't
	// serialize it directly otherwise. Clients should use `DeserializeAction` to get the actual action.
	//
	// TODO: rather than having "Actions", "HandsDealt", "RoundResults", "ActionOwnerPlayerIDs" as
	// separate fields, we should have a single "ActionLog" field that contains all of these.
	Actions []json.RawMessage `json:"actions"`

	// HandsDealt is the list of hands that were dealt in the game. Each element is a map of player IDs to
	// their respective hands.
	HandsDealt []map[int]*Hand `json:"handsDealt"`

	// RoundResults is the list of results of each round. Each element is a `RoundResult` struct.
	// This is useful for clients to see the results of each round, since `CurrentRoundResult` is reset
	// at the beginning of each round.
	RoundResults []RoundResult `json:"roundResults"`

	// ActionOwnerPlayerIDs is the list of player IDs who ran each action. There is no PlayerID field in
	// the `Action` struct, mostly because it would be annoying to have to distrust the client who sends
	// it.
	//
	// Note: this definitely should be inside an "ActionLog" slice, instead of here.
	ActionOwnerPlayerIDs []int `json:"actionOwnerPlayerIDs"`

	deck *deck `json:"-"`
}

func New(opts ...func(*GameState)) *GameState {
	// TODO: support taking player ids, ser/de, ...
	gs := &GameState{
		RoundTurnPlayerID: 1,
		RoundNumber:       0,
		Scores:            map[int]int{0: 0, 1: 0},
		Hands:             map[int]*Hand{0: nil, 1: nil},
		IsEnded:           false,
		WinnerPlayerID:    -1,
		Actions:           []json.RawMessage{},
		deck:              newDeck(),
	}

	for _, opt := range opts {
		opt(gs)
	}

	gs.startNewRound()

	return gs
}

func (g *GameState) startNewRound() {
	g.CurrentRoundResult = RoundResult{
		EnvidoWinnerPlayerID: -1,
		EnvidoPoints:         0,
		TrucoWinnerPlayerID:  -1,
		TrucoPoints:          0,
	}

	g.RoundJustStarted = true
	g.RoundTurnPlayerID = g.OpponentOf(g.RoundTurnPlayerID)
	g.RoundNumber++
	g.TurnPlayerID = g.RoundTurnPlayerID

	handPlayer0 := g.deck.dealHand()
	handPlayer1 := g.deck.dealHand()
	g.HandsDealt = append(g.HandsDealt, map[int]*Hand{
		g.RoundTurnPlayerID:               handPlayer0,
		g.OpponentOf(g.RoundTurnPlayerID): handPlayer1,
	})
	g.Hands = map[int]*Hand{
		g.RoundTurnPlayerID:               handPlayer0,
		g.OpponentOf(g.RoundTurnPlayerID): handPlayer1,
	}
	g.EnvidoWinnerPlayerID = -1
	g.EnvidoSequence = &EnvidoSequence{StartingPlayerID: -1}
	g.TrucoSequence = &TrucoSequence{}
	g.CardRevealSequence = &CardRevealSequence{}
	g.EnvidoFinished = false
	g.RoundFinished = false
	g.TrucoQuieroOwnerPlayerId = -1
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
		g.RoundResults = append(g.RoundResults, g.CurrentRoundResult)
		g.startNewRound()
		return nil
	}

	// Switch player turn within current round (unless current action doesn't yield turn)
	if !g.IsEnded && !g.RoundFinished && action.YieldsTurn(*g) {
		g.TurnPlayerID = g.OpponentOf(g.TurnPlayerID)
	}

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

	g.PossibleActions = g.CalculatePossibleActions()
	return nil
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
	case REVEAL_CARD:
		action = &ActionRevealCard{}
	case SAY_ENVIDO:
		action = &ActionSayEnvido{}
	case SAY_REAL_ENVIDO:
		action = &ActionSayRealEnvido{}
	case SAY_FALTA_ENVIDO:
		action = &ActionSayFaltaEnvido{}
	case SAY_ENVIDO_QUIERO:
		action = &ActionSayEnvidoQuiero{}
	case SAY_ENVIDO_NO_QUIERO:
		action = &ActionSayEnvidoNoQuiero{}
	case SAY_TRUCO:
		action = &ActionSayTruco{}
	case SAY_TRUCO_QUIERO:
		action = &ActionSayTrucoQuiero{}
	case SAY_TRUCO_NO_QUIERO:
		action = &ActionSayTrucoNoQuiero{}
	case SAY_QUIERO_RETRUCO:
		action = &ActionSayQuieroRetruco{}
	case SAY_QUIERO_VALE_CUATRO:
		action = &ActionSayQuieroValeCuatro{}
	case SAY_SON_BUENAS:
		action = &ActionSaySonBuenas{}
	case SAY_SON_MEJORES:
		action = &ActionSaySonMejores{}
	case SAY_ME_VOY_AL_MAZO:
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
	EnvidoWinnerPlayerID int `json:"envidoWinnerPlayerID"`
	EnvidoPoints         int `json:"envidoPoints"`
	TrucoWinnerPlayerID  int `json:"trucoWinnerPlayerID"`
	TrucoPoints          int `json:"trucoPoints"`
}
