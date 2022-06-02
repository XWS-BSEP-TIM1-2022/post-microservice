package config

import (
	"os"
	"time"
)

type Config struct {
	Port                  string
	PostDBHost            string
	PostDBPort            string
	PostServiceName       string
	ExpiresIn             time.Duration
	ConnectionServiceHost string
	ConnectionServicePort string
	MessageServiceHost    string
	MessageServicePort    string
}

func NewConfig() *Config {
	return &Config{
		Port:                  getEnv("POST_SERVICE_PORT", "8086"),
		PostDBHost:            getEnv("POST_DB_HOST", "dislinkt:WiYf6BvFmSpJS2Ob@xws.cjx50.mongodb.net/postsDB"),
		PostDBPort:            getEnv("POST_DB_PORT", ""),
		PostServiceName:       getEnv("POST_SERVICE_NAME", "post_service"),
		ExpiresIn:             30 * time.Minute,
		ConnectionServiceHost: getEnv("CONNECTION_SERVICE_HOST", "localhost"),
		ConnectionServicePort: getEnv("CONNECTION_SERVICE_PORT", "8087"),
		MessageServiceHost:    getEnv("MESSAGE_SERVICE_HOST", "localhost"),
		MessageServicePort:    getEnv("MESSAGE_SERVICE_PORT", "8089"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
