package main

import (
	"testing"

	"flashcards/test"
)

func TestIndexPage(t *testing.T) {
	wt := test.NewWebTest(t)
	wt.SendGet("/")
	defer wt.ShowBodyOnFail()

	homePage(wt.Response, wt.Request)

	wt.AssertSuccess()
	wt.AssertBodyContains("h1", "Flashcards")
}

func TestShowDeckPage(t *testing.T) {
	dataStore = platform.dataStore()
	setupTestData(dataStore)

	wt := test.NewWebTest(t)
	wt.SendGet("/deck/TEST-CODE")
	defer wt.ShowBodyOnFail()

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertSuccess()
	wt.AssertBodyContains("title", "Flashcards - Test flashcard deck")
	wt.AssertBodyContains("h1", "Flash Card Deck")
}

func TestDeckNotFound(t *testing.T) {
	dataStore = platform.dataStore()
	setupTestData(dataStore)

	wt := test.NewWebTest(t)
	wt.SendGet("/deck/BAD-CODE")
	defer wt.ShowBodyOnFail()

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertRedirectTo("/error?code=2001")
}

func TestErrorPage(t *testing.T) {
	wt := test.NewWebTest(t)
	wt.SendGet("/error?code=2002")
	defer wt.ShowBodyOnFail()

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertBodyContains(".error", "Card not found")
}

func TestUnknownError(t *testing.T) {
	wt := test.NewWebTest(t)
	wt.SendGet("/error?code=9999")
	defer wt.ShowBodyOnFail()

	applicationRouter().ServeHTTP(wt.Response, wt.Request)

	wt.AssertBodyContains(".error", "Unknown error 9999")
}
