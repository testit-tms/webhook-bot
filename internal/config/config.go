package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer  `yaml:"http_server"`
	Database    `yaml:"database"`
	TelegramBot `yaml:"telegram_bot"`
	LogLevel    string `yaml:"log_level" env-default:"Info"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int64  `yaml:"port" env-default:"5432"`
	DBName   string `yaml:"db_name" env-default:"bot"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-required:"true" env:"DB_PASSWORD"`
}

type TelegramBot struct {
	Token string `yaml:"token"  env-required:"true" env:"BOT_TOKEN"`
}

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
