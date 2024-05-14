package handlers

import (
	"log"
	"net/http"
	"time"
)

// LogoutHandler handles user logout by clearing the session
func LogoutHandler(writer http.ResponseWriter, request *http.Request) {
	sessionCookie, err := request.Cookie("session")
	if err == nil {
		sessionUUID := sessionCookie.Value
		err := clearSession(sessionUUID)
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

	// Redirect to the home page after logout
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
