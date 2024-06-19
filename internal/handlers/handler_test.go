package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"flashcards/internal/cards"
	"flashcards/internal/platform"
	"flashcards/internal/test"
)

var p platform.Platform

func TestIndexPage(t *testing.T) {
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))
	defer wt.ShowBodyOnFail()

	wt.SendGet("/")

	wt.AssertSuccess()
	wt.AssertBodyContains("h1", "Flashcards")
}

func TestShowDeckPage(t *testing.T) {
	setupPlatform()
	test.SetupTestData(context.Background(), dataStore, logs)

	wt := test.NewWebTest(t, *ApplicationRouter(p))
	defer wt.ShowBodyOnFail()

	wt.SendGet("/deck/TEST-CODE?share=true")

	wt.AssertSuccess()
	wt.AssertBodyContains("title", "Flashcards - Test flashcard deck")
	wt.AssertBodyContains("h1", "Flash Card Deck")
}

func setupPlatform() {
	p = platform.LocalPlatform(context.Background())
	logs = p.Logger()
	dataStore = p.DataStore()
	dataStore.Init(context.Background())
}

func TestDeckNotFound(t *testing.T) {
	setupPlatform()
	test.SetupTestData(context.Background(), dataStore, logs)

	wt := test.NewWebTest(t, *ApplicationRouter(p))
	defer wt.ShowBodyOnFail()

	wt.SendGet("/deck/BAD-CODE")

	wt.AssertRedirectTo("/error?code=2001")
}

func TestErrorPage(t *testing.T) {
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))
	defer wt.ShowBodyOnFail()

	wt.SendGet("/error?code=2002")

	wt.AssertBodyContains(".error", "Card not found")
}

func TestUnknownError(t *testing.T) {
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))
	defer wt.ShowBodyOnFail()

	wt.SendGet("/error?code=9999")

	wt.AssertBodyContains(".error", "Unknown error 9999")
}

func TestDeckRedirect(t *testing.T) {
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))

	wt.SendGet("/decks?deck=1234")

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
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))

	wt.SendPost("/newdeck", map[string]string{
		"title":  "testing",
		"author": "guessme",
	})

	wt.AssertRedirectToPrefix("/deck/")
}

func TestNewDeckBadAuthor(t *testing.T) {
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))

	wt.SendPost("/newdeck", map[string]string{
		"title":  "testing",
		"author": "badkey",
	})

	wt.AssertRedirectTo("/error?code=3001")
}

func TestCardPage(t *testing.T) {
	setupPlatform()
	test.SetupTestData(context.Background(), dataStore, logs)

	deck := dataStore.GetDeck(context.Background(), "TEST-CODE")
	card := deck.RandomCard()

	wt := test.NewWebTest(t, *ApplicationRouter(p))
	defer wt.ShowBodyOnFail()

	wt.SendGet("/deck/" + card.DeckID + "/card/" + card.ID)

	wt.AssertSuccess()
	wt.AssertBodyContains("title", "Flashcards - Test flashcard deck - Card")
	wt.AssertBodyContains("h1", "FlashCard")
}

func TestCardNotFound(t *testing.T) {
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))

	// Need to initialise datastore after creating the ApplicationRouter for the test
	dataStore = &platform.TestDataStore{}
	dataStore.PutDeck(context.Background(), "123", cards.Deck{ID: "123"})

	wt.SendGet("/deck/123/card/789")

	wt.AssertRedirectTo("/error?code=2002")
}

func TestCardDeckNotFound(t *testing.T) {
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))

	dataStore = &platform.TestDataStore{}

	wt.SendGet("/deck/234/card/789")

	wt.AssertRedirectTo("/error?code=2001")
}

func TestRandomCardPage(t *testing.T) {
	setupPlatform()
	test.SetupTestData(context.Background(), dataStore, logs)

	wt := test.NewWebTest(t, *ApplicationRouter(p))

	wt.SendGet("/random?deck=TEST-CODE")

	wt.AssertRedirectToPrefix("/deck/TEST-CODE/card")
}

func TestRandomCardBadDeck(t *testing.T) {
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))
	test.SetupTestData(context.Background(), dataStore, logs)

	wt.SendGet("/random?deck=BAD-CODE")

	wt.AssertRedirectTo("/")
}

func TestEditCardForm(t *testing.T) {
	setupPlatform()
	test.SetupTestData(context.Background(), dataStore, logs)

	deck := dataStore.GetDeck(context.Background(), "TEST-CODE")
	cardID := deck.RandomCard().ID

	wt := test.NewWebTest(t, *ApplicationRouter(p))

	wt.SendGet("/editcard?deck=TEST-CODE&card=" + cardID)

	wt.AssertSuccess()
}

func TestPostEditCard(t *testing.T) {
	setupPlatform()
	test.SetupTestData(context.Background(), dataStore, logs)

	deckID := "TEST-CODE"
	deck := dataStore.GetDeck(context.Background(), deckID)
	cardID := deck.RandomCard().ID

	wt := test.NewWebTest(t, *ApplicationRouter(p))

	wt.SendPost("/editcard", map[string]string{
		"deck_id":  deckID,
		"card_id":  cardID,
		"question": "NewQ",
		"answer":   "NewA",
		"hint":     "NewH",
	})

	wt.AssertRedirectTo("/deck/" + deckID + "/card/" + cardID + "?answer=show")

	deck = dataStore.GetDeck(context.Background(), deckID)
	card := deck.Cards[cardID]

	if card.Question != "NewQ" {
		t.Errorf("Unexpected question after edit: %s", card.Question)
	}
}

func TestAddCardForm(t *testing.T) {
	setupPlatform()
	test.SetupTestData(context.Background(), dataStore, logs)

	wt := test.NewWebTest(t, *ApplicationRouter(p))

	wt.SendGet("/newcard?deck=TEST-CODE")

	wt.AssertSuccess()
}

func TestPostAddCard(t *testing.T) {
	setupPlatform()
	test.SetupTestData(context.Background(), dataStore, logs)

	deckID := "TEST-CODE"
	wt := test.NewWebTest(t, *ApplicationRouter(p))

	wt.SendPost("/newcard", map[string]string{
		"deck_id":  deckID,
		"question": "NewQ",
		"answer":   "NewA",
		"hint":     "NewH",
	})

	wt.AssertRedirectTo("/deck/" + deckID)
}

func TestQrCodes(t *testing.T) {
	setupPlatform()
	wt := test.NewWebTest(t, *ApplicationRouter(p))

	wt.SendGet("/qrcode?deck=TEST-CODE")

	wt.AssertSuccess()
}

func TestTemplateNotFound(t *testing.T) {
	lr := test.LogRecorder{}
	logs = &lr
	w := httptest.NewRecorder()

	showTemplatePage("notfound", nil, w)

	if w.Result().StatusCode != 500 {
		t.Errorf("Unexpected http response: %d", w.Result().StatusCode)
	}
	if !lr.HasEntryWithPrefix("ERROR Error parsing template") {
		t.Error("Expected log entry not found")
	}
}

func TestTemplateExecutionError(t *testing.T) {
	lr := test.LogRecorder{}
	logs = &lr
	w := httptest.NewRecorder()

	showTemplatePage("deck", "broken", w)

	if !strings.Contains(w.Body.String(), "Internal Server Error") {
		t.Errorf("Unexpected http response: %d", w.Result().StatusCode)
		t.Log(w.Body)
	}
	if !lr.HasEntryWithPrefix("ERROR template.Execute") {
		t.Error("Expected log entry not found")
	}
}
