package main

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
