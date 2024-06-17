package main

import (
	"context"
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

	setupTestData(&ds)

	if ds.IsEmpty() {
		t.Error("Datastore empty after test data")
	}

	deck := ds.GetDeck(ctx, "TEST-CODE")

	if deck.ID != "TEST-CODE" {
		t.Errorf("Unexpected deck ID: %s", deck.ID)
	}
}
