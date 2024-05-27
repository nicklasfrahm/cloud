package zone

import (
	"log"

	"github.com/nicklasfrahm/cloud/pkg/util"
	"go.temporal.io/sdk/worker"
)

const (
	// TaskQueue is the Task Queue for the Zone Workflow.
	TaskQueue = "zone"
)

func NewWorker() error {
	client, err := util.NewTemporalClient()
	if err != nil {
		return nil
	}

	// This worker hosts both Workflow and Activity functions.
	w := worker.New(client, TaskQueue, worker.Options{})
	w.RegisterWorkflow(WorkflowZoneUp)
	w.RegisterActivity(ActivityPrintZoneName)

	// Start listening to the Task Queue
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("failed to start worker: %s", err)
	}

	return nil
}
