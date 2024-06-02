package handlers

// import (
// 	"real-forum/database"
// 	"real-forum/utils"
// 	"net/http"
// 	"strconv"
// )

// LikeCommentHandler handles liking or unliking a comment
// func LikeCommentHandler(writer http.ResponseWriter, request *http.Request) {
// 	// Check user authentication by verifying the session
// 	cookie, err := request.Cookie("session")
// 	if err != nil {
// 		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
// 		return
// 	}

// 	sessionUUID := cookie.Value
// 	userID, validSession := utils.VerifySession(sessionUUID)
// 	if !validSession {
// 		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
// 		return
// 	}

// 	if request.Method == http.MethodPost {
// 		err := request.ParseForm()
// 		if err != nil {
// 			http.Error(writer, "Error parsing form data", http.StatusInternalServerError)
// 			return
// 		}

// 		// Get comment ID from form
// 		commentIDStr := request.FormValue("comment_id")
// 		commentID, err := strconv.Atoi(commentIDStr)
// 		if err != nil {
// 			http.Error(writer, "Invalid comment ID", http.StatusBadRequest)
// 			return
// 		}

// 		alreadyLiked, err := CheckIfCommentLiked(userID, commentID)
// 		if err != nil {
// 			http.Error(writer, "Error checking if comment is already liked", http.StatusInternalServerError)
// 			return
// 		}

// 		if alreadyLiked {
// 			// Remove the like
// 			err = RemoveCommentLike(userID, commentID)
// 			if err != nil {
// 				http.Error(writer, "Error removing like", http.StatusInternalServerError)
// 				return
// 			}
// 		} else {
// 			// Add the like
// 			err = LikeComment(userID, commentID)
// 			if err != nil {
// 				http.Error(writer, "Error liking comment", http.StatusInternalServerError)
// 				return
// 			}
// 		}

// 		// Redirect back to the referring page after performing the action
// 		referer := request.Header.Get("Referer")
// 		if referer != "" {
// 			http.Redirect(writer, request, referer, http.StatusSeeOther)
// 		} else {
// 			// If no referring page is found, redirect to the home page
// 			http.Redirect(writer, request, "/", http.StatusSeeOther)
// 		}
// 		return
// 	}

// 	// Handle cases where the request method is not POST
// 	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
// }

// // DislikeCommentHandler handles disliking or removing dislike from a comment
// func DislikeCommentHandler(writer http.ResponseWriter, request *http.Request) {
// 	// Check user authentication by verifying the session
// 	cookie, err := request.Cookie("session")
// 	if err != nil {
// 		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
// 		return
// 	}

// 	sessionUUID := cookie.Value
// 	userID, validSession := utils.VerifySession(sessionUUID)
// 	if !validSession {
// 		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
// 		return
// 	}

// 	if request.Method == http.MethodPost {
// 		err := request.ParseForm()
// 		if err != nil {
// 			http.Error(writer, "Error parsing form data", http.StatusInternalServerError)
// 			return
// 		}

// 		// Get comment ID from form
// 		commentIDStr := request.FormValue("comment_id")
// 		commentID, err := strconv.Atoi(commentIDStr)
// 		if err != nil {
// 			http.Error(writer, "Invalid comment ID", http.StatusBadRequest)
// 			return
// 		}

// 		alreadyDisliked, err := checkIfCommentDisliked(userID, commentID)
// 		if err != nil {
// 			http.Error(writer, "Error checking if comment is already disliked", http.StatusInternalServerError)
// 			return
// 		}

// 		if alreadyDisliked {
// 			// Remove the dislike
// 			err = removeCommentDislike(userID, commentID)
// 			if err != nil {
// 				http.Error(writer, "Error removing dislike", http.StatusInternalServerError)
// 				return
// 			}
// 		} else {
// 			// Add the dislike
// 			err = dislikeComment(userID, commentID)
// 			if err != nil {
// 				http.Error(writer, "Error disliking comment", http.StatusInternalServerError)
// 				return
// 			}
// 		}

// 		// Redirect back to the referring page after performing the action
// 		referer := request.Header.Get("Referer")
// 		if referer != "" {
// 			http.Redirect(writer, request, referer, http.StatusSeeOther)
// 		} else {
// 			// If no referring page is found, redirect to the home page
// 			http.Redirect(writer, request, "/", http.StatusSeeOther)
// 		}
// 		return
// 	}

// 	// Handle cases where the request method is not POST
// 	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
// }


