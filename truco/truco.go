package truco

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// MaxPoints is the points a player must reach to win the game.
// It is set as a const in case support for 15 points games is needed in the future.
const MaxPoints = 30

// GameState represents the state of a Truco game. It is the central struct to this package.
//
// If you want to implement a client, you should look at ClientGameState instead.
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

	// Players is a map of player IDs to their respective hands and scores.
	// There are 2 players in a game. Use TurnPlayerID and TurnOpponentPlayerID to index
	// into this map, or iterate over it to discover player ids.
	Players map[int]*Player `json:"players"`

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
	// a player reaches MaxPoints points.
	IsGameEnded bool `json:"isGameEnded"`

	// WinnerPlayerID is the player ID of the player who won the game. This is only set when `IsGameEnded` is
	// `true`. Otherwise, it's -1.
	WinnerPlayerID int `json:"winnerPlayerID"`

	// RoundsLog is the ordered list of logs of each round that was played in the game.
	//
	// Use GameState.RoundNumber to index into this list (note thus that it's 1-indexed).
	// This means that there is an empty round at the beginning of the list.
	//
	// Note that there is a "live entry" for the current round. This entry will always have
	// the HandsDealt, but the other fields depend on the stage the round is in. You can
	// use the ActionsLog's length to determine if a round has just started.
	//
	// Note that, if the last action ran caused a round to finish, if you want to render
	// the screen with the last action, you have to check the last action of the previous round instead!
	RoundsLog []*RoundLog `json:"actionLog"`

	RoundFinishedConfirmedPlayerIDs map[int]bool `json:"roundFinishedConfirmedPlayerIDs"`

	deck *deck `json:"-"`
}

type Player struct {
	// Hands contains the revealed and unrevealed cards of the player.
	Hand *Hand `json:"hand"`

	// Score is the player's scores (from 0 to MaxPoints).
	Score int `json:"score"`
}

// RoundLog is a log of a round that was played in the game
type RoundLog struct {
	// HandsDealt is a map from PlayerID to its hand during this round.
	HandsDealt map[int]*Hand `json:"handsDealt"`

	// For envido/truco winners and points, note that there is still a
	// winner of 1 point if a player said "no quiero" to the envido/truco.
	//
	// If envido wasn't played at all, then EnvidoWinnerPlayerID is -1.
	//
	// At the end of a round, there will always be a TrucoWinnerPlayerID,
	// even if truco wasn't played, implicitly by revealing the cards.

	EnvidoWinnerPlayerID int `json:"envidoWinnerPlayerID"`
	EnvidoPoints         int `json:"envidoPoints"`
	TrucoWinnerPlayerID  int `json:"trucoWinnerPlayerID"`
	TrucoPoints          int `json:"trucoPoints"`

	// ActionsLog is the ordered list of actions of this round.
	ActionsLog []ActionLog `json:"actionsLog"`
}

// ActionLog is a log of an action that was run in a round.
type ActionLog struct {
	// PlayerID is the player ID of the player who ran the action.
	PlayerID int `json:"playerID"`

	// Action is a JSON-serialized action. This is because `Action` is an interface, and we can't
	// serialize it directly otherwise. Clients should use `truco.DeserializeAction`.`
	Action json.RawMessage `json:"action"`
}

func New(opts ...func(*GameState)) *GameState {
	gs := &GameState{
		RoundTurnPlayerID: 1,
		RoundNumber:       0,
		Players: map[int]*Player{
			0: {Hand: nil, Score: 0},
			1: {Hand: nil, Score: 0},
		},
		IsGameEnded:    false,
		WinnerPlayerID: -1,
		RoundsLog:      []*RoundLog{{}}, // initialised with an empty round to be 1-indexed
		deck:           newDeck(),
	}

	for _, opt := range opts {
		opt(gs)
	}

	gs.startNewRound()

	return gs
}

func (g *GameState) startNewRound() {
	g.deck.shuffle()
	g.RoundTurnPlayerID = g.OpponentOf(g.RoundTurnPlayerID)
	g.RoundNumber++
	g.TurnPlayerID = g.RoundTurnPlayerID
	g.TurnOpponentPlayerID = g.OpponentOf(g.TurnPlayerID)
	g.Players[g.TurnPlayerID].Hand = g.deck.dealHand()
	g.Players[g.TurnOpponentPlayerID].Hand = g.deck.dealHand()
	g.EnvidoSequence = &EnvidoSequence{StartingPlayerID: -1}
	g.TrucoSequence = &TrucoSequence{StartingPlayerID: -1, QuieroOwnerPlayerID: -1}
	g.CardRevealSequence = &CardRevealSequence{}
	g.IsEnvidoFinished = false
	g.IsRoundFinished = false
	g.RoundFinishedConfirmedPlayerIDs = map[int]bool{}
	g.RoundsLog = append(g.RoundsLog, &RoundLog{
		HandsDealt: map[int]*Hand{
			g.TurnPlayerID:         g.Players[g.TurnPlayerID].Hand,
			g.TurnOpponentPlayerID: g.Players[g.TurnOpponentPlayerID].Hand,
		},
		EnvidoWinnerPlayerID: -1,
		EnvidoPoints:         0,
		TrucoWinnerPlayerID:  -1,
		TrucoPoints:          0,
		ActionsLog:           []ActionLog{},
	})
	g.PossibleActions = _serializeActions(g.CalculatePossibleActions())
}

func (g *GameState) RunAction(action Action) error {
	if action == nil {
		return nil
	}

	if g.IsGameEnded {
		return fmt.Errorf("%w trying to run [%v]", errGameIsEnded, action)
	}

	if !g.IsRoundFinished && action.GetPlayerID() != g.TurnPlayerID {
		return errNotYourTurn
	}

	if !action.IsPossible(*g) {
		return errActionNotPossible
	}
	err := action.Run(g)
	if err != nil {
		return err
	}

	if action.GetName() != CONFIRM_ROUND_FINISHED {
		g.RoundsLog[g.RoundNumber].ActionsLog = append(g.RoundsLog[g.RoundNumber].ActionsLog, ActionLog{
			PlayerID: g.TurnPlayerID,
			Action:   SerializeAction(action),
		})
	}

	// Start new round if current round is finished
	if !g.IsGameEnded && g.IsRoundFinished && len(g.RoundFinishedConfirmedPlayerIDs) == 2 {
		fmt.Println("Starting new round...")
		g.startNewRound()
		return nil
	}

	// Switch player turn within current round (unless current action doesn't yield turn)
	if !g.IsGameEnded && !g.IsRoundFinished && action.YieldsTurn(*g) {
		g.TurnPlayerID, g.TurnOpponentPlayerID = g.TurnOpponentPlayerID, g.TurnPlayerID
	}

	if !g.IsGameEnded && g.IsRoundFinished && len(g.RoundFinishedConfirmedPlayerIDs) == 1 {
		if g.RoundFinishedConfirmedPlayerIDs[g.TurnPlayerID] {
			g.TurnPlayerID, g.TurnOpponentPlayerID = g.TurnOpponentPlayerID, g.TurnPlayerID
		}
	}

	// Handle end of game due to score
	for playerID := range g.Players {
		if g.Players[playerID].Score >= MaxPoints {
			g.Players[playerID].Score = MaxPoints
			g.IsGameEnded = true
			g.WinnerPlayerID = playerID
		}
	}

	possibleActions := g.CalculatePossibleActions()
	g.PossibleActions = _serializeActions(possibleActions)

	log.Printf("Possible actions: %v\n", possibleActions)

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

func (g *GameState) canAwardEnvidoPoints(revealedHand Hand) bool {
	wonBy := g.RoundsLog[g.RoundNumber].EnvidoWinnerPlayerID
	if wonBy == -1 {
		return false
	}
	if !g.EnvidoSequence.WasAccepted() {
		return false
	}
	if g.EnvidoSequence.EnvidoPointsAwarded {
		return false
	}
	if revealedHand.EnvidoScore() != g.Players[wonBy].Hand.EnvidoScore() {
		return false
	}
	return true
}

func (g *GameState) tryAwardEnvidoPoints(playerID int) bool {
	if !g.canAwardEnvidoPoints(Hand{Revealed: g.Players[playerID].Hand.Revealed}) {
		return false
	}
	wonBy := g.RoundsLog[g.RoundNumber].EnvidoWinnerPlayerID
	score := g.RoundsLog[g.RoundNumber].EnvidoPoints
	g.Players[wonBy].Score += score
	g.EnvidoSequence.EnvidoPointsAwarded = true
	return true
}

type Action interface {
	IsPossible(g GameState) bool
	Run(g *GameState) error
	GetName() string
	GetPlayerID() int
	YieldsTurn(g GameState) bool
	// Some actions need to be enriched with additional information.
	// e.g. a say_truco_quiero action is enriched with "RequiresReminder".
	// GameState.CalculatePossibleActions() must call this method on all actions.
	Enrich(g GameState)

	fmt.Stringer
}

var (
	errActionNotPossible = errors.New("action not possible")
	errEnvidoFinished    = errors.New("envido finished")
	errGameIsEnded       = errors.New("game is ended")
	errNotYourTurn       = errors.New("not your turn")
)

func (g GameState) CalculatePossibleActions() []Action {
	allActions := []Action{}
	allActions = append(allActions, NewActionsRevealCards(g.TurnPlayerID, g)...)
	allActions = append(allActions,
		NewActionSayEnvido(g.TurnPlayerID),
		NewActionSayRealEnvido(g.TurnPlayerID),
		NewActionSayFaltaEnvido(g.TurnPlayerID),
		NewActionSayEnvidoQuiero(g.TurnPlayerID),
		NewActionSayEnvidoScore(g.TurnPlayerID),
		NewActionSayEnvidoNoQuiero(g.TurnPlayerID),
		NewActionSayTruco(g.TurnPlayerID),
		NewActionSayTrucoQuiero(g.TurnPlayerID),
		NewActionSayTrucoNoQuiero(g.TurnPlayerID),
		NewActionSayQuieroRetruco(g.TurnPlayerID),
		NewActionSayQuieroValeCuatro(g.TurnPlayerID),
		NewActionSaySonBuenas(g.TurnPlayerID),
		NewActionSaySonMejores(g.TurnPlayerID),
		NewActionSayMeVoyAlMazo(g.TurnPlayerID),
		NewActionConfirmRoundFinished(g.TurnPlayerID),
		NewActionConfirmRoundFinished(g.TurnOpponentPlayerID),
		NewActionRevealEnvidoScore(g.TurnPlayerID),
	)

	// The reveal_envido_score action happens in two cases:
	// 1. Round has ended and player hasn't shown envido score yet
	// 2. Round is going on, but player would win the game by revealing the score
	//
	// In both cases, this should be the only action available. So don't check others.
	actionRevealEnvidoScore := NewActionRevealEnvidoScore(g.TurnPlayerID)
	actionRevealEnvidoScore.Enrich(g)
	if actionRevealEnvidoScore.IsPossible(g) {
		allActions = []Action{actionRevealEnvidoScore}
	}

	possibleActions := []Action{}
	for _, action := range allActions {
		action.Enrich(g)
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
	case SAY_ENVIDO_SCORE:
		action = &ActionSayEnvidoScore{}
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
	case CONFIRM_ROUND_FINISHED:
		action = &ActionConfirmRoundFinished{}
	case REVEAL_ENVIDO_SCORE:
		action = &ActionRevealEnvidoScore{}
	default:
		return nil, fmt.Errorf("unknown action: [%v]", string(bs))
	}

	err = json.Unmarshal(bs, action)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func _serializeActions(as []Action) []json.RawMessage {
	_as := []json.RawMessage{}
	for _, a := range as {
		_as = append(_as, json.RawMessage(SerializeAction(a)))
	}
	return _as
}

func _deserializeCurrentRoundLastAction(g GameState) Action {
	lastAction := g.RoundsLog[g.RoundNumber].ActionsLog[len(g.RoundsLog[g.RoundNumber].ActionsLog)-1].Action
	a, _ := DeserializeAction(lastAction)
	return a
}

func (g *GameState) ToClientGameState(youPlayerID int) ClientGameState {
	themPlayerID := g.OpponentOf(youPlayerID)

	// GameState may have possible game actions that this player can't take.
	filteredPossibleActions := []Action{}
	for _, a := range g.CalculatePossibleActions() {
		if a.GetPlayerID() == youPlayerID {
			filteredPossibleActions = append(filteredPossibleActions, a)
		}
	}

	cgs := ClientGameState{
		RoundTurnPlayerID:           g.RoundTurnPlayerID,
		RoundNumber:                 g.RoundNumber,
		TurnPlayerID:                g.TurnPlayerID,
		YouPlayerID:                 youPlayerID,
		ThemPlayerID:                themPlayerID,
		YourScore:                   g.Players[youPlayerID].Score,
		TheirScore:                  g.Players[themPlayerID].Score,
		YourRevealedCards:           g.Players[youPlayerID].Hand.Revealed,
		TheirRevealedCards:          g.Players[themPlayerID].Hand.Revealed,
		YourUnrevealedCards:         g.Players[youPlayerID].Hand.Unrevealed,
		PossibleActions:             _serializeActions(filteredPossibleActions),
		IsGameEnded:                 g.IsGameEnded,
		IsRoundFinished:             g.IsRoundFinished,
		WinnerPlayerID:              g.WinnerPlayerID,
		EnvidoWinnerPlayerID:        g.RoundsLog[g.RoundNumber].EnvidoWinnerPlayerID,
		WasEnvidoAccepted:           g.EnvidoSequence.WasAccepted(),
		EnvidoPoints:                g.RoundsLog[g.RoundNumber].EnvidoPoints,
		TrucoWinnerPlayerID:         g.RoundsLog[g.RoundNumber].TrucoWinnerPlayerID,
		TrucoPoints:                 g.RoundsLog[g.RoundNumber].TrucoPoints,
		WasTrucoAccepted:            g.TrucoSequence.WasAccepted(),
		YourDisplayUnrevealedCards:  g.Players[youPlayerID].Hand.prepareDisplayUnrevealedCards(true),
		TheirDisplayUnrevealedCards: g.Players[themPlayerID].Hand.prepareDisplayUnrevealedCards(false),
	}

	if len(g.RoundsLog[g.RoundNumber].ActionsLog) > 0 {
		actionsLog := g.RoundsLog[g.RoundNumber].ActionsLog
		cgs.LastActionLog = &actionsLog[len(actionsLog)-1]
	}

	return cgs
}

// ClientGameState represents the state of a Truco game as available to a client.
//
// It is returned by the server on every single call, so if you want to implement a client,
// you need to be very familiar with this struct.
type ClientGameState struct {
	// RoundTurnPlayerID is the player ID of the player who starts the round, or "mano".
	RoundTurnPlayerID int `json:"roundTurnPlayerID"`

	// RoundNumber is the number of the current round, starting from 1.
	RoundNumber int `json:"roundNumber"`

	// TurnPlayerID is the player ID of the player whose turn it is to play an action.
	// This is different from RoundTurnPlayerID, which is the player who starts the round.
	// They are the same at the beginning of the round.
	TurnPlayerID int `json:"turnPlayerID"`

	YouPlayerID         int    `json:"you"`
	ThemPlayerID        int    `json:"them"`
	YourScore           int    `json:"yourScore"`
	TheirScore          int    `json:"theirScore"`
	YourRevealedCards   []Card `json:"yourRevealedCards"`
	TheirRevealedCards  []Card `json:"theirRevealedCards"`
	YourUnrevealedCards []Card `json:"yourUnrevealedCards"`

	// YourDisplayUnrevealedCards is like YourUnrevealedCards, but it always has 3 cards
	// and it adds two properties: `IsBackwards` and `IsHole`.
	//
	// `IsBackwards` is true if the card is facing backwards (i.e. the client doesn't know what it is).
	// `IsHole` is true if the card was revealed by the opponent, and it used to be that card.
	//
	// Use this property to render the card
	YourDisplayUnrevealedCards []DisplayCard `json:"yourDisplayUnrevealedCards"`

	// TheirDisplayUnrevealedCards is like TheirUnrevealedCards, but it always has 3 cards
	// and it adds two properties: `IsBackwards` and `IsHole`.
	//
	// `IsBackwards` is true if the card is facing backwards (i.e. the client doesn't know what it is).
	// `IsHole` is true if the card was revealed by the opponent, and it used to be that card.
	//
	// Use this property to render the card
	TheirDisplayUnrevealedCards []DisplayCard `json:"theirDisplayUnrevealedCards"`

	// PossibleActions is a list of possible actions that the current player can take.
	// Possible actions are calculated based on game state at the beginning of the round and after
	// each action is run (i.e. GameState.RunAction).
	PossibleActions []json.RawMessage `json:"possibleActions"`

	// IsGameEnded is true if the whole game is ended, rather than an individual round. This happens when
	// a player reaches MaxPoints points.
	IsGameEnded bool `json:"isGameEnded"`

	IsRoundFinished bool `json:"isRoundFinished"`

	// WinnerPlayerID is the player ID of the player who won the game. This is only set when `IsGameEnded` is
	// `true`. Otherwise, it's -1.
	WinnerPlayerID int `json:"winnerPlayerID"`

	// Some state information about the current round, in case it's useful to the client.
	EnvidoWinnerPlayerID int  `json:"envidoWinnerPlayerID"`
	WasEnvidoAccepted    bool `json:"wasEnvidoAccepted"`
	EnvidoPoints         int  `json:"envidoPoints"`
	TrucoWinnerPlayerID  int  `json:"trucoWinnerPlayerID"`
	TrucoPoints          int  `json:"trucoPoints"`
	WasTrucoAccepted     bool `json:"wasTrucoAccepted"`

	// LastActionLog is the log of the last action that was run in the current round. If the round has
	// just started, this will be nil. Clients typically want to user this to show the current player
	// what the opponent just did.
	LastActionLog *ActionLog `json:"lastActionLog"`
}

type Bot interface {
	ChooseAction(ClientGameState) Action
}
