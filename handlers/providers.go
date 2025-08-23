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

// ProvidersHandler handles requests to the main providers page with integrated bouquets
func ProvidersHandler(w http.ResponseWriter, r *http.Request) {
	var data models.ProvidersWithBouquetsPageData
	data.Title = "Fuzzy - Providers & Bouquets"

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

func handleGetProviders(w http.ResponseWriter, r *http.Request, data *models.ProvidersWithBouquetsPageData) {
	// Check if this is a delete request for provider
	if deleteID := r.URL.Query().Get("delete-provider"); deleteID != "" {
		id, err := strconv.Atoi(deleteID)
		if err != nil {
			data.Error = "Invalid provider ID"
		} else if models.GlobalStore.DeleteProvider(id) {
			data.Message = "Provider deleted successfully"
		} else {
			data.Error = "Provider not found"
		}
	}

	// Check if this is a delete request for bouquet
	if deleteID := r.URL.Query().Get("delete-bouquet"); deleteID != "" {
		id, err := strconv.Atoi(deleteID)
		if err != nil {
			data.Error = "Invalid bouquet ID"
		} else if models.GlobalStore.DeleteBouquet(id) {
			data.Message = "Bouquet deleted successfully"
		} else {
			data.Error = "Bouquet not found"
		}
	}

	// Get all providers with their bouquets
	data.Providers = models.GlobalStore.GetProvidersWithBouquets()

	// Render template
	renderProvidersTemplate(w, data)
}

func handlePostProviders(w http.ResponseWriter, r *http.Request, data *models.ProvidersWithBouquetsPageData) {
	if err := r.ParseForm(); err != nil {
		data.Error = "Failed to parse form data"
		data.Providers = models.GlobalStore.GetProvidersWithBouquets()
		renderProvidersTemplate(w, data)
		return
	}

	action := r.FormValue("action")
	switch action {
	case "create-provider":
		handleCreateProvider(r, data)
	case "update-provider":
		handleUpdateProvider(r, data)
	case "create-bouquet":
		handleCreateBouquet(r, data)
	case "update-bouquet":
		handleUpdateBouquet(r, data)
	default:
		data.Error = "Invalid action"
	}

	// Get updated providers list
	data.Providers = models.GlobalStore.GetProvidersWithBouquets()
	renderProvidersTemplate(w, data)
}

func handleCreateProvider(r *http.Request, data *models.ProvidersWithBouquetsPageData) {
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

func handleUpdateProvider(r *http.Request, data *models.ProvidersWithBouquetsPageData) {
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

func handleCreateBouquet(r *http.Request, data *models.ProvidersWithBouquetsPageData) {
	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))
	providerIDStr := strings.TrimSpace(r.FormValue("provider_id"))
	
	// Validate input
	if name == "" {
		data.Error = "Bouquet name is required"
		return
	}
	
	providerID, err := strconv.Atoi(providerIDStr)
	if err != nil {
		data.Error = "Invalid provider ID"
		return
	}

	// Parse channels from form data
	var channels []models.Channel
	
	// Get channel data from form - we'll implement multiple channels later
	channelName := strings.TrimSpace(r.FormValue("channel_name"))
	channelManifest := strings.TrimSpace(r.FormValue("channel_manifest"))
	channelKeyKid := strings.TrimSpace(r.FormValue("channel_keykid"))
	videoCodec := strings.TrimSpace(r.FormValue("channel_video_codec"))
	audioCodec := strings.TrimSpace(r.FormValue("channel_audio_codec"))
	resolution := strings.TrimSpace(r.FormValue("channel_resolution"))
	videoBitrate := strings.TrimSpace(r.FormValue("channel_video_bitrate"))
	audioBitrate := strings.TrimSpace(r.FormValue("channel_audio_bitrate"))
	quality := strings.TrimSpace(r.FormValue("channel_quality"))
	
	if channelName != "" && channelManifest != "" && channelKeyKid != "" {
		// Set defaults for video encoding fields
		if videoCodec == "" {
			videoCodec = "x265"
		}
		if audioCodec == "" {
			audioCodec = "AAC"
		}
		if resolution == "" {
			resolution = "1080p"
		}
		if videoBitrate == "" {
			videoBitrate = "5000k"
		}
		if audioBitrate == "" {
			audioBitrate = "128k"
		}
		if quality == "" {
			quality = "High"
		}
		
		channels = append(channels, models.Channel{
			Name:         channelName,
			Manifest:     channelManifest,
			KeyKid:       channelKeyKid,
			VideoCodec:   videoCodec,
			AudioCodec:   audioCodec,
			Resolution:   resolution,
			VideoBitrate: videoBitrate,
			AudioBitrate: audioBitrate,
			Quality:      quality,
			Running:      false,
			RemuxPort:    0,
		})
	}

	// Create bouquet
	bouquet := models.Bouquet{
		Name:        name,
		Description: description,
		ProviderID:  providerID,
		Channels:    channels,
	}

	models.GlobalStore.CreateBouquet(bouquet)
	data.Message = "Bouquet created successfully"
}

func handleUpdateBouquet(r *http.Request, data *models.ProvidersWithBouquetsPageData) {
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

	// Update bouquet (keeping existing channels for now)
	updated := existing
	updated.Name = name
	updated.Description = description

	if models.GlobalStore.UpdateBouquet(updated) {
		data.Message = "Bouquet updated successfully"
	} else {
		data.Error = "Failed to update bouquet"
	}
}

func renderProvidersTemplate(w http.ResponseWriter, data *models.ProvidersWithBouquetsPageData) {
	tmplPath := filepath.Join("templates", "providers.html")
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

// ChannelStartHandler handles channel start requests
func ChannelStartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	channelIDStr := r.FormValue("channel_id")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	port, success := models.GlobalStore.StartChannel(channelID)
	if !success {
		http.Error(w, "Failed to start channel", http.StatusInternalServerError)
		return
	}

	// Update channel in bouquets
	models.GlobalStore.UpdateChannelInBouquet(channelID)

	log.Printf("Channel %d started on port %d", channelID, port)
	http.Redirect(w, r, "/providers", http.StatusSeeOther)
}

// ChannelStopHandler handles channel stop requests
func ChannelStopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	channelIDStr := r.FormValue("channel_id")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	success := models.GlobalStore.StopChannel(channelID)
	if !success {
		http.Error(w, "Failed to stop channel", http.StatusInternalServerError)
		return
	}

	// Update channel in bouquets
	models.GlobalStore.UpdateChannelInBouquet(channelID)

	log.Printf("Channel %d stopped", channelID)
	http.Redirect(w, r, "/providers", http.StatusSeeOther)
}