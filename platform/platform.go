package platform

import (
	"context"
	"os"
)

type Platform interface {
	Logger() Logger
	DataStore() DataStore
	ListenAddress() string
}

var platform Platform

type TestPlatform struct {
	logs  ConsoleLogger
	store TestDataStore
}

func LocalPlatform() Platform {
	if platform == nil {
		tp := &TestPlatform{}
		tp.logs = *new(ConsoleLogger)
		tp.store = *new(TestDataStore)
		tp.store.init(context.Background())
		platform = tp
	}
	return platform
}

func (platform *TestPlatform) Logger() Logger {
	return &platform.logs
}

func (platform *TestPlatform) DataStore() DataStore {
	return &platform.store
}

func (platform *TestPlatform) ListenAddress() string {
	return "127.0.0.1:8080"
}

func TemplateDir(logs Logger) string {
	for _, path := range []string{"template", "../template"} {
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			return path
		}
	}
	logs.Error(context.Background(), "Unable to locate template files")
	return ""
}
