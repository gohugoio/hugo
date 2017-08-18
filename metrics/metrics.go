// Copyright 2017 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package metrics provides simple metrics tracking features.
package metrics

import (
	"sort"
	"sync"
	"time"

	jww "github.com/spf13/jwalterweatherman"
)

// The MetricsProvider interface provides metric measuring features.
type MetricsProvider interface {
	MeasureSince(key string, start time.Time)
	LogFeedback(w *jww.Feedback)
	Reset()
}

// Store provides storage for a set of metrics.
type Store struct {
	metrics map[string][]time.Duration
	mtx     *sync.Mutex
}

// NewStore returns a new instance of a metric store.
func NewStore() *Store {
	return &Store{
		metrics: make(map[string][]time.Duration),
		mtx:     &sync.Mutex{},
	}
}

// Reset clears the metrics store.
func (s *Store) Reset() {
	s.mtx.Lock()
	s.metrics = make(map[string][]time.Duration)
	s.mtx.Unlock()
}

// MeasureSince adds a measurement for key to the metric store.
func (s *Store) MeasureSince(key string, start time.Time) {
	s.mtx.Lock()
	s.metrics[key] = append(s.metrics[key], time.Since(start))
	s.mtx.Unlock()
}

// LogFeedback prints metrics in a pretty format to a jwalterweatherman.Feedback.
func (s *Store) LogFeedback(w *jww.Feedback) {
	w.Printf("  %13s  %12s  %12s  %5s  %s\n", "cumulative", "average", "maximum", "", "")
	w.Printf("  %13s  %12s  %12s  %5s  %s\n", "duration", "duration", "duration", "count", "template")
	w.Printf("  %13s  %12s  %12s  %5s  %s\n", "----------", "--------", "--------", "-----", "--------")

	s.mtx.Lock()

	// sort keys
	var tkeys []string
	for k := range s.metrics {
		tkeys = append(tkeys, k)
	}
	sort.Strings(tkeys)

	for _, k := range tkeys {
		var sum time.Duration
		var max time.Duration

		for _, d := range s.metrics[k] {
			sum += d
			if d > max {
				max = d
			}
		}

		w.Printf("  %13s  %12s  %12s  %5d  %s\n",
			sum, time.Duration(int(sum)/len(s.metrics[k])), max, len(s.metrics[k]), k)
	}

	s.mtx.Unlock()
}
