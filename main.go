package main

import (
	"net/http"
	"os"

	"flashcards/gcp"
	"flashcards/platform"
	"flashcards/test"
)

var p platform.Platform
var logs platform.Logger
var dataStore platform.DataStore

// main starts an http server on the $PORT environment variable.
func main() {
	p = getPlatform()
	addr := p.ListenAddress()
	logs = p.Logger()
	dataStore = p.DataStore()
	test.SetupTestData(dataStore, logs)

	logs.Info("Server listening on port %s", addr)

	router := applicationRouter()

	if err := http.ListenAndServe(addr, router); err != nil {
		logs.Error("Server listening error: %+v", err)
		os.Exit(-5)
	}
}

func getPlatform() platform.Platform {
	if runningOnGCloud() {
		return &gcp.GooglePlatform{}
	} else {
		return &platform.TestPlatform{}
	}
}

func runningOnGCloud() bool {
	projectId := os.Getenv("GCLOUD_PROJECT")
	return len(projectId) > 0
}
