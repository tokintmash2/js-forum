package handlers

import (
	"forum-auth/database"
	"forum-auth/utils"
	"net/http"
	"strconv"
)

// LikePostHandler handles the like functionality for a post
func LikePostHandler(writer http.ResponseWriter, request *http.Request) {
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

		// Get post ID from form
		postIDStr := request.FormValue("post_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(writer, "Invalid post ID", http.StatusBadRequest)
			return
		}

		alreadyLiked, err := checkIfPostLiked(userID, postID)
		if err != nil {
			http.Error(writer, "Error checking if post is already liked", http.StatusInternalServerError)
			return
		}

		if alreadyLiked {
			// Remove the like
			err = removePostLike(userID, postID)
			if err != nil {
				http.Error(writer, "Error removing like", http.StatusInternalServerError)
				return
			}
		} else {
			// Add the like
			err = likePost(userID, postID)
			if err != nil {
				http.Error(writer, "Error liking post", http.StatusInternalServerError)
				return
			}
		}

		// Redirect to the referring page after performing the action
		referer := request.Header.Get("Referer")
		if referer != "" {
			http.Redirect(writer, request, referer, http.StatusSeeOther)
		} else {
			// If no referring page is found, redirect to the homepage ("/")
			http.Redirect(writer, request, "/", http.StatusSeeOther)
		}
		return
	}
}

// DislikePostHandler handles the dislike functionality for a post
func DislikePostHandler(writer http.ResponseWriter, request *http.Request) {
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

		// Get post ID from form
		postIDStr := request.FormValue("post_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(writer, "Invalid post ID", http.StatusBadRequest)
			return
		}

		alreadyDisliked, err := checkIfPostDisliked(userID, postID)
		if err != nil {
			http.Error(writer, "Error checking if post is already disliked", http.StatusInternalServerError)
			return
		}

		if alreadyDisliked {
			// Remove the dislike
			err = removePostDislike(userID, postID)
			if err != nil {
				http.Error(writer, "Error removing dislike", http.StatusInternalServerError)
				return
			}
		} else {
			// Add the dislike
			err = dislikePost(userID, postID)
			if err != nil {
				http.Error(writer, "Error disliking post", http.StatusInternalServerError)
				return
			}
		}

		referer := request.Header.Get("Referer")
		if referer != "" {
			http.Redirect(writer, request, referer, http.StatusSeeOther)
		} else {
			// If no referring page is found, redirect to the home page
			http.Redirect(writer, request, "/", http.StatusSeeOther)
		}
		return
	}
}

// removePostLike removes a like for a post from the database
func removePostLike(userID, postID int) error {
	_, err := database.DB.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND post_id = ? AND is_post_like = 1
    `, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

// removePostDislike removes a dislike for a post from the database
func removePostDislike(userID, postID int) error {
	_, err := database.DB.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND post_id = ? AND is_post_like = 0
    `, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

// likePost adds a like for a post to the database
func likePost(userID, postID int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing dislike if present
	_, err = tx.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND post_id = ? AND is_post_like = 0
    `, userID, postID)
	if err != nil {
		return err
	}

	// Check if the user already liked the post
	alreadyLiked, err := checkIfPostLiked(userID, postID)
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
        INSERT INTO likes (user_id, post_id, is_post_like)
        VALUES (?, ?, 1)
    `, userID, postID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// dislikePost adds a dislike for a post to the database
func dislikePost(userID, postID int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing like if present
	_, err = tx.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND post_id = ? AND is_post_like = 1
    `, userID, postID)
	if err != nil {
		return err
	}

	// Check if the user already disliked the post
	alreadyDisliked, err := checkIfPostDisliked(userID, postID)
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
        INSERT INTO likes (user_id, post_id, is_post_like)
        VALUES (?, ?, 0)
    `, userID, postID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// checkIfPostLiked checks if a user has already liked a post
func checkIfPostLiked(userID, postID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM likes
		WHERE user_id = ? AND post_id = ? AND is_post_like = 1
	`, userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// checkIfPostDisliked checks if a user has already disliked a post
func checkIfPostDisliked(userID, postID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM likes
		WHERE user_id = ? AND post_id = ? AND is_post_like = 0
	`, userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
