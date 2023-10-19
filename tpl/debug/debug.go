// Copyright 2020 The Hugo Authors. All rights reserved.
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

// Package debug provides template functions to help debugging templates.
package debug

import (
	"sort"
	"sync"
	"time"

	"github.com/bep/logg"
	"github.com/sanity-io/litter"
	"github.com/spf13/cast"
	"github.com/yuin/goldmark/util"

	"github.com/gohugoio/hugo/deps"
)

// New returns a new instance of the debug-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	var timers map[string][]*timer
	if d.Log.Level() <= logg.LevelInfo {
		timers = make(map[string][]*timer)
	}
	ns := &Namespace{
		timers: timers,
	}

	if ns.timers == nil {
		return ns
	}

	l := d.Log.InfoCommand("timer")

	d.BuildEndListeners.Add(func() {
		type nameCountDuration struct {
			Name     string
			Count    int
			Duration time.Duration
		}

		var timersSorted []nameCountDuration

		for k, v := range timers {
			var total time.Duration
			for _, t := range v {
				// Stop any running timers.
				t.Stop()
				total += t.elapsed
			}
			timersSorted = append(timersSorted, nameCountDuration{k, len(v), total})
		}

		sort.Slice(timersSorted, func(i, j int) bool {
			// Sort it so the slowest gets printed last.
			return timersSorted[i].Duration < timersSorted[j].Duration
		})

		for _, t := range timersSorted {
			l.WithField("name", t.Name).WithField("count", t.Count).WithField("duration", t.Duration).Logf("")
		}

		ns.timers = make(map[string][]*timer)
	})

	return ns
}

// Namespace provides template functions for the "debug" namespace.
type Namespace struct {
	timersMu sync.Mutex
	timers   map[string][]*timer
}

// Dump returns a object dump of val as a string.
// Note that not every value passed to Dump will print so nicely, but
// we'll improve on that.
//
// We recommend using the "go" Chroma lexer to format the output
// nicely.
//
// Also note that the output from Dump may change from Hugo version to the next,
// so don't depend on a specific output.
func (ns *Namespace) Dump(val any) string {
	return litter.Sdump(val)
}

// VisualizeSpaces returns a string with spaces replaced by a visible string.
func (ns *Namespace) VisualizeSpaces(val any) string {
	s := cast.ToString(val)
	return string(util.VisualizeSpaces([]byte(s)))
}

func (ns *Namespace) Timer(name string) Timer {
	if ns.timers == nil {
		return nopTimer
	}
	ns.timersMu.Lock()
	defer ns.timersMu.Unlock()
	t := &timer{start: time.Now()}
	ns.timers[name] = append(ns.timers[name], t)
	return t
}

var nopTimer = nopTimerImpl{}

type nopTimerImpl struct{}

func (nopTimerImpl) Stop() string {
	return ""
}

// Timer is a timer that can be stopped.
type Timer interface {
	// Stop stops the timer and returns an empty string.
	// Stop can be called multiple times, but only the first call will stop the timer.
	// If Stop is not called, the timer will be stopped when the build ends.
	Stop() string
}

type timer struct {
	start    time.Time
	elapsed  time.Duration
	stopOnce sync.Once
}

func (t *timer) Stop() string {
	t.stopOnce.Do(func() {
		t.elapsed = time.Since(t.start)
	})
	// This is used in templates, we need to return something.
	return ""
}
