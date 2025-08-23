package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Security SecurityConfig
	Database DatabaseConfig
	Logging  LoggingConfig
	UI       UIConfig
	Limits   LimitsConfig
	Features FeaturesConfig
}

type ServerConfig struct {
	Port     int
	AppName  string
	Version  string
	DevMode  bool
}

type SecurityConfig struct {
	SessionCookieName     string
	SessionDurationHours  int
	SecretKey             string
	HTTPSEnabled          bool
	CSRFEnabled           bool
}

type DatabaseConfig struct {
	Type     string
	DataFile string
}

type LoggingConfig struct {
	Level   string
	File    string
	Console bool
}

type UIConfig struct {
	Theme    string
	Language string
	DarkMode bool
}

type LimitsConfig struct {
	MaxLoginAttempts      int
	LoginTimeoutMinutes   int
	MaxUploadSizeMB       int
}

type FeaturesConfig struct {
	UserManagement     bool
	ProviderManagement bool
	ChannelManagement  bool
}

// Global configuration instance
var AppConfig *Config

// LoadConfig loads configuration from the specified file
func LoadConfig(filename string) (*Config, error) {
	config := &Config{
		// Set default values
		Server: ServerConfig{
			Port:    8080,
			AppName: "Fuzzy",
			Version: "1.0.0",
			DevMode: false,
		},
		Security: SecurityConfig{
			SessionCookieName:    "fuzzy_session",
			SessionDurationHours: 24,
			SecretKey:            "changeme_in_production",
			HTTPSEnabled:         false,
			CSRFEnabled:          true,
		},
		Database: DatabaseConfig{
			Type:     "memory",
			DataFile: "data/fuzzy.db",
		},
		Logging: LoggingConfig{
			Level:   "info",
			File:    "logs/fuzzy.log",
			Console: true,
		},
		UI: UIConfig{
			Theme:    "blue",
			Language: "fr",
			DarkMode: false,
		},
		Limits: LimitsConfig{
			MaxLoginAttempts:    5,
			LoginTimeoutMinutes: 15,
			MaxUploadSizeMB:     100,
		},
		Features: FeaturesConfig{
			UserManagement:     true,
			ProviderManagement: true,
			ChannelManagement:  true,
		},
	}

	// Check if config file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create default config file
		if err := SaveConfig(config, filename); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			continue
		}

		// Parse key-value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if err := setConfigValue(config, currentSection, key, value); err != nil {
				return nil, fmt.Errorf("error setting config value %s.%s: %v", currentSection, key, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	return config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()

	// Write config with comments
	fmt.Fprintln(file, "# Fuzzy Application Configuration")
	fmt.Fprintln(file, "# Ce fichier contient toutes les configurations de l'application")
	fmt.Fprintln(file, "# This file contains all application configurations")
	fmt.Fprintln(file, "")

	// Server section
	fmt.Fprintln(file, "[server]")
	fmt.Fprintln(file, "# Port d'écoute du serveur / Server listening port")
	fmt.Fprintf(file, "port = %d\n", config.Server.Port)
	fmt.Fprintln(file, "# Nom de l'application / Application name")
	fmt.Fprintf(file, "app_name = %s\n", config.Server.AppName)
	fmt.Fprintln(file, "# Version de l'application / Application version")
	fmt.Fprintf(file, "version = %s\n", config.Server.Version)
	fmt.Fprintln(file, "# Mode de développement / Development mode (true/false)")
	fmt.Fprintf(file, "dev_mode = %t\n", config.Server.DevMode)
	fmt.Fprintln(file, "")

	// Security section
	fmt.Fprintln(file, "[security]")
	fmt.Fprintln(file, "# Nom du cookie de session / Session cookie name")
	fmt.Fprintf(file, "session_cookie_name = %s\n", config.Security.SessionCookieName)
	fmt.Fprintln(file, "# Durée de session en heures / Session duration in hours")
	fmt.Fprintf(file, "session_duration_hours = %d\n", config.Security.SessionDurationHours)
	fmt.Fprintln(file, "# Clé secrète pour la sécurité / Secret key for security")
	fmt.Fprintf(file, "secret_key = %s\n", config.Security.SecretKey)
	fmt.Fprintln(file, "# HTTPS activé / HTTPS enabled (true/false)")
	fmt.Fprintf(file, "https_enabled = %t\n", config.Security.HTTPSEnabled)
	fmt.Fprintln(file, "# CSRF protection activé / CSRF protection enabled (true/false)")
	fmt.Fprintf(file, "csrf_enabled = %t\n", config.Security.CSRFEnabled)
	fmt.Fprintln(file, "")

	// Add other sections...
	return nil
}

// setConfigValue sets a configuration value based on section and key
func setConfigValue(config *Config, section, key, value string) error {
	switch section {
	case "server":
		return setServerConfig(&config.Server, key, value)
	case "security":
		return setSecurityConfig(&config.Security, key, value)
	case "database":
		return setDatabaseConfig(&config.Database, key, value)
	case "logging":
		return setLoggingConfig(&config.Logging, key, value)
	case "ui":
		return setUIConfig(&config.UI, key, value)
	case "limits":
		return setLimitsConfig(&config.Limits, key, value)
	case "features":
		return setFeaturesConfig(&config.Features, key, value)
	}
	return nil
}

func setServerConfig(config *ServerConfig, key, value string) error {
	switch key {
	case "port":
		port, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		config.Port = port
	case "app_name":
		config.AppName = value
	case "version":
		config.Version = value
	case "dev_mode":
		devMode, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		config.DevMode = devMode
	}
	return nil
}

func setSecurityConfig(config *SecurityConfig, key, value string) error {
	switch key {
	case "session_cookie_name":
		config.SessionCookieName = value
	case "session_duration_hours":
		hours, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		config.SessionDurationHours = hours
	case "secret_key":
		config.SecretKey = value
	case "https_enabled":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		config.HTTPSEnabled = enabled
	case "csrf_enabled":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		config.CSRFEnabled = enabled
	}
	return nil
}

func setDatabaseConfig(config *DatabaseConfig, key, value string) error {
	switch key {
	case "type":
		config.Type = value
	case "data_file":
		config.DataFile = value
	}
	return nil
}

func setLoggingConfig(config *LoggingConfig, key, value string) error {
	switch key {
	case "level":
		config.Level = value
	case "file":
		config.File = value
	case "console":
		console, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		config.Console = console
	}
	return nil
}

func setUIConfig(config *UIConfig, key, value string) error {
	switch key {
	case "theme":
		config.Theme = value
	case "language":
		config.Language = value
	case "dark_mode":
		darkMode, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		config.DarkMode = darkMode
	}
	return nil
}

func setLimitsConfig(config *LimitsConfig, key, value string) error {
	switch key {
	case "max_login_attempts":
		attempts, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		config.MaxLoginAttempts = attempts
	case "login_timeout_minutes":
		timeout, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		config.LoginTimeoutMinutes = timeout
	case "max_upload_size_mb":
		size, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		config.MaxUploadSizeMB = size
	}
	return nil
}

func setFeaturesConfig(config *FeaturesConfig, key, value string) error {
	switch key {
	case "user_management":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		config.UserManagement = enabled
	case "provider_management":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		config.ProviderManagement = enabled
	case "channel_management":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		config.ChannelManagement = enabled
	}
	return nil
}

// GetSessionDuration returns the session duration as time.Duration
func (c *Config) GetSessionDuration() time.Duration {
	return time.Duration(c.Security.SessionDurationHours) * time.Hour
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

// Initialize loads the global configuration
func Initialize() error {
	var err error
	AppConfig, err = LoadConfig("config/config.cfg")
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}
	return nil
}