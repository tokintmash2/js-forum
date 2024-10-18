package handlers

import (
	"html/template"
	"net/http"
	"path"
)

func NotFoundWrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request matches any known route
		if r.URL.Path != "/" &&
			r.URL.Path != "/create_post" &&
			r.URL.Path != "/liked-posts" &&
			r.URL.Path != "/like" &&
			r.URL.Path != "/dislike" &&
			r.URL.Path != "/my-posts" &&
			r.URL.Path != "/logout" &&
			// r.URL.Path != "/log-in" &&
			r.URL.Path != "/sign-in" &&
			r.URL.Path != "/register" &&
			r.URL.Path != "/sign-in-form" &&
			r.URL.Path != "/sign-up-form" &&
			r.URL.Path != "/category" &&
			r.URL.Path != "/add-comment" &&
			r.URL.Path != "/like-comment" &&
			r.URL.Path != "/dislike-comment" {

			// If not, invoke the custom 404 handler
			notFoundHandler(w, r)
			return
		}
		// Otherwise, call the original handler
		handler(w, r)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	tmplFile := path.Join("templates", "404.html")
	t, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
