package main

import (
	"log"
	"net/http"

	"fuzzy/handlers"
)

func main() {
	// Set up HTTP routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/setup", handlers.SetupHandler)
	
	// Protected routes
	http.HandleFunc("/providers", handlers.RequireSetupOrAuth(handlers.ProvidersHandler))
	http.HandleFunc("/channels", handlers.RequireSetupOrAuth(handlers.ChannelsHandler))
	http.HandleFunc("/users", handlers.RequireSetupOrAuth(handlers.UsersHandler))
	http.HandleFunc("/channel/start", handlers.RequireSetupOrAuth(handlers.ChannelStartHandler))
	http.HandleFunc("/channel/stop", handlers.RequireSetupOrAuth(handlers.ChannelStopHandler))

	// Server configuration
	const port = ":8080"
	
	// Log server startup
	log.Printf("Starting Fuzzy web server on port %s", port)
	log.Printf("Home page: http://localhost%s/", port)
	log.Printf("Health check: http://localhost%s/health", port)

	// Start the HTTP server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}