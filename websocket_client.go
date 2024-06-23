package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/marianogappa/truco/truco"
	"github.com/nsf/termbox-go"
)

func player(playerID int, address string) {
	ui := NewUI()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%v/ws", address), nil)
	if err != nil {
		log.Println("Failed to connect to WebSocket server:", err)
		return
	}
	defer conn.Close()

	if err := wsSend(conn, NewMessageHello(playerID)); err != nil {
		log.Println(err)
		return
	}

	lastRound := 0
	for {
		gameState, err := wsReadMessage[truco.GameState, MessageHeresGameState](conn, MessageTypeHeresGameState)
		if err != nil {
			log.Println(err)
			return
		}

		if gameState.IsEnded {
			_ = ui.printState(playerID, *gameState, PRINT_MODE_END)
			return
		}

		if gameState.RoundNumber != lastRound && lastRound != 0 {
			err := ui.printState(playerID, *gameState, PRINT_MODE_SHOW_ROUND_RESULT)
			if err != nil {
				log.Println(err)
				return
			}
		}
		lastRound = gameState.RoundNumber

		if gameState.TurnPlayerID != playerID {
			err := ui.printState(playerID, *gameState, PRINT_MODE_NORMAL)
			if err != nil {
				log.Println(err)
				return
			}
			continue
		}

		action, err := ui.play(playerID, *gameState)
		if err != nil {
			fmt.Println("Invalid action:", err)
			break
		}

		msg, _ := NewMessageAction(action)
		if err := wsSend(conn, msg); err != nil {
			log.Println(err)
			return
		}
	}
}
