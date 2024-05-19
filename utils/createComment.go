package utils

import (
	"forum-auth/database"
	"time"
)

// Comment represents the structure of a comment
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

// CreateComment creates a new comment in the database
func CreateComment(newComment Comment) error {
	_, err := database.DB.Exec(`
        INSERT INTO comments (user_id, post_id, content, created_at)
        VALUES (?, ?, ?, ?)`,
		newComment.UserID, newComment.PostID, newComment.Content, newComment.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
