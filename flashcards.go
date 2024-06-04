package main

import (
	"fmt"
	"math/rand"
)

type Card struct {
	ID       string
	DeckID   string
	Question string
	Answer   string
}

type Deck struct {
	ID    string
	Title string
	Cards map[string]Card
}

func randomCardId() string {
	return fmt.Sprintf("%8X", rand.Intn(0xFFFFFFFF))
}

func (deck *Deck) addCard(card Card) {
	if card.ID == "" {
		card.ID = randomCardId()
	}
	if card.DeckID == "" {
		card.DeckID = deck.ID
	}
	if deck.Cards == nil {
		deck.Cards = make(map[string]Card)
	}
	deck.Cards[card.ID] = card
}

func (deck *Deck) getCard(id string) Card {
	if deck.Cards == nil {
		deck.Cards = make(map[string]Card)
	}
	return deck.Cards[id]
}
