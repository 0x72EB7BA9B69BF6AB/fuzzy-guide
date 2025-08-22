package handlers

import (
	"net/http"
	
	"fuzzy/models"
)

// RequireAuth is middleware that requires user authentication
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		_, authenticated := GetCurrentUser(r)
		if !authenticated {
			// Redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// User is authenticated, proceed to next handler
		next(w, r)
	}
}

// RequireSetupOrAuth redirects to setup if no users exist, otherwise requires auth
func RequireSetupOrAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If no users exist, redirect to setup
		if !models.GlobalStore.HasUsers() {
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
			return
		}

		// Users exist, require authentication
		RequireAuth(next)(w, r)
	}
}

// RedirectIfAuthenticated redirects authenticated users away from login/setup pages
func RedirectIfAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		_, authenticated := GetCurrentUser(r)
		if authenticated {
			// User is already logged in, redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// User is not authenticated, proceed to next handler
		next(w, r)
	}
}