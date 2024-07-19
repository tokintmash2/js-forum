package utils

import (
	"errors"
	"fmt"
	"log"
	"real-forum/database"
	"real-forum/structs"

	"golang.org/x/crypto/bcrypt"
)

func VerifyUser(user structs.User) (int, bool) {
	var storedPassword string
	var userID int
	err := database.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", user.Email).Scan(&userID, &storedPassword)
	// fmt.Printf("Useri email %v", user.Email)
	if err != nil {
		return 0, false
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	return userID, err == nil
}

func CreateUser(newUser structs.User) error {
	var emailCount, usernameCount int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", newUser.Email).Scan(&emailCount)
	if err != nil {
		return err
	}
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", newUser.Username).Scan(&usernameCount)
	if err != nil {
		return err
	}

	if emailCount > 0 {
		return errors.New("email already exists")
	}
	if usernameCount > 0 {
		return errors.New("username already exists")
	}

	passwordHash, err := getPasswordHash(newUser.Password)
	if err != nil {
		return err
	}
	result, err := database.DB.Exec("INSERT INTO users (email, password, username, github_id, google_id) VALUES (?, ?, ?, ?, ?)", newUser.Email, passwordHash, newUser.Username, "", "")

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("user creation failed for email: %s", newUser.Email)
	}
	return nil
}

func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}

func SetUserOnline(userID int) error {
	_, err := database.DB.Exec("INSERT OR REPLACE INTO online_status (user_id, online) VALUES (?, ?)", userID, true)
	if err != nil {
		log.Println("Error setting user online:", err)
		return err
	}
	return nil
}

func SetUserOffline(userID int) error {
	_, err := database.DB.Exec("UPDATE online_status SET online = ? WHERE user_id = ?", false, userID)
	if err != nil {
		log.Println("Error setting user offline:", err)
		return err
	}
	return nil
}

func GetOnlineUsers() ([]structs.User, error) {
	rows, err := database.DB.Query("SELECT users.id, users.username FROM users JOIN online_status ON users.id = online_status.user_id WHERE online_status.online = ?", true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []structs.User
	for rows.Next() {
		var user structs.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// fmt.Println("Online users from GetOnlineUsers", users)
	return users, nil
}
