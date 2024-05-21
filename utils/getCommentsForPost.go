package utils

import (
	"database/sql"
	"real-forum/database"
)

// GetCommentsForPost retrieves comments associated with a specific post ID
func GetCommentsForPost(postID int) ([]Comment, error) {
	// Query to fetch comments associated with a post including likes, dislikes, and author details
	rows, err := database.DB.Query(`
        SELECT c.id, c.user_id, c.post_id, c.content, c.created_at,
        COALESCE(l.likes, 0) as comment_likes,
        COALESCE(d.dislikes, 0) as comment_dislikes,
        u.username
        FROM comments c
        LEFT JOIN (
            SELECT comment_id, COUNT(*) as likes FROM likes WHERE is_comment_like = 1 GROUP BY comment_id
        ) l ON c.id = l.comment_id
        LEFT JOIN (
            SELECT comment_id, COUNT(*) as dislikes FROM likes WHERE is_comment_like = 0 GROUP BY comment_id
        ) d ON c.id = d.comment_id
        LEFT JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var nullUsername sql.NullString
		if err := rows.Scan(
			&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.CreatedAt,
			&comment.Likes, &comment.Dislikes, &nullUsername,
		); err != nil {
			return nil, err
		}
		// Check if username is NULL before assigning
		if nullUsername.Valid {
			comment.Author = nullUsername.String
		} else {
			comment.Author = "Unknown"
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
