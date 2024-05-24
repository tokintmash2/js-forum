package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-forum/handlers"
	"real-forum/utils"
	"strconv"
)

// LikePostHandler handles the like functionality for a post
func LikePostHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("LikePostHandler started")
	// Check user authentication by verifying the session
	cookie, err := request.Cookie("session")
	if err != nil {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	sessionUUID := cookie.Value
	userID, validSession := utils.VerifySession(sessionUUID)
	if !validSession {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	if request.Method == http.MethodPost {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		// Get post ID from form
		postIDStr := request.FormValue("post_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(writer, "Invalid post ID", http.StatusBadRequest)
			return
		}

		alreadyLiked, err := handlers.CheckIfPostLiked(userID, postID)
		if err != nil {
			http.Error(writer, "Error checking if post is already liked", http.StatusInternalServerError)
			return
		}

		if alreadyLiked {
			// Remove the like
			err = handlers.RemovePostLike(userID, postID)
			if err != nil {
				http.Error(writer, "Error removing like", http.StatusInternalServerError)
				return
			}
		} else {
			// Add the like
			err = handlers.LikePost(userID, postID)
			if err != nil {
				http.Error(writer, "Error liking post", http.StatusInternalServerError)
				return
			}
		}

		// // Redirect to the referring page after performing the action
		// referer := request.Header.Get("Referer")
		// if referer != "" {
		// 	http.Redirect(writer, request, referer, http.StatusSeeOther)
		// } else {
		// 	// If no referring page is found, redirect to the homepage ("/")
		// 	http.Redirect(writer, request, "/", http.StatusSeeOther)
		// }
		// return

		response := map[string]interface{}{
			"success":      true,
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
		return
	}
}
