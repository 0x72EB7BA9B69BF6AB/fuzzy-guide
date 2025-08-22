package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"fuzzy/models"
)

// HomeHandler handles requests to the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests to the home page
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// If no users exist, redirect to setup
	if !models.GlobalStore.HasUsers() {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	// Check if user is authenticated
	user, authenticated := GetCurrentUser(r)
	if !authenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Prepare data for the home page template
	welcomeMsg := "Welcome to Fuzzy!"
	if user.FirstName != "" {
		welcomeMsg = "Welcome back, " + user.FirstName + "!"
	}

	data := models.HomePageData{
		Title:       "Fuzzy - Home",
		WelcomeMsg:  welcomeMsg,
		CurrentTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Parse the template file
	tmplPath := filepath.Join("templates", "home.html")
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