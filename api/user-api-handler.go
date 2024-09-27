package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-forum/structs"
	"real-forum/utils"
	"time"
)

// var userIDWS int

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user structs.User
	json.NewDecoder(r.Body).Decode(&user)

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
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "sessionUUID": sessionUUID, "userID": userID})
		// userIDWS = userID
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Wrong email or password"})
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	var newUser structs.User
	json.NewDecoder(r.Body).Decode(&newUser)

	fmt.Println(newUser)

	err := utils.CreateUser(newUser)
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
	cookie, err := request.Cookie("session")
	if err != nil {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	sessionUUID := cookie.Value
	userID, validSession := utils.VerifySession(sessionUUID, "LogoutHandler")
	if !validSession {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	err = utils.SetUserOffline(userID)
	if err != nil {
		http.Error(writer, "Error setting user offline", http.StatusInternalServerError)
		return
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
