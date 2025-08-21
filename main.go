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
	http.HandleFunc("/providers", handlers.ProvidersHandler)
	http.HandleFunc("/providers-management", handlers.ProvidersManagementHandler)
	http.HandleFunc("/channels", handlers.ChannelsHandler)
	http.HandleFunc("/users", handlers.UsersHandler)

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