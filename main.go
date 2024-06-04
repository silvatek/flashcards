package main

import (
	"context"
	"net/http"
	"os"
)

const defaultAddr = "127.0.0.1:8080"

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
	if !store.isEmpty() {
		logs.debug("Datastore is not empty so not adding test data")
		return
	}

	testDeck := Deck{
		ID:    "TEST-CODE",
		Title: "Test flashcard deck",
	}

	testDeck.addCard(Card{Question: "What Is the airspeed velocity of an unladen swallow?", Answer: "What do you mean? African or European swallow?"})
	testDeck.addCard(Card{Question: "Is the pope catholic?", Answer: "Probably"})
	testDeck.addCard(Card{Question: "What is the meaning of life?", Answer: "42"})
	testDeck.addCard(Card{Question: "Should I stay or should I go?", Answer: "If I stay there will be trouble"})
	testDeck.addCard(Card{Question: "How much wood would a woodchuck chuck if a woodchuck could chuck wood??", Answer: "Lots"})

	store.putDeck(context.Background(), testDeck.ID, testDeck)

	logs.debug("Test data created in %s", store.summary())
}
