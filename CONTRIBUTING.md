# Contributing

If you're a little lost on how to contribute and these notes are not helpful, reach me online and I'm happy to chat about it. I'm not hard to find.

## Making your own bot

### Getting started

The bot interface is very straightforward; you only need to implement a single function:

```go
type Bot interface {
	ChooseAction(ClientGameState) Action
}
```

Upon a given `ClientGameState`, choose an action. `ClientGameState` provides a `PossibleActions` property, and these are the only possible actions, so you just need to pick one. Use `truco.DeserializeAction` to be able to return it.

The `ClientGameState` struct is a "view" of the main `GameState` struct from the point of view of the bot. This prevents the bot from seeing the opponent's cards, but also simplifies the state.

There are subtleties to implementing this function, such as:
- There could be no possible actions. You must return `nil` in this case.
- There could be only one possible action. You must return this action in this case.
- Review the existing bot for inspiration. You'll have to figure out how to calculate envido/flor scores, the results of the card faceoffs, etc. I would clone the existing bot as a starting point.

### My bot is ready, how do I test it?

You should be able to instantiate your own bot instead of the existing one in this one line:

https://github.com/marianogappa/truco/blob/main/main.go#L51

### I don't know Go, can I still make a Bot?

The server implementation allows any client that understands it's WebSocket message implementation to play a game.

Here's the server code: https://github.com/marianogappa/truco/blob/main/server/websocket_server.go#L36

And here's an example bot client that implements the client WebSocket message implementation:

https://github.com/marianogappa/truco/blob/main/botclient/main.go

The implementation is quite straightforward, so I encourage you to implement yours in whichever language you want. As long as your code can address the server, you can make it work. If you're stuck, let me know.

## Making your own frontend

Writing a frontend is a more involved task, but in terms of the communication with the Truco engine, it's essentially exactly the same as making your own bot.

You will have to implement the same Websocket message implementation, you get the same state struct (`ClientGameState`), and you must send the actions that the user selects in the same fashion that you did for the bot, the only difference is that the user is picking them, rather than an algorithm.

The `ClientGameState` struct is designed to be straightforward for making a frontend implementation. Even the cards information is presented in a way that you're able to know how to animate the card from source to destination (as an example).

Please use the existing implementations to guide your own; let me know if you get stuck.

## Contributing guidelines

Feel free to contribute informally; please add tests if possible. Reach out if you need help.

## Basic Flow Diagram