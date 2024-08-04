package api

import (
	"encoding/json"
	"net/http"
	"real-forum/structs"
	"real-forum/utils"
)

// MyPostsHandler handles requests to display posts of a logged-in user
func MyPostsApiHandler(writer http.ResponseWriter, request *http.Request) {
	// Check if the user is logged in via session cookie
	sessionCookie, err := request.Cookie("session")
	loggedIn := false

	if err == nil {
		sessionUUID := sessionCookie.Value
		userID, validSession := utils.VerifySession(sessionUUID, "MyPostsApiHandler")
		if validSession {
			loggedIn = true

			// Fetch posts of the logged-in user
			userPosts, err := utils.GetUserPostsWithComments(userID)
			if err != nil {
				http.Error(writer, "Error fetching user posts: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Pass userPosts data to the template
			data := struct {
				LoggedIn  bool
				UserPosts []structs.PostDetails
			}{
				LoggedIn:  loggedIn,
				UserPosts: userPosts,
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				http.Error(writer, "Error generating JSON", http.StatusInternalServerError)
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			writer.Write(jsonData)
		}
	}

	// If not logged in, redirect to sign-in page
	// http.Redirect(writer, request, "/sign-in-form", http.StatusSeeOther)
}
