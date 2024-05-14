package utils

import "text/template"

// Templates stores parsed templates for use across the application
var Templates *template.Template

// init function initializes the Templates variable by parsing template files
func init() {
	// Define and parse templates from the specified files
	Templates = template.Must(template.ParseFiles(
		"templates/404.html",
		"templates/add_comment.html",
		"templates/base.html",
		"templates/categories.html",
		"templates/category_posts.html",
		"templates/comment_buttons.html",
		"templates/comments.html",
		"templates/create_post.html",
		"templates/display_comments.html",
		"templates/like_buttons.html",
		"templates/liked_posts.html",
		"templates/my_posts.html",
		"templates/navbar.html",
		"templates/recent_posts.html",
		"templates/sign_in.html",
		"templates/sign_up.html",
	))
}
