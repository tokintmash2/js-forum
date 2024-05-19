package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum-auth/database"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles various user-related requests
func UserHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/sign-in":
		if request.Method == http.MethodGet {
			SignInFormHandler(writer, request)
		} else if request.Method == http.MethodPost {
			SignInHandler(writer, request)
		}
	case "/sign-up":
		if request.Method == http.MethodGet {
			SignUpFormHandler(writer, request)
		} else if request.Method == http.MethodPost {
			SignUpHandler(writer, request)
		}
	case "/sign-in-form":
		SignInFormHandler(writer, request)
	case "/sign-up-form":
		SignUpFormHandler(writer, request)
	case "/github-login":
		GitHubLoginHandler(writer, request)
	case "/github-callback":
		GitHubCallbackHandler(writer, request)
	case "/google-login":
		GoogleLoginHandler(writer, request)
	case "/google-callback":
		GoogleCallbackHandler(writer, request)
	}
}

// SignInHandler handles the sign-in process
func SignInHandler(writer http.ResponseWriter, request *http.Request) {
	newUser := getUser(request)
	userID, verified := verifyUser(newUser)
	if verified {
		// If user is verified, create a session and set a cookie
		sessionUUID, err := createSessionForUser(writer, userID) // -------- this was createSession
		if err != nil {
			http.Error(writer, "Failed to create a session", http.StatusInternalServerError)
			return
		}
		cookie := createSessionCookie(sessionUUID)
		http.SetCookie(writer, cookie)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	} else {
		// Check if the email exists in the database
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", newUser.Email).Scan(&count)
		if err != nil {
			http.Error(writer, "Error checking user", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			// Email doesn't exist
			handleResponse(writer, "sign_in.html", false, "Wrong email or password")
			return
		}
		// Email exists, but password is incorrect
		handleResponse(writer, "sign_in.html", false, "Wrong email or password")
		return
	}
}

// SignUpHandler handles the sign-up process
func SignUpHandler(writer http.ResponseWriter, request *http.Request) {
	newUser := getUser(request)
	// Check if the user submitted the form without any input
	if newUser.Email == "" || newUser.Password == "" || newUser.Username == "" {
		message := "Please fill in all fields."
		handleResponse(writer, "sign_up.html", false, message)
		return
	}

	err := createUser(newUser)
	if err == nil {
		// If user creation is successful, create a session and redirect to sign-in page
		sessionUUID, err := createSession(newUser.ID)
		if err != nil {
			http.Error(writer, "Failed to create a session", http.StatusInternalServerError)
			return
		}
		cookie := createSessionCookie(sessionUUID)
		http.SetCookie(writer, cookie)
		http.Redirect(writer, request, "/sign-in-form", http.StatusSeeOther)
	} else {
		// If user creation fails, handle the error message based on the error type
		errorMessage := "Error creating user."
		if strings.Contains(err.Error(), "email already exists") {
			errorMessage = "Email is already taken."
		} else if strings.Contains(err.Error(), "username already exists") {
			errorMessage = "Username is already taken."
		}
		handleResponse(writer, "sign_up.html", false, errorMessage)
	}
}

// SignInFormHandler serves the sign-in form
func SignInFormHandler(writer http.ResponseWriter, request *http.Request) {
	tmplFile := path.Join("templates", "sign_in.html")
	t, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(writer, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

// SignUpFormHandler serves the sign-up form
func SignUpFormHandler(writer http.ResponseWriter, request *http.Request) {
	// Check if the form has been submitted
	if request.Method == http.MethodPost {
		SignUpHandler(writer, request) // Process the sign-up form
		return
	}

	// Display the sign-up form with a generic message
	message := "Please fill in this form to create an account."
	handleResponse(writer, "sign_up.html", false, message)
}

// getUser extracts user data from the request
func getUser(request *http.Request) User {
	email := request.FormValue("email")
	password := request.FormValue("password")
	username := request.FormValue("username")
	return User{
		Email:    email,
		Password: password,
		Username: username,
	}
}

// handleResponse handles template rendering based on success or failure
func handleResponse(writer http.ResponseWriter, filename string, success bool, message string) {
	statusMessage := message
	if success {
		statusMessage += " Success."
	}
	tmplFile := path.Join("templates", filename)
	t, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(writer, statusMessage)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

// User represents user data
type User struct {
	ID       int
	Email    string
	Password string
	Username string
	GitHubID string
	GoogleID string
}

// createUser creates a new user in the database
func createUser(newUser User) error {
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

// verifyUser verifies user credentials
func verifyUser(user User) (int, bool) {
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

// getPasswordHash generates a bcrypt hash for a password
func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}

// createSession creates a new session for a user
func createSession(userID int) (string, error) {
	// Check if there's an existing active session for the user
	var existingSessionID string
	err := database.DB.QueryRow(`
        SELECT session_uuid FROM sessions
        WHERE user_id = ? AND expiry > CURRENT_TIMESTAMP`,
		userID,
	).Scan(&existingSessionID)

	if err == nil {
		// If an active session exists, clear it before creating a new one
		err = clearSession(existingSessionID)
		if err != nil {
			return "", err
		}
	}

	// Create a new session
	sessionUUID := uuid.New().String()
	expiry := time.Now().Add(24 * time.Hour) // Set the session expiry time (24 hours)
	_, err = database.DB.Exec(`
        INSERT INTO sessions (user_id, session_uuid, expiry)
        VALUES (?, ?, ?)`,
		userID, sessionUUID, expiry,
	)
	if err != nil {
		return "", err
	}
	return sessionUUID, nil
}

// ------------------

func createSessionForUser(writer http.ResponseWriter, userID int) (string, error) {
	sessionUUID, err := createSession(userID)
	if err != nil {
		return "", err
	}

	cookie := createSessionCookie(sessionUUID)
	http.SetCookie(writer, cookie)

	return sessionUUID, nil
}

const (
	githubClientID     = "454f332c762d46733ed6"
	githubClientSecret = "34e0d9f3a4d8ac1f240553af074f7da7e4ad1340"
	githubRedirectURI  = "http://localhost:4000/github-callback"
)

type GitHubUserData struct {
	ID    int `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

func GitHubLoginHandler(writer http.ResponseWriter, request *http.Request) {
	githubAuthURL := "https://github.com/login/oauth/authorize"
	githubAuthURL += "?client_id=454f332c762d46733ed6"
	githubAuthURL += "&redirect_uri=" + githubRedirectURI
	githubAuthURL += "&scope=user:email" // Specify the required scope
	http.Redirect(writer, request, githubAuthURL, http.StatusSeeOther)
}

func GitHubCallbackHandler(writer http.ResponseWriter, request *http.Request) {
	// Retrieve the authorization code from the query parameters
	code := request.URL.Query().Get("code")
	if code == "" {
		http.Error(writer, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	accessToken, err := exchangeGitHubCodeForToken(code)
	if err != nil {
		http.Error(writer, "Failed to exchange code for access token", http.StatusInternalServerError)
		return
	}

	// Fetch user data from GitHub using the access token
	userData, err := fetchGitHubUserData(accessToken)
	if err != nil {
		http.Error(writer, "Failed to fetch user data from GitHub", http.StatusInternalServerError)
		return
	}

	// Print the received user data
	// fmt.Printf("GitHub User Data: %+v\n", userData)

	// Check if the user already exists in the database
	userID, verified := verifyUser(User{Email: userData.Email})
	if !verified {
		// User doesn't exist, create a new user
		newUser := User{
			Email:    userData.Email,
			Username: userData.Login,
		}

		err := createUser(newUser)
		if err != nil {
			http.Error(writer, "Failed to create a user", http.StatusInternalServerError)
			return
		}

		// Now, the user should exist, retrieve the userID again
		userID, verified = verifyUser(newUser)
		if !verified {
			http.Error(writer, "Failed to verify the user", http.StatusInternalServerError)
			return
		}
	}

	// Update the user's GitHubID in the database
	err = updateGitHubIDInDatabase(userID, strconv.Itoa(userData.ID))
	if err != nil {
		http.Error(writer, "Failed to update GitHub ID in the database", http.StatusInternalServerError)
		return
	}

	// Create a session for the user
	_, err = createSessionForUser(writer, userID)
	if err != nil {
		http.Error(writer, "Failed to create a session", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the main page or any desired page
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func exchangeGitHubCodeForToken(code string) (string, error) {
	// GitHub token endpoint
	tokenURL := "https://github.com/login/oauth/access_token"

	// Create a POST request to exchange the authorization code for an access token
	data := url.Values{}
	data.Set("client_id", githubClientID)
	data.Set("client_secret", githubClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", githubRedirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the response body to extract the access token
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}

	accessToken := values.Get("access_token")
	if accessToken == "" {
		return "", fmt.Errorf("GitHub did not provide an access token")
	}

	return accessToken, nil
}

func fetchGitHubUserData(accessToken string) (GitHubUserData, error) {
	// GitHub user data endpoint
	userDataURL := "https://api.github.com/user"

	// Create a GET request to fetch user data
	req, err := http.NewRequest("GET", userDataURL, nil)
	if err != nil {
		return GitHubUserData{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GitHubUserData{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GitHubUserData{}, err
	}

	// fmt.Printf("GitHub User Data Response: %s\n", body)

	// Parse the response body into a GitHubUserData struct
	var userData GitHubUserData
	err = json.Unmarshal(body, &userData)
	if err != nil {
		return GitHubUserData{}, err
	}

	// Retrieve the user's email separately using another GitHub API endpoint
	emailURL := "https://api.github.com/user/emails"
	req, err = http.NewRequest("GET", emailURL, nil)
	if err != nil {
		return GitHubUserData{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Perform the request
	resp, err = client.Do(req)
	if err != nil {
		return GitHubUserData{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	emailBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GitHubUserData{}, err
	}

	// fmt.Printf("GitHub User Email Response: %s\n", emailBody)

	// Parse the response body to extract the primary email
	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	err = json.Unmarshal(emailBody, &emails)
	if err != nil {
		return GitHubUserData{}, err
	}

	// Find the primary email
	for _, e := range emails {
		if e.Primary {
			userData.Email = e.Email
			break
		}
	}

	return userData, nil
}


func updateGitHubIDInDatabase(userID int, github_id string) error {

	// fmt.Printf("Updating GitHub ID for user ID %d to %s\n", userID, github_id)

	result, err := database.DB.Exec("UPDATE users SET github_id = ? WHERE id = ?", github_id, userID)
	if err != nil {
		// Check if the error is due to no rows being affected
		if strings.Contains(err.Error(), "no rows in result set") {
			// User doesn't exist, insert a new row
			return insertGitHubUser(userID, github_id)
		}
		return err
	}

	// Check the number of rows affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		// Handle the case where the update did not affect exactly one row
		return fmt.Errorf("expected 1 row to be affected, but %d rows were affected", rowsAffected)
	}

	return nil

}

func insertGitHubUser(userID int, githubID string) error {
	_, err := database.DB.Exec("INSERT INTO users (id, github_id) VALUES (?, ?)", userID, githubID)
	return err
}

const (
	googleClientID     = "677000148449-qt8kejfk3j1id34kqnmfn2p56lc6jonv.apps.googleusercontent.com"
	googleClientSecret = "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"
	googleRedirectURI  = "http://localhost:4000/google-callback"
)

// GoogleUserData represents user data from Google
type GoogleUserData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func GoogleLoginHandler(writer http.ResponseWriter, request *http.Request) {
	googleAuthURL := "https://accounts.google.com/o/oauth2/auth"
	googleAuthURL += "?client_id=" + googleClientID
	googleAuthURL += "&redirect_uri=" + googleRedirectURI
	googleAuthURL += "&scope=openid profile email" // Specify the required scope
	googleAuthURL += "&response_type=code"
	http.Redirect(writer, request, googleAuthURL, http.StatusSeeOther)
}

func GoogleCallbackHandler(writer http.ResponseWriter, request *http.Request) {
	// Retrieve the authorization code from the query parameters
	code := request.URL.Query().Get("code")
	if code == "" {
		http.Error(writer, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	accessToken, err := exchangeGoogleCodeForToken(code)
	if err != nil {
		http.Error(writer, "Failed to exchange code for access token", http.StatusInternalServerError)
		return
	}

	// Fetch user data from Google using the access token
	userData, err := fetchGoogleUserData(accessToken)
	if err != nil {
		http.Error(writer, "Failed to fetch user data from Google", http.StatusInternalServerError)
		return
	}

	// Check if the user already exists in the database
	userID, verified := verifyUser(User{Email: userData.Email})
	if !verified {
		// User doesn't exist, create a new user
		newUser := User{
			Email:    userData.Email,
			Username: userData.Name,
			// Set other fields as needed
		}

		err := createUser(newUser)
		if err != nil {
			http.Error(writer, "Failed to create a user", http.StatusInternalServerError)
			return
		}

		// Now, the user should exist, retrieve the userID again
		userID, verified = verifyUser(newUser)
		if !verified {
			http.Error(writer, "Failed to verify the user", http.StatusInternalServerError)
			return
		}
	}

	// Update the user's GoogleID in the database
	err = updateGoogleIDInDatabase(userID, userData.ID)
	if err != nil {
		http.Error(writer, "Failed to update Google ID in the database", http.StatusInternalServerError)
		return
	}

	// Create a session for the user
	_, err = createSessionForUser(writer, userID)
	if err != nil {
		http.Error(writer, "Failed to create a session", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the main page or any desired page
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func exchangeGoogleCodeForToken(code string) (string, error) {
	// Google token endpoint
	tokenURL := "https://oauth2.googleapis.com/token"

	// Create a POST request to exchange the authorization code for an access token
	data := url.Values{}
	data.Set("client_id", googleClientID)
	data.Set("client_secret", googleClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", googleRedirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the response body to extract the access token
	var tokenResponse map[string]interface{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("google did not provide an access token")
	}

	return accessToken, nil
}

func fetchGoogleUserData(accessToken string) (GoogleUserData, error) {
	// Google user data endpoint
	userDataURL := "https://www.googleapis.com/oauth2/v2/userinfo"

	// Create a GET request to fetch user data
	req, err := http.NewRequest("GET", userDataURL, nil)
	if err != nil {
		return GoogleUserData{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GoogleUserData{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GoogleUserData{}, err
	}

	// Parse the response body into a GoogleUserData struct
	var userData GoogleUserData
	err = json.Unmarshal(body, &userData)
	if err != nil {
		return GoogleUserData{}, err
	}

	return userData, nil
}

func updateGoogleIDInDatabase(userID int, googleID string) error {
	// fmt.Printf("Updating Google ID for user ID %d to %s\n", userID, googleID)

	result, err := database.DB.Exec("UPDATE users SET google_id = ? WHERE id = ?", googleID, userID)
	if err != nil {
		// Check if the error is due to no rows being affected
		if strings.Contains(err.Error(), "no rows in result set") {
			// User doesn't exist, insert a new row
			return insertGoogleUser(userID, googleID)
		}
		return err
	}

	// Check the number of rows affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		// Handle the case where the update did not affect exactly one row
		return fmt.Errorf("expected 1 row to be affected, but %d rows were affected", rowsAffected)
	}

	return nil
}

func insertGoogleUser(userID int, googleID string) error {
	_, err := database.DB.Exec("INSERT INTO users (id, google_id) VALUES (?, ?)", userID, googleID)
	return err
}

// createSessionCookie creates a cookie for the session
func createSessionCookie(sessionUUID string) *http.Cookie {
	expiration := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionUUID,
		Expires:  expiration,
		HttpOnly: true,
	}
	return cookie
}

// clearSession clears a session by its UUID
func clearSession(sessionUUID string) error {
	_, err := database.DB.Exec(`
		DELETE FROM sessions
		WHERE session_uuid = ?`,
		sessionUUID,
	)
	if err != nil {
		return err
	}
	return nil
}
