package main

import (
	"net/http"
	"os"

	"flashcards/gcp"
	"flashcards/handlers"
	"flashcards/platform"
	"flashcards/test"
)

// main starts an http server on the $PORT environment variable.
func main() {
	p := getPlatform()
	addr := p.ListenAddress()
	logs := p.Logger()

	test.SetupTestData(p.DataStore(), logs)

	logs.Info("Server listening on port %s", addr)

	router := handlers.ApplicationRouter(p)

	if err := http.ListenAndServe(addr, router); err != nil {
		logs.Error("Server listening error: %+v", err)
		os.Exit(-5)
	}
}

func getPlatform() platform.Platform {
	if gcp.RunningOnGCloud() {
		return &gcp.GooglePlatform{}
	} else {
		return platform.LocalPlatform()
	}
}
