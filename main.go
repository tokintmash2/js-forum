package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"real-forum/api"
	"real-forum/database"
	"real-forum/handlers"
	"real-forum/utils"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var upgrader = websocket.Upgrader{}
var db *sql.DB

func main() {
	var err error

	// Connect to the SQLite database
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal("Error while connecting to the database:", err)
	}
	defer db.Close() // Close the database connection when main() exits

	// Initialize the forum database and create necessary tables
	database.ConnectToForumDB()
	database.CreateTables()

	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Serve static files in the /static directory
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Define HTTP request handlers for various endpoints
	mux.HandleFunc("/", handlers.NotFoundWrapper(handlers.HomePageHandler))

	// -------- JAvaScript API ---------
	mux.HandleFunc("/ws", wsHandler)
	mux.HandleFunc("/api/categories", api.CategoriesHandler)
	mux.HandleFunc("/api/recents", api.RecentPostsHandler)
	mux.HandleFunc("/api/home", api.HomeJSONHandler)
	mux.HandleFunc("/api/category/", api.CategoryPostsApiHandler)
	mux.HandleFunc("/api/likePost", api.LikePostHandler)
	mux.HandleFunc("/api/dislikePost", api.DislikePostHandler)
	mux.HandleFunc("/api/likeComment", api.LikeCommentHandler)
	mux.HandleFunc("/api/dislikeComment", api.DislikeCommentHandler)
	mux.HandleFunc("/api/logout", api.LogoutHandler)
	mux.HandleFunc("/api/create-post", api.CreatePostApiHandler)
	mux.HandleFunc("/api/my-posts", api.MyPostsApiHandler)
	mux.HandleFunc("/api/liked-posts", api.LikedPostsApiHandler)
	mux.HandleFunc("/api/add-comment", api.AddCommentApiHandler)
	mux.HandleFunc("/api/login", api.LoginHandler)
	mux.HandleFunc("/api/signup", api.SignupHandler)
	mux.HandleFunc("/api/online-user", api.GetOnlineUsersHandler)

	// Start the server on port 4000 and log its status
	log.Println("Server started at :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	sessionCookie, err := r.Cookie("session")
	if err != nil {
		log.Println("Session cookie not foundz:", err)
		return
	}

	if sessionCookie == nil {
		log.Println("Session cookie is nil")
		return
	}

	sessionUUID := sessionCookie.Value
	userID, validSession := utils.VerifySession(sessionUUID, "wsHandler")

	fmt.Println("Session cookie:", sessionCookie.Value)
	fmt.Println("sessiomUUID:", sessionUUID)
	fmt.Println("Verify session return:", userID, validSession)

	if !validSession {
		log.Println("Invalid session")
		return
	}

	log.Println("User connected:", userID)

	err = utils.SetUserOnline(userID)
	if err != nil {
		log.Println("Error setting user online:", err)
		return
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			utils.SetUserOffline(userID)
			log.Println("User disconnected:", userID)
			break
		}
	}
}
