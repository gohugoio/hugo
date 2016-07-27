// Copyright 2016-present The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"time"

	"github.com/fsnotify/fsnotify"

	jww "github.com/spf13/jwalterweatherman"
)

// HugoSites represents the sites to build. Each site represents a language.
type HugoSites []*Site

// Reset resets the sites, making it ready for a full rebuild.
// TODO(bep) multilingo
func (h HugoSites) Reset() {
	for i, s := range h {
		h[i] = s.Reset()
	}
}

// Build builds all sites.
func (h HugoSites) Build(watching, printStats bool) error {
	t0 := time.Now()

	for _, site := range h {
		t1 := time.Now()

		site.RunMode.Watching = watching

		if err := site.Build(); err != nil {
			return err
		}
		if printStats {
			site.Stats(t1)
		}
	}

	if printStats {
		jww.FEEDBACK.Printf("total in %v ms\n", int(1000*time.Since(t0).Seconds()))
	}

	return nil

}

// Rebuild rebuilds all sites.
func (h HugoSites) Rebuild(events []fsnotify.Event, printStats bool) error {
	t0 := time.Now()

	for _, site := range h {
		t1 := time.Now()

		if err := site.ReBuild(events); err != nil {
			return err
		}

		if printStats {
			site.Stats(t1)
		}
	}

	if printStats {
		jww.FEEDBACK.Printf("total in %v ms\n", int(1000*time.Since(t0).Seconds()))
	}

	return nil

}
