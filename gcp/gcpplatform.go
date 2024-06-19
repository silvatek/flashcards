package gcp

import (
	"context"
	"os"

	"flashcards/platform"
)

type GooglePlatform struct {
	project string
	logs    GcpLogger
	store   FireDataStore
}

func RunningOnGCloud() bool {
	return len(os.Getenv("GCLOUD_PROJECT")) > 0
}

func GcpPlatform(ctx context.Context) *GooglePlatform {
	gcp := GooglePlatform{}
	gcp.store = *fireDataStore(&gcp.logs, ctx)
	gcp.store.init(ctx)
	gcp.project = os.Getenv("GCLOUD_PROJECT")
	return &gcp
}

func (platform *GooglePlatform) Logger() platform.Logger {
	return &platform.logs
}

func (platform *GooglePlatform) DataStore() platform.DataStore {
	return &platform.store
}

func (platform *GooglePlatform) ListenAddress() string {
	return ":" + os.Getenv("PORT")
}
