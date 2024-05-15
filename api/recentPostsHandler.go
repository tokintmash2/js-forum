package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-forum/database"
	"real-forum/structs"
	"real-forum/utils"
	"time"
)

// PostDetails struct holds information about a post
// type PostDetails struct {
// 	ID         int
// 	UserID     int
// 	Title      string
// 	Content    string
// 	CreatedAt  time.Time
// 	Author     string
// 	Likes      int
// 	Dislikes   int
// 	LoggedIn   bool
// 	Comments   []utils.Comment
// 	CategoryID int
// }

func RecentPostsHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}
	defer rows.Close()

	var recentPosts []structs.PostDetails
	for rows.Next() {
		var post structs.PostDetails
		var nullUsername sql.NullString
		if err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt,
			&nullUsername, &post.Likes, &post.Dislikes,
		); err != nil {
			return
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
			return
		}
		post.Comments = comments

		recentPosts = append(recentPosts, post)
	}

	if err := rows.Err(); err != nil {
		return
	}

	jsonData, err := json.Marshal(recentPosts)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}
