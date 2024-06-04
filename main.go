package main

import (
	"context"
	"net/http"
	"os"
)

const defaultAddr = "127.0.0.1:8080"

type pageData struct {
	Message string
	Error   string
	Deck    Deck
	Card    Card
}

var platform Platform
var logs Logger = platform.Logger()
var dataStore DataStore = platform.DataStore()

// main starts an http server on the $PORT environment variable.
func main() {
	setupTestData(dataStore)

	addr := defaultAddr
	// $PORT environment variable is provided in the Kubernetes deployment.
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	logs.info("Server listening on port %s", addr)

	addHandlers()

	if err := http.ListenAndServe(addr, nil); err != nil {
		logs.error("Server listening error: %+v", err)
		os.Exit(-5)
	}
}

func setupTestData(store DataStore) {
	testDeck := Deck{
		ID:    "TEST-CODE",
		Title: "Test flashcard deck",
	}

	testDeck.addCard(Card{Question: "Is the sun hot?", Answer: "Yes"})
	testDeck.addCard(Card{Question: "Is the pope catholic?", Answer: "Yes"})
	testDeck.addCard(Card{Question: "How many chucks?", Answer: "42"})

	store.putDeck(context.Background(), testDeck.ID, testDeck)

	logs.debug("Test data created in %s", store.summary())
}
