package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-forum/structs"
	"real-forum/utils"
	"time"
)

// var userIDWS int

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user structs.User
	json.NewDecoder(r.Body).Decode(&user)

	fmt.Printf("LoginHandler user: %+v\n", user)

	userID, verified := utils.VerifyUser(user)
	if verified {
		err := utils.SetUserOnline(userID)
		if err != nil {
			http.Error(w, "Error setting user online", http.StatusInternalServerError)
			return
		}
		sessionUUID, err := utils.CreateSession(userID)
		if err != nil {
			http.Error(w, "Failed to create a session", http.StatusInternalServerError)
			return
		}
		cookie := utils.CreateSessionCookie(sessionUUID)
		http.SetCookie(w, cookie)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "sessionUUID": sessionUUID})
		// userIDWS = userID
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Wrong email or password"})
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	var newUser structs.User
	json.NewDecoder(r.Body).Decode(&newUser)

	// Create the new user
	userID, err := utils.CreateUser(newUser)
	// Set the newly generated user ID
	newUser.ID = userID
	fmt.Println(newUser)
	if err == nil {
		sessionUUID, err := utils.CreateSession(newUser.ID)
		if err != nil {
			http.Error(w, "Failed to create a session", http.StatusInternalServerError)
			return
		}
		cookie := utils.CreateSessionCookie(sessionUUID)
		http.SetCookie(w, cookie)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
	}
}

// LogoutHandler handles user logout by clearing the session
func LogoutHandler(writer http.ResponseWriter, request *http.Request) {

	log.Println("LogoutHandler called")

	cookie, err := request.Cookie("session")
	if err != nil {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	// Verify session UUID
	sessionUUID := cookie.Value
	userID, validSession := utils.VerifySession(sessionUUID, "LogoutHandler")
	if !validSession {
		// If the session is not valid, redirect to the sign-in page
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	// Set user offline in the database
	err = utils.SetUserOffline(userID)
	if err != nil {
		http.Error(writer, "Error setting user offline", http.StatusInternalServerError)
		return
	}

	// Clear the session cookie by setting it to expire
	clearedCookie := &http.Cookie{
		Name:    "session",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,                             // Set MaxAge to -1 to immediately expire the cookie
		Expires: time.Now().Add(-1 * time.Hour), // Set the cookie's expiration to a past date
	}

	// Clear the session cookie entirely by setting it to expire immediately
	// http.SetCookie(writer, &http.Cookie{
	// 	Name:    "session",
	// 	Value:   "",
	// 	Expires: time.Now(), // Setting the cookie's expiration to immediately expire
	// })

	http.SetCookie(writer, clearedCookie)

	log.Println("Cookie value:", cookie.Value)

	// Return a JSON response instead of redirecting
	response := map[string]interface{}{
		"success": true,
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)

	// Redirect to the sign-in page after successfully logging out
	// http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
}

func GetOnlineUsersHandler(writer http.ResponseWriter, request *http.Request) {
	users, err := utils.GetOnlineUsers()
	if err != nil {
		http.Error(writer, "Error retrieving online users", http.StatusInternalServerError)
		return
	}

	// fmt.Println("Online users from OnlineUsersHandler", users)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}
