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

	fc, err := newFrontmatterConfig(newWarningLogger(), viper.New())
	assert.NoError(err)

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
		gotDate, gotSlug := fc.dateAndSlugFromBaseFilename(test.name)

		assert.Equal(expectedDate, gotDate, errMsg)
		assert.Equal(test.slug, gotSlug, errMsg)

	}
}
