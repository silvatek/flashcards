package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const DECK_COLLECTION = "Decks"
const KEYS_COLLECTION = "Keys"

type FireDataStore struct {
	Client   *firestore.Client
	Project  string
	Database string
	Err      error
}

func fireDataStore() *FireDataStore {
	store := new(FireDataStore)
	store.Project = os.Getenv("GCLOUD_PROJECT")
	store.Database = os.Getenv("FIRESTORE_DB_NAME")
	logs.info("Opening Firestore datastore %s, %s", store.Project, store.Database)
	return store
}

func (store *FireDataStore) summary() string {
	return fmt.Sprintf("FireDataStore(%s,%s)", store.Project, store.Database)
}

func createClient(ctx context.Context, projectID string, database string) (*firestore.Client, error) {
	client, err := firestore.NewClientWithDatabase(ctx, projectID, database)
	if err == nil {
		logs.info("Firestore client created: %s %s", projectID, database)
	} else {
		logs.error("Failed to create FireStore client: %v", err)
	}
	// Close client when done with "defer client.Close()"
	return client, err
}

func (store *FireDataStore) getDeck(ctx context.Context, id string) Deck {
	logs.debug1(ctx, "Fetching Firestore deck %s", id)

	var deck Deck

	doc := store.Client.Doc(DECK_COLLECTION + "/" + id)
	deckDoc, err := doc.Get(ctx)
	if err != nil {
		logs.error1(ctx, "Error fetching deck %s, %v", id, err)
	} else {
		logs.debug1(ctx, "Found game deck %s", id)

		deckDoc.DataTo(&deck)
	}

	return deck
}

func (store *FireDataStore) putDeck(ctx context.Context, id string, deck Deck) {
	logs.info1(ctx, "Writing Firestore deck %s", id)

	doc := store.Client.Doc(DECK_COLLECTION + "/" + id)
	_, err := doc.Set(ctx, deck)
	if err != nil {
		logs.error1(ctx, "Error writing deck %v", err)
	} else {
		logs.debug1(ctx, "Wrote deck document %s", id)
	}
}

func (store *FireDataStore) init() {
	store.Client, store.Err = createClient(context.Background(), store.Project, store.Database)
	logs.info("Initialised firestore")
}

func (store *FireDataStore) close() {
	store.Client.Close()
}

func (store *FireDataStore) isEmpty() bool {
	decks := store.Client.Collection(DECK_COLLECTION)
	_, err := decks.Documents(context.Background()).Next()
	return err == iterator.Done
}

func (store *FireDataStore) isValidAuthor(key string) bool {
	doc := store.Client.Doc(DECK_COLLECTION + "/" + strings.TrimSpace(key))
	keyDoc, err := doc.Get(context.Background())
	if err != nil {
		return false
	}
	return keyDoc.Data()["role"] == "author"
}
