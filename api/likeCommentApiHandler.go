package api

import (
	"encoding/json"
	"net/http"
	"real-forum/utils"
)

// LikeCommentHandler handles liking or unliking a comment
func LikeCommentHandler(writer http.ResponseWriter, request *http.Request) {
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
		var payload struct {
			CommentID int `json:"comment_id"`
		}

		err := json.NewDecoder(request.Body).Decode(&payload)
		if err != nil {
			http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		commentID := payload.CommentID

		alreadyLiked, err := utils.CheckIfCommentLiked(userID, commentID)
		if err != nil {
			http.Error(writer, "Error checking if comment is already liked", http.StatusInternalServerError)
			return
		}

		if alreadyLiked {
			// Remove the like
			err = utils.RemoveCommentLike(userID, commentID)
			if err != nil {
				http.Error(writer, "Error removing like", http.StatusInternalServerError)
				return
			}
		} else {
			// Add the like
			err = utils.LikeComment(userID, commentID)
			if err != nil {
				http.Error(writer, "Error liking comment", http.StatusInternalServerError)
				return
			}
		}

		newLikeCount, newDislikeCount, err := utils.GetCommentLikeCount(commentID)
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

	// Handle cases where the request method is not POST
	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
}

func DislikeCommentHandler(writer http.ResponseWriter, request *http.Request) {
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
		var payload struct {
			CommentID int `json:"comment_id"`
		}

		err := json.NewDecoder(request.Body).Decode(&payload)
		if err != nil {
			http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		commentID := payload.CommentID

		alreadyDisliked, err := utils.CheckIfCommentDisliked(userID, commentID)
		if err != nil {
			http.Error(writer, "Error checking if comment is already disliked", http.StatusInternalServerError)
			return
		}

		if alreadyDisliked {
			// Remove the dislike
			err = utils.RemoveCommentDislike(userID, commentID)
			if err != nil {
				http.Error(writer, "Error removing dislike", http.StatusInternalServerError)
				return
			}
		} else {
			// Add the dislike
			err = utils.DislikeComment(userID, commentID)
			if err != nil {
				http.Error(writer, "Error disliking comment", http.StatusInternalServerError)
				return
			}
		}

		newLikeCount, newDislikeCount, err := utils.GetCommentLikeCount(commentID)
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

	// Handle cases where the request method is not POST
	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
}
