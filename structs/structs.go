package structs

import (
	"real-forum/utils"
	"time"
)

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
	Comments   []utils.Comment
	CategoryID int
}
