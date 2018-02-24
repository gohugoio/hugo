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
	assert.Equal([]string{"publishdate", "lastmod"}, fc.date)
	assert.Equal([]string{"publishdate"}, fc.lastMod)
	assert.Equal([]string{"lastmod"}, fc.expiryDate)
	assert.Equal([]string{"date"}, fc.publishDate)

	// Default
	cfg = viper.New()
	fc, err = newFrontmatterConfig(cfg)
	assert.NoError(err)
	assert.Equal(3, len(fc.date))
	assert.Equal(2, len(fc.lastMod))
	assert.Equal(2, len(fc.publishDate))
	assert.Equal(1, len(fc.expiryDate))

}

func TestFrontMatterDatesConfigVariations(t *testing.T) {
	cfg := viper.New()

	cfg.Set("frontmatter", map[string]interface{}{
		"defaultDate": []string{"date"},
	})

	fmt.Println(">>", cfg)
}

func TestFrontMatterDates(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	cfg := viper.New()

	handler, err := NewFrontmatterHandler(nil, cfg)
	assert.NoError(err)

	testDate, err := time.Parse("2006-01-02", "2018-02-01")
	assert.NoError(err)

	sentinel := (time.Time{}).Add(1 * time.Hour)
	zero := time.Time{}

	// See http://www.imdb.com/title/tt0133093/
	for _, lastModKey := range []string{"lastmod", "modified"} {
		testDate = testDate.Add(24 * time.Hour)
		t.Log(lastModKey, testDate)
		for _, lastModDate := range []time.Time{testDate, sentinel} {
			for _, dateKey := range []string{"date"} {
				testDate = testDate.Add(24 * time.Hour)
				t.Log(dateKey, testDate)
				for _, dateDate := range []time.Time{testDate, sentinel} {
					for _, pubDateKey := range []string{"publishdate", "pubdate", "published"} {
						testDate = testDate.Add(24 * time.Hour)
						t.Log(pubDateKey, testDate)
						for _, pubDateDate := range []time.Time{testDate, sentinel} {
							for _, expiryDateKey := range []string{"expirydate", "unpublishdate"} {
								testDate = testDate.Add(24 * time.Hour)
								t.Log(expiryDateKey, testDate)
								for _, expiryDateDate := range []time.Time{testDate, sentinel} {
									d := FrontMatterDescriptor{
										Frontmatter: make(map[string]interface{}),
										Params:      make(map[string]interface{}),
										Dates:       &PageDates{},
										PageURLs:    &URLPath{},
									}

									var expLastMod, expDate, expPubDate, expExiryDate = zero, zero, zero, zero

									if dateDate != sentinel {
										d.Frontmatter[dateKey] = dateDate
										expDate = dateDate
									}

									if pubDateDate != sentinel {
										d.Frontmatter[pubDateKey] = pubDateDate
										expPubDate = pubDateDate
										if expDate.IsZero() {
											expDate = expPubDate
										}
									}

									if lastModDate != sentinel {
										d.Frontmatter[lastModKey] = lastModDate
										expLastMod = lastModDate

										if expDate.IsZero() {
											expDate = lastModDate
										}
									}

									if expiryDateDate != sentinel {
										d.Frontmatter[expiryDateKey] = expiryDateDate
										expExiryDate = expiryDateDate
									}

									if expLastMod.IsZero() {
										expLastMod = expDate
									}

									assert.NoError(handler.HandleDates(d))

									assertFrontMatterDate(assert, d, expDate, "date")
									assertFrontMatterDate(assert, d, expLastMod, "lastmod")
									assertFrontMatterDate(assert, d, expPubDate, "publishdate")
									assertFrontMatterDate(assert, d, expExiryDate, "expirydate")
								}
							}
						}
					}
				}
			}
		}
	}
}

func assertFrontMatterDate(assert *require.Assertions, d FrontMatterDescriptor, expected time.Time, dateField string) {
	switch dateField {
	case "date":
	case "lastmod":
	case "publishdate":
	case "expirydate":
	default:
		assert.Failf("Unknown datefield %s", dateField)
	}

	param, found := d.Params[dateField]

	if found && param.(time.Time).IsZero() {
		assert.Fail("Zero time in params", dateField)
	}

	message := fmt.Sprintf("[%s] Found: %t Expected: %v (%t) Param: %v Params: %v Front matter: %v",
		dateField, found, expected, expected.IsZero(), param, d.Params, d.Frontmatter)

	assert.True(found != expected.IsZero(), message)

	if found {
		if expected != param {
			assert.Fail("Params check failed", "[%s] Expected:\n%q\nGot:\n%q", dateField, expected, param)
		}
	}
}

func TestFrontMatterFieldHandlers(t *testing.T) {
	//handlers := &frontmatterFieldHandlers{}

}
