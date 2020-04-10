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
	"reflect"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/common/types"

	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/hugo/compare"
)

// The Provider interface defines an interface for measuring metrics.
type Provider interface {
	// MeasureSince adds a measurement for key to the metric store.
	// Used with defer and time.Now().
	MeasureSince(key string, start time.Time)

	// WriteMetrics will write a summary of the metrics to w.
	WriteMetrics(w io.Writer)

	// TrackValue tracks the value for diff calculations etc.
	TrackValue(key string, value interface{})

	// Reset clears the metric store.
	Reset()
}

type diff struct {
	baseline interface{}
	count    int
	simSum   int
}

var counter = 0

func (d *diff) add(v interface{}) *diff {
	if types.IsNil(d.baseline) {
		d.baseline = v
		d.count = 1
		d.simSum = 100 // If we get only one it is very cache friendly.
		return d
	}
	adder := howSimilar(v, d.baseline)
	d.simSum += adder
	d.count++

	return d
}

// Store provides storage for a set of metrics.
type Store struct {
	calculateHints bool
	metrics        map[string][]time.Duration
	mu             sync.Mutex
	diffs          map[string]*diff
	diffmu         sync.Mutex
}

// NewProvider returns a new instance of a metric store.
func NewProvider(calculateHints bool) Provider {
	return &Store{
		calculateHints: calculateHints,
		metrics:        make(map[string][]time.Duration),
		diffs:          make(map[string]*diff),
	}
}

// Reset clears the metrics store.
func (s *Store) Reset() {
	s.mu.Lock()
	s.metrics = make(map[string][]time.Duration)
	s.mu.Unlock()
	s.diffmu.Lock()
	s.diffs = make(map[string]*diff)
	s.diffmu.Unlock()
}

// TrackValue tracks the value for diff calculations etc.
func (s *Store) TrackValue(key string, value interface{}) {
	if !s.calculateHints {
		return
	}

	s.diffmu.Lock()
	var (
		d     *diff
		found bool
	)

	d, found = s.diffs[key]

	if !found {
		d = &diff{}
		s.diffs[key] = d
	}

	d.add(value)

	s.diffmu.Unlock()
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

		diff, found := s.diffs[k]

		cacheFactor := 0
		if found {
			cacheFactor = int(math.Floor(float64(diff.simSum) / float64(diff.count)))
		}

		for _, d := range v {
			sum += d
			if d > max {
				max = d
			}
		}

		avg := time.Duration(int(sum) / len(v))

		results[i] = result{key: k, count: len(v), max: max, sum: sum, avg: avg, cacheFactor: cacheFactor}
		i++
	}

	s.mu.Unlock()

	if s.calculateHints {
		fmt.Fprintf(w, "  %9s  %13s  %12s  %12s  %5s  %s\n", "cache", "cumulative", "average", "maximum", "", "")
		fmt.Fprintf(w, "  %9s  %13s  %12s  %12s  %5s  %s\n", "potential", "duration", "duration", "duration", "count", "template")
		fmt.Fprintf(w, "  %9s  %13s  %12s  %12s  %5s  %s\n", "-----", "----------", "--------", "--------", "-----", "--------")
	} else {
		fmt.Fprintf(w, "  %13s  %12s  %12s  %5s  %s\n", "cumulative", "average", "maximum", "", "")
		fmt.Fprintf(w, "  %13s  %12s  %12s  %5s  %s\n", "duration", "duration", "duration", "count", "template")
		fmt.Fprintf(w, "  %13s  %12s  %12s  %5s  %s\n", "----------", "--------", "--------", "-----", "--------")

	}

	sort.Sort(bySum(results))
	for _, v := range results {
		if s.calculateHints {
			fmt.Fprintf(w, "  %9d %13s  %12s  %12s  %5d  %s\n", v.cacheFactor, v.sum, v.avg, v.max, v.count, v.key)
		} else {
			fmt.Fprintf(w, "  %13s  %12s  %12s  %5d  %s\n", v.sum, v.avg, v.max, v.count, v.key)
		}
	}

}

// A result represents the calculated results for a given metric.
type result struct {
	key         string
	count       int
	cacheFactor int
	sum         time.Duration
	max         time.Duration
	avg         time.Duration
}

type bySum []result

func (b bySum) Len() int           { return len(b) }
func (b bySum) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b bySum) Less(i, j int) bool { return b[i].sum > b[j].sum }

// howSimilar is a naive diff implementation that returns
// a number between 0-100 indicating how similar a and b are.
func howSimilar(a, b interface{}) int {
	t1, t2 := reflect.TypeOf(a), reflect.TypeOf(b)
	if t1 != t2 {
		return 0
	}

	if t1.Comparable() && t2.Comparable() {
		if a == b {
			return 100
		}
	}

	as, ok1 := types.TypeToString(a)
	bs, ok2 := types.TypeToString(b)

	if ok1 && ok2 {
		return howSimilarStrings(as, bs)
	}

	if ok1 != ok2 {
		return 0
	}

	e1, ok1 := a.(compare.Eqer)
	e2, ok2 := b.(compare.Eqer)
	if ok1 && ok2 && e1.Eq(e2) {
		return 100
	}

	pe1, pok1 := a.(compare.ProbablyEqer)
	pe2, pok2 := b.(compare.ProbablyEqer)
	if pok1 && pok2 && pe1.ProbablyEq(pe2) {
		return 90
	}

	h1, h2 := helpers.HashString(a), helpers.HashString(b)
	if h1 == h2 {
		return 100
	}
	return 0

}

// howSimilar is a naive diff implementation that returns
// a number between 0-100 indicating how similar a and b are.
// 100 is when all words in a also exists in b.
func howSimilarStrings(a, b string) int {
	if a == b {
		return 100
	}

	// Give some weight to the word positions.
	const partitionSize = 4

	af, bf := strings.Fields(a), strings.Fields(b)
	if len(bf) > len(af) {
		af, bf = bf, af
	}

	m1 := make(map[string]bool)
	for i, x := range bf {
		partition := partition(i, partitionSize)
		key := x + "/" + strconv.Itoa(partition)
		m1[key] = true
	}

	common := 0
	for i, x := range af {
		partition := partition(i, partitionSize)
		key := x + "/" + strconv.Itoa(partition)
		if m1[key] {
			common++
		}
	}

	return int(math.Floor((float64(common) / float64(len(af)) * 100)))
}

func partition(d, scale int) int {
	return int(math.Floor((float64(d) / float64(scale)))) * scale
}
