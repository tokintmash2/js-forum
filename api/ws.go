package api

import (
	"encoding/json"
	"log"
	"net/http"
	"real-forum/structs"
	"real-forum/utils"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	connections = make(map[string]*Client)
	broadcast   = make(chan structs.SocketMessage)
	mu          sync.Mutex
)

func init() {
	go broadcastOnlineUsers()
}

// Upgrader to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{ // buffers missing
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin:     checkOrigin,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	connection  *websocket.Conn
	send        chan []byte
	mu          sync.Mutex
	connOwnerId string
	lastActive  time.Time
}

// type Message struct {
// 	Recipient string `json:"recipient"`
// 	Content   string `json:"content"`
// 	Sender int
// 	CreatedAt time.Time
// }

func HandleConnections(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from the request (e.g., from a query parameter or header)
    userID := r.URL.Query().Get("userID")
    if userID == "" {
        http.Error(w, "userID is required", http.StatusBadRequest)
        return
    }

    // Upgrade initial GET request to a websocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }

    // Add client data
    client := &Client{
        connection:  conn,
        connOwnerId: userID,
        send:        make(chan []byte, 256),
    }

    log.Println("Client ID: ", client.connOwnerId)

    mu.Lock()
    connections[userID] = client
    onlineUsers := make([]string, 0, len(connections))
    for uid := range connections {
        userInt, _ := strconv.Atoi(uid)
        user, _ := utils.GetUsername(userInt)
        onlineUsers = append(onlineUsers, user)
    }

    message := structs.SocketMessage{
        Type:        "online_users",
        OnlineUsers: onlineUsers,
    }

    broadcast <- message
    mu.Unlock()

    go handleClientRead(conn, client, userID)
    go handleClientWrite(client, userID)
}

func handleClientRead(conn *websocket.Conn, client *Client, userID string) {
    defer func() {
        // Cleanup: Close the WebSocket connection when the goroutine exits
        conn.Close()

        // Remove the client from the connections map
        mu.Lock()
        delete(connections, userID)
        mu.Unlock()
    }()
    for {
        // Read message from browser
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("error reading message", err)
            break
        }

        var msg structs.Message
        err = json.Unmarshal(message, &msg)
        if err != nil {
            log.Println("Error unmarshalling message:", err)
            continue
        }

        senderID, _ := strconv.Atoi(userID)
        prepMessage(&msg, senderID) // Add timestamp and sender

        mu.Lock()
        recipientConn, ok := connections[msg.Recipient]
        mu.Unlock()

        if ok {
            sendMessage(recipientConn, msg)
        } else {
            log.Println("Recipient is not connected")
        }
    }
}

func handleClientWrite(client *Client, userID string) {
    defer func() {
        client.connection.Close()
        mu.Lock()
        delete(connections, userID)
        mu.Unlock()
    }()

    // Listen on the send channel for outgoing messages
    for message := range client.send {
        client.mu.Lock()
        err := client.connection.WriteMessage(websocket.TextMessage, message)
        client.mu.Unlock()

        if err != nil {
            log.Println("Error writing message:", err)
            return
        }
    }
}


func prepMessage(m *structs.Message, senderID int) {
	m.Sender = senderID
	m.CreatedAt = time.Now()
}

func sendMessage(recipientConn *Client, msg structs.Message) {
	// Convert message to JSON
	messageData, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	// Send the message into the recipient's send channel (non-blocking)
	select {
	case recipientConn.send <- messageData:
		utils.SaveMessage(msg)
		// Successfully added the message to the recipient's send channel
	default:
		// Channel is full or blocked, handle disconnect, or log a warning
		log.Println("Recipient channel is full, disconnecting")
		close(recipientConn.send)
	}
}

func broadcastOnlineUsers() {
	for {
		message := <-broadcast

		messageJSON, _ := json.Marshal(message.OnlineUsers)

		log.Println(message)

		mu.Lock()
		for _, client := range connections {
			select {
			case client.send <- []byte(messageJSON):
			default:
				close(client.send)
				delete(connections, client.connOwnerId)
			}
		}
		mu.Unlock()
	}
}

// Only accept connections from localhost:4000
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "localhost:4000/":
		return true
	default:
		return false
	}
}
