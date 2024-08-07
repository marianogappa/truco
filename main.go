//go:build !tinygo
// +build !tinygo

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/marianogappa/truco/botclient"
	"github.com/marianogappa/truco/examplebot/newbot"
	"github.com/marianogappa/truco/exampleclient"
	"github.com/marianogappa/truco/server"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cmd := os.Args[1]

	address := fmt.Sprintf("localhost:%v", port)
	if len(os.Args) >= 4 {
		address = os.Args[3]
	}

	var (
		playerNum int
		err       error
	)
	if cmd == "player" || cmd == "bot" {
		playerNum, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid player number. Please provide a number.")
			usage()
		}
	}

	switch cmd {
	case "server":
		server.New(port).Start()
	case "player":
		exampleclient.Player(playerNum-1, address)
	case "bot":
		botclient.Bot(playerNum-1, address, newbot.New(newbot.WithDefaultLogger))
	default:
		fmt.Println("Invalid argument. Please provide either server or client.")
	}
}

func usage() {
	fmt.Println("usage: truco server")
	fmt.Println("usage: truco player %number [address]")
	fmt.Println("usage: truco bot %number [address]")
	fmt.Println("usage: e.g. truco player 1")
	fmt.Println("usage: e.g. truco player 2")
	fmt.Println("usage: e.g. truco player 1 localhost:8080")
	fmt.Println("usage: e.g. truco bot 1 localhost:8080")
	fmt.Println("usage: e.g. truco bot 2")
	fmt.Println("Define the PORT environment variable for truco server to change the default port (8080).")
	os.Exit(1)
}
