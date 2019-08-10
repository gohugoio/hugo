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

package pagemeta

import (
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
)

func TestDateAndSlugFromBaseFilename(t *testing.T) {

	t.Parallel()

	c := qt.New(t)

	tests := []struct {
		name string
		date string
		slug string
	}{
		{"page.md", "0001-01-01", ""},
		{"2012-09-12-page.md", "2012-09-12", "page"},
		{"2018-02-28-page.md", "2018-02-28", "page"},
		{"2018-02-28_page.md", "2018-02-28", "page"},
		{"2018-02-28 page.md", "2018-02-28", "page"},
		{"2018-02-28page.md", "2018-02-28", "page"},
		{"2018-02-28-.md", "2018-02-28", ""},
		{"2018-02-28-.md", "2018-02-28", ""},
		{"2018-02-28.md", "2018-02-28", ""},
		{"2018-02-28-page", "2018-02-28", "page"},
		{"2012-9-12-page.md", "0001-01-01", ""},
		{"asdfasdf.md", "0001-01-01", ""},
	}

	for _, test := range tests {
		expecteFDate, err := time.Parse("2006-01-02", test.date)
		c.Assert(err, qt.IsNil)

		gotDate, gotSlug := dateAndSlugFromBaseFilename(test.name)

		c.Assert(gotDate, qt.Equals, expecteFDate)
		c.Assert(gotSlug, qt.Equals, test.slug)

	}
}

func newTestFd() *FrontMatterDescriptor {
	return &FrontMatterDescriptor{
		Frontmatter: make(map[string]interface{}),
		Params:      make(map[string]interface{}),
		Dates:       &resource.Dates{},
		PageURLs:    &URLPath{},
	}
}

func TestFrontMatterNewConfig(t *testing.T) {
	c := qt.New(t)

	cfg := viper.New()

	cfg.Set("frontmatter", map[string]interface{}{
		"date":        []string{"publishDate", "LastMod"},
		"Lastmod":     []string{"publishDate"},
		"expiryDate":  []string{"lastMod"},
		"publishDate": []string{"date"},
	})

	fc, err := newFrontmatterConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(fc.date, qt.DeepEquals, []string{"publishdate", "pubdate", "published", "lastmod", "modified"})
	c.Assert(fc.lastmod, qt.DeepEquals, []string{"publishdate", "pubdate", "published"})
	c.Assert(fc.expiryDate, qt.DeepEquals, []string{"lastmod", "modified"})
	c.Assert(fc.publishDate, qt.DeepEquals, []string{"date"})

	// Default
	cfg = viper.New()
	fc, err = newFrontmatterConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(fc.date, qt.DeepEquals, []string{"date", "publishdate", "pubdate", "published", "lastmod", "modified"})
	c.Assert(fc.lastmod, qt.DeepEquals, []string{":git", "lastmod", "modified", "date", "publishdate", "pubdate", "published"})
	c.Assert(fc.expiryDate, qt.DeepEquals, []string{"expirydate", "unpublishdate"})
	c.Assert(fc.publishDate, qt.DeepEquals, []string{"publishdate", "pubdate", "published", "date"})

	// :default keyword
	cfg.Set("frontmatter", map[string]interface{}{
		"date":        []string{"d1", ":default"},
		"lastmod":     []string{"d2", ":default"},
		"expiryDate":  []string{"d3", ":default"},
		"publishDate": []string{"d4", ":default"},
	})
	fc, err = newFrontmatterConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(fc.date, qt.DeepEquals, []string{"d1", "date", "publishdate", "pubdate", "published", "lastmod", "modified"})
	c.Assert(fc.lastmod, qt.DeepEquals, []string{"d2", ":git", "lastmod", "modified", "date", "publishdate", "pubdate", "published"})
	c.Assert(fc.expiryDate, qt.DeepEquals, []string{"d3", "expirydate", "unpublishdate"})
	c.Assert(fc.publishDate, qt.DeepEquals, []string{"d4", "publishdate", "pubdate", "published", "date"})

}

func TestFrontMatterDatesHandlers(t *testing.T) {
	c := qt.New(t)

	for _, handlerID := range []string{":filename", ":fileModTime", ":git"} {

		cfg := viper.New()

		cfg.Set("frontmatter", map[string]interface{}{
			"date": []string{handlerID, "date"},
		})

		handler, err := NewFrontmatterHandler(nil, cfg)
		c.Assert(err, qt.IsNil)

		d1, _ := time.Parse("2006-01-02", "2018-02-01")
		d2, _ := time.Parse("2006-01-02", "2018-02-02")

		d := newTestFd()
		switch strings.ToLower(handlerID) {
		case ":filename":
			d.BaseFilename = "2018-02-01-page.md"
		case ":filemodtime":
			d.ModTime = d1
		case ":git":
			d.GitAuthorDate = d1
		}
		d.Frontmatter["date"] = d2
		c.Assert(handler.HandleDates(d), qt.IsNil)
		c.Assert(d.Dates.FDate, qt.Equals, d1)
		c.Assert(d.Params["date"], qt.Equals, d2)

		d = newTestFd()
		d.Frontmatter["date"] = d2
		c.Assert(handler.HandleDates(d), qt.IsNil)
		c.Assert(d.Dates.FDate, qt.Equals, d2)
		c.Assert(d.Params["date"], qt.Equals, d2)

	}
}

func TestFrontMatterDatesCustomConfig(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	cfg := viper.New()
	cfg.Set("frontmatter", map[string]interface{}{
		"date":        []string{"mydate"},
		"lastmod":     []string{"publishdate"},
		"publishdate": []string{"publishdate"},
	})

	handler, err := NewFrontmatterHandler(nil, cfg)
	c.Assert(err, qt.IsNil)

	testDate, err := time.Parse("2006-01-02", "2018-02-01")
	c.Assert(err, qt.IsNil)

	d := newTestFd()
	d.Frontmatter["mydate"] = testDate
	testDate = testDate.Add(24 * time.Hour)
	d.Frontmatter["date"] = testDate
	testDate = testDate.Add(24 * time.Hour)
	d.Frontmatter["lastmod"] = testDate
	testDate = testDate.Add(24 * time.Hour)
	d.Frontmatter["publishdate"] = testDate
	testDate = testDate.Add(24 * time.Hour)
	d.Frontmatter["expirydate"] = testDate

	c.Assert(handler.HandleDates(d), qt.IsNil)

	c.Assert(d.Dates.FDate.Day(), qt.Equals, 1)
	c.Assert(d.Dates.FLastmod.Day(), qt.Equals, 4)
	c.Assert(d.Dates.FPublishDate.Day(), qt.Equals, 4)
	c.Assert(d.Dates.FExpiryDate.Day(), qt.Equals, 5)

	c.Assert(d.Params["date"], qt.Equals, d.Dates.FDate)
	c.Assert(d.Params["mydate"], qt.Equals, d.Dates.FDate)
	c.Assert(d.Params["publishdate"], qt.Equals, d.Dates.FPublishDate)
	c.Assert(d.Params["expirydate"], qt.Equals, d.Dates.FExpiryDate)

	c.Assert(handler.IsDateKey("date"), qt.Equals, false) // This looks odd, but is configured like this.
	c.Assert(handler.IsDateKey("mydate"), qt.Equals, true)
	c.Assert(handler.IsDateKey("publishdate"), qt.Equals, true)
	c.Assert(handler.IsDateKey("pubdate"), qt.Equals, true)

}

func TestFrontMatterDatesDefaultKeyword(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	cfg := viper.New()

	cfg.Set("frontmatter", map[string]interface{}{
		"date":        []string{"mydate", ":default"},
		"publishdate": []string{":default", "mypubdate"},
	})

	handler, err := NewFrontmatterHandler(nil, cfg)
	c.Assert(err, qt.IsNil)

	testDate, _ := time.Parse("2006-01-02", "2018-02-01")
	d := newTestFd()
	d.Frontmatter["mydate"] = testDate
	d.Frontmatter["date"] = testDate.Add(1 * 24 * time.Hour)
	d.Frontmatter["mypubdate"] = testDate.Add(2 * 24 * time.Hour)
	d.Frontmatter["publishdate"] = testDate.Add(3 * 24 * time.Hour)

	c.Assert(handler.HandleDates(d), qt.IsNil)

	c.Assert(d.Dates.FDate.Day(), qt.Equals, 1)
	c.Assert(d.Dates.FLastmod.Day(), qt.Equals, 2)
	c.Assert(d.Dates.FPublishDate.Day(), qt.Equals, 4)
	c.Assert(d.Dates.FExpiryDate.IsZero(), qt.Equals, true)

}

func TestExpandDefaultValues(t *testing.T) {
	c := qt.New(t)
	c.Assert(expandDefaultValues([]string{"a", ":default", "d"}, []string{"b", "c"}), qt.DeepEquals, []string{"a", "b", "c", "d"})
	c.Assert(expandDefaultValues([]string{"a", "b", "c"}, []string{"a", "b", "c"}), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(expandDefaultValues([]string{":default", "a", ":default", "d"}, []string{"b", "c"}), qt.DeepEquals, []string{"b", "c", "a", "b", "c", "d"})

}

func TestFrontMatterDateFieldHandler(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	handlers := new(frontmatterFieldHandlers)

	fd := newTestFd()
	d, _ := time.Parse("2006-01-02", "2018-02-01")
	fd.Frontmatter["date"] = d
	h := handlers.newDateFieldHandler("date", func(d *FrontMatterDescriptor, t time.Time) { d.Dates.FDate = t })

	handled, err := h(fd)
	c.Assert(handled, qt.Equals, true)
	c.Assert(err, qt.IsNil)
	c.Assert(fd.Dates.FDate, qt.Equals, d)
}
