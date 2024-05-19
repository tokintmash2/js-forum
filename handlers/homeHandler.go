package handlers

import (
	"database/sql"
	"fmt"
	"forum-auth/database"
	"forum-auth/utils"
	"net/http"
	"strconv"
	"time"
)

// HomePageHandler manages the homepage and handles displaying recent posts
func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
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
				http.Error(writer, "Error fetching username", http.StatusInternalServerError)
				return
			}
		}
	}

	// Fetch recent posts
	recentPosts, err := getRecentPosts()
	if err != nil {
		http.Error(writer, "Error fetching recent posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set LoggedIn for posts and their comments
	for i := range recentPosts {
		recentPosts[i].LoggedIn = loggedIn

		// Fetch comments associated with this post
		comments, err := utils.GetCommentsForPost(recentPosts[i].ID)
		if err != nil {
			http.Error(writer, "Error fetching comments for post "+strconv.Itoa(recentPosts[i].ID)+": "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the LoggedIn field for each comment within the post
		for j := range comments {
			comments[j].LoggedIn = loggedIn
		}

		// Assign comments to the respective post
		recentPosts[i].Comments = comments
	}

	// Fetch all categories
	allCategories, err := utils.GetCategories()
	if err != nil {
		http.Error(writer, "Error fetching categories", http.StatusInternalServerError)
		return
	}

	// Prepare data for rendering the template
	data := struct {
		LoggedIn    bool
		Username    string
		RecentPosts []PostDetails
		Categories  []utils.Category
	}{
		LoggedIn:    loggedIn,
		Username:    username,
		RecentPosts: recentPosts,
		Categories:  allCategories,
	}

	// Render the template using the provided data
	err = utils.Templates.ExecuteTemplate(writer, "base", data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getUsername retrieves the username for a given user ID
func getUsername(userID int) (string, error) {
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

// getRecentPosts retrieves recent posts within the last 7 days
func getRecentPosts() ([]PostDetails, error) {
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
