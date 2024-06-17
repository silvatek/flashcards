package gcp

import (
	"os"

	"flashcards/platform"
)

const defaultAddr = "127.0.0.1:8080"

type GooglePlatform struct {
	logs GcpLogger
}

func RunningOnGCloud() bool {
	projectId := os.Getenv("GCLOUD_PROJECT")
	return len(projectId) > 0
}

func (platform *GooglePlatform) Logger() platform.Logger {
	platform.logs.init()
	return &platform.logs
}

func (platform *GooglePlatform) DataStore() platform.DataStore {
	store := fireDataStore(&platform.logs)
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
