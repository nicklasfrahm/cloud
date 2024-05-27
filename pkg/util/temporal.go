package util

import (
	"log"

	"go.temporal.io/sdk/client"
)

var (
	// temporalClient is a singleton Temporal client.
	temporalClient client.Client
)

// NewTemporalClient creates a new Temporal client if one
// does not exist or returns the existing singleton client.
func NewTemporalClient() (client.Client, error) {
	if temporalClient != nil {
		return temporalClient, nil
	}

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}

	return temporalClient, nil
}
