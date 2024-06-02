package utils

import (
	"real-forum/database"
)

// Returns the current like and dislike count for a specific post
func GetPostLikeCount(postID int) (int, int, error) {
	var likeCount int
	var dislikeCount int
	err := database.DB.QueryRow(`
        SELECT 
            COALESCE(SUM(CASE WHEN is_post_like = 1 THEN 1 ELSE 0 END), 0) AS likes,
            COALESCE(SUM(CASE WHEN is_post_like = 0 THEN 1 ELSE 0 END), 0) AS dislikes
        FROM likes
        WHERE post_id = ?
    `, postID).Scan(&likeCount, &dislikeCount)
	if err != nil {
		return 0, 0, err
	}

	return likeCount, dislikeCount, nil
}

func GetCommentLikeCount(commentID int) (int, int, error) {
	var likeCount int
	var dislikeCount int
	err := database.DB.QueryRow(`
        SELECT 
            COALESCE(SUM(CASE WHEN is_comment_like = 1 THEN 1 ELSE 0 END), 0) AS likes,
            COALESCE(SUM(CASE WHEN is_comment_like = 0 THEN 1 ELSE 0 END), 0) AS dislikes
        FROM likes
        WHERE comment_id = ?
    `, commentID).Scan(&likeCount, &dislikeCount)
	if err != nil {
		return 0, 0, err
	}
	return likeCount, dislikeCount, nil
}

// RemoveCommentLike removes a like from a comment
func RemoveCommentLike(userID, commentID int) error {
	_, err := database.DB.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND comment_id = ? AND is_comment_like = 1
    `, userID, commentID)
	if err != nil {
		return err
	}
	return nil
}

// RemoveCommentDislike removes a dislike from a comment
func RemoveCommentDislike(userID, commentID int) error {
	_, err := database.DB.Exec(`
        DELETE FROM likes
        WHERE user_id = ? AND comment_id = ? AND is_comment_like = 0
    `, userID, commentID)
	if err != nil {
		return err
	}
	return nil
}

// LikeComment adds a like to a comment
func LikeComment(userID, commentID int) error {
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
	alreadyLiked, err := CheckIfCommentLiked(userID, commentID)
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

// DislikeComment adds a dislike to a comment
func DislikeComment(userID, commentID int) error {
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
	alreadyDisliked, err := CheckIfCommentDisliked(userID, commentID)
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

// CheckIfCommentLiked checks if a comment is liked by a user
func CheckIfCommentLiked(userID, commentID int) (bool, error) {
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

// CheckIfCommentDisliked checks if a comment is disliked by a user
func CheckIfCommentDisliked(userID, commentID int) (bool, error) {
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
