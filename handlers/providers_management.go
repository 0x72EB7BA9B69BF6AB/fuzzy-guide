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

// ProvidersManagementHandler handles requests to the providers management page
func ProvidersManagementHandler(w http.ResponseWriter, r *http.Request) {
	var data models.ProvidersManagementPageData
	data.Title = "Fuzzy - Provider Management"

	switch r.Method {
	case http.MethodGet:
		handleGetProvidersManagement(w, r, &data)
	case http.MethodPost:
		handlePostProvidersManagement(w, r, &data)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleGetProvidersManagement(w http.ResponseWriter, r *http.Request, data *models.ProvidersManagementPageData) {
	// Check if this is a delete request
	if deleteID := r.URL.Query().Get("delete"); deleteID != "" {
		id, err := strconv.Atoi(deleteID)
		if err != nil {
			data.Error = "Invalid provider ID"
		} else if models.GlobalStore.DeleteProvider(id) {
			data.Message = "Provider deleted successfully"
		} else {
			data.Error = "Provider not found"
		}
	}

	// Get all providers
	data.Providers = models.GlobalStore.GetAllProviders()

	// Render template
	renderProvidersManagementTemplate(w, data)
}

func handlePostProvidersManagement(w http.ResponseWriter, r *http.Request, data *models.ProvidersManagementPageData) {
	if err := r.ParseForm(); err != nil {
		data.Error = "Failed to parse form data"
		data.Providers = models.GlobalStore.GetAllProviders()
		renderProvidersManagementTemplate(w, data)
		return
	}

	action := r.FormValue("action")
	switch action {
	case "create":
		handleCreateProvider(r, data)
	case "update":
		handleUpdateProvider(r, data)
	default:
		data.Error = "Invalid action"
	}

	// Get updated providers list
	data.Providers = models.GlobalStore.GetAllProviders()
	renderProvidersManagementTemplate(w, data)
}

func handleCreateProvider(r *http.Request, data *models.ProvidersManagementPageData) {
	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))
	url := strings.TrimSpace(r.FormValue("url"))
	apiKey := strings.TrimSpace(r.FormValue("api_key"))
	active := r.FormValue("active") == "on"
	
	// Validate input
	if name == "" {
		data.Error = "Provider name is required"
		return
	}
	if url == "" {
		data.Error = "Provider URL is required"
		return
	}

	// Create provider
	provider := models.Provider{
		Name:        name,
		Description: description,
		URL:         url,
		APIKey:      apiKey,
		Active:      active,
	}

	models.GlobalStore.CreateProvider(provider)
	data.Message = "Provider created successfully"
}

func handleUpdateProvider(r *http.Request, data *models.ProvidersManagementPageData) {
	idStr := strings.TrimSpace(r.FormValue("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		data.Error = "Invalid provider ID"
		return
	}

	// Get existing provider
	existing, exists := models.GlobalStore.GetProvider(id)
	if !exists {
		data.Error = "Provider not found"
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))
	url := strings.TrimSpace(r.FormValue("url"))
	apiKey := strings.TrimSpace(r.FormValue("api_key"))
	active := r.FormValue("active") == "on"

	// Validate input
	if name == "" {
		data.Error = "Provider name is required"
		return
	}
	if url == "" {
		data.Error = "Provider URL is required"
		return
	}

	// Update provider
	updated := existing
	updated.Name = name
	updated.Description = description
	updated.URL = url
	updated.APIKey = apiKey
	updated.Active = active

	if models.GlobalStore.UpdateProvider(updated) {
		data.Message = "Provider updated successfully"
	} else {
		data.Error = "Failed to update provider"
	}
}

func renderProvidersManagementTemplate(w http.ResponseWriter, data *models.ProvidersManagementPageData) {
	tmplPath := filepath.Join("templates", "providers_management.html")
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}