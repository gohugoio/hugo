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
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/types"
	"github.com/mitchellh/mapstructure"
)

var (
	_        Keyword = (*StringKeyword)(nil)
	zeroDate         = time.Time{}

	// DefaultConfig is the default related config.
	DefaultConfig = Config{
		Threshold: 80,
		Indices: IndexConfigs{
			IndexConfig{Name: "keywords", Weight: 100},
			IndexConfig{Name: "date", Weight: 10},
		},
	}
)

/*
Config is the top level configuration element used to configure how to retrieve
related content in Hugo.

An example site config.toml:

	[related]
	threshold = 1
	[[related.indices]]
	name = "keywords"
	weight = 200
	[[related.indices]]
	name  = "tags"
	weight = 100
	[[related.indices]]
	name  = "date"
	weight = 1
	pattern = "2006"
*/
type Config struct {
	// Only include matches >= threshold, a normalized rank between 0 and 100.
	Threshold int

	// To get stable "See also" sections we, by default, exclude newer related pages.
	IncludeNewer bool

	// Will lower case all string values and queries to the indices.
	// May get better results, but at a slight performance cost.
	ToLower bool

	Indices IndexConfigs
}

// Add adds a given index.
func (c *Config) Add(index IndexConfig) {
	if c.ToLower {
		index.ToLower = true
	}
	c.Indices = append(c.Indices, index)
}

// IndexConfigs holds a set of index configurations.
type IndexConfigs []IndexConfig

// IndexConfig configures an index.
type IndexConfig struct {
	// The index name. This directly maps to a field or Param name.
	Name string

	// Contextual pattern used to convert the Param value into a string.
	// Currently only used for dates. Can be used to, say, bump posts in the same
	// time frame when searching for related documents.
	// For dates it follows Go's time.Format patterns, i.e.
	// "2006" for YYYY and "200601" for YYYYMM.
	Pattern string

	// This field's weight when doing multi-index searches. Higher is "better".
	Weight int

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

// InvertedIndex holds an inverted index, also sometimes named posting list, which
// lists, for every possible search term, the documents that contain that term.
type InvertedIndex struct {
	cfg   Config
	index map[string]map[Keyword][]Document

	minWeight int
	maxWeight int
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
	idx := &InvertedIndex{index: make(map[string]map[Keyword][]Document), cfg: cfg}
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
func (idx *InvertedIndex) Add(docs ...Document) error {
	var err error
	for _, config := range idx.cfg.Indices {
		if config.Weight == 0 {
			// Disabled
			continue
		}
		setm := idx.index[config.Name]

		for _, doc := range docs {
			var words []Keyword
			words, err = doc.RelatedKeywords(config)
			if err != nil {
				continue
			}

			for _, keyword := range words {
				setm[keyword] = append(setm[keyword], doc)
			}
		}
	}

	return err

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

func newRank(doc Document, weight int) *rank {
	return &rank{Doc: doc, Weight: weight, Matches: 1}
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

// SearchDoc finds the documents matching any of the keywords in the given indices
// against the given document.
// The resulting document set will be sorted according to number of matches
// and the index weights, and any matches with a rank below the configured
// threshold (normalize to 0..100) will be removed.
// If an index name is provided, only that index will be queried.
func (idx *InvertedIndex) SearchDoc(doc Document, indices ...string) ([]Document, error) {
	var q []queryElement

	var configs IndexConfigs

	if len(indices) == 0 {
		configs = idx.cfg.Indices
	} else {
		configs = make(IndexConfigs, len(indices))
		for i, indexName := range indices {
			cfg, found := idx.getIndexCfg(indexName)
			if !found {
				return nil, fmt.Errorf("index %q not found", indexName)
			}
			configs[i] = cfg
		}
	}

	for _, cfg := range configs {
		keywords, err := doc.RelatedKeywords(cfg)
		if err != nil {
			return nil, err
		}

		q = append(q, newQueryElement(cfg.Name, keywords...))

	}

	return idx.searchDate(doc.PublishDate(), q...)
}

// ToKeywords returns a Keyword slice of the given input.
func (cfg IndexConfig) ToKeywords(v interface{}) ([]Keyword, error) {
	var (
		keywords []Keyword
		toLower  = cfg.ToLower
	)
	switch vv := v.(type) {
	case string:
		if toLower {
			vv = strings.ToLower(vv)
		}
		keywords = append(keywords, StringKeyword(vv))
	case []string:
		if toLower {
			for i := 0; i < len(vv); i++ {
				vv[i] = strings.ToLower(vv[i])
			}
		}
		keywords = append(keywords, StringsToKeywords(vv...)...)
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

// SearchKeyValues finds the documents matching any of the keywords in the given indices.
// The resulting document set will be sorted according to number of matches
// and the index weights, and any matches with a rank below the configured
// threshold (normalize to 0..100) will be removed.
func (idx *InvertedIndex) SearchKeyValues(args ...types.KeyValues) ([]Document, error) {
	q := make([]queryElement, len(args))

	for i, arg := range args {
		var keywords []Keyword
		key := arg.KeyString()
		if key == "" {
			return nil, fmt.Errorf("index %q not valid", arg.Key)
		}
		conf, found := idx.getIndexCfg(key)
		if !found {
			return nil, fmt.Errorf("index %q not found", key)
		}

		for _, val := range arg.Values {
			k, err := conf.ToKeywords(val)
			if err != nil {
				return nil, err
			}
			keywords = append(keywords, k...)
		}

		q[i] = newQueryElement(conf.Name, keywords...)

	}

	return idx.search(q...)
}

func (idx *InvertedIndex) search(query ...queryElement) ([]Document, error) {
	return idx.searchDate(zeroDate, query...)
}

func (idx *InvertedIndex) searchDate(upperDate time.Time, query ...queryElement) ([]Document, error) {
	matchm := make(map[Document]*rank, 200)
	applyDateFilter := !idx.cfg.IncludeNewer && !upperDate.IsZero()

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
					if applyDateFilter {
						// Exclude newer than the limit given
						if doc.PublishDate().After(upperDate) {
							continue
						}
					}
					r, found := matchm[doc]
					if !found {
						matchm[doc] = newRank(doc, config.Weight)
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

	result := make([]Document, len(matches))

	for i, m := range matches {
		result[i] = m.Doc
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
func DecodeConfig(in interface{}) (Config, error) {
	if in == nil {
		return Config{}, errors.New("no related config provided")
	}

	m, ok := in.(map[string]interface{})
	if !ok {
		return Config{}, fmt.Errorf("expected map[string]interface {} got %T", in)
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

	return c, nil
}

// StringKeyword is a string search keyword.
type StringKeyword string

func (s StringKeyword) String() string {
	return string(s)
}

// Keyword is the interface a keyword in the search index must implement.
type Keyword interface {
	String() string
}

// StringsToKeywords converts the given slice of strings to a slice of Keyword.
func StringsToKeywords(s ...string) []Keyword {
	kw := make([]Keyword, len(s))

	for i := 0; i < len(s); i++ {
		kw[i] = StringKeyword(s[i])
	}

	return kw
}
