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

package pagemeta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"

	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/gohugoio/hugo/resources/resource"

	qt "github.com/frankban/quicktest"
)

func newTestFd() *pagemeta.FrontMatterDescriptor {
	return &pagemeta.FrontMatterDescriptor{
		Params:   make(map[string]any),
		Dates:    &resource.Dates{},
		PageURLs: &pagemeta.URLPath{},
		Location: time.UTC,
	}
}

func TestFrontMatterNewConfig(t *testing.T) {
	c := qt.New(t)

	cfg := config.New()

	cfg.Set("frontmatter", map[string]any{
		"date":        []string{"publishDate", "LastMod"},
		"Lastmod":     []string{"publishDate"},
		"expiryDate":  []string{"lastMod"},
		"publishDate": []string{"date"},
	})

	fc, err := pagemeta.DecodeFrontMatterConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(fc.Date, qt.DeepEquals, []string{"publishdate", "pubdate", "published", "lastmod", "modified"})
	c.Assert(fc.Lastmod, qt.DeepEquals, []string{"publishdate", "pubdate", "published"})
	c.Assert(fc.ExpiryDate, qt.DeepEquals, []string{"lastmod", "modified"})
	c.Assert(fc.PublishDate, qt.DeepEquals, []string{"date"})

	// Default
	cfg = config.New()
	fc, err = pagemeta.DecodeFrontMatterConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(fc.Date, qt.DeepEquals, []string{"date", "publishdate", "pubdate", "published", "lastmod", "modified"})
	c.Assert(fc.Lastmod, qt.DeepEquals, []string{":git", "lastmod", "modified", "date", "publishdate", "pubdate", "published"})
	c.Assert(fc.ExpiryDate, qt.DeepEquals, []string{"expirydate", "unpublishdate"})
	c.Assert(fc.PublishDate, qt.DeepEquals, []string{"publishdate", "pubdate", "published", "date"})

	// :default keyword
	cfg.Set("frontmatter", map[string]any{
		"date":        []string{"d1", ":default"},
		"lastmod":     []string{"d2", ":default"},
		"expiryDate":  []string{"d3", ":default"},
		"publishDate": []string{"d4", ":default"},
	})
	fc, err = pagemeta.DecodeFrontMatterConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(fc.Date, qt.DeepEquals, []string{"d1", "date", "publishdate", "pubdate", "published", "lastmod", "modified"})
	c.Assert(fc.Lastmod, qt.DeepEquals, []string{"d2", ":git", "lastmod", "modified", "date", "publishdate", "pubdate", "published"})
	c.Assert(fc.ExpiryDate, qt.DeepEquals, []string{"d3", "expirydate", "unpublishdate"})
	c.Assert(fc.PublishDate, qt.DeepEquals, []string{"d4", "publishdate", "pubdate", "published", "date"})
}

func TestFrontMatterDatesHandlers(t *testing.T) {
	c := qt.New(t)

	for _, handlerID := range []string{":filename", ":fileModTime", ":git"} {

		cfg := config.New()

		cfg.Set("frontmatter", map[string]any{
			"date": []string{handlerID, "date"},
		})
		conf := testconfig.GetTestConfig(nil, cfg)
		handler, err := pagemeta.NewFrontmatterHandler(nil, conf.GetConfigSection("frontmatter").(pagemeta.FrontmatterConfig))
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
		d.Params["date"] = d2
		c.Assert(handler.HandleDates(d), qt.IsNil)
		c.Assert(d.Dates.FDate, qt.Equals, d1)
		c.Assert(d.Params["date"], qt.Equals, d2)

		d = newTestFd()
		d.Params["date"] = d2
		c.Assert(handler.HandleDates(d), qt.IsNil)
		c.Assert(d.Dates.FDate, qt.Equals, d2)
		c.Assert(d.Params["date"], qt.Equals, d2)

	}
}

func TestFrontMatterDatesDefaultKeyword(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	cfg := config.New()

	cfg.Set("frontmatter", map[string]any{
		"date":        []string{"mydate", ":default"},
		"publishdate": []string{":default", "mypubdate"},
	})

	conf := testconfig.GetTestConfig(nil, cfg)
	handler, err := pagemeta.NewFrontmatterHandler(nil, conf.GetConfigSection("frontmatter").(pagemeta.FrontmatterConfig))
	c.Assert(err, qt.IsNil)

	testDate, _ := time.Parse("2006-01-02", "2018-02-01")
	d := newTestFd()
	d.Params["mydate"] = testDate
	d.Params["date"] = testDate.Add(1 * 24 * time.Hour)
	d.Params["mypubdate"] = testDate.Add(2 * 24 * time.Hour)
	d.Params["publishdate"] = testDate.Add(3 * 24 * time.Hour)

	c.Assert(handler.HandleDates(d), qt.IsNil)

	c.Assert(d.Dates.FDate.Day(), qt.Equals, 1)
	c.Assert(d.Dates.FLastmod.Day(), qt.Equals, 2)
	c.Assert(d.Dates.FPublishDate.Day(), qt.Equals, 4)
	c.Assert(d.Dates.FExpiryDate.IsZero(), qt.Equals, true)
}
