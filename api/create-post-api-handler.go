package api

import (
	"encoding/json"
	"log"
	"net/http"
	"real-forum/structs"
	"real-forum/utils"
	"time"
)

// CreatePostHandler manages post creation functionality
func CreatePostApiHandler(writer http.ResponseWriter, request *http.Request) {
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
		var newPost struct {
			Title       string `json:"title"`
			Content     string `json:"content"`
			CategoryIDs []int  `json:"categoryIDs"`
		}

		err := json.NewDecoder(request.Body).Decode(&newPost)
		if err != nil {
			http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if newPost.Title == "" || newPost.Content == "" || len(newPost.CategoryIDs) < 1 {
			http.Error(writer, "Title, content or Category cannot be empty", http.StatusBadRequest)
			return
		}

		post := structs.Post{
			UserID:      userID,
			Title:       newPost.Title,
			Content:     newPost.Content,
			CreatedAt:   time.Now(),
			CategoryIDs: newPost.CategoryIDs,
		}

		err = utils.CreatePost(post)
		if err != nil {
			log.Printf("Error creating post: %v\n", err)
			http.Error(writer, "Error creating post", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
		return
	}

	http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
}
