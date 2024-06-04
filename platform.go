package main

import (
	"context"
	"os"
)

type Platform struct {
}

func (platform *Platform) runningOnGCloud() bool {
	projectId := os.Getenv("GCLOUD_PROJECT")
	return len(projectId) > 0
}

func (platform *Platform) Logger() Logger {
	logs := *new(Logger)
	logs.init()
	return logs
}

func (platform *Platform) DataStore() DataStore {
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
