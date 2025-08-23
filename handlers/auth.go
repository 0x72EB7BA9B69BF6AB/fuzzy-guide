package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"fuzzy/config"
	"fuzzy/models"
)

// Session management constants and variables
var (
	sessions = make(map[string]int) // sessionID -> userID
	loginAttempts = make(map[string][]time.Time) // IP -> attempt times
)

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
	cookie, err := r.Cookie(config.AppConfig.Security.SessionCookieName)
	if err == nil {
		// Remove session from store
		delete(sessions, cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     config.AppConfig.Security.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   config.AppConfig.Security.HTTPSEnabled,
		SameSite: http.SameSiteStrictMode,
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

	// Check rate limiting
	clientIP := getClientIP(r)
	if isRateLimited(clientIP) {
		data := models.LoginPageData{
			Title: "Fuzzy - Login",
			Error: "Too many login attempts. Please wait before trying again.",
		}
		renderLoginTemplate(w, &data)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	// Validate input
	if username == "" || password == "" {
		recordLoginAttempt(clientIP)
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
		recordLoginAttempt(clientIP)
		data := models.LoginPageData{
			Title: "Fuzzy - Login",
			Error: "Invalid username or password",
		}
		renderLoginTemplate(w, &data)
		return
	}

	// Check if user is active
	if !user.Active {
		recordLoginAttempt(clientIP)
		data := models.LoginPageData{
			Title: "Fuzzy - Login",
			Error: "Account is disabled",
		}
		renderLoginTemplate(w, &data)
		return
	}

	// Successful login - clear attempts
	clearLoginAttempts(clientIP)

	// Create session
	sessionID := generateSessionID()
	sessions[sessionID] = user.ID

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     config.AppConfig.Security.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(config.AppConfig.GetSessionDuration().Seconds()),
		HttpOnly: true,
		Secure:   config.AppConfig.Security.HTTPSEnabled,
		SameSite: http.SameSiteStrictMode,
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

	// Validate input
	if username == "" {
		data := models.SetupPageData{
			Title: "Fuzzy - First Time Setup",
			Error: "Username is required",
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

	// Validate password strength
	if err := validatePasswordStrength(password); err != nil {
		data := models.SetupPageData{
			Title: "Fuzzy - First Time Setup",
			Error: err.Error(),
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
		FirstName: "",
		LastName:  "",
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

// generateSessionID creates a cryptographically secure session ID
func generateSessionID() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to time-based ID if crypto/rand fails
		return time.Now().Format("20060102150405") + "-" + time.Now().Format("999999999")
	}
	return hex.EncodeToString(bytes)
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser(r *http.Request) (models.User, bool) {
	cookie, err := r.Cookie(config.AppConfig.Security.SessionCookieName)
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

// Rate limiting functions
func getClientIP(r *http.Request) string {
	// Check for forwarded headers first
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func isRateLimited(clientIP string) bool {
	attempts, exists := loginAttempts[clientIP]
	if !exists {
		return false
	}

	// Clean old attempts
	cutoff := time.Now().Add(-time.Duration(config.AppConfig.Limits.LoginTimeoutMinutes) * time.Minute)
	var validAttempts []time.Time
	for _, attempt := range attempts {
		if attempt.After(cutoff) {
			validAttempts = append(validAttempts, attempt)
		}
	}
	loginAttempts[clientIP] = validAttempts

	return len(validAttempts) >= config.AppConfig.Limits.MaxLoginAttempts
}

func recordLoginAttempt(clientIP string) {
	loginAttempts[clientIP] = append(loginAttempts[clientIP], time.Now())
}

func clearLoginAttempts(clientIP string) {
	delete(loginAttempts, clientIP)
}

// validatePasswordStrength checks if password meets security requirements
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("le mot de passe doit contenir au moins 8 caractères")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("le mot de passe doit contenir au moins une lettre majuscule")
	}
	if !hasLower {
		return fmt.Errorf("le mot de passe doit contenir au moins une lettre minuscule")
	}
	if !hasNumber {
		return fmt.Errorf("le mot de passe doit contenir au moins un chiffre")
	}
	if !hasSpecial {
		return fmt.Errorf("le mot de passe doit contenir au moins un caractère spécial (!@#$%^&*)")
	}

	return nil
}