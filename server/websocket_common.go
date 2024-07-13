//go:build !tinygo
// +build !tinygo

package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func WsSend(conn *websocket.Conn, message any) error {
	bs, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, bs); err != nil {
		log.Println("Failed to write message:", err)
	}
	return err
}

func WsReadMessage[U any, T IWebsocketMessage[U]](conn *websocket.Conn, expectedType int) (*U, error) {
	messageType, message, err := conn.ReadMessage()
	if messageType != websocket.TextMessage {
		return nil, fmt.Errorf("Expected text message, got %d with error %v", messageType, err)
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to read message from client: %v", err)
	}
	return WsDeserializeMessage[U, T](message, expectedType)
}

func WsDeserializeMessage[U any, T IWebsocketMessage[U]](message []byte, expectedType int) (*U, error) {
	var m T
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal message: %v", err)
	}
	if m.GetType() != expectedType {
		return nil, fmt.Errorf("Expected message type %d, got %d", expectedType, m.GetType())
	}
	elem, err := m.Deserialize()
	return &elem, err
}
