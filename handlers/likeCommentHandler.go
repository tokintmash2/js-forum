package handlers

import (
	"forum-auth/database"
	"forum-auth/utils"
	"net/http"
	"strconv"
)

// LikeCommentHandler handles liking or unliking a comment
func LikeCommentHandler(writer http.ResponseWriter, request *http.Request) {
	// Check user authentication by verifying the session
	cookie, err := request.Cookie("session")
	if err != nil {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	sessionUUID := cookie.Value
	userID, validSession := utils.VerifySession(sessionUUID)
	if !validSession {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	if request.Method == http.MethodPost {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		// Get comment ID from form
		commentIDStr := request.FormValue("comment_id")
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			http.Error(writer, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		alreadyLiked, err := checkIfCommentLiked(userID, commentID)
		if err != nil {
			http.Error(writer, "Error checking if comment is already liked", http.StatusInternalServerError)
			return
		}

		if alreadyLiked {
			// Remove the like
			err = removeCommentLike(userID, commentID)
			if err != nil {
				http.Error(writer, "Error removing like", http.StatusInternalServerError)
				return
			}
		} else {
			// Add the like
			err = likeComment(userID, commentID)
			if err != nil {
				http.Error(writer, "Error liking comment", http.StatusInternalServerError)
				return
			}
		}

		// Redirect back to the referring page after performing the action
		referer := request.Header.Get("Referer")
		if referer != "" {
			http.Redirect(writer, request, referer, http.StatusSeeOther)
		} else {
			// If no referring page is found, redirect to the home page
			http.Redirect(writer, request, "/", http.StatusSeeOther)
		}
		return
	}

	// Handle cases where the request method is not POST
	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
}

// DislikeCommentHandler handles disliking or removing dislike from a comment
func DislikeCommentHandler(writer http.ResponseWriter, request *http.Request) {
	// Check user authentication by verifying the session
	cookie, err := request.Cookie("session")
	if err != nil {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	sessionUUID := cookie.Value
	userID, validSession := utils.VerifySession(sessionUUID)
	if !validSession {
		http.Redirect(writer, request, "/sign-in", http.StatusSeeOther)
		return
	}

	if request.Method == http.MethodPost {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		// Get comment ID from form
		commentIDStr := request.FormValue("comment_id")
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			http.Error(writer, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		alreadyDisliked, err := checkIfCommentDisliked(userID, commentID)
		if err != nil {
			http.Error(writer, "Error checking if comment is already disliked", http.StatusInternalServerError)
			return
		}

		if alreadyDisliked {
			// Remove the dislike
			err = removeCommentDislike(userID, commentID)
			if err != nil {
				http.Error(writer, "Error removing dislike", http.StatusInternalServerError)
				return
			}
		} else {
			// Add the dislike
			err = dislikeComment(userID, commentID)
			if err != nil {
				http.Error(writer, "Error disliking comment", http.StatusInternalServerError)
				return
			}
		}

		// Redirect back to the referring page after performing the action
		referer := request.Header.Get("Referer")
		if referer != "" {
			http.Redirect(writer, request, referer, http.StatusSeeOther)
		} else {
			// If no referring page is found, redirect to the home page
			http.Redirect(writer, request, "/", http.StatusSeeOther)
		}
		return
	}

	// Handle cases where the request method is not POST
	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
}

// removeCommentLike removes a like from a comment
func removeCommentLike(userID, commentID int) error {
	_, err := database.DB.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND comment_id = ? AND is_comment_like = 1
    `, userID, commentID)
	if err != nil {
		return err
	}
	return nil
}

// removeCommentDislike removes a dislike from a comment
func removeCommentDislike(userID, commentID int) error {
	_, err := database.DB.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND comment_id = ? AND is_comment_like = 0
    `, userID, commentID)
	if err != nil {
		return err
	}
	return nil
}

// likeComment adds a like to a comment
func likeComment(userID, commentID int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing dislike if present
	_, err = tx.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND comment_id = ? AND is_comment_like = 0
    `, userID, commentID)
	if err != nil {
		return err
	}

	// Check if the user already liked the comment
	alreadyLiked, err := checkIfCommentLiked(userID, commentID)
	if err != nil {
		return err
	}

	if alreadyLiked {
		// User already liked, no need to re-add like, commit transaction and return
		err = tx.Commit()
		if err != nil {
			return err
		}
		return nil
	}

	// Add the like
	_, err = tx.Exec(`
        INSERT INTO likes (user_id, comment_id, is_comment_like)
        VALUES (?, ?, 1)
    `, userID, commentID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// dislikeComment adds a dislike to a comment
func dislikeComment(userID, commentID int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing like if present
	_, err = tx.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND comment_id = ? AND is_comment_like = 1
    `, userID, commentID)
	if err != nil {
		return err
	}

	// Check if the user already disliked the comment
	alreadyDisliked, err := checkIfCommentDisliked(userID, commentID)
	if err != nil {
		return err
	}

	if alreadyDisliked {
		// User already disliked, no need to re-add dislike, commit transaction and return
		err = tx.Commit()
		if err != nil {
			return err
		}
		return nil
	}

	// Add the dislike
	_, err = tx.Exec(`
        INSERT INTO likes (user_id, comment_id, is_comment_like)
        VALUES (?, ?, 0)
    `, userID, commentID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// checkIfCommentLiked checks if a comment is liked by a user
func checkIfCommentLiked(userID, commentID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM likes
		WHERE user_id = ? AND comment_id = ? AND is_comment_like = 1
	`, userID, commentID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// checkIfCommentDisliked checks if a comment is disliked by a user
func checkIfCommentDisliked(userID, commentID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM likes
		WHERE user_id = ? AND comment_id = ? AND is_comment_like = 0
	`, userID, commentID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
