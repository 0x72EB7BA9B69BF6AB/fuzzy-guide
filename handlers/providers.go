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

// ProvidersHandler handles requests to the providers page for managing bouquets
func ProvidersHandler(w http.ResponseWriter, r *http.Request) {
	var data models.ProvidersPageData
	data.Title = "Fuzzy - Administration Panel"

	switch r.Method {
	case http.MethodGet:
		handleGetProviders(w, r, &data)
	case http.MethodPost:
		handlePostProviders(w, r, &data)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleGetProviders(w http.ResponseWriter, r *http.Request, data *models.ProvidersPageData) {
	// Check if this is a delete request
	if deleteID := r.URL.Query().Get("delete"); deleteID != "" {
		id, err := strconv.Atoi(deleteID)
		if err != nil {
			data.Error = "Invalid bouquet ID"
		} else if models.GlobalStore.DeleteBouquet(id) {
			data.Message = "Bouquet deleted successfully"
		} else {
			data.Error = "Bouquet not found"
		}
	}

	// Get all bouquets
	data.Bouquets = models.GlobalStore.GetAllBouquets()

	// Render template
	renderProvidersTemplate(w, data)
}

func handlePostProviders(w http.ResponseWriter, r *http.Request, data *models.ProvidersPageData) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		data.Error = "Error parsing form data"
		data.Bouquets = models.GlobalStore.GetAllBouquets()
		renderProvidersTemplate(w, data)
		return
	}

	action := r.FormValue("action")
	
	switch action {
	case "create":
		handleCreateBouquet(r, data)
	case "update":
		handleUpdateBouquet(r, data)
	default:
		data.Error = "Invalid action"
	}

	// Get updated bouquets list
	data.Bouquets = models.GlobalStore.GetAllBouquets()
	renderProvidersTemplate(w, data)
}

func handleCreateBouquet(r *http.Request, data *models.ProvidersPageData) {
	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))
	
	// Validate input
	if name == "" {
		data.Error = "Bouquet name is required"
		return
	}

	// Parse channels from form data - for now, we'll use a simple format
	// Later we can enhance this with a more sophisticated channel input method
	var channels []models.Channel
	
	// Get channel data from form - we'll implement multiple channels later
	// For now, let's handle a single channel input
	channelName := strings.TrimSpace(r.FormValue("channel_name"))
	channelManifest := strings.TrimSpace(r.FormValue("channel_manifest"))
	channelKeyKid := strings.TrimSpace(r.FormValue("channel_keykid"))
	
	if channelName != "" && channelManifest != "" && channelKeyKid != "" {
		channels = append(channels, models.Channel{
			Name:     channelName,
			Manifest: channelManifest,
			KeyKid:   channelKeyKid,
		})
	}

	// Create bouquet
	bouquet := models.Bouquet{
		Name:        name,
		Description: description,
		Channels:    channels,
	}

	models.GlobalStore.CreateBouquet(bouquet)
	data.Message = "Bouquet created successfully"
}

func handleUpdateBouquet(r *http.Request, data *models.ProvidersPageData) {
	idStr := strings.TrimSpace(r.FormValue("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		data.Error = "Invalid bouquet ID"
		return
	}

	// Get existing bouquet
	existing, exists := models.GlobalStore.GetBouquet(id)
	if !exists {
		data.Error = "Bouquet not found"
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))

	// Validate input
	if name == "" {
		data.Error = "Bouquet name is required"
		return
	}

	// Parse channels from form data - for now, we'll use a simple format
	var channels []models.Channel
	
	// Get channel data from form - we'll implement multiple channels later
	channelName := strings.TrimSpace(r.FormValue("channel_name"))
	channelManifest := strings.TrimSpace(r.FormValue("channel_manifest"))
	channelKeyKid := strings.TrimSpace(r.FormValue("channel_keykid"))
	
	if channelName != "" && channelManifest != "" && channelKeyKid != "" {
		channels = append(channels, models.Channel{
			Name:     channelName,
			Manifest: channelManifest,
			KeyKid:   channelKeyKid,
		})
	}

	// Update bouquet
	updated := existing
	updated.Name = name
	updated.Description = description
	updated.Channels = channels

	if models.GlobalStore.UpdateBouquet(updated) {
		data.Message = "Bouquet updated successfully"
	} else {
		data.Error = "Failed to update bouquet"
	}
}

func renderProvidersTemplate(w http.ResponseWriter, data *models.ProvidersPageData) {
	// Parse the template file
	tmplPath := filepath.Join("templates", "providers.html")
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