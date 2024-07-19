package api

import (
	"encoding/json"
	"net/http"
	"real-forum/structs"
	"real-forum/utils"
)

// LikedPostsHandler fetches and displays posts liked by the user
func LikedPostsApiHandler(writer http.ResponseWriter, request *http.Request) {
	sessionCookie, err := request.Cookie("session")
	loggedIn := false

	if err == nil {
		sessionUUID := sessionCookie.Value
		userID, validSession := utils.VerifySession(sessionUUID)
		if validSession {
			loggedIn = true

			likedPosts, err := utils.GetLikedPosts(userID)
			if err != nil {
				http.Error(writer, "Error fetching liked posts: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Fetch comments for each liked post
			for i, post := range likedPosts {
				comments, err := utils.GetCommentsForPost(post.ID)
				if err != nil {
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

			jsonData, err := json.Marshal(data)
			if err != nil {
				http.Error(writer, "Error generating JSON", http.StatusInternalServerError)
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			writer.Write(jsonData)
		}
	}

	// If not logged in, redirect to sign-in form
	// http.Redirect(writer, request, "/sign-in-form", http.StatusSeeOther)
}
