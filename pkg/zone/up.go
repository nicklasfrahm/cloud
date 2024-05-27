package zone

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/nicklasfrahm/cloud/pkg/util"
	"go.temporal.io/server/temporal"
)

// Up creates or updates an availability zone.
func Up(zone Zone) error {
	// Check if a temporal service is running for this domain.

	// Perform a DNS query to check if the zone is already up.
	if _, err := net.LookupIP("temporal.example.com"); err != nil {
		log.Println("failed to resolve temporal.example.com, starting Temporal Server...")

		go func() {
			server, err := temporal.NewServer(
				temporal.ForServices(temporal.DefaultServices),
				temporal.InterruptOn(temporal.InterruptCh()),
			)
			if err != nil {
				log.Fatal(err)
			}

			if err := server.Start(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	// Start the worker.
	go func() {
		if err := NewWorker(); err != nil {
			log.Fatal(err)
		}
	}()

	client, err := util.NewTemporalClient()
	if err != nil {
		return err
	}

	// Start the workflow.
	run, err := client.ExecuteWorkflow(context.Background(), WHAT?, WorkflowZoneUp, zone)
	if err != nil {
		return err
	}

	// Wait for the workflow to complete.
	run.Get(context.Background(), nil)

	return nil
}
