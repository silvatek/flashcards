package main

import (
	"context"
	"net/http"
	"os"

	"flashcards/internal/gcp"
	"flashcards/internal/handlers"
	"flashcards/internal/platform"
	"flashcards/internal/test"
)

// main starts an http server on the $PORT environment variable.
func main() {
	ctx := platform.NewStartupContext()

	p := getPlatform(ctx)
	logs := p.Logger()
	logs.Info(ctx, "Starting instance")

	//logEnvironment(logs, ctx)

	p.DataStore().Init(ctx)
	if p.DataStore().Summary() == "TestDataStore" {
		test.SetupTestData(ctx, p.DataStore(), logs)
	}

	router := handlers.ApplicationRouter(p)

	addr := p.ListenAddress()
	logs.Info(ctx, "Server listening on port %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logs.Error(ctx, "Server listening error: %+v", err)
		os.Exit(-5)
	}
}

func logEnvironment(logs platform.Logger, ctx context.Context) {
	for _, envvar := range os.Environ() {
		logs.Debug(ctx, "EnvVar: %s", envvar)
	}
}

func getPlatform(ctx context.Context) platform.Platform {
	if gcp.RunningOnGCloud() {
		return gcp.GcpPlatform(ctx)
	} else {
		return platform.LocalPlatform(ctx)
	}
}
