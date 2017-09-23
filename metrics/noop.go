package metrics

import (
	"io"
	"time"
)

// Noop provides no-op storage provider.
type Noop struct{}

func (n *Noop) IsEnabled() bool                          { return false }
func (n *Noop) Reset()                                   {}
func (n *Noop) MeasureSince(key string, start time.Time) {}
func (n *Noop) WriteMetrics(w io.Writer)                 {}
