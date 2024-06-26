package platform

import (
	"context"
	"fmt"
	"math/rand"
	"os"
)

type Platform interface {
	Logger() Logger
	DataStore() DataStore
	ListenAddress() string
}

var platform Platform

type StartupContext string

const StartupContextKey StartupContext = "StartupContext"

type StartupContextData struct {
	TraceID string
	SpanID  string
}

type TestPlatform struct {
	logs  ConsoleLogger
	store TestDataStore
}

func LocalPlatform(ctx context.Context) Platform {
	if platform == nil {
		tp := &TestPlatform{}
		tp.logs = *new(ConsoleLogger)
		tp.store = *new(TestDataStore)
		tp.store.Init(ctx)
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
	for _, path := range []string{"template", "../template", "../../template", "../../../template",
		"web", "../web", "../../web", "../../../web"} {
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			return path
		}
	}
	logs.Error(context.Background(), "Unable to locate template files")
	return ""
}

func NewStartupContext() context.Context {
	data := StartupContextData{
		TraceID: fmt.Sprintf("%08X", rand.Intn(0xFFFFFFFF)),
		SpanID:  "0",
	}
	return context.WithValue(context.Background(), StartupContextKey, data)
}
