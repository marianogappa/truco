package exampleclient

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/marianogappa/truco/server"
	"github.com/marianogappa/truco/truco"
)

func Player(playerID int, address string) {
	ui := NewUI()
	defer ui.Close()

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%v/ws", address), nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	if err := server.WsSend(conn, server.NewMessageHello(playerID)); err != nil {
		log.Fatal(err)
	}

	lastRound := 0
	for {
		gameState, err := server.WsReadMessage[truco.GameState, server.MessageHeresGameState](conn, server.MessageTypeHeresGameState)
		if err != nil {
			log.Fatal(err)
		}

		if gameState.IsGameEnded {
			_ = ui.render(playerID, *gameState, PRINT_MODE_END)
			ui.pressAnyKey()
			return
		}

		if gameState.RoundNumber != lastRound && lastRound != 0 {
			err := ui.render(playerID, *gameState, PRINT_MODE_SHOW_ROUND_RESULT)
			if err != nil {
				log.Fatal(err)
			}
			ui.pressAnyKey()
		}
		lastRound = gameState.RoundNumber

		if err := ui.render(playerID, *gameState, PRINT_MODE_NORMAL); err != nil {
			log.Fatal(err)
		}
		if gameState.TurnPlayerID != playerID {
			continue
		}

		var (
			action          truco.Action
			possibleActions = _deserializeActions(gameState.PossibleActions)
		)
		for {
			num := ui.pressAnyNumber()
			if num > len(possibleActions) {
				continue
			}
			action = possibleActions[num-1]
			break
		}

		msg, _ := server.NewMessageAction(action)
		if err := server.WsSend(conn, msg); err != nil {
			log.Fatal(err)
		}
	}
}
