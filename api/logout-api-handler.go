package api

import (
	"encoding/json"
	"log"
	"net/http"
	"real-forum/utils"
	"time"
)

// LogoutHandler handles user logout by clearing the session
func LogoutHandler(writer http.ResponseWriter, request *http.Request) {
	sessionCookie, err := request.Cookie("session")
	if err == nil {
		sessionUUID := sessionCookie.Value
		err := utils.ClearSession(sessionUUID)
		if err != nil {
			log.Println("Error clearing session:", err)
		}
	}

	// Clear the session cookie entirely by setting it to expire immediately
	http.SetCookie(writer, &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Now(), // Setting the cookie's expiration to immediately expire
	})

	// Return a JSON response instead of redirecting
	response := map[string]interface{}{
		"success": true,
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}
