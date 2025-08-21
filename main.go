package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// HomePageData represents the data structure for the home page template
type HomePageData struct {
	Title       string
	WelcomeMsg  string
	CurrentTime string
}

// homeHandler handles requests to the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests to the home page
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Prepare data for the home page template
	data := HomePageData{
		Title:       "Fuzzy Guide - Home",
		WelcomeMsg:  "Welcome to Fuzzy Guide!",
		CurrentTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Parse and execute the home page template
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
        }
        p {
            color: #666;
            line-height: 1.6;
        }
        .time {
            color: #888;
            font-size: 0.9em;
            text-align: center;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.WelcomeMsg}}</h1>
        <p>This is a clean, simple, and efficient Go web server. The server is designed to be lightweight and easy to understand.</p>
        <p>Features:</p>
        <ul>
            <li>Clean Go code with proper error handling</li>
            <li>Simple HTML template rendering</li>
            <li>Responsive design</li>
            <li>Efficient HTTP routing</li>
        </ul>
        <div class="time">Current server time: {{.CurrentTime}}</div>
    </div>
</body>
</html>`

	// Parse the template
	t, err := template.New("home").Parse(tmpl)
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

// healthHandler provides a simple health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"healthy","service":"fuzzy-guide"}`)
}

func main() {
	// Set up HTTP routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)

	// Server configuration
	const port = ":8080"
	
	// Log server startup
	log.Printf("Starting Fuzzy Guide web server on port %s", port)
	log.Printf("Home page: http://localhost%s/", port)
	log.Printf("Health check: http://localhost%s/health", port)

	// Start the HTTP server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}