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

package pagemeta

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/require"
)

func TestDateAndSlugFromBaseFilename(t *testing.T) {

	t.Parallel()

	assert := require.New(t)

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

	for i, test := range tests {
		expectedDate, err := time.Parse("2006-01-02", test.date)
		assert.NoError(err)

		errMsg := fmt.Sprintf("Test %d", i)
		gotDate, gotSlug := dateAndSlugFromBaseFilename(test.name)

		assert.Equal(expectedDate, gotDate, errMsg)
		assert.Equal(test.slug, gotSlug, errMsg)

	}
}

func newTestFd() *FrontMatterDescriptor {
	return &FrontMatterDescriptor{
		Frontmatter: make(map[string]interface{}),
		Params:      make(map[string]interface{}),
		Dates:       &PageDates{},
		PageURLs:    &URLPath{},
	}
}

func TestFrontMatterNewConfig(t *testing.T) {
	assert := require.New(t)

	cfg := viper.New()

	cfg.Set("frontmatter", map[string]interface{}{
		"date":        []string{"publishDate", "LastMod"},
		"Lastmod":     []string{"publishDate"},
		"expiryDate":  []string{"lastMod"},
		"publishDate": []string{"date"},
	})

	fc, err := newFrontmatterConfig(cfg)
	assert.NoError(err)
	assert.Equal([]string{"publishdate", "pubdate", "published", "lastmod", "modified"}, fc.date)
	assert.Equal([]string{"publishdate", "pubdate", "published"}, fc.lastmod)
	assert.Equal([]string{"lastmod", "modified"}, fc.expiryDate)
	assert.Equal([]string{"date"}, fc.publishDate)

	// Default
	cfg = viper.New()
	fc, err = newFrontmatterConfig(cfg)
	assert.NoError(err)
	assert.Equal([]string{"date", "publishdate", "pubdate", "published", "lastmod", "modified"}, fc.date)
	assert.Equal([]string{":git", "lastmod", "modified", "date", "publishdate", "pubdate", "published"}, fc.lastmod)
	assert.Equal([]string{"expirydate", "unpublishdate"}, fc.expiryDate)
	assert.Equal([]string{"publishdate", "pubdate", "published", "date"}, fc.publishDate)

	// :default keyword
	cfg.Set("frontmatter", map[string]interface{}{
		"date":        []string{"d1", ":default"},
		"lastmod":     []string{"d2", ":default"},
		"expiryDate":  []string{"d3", ":default"},
		"publishDate": []string{"d4", ":default"},
	})
	fc, err = newFrontmatterConfig(cfg)
	assert.NoError(err)
	assert.Equal([]string{"d1", "date", "publishdate", "pubdate", "published", "lastmod", "modified"}, fc.date)
	assert.Equal([]string{"d2", ":git", "lastmod", "modified", "date", "publishdate", "pubdate", "published"}, fc.lastmod)
	assert.Equal([]string{"d3", "expirydate", "unpublishdate"}, fc.expiryDate)
	assert.Equal([]string{"d4", "publishdate", "pubdate", "published", "date"}, fc.publishDate)

}

func TestFrontMatterDatesHandlers(t *testing.T) {
	assert := require.New(t)

	for _, handlerID := range []string{":filename", ":fileModTime", ":git"} {

		cfg := viper.New()

		cfg.Set("frontmatter", map[string]interface{}{
			"date": []string{handlerID, "date"},
		})

		handler, err := NewFrontmatterHandler(nil, cfg)
		assert.NoError(err)

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
		assert.NoError(handler.HandleDates(d))
		assert.Equal(d1, d.Dates.Date)
		assert.Equal(d2, d.Params["date"])

		d = newTestFd()
		d.Frontmatter["date"] = d2
		assert.NoError(handler.HandleDates(d))
		assert.Equal(d2, d.Dates.Date)
		assert.Equal(d2, d.Params["date"])

	}
}

func TestFrontMatterDatesCustomConfig(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	cfg := viper.New()
	cfg.Set("frontmatter", map[string]interface{}{
		"date":        []string{"mydate"},
		"lastmod":     []string{"publishdate"},
		"publishdate": []string{"publishdate"},
	})

	handler, err := NewFrontmatterHandler(nil, cfg)
	assert.NoError(err)

	testDate, err := time.Parse("2006-01-02", "2018-02-01")
	assert.NoError(err)

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

	assert.NoError(handler.HandleDates(d))

	assert.Equal(1, d.Dates.Date.Day())
	assert.Equal(4, d.Dates.Lastmod.Day())
	assert.Equal(4, d.Dates.PublishDate.Day())
	assert.Equal(5, d.Dates.ExpiryDate.Day())

	assert.Equal(d.Dates.Date, d.Params["date"])
	assert.Equal(d.Dates.Date, d.Params["mydate"])
	assert.Equal(d.Dates.PublishDate, d.Params["publishdate"])
	assert.Equal(d.Dates.ExpiryDate, d.Params["expirydate"])

	assert.False(handler.IsDateKey("date")) // This looks odd, but is configured like this.
	assert.True(handler.IsDateKey("mydate"))
	assert.True(handler.IsDateKey("publishdate"))
	assert.True(handler.IsDateKey("pubdate"))

}

func TestFrontMatterDatesDefaultKeyword(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	cfg := viper.New()

	cfg.Set("frontmatter", map[string]interface{}{
		"date":        []string{"mydate", ":default"},
		"publishdate": []string{":default", "mypubdate"},
	})

	handler, err := NewFrontmatterHandler(nil, cfg)
	assert.NoError(err)

	testDate, _ := time.Parse("2006-01-02", "2018-02-01")
	d := newTestFd()
	d.Frontmatter["mydate"] = testDate
	d.Frontmatter["date"] = testDate.Add(1 * 24 * time.Hour)
	d.Frontmatter["mypubdate"] = testDate.Add(2 * 24 * time.Hour)
	d.Frontmatter["publishdate"] = testDate.Add(3 * 24 * time.Hour)

	assert.NoError(handler.HandleDates(d))

	assert.Equal(1, d.Dates.Date.Day())
	assert.Equal(2, d.Dates.Lastmod.Day())
	assert.Equal(4, d.Dates.PublishDate.Day())
	assert.True(d.Dates.ExpiryDate.IsZero())

}

func TestExpandDefaultValues(t *testing.T) {
	assert := require.New(t)
	assert.Equal([]string{"a", "b", "c", "d"}, expandDefaultValues([]string{"a", ":default", "d"}, []string{"b", "c"}))
	assert.Equal([]string{"a", "b", "c"}, expandDefaultValues([]string{"a", "b", "c"}, []string{"a", "b", "c"}))
	assert.Equal([]string{"b", "c", "a", "b", "c", "d"}, expandDefaultValues([]string{":default", "a", ":default", "d"}, []string{"b", "c"}))

}

func TestFrontMatterDateFieldHandler(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	handlers := new(frontmatterFieldHandlers)

	fd := newTestFd()
	d, _ := time.Parse("2006-01-02", "2018-02-01")
	fd.Frontmatter["date"] = d
	h := handlers.newDateFieldHandler("date", func(d *FrontMatterDescriptor, t time.Time) { d.Dates.Date = t })

	handled, err := h(fd)
	assert.True(handled)
	assert.NoError(err)
	assert.Equal(d, fd.Dates.Date)
}
