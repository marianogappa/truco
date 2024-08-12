# Truco

https://marianogappa.github.io/truco-argentino/

Truco argentino implementation, featuring:

- [websocket-based client/server architecture](https://github.com/marianogappa/truco/tree/main/server)
- [example terminal-based frontend](https://github.com/marianogappa/truco/blob/main/exampleclient/websocket_client.go)
- [example React-based frontend](https://github.com/marianogappa/truco-argentino), and a simple, [documented](https://github.com/marianogappa/truco/blob/main/CONTRIBUTING.md#making-your-own-frontend) interface for making your own frontend
- [example bot](https://github.com/marianogappa/truco/blob/main/examplebot/newbot/bot.go), and a simple, [documented](https://github.com/marianogappa/truco/blob/main/CONTRIBUTING.md#making-your-own-bot) interface for making your own

![game](https://github.com/user-attachments/assets/c85b7a6d-a6c9-4556-a1ac-da74150478e6)



<img width="1512" alt="Screenshot 2024-06-23 at 19 26 11" src="https://github.com/marianogappa/truco/assets/1078546/881e7204-f1a6-4de2-a0b5-60faa43b4fac">

### Installation

Either install using Go

```bash
$ go install https://github.com/marianogappa/truco@latest
```

Or download the [latest release binary](https://github.com/marianogappa/truco/releases) for your OS.

### Usage

Start a server

```bash
$ truco server
```

You may change the port (default is 8080) via environment variable

```bash
$ PORT=1234 truco server
```

If you want to play via example terminal-based frontend, start two clients on separate terminals

```bash
$ truco player 1
```

```bash
$ truco player 2
```

### Playing with someone else over the Internet

Whoever starts the server may expose it to the Internet somehow, e.g. via `cloudflared` tunnels

```bash
$ cloudflared tunnel --url localhost:8080
```

Then, the clients can connect to the address the tunnel provides, e.g. if tunnel says

```bash
...
2024-06-23T18:35:10Z INF +--------------------------------------------------------------------------------------------+
2024-06-23T18:35:10Z INF |  Your quick Tunnel has been created! Visit it at (it may take some time to be reachable):  |
2024-06-23T18:35:10Z INF |  https://retail-curves-bernard-affairs.trycloudflare.com                                   |
2024-06-23T18:35:10Z INF +--------------------------------------------------------------------------------------------+
```

Start the clients with

```bash
$ truco player 1 retail-curves-bernard-affairs.trycloudflare.com
```

```bash
$ truco player 2 retail-curves-bernard-affairs.trycloudflare.com
```

### Reconnect after issue

If the server dies, state is gone. If client dies, you can simply reconnect to the same server and game goes on.

### I don't like your UI

It's just an example UI. I encourage you to [implement your own frontend](https://github.com/marianogappa/truco/blob/main/CONTRIBUTING.md#making-your-own-frontend). You may [browse the documentation](https://github.com/marianogappa/truco/blob/main/CONTRIBUTING.md) and the [existing React-based UI code](https://github.com/marianogappa/truco-argentino) and [terminal UI code](https://github.com/marianogappa/truco/blob/main/exampleclient/ui.go) to guide your implementation.

### I don't like your Bot

It's just an example bot, which beats me. I encourage you to [implement your own bot](https://github.com/marianogappa/truco/blob/main/CONTRIBUTING.md#making-your-own-bot). You may [browse the documentation](https://github.com/marianogappa/truco/blob/main/CONTRIBUTING.md) and the [existing bot code](https://github.com/marianogappa/truco/blob/main/examplebot/newbot/bot.go) to guide your implementation.

## Technology stack

- This truco engine is written 100% in Go
- Terminal-based UI uses [Termbox](https://github.com/nsf/termbox-go)
- React-based UI uses [TinyGo](https://tinygo.org/) with WASM target to transpile to WebAssembly, and the frontend itself is built in React

### Known issues / limitations

- Don't resize your terminal. This is a go-termbox issue. Also, have a terminal with a decent viewport. That is on me mostly.

### Issues / Improvements

Please do [create issues](https://github.com/marianogappa/truco/issues) and send PRs. Also feel free to reach me for comments / discussions. I'm not hard to find.
