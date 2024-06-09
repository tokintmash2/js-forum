package utils

import (
	"log"
	"real-forum/database"
	"real-forum/structs"
	"strconv"
)

// getLastInsertedPostID retrieves the ID of the last inserted post
func GetLastInsertedPostID() (int, error) {
	var postID int
	err := database.DB.QueryRow("SELECT last_insert_rowid()").Scan(&postID)
	if err != nil {
		return 0, err
	}
	return postID, nil
}

// convertToIntSlice converts a string slice to an integer slice
func ConvertToIntSlice(strSlice []string) []int {
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

// createPost creates a new post in the database
func CreatePost(newPost structs.Post) error {
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
