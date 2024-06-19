package main

import (
	"context"
	"net/http"
	"os"

	"flashcards/gcp"
	"flashcards/handlers"
	"flashcards/platform"
	"flashcards/test"
)

// main starts an http server on the $PORT environment variable.
func main() {
	ctx := platform.NewStartupContext()

	p := getPlatform(ctx)
	addr := p.ListenAddress()
	logs := p.Logger()

	logs.Info(ctx, "Starting instance")

	//logEnvironment(logs, ctx)

	test.SetupTestData(ctx, p.DataStore(), logs)

	router := handlers.ApplicationRouter(p)

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
		return platform.LocalPlatform()
	}
}
