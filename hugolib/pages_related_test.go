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
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/deps"

	"github.com/stretchr/testify/require"
)

func TestRelated(t *testing.T) {
	assert := require.New(t)

	t.Parallel()

	var (
		cfg, fs = newTestCfg()
		//th      = testHelper{cfg, fs, t}
	)

	pageTmpl := `---
title: Page %d
keywords: [%s]
date: %s
---

Content
`

	writeSource(t, fs, filepath.Join("content", "page1.md"), fmt.Sprintf(pageTmpl, 1, "hugo, says", "2017-01-03"))
	writeSource(t, fs, filepath.Join("content", "page2.md"), fmt.Sprintf(pageTmpl, 2, "hugo, rocks", "2017-01-02"))
	writeSource(t, fs, filepath.Join("content", "page3.md"), fmt.Sprintf(pageTmpl, 3, "bep, says", "2017-01-01"))

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})
	assert.Len(s.RegularPages, 3)

	result, err := s.RegularPages.RelatedTo(types.NewKeyValuesStrings("keywords", "hugo", "rocks"))

	assert.NoError(err)
	assert.Len(result, 2)
	assert.Equal("Page 2", result[0].Title)
	assert.Equal("Page 1", result[1].Title)

	result, err = s.RegularPages.Related(s.RegularPages[0])
	assert.Len(result, 2)
	assert.Equal("Page 2", result[0].Title)
	assert.Equal("Page 3", result[1].Title)

	result, err = s.RegularPages.RelatedIndices(s.RegularPages[0], "keywords")
	assert.Len(result, 2)
	assert.Equal("Page 2", result[0].Title)
	assert.Equal("Page 3", result[1].Title)

	result, err = s.RegularPages.RelatedTo(types.NewKeyValuesStrings("keywords", "bep", "rocks"))
	assert.NoError(err)
	assert.Len(result, 2)
	assert.Equal("Page 2", result[0].Title)
	assert.Equal("Page 3", result[1].Title)
}
