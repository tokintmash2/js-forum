package handlers

import (
	"database/sql"
	"fmt"
	"forum-auth/database"
	"forum-auth/utils"
	"net/http"
	"strconv"
)

// CategoryPostsHandler handles requests for posts within a specific category
func CategoryPostsHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract category ID from the URL
	categoryIDStr := request.URL.Path[len("/category/"):]
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(writer, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Check if the user is logged in
	sessionCookie, err := request.Cookie("session")
	loggedIn := false
	var username string

	if err == nil {
		sessionUUID := sessionCookie.Value
		userID, validSession := utils.VerifySession(sessionUUID)
		if validSession {
			loggedIn = true

			// Fetch username for the logged-in user
			username, err = getUsername(userID)
			if err != nil {
				fmt.Println("Error fetching username for userID", userID, ":", err)
				http.Error(writer, "Error fetching username", http.StatusInternalServerError)
				return
			}
		}
	}

	// Fetch posts associated with the given category ID including likes and dislikes
	categoryPosts, err := getCategoryPosts(categoryID)
	if err != nil {
		http.Error(writer, "Error fetching category posts", http.StatusInternalServerError)
		return
	}

	// Set LoggedIn for category posts and their comments
	for i := range categoryPosts {
		categoryPosts[i].LoggedIn = loggedIn

		// Fetch comments associated with this post
		comments, err := utils.GetCommentsForPost(categoryPosts[i].ID)
		if err != nil {
			http.Error(writer, "Error fetching comments for post "+strconv.Itoa(categoryPosts[i].ID)+": "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the LoggedIn field for each comment within the post
		for j := range comments {
			comments[j].LoggedIn = loggedIn
		}

		// Assign comments to the respective post
		categoryPosts[i].Comments = comments
	}

	data := struct {
		LoggedIn      bool
		Username      string
		CategoryPosts []PostDetails
	}{
		LoggedIn:      loggedIn,
		Username:      username,
		CategoryPosts: categoryPosts,
	}

	// Set the content type before executing the template
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render the category_posts.html template with the fetched posts
	err = utils.Templates.ExecuteTemplate(writer, "category_posts.html", data)
	if err != nil {
		fmt.Println("Template execution error:", err)
		return
	}
}

// getCategoryPosts retrieves posts associated with a given category ID
func getCategoryPosts(categoryID int) ([]PostDetails, error) {
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

	var categoryPosts []PostDetails
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

		categoryPosts = append(categoryPosts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categoryPosts, nil
}
