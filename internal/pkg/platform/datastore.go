package platform

import (
	"context"

	"flashcards/internal/pkg/cards"
)

type DataStore interface {
	Summary() string
	Init(ctx context.Context)
	GetDeck(ctx context.Context, id string) cards.Deck
	PutDeck(ctx context.Context, id string, deck cards.Deck)
	IsEmpty() bool
	IsValidAuthor(key string) bool
}

type TestDataStore struct {
	decks map[string]cards.Deck
}

func (store *TestDataStore) Summary() string {
	return "TestDataStore"
}

func (store *TestDataStore) Init(ctx context.Context) {
	store.decks = make(map[string]cards.Deck)
}

func (store *TestDataStore) GetDeck(ctx context.Context, id string) cards.Deck {
	deck, ok := store.decks[id]
	if !ok {
		deck = *new(cards.Deck)
	}
	return deck
}

func (store *TestDataStore) PutDeck(ctx context.Context, id string, deck cards.Deck) {
	if store.decks == nil {
		store.Init(ctx)
	}
	store.decks[id] = deck
}

func (store *TestDataStore) IsEmpty() bool {
	return (store.decks == nil) || (len(store.decks) == 0)
}

func (store *TestDataStore) IsValidAuthor(key string) bool {
	return key == "guessme"
}
