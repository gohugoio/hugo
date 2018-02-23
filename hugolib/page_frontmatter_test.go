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

package hugolib

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

func TestFrontMatterDates(t *testing.T) {
	t.Parallel()

	defaultDateSettings := []string{"none", "file"}

	for _, defaultDateSetting := range defaultDateSettings {
		t.Run(fmt.Sprintf("defaultDate=%s", defaultDateSetting), func(t *testing.T) {
			doTestFrontMatterDates(t, defaultDateSetting)
		})
	}
}

func doTestFrontMatterDates(t *testing.T, defaultDateSetting string) {
	assert := require.New(t)

	cfg := viper.New()

	cfg.Set("frontmatter", map[string]interface{}{
		"defaultDate": []string{defaultDateSetting},
	})

	handler, err := newFrontmatterHandler(newWarningLogger(), cfg)
	assert.NoError(err)

	testDate, err := time.Parse("2006-01-02", "2018-02-01")
	assert.NoError(err)

	sentinel := (time.Time{}).Add(1 * time.Hour)
	zero := time.Time{}

	// See http://www.imdb.com/title/tt0133093/
	for _, lastModKey := range []string{"lastmod", "modified"} {
		for _, lastModDate := range []time.Time{testDate, sentinel} {
			for _, dateKey := range []string{"date"} {
				testDate = testDate.Add(24 * time.Hour)
				for _, dateDate := range []time.Time{testDate, sentinel} {
					for _, pubDateKey := range []string{"publishdate", "pubdate", "published"} {
						testDate = testDate.Add(24 * time.Hour)
						for _, pubDateDate := range []time.Time{testDate, sentinel} {
							for _, expiryDateKey := range []string{"expirydate", "unpublishdate"} {
								testDate = testDate.Add(24 * time.Hour)
								for _, expiryDateDate := range []time.Time{testDate, sentinel} {
									d := frontMatterDescriptor{
										frontmatter: make(map[string]interface{}),
										params:      make(map[string]interface{}),
										dates:       &PageDates{},
										pageURLs:    &URLPath{},
									}

									var (
										//	expLastModP, expDateP, expPubDateP, expExiryDateP = sentinel, sentinel, sentinel, sentinel
										expLastMod, expDate, expPubDate, expExiryDate = zero, zero, zero, zero
									)

									if lastModDate != sentinel {
										d.frontmatter[lastModKey] = lastModDate
										expLastMod = lastModDate
										expDate = lastModDate
									}

									if dateDate != sentinel {
										d.frontmatter[dateKey] = dateDate
										expDate = dateDate
									}

									if pubDateDate != sentinel {
										d.frontmatter[pubDateKey] = pubDateDate
										expPubDate = pubDateDate
									}

									if expiryDateDate != sentinel {
										d.frontmatter[expiryDateKey] = expiryDateDate
										expExiryDate = expiryDateDate
									}

									assert.NoError(handler.handleDates(d))

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

func assertFrontMatterDate(assert *require.Assertions, d frontMatterDescriptor, expected time.Time, dateField string) {
	switch dateField {
	case "date":
	case "lastmod":
	case "publishdate":
	case "expirydate":
	default:
		assert.Failf("Unknown datefield %s", dateField)
	}

	param, found := d.params[dateField]

	message := fmt.Sprintf("[%s] Found: %t Expected: %v Params: %v Front matter: %v", dateField, found, expected, d.params, d.frontmatter)

	assert.True(found != expected.IsZero(), message)

	if found {
		if expected != param {
			assert.Fail("Params check failed", "[%s] Expected:\n%q\nGot:\n%q", dateField, expected, param)
		}
		assert.Equal(expected, param)

	}
}

type dateTestHelper struct {
	name string

	dates PageDates
}

func (d dateTestHelper) descriptor() frontMatterDescriptor {
	return frontMatterDescriptor{dates: &d.dates}
}

func (d dateTestHelper) assert(t *testing.T) {

}
