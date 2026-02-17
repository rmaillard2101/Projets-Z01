[GitHub deposit link (Oscar Hernandez)](https://github.com/IronBeagle404/lem-in)

TCPChat

## 💬 Overview

    TCPChat is a TCP-based chat server written in Go.
    It recreates a simplified NetCat-like group chat, allowing multiple clients to communicate in real time through a single server.

    The goal of the project is to have a server that handle concurrent TCP connections, manage users safely, and broadcast messages reliably while respecting strict technical constraints.

## 💬 Client Interaction Rules

    Each connected client must follow a set of rules enforced by the server:
    -A client must provide a non-empty username upon connection
    -A maximum of 10 clients can be connected simultaneously
    -Empty messages are not broadcast
    -Each message is:
        -Timestamped
        -Tagged with the sender’s username
    -All clients receive messages sent by other clients
    -When a client joins the server:
    -All previous messages are sent to them
    -Other clients are notified
    -When a client leaves:
        -Other clients are notified
        -Remaining clients stay connected
    These rules ensure a stable, synchronized, and user-friendly group chat experience.

## ⚙️ Features

    -TCP server handling multiple concurrent clients
    -Username management
    -Message broadcasting with timestamps
    -Join / leave notifications
    -Message history synchronization for new clients
    -Connection limit enforcement (max 10 clients)
    -Logging system for server activity
    -Graceful error handling on server and client sides

    -Message format: [YYYY-MM-DD HH:MM:SS][username]:message

## 🚀 Usage

    ```
    go build -o TCPChat main.go
    ./TCPChat [port]
    ```

    If no port is provided, the default port is 8989

    If invalid arguments are given, the program outputs: [USAGE]: ./TCPChat $port

    The server listens for incoming TCP connections and handles clients using goroutines.

## 📥 Client Connection

    Clients can connect using NetCat:
    ```
    nc <host> <port>
    ```

    Upon connection, the server displays a welcome message and prompts for a username:

```
Welcome to TCP-Chat!
_nnnn_
dGGGGMMb
@p~qp~~qMb
M|@||@) M|
@,----.JM|
JS^\__/ qKL
dZP qKRb
dZP qKKb
fZP SMMb
HZM MMMM
FqM MMMM
\_\_| ". |\dS"qML
| `.       | `' \Zq
_) \.**\_.,| .'
\_\_** )MMMMMP| .'
`-'       `--'
[ENTER YOUR NAME]:
```

## 🧑‍💻 Client Commands

    Clients can interact with the server using the following commands:
    -/help — display available commands
    -/quit — leave the server
    -/newname <username> — change username
    -/join <room> — switch chat room
    -/list — list available rooms

## 📝 Logs

    Server logs are stored in the logfiles directory.

    Log file format: log_YYYY-MM-DD_HH-MM-SS.log

    Logs include connection events, messages, and disconnections.

🗂 Project Structure

```
.
├── main.go # Program entry point
├── server/ # TCP server and connection handling
├── client/ # Client-related logic
├── chat/ # Message broadcasting and rooms
├── logger/ # Logging system
├── errors/ # Custom error handling
└── tests/ # Unit and integration tests
```

👤 Authors

Author details are available in the creditentials.md file.
