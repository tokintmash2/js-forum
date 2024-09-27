package structs

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Chat message
type Message struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	Sender    int
	CreatedAt time.Time
}

// User represents user data
type User struct {
	ID       int
	Email    string
	Password string
	Username string
	GitHubID string
	GoogleID string
}

type SocketMessage struct {
	Type string
	OnlineUsers []string
}

type Client struct {
	Connection  *websocket.Conn
	Send        chan []byte
	Mu          sync.Mutex
	ConnOwnerId string
	LastActive  time.Time
}

type PostDetails struct {
	ID         int
	UserID     int
	Title      string
	Content    string
	CreatedAt  time.Time
	Author     string
	Likes      int
	Dislikes   int
	LoggedIn   bool
	Comments   []Comment
	CategoryID int
}

type Post struct {
	ID          int
	PostID      int
	UserID      int
	Title       string
	Content     string
	CreatedAt   time.Time
	Author      string
	Likes       int
	Dislikes    int
	CategoryIDs []int
	Comments    []Comment
}

type Comment struct {
	ID        int
	UserID    int
	PostID    int
	Content   string
	CreatedAt time.Time
	Author    string
	Likes     int
	Dislikes  int
	LoggedIn  bool
}

type Category struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}
