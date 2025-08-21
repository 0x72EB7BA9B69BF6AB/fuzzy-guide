package models

import "time"

// HomePageData represents the data structure for the home page template
type HomePageData struct {
	Title       string
	WelcomeMsg  string
	CurrentTime string
}

// Bouquet represents a package/bundle that can be configured by providers
type Bouquet struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Channels    []string  `json:"channels"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// User represents a system user
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProvidersPageData represents the data structure for the providers page template
type ProvidersPageData struct {
	Title    string
	Bouquets []Bouquet
	Message  string
	Error    string
}

// UsersPageData represents the data structure for the users page template
type UsersPageData struct {
	Title   string
	Users   []User
	Message string
	Error   string
}