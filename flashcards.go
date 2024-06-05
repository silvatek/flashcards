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
	return fmt.Sprintf("%08X", rand.Intn(0xFFFFFFFF))
}

func (deck *Deck) addCard(card Card) {
	if card.ID == "" {
		card.ID = randomCardId()
	}
	deck.putCard(card.ID, card)
}

func (deck *Deck) putCard(cardID string, card Card) {
	if card.DeckID == "" {
		card.DeckID = deck.ID
	}
	if deck.Cards == nil {
		deck.Cards = make(map[string]Card)
	}
	deck.Cards[cardID] = card
}
func (deck *Deck) getCard(id string) Card {
	if deck.Cards == nil {
		deck.Cards = make(map[string]Card)
	}
	return deck.Cards[id]
}

func (deck *Deck) randomCard() Card {
	cardCount := len(deck.Cards)
	logs.debug("Deck %s has %d cards", deck.ID, cardCount)
	counter := rand.Intn(cardCount)
	for _, card := range deck.Cards {
		if counter == 0 {
			return card
		}
		counter--
	}
	panic("Unable to pick a random card")
}
