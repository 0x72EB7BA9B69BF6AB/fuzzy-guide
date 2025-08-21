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

// ChannelsHandler handles requests to the channels page for managing individual channels
func ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	var data models.ChannelsPageData
	data.Title = "Fuzzy - Channel Management"

	switch r.Method {
	case http.MethodGet:
		handleGetChannels(w, r, &data)
	case http.MethodPost:
		handlePostChannels(w, r, &data)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleGetChannels(w http.ResponseWriter, r *http.Request, data *models.ChannelsPageData) {
	// Check if this is a delete request
	if deleteID := r.URL.Query().Get("delete"); deleteID != "" {
		id, err := strconv.Atoi(deleteID)
		if err != nil {
			data.Error = "Invalid channel ID"
		} else if models.GlobalStore.DeleteChannel(id) {
			data.Message = "Channel deleted successfully"
		} else {
			data.Error = "Channel not found"
		}
	}

	// Get all channels
	data.Channels = models.GlobalStore.GetAllChannels()

	// Render template
	renderChannelsTemplate(w, data)
}

func handlePostChannels(w http.ResponseWriter, r *http.Request, data *models.ChannelsPageData) {
	if err := r.ParseForm(); err != nil {
		data.Error = "Failed to parse form data"
		data.Channels = models.GlobalStore.GetAllChannels()
		renderChannelsTemplate(w, data)
		return
	}

	action := r.FormValue("action")
	switch action {
	case "create":
		handleCreateChannel(r, data)
	case "update":
		handleUpdateChannel(r, data)
	default:
		data.Error = "Invalid action"
	}

	// Get updated channels list
	data.Channels = models.GlobalStore.GetAllChannels()
	renderChannelsTemplate(w, data)
}

func handleCreateChannel(r *http.Request, data *models.ChannelsPageData) {
	name := strings.TrimSpace(r.FormValue("name"))
	manifest := strings.TrimSpace(r.FormValue("manifest"))
	keyKid := strings.TrimSpace(r.FormValue("key_kid"))
	
	// Validate input
	if name == "" {
		data.Error = "Channel name is required"
		return
	}
	if manifest == "" {
		data.Error = "Channel manifest URL is required"
		return
	}
	if keyKid == "" {
		data.Error = "Channel Key:Kid is required"
		return
	}

	// Create channel
	channel := models.Channel{
		Name:     name,
		Manifest: manifest,
		KeyKid:   keyKid,
	}

	models.GlobalStore.CreateChannel(channel)
	data.Message = "Channel created successfully"
}

func handleUpdateChannel(r *http.Request, data *models.ChannelsPageData) {
	idStr := strings.TrimSpace(r.FormValue("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		data.Error = "Invalid channel ID"
		return
	}

	// Get existing channel
	existing, exists := models.GlobalStore.GetChannel(id)
	if !exists {
		data.Error = "Channel not found"
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	manifest := strings.TrimSpace(r.FormValue("manifest"))
	keyKid := strings.TrimSpace(r.FormValue("key_kid"))

	// Validate input
	if name == "" {
		data.Error = "Channel name is required"
		return
	}
	if manifest == "" {
		data.Error = "Channel manifest URL is required"
		return
	}
	if keyKid == "" {
		data.Error = "Channel Key:Kid is required"
		return
	}

	// Update channel
	updated := existing
	updated.Name = name
	updated.Manifest = manifest
	updated.KeyKid = keyKid

	if models.GlobalStore.UpdateChannel(updated) {
		data.Message = "Channel updated successfully"
	} else {
		data.Error = "Failed to update channel"
	}
}

func renderChannelsTemplate(w http.ResponseWriter, data *models.ChannelsPageData) {
	tmplPath := filepath.Join("templates", "channels.html")
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