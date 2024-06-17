package platform

import (
	"context"
	"flashcards/cards"
	"testing"
)

func TestDataStoreFunctionality(t *testing.T) {
	ctx := context.Background()

	var ds TestDataStore

	if ds.Summary() != "TestDataStore" {
		t.Errorf("Unexpected datastore summary: %s", ds.Summary())
	}

	ds.init(ctx)

	if !ds.IsEmpty() {
		t.Error("Datastore was not empty at start of test")
	}

	deck := cards.Deck{
		ID:    "TEST-CODE",
		Title: "Testing",
	}

	ds.PutDeck(ctx, deck.ID, deck)

	if ds.IsEmpty() {
		t.Error("Datastore empty after test data")
	}

	deck2 := ds.GetDeck(ctx, "TEST-CODE")

	if deck2.ID != "TEST-CODE" {
		t.Errorf("Unexpected deck ID: %s", deck2.ID)
	}
}
