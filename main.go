package main

import (
	"database/sql"
	"log"
	"net/http"
	"real-forum/api"
	"real-forum/database"
	"real-forum/handlers"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error

	// server := &WebSocketServer{
	// 	clients: make(map[string]*Client),
	// }

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

	//ws manager test

	// -------- JAvaScript API ---------
	mux.HandleFunc("/ws", api.HandleConnections)
	mux.HandleFunc("/api/conversation/", api.ChatHandler)
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
