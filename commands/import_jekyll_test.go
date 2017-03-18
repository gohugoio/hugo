// Copyright 2015 The Hugo Authors. All rights reserved.
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

package commands

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseJekyllFilename(t *testing.T) {
	filenameArray := []string{
		"2015-01-02-test.md",
		"2012-03-15-中文.markup",
	}

	expectResult := []struct {
		postDate time.Time
		postName string
	}{
		{time.Date(2015, time.January, 2, 0, 0, 0, 0, time.UTC), "test"},
		{time.Date(2012, time.March, 15, 0, 0, 0, 0, time.UTC), "中文"},
	}

	for i, filename := range filenameArray {
		postDate, postName, err := parseJekyllFilename(filename)
		assert.Equal(t, err, nil)
		assert.Equal(t, expectResult[i].postDate.Format("2006-01-02"), postDate.Format("2006-01-02"))
		assert.Equal(t, expectResult[i].postName, postName)
	}
}

func TestConvertJekyllMetadata(t *testing.T) {
	testDataList := []struct {
		metadata interface{}
		postName string
		postDate time.Time
		draft    bool
		expect   string
	}{
		{map[interface{}]interface{}{}, "testPost", time.Date(2015, 10, 1, 0, 0, 0, 0, time.UTC), false,
			`{"date":"2015-10-01T00:00:00Z","url":"/2015/10/01/testPost/"}`},
		{map[interface{}]interface{}{}, "testPost", time.Date(2015, 10, 1, 0, 0, 0, 0, time.UTC), true,
			`{"date":"2015-10-01T00:00:00Z","draft":true,"url":"/2015/10/01/testPost/"}`},
		{map[interface{}]interface{}{"Permalink": "/permalink.html", "layout": "post"},
			"testPost", time.Date(2015, 10, 1, 0, 0, 0, 0, time.UTC), false,
			`{"date":"2015-10-01T00:00:00Z","url":"/permalink.html"}`},
		{map[interface{}]interface{}{"permalink": "/permalink.html"},
			"testPost", time.Date(2015, 10, 1, 0, 0, 0, 0, time.UTC), false,
			`{"date":"2015-10-01T00:00:00Z","url":"/permalink.html"}`},
		{map[interface{}]interface{}{"category": nil, "permalink": 123},
			"testPost", time.Date(2015, 10, 1, 0, 0, 0, 0, time.UTC), false,
			`{"date":"2015-10-01T00:00:00Z","url":"/2015/10/01/testPost/"}`},
		{map[interface{}]interface{}{"Excerpt_Separator": "sep"},
			"testPost", time.Date(2015, 10, 1, 0, 0, 0, 0, time.UTC), false,
			`{"date":"2015-10-01T00:00:00Z","excerpt_separator":"sep","url":"/2015/10/01/testPost/"}`},
		{map[interface{}]interface{}{"category": "book", "layout": "post", "Others": "Goods", "Date": "2015-10-01 12:13:11"},
			"testPost", time.Date(2015, 10, 1, 0, 0, 0, 0, time.UTC), false,
			`{"Others":"Goods","categories":["book"],"date":"2015-10-01T12:13:11Z","url":"/2015/10/01/testPost/"}`},
	}

	for _, data := range testDataList {
		result, err := convertJekyllMetaData(data.metadata, data.postName, data.postDate, data.draft)
		assert.Equal(t, nil, err)
		jsonResult, err := json.Marshal(result)
		assert.Equal(t, nil, err)
		assert.Equal(t, data.expect, string(jsonResult))
	}
}

func TestConvertJekyllContent(t *testing.T) {
	testDataList := []struct {
		metadata interface{}
		content  string
		expect   string
	}{
		{map[interface{}]interface{}{},
			`Test content\n<!-- more -->\npart2 content`, `Test content\n<!--more-->\npart2 content`},
		{map[interface{}]interface{}{},
			`Test content\n<!-- More -->\npart2 content`, `Test content\n<!--more-->\npart2 content`},
		{map[interface{}]interface{}{"excerpt_separator": "<!--sep-->"},
			`Test content\n<!--sep-->\npart2 content`, `Test content\n<!--more-->\npart2 content`},
		{map[interface{}]interface{}{}, "{% raw %}text{% endraw %}", "text"},
		{map[interface{}]interface{}{}, "{%raw%} text2 {%endraw %}", "text2"},
		{map[interface{}]interface{}{},
			"{% highlight go %}\nvar s int\n{% endhighlight %}",
			"{{< highlight go >}}\nvar s int\n{{< / highlight >}}"},
	}

	for _, data := range testDataList {
		result := convertJekyllContent(data.metadata, data.content)
		assert.Equal(t, data.expect, result)
	}
}
