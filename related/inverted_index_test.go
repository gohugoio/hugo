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

package related

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

type testDoc struct {
	keywords map[string][]Keyword
	date     time.Time
	name     string
}

func (d *testDoc) String() string {
	s := "\n"
	for k, v := range d.keywords {
		s += k + ":\t\t"
		for _, vv := range v {
			s += "  " + vv.String()
		}
		s += "\n"
	}
	return s
}

func (d *testDoc) Name() string {
	return d.name
}

func newTestDoc(name string, keywords ...string) *testDoc {
	time.Sleep(1 * time.Millisecond)
	return newTestDocWithDate(name, time.Now(), keywords...)
}

func newTestDocWithDate(name string, date time.Time, keywords ...string) *testDoc {
	km := make(map[string][]Keyword)

	kw := &testDoc{keywords: km, date: date}

	kw.addKeywords(name, keywords...)
	return kw
}

func (d *testDoc) addKeywords(name string, keywords ...string) *testDoc {
	keywordm := createTestKeywords(name, keywords...)

	for k, v := range keywordm {
		keywords := make([]Keyword, len(v))
		for i := 0; i < len(v); i++ {
			keywords[i] = StringKeyword(v[i])
		}
		d.keywords[k] = keywords
	}
	return d
}

func createTestKeywords(name string, keywords ...string) map[string][]string {
	return map[string][]string{
		name: keywords,
	}
}

func (d *testDoc) RelatedKeywords(cfg IndexConfig) ([]Keyword, error) {
	return d.keywords[cfg.Name], nil
}

func (d *testDoc) PublishDate() time.Time {
	return d.date
}

func TestSearch(t *testing.T) {

	config := Config{
		Threshold:    90,
		IncludeNewer: false,
		Indices: IndexConfigs{
			IndexConfig{Name: "tags", Weight: 50},
			IndexConfig{Name: "keywords", Weight: 65},
		},
	}

	idx := NewInvertedIndex(config)
	//idx.debug = true

	docs := []Document{
		newTestDoc("tags", "a", "b", "c", "d"),
		newTestDoc("tags", "b", "d", "g"),
		newTestDoc("tags", "b", "h").addKeywords("keywords", "a"),
		newTestDoc("tags", "g", "h").addKeywords("keywords", "a", "b"),
	}

	idx.Add(docs...)

	t.Run("count", func(t *testing.T) {
		c := qt.New(t)
		c.Assert(len(idx.index), qt.Equals, 2)
		set1, found := idx.index["tags"]
		c.Assert(found, qt.Equals, true)
		// 6 tags
		c.Assert(len(set1), qt.Equals, 6)

		set2, found := idx.index["keywords"]
		c.Assert(found, qt.Equals, true)
		c.Assert(len(set2), qt.Equals, 2)

	})

	t.Run("search-tags", func(t *testing.T) {
		c := qt.New(t)
		m, err := idx.search(newQueryElement("tags", StringsToKeywords("a", "b", "d", "z")...))
		c.Assert(err, qt.IsNil)
		c.Assert(len(m), qt.Equals, 2)
		c.Assert(m[0], qt.Equals, docs[0])
		c.Assert(m[1], qt.Equals, docs[1])
	})

	t.Run("search-tags-and-keywords", func(t *testing.T) {
		c := qt.New(t)
		m, err := idx.search(
			newQueryElement("tags", StringsToKeywords("a", "b", "z")...),
			newQueryElement("keywords", StringsToKeywords("a", "b")...))
		c.Assert(err, qt.IsNil)
		c.Assert(len(m), qt.Equals, 3)
		c.Assert(m[0], qt.Equals, docs[3])
		c.Assert(m[1], qt.Equals, docs[2])
		c.Assert(m[2], qt.Equals, docs[0])
	})

	t.Run("searchdoc-all", func(t *testing.T) {
		c := qt.New(t)
		doc := newTestDoc("tags", "a").addKeywords("keywords", "a")
		m, err := idx.SearchDoc(doc)
		c.Assert(err, qt.IsNil)
		c.Assert(len(m), qt.Equals, 2)
		c.Assert(m[0], qt.Equals, docs[3])
		c.Assert(m[1], qt.Equals, docs[2])
	})

	t.Run("searchdoc-tags", func(t *testing.T) {
		c := qt.New(t)
		doc := newTestDoc("tags", "a", "b", "d", "z").addKeywords("keywords", "a", "b")
		m, err := idx.SearchDoc(doc, "tags")
		c.Assert(err, qt.IsNil)
		c.Assert(len(m), qt.Equals, 2)
		c.Assert(m[0], qt.Equals, docs[0])
		c.Assert(m[1], qt.Equals, docs[1])
	})

	t.Run("searchdoc-keywords-date", func(t *testing.T) {
		c := qt.New(t)
		doc := newTestDoc("tags", "a", "b", "d", "z").addKeywords("keywords", "a", "b")
		// This will get a date newer than the others.
		newDoc := newTestDoc("keywords", "a", "b")
		idx.Add(newDoc)

		m, err := idx.SearchDoc(doc, "keywords")
		c.Assert(err, qt.IsNil)
		c.Assert(len(m), qt.Equals, 2)
		c.Assert(m[0], qt.Equals, docs[3])
	})

	t.Run("searchdoc-keywords-same-date", func(t *testing.T) {
		c := qt.New(t)
		idx := NewInvertedIndex(config)

		date := time.Now()

		doc := newTestDocWithDate("keywords", date, "a", "b")
		doc.name = "thedoc"

		for i := 0; i < 10; i++ {
			docc := *doc
			docc.name = fmt.Sprintf("doc%d", i)
			idx.Add(&docc)
		}

		m, err := idx.SearchDoc(doc, "keywords")
		c.Assert(err, qt.IsNil)
		c.Assert(len(m), qt.Equals, 10)
		for i := 0; i < 10; i++ {
			c.Assert(m[i].Name(), qt.Equals, fmt.Sprintf("doc%d", i))
		}
	})

}

func BenchmarkRelatedNewIndex(b *testing.B) {

	pages := make([]*testDoc, 100)
	numkeywords := 30
	allKeywords := make([]string, numkeywords)
	for i := 0; i < numkeywords; i++ {
		allKeywords[i] = fmt.Sprintf("keyword%d", i+1)
	}

	for i := 0; i < len(pages); i++ {
		start := rand.Intn(len(allKeywords))
		end := start + 3
		if end >= len(allKeywords) {
			end = start + 1
		}

		kw := newTestDoc("tags", allKeywords[start:end]...)
		if i%5 == 0 {
			start := rand.Intn(len(allKeywords))
			end := start + 3
			if end >= len(allKeywords) {
				end = start + 1
			}
			kw.addKeywords("keywords", allKeywords[start:end]...)
		}

		pages[i] = kw
	}

	cfg := Config{
		Threshold: 50,
		Indices: IndexConfigs{
			IndexConfig{Name: "tags", Weight: 100},
			IndexConfig{Name: "keywords", Weight: 200},
		},
	}

	b.Run("singles", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := NewInvertedIndex(cfg)
			for _, doc := range pages {
				idx.Add(doc)
			}
		}
	})

	b.Run("all", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := NewInvertedIndex(cfg)
			docs := make([]Document, len(pages))
			for i := 0; i < len(pages); i++ {
				docs[i] = pages[i]
			}
			idx.Add(docs...)
		}
	})

}

func BenchmarkRelatedMatchesIn(b *testing.B) {

	q1 := newQueryElement("tags", StringsToKeywords("keyword2", "keyword5", "keyword32", "asdf")...)
	q2 := newQueryElement("keywords", StringsToKeywords("keyword3", "keyword4")...)

	docs := make([]*testDoc, 1000)
	numkeywords := 20
	allKeywords := make([]string, numkeywords)
	for i := 0; i < numkeywords; i++ {
		allKeywords[i] = fmt.Sprintf("keyword%d", i+1)
	}

	cfg := Config{
		Threshold: 20,
		Indices: IndexConfigs{
			IndexConfig{Name: "tags", Weight: 100},
			IndexConfig{Name: "keywords", Weight: 200},
		},
	}

	idx := NewInvertedIndex(cfg)

	for i := 0; i < len(docs); i++ {
		start := rand.Intn(len(allKeywords))
		end := start + 3
		if end >= len(allKeywords) {
			end = start + 1
		}

		index := "tags"
		if i%5 == 0 {
			index = "keywords"
		}

		idx.Add(newTestDoc(index, allKeywords[start:end]...))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%10 == 0 {
			idx.search(q2)
		} else {
			idx.search(q1)
		}
	}
}
