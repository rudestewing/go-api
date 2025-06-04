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
	LogDir         string
	LogMaxSize     int64
	LogMaxAge      int
	EnableDailyLog bool
	// Debug/Development logging
	EnableGormLog  bool
	EnableFiberLog bool
	LogLevel       string
	// Email configurations
	SMTPHost      string
	SMTPPort      int
	EmailUsername string
	EmailPassword string
	FromName      string
	FromEmail     string
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

	// Security defaults
	viper.SetDefault("security.jwt_secret", "your-super-secret-jwt-key-here")
	viper.SetDefault("security.jwt_expiry", 24*time.Hour)
	viper.SetDefault("security.headers_enabled", true)
	viper.SetDefault("security.trusted_proxies", "")

	// App defaults
	viper.SetDefault("app.port", "8000")
	viper.SetDefault("app.environment", "development")
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

	// Logging defaults
	viper.SetDefault("logging.dir", "storage/logs")
	viper.SetDefault("logging.max_size", 10*1024*1024) // 10MB
	viper.SetDefault("logging.max_age", 30)            // 30 days
	viper.SetDefault("logging.enable_daily", true)
	viper.SetDefault("logging.enable_gorm_log", false)  // Default off untuk production
	viper.SetDefault("logging.enable_fiber_log", false) // Default off untuk production
	viper.SetDefault("logging.level", "info")           // info, debug, warn, error

	// Email defaults
	viper.SetDefault("email.smtp_host", "smtp.gmail.com")
	viper.SetDefault("email.smtp_port", 587)
	viper.SetDefault("email.username", "")
	viper.SetDefault("email.password", "")
	viper.SetDefault("email.from_name", "Go API App")
	viper.SetDefault("email.from_email", "")
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
		LogDir:         viper.GetString("logging.dir"),
		LogMaxSize:     viper.GetInt64("logging.max_size"),
		LogMaxAge:      viper.GetInt("logging.max_age"),
		EnableDailyLog: viper.GetBool("logging.enable_daily"),
		EnableGormLog:  viper.GetBool("logging.enable_gorm_log"),
		EnableFiberLog: viper.GetBool("logging.enable_fiber_log"),
		LogLevel:       viper.GetString("logging.level"),

		// Email configurations
		SMTPHost:      viper.GetString("email.smtp_host"),
		SMTPPort:      viper.GetInt("email.smtp_port"),
		EmailUsername: viper.GetString("email.username"),
		EmailPassword: viper.GetString("email.password"),
		FromName:      viper.GetString("email.from_name"),
		FromEmail:     viper.GetString("email.from_email"),
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
			value == "postgres://username:password@localhost:5432/database_name" ||
			value == "your-super-secret-jwt-key-here" ||
			value == "your-super-secret-jwt-key-here-minimum-32-characters" {
			missingConfigs = append(missingConfigs, key)
		}
	}

	if len(missingConfigs) > 0 {
		log.Fatalf("Required configurations need to be updated in config.yaml: %v", missingConfigs)
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
