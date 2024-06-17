package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"flashcards/platform"
	"flashcards/test"
)

func TestIndexPage(t *testing.T) {
	logs = platform.LocalPlatform().Logger()
	wt := test.NewWebTest(t)
	wt.SendGet("/")
	defer wt.ShowBodyOnFail()

	homePage(wt.Response, wt.Request)

	wt.AssertSuccess()
	wt.AssertBodyContains("h1", "Flashcards")
}

func TestShowDeckPage(t *testing.T) {
	logs = platform.LocalPlatform().Logger()
	dataStore = platform.LocalPlatform().DataStore()
	test.SetupTestData(dataStore, logs)

	wt := test.NewWebTest(t)
	wt.SendGet("/deck/TEST-CODE")
	defer wt.ShowBodyOnFail()

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertSuccess()
	wt.AssertBodyContains("title", "Flashcards - Test flashcard deck")
	wt.AssertBodyContains("h1", "Flash Card Deck")
}

func TestDeckNotFound(t *testing.T) {
	logs = platform.LocalPlatform().Logger()
	dataStore = platform.LocalPlatform().DataStore()
	test.SetupTestData(dataStore, logs)

	wt := test.NewWebTest(t)
	wt.SendGet("/deck/BAD-CODE")
	defer wt.ShowBodyOnFail()

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertRedirectTo("/error?code=2001")
}

func TestErrorPage(t *testing.T) {
	logs = platform.LocalPlatform().Logger()
	wt := test.NewWebTest(t)
	wt.SendGet("/error?code=2002")
	defer wt.ShowBodyOnFail()

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertBodyContains(".error", "Card not found")
}

func TestUnknownError(t *testing.T) {
	logs = platform.LocalPlatform().Logger()
	wt := test.NewWebTest(t)
	wt.SendGet("/error?code=9999")
	defer wt.ShowBodyOnFail()

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertBodyContains(".error", "Unknown error 9999")
}

func TestDeckRedirect(t *testing.T) {
	logs = platform.LocalPlatform().Logger()
	wt := test.NewWebTest(t)
	wt.SendGet("/decks?deck=1234")

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertRedirectTo("/deck/1234")
}

func TestDeckUrl(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Host = "localhost:8080"
	if deckUrl(r, "1234") != "http://localhost:8080/deck/1234" {
		t.Errorf("Unexpected local deck URL: %s", deckUrl(r, "1234"))
	}

	r.Host = "someserver"
	if deckUrl(r, "1234") != "https://someserver/deck/1234" {
		t.Errorf("Unexpected hosted deck URL: %s", deckUrl(r, "1234"))
	}
}

func TestNewDeck(t *testing.T) {
	logs = platform.LocalPlatform().Logger()
	dataStore = platform.LocalPlatform().DataStore()
	wt := test.NewWebTest(t)
	wt.SendPost("/newdeck", map[string]string{
		"title":  "testing",
		"author": "guessme",
	})

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	if !strings.HasPrefix(wt.RedirectTarget(), "/deck/") {
		t.Errorf("Posting new deck did not result in redirect to deck URL: %s", wt.RedirectTarget())
	}
}

func TestNewDeckBadAuthor(t *testing.T) {
	logs = platform.LocalPlatform().Logger()
	dataStore = platform.LocalPlatform().DataStore()
	wt := test.NewWebTest(t)
	wt.SendPost("/newdeck", map[string]string{
		"title":  "testing",
		"author": "badkey",
	})

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertRedirectTo("/error?code=3001")
}

func TestCardPage(t *testing.T) {
	logs = platform.LocalPlatform().Logger()
	dataStore = platform.LocalPlatform().DataStore()
	test.SetupTestData(dataStore, logs)

	deck := dataStore.GetDeck(context.Background(), "TEST-CODE")
	card := deck.RandomCard()

	wt := test.NewWebTest(t)
	wt.SendGet("/deck/" + card.DeckID + "/card/" + card.ID)
	defer wt.ShowBodyOnFail()

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertSuccess()
	wt.AssertBodyContains("title", "Flashcards - Test flashcard deck - Card")
	wt.AssertBodyContains("h1", "FlashCard")
}
