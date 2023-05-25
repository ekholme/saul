package saul

import (
	"context"
	"log"

	firestore "cloud.google.com/go/firestore"
)

// create a new firestore client
func NewFirestoreClient() *firestore.Client {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Couldn't create firestore client: %v", err)
	}

	return client
}
