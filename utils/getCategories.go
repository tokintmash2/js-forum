package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"real-forum/database"
	"real-forum/structs"
)

// GetCategories fetches all categories from the database
func GetCategories() (string, error) {
	var categories []structs.Category

	rows, err := database.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		log.Fatal("Error fetching categories:", err)
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Fatal("Error scanning category row:", err)
			return "", err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating through categories:", err)
		return "", err
	}

	jsonData, err := json.Marshal(categories)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}
	// fmt.Println(categories)
	return string(jsonData), nil
}
