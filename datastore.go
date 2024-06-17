package main

import (
	"context"
	"fmt"
	"math/rand"

	"flashcards/cards"
)

type DataStore interface {
	summary() string
	getDeck(ctx context.Context, id string) cards.Deck
	putDeck(ctx context.Context, id string, deck cards.Deck)
	isEmpty() bool
	isValidAuthor(key string) bool
}

type TestDataStore struct {
	decks map[string]cards.Deck
}

func randomDeckId() string {
	return fmt.Sprintf("%04X-%04X", rand.Intn(0xFFFF), rand.Intn(0xFFFF))
}

func (store *TestDataStore) summary() string {
	return "TestDataStore"
}

func (store *TestDataStore) init(ctx context.Context) {
	store.decks = make(map[string]cards.Deck)
}

func (store *TestDataStore) getDeck(ctx context.Context, id string) cards.Deck {
	deck, ok := store.decks[id]
	if !ok {
		deck = *new(cards.Deck)
	}
	return deck
}

func (store *TestDataStore) putDeck(ctx context.Context, id string, deck cards.Deck) {
	store.decks[id] = deck
}

func (store *TestDataStore) isEmpty() bool {
	return (store.decks == nil) || (len(store.decks) == 0)
}

func (store *TestDataStore) isValidAuthor(key string) bool {
	return key == "guessme"
}
