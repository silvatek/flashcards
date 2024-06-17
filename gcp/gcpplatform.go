package gcp

import (
	"os"

	"flashcards/platform"
)

const defaultAddr = "127.0.0.1:8080"

type GooglePlatform struct {
	logs platform.Logger
}

func (platform *GooglePlatform) Logger() platform.Logger {
	platform.logs.Init()
	return platform.logs
}

func (platform *GooglePlatform) DataStore() platform.DataStore {
	store := fireDataStore(platform.logs)
	store.init()
	return store
}

func (platform *GooglePlatform) ListenAddress() string {
	// $PORT environment variable is provided in the Kubernetes deployment.
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	} else {
		return defaultAddr
	}
}
