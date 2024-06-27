package examplebot

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gorilla/websocket"
	"github.com/marianogappa/truco/server"
	"github.com/marianogappa/truco/truco"
)

func Bot(playerID int, address string) {
	// Open the WebSocket connection, and send a hello message.

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

	// On each iteration
	for {
		clientGameState, err := server.WsReadMessage[truco.ClientGameState, server.MessageHeresGameState](conn, server.MessageTypeHeresGameState)
		if err != nil {
			log.Fatal(err)
		}

		if clientGameState.IsGameEnded {
			return
		}

		if clientGameState.TurnPlayerID != playerID {
			continue
		}

		// Get a random element from clientGameState.PossibleActions.
		randomAction := clientGameState.PossibleActions[rand.Intn(len(clientGameState.PossibleActions))]

		// Send the action to the server.
		if err := server.WsSend(conn, server.MessageAction{WebsocketMessage: server.WebsocketMessage{Type: server.MessageTypeAction}, Action: randomAction}); err != nil {
			log.Fatal(err)
		}
	}
}
