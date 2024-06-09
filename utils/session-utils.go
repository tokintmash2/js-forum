package utils

import (
	"fmt"
	"real-forum/database"
)

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
