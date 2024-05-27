package zone

import (
	"context"
	"encoding/json"

	"go.temporal.io/sdk/activity"
)

// ActivityPrintZoneName prints the name of the zone.
func ActivityPrintZoneName(ctx context.Context, zone *Zone) error {
	logger := activity.GetLogger(ctx)

	jsonZone, err := json.Marshal(zone)
	if err != nil {
		logger.Error("failed to marshal zone", "err", err)
		return err
	}

	msg := "Zone name: " + zone.Name
	logger.Info(msg, "zone", string(jsonZone))

	return nil
}
