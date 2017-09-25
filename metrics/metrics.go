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
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

// The Provider interface defines an interface for measuring metrics.
type Provider interface {
	// MeasureSince adds a measurement for key to the metric store.
	// Used with defer and time.Now().
	MeasureSince(key string, start time.Time)

	// WriteMetrics will write a summary of the metrics to w.
	WriteMetrics(w io.Writer)

	// Reset clears the metric store.
	Reset()
}

// Store provides storage for a set of metrics.
type Store struct {
	metrics map[string][]time.Duration
	mu      *sync.Mutex
}

// NewProvider returns a new instance of a metric store.
func NewProvider() Provider {
	return &Store{
		metrics: make(map[string][]time.Duration),
		mu:      &sync.Mutex{},
	}
}

// Reset clears the metrics store.
func (s *Store) Reset() {
	s.mu.Lock()
	s.metrics = make(map[string][]time.Duration)
	s.mu.Unlock()
}

// MeasureSince adds a measurement for key to the metric store.
func (s *Store) MeasureSince(key string, start time.Time) {
	s.mu.Lock()
	s.metrics[key] = append(s.metrics[key], time.Since(start))
	s.mu.Unlock()
}

// WriteMetrics writes a summary of the metrics to w.
func (s *Store) WriteMetrics(w io.Writer) {
	s.mu.Lock()

	results := make([]result, len(s.metrics))

	var i int
	for k, v := range s.metrics {
		var sum time.Duration
		var max time.Duration

		for _, d := range v {
			sum += d
			if d > max {
				max = d
			}
		}

		avg := time.Duration(int(sum) / len(v))

		results[i] = result{key: k, count: len(v), max: max, sum: sum, avg: avg}
		i++
	}

	s.mu.Unlock()

	// sort and print results
	fmt.Fprintf(w, "  %13s  %12s  %12s  %5s  %s\n", "cumulative", "average", "maximum", "", "")
	fmt.Fprintf(w, "  %13s  %12s  %12s  %5s  %s\n", "duration", "duration", "duration", "count", "template")
	fmt.Fprintf(w, "  %13s  %12s  %12s  %5s  %s\n", "----------", "--------", "--------", "-----", "--------")

	sort.Sort(bySum(results))
	for _, v := range results {
		fmt.Fprintf(w, "  %13s  %12s  %12s  %5d  %s\n", v.sum, v.avg, v.max, v.count, v.key)
	}

}

// A result represents the calculated results for a given metric.
type result struct {
	key   string
	count int
	sum   time.Duration
	max   time.Duration
	avg   time.Duration
}

type bySum []result

func (b bySum) Len() int           { return len(b) }
func (b bySum) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b bySum) Less(i, j int) bool { return b[i].sum < b[j].sum }
