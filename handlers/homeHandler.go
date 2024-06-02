package handlers

import (
	"net/http"
)

// HomePageHandler manages the homepage and handles displaying recent posts
func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "index.html")

}
