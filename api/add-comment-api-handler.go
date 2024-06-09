package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-forum/structs"
	"real-forum/utils"
	"time"
)

// AddCommentHandler handles adding a comment to a post
func AddCommentApiHandler(writer http.ResponseWriter, request *http.Request) {

	fmt.Println("AdCOmment Go Handler triggered")

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
		var comment struct {
			Content string `json:"content"`
			PostID  int    `json:"post_id"`
		}

		fmt.Println(request.Body)
		
		err := json.NewDecoder(request.Body).Decode(&comment)
		if err != nil {
			http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		

		// Get comment content and post ID from the form
		// commentContent := request.FormValue("comment_content")
		// postIDStr := request.FormValue("post_id")
		// postID, err := strconv.Atoi(postIDStr)
		// if err != nil {
		// 	http.Error(writer, "Invalid post ID", http.StatusBadRequest)
		// 	return
		// }

		// Ensure comment content is not empty
		if comment.Content == "" {
			RedirectTo404(writer, request)
			return
		}

		// Create a new comment associated with the post
		newComment := structs.Comment{
			UserID:    userID,
			PostID:    comment.PostID,
			Content:   comment.Content,
			CreatedAt: time.Now(),
		}

		fmt.Println("Comment data:", newComment)

		err = utils.CreateComment(newComment)
		if err != nil {
			http.Error(writer, "Error creating comment", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
		return
	}

	// Handle cases where the request method is not POST
	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
}

// RedirectTo404 redirects the user to the 404 page
func RedirectTo404(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/404", http.StatusSeeOther)
}
