package main

import (
	"log"
	"net/http"
	"os"

	"fuzzy/config"
	"fuzzy/handlers"
)

func main() {
	// Initialize configuration
	if err := config.Initialize(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Create necessary directories
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Warning: Failed to create logs directory: %v", err)
	}
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Printf("Warning: Failed to create data directory: %v", err)
	}

	// Set up HTTP routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/setup", handlers.SetupHandler)
	
	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	
	// Protected routes
	http.HandleFunc("/providers", handlers.RequireSetupOrAuth(handlers.ProvidersHandler))
	http.HandleFunc("/channels", handlers.RequireSetupOrAuth(handlers.ChannelsHandler))
	http.HandleFunc("/users", handlers.RequireSetupOrAuth(handlers.UsersHandler))
	http.HandleFunc("/channel/start", handlers.RequireSetupOrAuth(handlers.ChannelStartHandler))
	http.HandleFunc("/channel/stop", handlers.RequireSetupOrAuth(handlers.ChannelStopHandler))

	// Get server configuration
	serverAddr := config.AppConfig.GetServerAddress()
	appName := config.AppConfig.Server.AppName
	version := config.AppConfig.Server.Version
	
	// Log server startup
	log.Printf("Starting %s v%s web server on port %s", appName, version, serverAddr)
	log.Printf("Home page: http://localhost%s/", serverAddr)
	log.Printf("Health check: http://localhost%s/health", serverAddr)
	log.Printf("Configuration loaded from: config/config.cfg")

	// Start the HTTP server
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}