//go:build tinygo
// +build tinygo

package main

import "github.com/marianogappa/truco/truco"

func NewTruco() *truco.GameState {
	return truco.New()
}

func multiply(x, y int) int {
	return x * y
}

func main() {

	select {}
}
