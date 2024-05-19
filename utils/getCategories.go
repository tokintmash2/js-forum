package utils

import (
	"forum-auth/database"
	"log"
)

// Category represents a post category
type Category struct {
	ID   int
	Name string
}

// GetCategories fetches all categories from the database
func GetCategories() ([]Category, error) {
	var categories []Category

	rows, err := database.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		log.Fatal("Error fetching categories:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Fatal("Error scanning category row:", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating through categories:", err)
		return nil, err
	}

	return categories, nil
}
