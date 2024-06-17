package cards

import (
	"fmt"
	"math/rand"
)

type Card struct {
	ID       string
	DeckID   string
	Question string
	Answer   string
	Hint     string
}

type Deck struct {
	ID    string
	Title string
	Cards map[string]Card
}

func RandomDeckId() string {
	return fmt.Sprintf("%04X-%04X", rand.Intn(0xFFFF), rand.Intn(0xFFFF))
}

func RandomCardId() string {
	return fmt.Sprintf("%08X", rand.Intn(0xFFFFFFFF))
}

func (deck *Deck) AddCard(card Card) {
	if card.ID == "" {
		card.ID = RandomCardId()
	}
	deck.PutCard(card.ID, card)
}

func (deck *Deck) PutCard(cardID string, card Card) {
	if card.DeckID == "" {
		card.DeckID = deck.ID
	}
	if deck.Cards == nil {
		deck.Cards = make(map[string]Card)
	}
	deck.Cards[cardID] = card
}
func (deck *Deck) GetCard(id string) Card {
	if deck.Cards == nil {
		deck.Cards = make(map[string]Card)
	}
	return deck.Cards[id]
}

func (deck *Deck) RandomCard() Card {
	cardCount := len(deck.Cards)
	counter := rand.Intn(cardCount)
	for _, card := range deck.Cards {
		if counter == 0 {
			return card
		}
		counter--
	}
	panic("Unable to pick a random card")
}
