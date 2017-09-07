// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"sync"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/related"
	"github.com/spf13/cast"
)

var (
	// Assert that Pages and PageGroup implements the PageGenealogist interface.
	_ PageGenealogist = (Pages)(nil)
	_ PageGenealogist = PageGroup{}
)

// A PageGenealogist finds related pages in a page collection. This interface is implemented
// by Pages and PageGroup, which makes it available as `{{ .RegularPages.Related . }}` etc.
type PageGenealogist interface {

	// Template example:
	// {{ $related := .RegularPages.Related . }}
	Related(doc related.Document) (Pages, error)

	// Template example:
	// {{ $related := .RegularPages.RelatedIndices . "tags" "date" }}
	RelatedIndices(doc related.Document, indices ...interface{}) (Pages, error)

	// Template example:
	// {{ $related := .RegularPages.RelatedTo ( keyVals "tags" "hugo", "rocks")  ( keyVals "date" .Date ) }}
	RelatedTo(args ...types.KeyValues) (Pages, error)
}

// Related searches all the configured indices with the search keywords from the
// supplied document.
func (p Pages) Related(doc related.Document) (Pages, error) {
	page, err := unwrapPage(doc)
	if err != nil {
		return nil, err
	}

	result, err := p.searchDoc(page)
	if err != nil {
		return nil, err
	}

	return result.removeFirstIfFound(page), nil
}

// RelatedIndices searches the given indices with the search keywords from the
// supplied document.
func (p Pages) RelatedIndices(doc related.Document, indices ...interface{}) (Pages, error) {
	page, err := unwrapPage(doc)
	if err != nil {
		return nil, err
	}

	indicesStr, err := cast.ToStringSliceE(indices)
	if err != nil {
		return nil, err
	}

	result, err := p.searchDoc(page, indicesStr...)
	if err != nil {
		return nil, err
	}

	return result.removeFirstIfFound(page), nil

}

// RelatedTo searches the given indices with the corresponding values.
func (p Pages) RelatedTo(args ...types.KeyValues) (Pages, error) {
	if len(p) == 0 {
		return nil, nil
	}

	return p.search(args...)

}

func (p Pages) search(args ...types.KeyValues) (Pages, error) {
	return p.withInvertedIndex(func(idx *related.InvertedIndex) ([]related.Document, error) {
		return idx.SearchKeyValues(args...)
	})

}

func (p Pages) searchDoc(doc related.Document, indices ...string) (Pages, error) {
	return p.withInvertedIndex(func(idx *related.InvertedIndex) ([]related.Document, error) {
		return idx.SearchDoc(doc, indices...)
	})
}

func (p Pages) withInvertedIndex(search func(idx *related.InvertedIndex) ([]related.Document, error)) (Pages, error) {
	if len(p) == 0 {
		return nil, nil
	}

	cache := p[0].s.relatedDocsHandler

	searchIndex, err := cache.getOrCreateIndex(p)
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
			mp[i] = match.(*Page)
		}
		return mp, nil
	}

	return nil, nil
}

type cachedPostingList struct {
	p Pages

	postingList *related.InvertedIndex
}

type relatedDocsHandler struct {
	// This is configured in site or langugage config.
	cfg related.Config

	postingLists []*cachedPostingList
	mu           sync.RWMutex
}

func newSearchIndexHandler(cfg related.Config) *relatedDocsHandler {
	return &relatedDocsHandler{cfg: cfg}
}

// This assumes that a lock has been acquired.
func (s *relatedDocsHandler) getIndex(p Pages) *related.InvertedIndex {
	for _, ci := range s.postingLists {
		if fastEqualPages(p, ci.p) {
			return ci.postingList
		}
	}
	return nil
}

func (s *relatedDocsHandler) getOrCreateIndex(p Pages) (*related.InvertedIndex, error) {
	s.mu.RLock()
	cachedIndex := s.getIndex(p)
	if cachedIndex != nil {
		s.mu.RUnlock()
		return cachedIndex, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if cachedIndex := s.getIndex(p); cachedIndex != nil {
		return cachedIndex, nil
	}

	searchIndex := related.NewInvertedIndex(s.cfg)

	for _, page := range p {
		if err := searchIndex.Add(page); err != nil {
			return nil, err
		}
	}

	s.postingLists = append(s.postingLists, &cachedPostingList{p: p, postingList: searchIndex})

	return searchIndex, nil
}
