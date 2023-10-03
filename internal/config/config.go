package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents the configuration for the application.
type Config struct {
	HTTPServer  `yaml:"http_server"`
	Database    `yaml:"database"`
	TelegramBot `yaml:"telegram_bot"`
	LogLevel    string `yaml:"log_level" env-default:"Info" env:"LOG_LEVEL"`
}

// HTTPServer represents the configuration for the HTTP server.
type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s" env:"TIMEOUT"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s" env:"IDLE_TIMEOUT"`
}

// Database represents the configuration for the PostgreSQL database.
type Database struct {
	Host     string `yaml:"host" env-default:"localhost" env:"DB_HOST"`
	Port     int64  `yaml:"port" env-default:"5432" env:"DB_PORT"`
	DBName   string `yaml:"db_name" env-default:"bot" env:"DB_NAME"`
	User     string `yaml:"user" env-default:"postgres" env:"DB_USER"`
	Password string `yaml:"password" env-required:"true" env:"DB_PASSWORD"`
}

// TelegramBot represents the configuration for the Telegram bot.
type TelegramBot struct {
	Token string `yaml:"token"  env-required:"true" env:"BOT_TOKEN"`
}

// MustLoad loads the configuration from the file specified in the CONFIG_PATH environment variable.
// It returns a pointer to the loaded Config struct.
// If CONFIG_PATH is not set or the file does not exist, it logs a fatal error.
// If there is an error reading the config file, it logs a fatal error.
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
