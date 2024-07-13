//go:build !tinygo
// +build !tinygo

package botclient

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/marianogappa/truco/server"
	"github.com/marianogappa/truco/truco"
)

func Bot(playerID int, address string, bot truco.Bot) {
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

		botAction := bot.ChooseAction(*clientGameState)

		if botAction == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		bs, _ := json.Marshal(botAction)

		// Send the action to the server.
		if err := server.WsSend(conn, server.MessageAction{WebsocketMessage: server.WebsocketMessage{Type: server.MessageTypeAction}, Action: bs}); err != nil {
			log.Fatal(err)
		}
	}
}
