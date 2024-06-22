package main

import (
	"fmt"

	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: truco server|player1|player2 [address]")
		fmt.Println("Define the PORT environment variable to change the default port (8080).")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf("localhost:%v", port)
	if len(os.Args) >= 3 {
		address = os.Args[2]
	}

	arg := os.Args[1]
	switch arg {
	case "server":
		serve(port)
	case "player1":
		player(0, address)
	case "player2":
		player(1, address)
	default:
		fmt.Println("Invalid argument. Please provide either server or client.")
	}
}
