package utils

import (
	"real-forum/database"
	"real-forum/structs"
)

// CreateComment creates a new comment in the database
func CreateComment(newComment structs.Comment) error {
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
