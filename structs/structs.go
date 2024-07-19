package structs

import (
	"time"
)

// User represents user data
type User struct {
	ID       int
	Email    string
	Password string
	Username string
	GitHubID string
	GoogleID string
}

type Message struct {
    SenderID   int    `json:"sender_id"`
    ReceiverID int    `json:"receiver_id"`
    Content    string `json:"content"`
    Timestamp  string `json:"timestamp"`
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