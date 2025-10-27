package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	// Application
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
	Port        string `env:"PORT" envDefault:"3000"`
	Debug       bool   `env:"DEBUG" envDefault:"false"`

	// Security
	JwtSecretKey       string `env:"JWT_SECRET_KEY,required"`
	JwtRefeshSecretKey string `env:"JWT_REFRESH_SECRET_KEY,required"`

	// Database
	DatabasePath string `env:"DATABASE_PATH" envDefault:"./internal/database/app_main.sqlite3"`

	// Cookies / Sessions
	CookieDomain string `env:"COOKIE_DOMAIN" envDefault:"localhost"`
	CookieSecure bool   `env:"COOKIE_SECURE" envDefault:"false"`

	// CORS / Allowed Hosts
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"*"`

	// Other (Optional)
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

// Cfg holds the globally available configuration.
var Cfg *Config

// Load reads environment variables and initializes the global config.
// It reads .env (if present) and parses environment variables.
func Load() *Config {
	_ = godotenv.Load() // optional; non-fatal if .env missing

	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ no .env file found, using system environment")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("❌ failed to load environment variables: %v", err)
	}

	log.Printf("✅ configuration loaded (env: %s)\n", cfg.Environment)
	Cfg = cfg

	return cfg
}

// IsProduction returns true if the app is running in production.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// Addr returns the address string used by Echo or HTTP server.
func (c *Config) Addr() string {
	return ":" + c.Port
}

// Must returns the global config or panics if not loaded.
func Must() *Config {
	if Cfg == nil {
		log.Fatal("config.Must(): configuration not loaded. Call config.Load() first.")
	}
	return Cfg
}

// GetEnv is a helper for raw lookups (used rarely).
func GetEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
