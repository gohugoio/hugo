// Copyright 2018 The Hugo Authors. All rights reserved.
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

package langs

import (
	"sync"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/loggers"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func TestCollator(t *testing.T) {
	c := qt.New(t)

	var wg sync.WaitGroup

	coll := &Collator{c: collate.New(language.English, collate.Loose)}

	for range 10 {
		wg.Add(1)
		go func() {
			coll.Lock()
			defer coll.Unlock()
			defer wg.Done()
			for range 10 {
				k := coll.CompareStrings("abc", "def")
				c.Assert(k, qt.Equals, -1)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkCollator(b *testing.B) {
	s := []string{"foo", "bar", "éntre", "baz", "qux", "quux", "corge", "grault", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"}

	doWork := func(coll *Collator) {
		for i := range s {
			for j := i + 1; j < len(s); j++ {
				_ = coll.CompareStrings(s[i], s[j])
			}
		}
	}

	b.Run("Single", func(b *testing.B) {
		coll := &Collator{c: collate.New(language.English, collate.Loose)}
		for b.Loop() {
			doWork(coll)
		}
	})

	b.Run("Para", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			coll := &Collator{c: collate.New(language.English, collate.Loose)}

			for pb.Next() {
				coll.Lock()
				doWork(coll)
				coll.Unlock()
			}
		})
	})
}

// TestLanguageLegacyFieldFallbacks verifies that Locale(), Direction(), and
// Label() fall back to legacy LanguageConfig fields when the canonical fields
// are not set. This matters for programmatic construction that bypasses the
// allconfig migration (which normally copies legacy→canonical and clears them).
func TestLanguageLegacyFieldFallbacks(t *testing.T) {
	c := qt.New(t)

	l, err := NewLanguage("en", "en", "UTC", LanguageConfig{
		LanguageCode:      "en-US",
		LanguageName:      "English",
		LanguageDirection: "ltr",
	}, loggers.NewDefault())
	c.Assert(err, qt.IsNil)
	c.Assert(l.Locale(), qt.Equals, "en-US")
	c.Assert(l.Label(), qt.Equals, "English")
	c.Assert(l.Direction(), qt.Equals, "ltr")

	// Deprecated methods must not panic when logger is nil (no logger provided
	// at construction). They delegate through Logger() which has a nil guard.
	lNoLogger, err := NewLanguage("en", "en", "UTC", LanguageConfig{}, nil)
	c.Assert(err, qt.IsNil)
	_ = lNoLogger.LanguageCode()
	_ = lNoLogger.LanguageName()
	_ = lNoLogger.LanguageDirection()
}
