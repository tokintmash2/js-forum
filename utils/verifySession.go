package utils

import (
	"forum-auth/database"
	"time"
)

// VerifySession checks if the provided session UUID exists and is still valid
func VerifySession(requestUUID string) (int, bool) {
	var userID int
	var expiry time.Time

	// Query the database to fetch user ID and session expiry based on the session UUID
	err := database.DB.QueryRow(`
		SELECT user_id, expiry FROM sessions
		WHERE session_uuid = ?`,
		requestUUID,
	).Scan(&userID, &expiry)

	// If there's an error during the query (session not found or other issues), return false
	if err != nil {
		return 0, false
	}

	// Check if the session expiry is after the current time to validate the session
	return userID, expiry.After(time.Now())
}
