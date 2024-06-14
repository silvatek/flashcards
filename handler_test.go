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

	deckPage(wt.w, wt.r)

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

	deckPage(wt.w, wt.r)

	wt.assertRedirectTo("/error?code=2001")
}

func TestQueryParam(t *testing.T) {
	assertQueryParam(t, "/resource?key=value", "key", "value")
	assertQueryParam(t, "/resource", "key", "")
	assertQueryParam(t, "/resource?name=nothing", "key", "")
	assertQueryParam(t, "/resource?a=1&b=2", "a", "1")
	assertQueryParam(t, "/resource?a=1&b=2", "b", "2")
}

func assertQueryParam(t *testing.T, uri string, key string, expectedValue string) {
	value := queryParam(uri, key)
	if value != expectedValue {
		t.Errorf("Unexpected query parameter value `%s` (%s %s %s)", value, uri, key, value)
	}
}

func TestLastPathElement(t *testing.T) {
	assertPathElement(t, "/abc/xyz", "xyz")
	assertPathElement(t, "/123", "123")
	assertPathElement(t, "noslash", "noslash")
	assertPathElement(t, "/abc/xyz?test=1", "xyz")
}

func assertPathElement(t *testing.T, path string, expected string) {
	result := lastPathElement(path)
	if result != expected {
		t.Errorf("lastPathElement returned [%s], expected [%s]", result, expected)
	}
}