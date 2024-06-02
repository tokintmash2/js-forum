package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-forum/structs"
	"real-forum/utils"
	"strconv"
)

func HomeJSONHandler(writer http.ResponseWriter, request *http.Request) {
	sessionCookie, err := request.Cookie("session")
	loggedIn := false
	var username string

	if err == nil {
		sessionUUID := sessionCookie.Value
		userID, validSession := utils.VerifySession(sessionUUID)
		if validSession {
			loggedIn = true

			// Fetch username for the logged-in user
			username, err = utils.GetUsername(userID)
			if err != nil {
				http.Error(writer, "Error fetching username", http.StatusInternalServerError)
				return
			}
		}
	}

	// Fetch recent posts
	recentPosts, err := utils.GetRecentPosts()
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
	// allCategories, err := utils.GetCategories()
	// if err != nil {
	// 	http.Error(writer, "Error fetching categories", http.StatusInternalServerError)
	// 	return
	// }

	jsonCategories, err := utils.GetCategories()
	if err != nil {
		http.Error(writer, "Error fetching categories", http.StatusInternalServerError)
		return
	}

	var allCategories []structs.Category
	err = json.Unmarshal([]byte(jsonCategories), &allCategories)
	if err != nil {
		fmt.Println("Error unmarshaling in HomeHandler", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare data for rendering the template
	data := struct {
		LoggedIn    bool
		Username    string
		RecentPosts []structs.PostDetails
		Categories  []structs.Category
	}{
		LoggedIn:    loggedIn,
		Username:    username,
		RecentPosts: recentPosts,
		Categories:  allCategories,
	}

	// Render the template using the provided data
	// err = utils.Templates.ExecuteTemplate(writer, "base", data)
	// if err != nil {
	// 	http.Error(writer, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonData)

}
