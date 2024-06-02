package utils

import (
	"database/sql"
	"fmt"
	"real-forum/database"
)

// GetUsername retrieves the username for a given user ID
func GetUsername(userID int) (string, error) {
	var username sql.NullString

	err := database.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		// Check for 'no rows in result set' error specifically
		if err == sql.ErrNoRows {
			return "Unknown", nil
		}
		return "", fmt.Errorf("error fetching username for userID %d: %w", userID, err)
	}

	// Check if the username is NULL before assigning
	if username.Valid {
		return username.String, nil
	}

	return "Unknown", nil
}
