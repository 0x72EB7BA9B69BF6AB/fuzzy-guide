package models

import (
	"time"
	"golang.org/x/crypto/bcrypt"
)

// HomePageData represents the data structure for the home page template
type HomePageData struct {
	Title       string
	WelcomeMsg  string
	CurrentTime string
}

// Channel represents a channel with its execution properties and video encoding settings
type Channel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Manifest    string `json:"manifest"`
	KeyKid      string `json:"key_kid"`
	
	// Video encoding properties
	VideoCodec    string `json:"video_codec"`    // x264, x265, AV1, VP9
	AudioCodec    string `json:"audio_codec"`    // AAC, MP3, AC3, DTS
	Resolution    string `json:"resolution"`     // 720p, 1080p, 1440p, 2160p
	VideoBitrate  string `json:"video_bitrate"`  // Video bitrate (e.g., "5000k", "8000k")
	AudioBitrate  string `json:"audio_bitrate"`  // Audio bitrate (e.g., "128k", "256k")
	Quality       string `json:"quality"`        // Low, Medium, High, Ultra
	
	// Channel state for remuxer control
	Running   bool `json:"running"`   // Whether the channel is currently running
	RemuxPort int  `json:"remux_port"` // Port for internal remuxer (0 if not running)
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Provider represents a streaming provider
type Provider struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	APIKey      string    `json:"api_key"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Bouquet represents a package/bundle that can be configured by administrators
type Bouquet struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProviderID  int       `json:"provider_id"`  // Links bouquet to a provider
	Channels    []Channel `json:"channels"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// User represents a system user
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't include in JSON output
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ProvidersPageData represents the data structure for the providers page template
type ProvidersPageData struct {
	Title    string
	Bouquets []Bouquet
	Message  string
	Error    string
}

// ProvidersWithBouquetsPageData represents the data structure for the new providers page with integrated bouquets
type ProvidersWithBouquetsPageData struct {
	Title     string
	Providers []ProviderWithBouquets
	Message   string
	Error     string
}

// ProviderWithBouquets combines provider info with its bouquets
type ProviderWithBouquets struct {
	Provider
	Bouquets []Bouquet `json:"bouquets"`
}

// UsersPageData represents the data structure for the users page template
type UsersPageData struct {
	Title   string
	Users   []User
	Message string
	Error   string
}

// ChannelsPageData represents the data structure for the channels page template
type ChannelsPageData struct {
	Title    string
	Channels []Channel
	Message  string
	Error    string
}

// ProvidersManagementPageData represents the data structure for the providers management page template
type ProvidersManagementPageData struct {
	Title     string
	Providers []Provider
	Message   string
	Error     string
}

// LoginPageData represents the data structure for the login page template
type LoginPageData struct {
	Title   string
	Message string
	Error   string
}

// SetupPageData represents the data structure for the first-time setup page template
type SetupPageData struct {
	Title   string
	Message string
	Error   string
}