package database

import (
	"fmt"
	"log"
)

// Creates tables for holding the data in forum.db database
func CreateTables() {
	// Create Users table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE,
			username TEXT,
			password TEXT,
			github_id TEXT,
    		google_id TEXT
		)
	`)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	// Create Posts table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal("Error creating posts table:", err)
	}
	// Create Comments table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			post_id INTEGER,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id)
		)
	`)
	if err != nil {
		log.Fatal("Error creating comments table:", err)
	}

	// Create Categories table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		)
	`)
	if err != nil {
		log.Fatal("Error creating categories table:", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS post_categories (
			post_id INTEGER,
			category_id INTEGER,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
		)
	`)
	if err != nil {
		log.Fatal("Error creating post_categories table:", err)
	}

	// Create Likes table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			post_id INTEGER,
			comment_id INTEGER,
			is_post_like INTEGER,
			is_comment_like INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (comment_id) REFERENCES comments(id),
			UNIQUE(user_id, post_id, comment_id, is_post_like, is_comment_like)
		)
	`)
	if err != nil {
		log.Fatal("Error creating likes table:", err)
	}

	// Create Sessions table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			session_uuid TEXT,
			expiry TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal("Error creating sessions table:", err)
	}

	err = InsertCategories()
	if err != nil {
		log.Fatal("Error inserting categories:", err)
	}
}

// InsertCategories inserts predefined categories into the database
func InsertCategories() error {
	categories := []string{
		"Regular Dreams", "Lucid Dreams", "Nightmares", "Night Terrors",
		"Recurring Dreams", "Daydreams", "Fantasy Dreams", "Adventure Dreams",
		"Historic Dreams", "Romantic Dreams", "Prophetic Dreams", "Just-Bizarre Dreams",
	}

	// Begin a transaction
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare the insert statement
	stmt, err := tx.Prepare("INSERT INTO categories (name) VALUES (?)")
	if err != nil {
		return fmt.Errorf("error preparing category insert statement: %w", err)
	}
	defer stmt.Close()

	// Insert categories
	for _, category := range categories {
		count, err := countCategory(category)
		if err != nil {
			return fmt.Errorf("error checking existing category: %w", err)
		}
		if count > 0 {
			continue
		}

		_, err = stmt.Exec(category)
		if err != nil {
			return fmt.Errorf("error inserting category: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// countCategory checks if a category already exists in the database
func countCategory(category string) (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM categories WHERE name = ?", category).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
