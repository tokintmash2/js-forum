package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"real-forum/utils"
)

func ChatHandler(writer http.ResponseWriter, request *http.Request) {

	log.Println("ChatHandler triggered")

	// Check if the user is logged in
	sessionCookie, err := request.Cookie("session")
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionUUID := sessionCookie.Value
	userID, validSession := utils.VerifySession(sessionUUID, "ChatHandler")
	if validSession {
		// Extract partner's ID from URL
		pathSegments := strings.Split(request.URL.Path, "/")
		if len(pathSegments) < 4 {
			http.Error(writer, "Invalid URL", http.StatusBadRequest)
			return
		}

		partnerIDStr := pathSegments[3]
		partnerID, err := strconv.Atoi(partnerIDStr)
		if err != nil || partnerID == 0 {
			http.Error(writer, "Invalid partner ID", http.StatusBadRequest)
			return
		}

		log.Println(userID, partnerID)

		data := utils.GetRecentMessages(userID, partnerID)

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(writer, "Error generating JSON", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(jsonData)

		if err != nil {
			fmt.Println("error fetching messages for Sender", userID, ":", err)
			http.Error(writer, "error fetching messages", http.StatusInternalServerError)
			return
		}
	}
}
