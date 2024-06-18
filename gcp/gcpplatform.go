package gcp

import (
	"os"

	"flashcards/platform"
)

type GooglePlatform struct {
	logs GcpLogger2
}

func RunningOnGCloud() bool {
	return len(os.Getenv("GCLOUD_PROJECT")) > 0
}

func (platform *GooglePlatform) Logger() platform.Logger {
	//platform.logs.init()
	return &platform.logs
}

func (platform *GooglePlatform) DataStore() platform.DataStore {
	store := fireDataStore(&platform.logs)
	store.init()
	return store
}

func (platform *GooglePlatform) ListenAddress() string {
	return ":" + os.Getenv("PORT")
}
