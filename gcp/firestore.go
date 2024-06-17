package gcp

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"flashcards/cards"
	"flashcards/platform"
)

const DECK_COLLECTION = "Decks"
const KEYS_COLLECTION = "Keys"

type FireDataStore struct {
	Client   *firestore.Client
	Project  string
	Database string
	Err      error
	logs     platform.Logger
}

func fireDataStore(logs platform.Logger) *FireDataStore {
	store := new(FireDataStore)
	store.Project = os.Getenv("GCLOUD_PROJECT")
	store.Database = os.Getenv("FIRESTORE_DB_NAME")
	store.logs = logs
	logs.Info("Opening Firestore datastore %s, %s", store.Project, store.Database)
	return store
}

func (store *FireDataStore) Summary() string {
	return fmt.Sprintf("FireDataStore(%s,%s)", store.Project, store.Database)
}

func (store *FireDataStore) createClient(ctx context.Context, projectID string, database string) (*firestore.Client, error) {
	client, err := firestore.NewClientWithDatabase(ctx, projectID, database)
	if err == nil {
		store.logs.Info("Firestore client created: %s %s", projectID, database)
	} else {
		store.logs.Error("Failed to create FireStore client: %v", err)
	}
	// Close client when done with "defer client.Close()"
	return client, err
}

func (store *FireDataStore) GetDeck(ctx context.Context, id string) cards.Deck {
	store.logs.Debug("Fetching Firestore deck %s", id)

	var deck cards.Deck

	doc := store.Client.Doc(DECK_COLLECTION + "/" + id)
	deckDoc, err := doc.Get(ctx)
	if err != nil {
		store.logs.Error("Error fetching deck %s, %v", id, err)
	} else {
		store.logs.Debug("Found game deck %s", id)

		deckDoc.DataTo(&deck)
	}

	return deck
}

func (store *FireDataStore) PutDeck(ctx context.Context, id string, deck cards.Deck) {
	store.logs.Info("Writing Firestore deck %s", id)

	doc := store.Client.Doc(DECK_COLLECTION + "/" + id)
	_, err := doc.Set(ctx, deck)
	if err != nil {
		store.logs.Error("Error writing deck %v", err)
	} else {
		store.logs.Debug("Wrote deck document %s", id)
	}
}

func (store *FireDataStore) init() {
	store.Client, store.Err = store.createClient(context.Background(), store.Project, store.Database)
	store.logs.Info("Initialised firestore")
}

func (store *FireDataStore) close() {
	store.Client.Close()
}

func (store *FireDataStore) IsEmpty() bool {
	decks := store.Client.Collection(DECK_COLLECTION)
	_, err := decks.Documents(context.Background()).Next()
	return err == iterator.Done
}

func (store *FireDataStore) IsValidAuthor(key string) bool {
	doc := store.Client.Doc(KEYS_COLLECTION + "/" + strings.TrimSpace(key))
	keyDoc, err := doc.Get(context.Background())
	if err != nil {
		store.logs.Info("Author key not found")
		return false
	}
	if keyDoc.Data()["role"] != "author" {
		store.logs.Info("Key does not have author role")
		return false
	}
	return true
}