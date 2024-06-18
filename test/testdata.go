package test

import (
	"context"
	"flashcards/cards"
	"flashcards/platform"
)

func SetupTestData(ctx context.Context, store platform.DataStore, logs platform.Logger) {
	if !store.IsEmpty() {
		logs.Debug(ctx, "Datastore is not empty so not adding test data")
		return
	}

	testDeck := cards.Deck{
		ID:    "TEST-CODE",
		Title: "Test flashcard deck",
	}

	testDeck.AddCard(cards.Card{Question: "What Is the airspeed velocity of an unladen swallow?", Answer: "What do you mean? African or European swallow?", Hint: "Question"})
	testDeck.AddCard(cards.Card{Question: "What is the meaning of life?", Answer: "42", Hint: "Number"})
	testDeck.AddCard(cards.Card{Question: "Should I stay or should I go?", Answer: "If I stay there will be trouble", Hint: "Clash"})
	testDeck.AddCard(cards.Card{Question: "How much wood would a woodchuck chuck if a woodchuck could chuck wood?", Answer: "Much wood would be chucked"})

	testDeck.AddCard(cards.Card{Question: "Does `Markdown` work?",
		Answer: `Some features do, including...

* Bulleted lists
* **Bold**
* _Italics_

1. Numbered
2. Lists
3. Also
4. Work

But [links](http://some.bad.site/) are disabled`, Hint: "Formatting"})

	store.PutDeck(ctx, testDeck.ID, testDeck)

	logs.Debug(ctx, "Test data created in %s", store.Summary())
}
