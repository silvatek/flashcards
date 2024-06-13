package main

import (
	"context"
	"os"
)

const defaultAddr = "127.0.0.1:8080"

type Platform struct {
}

func (platform *Platform) runningOnGCloud() bool {
	projectId := os.Getenv("GCLOUD_PROJECT")
	return len(projectId) > 0
}

func (platform *Platform) logger() Logger {
	logs := *new(Logger)
	logs.init()
	return logs
}

func (platform *Platform) dataStore() DataStore {
	if platform.runningOnGCloud() {
		store := fireDataStore()
		store.init()
		return store
	} else {
		store := new(TestDataStore)
		store.init(context.Background())
		return store
	}
}

func (platform *Platform) listenAddress() string {
	// $PORT environment variable is provided in the Kubernetes deployment.
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	} else {
		return defaultAddr
	}
}
