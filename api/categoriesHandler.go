package api

import (
	"encoding/json"
	"log"
	"net/http"
	"real-forum/database"
	"real-forum/utils"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {

	var categories []utils.Category

	rows, err := database.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		log.Fatal("Error fetching categories:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category utils.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Fatal("Error scanning category row:", err)
			return
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating through categories:", err)
		return
	}

	jsonData, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
