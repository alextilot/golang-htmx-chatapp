package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment         string   `env:"ENVIRONMENT,required"`
	Port                int      `env:"PORT" envDefault:"8080"`
	Hosts               []string `env:"HOSTS,envSeparator:;"`
	Debug               bool     `env:"DEBUG" envDefault:"false"`
	JwtSecretKey        string   `env:"JWT_SECRET_KEY,required"`
	JwtResfeshSecretKey string   `env:"JWT_REFREST_SECRET_KEY,required"`
	DatabasePath        string   `env:"DATABASE_PATH" envDefault:"./app_main.sqlite3"`
}

var Cfg *Config // exported global config instance

// Load reads environment variables and initializes the global config.
func init() {
	_ = godotenv.Load()

	c := &Config{}
	if err := env.Parse(&c); err != nil {
		log.Fatalf("❌ failed to load environment variables: %v", err)
	}

	log.Printf("✅ configuration loaded (env: %s)\n", c.Environment)
	Cfg = c
}
