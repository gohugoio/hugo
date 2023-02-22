// Copyright 2019 The Hugo Authors. All rights reserved.
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

package page

import (
	"context"
	"fmt"
	"sync"

	"github.com/gohugoio/hugo/common/para"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/related"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

var (
	// Assert that Pages and PageGroup implements the PageGenealogist interface.
	_ PageGenealogist = (Pages)(nil)
	_ PageGenealogist = PageGroup{}
)

// A PageGenealogist finds related pages in a page collection. This interface is implemented
// by Pages and PageGroup, which makes it available as `{{ .RegularRelated . }}` etc.
type PageGenealogist interface {

	// Template example:
	// {{ $related := .RegularPages.Related . }}
	Related(ctx context.Context, opts any) (Pages, error)

	// Template example:
	// {{ $related := .RegularPages.RelatedIndices . "tags" "date" }}
	// Deprecated: Use Related instead.
	RelatedIndices(ctx context.Context, doc related.Document, indices ...any) (Pages, error)

	// Template example:
	// {{ $related := .RegularPages.RelatedTo ( keyVals "tags" "hugo", "rocks")  ( keyVals "date" .Date ) }}
	// Deprecated: Use Related instead.
	RelatedTo(ctx context.Context, args ...types.KeyValues) (Pages, error)
}

// Related searches all the configured indices with the search keywords from the
// supplied document.
func (p Pages) Related(ctx context.Context, optsv any) (Pages, error) {
	if len(p) == 0 {
		return nil, nil
	}

	var opts related.SearchOpts
	switch v := optsv.(type) {
	case related.Document:
		opts.Document = v
	case map[string]any:
		if err := mapstructure.WeakDecode(v, &opts); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid argument type %T", optsv)
	}

	result, err := p.search(ctx, opts)
	if err != nil {
		return nil, err
	}

	return result, nil

}

// RelatedIndices searches the given indices with the search keywords from the
// supplied document.
// Deprecated: Use Related instead.
func (p Pages) RelatedIndices(ctx context.Context, doc related.Document, indices ...any) (Pages, error) {
	indicesStr, err := cast.ToStringSliceE(indices)
	if err != nil {
		return nil, err
	}

	opts := related.SearchOpts{
		Document: doc,
		Indices:  indicesStr,
	}

	result, err := p.search(ctx, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RelatedTo searches the given indices with the corresponding values.
// Deprecated: Use Related instead.
func (p Pages) RelatedTo(ctx context.Context, args ...types.KeyValues) (Pages, error) {
	if len(p) == 0 {
		return nil, nil
	}

	opts := related.SearchOpts{
		NamedSlices: args,
	}

	return p.search(ctx, opts)
}

func (p Pages) search(ctx context.Context, opts related.SearchOpts) (Pages, error) {
	return p.withInvertedIndex(ctx, func(idx *related.InvertedIndex) ([]related.Document, error) {
		return idx.Search(ctx, opts)
	})
}

func (p Pages) withInvertedIndex(ctx context.Context, search func(idx *related.InvertedIndex) ([]related.Document, error)) (Pages, error) {
	if len(p) == 0 {
		return nil, nil
	}

	d, ok := p[0].(InternalDependencies)
	if !ok {
		return nil, fmt.Errorf("invalid type %T in related search", p[0])
	}

	cache := d.GetRelatedDocsHandler()

	searchIndex, err := cache.getOrCreateIndex(ctx, p)
	if err != nil {
		return nil, err
	}

	result, err := search(searchIndex)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		mp := make(Pages, len(result))
		for i, match := range result {
			mp[i] = match.(Page)
		}
		return mp, nil
	}

	return nil, nil
}

type cachedPostingList struct {
	p Pages

	postingList *related.InvertedIndex
}

type RelatedDocsHandler struct {
	cfg related.Config

	postingLists []*cachedPostingList
	mu           sync.RWMutex

	workers *para.Workers
}

func NewRelatedDocsHandler(cfg related.Config) *RelatedDocsHandler {
	return &RelatedDocsHandler{cfg: cfg, workers: para.New(config.GetNumWorkerMultiplier())}
}

func (s *RelatedDocsHandler) Clone() *RelatedDocsHandler {
	return NewRelatedDocsHandler(s.cfg)
}

// This assumes that a lock has been acquired.
func (s *RelatedDocsHandler) getIndex(p Pages) *related.InvertedIndex {
	for _, ci := range s.postingLists {
		if pagesEqual(p, ci.p) {
			return ci.postingList
		}
	}
	return nil
}
func (s *RelatedDocsHandler) getOrCreateIndex(ctx context.Context, p Pages) (*related.InvertedIndex, error) {
	s.mu.RLock()
	cachedIndex := s.getIndex(p)
	if cachedIndex != nil {
		s.mu.RUnlock()
		return cachedIndex, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check.
	if cachedIndex := s.getIndex(p); cachedIndex != nil {
		return cachedIndex, nil
	}

	for _, c := range s.cfg.Indices {
		if c.Type == related.TypeFragments {
			// This will trigger building the Pages' fragment map.
			g, _ := s.workers.Start(ctx)
			for _, page := range p {
				fp, ok := page.(related.FragmentProvider)
				if !ok {
					continue
				}
				g.Run(func() error {
					fp.Fragments(ctx)
					return nil
				})
			}

			if err := g.Wait(); err != nil {
				return nil, err
			}

			break
		}
	}

	searchIndex := related.NewInvertedIndex(s.cfg)

	for _, page := range p {
		if err := searchIndex.Add(ctx, page); err != nil {
			return nil, err
		}
	}

	s.postingLists = append(s.postingLists, &cachedPostingList{p: p, postingList: searchIndex})

	return searchIndex, nil
}
