package handlers

import (
	"database/sql"
	"net/http"
	"real-forum/database"
	"real-forum/structs"
	"real-forum/utils"
)

// LikedPostsHandler fetches and displays posts liked by the user
func LikedPostsHandler(writer http.ResponseWriter, request *http.Request) {
	sessionCookie, err := request.Cookie("session")
	loggedIn := false

	if err == nil {
		sessionUUID := sessionCookie.Value
		userID, validSession := utils.VerifySession(sessionUUID)
		if validSession {
			loggedIn = true

			likedPosts, err := getLikedPosts(userID)
			if err != nil {
				http.Error(writer, "Error fetching liked posts: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Fetch comments for each liked post
			for i, post := range likedPosts {
				comments, err := utils.GetCommentsForPost(post.ID)
				if err != nil {
					// Log the error to understand what went wrong with fetching comments
					http.Error(writer, "Error fetching comments for post: "+err.Error(), http.StatusInternalServerError)
					return
				}
				likedPosts[i].Comments = comments
			}

			// Pass likedPosts data to template
			data := struct {
				LoggedIn   bool
				LikedPosts []structs.PostDetails
			}{
				LoggedIn:   loggedIn,
				LikedPosts: likedPosts,
			}

			// Use the centralized template rendering from the utils package
			err = utils.Templates.ExecuteTemplate(writer, "liked_posts.html", data)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}
	}

	// If not logged in, redirect to sign-in form
	http.Redirect(writer, request, "/sign-in-form", http.StatusSeeOther)
}

// getLikedPosts retrieves posts liked by the user from the database
func getLikedPosts(userID int) ([]structs.PostDetails, error) {
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
		comments, err := utils.GetCommentsForPost(post.ID)
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
