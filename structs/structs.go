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

type GitHubUserData struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type GoogleUserData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
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