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
		Expires:  expiration,
		HttpOnly: true,
	}
	return cookie
}

// clearSession clears a session by its UUID
func ClearSession(sessionUUID string) error {
	fmt.Println("ClearSessino triggered")
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
