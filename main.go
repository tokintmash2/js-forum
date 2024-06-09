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
	// mux.HandleFunc("/liked-posts", handlers.LikedPostsHandler)
	// mux.HandleFunc("/my-posts", handlers.MyPostsHandler)
	mux.HandleFunc("/sign-in", handlers.SignInHandler)
	mux.HandleFunc("/sign-up", handlers.SignUpHandler)
	mux.HandleFunc("/sign-in-form", handlers.SignInFormHandler)
	mux.HandleFunc("/sign-up-form", handlers.SignUpFormHandler)
	// mux.HandleFunc("/add-comment", handlers.AddCommentHandler)
	mux.HandleFunc("/github-login", handlers.GitHubLoginHandler)
	mux.HandleFunc("/github-sign-up", handlers.GitHubLoginHandler)
	mux.HandleFunc("/github-callback", handlers.GitHubCallbackHandler)
	mux.HandleFunc("/google-login", handlers.GoogleLoginHandler)
	mux.HandleFunc("/google-sign-up", handlers.GoogleLoginHandler)
	mux.HandleFunc("/google-callback", handlers.GoogleCallbackHandler)

	// -------- JAvaScript API ---------
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

	// Start the server on port 4000 and log its status
	log.Println("Server started at :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
