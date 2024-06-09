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

// GetLikedPosts retrieves posts liked by the user from the database
func GetLikedPosts(userID int) ([]structs.PostDetails, error) {
	rows, err := database.DB.Query(`
        SELECT 
            p.id, p.user_id, p.title, p.content, p.created_at,
            u.username,
            COALESCE(l.likes, 0) as likes,
            COALESCE(d.dislikes, 0) as dislikes
        FROM posts p
        LEFT JOIN users u ON p.user_id = u.id
        LEFT JOIN (
            SELECT post_id, COUNT(*) as likes FROM likes WHERE is_post_like = 1 AND user_id = ? GROUP BY post_id
        ) l ON p.id = l.post_id
        LEFT JOIN (
            SELECT post_id, COUNT(*) as dislikes FROM likes WHERE is_post_like = 0 GROUP BY post_id
        ) d ON p.id = d.post_id
        WHERE p.id IN (
            SELECT post_id FROM likes WHERE is_post_like = 1 AND user_id = ?
        )
        ORDER BY p.created_at DESC
    `, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likedPosts []structs.PostDetails
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

		likedPosts = append(likedPosts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return likedPosts, nil
}

// GetUserPostsWithComments retrieves user posts along with their comments
func GetUserPostsWithComments(userID int) ([]structs.PostDetails, error) {
	userPosts, err := GetUserPosts(userID)
	if err != nil {
		return nil, err
	}

	// Fetch comments for each user post
	for i := range userPosts {
		comments, err := GetCommentsForPost(userPosts[i].ID)
		if err != nil {
			return nil, err
		}

		for j := range comments {
			comments[j].LoggedIn = true
		}

		userPosts[i].Comments = comments
	}

	return userPosts, nil
}

// GetUserPosts fetches posts of a user from the database
func GetUserPosts(userID int) ([]structs.PostDetails, error) {
	// Query the database to fetch posts for the given user ID
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
        WHERE p.user_id = ?
        ORDER BY p.created_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userPosts []structs.PostDetails
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

		userPosts = append(userPosts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userPosts, nil
}