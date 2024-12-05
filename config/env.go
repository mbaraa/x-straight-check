package config

import (
	"os"

	"x-straight-check/log"
)

var (
	configInstance = config{}
)

func init() {
	configInstance = config{
		Port:        getEnv("PORT"),
		ApifyToken:  getEnv("APIFY_TOKEN"),
		GeminiToken: getEnv("GEMINI_TOKEN"),
		Version:     getEnv("VERSION"),
		Cache: struct {
			Address  string
			Password string
		}{
			Address:  getEnv("REDIS_ADDRESS"),
			Password: getEnv("REDIS_PASSWORD"),
		},
	}
}

type config struct {
	Port        string
	ApifyToken  string
	GeminiToken string
	Version     string
	Cache       struct {
		Address  string
		Password string
	}
}

// Env returns the thing's config values :)
func Env() config {
	return configInstance
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalln(log.ErrorLevel, "The \""+key+"\" variable is missing.")
	}
	return value
}
