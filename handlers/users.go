package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"fuzzy/models"
)

// UsersHandler handles requests to the users page for managing users
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	var data models.UsersPageData
	data.Title = "Fuzzy - Users"

	switch r.Method {
	case http.MethodGet:
		handleGetUsers(w, r, &data)
	case http.MethodPost:
		handlePostUsers(w, r, &data)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleGetUsers(w http.ResponseWriter, r *http.Request, data *models.UsersPageData) {
	// Check if this is a delete request
	if deleteID := r.URL.Query().Get("delete"); deleteID != "" {
		id, err := strconv.Atoi(deleteID)
		if err != nil {
			data.Error = "Invalid user ID"
		} else if models.GlobalStore.DeleteUser(id) {
			data.Message = "User deleted successfully"
		} else {
			data.Error = "User not found"
		}
	}

	// Get all users
	data.Users = models.GlobalStore.GetAllUsers()

	// Render template
	renderUsersTemplate(w, data)
}

func handlePostUsers(w http.ResponseWriter, r *http.Request, data *models.UsersPageData) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		data.Error = "Error parsing form data"
		data.Users = models.GlobalStore.GetAllUsers()
		renderUsersTemplate(w, data)
		return
	}

	action := r.FormValue("action")
	
	switch action {
	case "create":
		handleCreateUser(r, data)
	case "update":
		handleUpdateUser(r, data)
	default:
		data.Error = "Invalid action"
	}

	// Get updated users list
	data.Users = models.GlobalStore.GetAllUsers()
	renderUsersTemplate(w, data)
}

func handleCreateUser(r *http.Request, data *models.UsersPageData) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	role := strings.TrimSpace(r.FormValue("role"))
	activeStr := r.FormValue("active")

	// Validate input
	if username == "" {
		data.Error = "Username is required"
		return
	}
	if email == "" {
		data.Error = "Email is required"
		return
	}
	if password == "" {
		data.Error = "Password is required"
		return
	}
	if role == "" {
		role = "User" // Default role
	}

	// Check if username already exists
	if _, exists := models.GlobalStore.GetUserByUsername(username); exists {
		data.Error = "Username already exists"
		return
	}

	active := activeStr == "on" || activeStr == "true"

	// Create user
	user := models.User{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		Active:    active,
	}

	// Set password
	if err := user.SetPassword(password); err != nil {
		data.Error = "Error encrypting password"
		return
	}

	models.GlobalStore.CreateUser(user)
	data.Message = "User created successfully"
}

func handleUpdateUser(r *http.Request, data *models.UsersPageData) {
	idStr := strings.TrimSpace(r.FormValue("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		data.Error = "Invalid user ID"
		return
	}

	// Get existing user
	existing, exists := models.GlobalStore.GetUser(id)
	if !exists {
		data.Error = "User not found"
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	role := strings.TrimSpace(r.FormValue("role"))
	activeStr := r.FormValue("active")

	// Validate input
	if username == "" {
		data.Error = "Username is required"
		return
	}
	if email == "" {
		data.Error = "Email is required"
		return
	}
	if role == "" {
		role = "User" // Default role
	}

	active := activeStr == "on" || activeStr == "true"

	// Update user
	updated := existing
	updated.Username = username
	updated.Email = email
	updated.FirstName = firstName
	updated.LastName = lastName
	updated.Role = role
	updated.Active = active

	if models.GlobalStore.UpdateUser(updated) {
		data.Message = "User updated successfully"
	} else {
		data.Error = "Failed to update user"
	}
}

func renderUsersTemplate(w http.ResponseWriter, data *models.UsersPageData) {
	// Parse the template file
	tmplPath := filepath.Join("templates", "users.html")
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set content type and execute template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}