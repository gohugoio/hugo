// Copyright 2026 The Hugo Authors. All rights reserved.
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
	"cmp"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/page"
	"golang.org/x/sync/errgroup"
)

// pendingAlias is an alias that will be published once all pages are rendered.
type pendingAlias struct {
	s            *Site
	alias        string
	targetPath   string
	permalink    string
	outputFormat output.Format
	p            page.Page
}

// pendingAliases collects the aliases and the target paths of the rendered
// pages during the render step.
//
// Pages and aliases used to be written concurrently, so when an alias had the
// same target path as a page, whichever was written last won.
// Publishing the aliases after the pages lets us skip an alias that would
// overwrite a page. See issue #15645.
type pendingAliases struct {
	mu        sync.Mutex
	aliases   []pendingAlias
	pagePaths map[string]bool
}

func newPendingAliases() *pendingAliases {
	return &pendingAliases{pagePaths: make(map[string]bool)}
}

// publishPathKey normalizes a target path so that page paths (/a/index.html)
// and alias paths (a/index.html) can be compared.
func publishPathKey(targetPath string) string {
	return string(filepath.Separator) + strings.TrimLeft(filepath.Clean(targetPath), string(filepath.Separator)+"/")
}

func (a *pendingAliases) add(pa pendingAlias) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.aliases = append(a.aliases, pa)
}

func (a *pendingAliases) addPagePath(targetPath string) {
	key := publishPathKey(targetPath)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.pagePaths[key] = true
}

// publishPendingAliases publishes the aliases collected during rendering,
// skipping those that would overwrite a rendered page.
func (h *HugoSites) publishPendingAliases() error {
	pa := h.pendingAliases
	pa.mu.Lock()
	aliases := pa.aliases
	pagePaths := pa.pagePaths
	pa.aliases = nil
	pa.mu.Unlock()

	if len(aliases) == 0 {
		return nil
	}

	// Group the aliases by target path. Aliases sharing a target path are
	// written in a stable order, one after the other.
	type group struct {
		key     string
		aliases []pendingAlias
	}

	byKey := make(map[string]*group)
	var groups []*group
	for _, a := range aliases {
		key := publishPathKey(a.targetPath)
		g, found := byKey[key]
		if !found {
			g = &group{key: key}
			byKey[key] = g
			groups = append(groups, g)
		}
		g.aliases = append(g.aliases, a)
	}

	slices.SortFunc(groups, func(a, b *group) int {
		return cmp.Compare(a.key, b.key)
	})

	var eg errgroup.Group
	eg.SetLimit(max(1, h.numWorkers))

	for _, g := range groups {
		slices.SortStableFunc(g.aliases, func(a, b pendingAlias) int {
			return cmp.Compare(a.p.Path(), b.p.Path())
		})

		if pagePaths[g.key] {
			for _, a := range g.aliases {
				a.s.Log.Warnf("Alias %q of page %q was not published: a page is already published to %q.", a.alias, a.p.Path(), filepath.ToSlash(a.targetPath))
			}
			continue
		}

		eg.Go(func() error {
			for _, a := range g.aliases {
				handler := newAliasHandler(a.s.GetTemplateStore(), a.s.Log, false)
				if err := a.s.publishAliasToTargetPath(handler, a.targetPath, a.permalink, a.outputFormat, a.p); err != nil {
					return err
				}
			}
			return nil
		})
	}

	return eg.Wait()
}
