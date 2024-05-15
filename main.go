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
	mux.HandleFunc("/create-post", handlers.CreatePostHandler)
	mux.HandleFunc("/liked-posts", handlers.LikedPostsHandler)
	mux.HandleFunc("/like", handlers.LikePostHandler)
	mux.HandleFunc("/dislike", handlers.DislikePostHandler)
	mux.HandleFunc("/my-posts", handlers.MyPostsHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	mux.HandleFunc("/sign-in", handlers.SignInHandler)
	mux.HandleFunc("/sign-up", handlers.SignUpHandler)
	mux.HandleFunc("/sign-in-form", handlers.SignInFormHandler)
	mux.HandleFunc("/sign-up-form", handlers.SignUpFormHandler)
	// mux.HandleFunc("/category/", handlers.CategoryPostsHandler)
	mux.HandleFunc("/add-comment", handlers.AddCommentHandler)
	mux.HandleFunc("/like-comment", handlers.LikeCommentHandler)
	mux.HandleFunc("/dislike-comment", handlers.DislikeCommentHandler)
	mux.HandleFunc("/github-login", handlers.GitHubLoginHandler)
	mux.HandleFunc("/github-sign-up", handlers.GitHubLoginHandler)
	mux.HandleFunc("/github-callback", handlers.GitHubCallbackHandler)
	mux.HandleFunc("/google-login", handlers.GoogleLoginHandler)
	mux.HandleFunc("/google-sign-up", handlers.GoogleLoginHandler)
	mux.HandleFunc("/google-callback", handlers.GoogleCallbackHandler)
	//--------- TESTING ----------------
	mux.HandleFunc("/api/test", api.TestHandler)
	// -------- JAvaScript API ---------
	mux.HandleFunc("/api/categories", api.CategoriesHandler)
	mux.HandleFunc("/api/recents", api.RecentPostsHandler)
	mux.HandleFunc("/api/home", api.HomeJSONHandler)
	mux.HandleFunc("/api/category/", api.CategoryPostsApiHandler)

	// Start the server on port 4000 and log its status
	log.Println("Server started at :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
