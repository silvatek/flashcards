package main

import (
	"context"
	"testing"
)

func TestDataStoreFunctionality(t *testing.T) {
	ctx := context.Background()

	var ds TestDataStore

	if ds.summary() != "TestDataStore" {
		t.Errorf("Unexpected datastore summary: %s", ds.summary())
	}

	ds.init(ctx)

	if !ds.isEmpty() {
		t.Error("Datastore was not empty at start of test")
	}

	setupTestData(&ds)

	if ds.isEmpty() {
		t.Error("Datastore empty after test data")
	}

	deck := ds.getDeck(ctx, "TEST-CODE")

	if deck.ID != "TEST-CODE" {
		t.Errorf("Unexpected deck ID: %s", deck.ID)
	}
}
