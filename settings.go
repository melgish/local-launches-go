package main

import (
	"log"
	"os"
	"time"
)

const (
	DefHttpPort        = "8080"
	DefLogRequests     = false
	DefRefreshInterval = 4 * time.Hour

	EnvHttpPort        = "HTTP_PORT"
	EnvLogRequests     = "LOG_REQUESTS"
	EnvRefreshInterval = "REFRESH"
)

type Settings struct {
	HttpPort        string
	LogRequests     bool
	RefreshInterval time.Duration
}

func loadSettings() Settings {
	settings := Settings{
		HttpPort:        DefHttpPort,
		LogRequests:     DefLogRequests,
		RefreshInterval: DefRefreshInterval,
	}

	str := os.Getenv(EnvHttpPort)
	if str != "" {
		settings.HttpPort = str
	}

	str = os.Getenv(EnvLogRequests)
	if str != "" {
		// Anything other than "true" will be considered false.
		settings.LogRequests = str == "true"
	}

	str = os.Getenv(EnvRefreshInterval)
	if str != "" {
		duration, err := time.ParseDuration(str)
		if err != nil {
			settings.RefreshInterval = duration
		} else {
			// Log the error and fallback to the default refresh interval
			// if the parsing fails.
			log.Println("Invalid refresh interval. Using default.")
		}
	}

	return settings
}
