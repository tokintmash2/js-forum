package api

import (
	"fmt"
	"net/http"
	"real-forum/handlers"
	"real-forum/utils"
	"strconv"
)

// CategoryPostsHandler handles requests for posts within a specific category
func CategoryPostsApiHandler(writer http.ResponseWriter, request *http.Request) {
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
			username, err = handlers.GetUsername(userID)
			if err != nil {
				fmt.Println("Error fetching username for userID", userID, ":", err)
				http.Error(writer, "Error fetching username", http.StatusInternalServerError)
				return
			}
		}
	}

	// Fetch posts associated with the given category ID including likes and dislikes
	categoryPosts, err := handlers.GetCategoryPosts(categoryID)
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
