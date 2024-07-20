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

// Package related holds code to help finding related content.
package related

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	xmaps "golang.org/x/exp/maps"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/compare"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/types"
	"github.com/mitchellh/mapstructure"
)

const (
	TypeBasic     = "basic"
	TypeFragments = "fragments"
)

var validTypes = map[string]bool{
	TypeBasic:     true,
	TypeFragments: true,
}

var (
	_        Keyword = (*StringKeyword)(nil)
	zeroDate         = time.Time{}

	// DefaultConfig is the default related config.
	DefaultConfig = Config{
		Threshold: 80,
		Indices: IndicesConfig{
			IndexConfig{Name: "keywords", Weight: 100, Type: TypeBasic},
			IndexConfig{Name: "date", Weight: 10, Type: TypeBasic},
		},
	}
)

// Config is the top level configuration element used to configure how to retrieve
// related content in Hugo.
type Config struct {
	// Only include matches >= threshold, a normalized rank between 0 and 100.
	Threshold int

	// To get stable "See also" sections we, by default, exclude newer related pages.
	IncludeNewer bool

	// Will lower case all string values and queries to the indices.
	// May get better results, but at a slight performance cost.
	ToLower bool

	Indices IndicesConfig
}

// Add adds a given index.
func (c *Config) Add(index IndexConfig) {
	if c.ToLower {
		index.ToLower = true
	}
	c.Indices = append(c.Indices, index)
}

func (c *Config) HasType(s string) bool {
	for _, i := range c.Indices {
		if i.Type == s {
			return true
		}
	}
	return false
}

// IndicesConfig holds a set of index configurations.
type IndicesConfig []IndexConfig

// IndexConfig configures an index.
type IndexConfig struct {
	// The index name. This directly maps to a field or Param name.
	Name string

	// The index type.
	Type string

	// Enable to apply a type specific filter to the results.
	// This is currently only used for the "fragments" type.
	ApplyFilter bool

	// Contextual pattern used to convert the Param value into a string.
	// Currently only used for dates. Can be used to, say, bump posts in the same
	// time frame when searching for related documents.
	// For dates it follows Go's time.Format patterns, i.e.
	// "2006" for YYYY and "200601" for YYYYMM.
	Pattern string

	// This field's weight when doing multi-index searches. Higher is "better".
	Weight int

	// A percentage (0-100) used to remove common keywords from the index.
	// As an example, setting this to 50 will remove all keywords that are
	// used in more than 50% of the documents in the index.
	CardinalityThreshold int

	// Will lower case all string values in and queries tothis index.
	// May get better accurate results, but at a slight performance cost.
	ToLower bool
}

// Document is the interface an indexable document in Hugo must fulfill.
type Document interface {
	// RelatedKeywords returns a list of keywords for the given index config.
	RelatedKeywords(cfg IndexConfig) ([]Keyword, error)

	// When this document was or will be published.
	PublishDate() time.Time

	// Name is used as an tiebreaker if both Weight and PublishDate are
	// the same.
	Name() string
}

// FragmentProvider is an optional interface that can be implemented by a Document.
type FragmentProvider interface {
	Fragments(context.Context) *tableofcontents.Fragments

	// For internal use.
	ApplyFilterToHeadings(context.Context, func(*tableofcontents.Heading) bool) Document
}

// InvertedIndex holds an inverted index, also sometimes named posting list, which
// lists, for every possible search term, the documents that contain that term.
type InvertedIndex struct {
	cfg   Config
	index map[string]map[Keyword][]Document
	// Counts the number of documents added to each index.
	indexDocCount map[string]int

	minWeight int
	maxWeight int

	// No modifications after this is set.
	finalized bool
}

func (idx *InvertedIndex) getIndexCfg(name string) (IndexConfig, bool) {
	for _, conf := range idx.cfg.Indices {
		if conf.Name == name {
			return conf, true
		}
	}

	return IndexConfig{}, false
}

// NewInvertedIndex creates a new InvertedIndex.
// Documents to index must be added in Add.
func NewInvertedIndex(cfg Config) *InvertedIndex {
	idx := &InvertedIndex{index: make(map[string]map[Keyword][]Document), indexDocCount: make(map[string]int), cfg: cfg}
	for _, conf := range cfg.Indices {
		idx.index[conf.Name] = make(map[Keyword][]Document)
		if conf.Weight < idx.minWeight {
			// By default, the weight scale starts at 0, but we allow
			// negative weights.
			idx.minWeight = conf.Weight
		}
		if conf.Weight > idx.maxWeight {
			idx.maxWeight = conf.Weight
		}
	}
	return idx
}

// Add documents to the inverted index.
// The value must support == and !=.
func (idx *InvertedIndex) Add(ctx context.Context, docs ...Document) error {
	if idx.finalized {
		panic("index is finalized")
	}
	var err error
	for _, config := range idx.cfg.Indices {
		if config.Weight == 0 {
			// Disabled
			continue
		}
		setm := idx.index[config.Name]

		for _, doc := range docs {
			var added bool
			var words []Keyword
			words, err = doc.RelatedKeywords(config)
			if err != nil {
				continue
			}

			for _, keyword := range words {
				added = true
				setm[keyword] = append(setm[keyword], doc)
			}

			if config.Type == TypeFragments {
				if fp, ok := doc.(FragmentProvider); ok {
					for _, fragment := range fp.Fragments(ctx).Identifiers {
						added = true
						setm[FragmentKeyword(fragment)] = append(setm[FragmentKeyword(fragment)], doc)
					}
				}
			}

			if added {
				idx.indexDocCount[config.Name]++
			}
		}
	}

	return err
}

func (idx *InvertedIndex) Finalize(ctx context.Context) error {
	if idx.finalized {
		return nil
	}

	for _, config := range idx.cfg.Indices {
		if config.CardinalityThreshold == 0 {
			continue
		}
		setm := idx.index[config.Name]
		if idx.indexDocCount[config.Name] == 0 {
			continue
		}

		// Remove high cardinality terms.
		numDocs := idx.indexDocCount[config.Name]
		for k, v := range setm {
			percentageWithKeyword := int(math.Ceil(float64(len(v)) / float64(numDocs) * 100))
			if percentageWithKeyword > config.CardinalityThreshold {
				delete(setm, k)
			}
		}

	}

	idx.finalized = true

	return nil
}

// queryElement holds the index name and keywords that can be used to compose a
// search for related content.
type queryElement struct {
	Index    string
	Keywords []Keyword
}

func newQueryElement(index string, keywords ...Keyword) queryElement {
	return queryElement{Index: index, Keywords: keywords}
}

type ranks []*rank

type rank struct {
	Doc     Document
	Weight  int
	Matches int
}

func (r *rank) addWeight(w int) {
	r.Weight += w
	r.Matches++
}

var rankPool = sync.Pool{
	New: func() interface{} {
		return &rank{}
	},
}

func getRank(doc Document, weight int) *rank {
	r := rankPool.Get().(*rank)
	r.Doc = doc
	r.Weight = weight
	r.Matches = 1
	return r
}

func putRank(r *rank) {
	rankPool.Put(r)
}

func (r ranks) Len() int      { return len(r) }
func (r ranks) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ranks) Less(i, j int) bool {
	if r[i].Weight == r[j].Weight {
		if r[i].Doc.PublishDate() == r[j].Doc.PublishDate() {
			return r[i].Doc.Name() < r[j].Doc.Name()
		}
		return r[i].Doc.PublishDate().After(r[j].Doc.PublishDate())
	}
	return r[i].Weight > r[j].Weight
}

// SearchOpts holds the options for a related search.
type SearchOpts struct {
	// The Document to search for related content for.
	Document Document

	// The keywords to search for.
	NamedSlices []types.KeyValues

	// The indices to search in.
	Indices []string

	// Fragments holds a a list of special keywords that is used
	// for indices configured as type "fragments".
	// This will match the fragment identifiers of the documents.
	Fragments []string
}

// Search finds the documents matching any of the keywords in the given indices
// against query options in opts.
// The resulting document set will be sorted according to number of matches
// and the index weights, and any matches with a rank below the configured
// threshold (normalize to 0..100) will be removed.
// If an index name is provided, only that index will be queried.
func (idx *InvertedIndex) Search(ctx context.Context, opts SearchOpts) ([]Document, error) {
	var (
		queryElements []queryElement
		configs       IndicesConfig
	)

	if len(opts.Indices) == 0 {
		configs = idx.cfg.Indices
	} else {
		configs = make(IndicesConfig, len(opts.Indices))
		for i, indexName := range opts.Indices {
			cfg, found := idx.getIndexCfg(indexName)
			if !found {
				return nil, fmt.Errorf("index %q not found", indexName)
			}
			configs[i] = cfg
		}
	}

	for _, cfg := range configs {
		var keywords []Keyword
		if opts.Document != nil {
			k, err := opts.Document.RelatedKeywords(cfg)
			if err != nil {
				return nil, err
			}
			keywords = append(keywords, k...)
		}
		if cfg.Type == TypeFragments {
			for _, fragment := range opts.Fragments {
				keywords = append(keywords, FragmentKeyword(fragment))
			}
			if opts.Document != nil {
				if fp, ok := opts.Document.(FragmentProvider); ok {
					for _, fragment := range fp.Fragments(ctx).Identifiers {
						keywords = append(keywords, FragmentKeyword(fragment))
					}
				}
			}

		}
		queryElements = append(queryElements, newQueryElement(cfg.Name, keywords...))
	}
	for _, slice := range opts.NamedSlices {
		var keywords []Keyword
		key := slice.KeyString()
		if key == "" {
			return nil, fmt.Errorf("index %q not valid", slice.Key)
		}
		conf, found := idx.getIndexCfg(key)
		if !found {
			return nil, fmt.Errorf("index %q not found", key)
		}

		for _, val := range slice.Values {
			k, err := conf.ToKeywords(val)
			if err != nil {
				return nil, err
			}
			keywords = append(keywords, k...)
		}
		queryElements = append(queryElements, newQueryElement(conf.Name, keywords...))
	}

	if opts.Document != nil {
		return idx.searchDate(ctx, opts.Document, opts.Document.PublishDate(), queryElements...)
	}
	return idx.search(ctx, queryElements...)
}

func (cfg IndexConfig) stringToKeyword(s string) Keyword {
	if cfg.ToLower {
		s = strings.ToLower(s)
	}
	if cfg.Type == TypeFragments {
		return FragmentKeyword(s)
	}
	return StringKeyword(s)
}

// ToKeywords returns a Keyword slice of the given input.
func (cfg IndexConfig) ToKeywords(v any) ([]Keyword, error) {
	var keywords []Keyword

	switch vv := v.(type) {
	case string:
		keywords = append(keywords, cfg.stringToKeyword(vv))
	case []string:
		vvv := make([]Keyword, len(vv))
		for i := 0; i < len(vvv); i++ {
			vvv[i] = cfg.stringToKeyword(vv[i])
		}
		keywords = append(keywords, vvv...)
	case []any:
		return cfg.ToKeywords(cast.ToStringSlice(vv))
	case time.Time:
		layout := "2006"
		if cfg.Pattern != "" {
			layout = cfg.Pattern
		}
		keywords = append(keywords, StringKeyword(vv.Format(layout)))
	case nil:
		return keywords, nil
	default:
		return keywords, fmt.Errorf("indexing currently not supported for index %q and type %T", cfg.Name, vv)
	}

	return keywords, nil
}

func (idx *InvertedIndex) search(ctx context.Context, query ...queryElement) ([]Document, error) {
	return idx.searchDate(ctx, nil, zeroDate, query...)
}

func (idx *InvertedIndex) searchDate(ctx context.Context, self Document, upperDate time.Time, query ...queryElement) ([]Document, error) {
	matchm := make(map[Document]*rank, 200)
	defer func() {
		for _, r := range matchm {
			putRank(r)
		}
	}()

	applyDateFilter := !idx.cfg.IncludeNewer && !upperDate.IsZero()
	var fragmentsFilter collections.SortedStringSlice

	for _, el := range query {
		setm, found := idx.index[el.Index]
		if !found {
			return []Document{}, fmt.Errorf("index for %q not found", el.Index)
		}

		config, found := idx.getIndexCfg(el.Index)
		if !found {
			return []Document{}, fmt.Errorf("index config for %q not found", el.Index)
		}

		for _, kw := range el.Keywords {
			if docs, found := setm[kw]; found {
				for _, doc := range docs {
					if compare.Eq(doc, self) {
						continue
					}

					if applyDateFilter {
						// Exclude newer than the limit given
						if doc.PublishDate().After(upperDate) {
							continue
						}
					}

					if config.Type == TypeFragments && config.ApplyFilter {
						if fkw, ok := kw.(FragmentKeyword); ok {
							fragmentsFilter = append(fragmentsFilter, string(fkw))
						}
					}

					r, found := matchm[doc]
					if !found {
						r = getRank(doc, config.Weight)
						matchm[doc] = r
					} else {
						r.addWeight(config.Weight)
					}
				}
			}
		}
	}

	if len(matchm) == 0 {
		return []Document{}, nil
	}

	matches := make(ranks, 0, 100)

	for _, v := range matchm {
		avgWeight := v.Weight / v.Matches
		weight := norm(avgWeight, idx.minWeight, idx.maxWeight)
		threshold := idx.cfg.Threshold / v.Matches

		if weight >= threshold {
			matches = append(matches, v)
		}
	}

	sort.Stable(matches)
	sort.Strings(fragmentsFilter)

	result := make([]Document, len(matches))

	for i, m := range matches {
		result[i] = m.Doc

		if len(fragmentsFilter) > 0 {
			if dp, ok := result[i].(FragmentProvider); ok {
				result[i] = dp.ApplyFilterToHeadings(ctx, func(h *tableofcontents.Heading) bool {
					return fragmentsFilter.Contains(h.ID)
				})
			}
		}
	}

	return result, nil
}

// normalizes num to a number between 0 and 100.
func norm(num, min, max int) int {
	if min > max {
		panic("min > max")
	}
	return int(math.Floor((float64(num-min) / float64(max-min) * 100) + 0.5))
}

// DecodeConfig decodes a slice of map into Config.
func DecodeConfig(m maps.Params) (Config, error) {
	if m == nil {
		return Config{}, errors.New("no related config provided")
	}

	if len(m) == 0 {
		return Config{}, errors.New("empty related config provided")
	}

	var c Config

	if err := mapstructure.WeakDecode(m, &c); err != nil {
		return c, err
	}

	if c.Threshold < 0 || c.Threshold > 100 {
		return Config{}, errors.New("related threshold must be between 0 and 100")
	}

	if c.ToLower {
		for i := range c.Indices {
			c.Indices[i].ToLower = true
		}
	}
	for i := range c.Indices {
		icfg := c.Indices[i]
		if icfg.Type == "" {
			c.Indices[i].Type = TypeBasic
		}
		if !validTypes[c.Indices[i].Type] {
			return c, fmt.Errorf("invalid index type %q. Must be one of %v", c.Indices[i].Type, xmaps.Keys(validTypes))
		}
		if icfg.CardinalityThreshold < 0 || icfg.CardinalityThreshold > 100 {
			return Config{}, errors.New("cardinalityThreshold threshold must be between 0 and 100")
		}
	}

	return c, nil
}

// StringKeyword is a string search keyword.
type StringKeyword string

func (s StringKeyword) String() string {
	return string(s)
}

// FragmentKeyword represents a document fragment.
type FragmentKeyword string

func (f FragmentKeyword) String() string {
	return string(f)
}

// Keyword is the interface a keyword in the search index must implement.
type Keyword interface {
	String() string
}

// StringsToKeywords converts the given slice of strings to a slice of Keyword.
func (cfg IndexConfig) StringsToKeywords(s ...string) []Keyword {
	kw := make([]Keyword, len(s))

	for i := 0; i < len(s); i++ {
		kw[i] = cfg.stringToKeyword(s[i])
	}

	return kw
}
