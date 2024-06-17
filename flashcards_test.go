package main

import "testing"

func TestGetCard(t *testing.T) {
	deck := Deck{
		ID: "TEST-CODE",
	}
	deck.putCard("1", Card{Question: "Q"})

	if deck.getCard("1").Question != "Q" {
		t.Errorf("getCard failed")
	}
}

func TestRandomCard(t *testing.T) {
	deck := Deck{
		ID: "TEST-CODE",
	}
	deck.addCard(Card{Question: "1"})
	deck.addCard(Card{Question: "2"})
	deck.addCard(Card{Question: "3"})

	counts := make(map[string]int)

	iterations := 1000

	for i := 0; i < iterations; i++ {
		card := deck.randomCard()
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
