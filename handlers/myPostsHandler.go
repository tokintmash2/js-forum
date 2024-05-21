package handlers

import (
	"database/sql"
	"real-forum/database"
	"real-forum/utils"
	"net/http"
)

// MyPostsHandler handles requests to display posts of a logged-in user
func MyPostsHandler(writer http.ResponseWriter, request *http.Request) {
	// Check if the user is logged in via session cookie
	sessionCookie, err := request.Cookie("session")
	loggedIn := false

	if err == nil {
		sessionUUID := sessionCookie.Value
		userID, validSession := utils.VerifySession(sessionUUID)
		if validSession {
			loggedIn = true

			// Fetch posts of the logged-in user
			userPosts, err := getUserPostsWithComments(userID)
			if err != nil {
				http.Error(writer, "Error fetching user posts: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Pass userPosts data to the template
			data := struct {
				LoggedIn  bool
				UserPosts []PostDetails
			}{
				LoggedIn:  loggedIn,
				UserPosts: userPosts,
			}

			// Use the centralized template rendering from the utils package
			err = utils.Templates.ExecuteTemplate(writer, "my_posts.html", data)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}
	}

	// If not logged in, redirect to sign-in page
	http.Redirect(writer, request, "/sign-in-form", http.StatusSeeOther)
}

// getUserPostsWithComments retrieves user posts along with their comments
func getUserPostsWithComments(userID int) ([]PostDetails, error) {
	userPosts, err := getUserPosts(userID)
	if err != nil {
		return nil, err
	}

	// Fetch comments for each user post
	for i := range userPosts {
		comments, err := utils.GetCommentsForPost(userPosts[i].ID)
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

// getUserPosts fetches posts of a user from the database
func getUserPosts(userID int) ([]PostDetails, error) {
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

	var userPosts []PostDetails
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

		userPosts = append(userPosts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userPosts, nil
}
