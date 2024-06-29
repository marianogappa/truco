//go:build !tinygo
// +build !tinygo

package exampleclient

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/marianogappa/truco/server"
	"github.com/marianogappa/truco/truco"
)

func Player(playerID int, address string) {
	// Create a UI, open the WebSocket connection, and send a hello message.
	ui := NewUI()
	defer ui.Close()

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%v/ws", address), nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Hello message is meant to tell the server who we are, and request game state.
	// Game could be in progress (this could be a reconnection).
	if err := server.WsSend(conn, server.NewMessageHello(playerID)); err != nil {
		log.Fatal(err)
	}

	lastRound := 0
	// On each iteration
	for {
		// Read the game state from the server.
		clientGameState, err := server.WsReadMessage[truco.ClientGameState, server.MessageHeresGameState](conn, server.MessageTypeHeresGameState)
		if err != nil {
			log.Fatal(err)
		}

		// If the game has ended, render the game state (with a final message) and exit.
		if clientGameState.IsGameEnded {
			_ = ui.render(*clientGameState, PRINT_MODE_END)
			ui.pressAnyKey()
			return
		}

		// If the round has ended, render the game state (with a summary) and wait for a key press.
		if clientGameState.RoundNumber != lastRound && lastRound != 0 {
			err := ui.render(*clientGameState, PRINT_MODE_SHOW_ROUND_RESULT)
			if err != nil {
				log.Fatal(err)
			}
			ui.pressAnyKey()
		}
		lastRound = clientGameState.RoundNumber

		// Render the game state. One could arrive here in 2 ways:
		// 1. A round had just ended, the player pressed a key, and we're here to render the new round.
		// 2. A turn starts (maybe even the game's first turn), and we're here to render the new turn.
		if err := ui.render(*clientGameState, PRINT_MODE_NORMAL); err != nil {
			log.Fatal(err)
		}

		// If it's not the current player's turn, wait for the next game state.
		if clientGameState.TurnPlayerID != playerID {
			continue
		}

		// If it's the current player's turn, wait for the player to choose action.
		var (
			action          truco.Action
			possibleActions = _deserializeActions(clientGameState.PossibleActions)
		)
		for {
			num := ui.pressAnyNumber()
			if num > len(possibleActions) {
				continue
			}
			action = possibleActions[num-1]
			break
		}

		// Send the action to the server.
		msg, _ := server.NewMessageAction(action)
		if err := server.WsSend(conn, msg); err != nil {
			log.Fatal(err)
		}
	}
}
