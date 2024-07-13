//go:build tinygo
// +build tinygo

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/marianogappa/truco/examplebot"
	"github.com/marianogappa/truco/truco"
)

func main() {
	js.Global().Set("trucoNew", js.FuncOf(trucoNew))
	js.Global().Set("trucoRunAction", js.FuncOf(trucoRunAction))
	js.Global().Set("trucoBotRunAction", js.FuncOf(trucoBotRunAction))
	select {}
}

var (
	state *truco.GameState
	bot   truco.Bot
)

func trucoNew(this js.Value, p []js.Value) interface{} {
	state = truco.New()
	bot = examplebot.New()

	nbs, err := json.Marshal(state.ToClientGameState(0))
	if err != nil {
		panic(err)
	}

	buffer := js.Global().Get("Uint8Array").New(len(nbs))
	js.CopyBytesToJS(buffer, nbs)
	return buffer
}

func trucoRunAction(this js.Value, p []js.Value) interface{} {
	jsonBytes := make([]byte, p[0].Length())
	js.CopyBytesToGo(jsonBytes, p[0])

	newBytes := _runAction(jsonBytes)

	buffer := js.Global().Get("Uint8Array").New(len(newBytes))
	js.CopyBytesToJS(buffer, newBytes)
	return buffer
}

func trucoBotRunAction(this js.Value, p []js.Value) interface{} {
	action := bot.ChooseAction(state.ToClientGameState(1))

	err := state.RunAction(action)
	if err != nil {
		panic(err)
	}
	nbs, err := json.Marshal(state.ToClientGameState(0))
	if err != nil {
		panic(err)
	}

	buffer := js.Global().Get("Uint8Array").New(len(nbs))
	js.CopyBytesToJS(buffer, nbs)
	return buffer
}

func _runAction(bs []byte) []byte {
	action, err := truco.DeserializeAction(bs)
	if err != nil {
		panic(err)
	}
	err = state.RunAction(action)
	if err != nil {
		panic(err)
	}
	fmt.Println("Ran action:", string(bs))
	nbs, err := json.Marshal(state.ToClientGameState(0))
	if err != nil {
		panic(err)
	}
	return nbs
}
