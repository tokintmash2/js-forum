package handlers

import (
	"real-forum/database"
)

// RemovePostLike removes a like for a post from the database
func RemovePostLike(userID, postID int) error {
	_, err := database.DB.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND post_id = ? AND is_post_like = 1
    `, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

// RemovePostDislike removes a dislike for a post from the database
func RemovePostDislike(userID, postID int) error {
	_, err := database.DB.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND post_id = ? AND is_post_like = 0
    `, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

// LikePost adds a like for a post to the database
func LikePost(userID, postID int) error {
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
	alreadyLiked, err := CheckIfPostLiked(userID, postID)
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
func DislikePost(userID, postID int) error {
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
	alreadyDisliked, err := CheckIfPostDisliked(userID, postID)
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
func CheckIfPostLiked(userID, postID int) (bool, error) {
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

// CheckIfPostDisliked checks if a user has already disliked a post
func CheckIfPostDisliked(userID, postID int) (bool, error) {
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
