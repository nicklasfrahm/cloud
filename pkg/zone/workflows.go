package zone

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// WorkflowZoneUp creates or updates an availability zone.
func WorkflowZoneUp(ctx workflow.Context, zone Zone) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	return workflow.ExecuteActivity(ctx, ActivityPrintZoneName, zone).Get(ctx, nil)
}
