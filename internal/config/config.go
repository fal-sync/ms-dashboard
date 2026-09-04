package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	AppName         string
	HTTPPort        string
	ShutdownTimeout time.Duration
	External        ExternalConfig
}

type ExternalConfig struct {
	Service ServiceConfig
	Kafka   KafkaConfig
}

type ServiceConfig struct {
	UserSyncBaseURL string
	UserSyncPath    string
	UserSyncTimeout time.Duration
}

type KafkaConfig struct {
	Enabled          bool
	Brokers          []string
	UserCreatedTopic string
}

func Load() Config {
	return Config{
		AppName:         getEnv("APP_NAME", "clean-architecture-boilerplate"),
		HTTPPort:        getEnv("APP_PORT", "8080"),
		ShutdownTimeout: getDurationEnv("APP_SHUTDOWN_TIMEOUT", 10*time.Second),
		External: ExternalConfig{
			Service: ServiceConfig{
				UserSyncBaseURL: getEnv("EXTERNAL_USER_SYNC_BASE_URL", ""),
				UserSyncPath:    getEnv("EXTERNAL_USER_SYNC_PATH", "/internal/users/sync"),
				UserSyncTimeout: getDurationEnv("EXTERNAL_USER_SYNC_TIMEOUT", 3*time.Second),
			},
			Kafka: KafkaConfig{
				Enabled:          getBoolEnv("KAFKA_ENABLED", false),
				Brokers:          getListEnv("KAFKA_BROKERS"),
				UserCreatedTopic: getEnv("KAFKA_USER_CREATED_TOPIC", "user.created"),
			},
		},
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return duration
}

func getBoolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	switch strings.ToLower(value) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func getListEnv(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}

	rawItems := strings.Split(value, ",")
	items := make([]string, 0, len(rawItems))
	for _, item := range rawItems {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}

	return items
}
