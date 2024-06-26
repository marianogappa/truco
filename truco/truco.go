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

	// TurnOpponentPlayerID is the player ID of the opponent of the player whose turn it is.
	TurnOpponentPlayerID int `json:"turnOpponentPlayerID"`

	Players map[int]*Player `json:"players"`

	// // Hands is a map of player IDs to their respective hands.
	// Hands map[int]*Hand `json:"hands"`

	// // Scores is a map of player IDs to their respective scores.
	// // Scores go from 0 to 30.
	// Scores map[int]int `json:"scores"`

	// PossibleActions is a list of possible actions that the current player can take.
	// Possible actions are calculated based on game state at the beginnin of the round and after
	// each action is run (i.e. GameState.RunAction).
	// The actions are strings, which are the names of the actions. In the case of REVEAL_CARD,
	// the card is not specified.
	PossibleActions []json.RawMessage `json:"possibleActions"`

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

	// IsEnvidoFinished is true if the envido sequence is finished, or can no longer be continued.
	// TODO: can we remove this? Looks redundant to other state. But need tests first.
	IsEnvidoFinished bool `json:"isEnvidoFinished"`

	// IsRoundFinished is true if the current round is finished. Each action's `Run()` method is responsible
	// for setting this. During `GameState.RunAction()`, If the action's `Run()` method sets this to true,
	// then `GameState.startNewRound()` will be called.
	//
	// Clients are not really notified of a round change, so they should keep track of the "last round
	// number" to see if it changes.
	IsRoundFinished bool `json:"isRoundFinished"`

	// IsGameEnded is true if the whole game is ended, rather than an individual round. This happens when
	// a player reaches 30 points.
	IsGameEnded bool `json:"isGameEnded"`

	// IsRoundJustStarted is true if a round has just started (i.e. no actions have been run on this round).
	// This isn't used at all by the engine. It's strictly for clients to know, since there's no way to
	// relate actions to rounds. TODO: that's probably a problem?
	IsRoundJustStarted bool `json:"isRoundJustStarted"`

	// WinnerPlayerID is the player ID of the player who won the game. This is only set when `IsEnded` is
	// `true`. Otherwise, it's -1.
	WinnerPlayerID int `json:"winnerPlayerID"`

	// CurrentRoundResult contains the live results of the ongoing round for envido/truco winners & points.
	// It is set when actions are run, and is reset at the beginning of each round. Be careful when using
	// this in the client, because if the last action caused the round to finish, it will be reset before
	// you can use it to summarise what happened. You should use `RoundResults` in this case.
	CurrentRoundResult RoundResult `json:"currentRoundResult"`

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

type Player struct {
	// Hands contains the revealed and unrevealed cards of the player.
	Hand *Hand `json:"hand"`

	// Score is the player's scores (from 0 to 30).
	Score int `json:"score"`
}

func New(opts ...func(*GameState)) *GameState {
	// TODO: support taking player ids, ser/de, ...
	gs := &GameState{
		RoundTurnPlayerID: 1,
		RoundNumber:       0,
		Players: map[int]*Player{
			0: {Hand: nil, Score: 0},
			1: {Hand: nil, Score: 0},
		},
		IsGameEnded:    false,
		WinnerPlayerID: -1,
		Actions:        []json.RawMessage{},
		deck:           newDeck(),
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

	g.IsRoundJustStarted = true
	g.RoundTurnPlayerID = g.OpponentOf(g.RoundTurnPlayerID)
	g.RoundNumber++
	g.TurnPlayerID = g.RoundTurnPlayerID
	g.TurnOpponentPlayerID = g.OpponentOf(g.TurnPlayerID)
	g.Players[g.TurnPlayerID].Hand = g.deck.dealHand()
	g.Players[g.TurnOpponentPlayerID].Hand = g.deck.dealHand()
	g.HandsDealt = append(g.HandsDealt, map[int]*Hand{
		g.TurnPlayerID:         g.Players[g.TurnPlayerID].Hand,
		g.TurnOpponentPlayerID: g.Players[g.TurnOpponentPlayerID].Hand,
	})
	g.EnvidoSequence = &EnvidoSequence{StartingPlayerID: -1}
	g.TrucoSequence = &TrucoSequence{StartingPlayerID: -1}
	g.CardRevealSequence = &CardRevealSequence{}
	g.IsEnvidoFinished = false
	g.IsRoundFinished = false
	g.TrucoQuieroOwnerPlayerId = -1
	g.PossibleActions = _serializeActions(g.CalculatePossibleActions())
}

func (g *GameState) RunAction(action Action) error {
	if g.IsGameEnded {
		return errGameIsEnded
	}

	if !action.IsPossible(*g) {
		return errActionNotPossible
	}
	err := action.Run(g)
	if err != nil {
		return err
	}
	g.IsRoundJustStarted = false
	bs := SerializeAction(action)
	g.Actions = append(g.Actions, bs)
	g.ActionOwnerPlayerIDs = append(g.ActionOwnerPlayerIDs, g.TurnPlayerID)

	// Start new round if current round is finished
	if !g.IsGameEnded && g.IsRoundFinished {
		g.RoundResults = append(g.RoundResults, g.CurrentRoundResult)
		g.startNewRound()
		return nil
	}

	// Switch player turn within current round (unless current action doesn't yield turn)
	if !g.IsGameEnded && !g.IsRoundFinished && action.YieldsTurn(*g) {
		g.TurnPlayerID, g.TurnOpponentPlayerID = g.TurnOpponentPlayerID, g.TurnPlayerID
	}

	// Handle end of game due to score
	for playerID := range g.Players {
		if g.Players[playerID].Score >= 30 {
			g.Players[playerID].Score = 30
			g.IsGameEnded = true
			g.WinnerPlayerID = playerID
		}
	}

	g.PossibleActions = _serializeActions(g.CalculatePossibleActions())
	return nil
}

func (g GameState) OpponentOf(playerID int) int {
	for id := range g.Players {
		if id != playerID {
			return id
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

func (g GameState) CalculatePossibleActions() []Action {
	envidoScore := g.Players[g.TurnPlayerID].Hand.EnvidoScore()

	allActions := []Action{}

	for _, card := range g.Players[g.TurnPlayerID].Hand.Unrevealed {
		allActions = append(allActions, newActionRevealCard(card))
	}

	allActions = append(allActions,
		newActionSayEnvido(),
		newActionSayRealEnvido(),
		newActionSayFaltaEnvido(),
		newActionSayEnvidoQuiero(envidoScore),
		newActionSayEnvidoNoQuiero(),
		newActionSayTruco(),
		newActionSayTrucoQuiero(),
		newActionSayTrucoNoQuiero(),
		newActionSayQuieroRetruco(),
		newActionSayQuieroValeCuatro(),
		newActionSaySonBuenas(),
		newActionSaySonMejores(envidoScore),
		newActionSayMeVoyAlMazo(),
	)

	possibleActions := []Action{}
	for _, action := range allActions {
		if action.IsPossible(g) {
			possibleActions = append(possibleActions, action)
		}
	}
	return possibleActions
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

func _serializeActions(as []Action) []json.RawMessage {
	_as := []json.RawMessage{}
	for _, a := range as {
		_as = append(_as, json.RawMessage(SerializeAction(a)))
	}
	return _as
}
