package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-forum/handlers"
	"real-forum/utils"
	"strconv"
)

// CategoryPostsHandler handles requests for posts within a specific category
func CategoryPostsApiHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract category ID from the URL
	categoryIDStr := request.URL.Path[len("/api/category/"):]
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(writer, "Invalid category ID", http.StatusBadRequest)
		return
	}

	fmt.Println("Requested Path:", request.URL.Path)
	fmt.Println("Requested ID:", categoryID)

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
		CategoryPosts []handlers.PostDetails
	}{
		LoggedIn:      loggedIn,
		Username:      username,
		CategoryPosts: categoryPosts,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonData)
}
