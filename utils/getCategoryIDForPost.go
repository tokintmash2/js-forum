package utils

import (
	"forum-auth/database"
)

// GetCategoryIDForPost retrieves the category ID associated with a post
func GetCategoryIDForPost(postID int) (int, error) {
	var categoryID int

	err := database.DB.QueryRow("SELECT category_id FROM post_categories WHERE post_id = ?", postID).Scan(&categoryID)
	if err != nil {
		return 0, err
	}

	return categoryID, nil
}
