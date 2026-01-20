package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	HOST = "0.0.0.0"
	TYPE = "tcp"
)

var (
	PORT        = "8989"
	ConnsCount  = 0
	logFileName string

	clients = make(map[string]*Client) // key = addr string
	rooms   = make(map[string]*Room)   // key = room name
)

// --- struct Client ---
type Client struct {
	Conn     net.Conn
	Username string
	Room     string
}

type Room struct {
	Name         string
	Messages     chan Message
	PastMessages []string
}

// Message struct with type
type Message struct {
	Content string
	Type    string // "normal", "join", "leave"
}

func main() {

	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	if len(os.Args) == 2 {
		PORT = os.Args[1]
	}

	logDir := "./logfiles"
	os.MkdirAll(logDir, 0755)

	dateStr := time.Now().Format("2006-01-02 15:04:05")
	logFileName = logDir + "/log_" + time.Now().Format("2006-01-02_15-04-05") + ".log"

	logFile, err := os.Create(logFileName)
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return
	}
	logFile.WriteString("Server started here " + dateStr + ":\n")
	logFile.Close()

	addr, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	listener, err := net.ListenTCP(TYPE, addr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on the port :" + PORT)

	// Init default room
	initRoom("default")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept failed : %v\n", err)
			continue
		}

		if ConnsCount < 10 {
			go handleConnection(conn)
			ConnsCount++
		} else {
			conn.Write([]byte("Server is full, impossible to connect.\n"))
			conn.Close()
		}
	}
}

// Init room
func initRoom(name string) {
	if _, exists := rooms[name]; !exists {
		rooms[name] = &Room{
			Name:     name,
			Messages: make(chan Message, 10),
		}
		go BroadcastMessages(rooms[name])
	}
}

// Broadcast in room only, with log in global file
func BroadcastMessages(room *Room) {
	for {
		msg := <-room.Messages

		var fullMsg string
		if msg.Type == "normal" {
			timeNow := time.Now().Format("2006-01-02 15:04:05")
			fullMsg = "[" + timeNow + "] [" + room.Name + "] " + msg.Content + "\n"
		} else {
			fullMsg = "[" + room.Name + "] " + msg.Content + "\n"
		}

		logPrint(logFileName, fullMsg)

		for _, client := range clients {
			if client.Room == room.Name {
				_, err := client.Conn.Write([]byte(fullMsg))
				if err != nil {
					log.Printf("broadcast message to %v has failed : %v\n", client.Username, err)
					client.Conn.Close()
					delete(clients, client.Conn.RemoteAddr().String())
					ConnsCount--
				}
			}
		}

		room.PastMessages = append(room.PastMessages, fullMsg)
	}
}

func logPrint(logFileName string, msg string) {
	fmt.Print(msg)

	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(msg)
	if err != nil {
		fmt.Printf("Failed to write to log file: %v\n", err)
	}
}

// simplified handleConnection with room management
func handleConnection(conn net.Conn) {
	conn.Write([]byte("Welcome to TCP-chat!\n"))

	tux := "         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		" |    `.       | `' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     `-'       `--'"
	conn.Write([]byte(tux + "\n"))

	buffer := make([]byte, 1024)

	var username string

	for {
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		n, err := conn.Read(buffer)
		if err != nil || n <= 1 {
			conn.Write([]byte("Invalid username, please choose another name.\n"))
			continue
		}
		username = strings.TrimSpace(string(buffer[:n]))
		if username == "" {
			conn.Write([]byte("Username cannot be empty, please choose a name.\n"))
			continue
		}

		nameTaken := false
		for _, c := range clients {
			if c.Username == username {
				nameTaken = true
				break
			}
		}
		if nameTaken {
			conn.Write([]byte("Username already taken, please choose another name.\n"))
			continue
		}
		if len(username) > 64 {
			conn.Write([]byte("Invalid username, please choose another name.\n"))
			continue
		}
		break
	}

	// default room name
	roomName := "default"
	initRoom(roomName) // makes sure room exists

	client := &Client{
		Conn:     conn,
		Username: username,
		Room:     roomName,
	}
	clients[conn.RemoteAddr().String()] = client

	// send room history
	for _, msg := range rooms[roomName].PastMessages {
		conn.Write([]byte(msg))
	}

	// notify arrival
	rooms[roomName].Messages <- Message{
		Content: username + " has joined the room " + roomName,
		Type:    "join",
	}

	for {
		length, err := conn.Read(buffer)

		if err != nil {
			if err != io.EOF {
				logPrint(logFileName, "read client message failed : "+err.Error())
			}
			break
		}

		if length <= 1 {
			continue
		}

		recvStr := strings.TrimSpace(string(buffer[:length]))

		if strings.HasPrefix(recvStr, "/") {
			cmd := recvStr[1:]
			if strings.HasPrefix(cmd, "newname ") {
				newName := strings.TrimSpace(cmd[len("newname "):])
				if newName == "" {
					conn.Write([]byte("New name cannot be empty.\n"))
					continue
				}
				nameTaken := false
				for _, c := range clients {
					if c.Username == newName {
						nameTaken = true
						break
					}
				}
				if nameTaken {
					conn.Write([]byte("This name is already taken. Choose another one.\n"))
					continue
				}
				oldName := client.Username
				client.Username = newName
				conn.Write([]byte("Your username has been changed to " + newName + ".\n"))
				rooms[client.Room].Messages <- Message{
					Content: oldName + " has changed their name to " + newName + ".",
					Type:    "normal",
				}

			} else if strings.HasPrefix(cmd, "join ") {
				// Switch room
				newRoom := strings.TrimSpace(cmd[len("join "):])
				if newRoom == "" {
					conn.Write([]byte("Usage: /join roomname\n"))
					continue
				}
				if newRoom == client.Room {
					conn.Write([]byte("You are already in room " + newRoom + "\n"))
					continue
				}

				// Leave old room
				oldRoom := client.Room
				rooms[oldRoom].Messages <- Message{
					Content: client.Username + " has left the room.",
					Type:    "leave",
				}

				// Join new room
				initRoom(newRoom)
				client.Room = newRoom

				// Send new room history
				for _, msg := range rooms[newRoom].PastMessages {
					conn.Write([]byte(msg))
				}

				rooms[newRoom].Messages <- Message{
					Content: client.Username + " has joined the room " + newRoom,
					Type:    "join",
				}

			} else if cmd == "list" {
				var roomNames []string
				for roomName := range rooms {
					roomNames = append(roomNames, roomName)
				}
				conn.Write([]byte("Rooms available:\n"))
				for _, rn := range roomNames {
					conn.Write([]byte(" - " + rn + "\n"))
				}
			} else if cmd == "help" {
				conn.Write([]byte("Available commands: /help, /quit, /newname yourname, /join roomname, /list\n"))
			} else if cmd == "quit" {
				conn.Write([]byte("Disconnecting...\n"))
				break
			} else {
				conn.Write([]byte("Unknown command.\n"))
			}
		} else {
			rooms[client.Room].Messages <- Message{
				Content: "[" + client.Username + "]: " + recvStr,
				Type:    "normal",
			}
		}
	}

	// cleanup client
	delete(clients, conn.RemoteAddr().String())
	rooms[client.Room].Messages <- Message{
		Content: client.Username + " has left the room.",
		Type:    "leave",
	}
	conn.Close()
	ConnsCount--
}
