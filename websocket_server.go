package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/marianogappa/truco/truco"
)

// Define WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serve(port string) {
	NewServer(port).Start()
}

type server struct {
	gameState *truco.GameState
	port      string
	players   []*websocket.Conn
}

func NewServer(port string) *server {
	return &server{gameState: truco.New(), port: port, players: []*websocket.Conn{nil, nil}}
}

func (s *server) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", s.handleWebSocket)
	log.Printf("Server running on port %v\n", s.port)
	log.Fatal(http.ListenAndServe(":"+s.port, router))
}

func (s *server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	playerID, err := wsReadMessage[int, MessageHello](conn, MessageTypeHello)
	if err != nil {
		log.Println(err)
		return
	}

	if *playerID < 0 || *playerID > 1 {
		log.Println("Invalid player ID")
		return
	}
	if s.players[*playerID] != nil {
		log.Println("Player already connected")
		return
	}
	s.players[*playerID] = conn

	msg, _ := NewMessageHeresGameState(*s.gameState)
	if err := wsSend(conn, msg); err != nil {
		log.Println(err)
		return
	}
	log.Println("Player", *playerID, "connected")

	for {
		log.Println("Waiting for action/state_request from player", *playerID)
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message from client, freeing slot:", err)
			s.players[*playerID] = nil
			break
		}

		var wsMessage WebsocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			log.Println("Failed to unmarshal message:", err)
			break
		}

		switch wsMessage.Type {
		case MessageTypeAction:
			log.Println("Got action message:", string(message))
			action, err := wsDeserializeMessage[truco.Action, MessageAction](message, MessageTypeAction)
			if err != nil {
				log.Println(err)
				return
			}
			err = s.gameState.RunAction(*action)
			if err != nil {
				// TODO write back to the connection
				log.Println("Failed to run action:", err)
				break
			}

			log.Println("Ran action message:", string(message))

			msg, _ := NewMessageHeresGameState(*s.gameState)
			for i, playerConn := range s.players {
				log.Println("Sending game state to player", i)
				if err := wsSend(playerConn, msg); err != nil {
					log.Println(err)
					return
				}
			}
		case MessageTypeGimmeGameState:
			log.Println("Got state request message:", string(message))

			msg, _ := NewMessageHeresGameState(*s.gameState)
			if err := wsSend(conn, msg); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
