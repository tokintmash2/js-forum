package handlers

import (
	"forum-auth/utils"
	"net/http"
	"strconv"
	"time"
)

// AddCommentHandler handles adding a comment to a post
func AddCommentHandler(writer http.ResponseWriter, request *http.Request) {
	// Check if the user is authenticated by verifying the session
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
		// Parse form data
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		// Get comment content and post ID from the form
		commentContent := request.FormValue("comment_content")
		postIDStr := request.FormValue("post_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(writer, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Ensure comment content is not empty
		if commentContent == "" {
			redirectTo404(writer, request)
			return
		}

		// Create a new comment associated with the post
		newComment := utils.Comment{
			UserID:    userID,
			PostID:    postID,
			Content:   commentContent,
			CreatedAt: time.Now(),
		}

		err = utils.CreateComment(newComment)
		if err != nil {
			http.Error(writer, "Error creating comment", http.StatusInternalServerError)
			return
		}

		// Redirect the user after adding the comment, back to the referring page
		returnURL := request.Header.Get("Referer") // Get the referring URL
		http.Redirect(writer, request, returnURL, http.StatusSeeOther)
		return
	}

	// Handle cases where the request method is not POST
	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
}

// redirectTo404 redirects the user to the 404 page
func redirectTo404(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/404", http.StatusSeeOther)
}
