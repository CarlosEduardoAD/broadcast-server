# Broadcast Server

A toy app to illustrate a simulated realmtime chat that you can run in your terminals!

## Installation

Clone this repo.

```bash
git clone https://github.com/CarlosEduardoAD/broadcast-server.git
```

Then after build the project

```bash
go build -o ./broadcast-server main.go
```

## Usage

```bash
./broadcast-server -h
```

It should output

```bash
A toy app to represent a simulated realtime chat in your terminal!

Usage:
  broadcast-server [command]

Available Commands:
  client      All the commands a client can use to interact with the chat
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  server      Commands to manage the websocket server

Flags:
  -a, --basicAuth string   Authenticate connection with basic auth
  -h, --help               help for broadcast-server
  -t, --toggle             Help message for toggle
```

## Quickstart

Start your server

```bash
./broadcast-server server start
```

Spawn 2 or more terminals, run the command below and begin to send messages!

```bash
./broadcast-server client join
```

## .env

To change the basic auth password and the allowed origin, head up to these 2 variables.

```bash
ALLOWED_ORIGIN=http://your-host:your-port/ 
BASIC_AUTH_CREDENTIALS=you:yourpassword123
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Known Issues

- Bug that keeps sending messages to closed connections

## License

[MIT](https://opensource.org/license/mit)