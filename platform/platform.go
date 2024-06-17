package platform

import (
	"context"
	"os"
)

const defaultAddr = "127.0.0.1:8080"

type Platform interface {
	Logger() Logger
	DataStore() DataStore
	ListenAddress() string
}

type TestPlatform struct {
}

func GetPlatform() Platform {
	return &TestPlatform{}
}

func (platform *TestPlatform) Logger() Logger {
	logs := *new(ConsoleLogger)
	logs.Init()
	return &logs
}

func (platform *TestPlatform) DataStore() DataStore {
	store := new(TestDataStore)
	store.init(context.Background())
	return store
}

func (platform *TestPlatform) ListenAddress() string {
	// $PORT environment variable is provided in the Kubernetes deployment.
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	} else {
		return defaultAddr
	}
}
