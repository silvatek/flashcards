package main

import (
	"context"
	"fmt"
	"math/rand"
)

type DataStore interface {
	init(ctx context.Context)
	getDeck(ctx context.Context, id string) Deck
	putDeck(ctx context.Context, id string, deck Deck)
}

type TestDataStore struct {
	decks map[string]Deck
}

func randomId() string {
	return fmt.Sprintf("%04X-%04X", rand.Intn(0xFFFF), rand.Intn(0xFFFF))
}

func (store *TestDataStore) init(ctx context.Context) {
	store.decks = make(map[string]Deck)
}

func (store *TestDataStore) getDeck(ctx context.Context, id string) Deck {
	deck, ok := store.decks[id]
	if !ok {
		deck = *new(Deck)
		deck.ID = id
		deck.Title = "Test deck"
		store.decks[deck.ID] = deck
	}
	return deck
}

func (store *TestDataStore) putDeck(ctx context.Context, id string, deck Deck) {
	store.decks[id] = deck
}
