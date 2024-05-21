package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"real-forum/database"
	"real-forum/utils"
	"time"
)

// HomePageHandler manages the homepage and handles displaying recent posts
func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "index.html")

}

// GetUsername retrieves the username for a given user ID
func GetUsername(userID int) (string, error) {
	var username sql.NullString

	err := database.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		// Check for 'no rows in result set' error specifically
		if err == sql.ErrNoRows {
			return "Unknown", nil
		}
		return "", fmt.Errorf("error fetching username for userID %d: %w", userID, err)
	}

	// Check if the username is NULL before assigning
	if username.Valid {
		return username.String, nil
	}

	return "Unknown", nil
}

// PostDetails struct holds information about a post
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

// GetRecentPosts retrieves recent posts within the last 7 days
func GetRecentPosts() ([]PostDetails, error) {
	threshold := time.Now().AddDate(0, 0, -365)
	rows, err := database.DB.Query(`
        SELECT 
            p.id, p.user_id, p.title, p.content, p.created_at,
            u.username,
            COALESCE(l.likes, 0) as likes,
            COALESCE(d.dislikes, 0) as dislikes
        FROM posts p
        LEFT JOIN users u ON p.user_id = u.id
        LEFT JOIN (
            SELECT post_id, COUNT(*) as likes FROM likes WHERE is_post_like = 1 GROUP BY post_id
        ) l ON p.id = l.post_id
        LEFT JOIN (
            SELECT post_id, COUNT(*) as dislikes FROM likes WHERE is_post_like = 0 GROUP BY post_id
        ) d ON p.id = d.post_id
        WHERE p.created_at > ?
        ORDER BY p.created_at DESC
    `, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recentPosts []PostDetails
	for rows.Next() {
		var post PostDetails
		var nullUsername sql.NullString
		if err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt,
			&nullUsername, &post.Likes, &post.Dislikes,
		); err != nil {
			return nil, err
		}
		// Check if username is NULL before assigning
		if nullUsername.Valid {
			post.Author = nullUsername.String
		} else {
			post.Author = "Unknown"
		}

		// Fetch comments associated with this post
		comments, err := utils.GetCommentsForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		recentPosts = append(recentPosts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recentPosts, nil
}
