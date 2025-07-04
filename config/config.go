package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	AppPort     string
	JWTSecret   string
	// Security configurations
	JWTExpiry   time.Duration
	Environment string
	// Timezone configuration
	Timezone    string
	TimezoneLoc *time.Location
	// CORS configurations
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
	// Rate limiting configurations
	RateLimitMax     int
	RateLimitWindow  time.Duration
	RateLimitEnabled bool
	// Security configurations
	SecurityHeadersEnabled bool
	TrustedProxies         string
	// Database configurations
	DBMaxIdleConns int
	DBMaxOpenConns int
	DBMaxLifetime  time.Duration
	// Server timeout configurations
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	// Logging configurations
	// Debug/Development logging
	EnableDatabaseLog bool
	EnableAppLog      bool
	LogLevel          string
	// Mail configurations
	SMTPHost    string
	SMTPPort    int
	MailUsername string
	MailPassword string
	FromName     string
	FromEmail    string
}

var GlobalConfig *Config

// InitConfig initializes viper configuration from config.yaml only
func InitConfig() {
	// Set config file name and paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.go-api")

	// Set default values
	setDefaults()

	// Read config file (required)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("Config file 'config.yaml' not found. Please copy from config.example.yaml")
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Build config struct
	buildConfig()

	// Validate critical configurations
	validateConfig()

	log.Println("Configuration loaded successfully")
}

func setDefaults() {
	// Database defaults
	viper.SetDefault("database.url", "postgres://username:password@localhost:5432/database_name")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_lifetime", time.Hour)

	// Security defaults - NO HARDCODED SECRETS
	// These MUST be set in config.yaml or environment variables
	viper.SetDefault("security.jwt_secret", "")
	viper.SetDefault("security.jwt_expiry", 24*time.Hour)
	viper.SetDefault("security.headers_enabled", true)
	viper.SetDefault("security.trusted_proxies", "")

	// App defaults - Secure defaults for production
	viper.SetDefault("app.port", "8000")
	viper.SetDefault("app.environment", "production") // Default to production for security
	viper.SetDefault("app.timezone", "Asia/Jakarta")
	viper.SetDefault("app.read_timeout", 30*time.Second)
	viper.SetDefault("app.write_timeout", 30*time.Second)
	viper.SetDefault("app.idle_timeout", 120*time.Second)
	viper.SetDefault("app.shutdown_timeout", 10*time.Second)

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", "http://localhost:3000")
	viper.SetDefault("cors.allowed_methods", "GET,POST,PUT,DELETE,OPTIONS")
	viper.SetDefault("cors.allowed_headers", "Origin,Content-Type,Accept,Authorization")

	// Rate limiting defaults
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.max", 100)
	viper.SetDefault("rate_limit.window", time.Minute)

	// Logging defaults - Secure defaults
	viper.SetDefault("logging.enable_database_log", false) // Default off for production
	viper.SetDefault("logging.enable_app_log", false)      // Default off for production
	viper.SetDefault("logging.level", "info")              // info level default

	// Mail configuration defaults
	viper.SetDefault("mail.smtp_host", "smtp.gmail.com")
	viper.SetDefault("mail.smtp_port", 587)
	viper.SetDefault("mail.username", "")
	viper.SetDefault("mail.password", "")
	viper.SetDefault("mail.from_name", "Go API App")
	viper.SetDefault("mail.from_email", "")
}

func buildConfig() {
	GlobalConfig = &Config{
		// Database from config.yaml
		DatabaseURL: viper.GetString("database.url"),

		// Security from config.yaml
		JWTSecret:              viper.GetString("security.jwt_secret"),
		JWTExpiry:              viper.GetDuration("security.jwt_expiry"),
		SecurityHeadersEnabled: viper.GetBool("security.headers_enabled"),
		TrustedProxies:         viper.GetString("security.trusted_proxies"),

		// App configurations
		AppPort:         viper.GetString("app.port"),
		Environment:     viper.GetString("app.environment"),
		Timezone:        viper.GetString("app.timezone"),
		ReadTimeout:     viper.GetDuration("app.read_timeout"),
		WriteTimeout:    viper.GetDuration("app.write_timeout"),
		IdleTimeout:     viper.GetDuration("app.idle_timeout"),
		ShutdownTimeout: viper.GetDuration("app.shutdown_timeout"),

		// CORS configurations
		AllowedOrigins: viper.GetString("cors.allowed_origins"),
		AllowedMethods: viper.GetString("cors.allowed_methods"),
		AllowedHeaders: viper.GetString("cors.allowed_headers"),

		// Rate limiting configurations
		RateLimitEnabled: viper.GetBool("rate_limit.enabled"),
		RateLimitMax:     viper.GetInt("rate_limit.max"),
		RateLimitWindow:  viper.GetDuration("rate_limit.window"),

		// Database configurations
		DBMaxIdleConns: viper.GetInt("database.max_idle_conns"),
		DBMaxOpenConns: viper.GetInt("database.max_open_conns"),
		DBMaxLifetime:  viper.GetDuration("database.max_lifetime"),

		// Logging configurations
		EnableDatabaseLog: viper.GetBool("logging.enable_database_log"),
		EnableAppLog:      viper.GetBool("logging.enable_app_log"),
		LogLevel:          viper.GetString("logging.level"),

		// Mail configurations
		SMTPHost:     viper.GetString("mail.smtp_host"),
		SMTPPort:     viper.GetInt("mail.smtp_port"),
		MailUsername: viper.GetString("mail.username"),
		MailPassword: viper.GetString("mail.password"),
		FromName:     viper.GetString("mail.from_name"),
		FromEmail:    viper.GetString("mail.from_email"),
	}

	// Load timezone location
	if loc, err := time.LoadLocation(GlobalConfig.Timezone); err != nil {
		log.Printf("Warning: Invalid timezone '%s', using UTC instead. Error: %v", GlobalConfig.Timezone, err)
		GlobalConfig.TimezoneLoc = time.UTC
		GlobalConfig.Timezone = "UTC"
	} else {
		GlobalConfig.TimezoneLoc = loc
		log.Printf("Timezone loaded successfully: %s", GlobalConfig.Timezone)
	}
}

func validateConfig() {
	requiredConfigs := map[string]string{
		"database.url":        GlobalConfig.DatabaseURL,
		"security.jwt_secret": GlobalConfig.JWTSecret,
	}

	var missingConfigs []string
	for key, value := range requiredConfigs {
		// Check if values are empty or still contain default placeholder values
		if value == "" ||
			value == "postgres://username:password@localhost:5432/database_name" {
			missingConfigs = append(missingConfigs, key)
		}
	}

	if len(missingConfigs) > 0 {
		log.Fatalf("Required configurations must be set in config.yaml: %v\n"+
			"These cannot be empty or use default values for security reasons.", missingConfigs)
	}

	// Additional security validations
	if len(GlobalConfig.JWTSecret) < 32 {
		log.Fatalf("JWT secret must be at least 32 characters long for security")
	}

	// Validate database URL format and SSL requirements
	if !strings.Contains(GlobalConfig.DatabaseURL, "sslmode") {
		log.Printf("Warning: Database connection should specify SSL mode for production")
	}

	// Validate environment-specific settings
	if GlobalConfig.Environment == "production" {
		if !GlobalConfig.SecurityHeadersEnabled {
			log.Printf("Warning: Security headers should be enabled in production")
		}
		if !GlobalConfig.RateLimitEnabled {
			log.Printf("Warning: Rate limiting should be enabled in production")
		}
	}
}

// Get returns the global config instance
func Get() *Config {
	if GlobalConfig == nil {
		log.Fatal("Config not initialized. Call InitConfig() first.")
	}
	return GlobalConfig
}

// Timezone helper functions

// Now returns the current time in the configured timezone
func Now() time.Time {
	if GlobalConfig == nil || GlobalConfig.TimezoneLoc == nil {
		return time.Now().UTC()
	}
	return time.Now().In(GlobalConfig.TimezoneLoc)
}

// NowUTC returns the current time in UTC
func NowUTC() time.Time {
	return time.Now().UTC()
}

// ToLocalTime converts a UTC time to the configured timezone
func ToLocalTime(t time.Time) time.Time {
	if GlobalConfig == nil || GlobalConfig.TimezoneLoc == nil {
		return t.UTC()
	}
	return t.In(GlobalConfig.TimezoneLoc)
}

// ToUTC converts a time from the configured timezone to UTC
func ToUTC(t time.Time) time.Time {
	return t.UTC()
}

// GetTimezone returns the configured timezone string
func GetTimezone() string {
	if GlobalConfig == nil {
		return "UTC"
	}
	return GlobalConfig.Timezone
}

// GetTimezoneLocation returns the configured timezone location
func GetTimezoneLocation() *time.Location {
	if GlobalConfig == nil || GlobalConfig.TimezoneLoc == nil {
		return time.UTC
	}
	return GlobalConfig.TimezoneLoc
}
