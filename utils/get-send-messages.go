package utils

import (
	"fmt"
	"log"
	"real-forum/database"
	"real-forum/structs"
	"time"
)

// Save messages to database
func SaveMessage(m structs.Message) {

	log.Println("MSG:", m)

	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("database error:", err)
	}
	defer tx.Rollback()

	// Insert message
	_, err = tx.Exec(`
        INSERT INTO messages (sender_id, receiver_id, content, timestamp)
        VALUES (?, ?, ?, ?)`,
		m.Sender, m.Recipient, m.Content, m.CreatedAt,
	)
	if err != nil {
		log.Printf("Error inserting message: %v\n", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return
	}
	log.Println("Message successfully saved to the database")
}

// func GetRecentMessages(user, partner int, time time.Time) []structs.Message {

// 	log.Println("Getting recent messages")
// 	log.Println("user:", user)
// 	log.Println("partner:", partner)

// 	rows, err := database.DB.Query(`
//     SELECT m.sender_id, m.receiver_id, m.content, m.timestamp,
//            s.username AS sender_username, r.username AS receiver_username
//     FROM messages m
//     JOIN users s ON m.sender_id = s.id
//     JOIN users r ON m.receiver_id = r.id
//     WHERE (m.sender_id = ? AND m.receiver_id = ?)
//        OR (m.sender_id = ? AND m.receiver_id = ?)
//     ORDER BY m.timestamp DESC 
//     LIMIT 10;`,
// 		user, partner, partner, user)

// 	if err != nil {
// 		log.Println("error getting recent messages", err)
// 	}
// 	defer rows.Close()

// 	var recentMessages []structs.Message
// 	for rows.Next() {
// 		var msg structs.Message
// 		msg.Type = "conversation"

// 		if err := rows.Scan(
// 			&msg.Sender, &msg.Recipient, &msg.Content, &msg.CreatedAt,
// 			&msg.SenderUsername, &msg.ReceiverUsername,
// 		); err != nil {
// 			log.Println("error scanning rec messages", err)
// 		}
// 		recentMessages = append(recentMessages, msg)
// 	}

// 	fmt.Println("Recent messages:", recentMessages)
// 	return recentMessages
// }

func GetRecentMessages(user, partner int, time time.Time) []structs.Message {

	log.Println("Getting recent messages older than:", time)
	log.Println("user:", user)
	log.Println("partner:", partner)

	// Modify the SQL query to include the timestamp condition
	rows, err := database.DB.Query(`
    SELECT m.sender_id, m.receiver_id, m.content, m.timestamp,
           s.username AS sender_username, r.username AS receiver_username
    FROM messages m
    JOIN users s ON m.sender_id = s.id
    JOIN users r ON m.receiver_id = r.id
    WHERE ((m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?))
      AND m.timestamp < ?
    ORDER BY m.timestamp DESC 
    LIMIT 10;`,
		user, partner, partner, user, time)

	if err != nil {
		log.Println("error getting recent messages", err)
	}
	defer rows.Close()

	var recentMessages []structs.Message
	for rows.Next() {
		var msg structs.Message
		msg.Type = "conversation"

		if err := rows.Scan(
			&msg.Sender, &msg.Recipient, &msg.Content, &msg.CreatedAt,
			&msg.SenderUsername, &msg.ReceiverUsername,
		); err != nil {
			log.Println("error scanning recent messages", err)
		}
		recentMessages = append(recentMessages, msg)
	}

	fmt.Println("Recent messages:", recentMessages)
	return recentMessages
}
