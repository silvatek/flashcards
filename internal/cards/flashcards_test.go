package cards

import "testing"

func TestGetCard(t *testing.T) {
	deck := Deck{
		ID: "TEST-CODE",
	}
	deck.PutCard("1", Card{Question: "Q"})

	if deck.GetCard("1").Question != "Q" {
		t.Errorf("getCard failed")
	}
}

func TestGetMissingCard(t *testing.T) {
	deck := Deck{
		ID: "TEST-CODE",
	}

	if deck.GetCard("1").ID != "" {
		t.Errorf("getCard returned a card with an ID when getting missing card")
	}
}
func TestRandomCard(t *testing.T) {
	deck := Deck{
		ID: "TEST-CODE",
	}
	deck.AddCard(Card{Question: "1"})
	deck.AddCard(Card{Question: "2"})
	deck.AddCard(Card{Question: "3"})

	counts := make(map[string]int)

	iterations := 1000

	for i := 0; i < iterations; i++ {
		card := deck.RandomCard()
		counts[card.Question] = counts[card.Question] + 1
	}

	for _, c := range []string{"1", "2", "3"} {
		if counts[c] == 0 {
			t.Errorf("No counts for %s in %d iterations", c, iterations)
		}
		if counts[c] == 1000 {
			t.Errorf("All %d counts were %s", iterations, c)
		}
	}
}

func TestRandomDeckId(t *testing.T) {
	id := RandomDeckId()

	if len(id) != 9 {
		t.Errorf("Random deck ID was not 9 characters (%d)", len(id))
	}
}
