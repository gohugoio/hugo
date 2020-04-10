// Copyright 2020 The Hugo Authors. All rights reserved.
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

package postcss

import (
	"regexp"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/htesting/hqt"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/helpers"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

// Issue 6166
func TestDecodeOptions(t *testing.T) {
	c := qt.New(t)
	opts1, err := DecodeOptions(map[string]interface{}{
		"no-map": true,
	})

	c.Assert(err, qt.IsNil)
	c.Assert(opts1.NoMap, qt.Equals, true)

	opts2, err := DecodeOptions(map[string]interface{}{
		"noMap": true,
	})

	c.Assert(err, qt.IsNil)
	c.Assert(opts2.NoMap, qt.Equals, true)

}

func TestShouldImport(t *testing.T) {
	c := qt.New(t)
	var imp *importResolver

	for _, test := range []struct {
		input  string
		expect bool
	}{
		{input: `@import "navigation.css";`, expect: true},
		{input: `@import "navigation.css"; /* Using a string */`, expect: true},
		{input: `@import "navigation.css"`, expect: true},
		{input: `@import 'navigation.css';`, expect: true},
		{input: `@import url("navigation.css");`, expect: false},
		{input: `@import url('https://fonts.googleapis.com/css?family=Open+Sans:400,400i,800,800i&display=swap');`, expect: false},
	} {
		c.Assert(imp.shouldImport(test.input), qt.Equals, test.expect)
	}
}

func TestImportResolver(t *testing.T) {
	c := qt.New(t)
	fs := afero.NewMemMapFs()

	writeFile := func(name, content string) {
		c.Assert(afero.WriteFile(fs, name, []byte(content), 0777), qt.IsNil)
	}

	writeFile("a.css", `@import "b.css";
@import "c.css";
A_STYLE1
A_STYLE2
`)

	writeFile("b.css", `B_STYLE`)
	writeFile("c.css", "@import \"d.css\"\nC_STYLE")
	writeFile("d.css", "@import \"a.css\"\n\nD_STYLE")
	writeFile("e.css", "E_STYLE")

	mainStyles := strings.NewReader(`@import "a.css";
@import "b.css";
LOCAL_STYLE
@import "c.css";
@import "e.css";
@import "missing.css";`)

	imp := newImportResolver(
		mainStyles,
		"styles.css",
		fs, loggers.NewErrorLogger(),
	)

	r, err := imp.resolve()
	c.Assert(err, qt.IsNil)
	rs := helpers.ReaderToString(r)
	result := regexp.MustCompile(`\n+`).ReplaceAllString(rs, "\n")

	c.Assert(result, hqt.IsSameString, `B_STYLE
D_STYLE
C_STYLE
A_STYLE1
A_STYLE2
LOCAL_STYLE
E_STYLE
@import "missing.css";`)

	dline := imp.linemap[3]
	c.Assert(dline, qt.DeepEquals, fileOffset{
		Offset:   1,
		Filename: "d.css",
	})

}

func BenchmarkImportResolver(b *testing.B) {
	c := qt.New(b)
	fs := afero.NewMemMapFs()

	writeFile := func(name, content string) {
		c.Assert(afero.WriteFile(fs, name, []byte(content), 0777), qt.IsNil)
	}

	writeFile("a.css", `@import "b.css";
@import "c.css";
A_STYLE1
A_STYLE2
`)

	writeFile("b.css", `B_STYLE`)
	writeFile("c.css", "@import \"d.css\"\nC_STYLE"+strings.Repeat("\nSTYLE", 12))
	writeFile("d.css", "@import \"a.css\"\n\nD_STYLE"+strings.Repeat("\nSTYLE", 55))
	writeFile("e.css", "E_STYLE")

	mainStyles := `@import "a.css";
@import "b.css";
LOCAL_STYLE
@import "c.css";
@import "e.css";
@import "missing.css";`

	logger := loggers.NewErrorLogger()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		imp := newImportResolver(
			strings.NewReader(mainStyles),
			"styles.css",
			fs, logger,
		)

		b.StartTimer()

		_, err := imp.resolve()
		if err != nil {
			b.Fatal(err)
		}

	}
}
