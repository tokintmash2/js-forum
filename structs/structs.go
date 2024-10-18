package structs

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Chat message
type Message struct {
	Type             string     `json:"type"`
	OnlineUsers      []UserInfo `json:"online_users"`
	Recipient        string     `json:"recipient"`
	SenderUsername   string     `json:"sender_username"`
	ReceiverUsername string     `json:"receiver_username"`
	Content          string     `json:"content"`
	Sender           int        `json:"sender"`
	CreatedAt        time.Time  `json:"created_at"`
}

// User represents user data
type User struct {
	ID         int
	Email      string
	Password   string
	Username   string
	FirstName  string
	LastName   string
	Age        int
	Gender     string
	Identifier string
}

type SocketMessage struct {
	Type        string     `json:"type"`
	OnlineUsers []UserInfo `json:"online_users"`
	Content     string     `json:"content"`
}

type UserInfo struct {
	ID       int    `json:"ID"`
	Username string `json:"Username"`
}

// type SocketMessage struct {
// 	Type string
// 	OnlineUsers []UserInfo
// }

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
