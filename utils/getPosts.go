package utils

import (
	"database/sql"
	"real-forum/database"
	"real-forum/structs"
	"time"
)

// GetRecentPosts retrieves recent posts within the last 7 days
func GetRecentPosts() ([]structs.PostDetails, error) {
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

	var recentPosts []structs.PostDetails
	for rows.Next() {
		var post structs.PostDetails
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
		comments, err := GetCommentsForPost(post.ID)
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

// GetCategoryPosts retrieves posts associated with a given category ID
func GetCategoryPosts(categoryID int) ([]structs.PostDetails, error) {
	// Fetch posts associated with the given category ID including likes and dislikes
	rows, err := database.DB.Query(`
        SELECT 
            p.id, p.user_id, p.title, p.content, p.created_at,
            u.username,
            COALESCE(l.likes, 0) as likes,
            COALESCE(d.dislikes, 0) as dislikes
        FROM posts p
        INNER JOIN post_categories pc ON p.id = pc.post_id
        LEFT JOIN users u ON p.user_id = u.id
        LEFT JOIN (
            SELECT post_id, COUNT(*) as likes FROM likes WHERE is_post_like = 1 GROUP BY post_id
        ) l ON p.id = l.post_id
        LEFT JOIN (
            SELECT post_id, COUNT(*) as dislikes FROM likes WHERE is_post_like = 0 GROUP BY post_id
        ) d ON p.id = d.post_id
        WHERE pc.category_id = ?
        ORDER BY p.created_at DESC
    `, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryPosts []structs.PostDetails
	for rows.Next() {
		var post structs.PostDetails
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
		comments, err := GetCommentsForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		categoryPosts = append(categoryPosts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categoryPosts, nil
}
