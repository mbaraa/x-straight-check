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
	}
}

type config struct {
	Port        string
	ApifyToken  string
	GeminiToken string
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
