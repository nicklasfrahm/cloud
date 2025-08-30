package scaleup

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"go.uber.org/zap"
)

// Event represents a single event in the timeline with a name, timestamp, duration, and extra information.
type Event struct {
	Name      string
	Timestamp time.Time
	Group     string
}

// Log logs the event to the console.
func (e Event) Log(logger *zap.Logger, fields ...zap.Field) {
	logger.Info(e.Name, append(fields, zap.String("group", e.Group))...)
}

// Timeline represents a sequence of steps in a process, each with a timestamp and duration.
type Timeline struct {
	sync.Mutex
	Steps []Event
}

func NewTimeline() *Timeline {
	return &Timeline{}
}

func (t *Timeline) Add(now time.Time, group string, name string) Event {
	event := Event{
		Name:      name,
		Timestamp: now,
		Group:     group,
	}

	t.Lock()
	defer t.Unlock()

	t.Steps = append(t.Steps, event)

	return event
}

func (t *Timeline) AddEvents(now time.Time, events ...Event) {
	t.Lock()
	defer t.Unlock()

	for _, event := range events {
		event.Timestamp = now
		t.Steps = append(t.Steps, event)
	}
}

func (t *Timeline) Print() {
	if len(t.Steps) == 0 {
		fmt.Println("No steps recorded")

		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	fmt.Fprintf(w, "GROUP\tSTEP\tRELATIVE (ms)\tΔ DURATION (ms)\tTIMESTAMP\n")

	var previous time.Time
	var start time.Time

	for i, s := range t.Steps {
		if i == 0 {
			previous = s.Timestamp
			start = s.Timestamp
		}

		duration := s.Timestamp.Sub(previous)
		relative := s.Timestamp.Sub(start)
		previous = s.Timestamp

		fmt.Fprintf(w, "%s\t%s\t%15d\t%15d\t%s\n",
			s.Group,
			s.Name,
			int(relative.Milliseconds()),
			int(duration.Milliseconds()),
			s.Timestamp.Format(time.RFC3339Nano),
		)
	}

	w.Flush()
}
