package config

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

const CONFIG_PATH = "CONFIG_PATH"

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
	Port     int           `yaml:"port"`
	Host     string        `yaml:"host"`
	Timeouts TimeoutConfig `yaml:"timeouts"`
}

type TimeoutConfig struct {
	Read    time.Duration `yaml:"read"`
	Write   time.Duration `yaml:"write"`
	Connect time.Duration `yaml:"connect"`
}

type DatabaseConfig struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	TimeOut  time.Duration `yaml:"timeout"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	DBName   string        `yaml:"dbname"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	File   string `yaml:"file"`
}

func MustLoad() *Config {
	const op = "config.MustLoad"

	if err := godotenv.Load(); err != nil {
		log.Fatalf("%s: Cannot load .env file: %v", op, err)
	}

	configPath := os.Getenv(CONFIG_PATH)
	configPath = filepath.Clean(configPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("%s: config file does not exist: %s, err: %v", op, configPath, err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("%s: Error load config file: %v", op, err)
	}

	fileWithEnvVars := os.ExpandEnv(string(data))
	// Разбор конфигурации из YAML
	var config Config
	if err := yaml.Unmarshal([]byte(fileWithEnvVars), &config); err != nil {
		log.Fatalf("%s: Error load config to struct: %v", op, err)
	}
	return &config
}
