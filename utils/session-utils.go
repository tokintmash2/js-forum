package utils

import (
	"fmt"
	"net/http"
	"real-forum/database"
	"time"

	"github.com/google/uuid"
)

// createSession creates a new session for a user
func CreateSession(userID int) (string, error) {
	// Check if there's an existing active session for the user
	var existingSessionID string
	err := database.DB.QueryRow(`
        SELECT session_uuid FROM sessions
        WHERE user_id = ? AND expiry > CURRENT_TIMESTAMP`,
		userID,
	).Scan(&existingSessionID)

	if err == nil {
		// If an active session exists, clear it before creating a new one
		err = ClearSession(existingSessionID)
		if err != nil {
			return "", err
		}
	}

	// Create a new session
	sessionUUID := uuid.New().String()
	expiry := time.Now().Add(24 * time.Hour) // Set the session expiry time (24 hours)
	_, err = database.DB.Exec(`
        INSERT INTO sessions (user_id, session_uuid, expiry)
        VALUES (?, ?, ?)`,
		userID, sessionUUID, expiry,
	)
	if err != nil {
		return "", err
	}
	return sessionUUID, nil
}

// createSessionCookie creates a cookie for the session
func CreateSessionCookie(sessionUUID string) *http.Cookie {
	expiration := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionUUID,
		Path:     "/",
		Expires:  expiration,
		// HttpOnly: true,
	}
	return cookie
}

// clearSession clears a session by its UUID
func ClearSession(sessionUUID string) error {
	fmt.Println("ClearSession triggered")
	_, err := database.DB.Exec(`
		DELETE FROM sessions
		WHERE session_uuid = ?`,
		sessionUUID,
	)
	if err != nil {
		return err
	}
	return nil
}

// VerifySession checks if the provided session UUID exists and is still valid
func VerifySession(requestUUID string, caller string) (int, bool) {

	var userID int
	var expiry time.Time

	fmt.Println("VerifySession:", requestUUID, "caller:", caller)

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
