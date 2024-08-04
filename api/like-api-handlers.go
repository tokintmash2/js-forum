package api

import (
	"encoding/json"
	"net/http"
	"real-forum/utils"
)

// LikePostHandler handles the like functionality for a post
func LikePostHandler(writer http.ResponseWriter, request *http.Request) {
	// Check user authentication by verifying the session
	cookie, err := request.Cookie("session")
	if err != nil {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	sessionUUID := cookie.Value
	userID, validSession := utils.VerifySession(sessionUUID, "LikePostHandler")
	if !validSession {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	if request.Method == http.MethodPost {
		var payload struct {
			PostID int `json:"post_id"`
		}

		err := json.NewDecoder(request.Body).Decode(&payload)
		if err != nil {
			http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		postID := payload.PostID

		alreadyLiked, err := utils.CheckIfPostLiked(userID, postID)
		if err != nil {
			http.Error(writer, "Error checking if post is already liked", http.StatusInternalServerError)
			return
		}

		if alreadyLiked {
			// Remove the like
			err = utils.RemovePostLike(userID, postID)
			if err != nil {
				http.Error(writer, "Error removing like", http.StatusInternalServerError)
				return
			}
		} else {
			// Add the like
			err = utils.LikePost(userID, postID)
			if err != nil {
				http.Error(writer, "Error liking post", http.StatusInternalServerError)
				return
			}
		}

		newLikeCount, newDislikeCount, err := utils.GetPostLikeCount(postID)
		if err != nil {
			http.Error(writer, "Error retrieving new like count", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":         true,
			"newLikeCount":    newLikeCount,
			"newDislikeCount": newDislikeCount,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
		return
	}
}

func DislikePostHandler(writer http.ResponseWriter, request *http.Request) {
	// Check user authentication by verifying the session
	cookie, err := request.Cookie("session")
	if err != nil {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	sessionUUID := cookie.Value
	userID, validSession := utils.VerifySession(sessionUUID, "DislikePostHandler")
	if !validSession {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	if request.Method == http.MethodPost {
		var payload struct {
			PostID int `json:"post_id"`
		}

		err := json.NewDecoder(request.Body).Decode(&payload)
		if err != nil {
			http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		postID := payload.PostID

		alreadyDisliked, err := utils.CheckIfPostDisliked(userID, postID)
		if err != nil {
			http.Error(writer, "Error checking if post is already disliked", http.StatusInternalServerError)
			return
		}

		if alreadyDisliked {
			// Remove the dislike
			err = utils.RemovePostDislike(userID, postID)
			if err != nil {
				http.Error(writer, "Error removing dislike", http.StatusInternalServerError)
				return
			}
		} else {
			// Add the dislike
			err = utils.DislikePost(userID, postID)
			if err != nil {
				http.Error(writer, "Error disliking post", http.StatusInternalServerError)
				return
			}
		}

		newLikeCount, newDislikeCount, err := utils.GetPostLikeCount(postID)
		if err != nil {
			http.Error(writer, "Error retrieving new like count", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":         true,
			"newLikeCount":    newLikeCount,
			"newDislikeCount": newDislikeCount,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
		return
	}
}
