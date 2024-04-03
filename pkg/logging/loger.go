package logging

import (
	"github.com/alimzhanoff/property-finder/internal/config"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

func MustLoad(config *config.Config) *logrus.Logger {
	const op = "pkg.logging.MustLoad"

	logger := logrus.New()
	level, err := logrus.ParseLevel(config.Logging.Level)
	if err != nil {
		log.Fatalf("%s: Cannot set log lever: %v", op, err)
	}
	logger.SetLevel(level)

	// Установка формата вывода
	switch strings.ToLower(config.Logging.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		log.Fatalf("%s: Unknown logging format: %v", op, config.Logging.Format)
	}

	// Установка файла для записи логов
	if config.Logging.File == "" {
		log.Fatalf("%s: Emty log file: %v", op, err)
	}

	file, err := os.OpenFile(config.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%s: Cannot create log file: %v", op, err)
	}

	logger.SetOutput(file)
	return logger
}
