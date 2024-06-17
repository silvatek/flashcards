package main

import (
	"testing"
)

func TestIndexPage(t *testing.T) {
	wt := webTest(t)
	wt.sendGet("/")
	defer wt.showBodyOnFail()

	homePage(wt.w, wt.r)

	wt.assertSuccess()
	wt.assertBodyContains("h1", "Flashcards")
}

func TestShowDeckPage(t *testing.T) {
	dataStore = platform.dataStore()
	setupTestData(dataStore)

	wt := webTest(t)
	wt.sendGet("/deck/TEST-CODE")
	defer wt.showBodyOnFail()

	applicationRouter().ServeHTTP(wt.w, wt.r)

	wt.assertSuccess()
	wt.assertBodyContains("title", "Flashcards - Test flashcard deck")
	wt.assertBodyContains("h1", "Flash Card Deck")
}

func TestDeckNotFound(t *testing.T) {
	dataStore = platform.dataStore()
	setupTestData(dataStore)

	wt := webTest(t)
	wt.sendGet("/deck/BAD-CODE")
	defer wt.showBodyOnFail()

	applicationRouter().ServeHTTP(wt.w, wt.r)

	wt.assertRedirectTo("/error?code=2001")
}

func TestErrorPage(t *testing.T) {
	wt := webTest(t)
	wt.sendGet("/error?code=2002")
	defer wt.showBodyOnFail()

	applicationRouter().ServeHTTP(wt.w, wt.r)

	wt.assertBodyContains(".error", "Card not found")
}

func TestUnknownError(t *testing.T) {
	wt := webTest(t)
	wt.sendGet("/error?code=9999")
	defer wt.showBodyOnFail()

	applicationRouter().ServeHTTP(wt.w, wt.r)

	wt.assertBodyContains(".error", "Unknown error 9999")
}
