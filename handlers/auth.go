package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"fuzzy/models"
)

// Session management constants
const (
	SessionCookieName = "fuzzy_session"
	SessionDuration   = 24 * time.Hour
)

// Simple in-memory session store
var sessions = make(map[string]int) // sessionID -> userID

// LoginHandler handles the login page and authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetLogin(w, r)
	case http.MethodPost:
		handlePostLogin(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session cookie
	cookie, err := r.Cookie(SessionCookieName)
	if err == nil {
		// Remove session from store
		delete(sessions, cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	// Redirect to login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// SetupHandler handles the first-time setup page
func SetupHandler(w http.ResponseWriter, r *http.Request) {
	// If users already exist, redirect to login
	if models.GlobalStore.HasUsers() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetSetup(w, r)
	case http.MethodPost:
		handlePostSetup(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetLogin(w http.ResponseWriter, r *http.Request) {
	data := models.LoginPageData{
		Title: "Fuzzy - Login",
	}

	renderLoginTemplate(w, &data)
}

func handlePostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		data := models.LoginPageData{
			Title: "Fuzzy - Login",
			Error: "Error parsing form data",
		}
		renderLoginTemplate(w, &data)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	// Validate input
	if username == "" || password == "" {
		data := models.LoginPageData{
			Title: "Fuzzy - Login",
			Error: "Username and password are required",
		}
		renderLoginTemplate(w, &data)
		return
	}

	// Check user credentials
	user, exists := models.GlobalStore.GetUserByUsername(username)
	if !exists || !user.CheckPassword(password) {
		data := models.LoginPageData{
			Title: "Fuzzy - Login",
			Error: "Invalid username or password",
		}
		renderLoginTemplate(w, &data)
		return
	}

	// Check if user is active
	if !user.Active {
		data := models.LoginPageData{
			Title: "Fuzzy - Login",
			Error: "Account is disabled",
		}
		renderLoginTemplate(w, &data)
		return
	}

	// Create session
	sessionID := generateSessionID()
	sessions[sessionID] = user.ID

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(SessionDuration.Seconds()),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleGetSetup(w http.ResponseWriter, r *http.Request) {
	data := models.SetupPageData{
		Title: "Fuzzy - First Time Setup",
	}

	renderSetupTemplate(w, &data)
}

func handlePostSetup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		data := models.SetupPageData{
			Title: "Fuzzy - First Time Setup",
			Error: "Error parsing form data",
		}
		renderSetupTemplate(w, &data)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))

	// Validate input
	if username == "" {
		data := models.SetupPageData{
			Title: "Fuzzy - First Time Setup",
			Error: "Username is required",
		}
		renderSetupTemplate(w, &data)
		return
	}

	if email == "" {
		data := models.SetupPageData{
			Title: "Fuzzy - First Time Setup",
			Error: "Email is required",
		}
		renderSetupTemplate(w, &data)
		return
	}

	if password == "" {
		data := models.SetupPageData{
			Title: "Fuzzy - First Time Setup",
			Error: "Password is required",
		}
		renderSetupTemplate(w, &data)
		return
	}

	if password != confirmPassword {
		data := models.SetupPageData{
			Title: "Fuzzy - First Time Setup",
			Error: "Passwords do not match",
		}
		renderSetupTemplate(w, &data)
		return
	}

	// Create admin user
	user := models.User{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      "Administrator",
		Active:    true,
	}

	if err := user.SetPassword(password); err != nil {
		data := models.SetupPageData{
			Title: "Fuzzy - First Time Setup",
			Error: "Error encrypting password",
		}
		renderSetupTemplate(w, &data)
		return
	}

	// Save user
	models.GlobalStore.CreateUser(user)

	// Redirect to login with success message
	http.Redirect(w, r, "/login?setup=complete", http.StatusSeeOther)
}

func renderLoginTemplate(w http.ResponseWriter, data *models.LoginPageData) {
	tmplPath := filepath.Join("templates", "login.html")
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func renderSetupTemplate(w http.ResponseWriter, data *models.SetupPageData) {
	tmplPath := filepath.Join("templates", "setup.html")
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing setup template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing setup template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// generateSessionID creates a simple session ID (in production, use crypto/rand)
func generateSessionID() string {
	return time.Now().Format("20060102150405") + "-" + time.Now().Format("999999999")
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser(r *http.Request) (models.User, bool) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return models.User{}, false
	}

	userID, exists := sessions[cookie.Value]
	if !exists {
		return models.User{}, false
	}

	user, exists := models.GlobalStore.GetUser(userID)
	return user, exists
}