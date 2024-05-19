package handlers

import (
	"forum-auth/database"
	"forum-auth/utils"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)

// CreatePostHandler manages post creation functionality
func CreatePostHandler(writer http.ResponseWriter, request *http.Request) {
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

	// Fetch available categories from the database
	categories, err := utils.GetCategories()
	if err != nil {
		http.Error(writer, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []utils.Category
	}{
		Categories: categories,
	}

	if request.Method == http.MethodPost {
		err := request.ParseForm()
		if err != nil {
			log.Printf("Error parsing form data: %v\n", err)
			http.Error(writer, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		title := request.FormValue("title")
		content := request.FormValue("content")

		if title == "" || content == "" {
			redirectTo404(writer, request)
			return
		}

		// Retrieve selected category IDs from the form
		categoryIDs := request.Form["categoryIDs"]
		// Convert category IDs to integers
		selectedCategoryIDs := convertToIntSlice(categoryIDs)

		newPost := Post{
			UserID:      userID,
			Title:       title,
			Content:     content,
			CreatedAt:   time.Now(),
			CategoryIDs: selectedCategoryIDs,
		}

		err = createPost(newPost)
		if err != nil {
			log.Printf("Error creating post: %v\n", err)
			http.Error(writer, "Error creating post", http.StatusInternalServerError)
			return
		}

		// Check if the form contains comment data
		commentContent := request.FormValue("comment_content")
		if commentContent != "" {
			postID, err := getLastInsertedPostID()
			if err != nil {
				http.Error(writer, "Error retrieving post ID", http.StatusInternalServerError)
				return
			}

			// Fetch username for the logged-in user
			username, err := getUsername(userID)
			if err != nil {
				http.Error(writer, "Error fetching username", http.StatusInternalServerError)
				return
			}

			newComment := utils.Comment{
				UserID:    userID,
				PostID:    postID,
				Content:   commentContent,
				CreatedAt: time.Now(),
				Author:    username,
			}

			err = utils.CreateComment(newComment)
			if err != nil {
				http.Error(writer, "Error creating comment", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	tmplFile := path.Join("templates", "create_post.html")
	t, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(writer, data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

// getLastInsertedPostID retrieves the ID of the last inserted post
func getLastInsertedPostID() (int, error) {
	var postID int
	err := database.DB.QueryRow("SELECT last_insert_rowid()").Scan(&postID)
	if err != nil {
		return 0, err
	}
	return postID, nil
}

// convertToIntSlice converts a string slice to an integer slice
func convertToIntSlice(strSlice []string) []int {
	intSlice := make([]int, len(strSlice))
	for i, str := range strSlice {
		num, err := strconv.Atoi(str)
		if err != nil {
			continue
		}
		intSlice[i] = num
	}
	return intSlice
}

// Post holds information about a post
type Post struct {
	ID          int
	PostID      int
	UserID      int
	Title       string
	Content     string
	CreatedAt   time.Time
	Author      string
	Likes       int
	Dislikes    int
	CategoryIDs []int
	Comments    []utils.Comment
}

// createPost creates a new post in the database
func createPost(newPost Post) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert post
	_, err = tx.Exec(`
        INSERT INTO posts (user_id, title, content, created_at)
        VALUES (?, ?, ?, ?)`,
		newPost.UserID, newPost.Title, newPost.Content, newPost.CreatedAt,
	)
	if err != nil {
		log.Printf("Error inserting post: %v\n", err)
		return err
	}

	// Retrieve last inserted post ID
	var postID int
	err = tx.QueryRow("SELECT last_insert_rowid()").Scan(&postID)
	if err != nil {
		log.Printf("Error getting last inserted post ID: %v\n", err)
		return err
	}

	// Insert associated categories
	for _, categoryID := range newPost.CategoryIDs {
		_, err = tx.Exec(`
            INSERT INTO post_categories (post_id, category_id)
            VALUES (?, ?)
        `, postID, categoryID)
		if err != nil {
			log.Printf("Error associating category: %v\n", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return err
	}

	return nil
}
