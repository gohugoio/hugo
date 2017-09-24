package metrics

import (
	"io"
	"time"
)

// Noop is a no-op metrics provider.
type Noop struct{}

// MeasureSince is a no-op function.
func (n *Noop) MeasureSince(string, time.Time) {}

// Reset is a no-op function.
func (n *Noop) Reset() {}

// WriteMetrics is a no-op function.
func (n *Noop) WriteMetrics(io.Writer) {}
