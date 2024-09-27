package utils

import (
	"fmt"
	"log"
	"real-forum/database"
	"real-forum/structs"
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

func GetRecentMessages(user, partner int) []structs.Message { // Missing userIDs

	rows, err := database.DB.Query(`SELECT sender_id, receiver_id, content, timestamp
		FROM messages 
		WHERE (
			sender_id = ? AND receiver_id = ?)
		OR (
			sender_id = ? AND receiver_id = ?)
		ORDER BY timestamp DESC 
		LIMIT 10;`,
	user, partner, partner, user)
	if err != nil {
		log.Println("error getting recent messages", err)
	}
	defer rows.Close()

	var recentMessages []structs.Message
	for rows.Next() {
		var msg structs.Message

		if err := rows.Scan(
			&msg.Sender, &msg.Recipient, &msg.Content, &msg.CreatedAt,
		); err != nil {
			log.Println("error scanning rec messages", err)
		}
		// fmt.Println(msg.Content)
		recentMessages = append(recentMessages, msg)
	}
	fmt.Println(recentMessages)
	return recentMessages
}
