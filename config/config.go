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
	// CORS configurations
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
	// Rate limiting configurations
	RateLimitMax        int
	RateLimitWindow     time.Duration
	RateLimitEnabled    bool
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
}

var GlobalConfig *Config

// InitConfig initializes viper configuration
func InitConfig() {
	// Set config file name and paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.go-api")

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO_API")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind specific environment variables for backward compatibility
	bindEnvironmentVariables()

	// Set default values
	setDefaults()

	// Read config file (optional, fallback to env vars and defaults)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables and defaults")
		} else {
			log.Printf("Error reading config file: %v", err)
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

func bindEnvironmentVariables() {
	// Bind only required environment variables for security
	viper.BindEnv("database_url", "DATABASE_URL")
	viper.BindEnv("jwt_secret", "JWT_SECRET")
	
	// Optional: Allow GO_API_ prefixed environment variables to override config.yaml
	// Example: GO_API_APP_PORT=8080 will override app.port in config.yaml
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.port", "8000")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.read_timeout", 30*time.Second)
	viper.SetDefault("app.write_timeout", 30*time.Second)
	viper.SetDefault("app.idle_timeout", 120*time.Second)
	viper.SetDefault("app.shutdown_timeout", 10*time.Second)

	// Security defaults
	viper.SetDefault("security.jwt_expiry", 24*time.Hour)
	viper.SetDefault("security.headers_enabled", true)
	viper.SetDefault("security.trusted_proxies", "")

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", "http://localhost:3000")
	viper.SetDefault("cors.allowed_methods", "GET,POST,PUT,DELETE,OPTIONS")
	viper.SetDefault("cors.allowed_headers", "Origin,Content-Type,Accept,Authorization")

	// Rate limiting defaults
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.max", 100)
	viper.SetDefault("rate_limit.window", time.Minute)

	// Database defaults
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_lifetime", time.Hour)

	// Logging defaults
	viper.SetDefault("logging.dir", "storage/logs")
	viper.SetDefault("logging.max_size", 10*1024*1024) // 10MB
	viper.SetDefault("logging.max_age", 30)            // 30 days
	viper.SetDefault("logging.enable_daily", true)
}

func buildConfig() {
	GlobalConfig = &Config{
		// Required from env
		DatabaseURL: viper.GetString("database_url"),
		JWTSecret:   viper.GetString("jwt_secret"),

		// App configurations
		AppPort:         viper.GetString("app.port"),
		Environment:     viper.GetString("app.environment"),
		ReadTimeout:     viper.GetDuration("app.read_timeout"),
		WriteTimeout:    viper.GetDuration("app.write_timeout"),
		IdleTimeout:     viper.GetDuration("app.idle_timeout"),
		ShutdownTimeout: viper.GetDuration("app.shutdown_timeout"),

		// Security configurations
		JWTExpiry:              viper.GetDuration("security.jwt_expiry"),
		SecurityHeadersEnabled: viper.GetBool("security.headers_enabled"),
		TrustedProxies:         viper.GetString("security.trusted_proxies"),

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
	}
}

func validateConfig() {
	requiredConfigs := map[string]string{
		"DATABASE_URL": GlobalConfig.DatabaseURL,
		"JWT_SECRET":   GlobalConfig.JWTSecret,
	}

	var missingConfigs []string
	for key, value := range requiredConfigs {
		if value == "" {
			missingConfigs = append(missingConfigs, key)
		}
	}

	if len(missingConfigs) > 0 {
		log.Fatalf("Required configurations are missing: %v", missingConfigs)
	}
}

// Get returns the global config instance
func Get() *Config {
	if GlobalConfig == nil {
		log.Fatal("Config not initialized. Call InitConfig() first.")
	}
	return GlobalConfig
}

// Viper helper functions for direct access to viper instance
func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

func IsSet(key string) bool {
	return viper.IsSet(key)
}

func AllKeys() []string {
	return viper.AllKeys()
}

func WriteConfig() error {
	return viper.WriteConfig()
}

func WriteConfigAs(filename string) error {
	return viper.WriteConfigAs(filename)
}
