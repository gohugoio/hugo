// Copyright 2016 The Hugo Authors. All rights reserved.
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

package tplimpl

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/spf13/hugo/tpl"

	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/helpers"

	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/spf13/hugo/config"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/i18n"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	logger = jww.NewNotepad(jww.LevelFatal, jww.LevelFatal, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
)

func newDepsConfig(cfg config.Provider) deps.DepsCfg {
	l := helpers.NewLanguage("en", cfg)
	l.Set("i18nDir", "i18n")
	return deps.DepsCfg{
		Language:            l,
		Cfg:                 cfg,
		Fs:                  hugofs.NewMem(l),
		Logger:              logger,
		TemplateProvider:    DefaultTemplateProvider,
		TranslationProvider: i18n.NewTranslationProvider(),
	}
}

type tstNoStringer struct {
}

type tstCompareType int

const (
	tstEq tstCompareType = iota
	tstNe
	tstGt
	tstGe
	tstLt
	tstLe
)

func tstIsEq(tp tstCompareType) bool {
	return tp == tstEq || tp == tstGe || tp == tstLe
}

func tstIsGt(tp tstCompareType) bool {
	return tp == tstGt || tp == tstGe
}

func tstIsLt(tp tstCompareType) bool {
	return tp == tstLt || tp == tstLe
}

func TestFuncsInTemplate(t *testing.T) {
	t.Parallel()

	workingDir := "/home/hugo"

	v := viper.New()

	v.Set("workingDir", workingDir)
	v.Set("multilingual", true)

	fs := hugofs.NewMem(v)

	afero.WriteFile(fs.Source, filepath.Join(workingDir, "README.txt"), []byte("Hugo Rocks!"), 0755)

	// Add the examples from the docs: As a smoke test and to make sure the examples work.
	// TODO(bep): docs: fix title example
	in :=
		`absLangURL: {{ "index.html" | absLangURL }}
absURL: {{ "http://gohugo.io/" | absURL }}
absURL: {{ "mystyle.css" | absURL }}
absURL: {{ 42 | absURL }}
add: {{add 1 2}}
base64Decode 1: {{ "SGVsbG8gd29ybGQ=" | base64Decode }}
base64Decode 2: {{ 42 | base64Encode | base64Decode }}
base64Encode: {{ "Hello world" | base64Encode }}
chomp: {{chomp "<p>Blockhead</p>\n" }}
dateFormat: {{ dateFormat "Monday, Jan 2, 2006" "2015-01-21" }}
delimit: {{ delimit (slice "A" "B" "C") ", " " and " }}
div: {{div 6 3}}
echoParam: {{ echoParam .Params "langCode" }}
emojify: {{ "I :heart: Hugo" | emojify }}
eq: {{ if eq .Section "blog" }}current{{ end }}
findRE: {{ findRE "[G|g]o" "Hugo is a static side generator written in Go." "1" }}
hasPrefix 1: {{ hasPrefix "Hugo" "Hu" }}
hasPrefix 2: {{ hasPrefix "Hugo" "Fu" }}
htmlEscape 1: {{ htmlEscape "Cathal Garvey & The Sunshine Band <cathal@foo.bar>" | safeHTML}}
htmlEscape 2: {{ htmlEscape "Cathal Garvey & The Sunshine Band <cathal@foo.bar>"}}
htmlUnescape 1: {{htmlUnescape "Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;" | safeHTML}}
htmlUnescape 2: {{"Cathal Garvey &amp;amp; The Sunshine Band &amp;lt;cathal@foo.bar&amp;gt;" | htmlUnescape | htmlUnescape | safeHTML}}
htmlUnescape 3: {{"Cathal Garvey &amp;amp; The Sunshine Band &amp;lt;cathal@foo.bar&amp;gt;" | htmlUnescape | htmlUnescape }}
htmlUnescape 4: {{ htmlEscape "Cathal Garvey & The Sunshine Band <cathal@foo.bar>" | htmlUnescape | safeHTML }}
htmlUnescape 5: {{ htmlUnescape "Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;" | htmlEscape | safeHTML }}
humanize 1: {{ humanize "my-first-post" }}
humanize 2: {{ humanize "myCamelPost" }}
humanize 3: {{ humanize "52" }}
humanize 4: {{ humanize 103 }}
in: {{ if in "this string contains a substring" "substring" }}Substring found!{{ end }}
jsonify: {{ (slice "A" "B" "C") | jsonify }}
lower: {{lower "BatMan"}}
markdownify: {{ .Title | markdownify}}
md5: {{ md5 "Hello world, gophers!" }}
mod: {{mod 15 3}}
modBool: {{modBool 15 3}}
mul: {{mul 2 3}}
print: {{ print "works!" }}
printf: {{ printf "%s!" "works" }}
println: {{ println "works!" -}}
plainify: {{ plainify  "Hello <strong>world</strong>, gophers!" }}
pluralize: {{ "cat" | pluralize }}
querify 1: {{ (querify "foo" 1 "bar" 2 "baz" "with spaces" "qux" "this&that=those") | safeHTML }}
querify 2: <a href="https://www.google.com?{{ (querify "q" "test" "page" 3) | safeURL }}">Search</a>
readDir: {{ range (readDir ".") }}{{ .Name }}{{ end }}
readFile: {{ readFile "README.txt" }}
relLangURL: {{ "index.html" | relLangURL }}
relURL 1: {{ "http://gohugo.io/" | relURL }}
relURL 2: {{ "mystyle.css" | relURL }}
relURL 3: {{ mul 2 21 | relURL }}
replace: {{ replace "Batman and Robin" "Robin" "Catwoman" }}
replaceRE: {{ "http://gohugo.io/docs" | replaceRE "^https?://([^/]+).*" "$1" }}
safeCSS: {{ "Bat&Man" | safeCSS | safeCSS }}
safeHTML: {{ "Bat&Man" | safeHTML | safeHTML }}
safeHTML: {{ "Bat&Man" | safeHTML }}
safeJS: {{ "(1*2)" | safeJS | safeJS }}
safeURL: {{ "http://gohugo.io" | safeURL | safeURL }}
seq: {{ seq 3 }}
sha1: {{ sha1 "Hello world, gophers!" }}
sha256: {{ sha256 "Hello world, gophers!" }}
singularize: {{ "cats" | singularize }}
slicestr: {{slicestr "BatMan" 0 3}}
slicestr: {{slicestr "BatMan" 3}}
sort: {{ slice "B" "C" "A" | sort }}
sub: {{sub 3 2}}
substr: {{substr "BatMan" 0 -3}}
substr: {{substr "BatMan" 3 3}}
title: {{title "Bat man"}}
time: {{ (time "2015-01-21").Year }}
trim: {{ trim "++Batman--" "+-" }}
truncate: {{ "this is a very long text" | truncate 10 " ..." }}
truncate: {{ "With [Markdown](/markdown) inside." | markdownify | truncate 14 }}
upper: {{upper "BatMan"}}
union: {{ union (slice 1 2 3) (slice 3 4 5) }}
urlize: {{ "Bat Man" | urlize }}
`

	expected := `absLangURL: http://mysite.com/hugo/en/index.html
absURL: http://gohugo.io/
absURL: http://mysite.com/hugo/mystyle.css
absURL: http://mysite.com/hugo/42
add: 3
base64Decode 1: Hello world
base64Decode 2: 42
base64Encode: SGVsbG8gd29ybGQ=
chomp: <p>Blockhead</p>
dateFormat: Wednesday, Jan 21, 2015
delimit: A, B and C
div: 2
echoParam: en
emojify: I ❤️ Hugo
eq: current
findRE: [go]
hasPrefix 1: true
hasPrefix 2: false
htmlEscape 1: Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;
htmlEscape 2: Cathal Garvey &amp;amp; The Sunshine Band &amp;lt;cathal@foo.bar&amp;gt;
htmlUnescape 1: Cathal Garvey & The Sunshine Band <cathal@foo.bar>
htmlUnescape 2: Cathal Garvey & The Sunshine Band <cathal@foo.bar>
htmlUnescape 3: Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;
htmlUnescape 4: Cathal Garvey & The Sunshine Band <cathal@foo.bar>
htmlUnescape 5: Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;
humanize 1: My first post
humanize 2: My camel post
humanize 3: 52nd
humanize 4: 103rd
in: Substring found!
jsonify: ["A","B","C"]
lower: batman
markdownify: <strong>BatMan</strong>
md5: b3029f756f98f79e7f1b7f1d1f0dd53b
mod: 0
modBool: true
mul: 6
print: works!
printf: works!
println: works!
plainify: Hello world, gophers!
pluralize: cats
querify 1: bar=2&baz=with+spaces&foo=1&qux=this%26that%3Dthose
querify 2: <a href="https://www.google.com?page=3&amp;q=test">Search</a>
readDir: README.txt
readFile: Hugo Rocks!
relLangURL: /hugo/en/index.html
relURL 1: http://gohugo.io/
relURL 2: /hugo/mystyle.css
relURL 3: /hugo/42
replace: Batman and Catwoman
replaceRE: gohugo.io
safeCSS: Bat&amp;Man
safeHTML: Bat&Man
safeHTML: Bat&Man
safeJS: (1*2)
safeURL: http://gohugo.io
seq: [1 2 3]
sha1: c8b5b0e33d408246e30f53e32b8f7627a7a649d4
sha256: 6ec43b78da9669f50e4e422575c54bf87536954ccd58280219c393f2ce352b46
singularize: cat
slicestr: Bat
slicestr: Man
sort: [A B C]
sub: 1
substr: Bat
substr: Man
title: Bat Man
time: 2015
trim: Batman
truncate: this is a ...
truncate: With <a href="/markdown">Markdown …</a>
upper: BATMAN
union: [1 2 3 4 5]
urlize: bat-man
`

	var b bytes.Buffer

	var data struct {
		Title   string
		Section string
		Params  map[string]interface{}
	}

	data.Title = "**BatMan**"
	data.Section = "blog"
	data.Params = map[string]interface{}{"langCode": "en"}

	v.Set("baseURL", "http://mysite.com/hugo/")
	v.Set("CurrentContentLanguage", helpers.NewLanguage("en", v))

	config := newDepsConfig(v)
	config.WithTemplate = func(templ tpl.TemplateHandler) error {
		if err := templ.AddTemplate("test", in); err != nil {
			t.Fatal("Got error on parse", err)
		}
		return nil
	}
	config.Fs = fs

	d, err := deps.New(config)
	if err != nil {
		t.Fatal(err)
	}

	if err := d.LoadResources(); err != nil {
		t.Fatal(err)
	}

	err = d.Tmpl.Lookup("test").Execute(&b, &data)

	if err != nil {
		t.Fatal("Got error on execute", err)
	}

	if b.String() != expected {
		sl1 := strings.Split(b.String(), "\n")
		sl2 := strings.Split(expected, "\n")
		t.Errorf("Diff:\n%q", helpers.DiffStringSlices(sl1, sl2))
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()
	for _, this := range []struct {
		tstCompareType
		funcUnderTest func(a, b interface{}) bool
	}{
		{tstGt, gt},
		{tstLt, lt},
		{tstGe, ge},
		{tstLe, le},
		{tstEq, eq},
		{tstNe, ne},
	} {
		doTestCompare(t, this.tstCompareType, this.funcUnderTest)
	}
}

func doTestCompare(t *testing.T, tp tstCompareType, funcUnderTest func(a, b interface{}) bool) {
	for i, this := range []struct {
		left            interface{}
		right           interface{}
		expectIndicator int
	}{
		{5, 8, -1},
		{8, 5, 1},
		{5, 5, 0},
		{int(5), int64(5), 0},
		{int32(5), int(5), 0},
		{int16(4), int(5), -1},
		{uint(15), uint64(15), 0},
		{-2, 1, -1},
		{2, -5, 1},
		{0.0, 1.23, -1},
		{1.1, 1.1, 0},
		{float32(1.0), float64(1.0), 0},
		{1.23, 0.0, 1},
		{"5", "5", 0},
		{"8", "5", 1},
		{"5", "0001", 1},
		{[]int{100, 99}, []int{1, 2, 3, 4}, -1},
		{cast.ToTime("2015-11-20"), cast.ToTime("2015-11-20"), 0},
		{cast.ToTime("2015-11-19"), cast.ToTime("2015-11-20"), -1},
		{cast.ToTime("2015-11-20"), cast.ToTime("2015-11-19"), 1},
	} {
		result := funcUnderTest(this.left, this.right)
		success := false

		if this.expectIndicator == 0 {
			if tstIsEq(tp) {
				success = result
			} else {
				success = !result
			}
		}

		if this.expectIndicator < 0 {
			success = result && (tstIsLt(tp) || tp == tstNe)
			success = success || (!result && !tstIsLt(tp))
		}

		if this.expectIndicator > 0 {
			success = result && (tstIsGt(tp) || tp == tstNe)
			success = success || (!result && (!tstIsGt(tp) || tp != tstNe))
		}

		if !success {
			t.Errorf("[%d][%s] %v compared to %v: %t", i, path.Base(runtime.FuncForPC(reflect.ValueOf(funcUnderTest).Pointer()).Name()), this.left, this.right, result)
		}
	}
}

func TestMod(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		a      interface{}
		b      interface{}
		expect interface{}
	}{
		{3, 2, int64(1)},
		{3, 1, int64(0)},
		{3, 0, false},
		{0, 3, int64(0)},
		{3.1, 2, false},
		{3, 2.1, false},
		{3.1, 2.1, false},
		{int8(3), int8(2), int64(1)},
		{int16(3), int16(2), int64(1)},
		{int32(3), int32(2), int64(1)},
		{int64(3), int64(2), int64(1)},
	} {
		result, err := mod(this.a, this.b)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] modulo didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] modulo got %v but expected %v", i, result, this.expect)
			}
		}
	}
}

func TestModBool(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		a      interface{}
		b      interface{}
		expect interface{}
	}{
		{3, 3, true},
		{3, 2, false},
		{3, 1, true},
		{3, 0, nil},
		{0, 3, true},
		{3.1, 2, nil},
		{3, 2.1, nil},
		{3.1, 2.1, nil},
		{int8(3), int8(3), true},
		{int8(3), int8(2), false},
		{int16(3), int16(3), true},
		{int16(3), int16(2), false},
		{int32(3), int32(3), true},
		{int32(3), int32(2), false},
		{int64(3), int64(3), true},
		{int64(3), int64(2), false},
	} {
		result, err := modBool(this.a, this.b)
		if this.expect == nil {
			if err == nil {
				t.Errorf("[%d] modulo didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] modulo got %v but expected %v", i, result, this.expect)
			}
		}
	}
}

func TestFirst(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		count    interface{}
		sequence interface{}
		expect   interface{}
	}{
		{int(2), []string{"a", "b", "c"}, []string{"a", "b"}},
		{int32(3), []string{"a", "b"}, []string{"a", "b"}},
		{int64(2), []int{100, 200, 300}, []int{100, 200}},
		{100, []int{100, 200}, []int{100, 200}},
		{"1", []int{100, 200, 300}, []int{100}},
		{int64(-1), []int{100, 200, 300}, false},
		{"noint", []int{100, 200, 300}, false},
		{1, nil, false},
		{nil, []int{100}, false},
		{1, t, false},
		{1, (*string)(nil), false},
	} {
		results, err := first(this.count, this.sequence)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] First didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, this.expect) {
				t.Errorf("[%d] First %d items, got %v but expected %v", i, this.count, results, this.expect)
			}
		}
	}
}

func TestLast(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		count    interface{}
		sequence interface{}
		expect   interface{}
	}{
		{int(2), []string{"a", "b", "c"}, []string{"b", "c"}},
		{int32(3), []string{"a", "b"}, []string{"a", "b"}},
		{int64(2), []int{100, 200, 300}, []int{200, 300}},
		{100, []int{100, 200}, []int{100, 200}},
		{"1", []int{100, 200, 300}, []int{300}},
		{int64(-1), []int{100, 200, 300}, false},
		{"noint", []int{100, 200, 300}, false},
		{1, nil, false},
		{nil, []int{100}, false},
		{1, t, false},
		{1, (*string)(nil), false},
	} {
		results, err := last(this.count, this.sequence)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] First didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, this.expect) {
				t.Errorf("[%d] First %d items, got %v but expected %v", i, this.count, results, this.expect)
			}
		}
	}
}

func TestAfter(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		count    interface{}
		sequence interface{}
		expect   interface{}
	}{
		{int(2), []string{"a", "b", "c", "d"}, []string{"c", "d"}},
		{int32(3), []string{"a", "b"}, false},
		{int64(2), []int{100, 200, 300}, []int{300}},
		{100, []int{100, 200}, false},
		{"1", []int{100, 200, 300}, []int{200, 300}},
		{int64(-1), []int{100, 200, 300}, false},
		{"noint", []int{100, 200, 300}, false},
		{1, nil, false},
		{nil, []int{100}, false},
		{1, t, false},
		{1, (*string)(nil), false},
	} {
		results, err := after(this.count, this.sequence)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] First didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, this.expect) {
				t.Errorf("[%d] First %d items, got %v but expected %v", i, this.count, results, this.expect)
			}
		}
	}
}

func TestShuffleInputAndOutputFormat(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		sequence interface{}
		success  bool
	}{
		{[]string{"a", "b", "c", "d"}, true},
		{[]int{100, 200, 300}, true},
		{[]int{100, 200, 300}, true},
		{[]int{100, 200}, true},
		{[]string{"a", "b"}, true},
		{[]int{100, 200, 300}, true},
		{[]int{100, 200, 300}, true},
		{[]int{100}, true},
		{nil, false},
		{t, false},
		{(*string)(nil), false},
	} {
		results, err := shuffle(this.sequence)
		if !this.success {
			if err == nil {
				t.Errorf("[%d] First didn't return an expected error", i)
			}
		} else {
			resultsv := reflect.ValueOf(results)
			sequencev := reflect.ValueOf(this.sequence)

			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}

			if resultsv.Len() != sequencev.Len() {
				t.Errorf("Expected %d items, got %d items", sequencev.Len(), resultsv.Len())
			}
		}
	}
}

func TestShuffleRandomising(t *testing.T) {
	t.Parallel()
	// Note that this test can fail with false negative result if the shuffle
	// of the sequence happens to be the same as the original sequence. However
	// the propability of the event is 10^-158 which is negligible.
	sequenceLength := 100
	rand.Seed(time.Now().UTC().UnixNano())

	for _, this := range []struct {
		sequence []int
	}{
		{rand.Perm(sequenceLength)},
	} {
		results, _ := shuffle(this.sequence)

		resultsv := reflect.ValueOf(results)

		allSame := true
		for index, value := range this.sequence {
			allSame = allSame && (resultsv.Index(index).Interface() == value)
		}

		if allSame {
			t.Error("Expected sequence to be shuffled but was in the same order")
		}
	}
}

func TestDictionary(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		v1            []interface{}
		expecterr     bool
		expectedValue map[string]interface{}
	}{
		{[]interface{}{"a", "b"}, false, map[string]interface{}{"a": "b"}},
		{[]interface{}{5, "b"}, true, nil},
		{[]interface{}{"a", 12, "b", []int{4}}, false, map[string]interface{}{"a": 12, "b": []int{4}}},
		{[]interface{}{"a", "b", "c"}, true, nil},
	} {
		r, e := dictionary(this.v1...)

		if (this.expecterr && e == nil) || (!this.expecterr && e != nil) {
			t.Errorf("[%d] got an unexpected error: %s", i, e)
		} else if !this.expecterr {
			if !reflect.DeepEqual(r, this.expectedValue) {
				t.Errorf("[%d] got %v but expected %v", i, r, this.expectedValue)
			}
		}
	}
}

func blankImage(width, height int) []byte {
	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TestImageConfig(t *testing.T) {
	t.Parallel()

	workingDir := "/home/hugo"

	v := viper.New()

	v.Set("workingDir", workingDir)

	f := newTestFuncsterWithViper(v)

	for i, this := range []struct {
		path     string
		input    []byte
		expected image.Config
	}{
		{
			path:  "a.png",
			input: blankImage(10, 10),
			expected: image.Config{
				Width:      10,
				Height:     10,
				ColorModel: color.NRGBAModel,
			},
		},
		{
			path:  "a.png",
			input: blankImage(10, 10),
			expected: image.Config{
				Width:      10,
				Height:     10,
				ColorModel: color.NRGBAModel,
			},
		},
		{
			path:  "b.png",
			input: blankImage(20, 15),
			expected: image.Config{
				Width:      20,
				Height:     15,
				ColorModel: color.NRGBAModel,
			},
		},
		{
			path:  "a.png",
			input: blankImage(20, 15),
			expected: image.Config{
				Width:      10,
				Height:     10,
				ColorModel: color.NRGBAModel,
			},
		},
	} {
		afero.WriteFile(f.Fs.Source, filepath.Join(workingDir, this.path), this.input, 0755)

		result, err := f.image.config(this.path)
		if err != nil {
			t.Errorf("imageConfig returned error: %s", err)
		}

		if !reflect.DeepEqual(result, this.expected) {
			t.Errorf("[%d] imageConfig: expected '%v', got '%v'", i, this.expected, result)
		}

		if len(f.image.imageConfigCache) == 0 {
			t.Error("defaultImageConfigCache should have at least 1 item")
		}
	}

	if _, err := f.image.config(t); err == nil {
		t.Error("Expected error from imageConfig when passed invalid path")
	}

	if _, err := f.image.config("non-existent.png"); err == nil {
		t.Error("Expected error from imageConfig when passed non-existent file")
	}

	if _, err := f.image.config(""); err == nil {
		t.Error("Expected error from imageConfig when passed empty path")
	}

}

func TestIn(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		v1     interface{}
		v2     interface{}
		expect bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]interface{}{"a", "b", "c"}, "b", true},
		{[]interface{}{"a", "b", "c"}, "d", false},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "12", "c"}, 12, false},
		{[]int{1, 2, 4}, 2, true},
		{[]interface{}{1, 2, 4}, 2, true},
		{[]interface{}{1, 2, 4}, nil, false},
		{[]interface{}{nil}, nil, false},
		{[]int{1, 2, 4}, 3, false},
		{[]float64{1.23, 2.45, 4.67}, 1.23, true},
		{[]float64{1.234567, 2.45, 4.67}, 1.234568, false},
		{"this substring should be found", "substring", true},
		{"this substring should not be found", "subseastring", false},
	} {
		result := in(this.v1, this.v2)

		if result != this.expect {
			t.Errorf("[%d] got %v but expected %v", i, result, this.expect)
		}
	}
}

func TestSlicestr(t *testing.T) {
	t.Parallel()
	var err error
	for i, this := range []struct {
		v1     interface{}
		v2     interface{}
		v3     interface{}
		expect interface{}
	}{
		{"abc", 1, 2, "b"},
		{"abc", 1, 3, "bc"},
		{"abcdef", 1, int8(3), "bc"},
		{"abcdef", 1, int16(3), "bc"},
		{"abcdef", 1, int32(3), "bc"},
		{"abcdef", 1, int64(3), "bc"},
		{"abc", 0, 1, "a"},
		{"abcdef", nil, nil, "abcdef"},
		{"abcdef", 0, 6, "abcdef"},
		{"abcdef", 0, 2, "ab"},
		{"abcdef", 2, nil, "cdef"},
		{"abcdef", int8(2), nil, "cdef"},
		{"abcdef", int16(2), nil, "cdef"},
		{"abcdef", int32(2), nil, "cdef"},
		{"abcdef", int64(2), nil, "cdef"},
		{123, 1, 3, "23"},
		{"abcdef", 6, nil, false},
		{"abcdef", 4, 7, false},
		{"abcdef", -1, nil, false},
		{"abcdef", -1, 7, false},
		{"abcdef", 1, -1, false},
		{tstNoStringer{}, 0, 1, false},
		{"ĀĀĀ", 0, 1, "Ā"}, // issue #1333
		{"a", t, nil, false},
		{"a", 1, t, false},
	} {
		var result string
		if this.v2 == nil {
			result, err = slicestr(this.v1)
		} else if this.v3 == nil {
			result, err = slicestr(this.v1, this.v2)
		} else {
			result, err = slicestr(this.v1, this.v2, this.v3)
		}

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Slice didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
			}
		}
	}

	// Too many arguments
	_, err = slicestr("a", 1, 2, 3)
	if err == nil {
		t.Errorf("Should have errored")
	}
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()
	cases := []struct {
		s      interface{}
		prefix interface{}
		want   interface{}
		isErr  bool
	}{
		{"abcd", "ab", true, false},
		{"abcd", "cd", false, false},
		{template.HTML("abcd"), "ab", true, false},
		{template.HTML("abcd"), "cd", false, false},
		{template.HTML("1234"), 12, true, false},
		{template.HTML("1234"), 34, false, false},
		{[]byte("abcd"), "ab", true, false},
	}

	for i, c := range cases {
		res, err := hasPrefix(c.s, c.prefix)
		if (err != nil) != c.isErr {
			t.Fatalf("[%d] unexpected isErr state: want %v, got %v, err = %v", i, c.isErr, err != nil, err)
		}
		if res != c.want {
			t.Errorf("[%d] want %v, got %v", i, c.want, res)
		}
	}
}

func TestSubstr(t *testing.T) {
	t.Parallel()
	var err error
	var n int
	for i, this := range []struct {
		v1     interface{}
		v2     interface{}
		v3     interface{}
		expect interface{}
	}{
		{"abc", 1, 2, "bc"},
		{"abc", 0, 1, "a"},
		{"abcdef", -1, 2, "ef"},
		{"abcdef", -3, 3, "bcd"},
		{"abcdef", 0, -1, "abcde"},
		{"abcdef", 2, -1, "cde"},
		{"abcdef", 4, -4, false},
		{"abcdef", 7, 1, false},
		{"abcdef", 1, 100, "bcdef"},
		{"abcdef", -100, 3, "abc"},
		{"abcdef", -3, -1, "de"},
		{"abcdef", 2, nil, "cdef"},
		{"abcdef", int8(2), nil, "cdef"},
		{"abcdef", int16(2), nil, "cdef"},
		{"abcdef", int32(2), nil, "cdef"},
		{"abcdef", int64(2), nil, "cdef"},
		{"abcdef", 2, int8(3), "cde"},
		{"abcdef", 2, int16(3), "cde"},
		{"abcdef", 2, int32(3), "cde"},
		{"abcdef", 2, int64(3), "cde"},
		{123, 1, 3, "23"},
		{1.2e3, 0, 4, "1200"},
		{tstNoStringer{}, 0, 1, false},
		{"abcdef", 2.0, nil, "cdef"},
		{"abcdef", 2.0, 2, "cd"},
		{"abcdef", 2, 2.0, "cd"},
		{"ĀĀĀ", 1, 2, "ĀĀ"}, // # issue 1333
		{"abcdef", "doo", nil, false},
		{"abcdef", "doo", "doo", false},
		{"abcdef", 1, "doo", false},
	} {
		var result string
		n = i

		if this.v3 == nil {
			result, err = substr(this.v1, this.v2)
		} else {
			result, err = substr(this.v1, this.v2, this.v3)
		}

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Substr didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
			}
		}
	}

	n++
	_, err = substr("abcdef")
	if err == nil {
		t.Errorf("[%d] Substr didn't return an expected error", n)
	}

	n++
	_, err = substr("abcdef", 1, 2, 3)
	if err == nil {
		t.Errorf("[%d] Substr didn't return an expected error", n)
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		v1     interface{}
		v2     string
		expect interface{}
	}{
		{"a, b", ", ", []string{"a", "b"}},
		{"a & b & c", " & ", []string{"a", "b", "c"}},
		{"http://example.com", "http://", []string{"", "example.com"}},
		{123, "2", []string{"1", "3"}},
		{tstNoStringer{}, ",", false},
	} {
		result, err := split(this.v1, this.v2)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Split didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
			}
		}
	}
}

func TestIntersect(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		sequence1 interface{}
		sequence2 interface{}
		expect    interface{}
	}{
		{[]string{"a", "b", "c", "c"}, []string{"a", "b", "b"}, []string{"a", "b"}},
		{[]string{"a", "b"}, []string{"a", "b", "c"}, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, []string{"d", "e"}, []string{}},
		{[]string{}, []string{}, []string{}},
		{nil, nil, make([]interface{}, 0)},
		{[]string{"1", "2"}, []int{1, 2}, []string{}},
		{[]int{1, 2}, []string{"1", "2"}, []int{}},
		{[]int{1, 2, 4}, []int{2, 4}, []int{2, 4}},
		{[]int{2, 4}, []int{1, 2, 4}, []int{2, 4}},
		{[]int{1, 2, 4}, []int{3, 6}, []int{}},
		{[]float64{2.2, 4.4}, []float64{1.1, 2.2, 4.4}, []float64{2.2, 4.4}},
	} {
		results, err := intersect(this.sequence1, this.sequence2)
		if err != nil {
			t.Errorf("[%d] failed: %s", i, err)
			continue
		}
		if !reflect.DeepEqual(results, this.expect) {
			t.Errorf("[%d] got %v but expected %v", i, results, this.expect)
		}
	}

	_, err1 := intersect("not an array or slice", []string{"a"})

	if err1 == nil {
		t.Error("Expected error for non array as first arg")
	}

	_, err2 := intersect([]string{"a"}, "not an array or slice")

	if err2 == nil {
		t.Error("Expected error for non array as second arg")
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		sequence1 interface{}
		sequence2 interface{}
		expect    interface{}
		isErr     bool
	}{
		{[]string{"a", "b", "c", "c"}, []string{"a", "b", "b"}, []string{"a", "b", "c"}, false},
		{[]string{"a", "b"}, []string{"a", "b", "c"}, []string{"a", "b", "c"}, false},
		{[]string{"a", "b", "c"}, []string{"d", "e"}, []string{"a", "b", "c", "d", "e"}, false},
		{[]string{}, []string{}, []string{}, false},
		{[]string{"a", "b"}, nil, []string{"a", "b"}, false},
		{nil, []string{"a", "b"}, []string{"a", "b"}, false},
		{nil, nil, make([]interface{}, 0), true},
		{[]string{"1", "2"}, []int{1, 2}, make([]string, 0), false},
		{[]int{1, 2}, []string{"1", "2"}, make([]int, 0), false},
		{[]int{1, 2, 3}, []int{3, 4, 5}, []int{1, 2, 3, 4, 5}, false},
		{[]int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}, false},
		{[]int{1, 2, 4}, []int{2, 4}, []int{1, 2, 4}, false},
		{[]int{2, 4}, []int{1, 2, 4}, []int{2, 4, 1}, false},
		{[]int{1, 2, 4}, []int{3, 6}, []int{1, 2, 4, 3, 6}, false},
		{[]float64{2.2, 4.4}, []float64{1.1, 2.2, 4.4}, []float64{2.2, 4.4, 1.1}, false},
	} {
		results, err := union(this.sequence1, this.sequence2)
		if err != nil && !this.isErr {
			t.Errorf("[%d] failed: %s", i, err)
			continue
		}
		if !reflect.DeepEqual(results, this.expect) && !this.isErr {
			t.Errorf("[%d] got %v but expected %v", i, results, this.expect)
		}
	}

	_, err1 := union("not an array or slice", []string{"a"})

	if err1 == nil {
		t.Error("Expected error for non array as first arg")
	}

	_, err2 := union([]string{"a"}, "not an array or slice")

	if err2 == nil {
		t.Error("Expected error for non array as second arg")
	}
}

func TestIsSet(t *testing.T) {
	t.Parallel()
	aSlice := []interface{}{1, 2, 3, 5}
	aMap := map[string]interface{}{"a": 1, "b": 2}

	assert.True(t, isSet(aSlice, 2))
	assert.True(t, isSet(aMap, "b"))
	assert.False(t, isSet(aSlice, 22))
	assert.False(t, isSet(aMap, "bc"))
}

func (x *TstX) TstRp() string {
	return "r" + x.A
}

func (x TstX) TstRv() string {
	return "r" + x.B
}

func (x TstX) unexportedMethod() string {
	return x.unexported
}

func (x TstX) MethodWithArg(s string) string {
	return s
}

func (x TstX) MethodReturnNothing() {}

func (x TstX) MethodReturnErrorOnly() error {
	return errors.New("some error occurred")
}

func (x TstX) MethodReturnTwoValues() (string, string) {
	return "foo", "bar"
}

func (x TstX) MethodReturnValueWithError() (string, error) {
	return "", errors.New("some error occurred")
}

func (x TstX) String() string {
	return fmt.Sprintf("A: %s, B: %s", x.A, x.B)
}

type TstX struct {
	A, B       string
	unexported string
}

func TestTimeUnix(t *testing.T) {
	t.Parallel()
	var sec int64 = 1234567890
	tv := reflect.ValueOf(time.Unix(sec, 0))
	i := 1

	res := toTimeUnix(tv)
	if sec != res {
		t.Errorf("[%d] timeUnix got %v but expected %v", i, res, sec)
	}

	i++
	func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("[%d] timeUnix didn't return an expected error", i)
			}
		}()
		iv := reflect.ValueOf(sec)
		toTimeUnix(iv)
	}(t)
}

func TestEvaluateSubElem(t *testing.T) {
	t.Parallel()
	tstx := TstX{A: "foo", B: "bar"}
	var inner struct {
		S fmt.Stringer
	}
	inner.S = tstx
	interfaceValue := reflect.ValueOf(&inner).Elem().Field(0)

	for i, this := range []struct {
		value  reflect.Value
		key    string
		expect interface{}
	}{
		{reflect.ValueOf(tstx), "A", "foo"},
		{reflect.ValueOf(&tstx), "TstRp", "rfoo"},
		{reflect.ValueOf(tstx), "TstRv", "rbar"},
		//{reflect.ValueOf(map[int]string{1: "foo", 2: "bar"}), 1, "foo"},
		{reflect.ValueOf(map[string]string{"key1": "foo", "key2": "bar"}), "key1", "foo"},
		{interfaceValue, "String", "A: foo, B: bar"},
		{reflect.Value{}, "foo", false},
		//{reflect.ValueOf(map[int]string{1: "foo", 2: "bar"}), 1.2, false},
		{reflect.ValueOf(tstx), "unexported", false},
		{reflect.ValueOf(tstx), "unexportedMethod", false},
		{reflect.ValueOf(tstx), "MethodWithArg", false},
		{reflect.ValueOf(tstx), "MethodReturnNothing", false},
		{reflect.ValueOf(tstx), "MethodReturnErrorOnly", false},
		{reflect.ValueOf(tstx), "MethodReturnTwoValues", false},
		{reflect.ValueOf(tstx), "MethodReturnValueWithError", false},
		{reflect.ValueOf((*TstX)(nil)), "A", false},
		{reflect.ValueOf(tstx), "C", false},
		{reflect.ValueOf(map[int]string{1: "foo", 2: "bar"}), "1", false},
		{reflect.ValueOf([]string{"foo", "bar"}), "1", false},
	} {
		result, err := evaluateSubElem(this.value, this.key)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] evaluateSubElem didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if result.Kind() != reflect.String || result.String() != this.expect {
				t.Errorf("[%d] evaluateSubElem with %v got %v but expected %v", i, this.key, result, this.expect)
			}
		}
	}
}

func TestCheckCondition(t *testing.T) {
	t.Parallel()
	type expect struct {
		result  bool
		isError bool
	}

	for i, this := range []struct {
		value reflect.Value
		match reflect.Value
		op    string
		expect
	}{
		{reflect.ValueOf(123), reflect.ValueOf(123), "", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("foo"), "", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			"",
			expect{true, false},
		},
		{reflect.ValueOf(true), reflect.ValueOf(true), "", expect{true, false}},
		{reflect.ValueOf(nil), reflect.ValueOf(nil), "", expect{true, false}},
		{reflect.ValueOf(123), reflect.ValueOf(456), "!=", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar"), "!=", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			"!=",
			expect{true, false},
		},
		{reflect.ValueOf(true), reflect.ValueOf(false), "!=", expect{true, false}},
		{reflect.ValueOf(123), reflect.ValueOf(nil), "!=", expect{true, false}},
		{reflect.ValueOf(456), reflect.ValueOf(123), ">=", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar"), ">=", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			">=",
			expect{true, false},
		},
		{reflect.ValueOf(456), reflect.ValueOf(123), ">", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar"), ">", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			">",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf(456), "<=", expect{true, false}},
		{reflect.ValueOf("bar"), reflect.ValueOf("foo"), "<=", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			"<=",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf(456), "<", expect{true, false}},
		{reflect.ValueOf("bar"), reflect.ValueOf("foo"), "<", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			"<",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf([]int{123, 45, 678}), "in", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf([]string{"foo", "bar", "baz"}), "in", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf([]time.Time{
				time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.June, 26, 19, 18, 56, 12345, time.UTC),
			}),
			"in",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf([]int{45, 678}), "not in", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf([]string{"bar", "baz"}), "not in", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf([]time.Time{
				time.Date(2015, time.February, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.March, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC),
			}),
			"not in",
			expect{true, false},
		},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar-foo-baz"), "in", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar--baz"), "not in", expect{true, false}},
		{reflect.Value{}, reflect.ValueOf("foo"), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.Value{}, "", expect{false, false}},
		{reflect.ValueOf((*TstX)(nil)), reflect.ValueOf("foo"), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf((*TstX)(nil)), "", expect{false, false}},
		{reflect.ValueOf(true), reflect.ValueOf("foo"), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf(true), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf(map[int]string{}), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf([]int{1, 2}), "", expect{false, false}},
		{reflect.ValueOf((*TstX)(nil)), reflect.ValueOf((*TstX)(nil)), ">", expect{false, false}},
		{reflect.ValueOf(true), reflect.ValueOf(false), ">", expect{false, false}},
		{reflect.ValueOf(123), reflect.ValueOf([]int{}), "in", expect{false, false}},
		{reflect.ValueOf(123), reflect.ValueOf(123), "op", expect{false, true}},
	} {
		result, err := checkCondition(this.value, this.match, this.op)
		if this.expect.isError {
			if err == nil {
				t.Errorf("[%d] checkCondition didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if result != this.expect.result {
				t.Errorf("[%d] check condition %v %s %v, got %v but expected %v", i, this.value, this.op, this.match, result, this.expect.result)
			}
		}
	}
}

func TestWhere(t *testing.T) {
	t.Parallel()

	type Mid struct {
		Tst TstX
	}

	d1 := time.Now()
	d2 := d1.Add(1 * time.Hour)
	d3 := d2.Add(1 * time.Hour)
	d4 := d3.Add(1 * time.Hour)
	d5 := d4.Add(1 * time.Hour)
	d6 := d5.Add(1 * time.Hour)

	for i, this := range []struct {
		sequence interface{}
		key      interface{}
		op       string
		match    interface{}
		expect   interface{}
	}{
		{
			sequence: []map[int]string{
				{1: "a", 2: "m"}, {1: "c", 2: "d"}, {1: "e", 3: "m"},
			},
			key: 2, match: "m",
			expect: []map[int]string{
				{1: "a", 2: "m"},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4,
			expect: []map[string]int{
				{"a": 3, "b": 4},
			},
		},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", match: "f",
			expect: []TstX{
				{A: "e", B: "f"},
			},
		},
		{
			sequence: []*map[int]string{
				{1: "a", 2: "m"}, {1: "c", 2: "d"}, {1: "e", 3: "m"},
			},
			key: 2, match: "m",
			expect: []*map[int]string{
				{1: "a", 2: "m"},
			},
		},
		{
			sequence: []*TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", match: "f",
			expect: []*TstX{
				{A: "e", B: "f"},
			},
		},
		{
			sequence: []*TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "c"},
			},
			key: "TstRp", match: "rc",
			expect: []*TstX{
				{A: "c", B: "d"},
			},
		},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "c"},
			},
			key: "TstRv", match: "rc",
			expect: []TstX{
				{A: "e", B: "c"},
			},
		},
		{
			sequence: []map[string]TstX{
				{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: "foo.B", match: "d",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			sequence: []map[string]TstX{
				{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: ".foo.B", match: "d",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			sequence: []map[string]TstX{
				{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: "foo.TstRv", match: "rd",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			sequence: []map[string]*TstX{
				{"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}, {"foo": &TstX{A: "e", B: "f"}},
			},
			key: "foo.TstRp", match: "rc",
			expect: []map[string]*TstX{
				{"foo": &TstX{A: "c", B: "d"}},
			},
		},
		{
			sequence: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": Mid{Tst: TstX{A: "e", B: "f"}}},
			},
			key: "foo.Tst.B", match: "d",
			expect: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "c", B: "d"}}},
			},
		},
		{
			sequence: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": Mid{Tst: TstX{A: "e", B: "f"}}},
			},
			key: "foo.Tst.TstRv", match: "rd",
			expect: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "c", B: "d"}}},
			},
		},
		{
			sequence: []map[string]*Mid{
				{"foo": &Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": &Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": &Mid{Tst: TstX{A: "e", B: "f"}}},
			},
			key: "foo.Tst.TstRp", match: "rc",
			expect: []map[string]*Mid{
				{"foo": &Mid{Tst: TstX{A: "c", B: "d"}}},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: ">", match: 3,
			expect: []map[string]int{
				{"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
		},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "!=", match: "f",
			expect: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: "in", match: []int{3, 4, 5},
			expect: []map[string]int{
				{"a": 3, "b": 4},
			},
		},
		{
			sequence: []map[string][]string{
				{"a": []string{"A", "B", "C"}, "b": []string{"D", "E", "F"}}, {"a": []string{"G", "H", "I"}, "b": []string{"J", "K", "L"}}, {"a": []string{"M", "N", "O"}, "b": []string{"P", "Q", "R"}},
			},
			key: "b", op: "intersect", match: []string{"D", "P", "Q"},
			expect: []map[string][]string{
				{"a": []string{"A", "B", "C"}, "b": []string{"D", "E", "F"}}, {"a": []string{"M", "N", "O"}, "b": []string{"P", "Q", "R"}},
			},
		},
		{
			sequence: []map[string][]int{
				{"a": []int{1, 2, 3}, "b": []int{4, 5, 6}}, {"a": []int{7, 8, 9}, "b": []int{10, 11, 12}}, {"a": []int{13, 14, 15}, "b": []int{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int{4, 10, 12},
			expect: []map[string][]int{
				{"a": []int{1, 2, 3}, "b": []int{4, 5, 6}}, {"a": []int{7, 8, 9}, "b": []int{10, 11, 12}},
			},
		},
		{
			sequence: []map[string][]int8{
				{"a": []int8{1, 2, 3}, "b": []int8{4, 5, 6}}, {"a": []int8{7, 8, 9}, "b": []int8{10, 11, 12}}, {"a": []int8{13, 14, 15}, "b": []int8{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int8{4, 10, 12},
			expect: []map[string][]int8{
				{"a": []int8{1, 2, 3}, "b": []int8{4, 5, 6}}, {"a": []int8{7, 8, 9}, "b": []int8{10, 11, 12}},
			},
		},
		{
			sequence: []map[string][]int16{
				{"a": []int16{1, 2, 3}, "b": []int16{4, 5, 6}}, {"a": []int16{7, 8, 9}, "b": []int16{10, 11, 12}}, {"a": []int16{13, 14, 15}, "b": []int16{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int16{4, 10, 12},
			expect: []map[string][]int16{
				{"a": []int16{1, 2, 3}, "b": []int16{4, 5, 6}}, {"a": []int16{7, 8, 9}, "b": []int16{10, 11, 12}},
			},
		},
		{
			sequence: []map[string][]int32{
				{"a": []int32{1, 2, 3}, "b": []int32{4, 5, 6}}, {"a": []int32{7, 8, 9}, "b": []int32{10, 11, 12}}, {"a": []int32{13, 14, 15}, "b": []int32{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int32{4, 10, 12},
			expect: []map[string][]int32{
				{"a": []int32{1, 2, 3}, "b": []int32{4, 5, 6}}, {"a": []int32{7, 8, 9}, "b": []int32{10, 11, 12}},
			},
		},
		{
			sequence: []map[string][]int64{
				{"a": []int64{1, 2, 3}, "b": []int64{4, 5, 6}}, {"a": []int64{7, 8, 9}, "b": []int64{10, 11, 12}}, {"a": []int64{13, 14, 15}, "b": []int64{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int64{4, 10, 12},
			expect: []map[string][]int64{
				{"a": []int64{1, 2, 3}, "b": []int64{4, 5, 6}}, {"a": []int64{7, 8, 9}, "b": []int64{10, 11, 12}},
			},
		},
		{
			sequence: []map[string][]float32{
				{"a": []float32{1.0, 2.0, 3.0}, "b": []float32{4.0, 5.0, 6.0}}, {"a": []float32{7.0, 8.0, 9.0}, "b": []float32{10.0, 11.0, 12.0}}, {"a": []float32{13.0, 14.0, 15.0}, "b": []float32{16.0, 17.0, 18.0}},
			},
			key: "b", op: "intersect", match: []float32{4, 10, 12},
			expect: []map[string][]float32{
				{"a": []float32{1.0, 2.0, 3.0}, "b": []float32{4.0, 5.0, 6.0}}, {"a": []float32{7.0, 8.0, 9.0}, "b": []float32{10.0, 11.0, 12.0}},
			},
		},
		{
			sequence: []map[string][]float64{
				{"a": []float64{1.0, 2.0, 3.0}, "b": []float64{4.0, 5.0, 6.0}}, {"a": []float64{7.0, 8.0, 9.0}, "b": []float64{10.0, 11.0, 12.0}}, {"a": []float64{13.0, 14.0, 15.0}, "b": []float64{16.0, 17.0, 18.0}},
			},
			key: "b", op: "intersect", match: []float64{4, 10, 12},
			expect: []map[string][]float64{
				{"a": []float64{1.0, 2.0, 3.0}, "b": []float64{4.0, 5.0, 6.0}}, {"a": []float64{7.0, 8.0, 9.0}, "b": []float64{10.0, 11.0, 12.0}},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: "in", match: slice(3, 4, 5),
			expect: []map[string]int{
				{"a": 3, "b": 4},
			},
		},
		{
			sequence: []map[string]time.Time{
				{"a": d1, "b": d2}, {"a": d3, "b": d4}, {"a": d5, "b": d6},
			},
			key: "b", op: "in", match: slice(d3, d4, d5),
			expect: []map[string]time.Time{
				{"a": d3, "b": d4},
			},
		},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "not in", match: []string{"c", "d", "e"},
			expect: []TstX{
				{A: "a", B: "b"}, {A: "e", B: "f"},
			},
		},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "not in", match: slice("c", t, "d", "e"),
			expect: []TstX{
				{A: "a", B: "b"}, {A: "e", B: "f"},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: "", match: nil,
			expect: []map[string]int{
				{"a": 3},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: "!=", match: nil,
			expect: []map[string]int{
				{"a": 1, "b": 2}, {"a": 5, "b": 6},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: ">", match: nil,
			expect: []map[string]int{},
		},
		{
			sequence: []map[string]bool{
				{"a": true, "b": false}, {"c": true, "b": true}, {"d": true, "b": false},
			},
			key: "b", op: "", match: true,
			expect: []map[string]bool{
				{"c": true, "b": true},
			},
		},
		{
			sequence: []map[string]bool{
				{"a": true, "b": false}, {"c": true, "b": true}, {"d": true, "b": false},
			},
			key: "b", op: "!=", match: true,
			expect: []map[string]bool{
				{"a": true, "b": false}, {"d": true, "b": false},
			},
		},
		{
			sequence: []map[string]bool{
				{"a": true, "b": false}, {"c": true, "b": true}, {"d": true, "b": false},
			},
			key: "b", op: ">", match: false,
			expect: []map[string]bool{},
		},
		{sequence: (*[]TstX)(nil), key: "A", match: "a", expect: false},
		{sequence: TstX{A: "a", B: "b"}, key: "A", match: "a", expect: false},
		{sequence: []map[string]*TstX{{"foo": nil}}, key: "foo.B", match: "d", expect: false},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "op", match: "f",
			expect: false,
		},
		{
			sequence: map[string]interface{}{
				"foo": []interface{}{map[interface{}]interface{}{"a": 1, "b": 2}},
				"bar": []interface{}{map[interface{}]interface{}{"a": 3, "b": 4}},
				"zap": []interface{}{map[interface{}]interface{}{"a": 5, "b": 6}},
			},
			key: "b", op: "in", match: slice(3, 4, 5),
			expect: map[string]interface{}{
				"bar": []interface{}{map[interface{}]interface{}{"a": 3, "b": 4}},
			},
		},
		{
			sequence: map[string]interface{}{
				"foo": []interface{}{map[interface{}]interface{}{"a": 1, "b": 2}},
				"bar": []interface{}{map[interface{}]interface{}{"a": 3, "b": 4}},
				"zap": []interface{}{map[interface{}]interface{}{"a": 5, "b": 6}},
			},
			key: "b", op: ">", match: 3,
			expect: map[string]interface{}{
				"bar": []interface{}{map[interface{}]interface{}{"a": 3, "b": 4}},
				"zap": []interface{}{map[interface{}]interface{}{"a": 5, "b": 6}},
			},
		},
	} {
		var results interface{}
		var err error

		if len(this.op) > 0 {
			results, err = where(this.sequence, this.key, this.op, this.match)
		} else {
			results, err = where(this.sequence, this.key, this.match)
		}
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Where didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, this.expect) {
				t.Errorf("[%d] Where clause matching %v with %v, got %v but expected %v", i, this.key, this.match, results, this.expect)
			}
		}
	}

	var err error
	_, err = where(map[string]int{"a": 1, "b": 2}, "a", []byte("="), 1)
	if err == nil {
		t.Errorf("Where called with none string op value didn't return an expected error")
	}

	_, err = where(map[string]int{"a": 1, "b": 2}, "a", []byte("="), 1, 2)
	if err == nil {
		t.Errorf("Where called with more than two variable arguments didn't return an expected error")
	}

	_, err = where(map[string]int{"a": 1, "b": 2}, "a")
	if err == nil {
		t.Errorf("Where called with no variable arguments didn't return an expected error")
	}
}

func TestDelimit(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		sequence  interface{}
		delimiter interface{}
		last      interface{}
		expect    template.HTML
	}{
		{[]string{"class1", "class2", "class3"}, " ", nil, "class1 class2 class3"},
		{[]int{1, 2, 3, 4, 5}, ",", nil, "1,2,3,4,5"},
		{[]int{1, 2, 3, 4, 5}, ", ", nil, "1, 2, 3, 4, 5"},
		{[]string{"class1", "class2", "class3"}, " ", " and ", "class1 class2 and class3"},
		{[]int{1, 2, 3, 4, 5}, ",", ",", "1,2,3,4,5"},
		{[]int{1, 2, 3, 4, 5}, ", ", ", and ", "1, 2, 3, 4, and 5"},
		// test maps with and without sorting required
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, "--", nil, "10--20--30--40--50"},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, "--", nil, "30--20--10--40--50"},
		{map[string]string{"1": "10", "2": "20", "3": "30", "4": "40", "5": "50"}, "--", nil, "10--20--30--40--50"},
		{map[string]string{"3": "10", "2": "20", "1": "30", "4": "40", "5": "50"}, "--", nil, "30--20--10--40--50"},
		{map[string]string{"one": "10", "two": "20", "three": "30", "four": "40", "five": "50"}, "--", nil, "50--40--10--30--20"},
		{map[int]string{1: "10", 2: "20", 3: "30", 4: "40", 5: "50"}, "--", nil, "10--20--30--40--50"},
		{map[int]string{3: "10", 2: "20", 1: "30", 4: "40", 5: "50"}, "--", nil, "30--20--10--40--50"},
		{map[float64]string{3.3: "10", 2.3: "20", 1.3: "30", 4.3: "40", 5.3: "50"}, "--", nil, "30--20--10--40--50"},
		// test maps with a last delimiter
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, "--", "--and--", "10--20--30--40--and--50"},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, "--", "--and--", "30--20--10--40--and--50"},
		{map[string]string{"1": "10", "2": "20", "3": "30", "4": "40", "5": "50"}, "--", "--and--", "10--20--30--40--and--50"},
		{map[string]string{"3": "10", "2": "20", "1": "30", "4": "40", "5": "50"}, "--", "--and--", "30--20--10--40--and--50"},
		{map[string]string{"one": "10", "two": "20", "three": "30", "four": "40", "five": "50"}, "--", "--and--", "50--40--10--30--and--20"},
		{map[int]string{1: "10", 2: "20", 3: "30", 4: "40", 5: "50"}, "--", "--and--", "10--20--30--40--and--50"},
		{map[int]string{3: "10", 2: "20", 1: "30", 4: "40", 5: "50"}, "--", "--and--", "30--20--10--40--and--50"},
		{map[float64]string{3.5: "10", 2.5: "20", 1.5: "30", 4.5: "40", 5.5: "50"}, "--", "--and--", "30--20--10--40--and--50"},
	} {
		var result template.HTML
		var err error
		if this.last == nil {
			result, err = delimit(this.sequence, this.delimiter)
		} else {
			result, err = delimit(this.sequence, this.delimiter, this.last)
		}
		if err != nil {
			t.Errorf("[%d] failed: %s", i, err)
			continue
		}
		if !reflect.DeepEqual(result, this.expect) {
			t.Errorf("[%d] Delimit called on sequence: %v | delimiter: `%v` | last: `%v`, got %v but expected %v", i, this.sequence, this.delimiter, this.last, result, this.expect)
		}
	}
}

func TestSort(t *testing.T) {
	t.Parallel()
	type ts struct {
		MyInt    int
		MyFloat  float64
		MyString string
	}
	type mid struct {
		Tst TstX
	}

	for i, this := range []struct {
		sequence    interface{}
		sortByField interface{}
		sortAsc     string
		expect      interface{}
	}{
		{[]string{"class1", "class2", "class3"}, nil, "asc", []string{"class1", "class2", "class3"}},
		{[]string{"class3", "class1", "class2"}, nil, "asc", []string{"class1", "class2", "class3"}},
		{[]int{1, 2, 3, 4, 5}, nil, "asc", []int{1, 2, 3, 4, 5}},
		{[]int{5, 4, 3, 1, 2}, nil, "asc", []int{1, 2, 3, 4, 5}},
		// test sort key parameter is focibly set empty
		{[]string{"class3", "class1", "class2"}, map[int]string{1: "a"}, "asc", []string{"class1", "class2", "class3"}},
		// test map sorting by keys
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, nil, "asc", []int{10, 20, 30, 40, 50}},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, nil, "asc", []int{30, 20, 10, 40, 50}},
		{map[string]string{"1": "10", "2": "20", "3": "30", "4": "40", "5": "50"}, nil, "asc", []string{"10", "20", "30", "40", "50"}},
		{map[string]string{"3": "10", "2": "20", "1": "30", "4": "40", "5": "50"}, nil, "asc", []string{"30", "20", "10", "40", "50"}},
		{map[string]string{"one": "10", "two": "20", "three": "30", "four": "40", "five": "50"}, nil, "asc", []string{"50", "40", "10", "30", "20"}},
		{map[int]string{1: "10", 2: "20", 3: "30", 4: "40", 5: "50"}, nil, "asc", []string{"10", "20", "30", "40", "50"}},
		{map[int]string{3: "10", 2: "20", 1: "30", 4: "40", 5: "50"}, nil, "asc", []string{"30", "20", "10", "40", "50"}},
		{map[float64]string{3.3: "10", 2.3: "20", 1.3: "30", 4.3: "40", 5.3: "50"}, nil, "asc", []string{"30", "20", "10", "40", "50"}},
		// test map sorting by value
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, "value", "asc", []int{10, 20, 30, 40, 50}},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, "value", "asc", []int{10, 20, 30, 40, 50}},
		// test map sorting by field value
		{
			map[string]ts{"1": {10, 10.5, "ten"}, "2": {20, 20.5, "twenty"}, "3": {30, 30.5, "thirty"}, "4": {40, 40.5, "forty"}, "5": {50, 50.5, "fifty"}},
			"MyInt",
			"asc",
			[]ts{{10, 10.5, "ten"}, {20, 20.5, "twenty"}, {30, 30.5, "thirty"}, {40, 40.5, "forty"}, {50, 50.5, "fifty"}},
		},
		{
			map[string]ts{"1": {10, 10.5, "ten"}, "2": {20, 20.5, "twenty"}, "3": {30, 30.5, "thirty"}, "4": {40, 40.5, "forty"}, "5": {50, 50.5, "fifty"}},
			"MyFloat",
			"asc",
			[]ts{{10, 10.5, "ten"}, {20, 20.5, "twenty"}, {30, 30.5, "thirty"}, {40, 40.5, "forty"}, {50, 50.5, "fifty"}},
		},
		{
			map[string]ts{"1": {10, 10.5, "ten"}, "2": {20, 20.5, "twenty"}, "3": {30, 30.5, "thirty"}, "4": {40, 40.5, "forty"}, "5": {50, 50.5, "fifty"}},
			"MyString",
			"asc",
			[]ts{{50, 50.5, "fifty"}, {40, 40.5, "forty"}, {10, 10.5, "ten"}, {30, 30.5, "thirty"}, {20, 20.5, "twenty"}},
		},
		// test sort desc
		{[]string{"class1", "class2", "class3"}, "value", "desc", []string{"class3", "class2", "class1"}},
		{[]string{"class3", "class1", "class2"}, "value", "desc", []string{"class3", "class2", "class1"}},
		// test sort by struct's method
		{
			[]TstX{{A: "i", B: "j"}, {A: "e", B: "f"}, {A: "c", B: "d"}, {A: "g", B: "h"}, {A: "a", B: "b"}},
			"TstRv",
			"asc",
			[]TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}, {A: "g", B: "h"}, {A: "i", B: "j"}},
		},
		{
			[]*TstX{{A: "i", B: "j"}, {A: "e", B: "f"}, {A: "c", B: "d"}, {A: "g", B: "h"}, {A: "a", B: "b"}},
			"TstRp",
			"asc",
			[]*TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}, {A: "g", B: "h"}, {A: "i", B: "j"}},
		},
		// test map sorting by struct's method
		{
			map[string]TstX{"1": {A: "i", B: "j"}, "2": {A: "e", B: "f"}, "3": {A: "c", B: "d"}, "4": {A: "g", B: "h"}, "5": {A: "a", B: "b"}},
			"TstRv",
			"asc",
			[]TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}, {A: "g", B: "h"}, {A: "i", B: "j"}},
		},
		{
			map[string]*TstX{"1": {A: "i", B: "j"}, "2": {A: "e", B: "f"}, "3": {A: "c", B: "d"}, "4": {A: "g", B: "h"}, "5": {A: "a", B: "b"}},
			"TstRp",
			"asc",
			[]*TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}, {A: "g", B: "h"}, {A: "i", B: "j"}},
		},
		// test sort by dot chaining key argument
		{
			[]map[string]TstX{{"foo": TstX{A: "e", B: "f"}}, {"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}},
			"foo.A",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			[]map[string]TstX{{"foo": TstX{A: "e", B: "f"}}, {"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}},
			".foo.A",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			[]map[string]TstX{{"foo": TstX{A: "e", B: "f"}}, {"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}},
			"foo.TstRv",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			[]map[string]*TstX{{"foo": &TstX{A: "e", B: "f"}}, {"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}},
			"foo.TstRp",
			"asc",
			[]map[string]*TstX{{"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}, {"foo": &TstX{A: "e", B: "f"}}},
		},
		{
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "e", B: "f"}}}, {"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}},
			"foo.Tst.A",
			"asc",
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": mid{Tst: TstX{A: "e", B: "f"}}}},
		},
		{
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "e", B: "f"}}}, {"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}},
			"foo.Tst.TstRv",
			"asc",
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": mid{Tst: TstX{A: "e", B: "f"}}}},
		},
		// test map sorting by dot chaining key argument
		{
			map[string]map[string]TstX{"1": {"foo": TstX{A: "e", B: "f"}}, "2": {"foo": TstX{A: "a", B: "b"}}, "3": {"foo": TstX{A: "c", B: "d"}}},
			"foo.A",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			map[string]map[string]TstX{"1": {"foo": TstX{A: "e", B: "f"}}, "2": {"foo": TstX{A: "a", B: "b"}}, "3": {"foo": TstX{A: "c", B: "d"}}},
			".foo.A",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			map[string]map[string]TstX{"1": {"foo": TstX{A: "e", B: "f"}}, "2": {"foo": TstX{A: "a", B: "b"}}, "3": {"foo": TstX{A: "c", B: "d"}}},
			"foo.TstRv",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			map[string]map[string]*TstX{"1": {"foo": &TstX{A: "e", B: "f"}}, "2": {"foo": &TstX{A: "a", B: "b"}}, "3": {"foo": &TstX{A: "c", B: "d"}}},
			"foo.TstRp",
			"asc",
			[]map[string]*TstX{{"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}, {"foo": &TstX{A: "e", B: "f"}}},
		},
		{
			map[string]map[string]mid{"1": {"foo": mid{Tst: TstX{A: "e", B: "f"}}}, "2": {"foo": mid{Tst: TstX{A: "a", B: "b"}}}, "3": {"foo": mid{Tst: TstX{A: "c", B: "d"}}}},
			"foo.Tst.A",
			"asc",
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": mid{Tst: TstX{A: "e", B: "f"}}}},
		},
		{
			map[string]map[string]mid{"1": {"foo": mid{Tst: TstX{A: "e", B: "f"}}}, "2": {"foo": mid{Tst: TstX{A: "a", B: "b"}}}, "3": {"foo": mid{Tst: TstX{A: "c", B: "d"}}}},
			"foo.Tst.TstRv",
			"asc",
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": mid{Tst: TstX{A: "e", B: "f"}}}},
		},
		// interface slice with missing elements
		{
			[]interface{}{
				map[interface{}]interface{}{"Title": "Foo", "Weight": 10},
				map[interface{}]interface{}{"Title": "Bar"},
				map[interface{}]interface{}{"Title": "Zap", "Weight": 5},
			},
			"Weight",
			"asc",
			[]interface{}{
				map[interface{}]interface{}{"Title": "Bar"},
				map[interface{}]interface{}{"Title": "Zap", "Weight": 5},
				map[interface{}]interface{}{"Title": "Foo", "Weight": 10},
			},
		},
		// test error cases
		{(*[]TstX)(nil), nil, "asc", false},
		{TstX{A: "a", B: "b"}, nil, "asc", false},
		{
			[]map[string]TstX{{"foo": TstX{A: "e", B: "f"}}, {"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}},
			"foo.NotAvailable",
			"asc",
			false,
		},
		{
			map[string]map[string]TstX{"1": {"foo": TstX{A: "e", B: "f"}}, "2": {"foo": TstX{A: "a", B: "b"}}, "3": {"foo": TstX{A: "c", B: "d"}}},
			"foo.NotAvailable",
			"asc",
			false,
		},
		{nil, nil, "asc", false},
	} {
		var result interface{}
		var err error
		if this.sortByField == nil {
			result, err = sortSeq(this.sequence)
		} else {
			result, err = sortSeq(this.sequence, this.sortByField, this.sortAsc)
		}

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Sort didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] Sort called on sequence: %v | sortByField: `%v` | got %v but expected %v", i, this.sequence, this.sortByField, result, this.expect)
			}
		}
	}
}

func TestReturnWhenSet(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		data   interface{}
		key    interface{}
		expect interface{}
	}{
		{[]int{1, 2, 3}, 1, int64(2)},
		{[]uint{1, 2, 3}, 1, uint64(2)},
		{[]float64{1.1, 2.2, 3.3}, 1, float64(2.2)},
		{[]string{"foo", "bar", "baz"}, 1, "bar"},
		{[]TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}}, 1, ""},
		{map[string]int{"foo": 1, "bar": 2, "baz": 3}, "bar", int64(2)},
		{map[string]uint{"foo": 1, "bar": 2, "baz": 3}, "bar", uint64(2)},
		{map[string]float64{"foo": 1.1, "bar": 2.2, "baz": 3.3}, "bar", float64(2.2)},
		{map[string]string{"foo": "FOO", "bar": "BAR", "baz": "BAZ"}, "bar", "BAR"},
		{map[string]TstX{"foo": {A: "a", B: "b"}, "bar": {A: "c", B: "d"}, "baz": {A: "e", B: "f"}}, "bar", ""},
		{(*[]string)(nil), "bar", ""},
	} {
		result := returnWhenSet(this.data, this.key)
		if !reflect.DeepEqual(result, this.expect) {
			t.Errorf("[%d] ReturnWhenSet got %v (type %v) but expected %v (type %v)", i, result, reflect.TypeOf(result), this.expect, reflect.TypeOf(this.expect))
		}
	}
}

func TestMarkdownify(t *testing.T) {
	t.Parallel()
	v := viper.New()

	f := newTestFuncsterWithViper(v)

	for i, this := range []struct {
		in     interface{}
		expect interface{}
	}{
		{"Hello **World!**", template.HTML("Hello <strong>World!</strong>")},
		{[]byte("Hello Bytes **World!**"), template.HTML("Hello Bytes <strong>World!</strong>")},
	} {
		result, err := f.markdownify(this.in)
		if err != nil {
			t.Fatalf("[%d] unexpected error in markdownify: %s", i, err)
		}
		if !reflect.DeepEqual(result, this.expect) {
			t.Errorf("[%d] markdownify got %v (type %v) but expected %v (type %v)", i, result, reflect.TypeOf(result), this.expect, reflect.TypeOf(this.expect))
		}
	}

	if _, err := f.markdownify(t); err == nil {
		t.Fatalf("markdownify should have errored")
	}
}

func TestApply(t *testing.T) {
	t.Parallel()

	f := newTestFuncster()

	strings := []interface{}{"a\n", "b\n"}
	noStringers := []interface{}{tstNoStringer{}, tstNoStringer{}}

	chomped, _ := f.apply(strings, "chomp", ".")
	assert.Equal(t, []interface{}{template.HTML("a"), template.HTML("b")}, chomped)

	chomped, _ = f.apply(strings, "chomp", "c\n")
	assert.Equal(t, []interface{}{template.HTML("c"), template.HTML("c")}, chomped)

	chomped, _ = f.apply(nil, "chomp", ".")
	assert.Equal(t, []interface{}{}, chomped)

	_, err := f.apply(strings, "apply", ".")
	if err == nil {
		t.Errorf("apply with apply should fail")
	}

	var nilErr *error
	_, err = f.apply(nilErr, "chomp", ".")
	if err == nil {
		t.Errorf("apply with nil in seq should fail")
	}

	_, err = f.apply(strings, "dobedobedo", ".")
	if err == nil {
		t.Errorf("apply with unknown func should fail")
	}

	_, err = f.apply(noStringers, "chomp", ".")
	if err == nil {
		t.Errorf("apply when func fails should fail")
	}

	_, err = f.apply(tstNoStringer{}, "chomp", ".")
	if err == nil {
		t.Errorf("apply with non-sequence should fail")
	}
}

func TestChomp(t *testing.T) {
	t.Parallel()
	base := "\n This is\na story "
	for i, item := range []string{
		"\n", "\n\n",
		"\r", "\r\r",
		"\r\n", "\r\n\r\n",
	} {
		c, _ := chomp(base + item)
		chomped := string(c)

		if chomped != base {
			t.Errorf("[%d] Chomp failed, got '%v'", i, chomped)
		}

		_, err := chomp(tstNoStringer{})

		if err == nil {
			t.Errorf("Chomp should fail")
		}
	}
}

func TestLower(t *testing.T) {
	t.Parallel()
	cases := []struct {
		s     interface{}
		want  string
		isErr bool
	}{
		{"TEST", "test", false},
		{template.HTML("LoWeR"), "lower", false},
		{[]byte("BYTES"), "bytes", false},
	}

	for i, c := range cases {
		res, err := lower(c.s)
		if (err != nil) != c.isErr {
			t.Fatalf("[%d] unexpected isErr state: want %v, got %v, err = %v", i, c.want, (err != nil), err)
		}

		if res != c.want {
			t.Errorf("[%d] lower failed: want %v, got %v", i, c.want, res)
		}
	}
}

func TestTitle(t *testing.T) {
	t.Parallel()
	cases := []struct {
		s     interface{}
		want  string
		isErr bool
	}{
		{"test", "Test", false},
		{template.HTML("hypertext"), "Hypertext", false},
		{[]byte("bytes"), "Bytes", false},
	}

	for i, c := range cases {
		res, err := title(c.s)
		if (err != nil) != c.isErr {
			t.Fatalf("[%d] unexpected isErr state: want %v, got %v, err = %v", i, c.want, (err != nil), err)
		}

		if res != c.want {
			t.Errorf("[%d] title failed: want %v, got %v", i, c.want, res)
		}
	}
}

func TestUpper(t *testing.T) {
	t.Parallel()
	cases := []struct {
		s     interface{}
		want  string
		isErr bool
	}{
		{"test", "TEST", false},
		{template.HTML("UpPeR"), "UPPER", false},
		{[]byte("bytes"), "BYTES", false},
	}

	for i, c := range cases {
		res, err := upper(c.s)
		if (err != nil) != c.isErr {
			t.Fatalf("[%d] unexpected isErr state: want %v, got %v, err = %v", i, c.want, (err != nil), err)
		}

		if res != c.want {
			t.Errorf("[%d] upper failed: want %v, got %v", i, c.want, res)
		}
	}
}

func TestHighlight(t *testing.T) {
	t.Parallel()
	code := "func boo() {}"

	f := newTestFuncster()

	highlighted, err := f.highlight(code, "go", "")

	if err != nil {
		t.Fatal("Highlight returned error:", err)
	}

	// this depends on a Pygments installation, but will always contain the function name.
	if !strings.Contains(string(highlighted), "boo") {
		t.Errorf("Highlight mismatch,  got %v", highlighted)
	}

	_, err = f.highlight(t, "go", "")

	if err == nil {
		t.Error("Expected highlight error")
	}
}

func TestInflect(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		inflectFunc func(i interface{}) (string, error)
		in          interface{}
		expected    string
	}{
		{humanize, "MyCamel", "My camel"},
		{humanize, "", ""},
		{humanize, "103", "103rd"},
		{humanize, "41", "41st"},
		{humanize, 103, "103rd"},
		{humanize, int64(92), "92nd"},
		{humanize, "5.5", "5.5"},
		{pluralize, "cat", "cats"},
		{pluralize, "", ""},
		{singularize, "cats", "cat"},
		{singularize, "", ""},
	} {

		result, err := this.inflectFunc(this.in)

		if err != nil {
			t.Errorf("[%d] Unexpected Inflect error: %s", i, err)
		} else if result != this.expected {
			t.Errorf("[%d] Inflect method error, got %v expected %v", i, result, this.expected)
		}

		_, err = this.inflectFunc(t)
		if err == nil {
			t.Errorf("[%d] Expected Inflect error", i)
		}
	}
}

func TestCounterFuncs(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		countFunc func(i interface{}) (int, error)
		in        string
		expected  int
	}{
		{countWords, "Do Be Do Be Do", 5},
		{countWords, "旁边", 2},
		{countRunes, "旁边", 2},
	} {

		result, err := this.countFunc(this.in)

		if err != nil {
			t.Errorf("[%d] Unexpected counter error: %s", i, err)
		} else if result != this.expected {
			t.Errorf("[%d] Count method error, got %v expected %v", i, result, this.expected)
		}

		_, err = this.countFunc(t)
		if err == nil {
			t.Errorf("[%d] Expected Count error", i)
		}
	}
}

func TestReplace(t *testing.T) {
	t.Parallel()
	v, _ := replace("aab", "a", "b")
	assert.Equal(t, "bbb", v)
	v, _ = replace("11a11", 1, 2)
	assert.Equal(t, "22a22", v)
	v, _ = replace(12345, 1, 2)
	assert.Equal(t, "22345", v)
	_, e := replace(tstNoStringer{}, "a", "b")
	assert.NotNil(t, e, "tstNoStringer isn't trimmable")
	_, e = replace("a", tstNoStringer{}, "b")
	assert.NotNil(t, e, "tstNoStringer cannot be converted to string")
	_, e = replace("a", "b", tstNoStringer{})
	assert.NotNil(t, e, "tstNoStringer cannot be converted to string")
}

func TestReplaceRE(t *testing.T) {
	t.Parallel()
	for i, val := range []struct {
		pattern interface{}
		repl    interface{}
		src     interface{}
		expect  string
		ok      bool
	}{
		{"^https?://([^/]+).*", "$1", "http://gohugo.io/docs", "gohugo.io", true},
		{"^https?://([^/]+).*", "$2", "http://gohugo.io/docs", "", true},
		{tstNoStringer{}, "$2", "http://gohugo.io/docs", "", false},
		{"^https?://([^/]+).*", tstNoStringer{}, "http://gohugo.io/docs", "", false},
		{"^https?://([^/]+).*", "$2", tstNoStringer{}, "", false},
		{"(ab)", "AB", "aabbaab", "aABbaAB", true},
		{"(ab", "AB", "aabb", "", false}, // invalid re
	} {
		v, err := replaceRE(val.pattern, val.repl, val.src)
		if (err == nil) != val.ok {
			t.Errorf("[%d] %s", i, err)
		}
		assert.Equal(t, val.expect, v)
	}
}

func TestFindRE(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		expr    string
		content interface{}
		limit   interface{}
		expect  []string
		ok      bool
	}{
		{"[G|g]o", "Hugo is a static site generator written in Go.", 2, []string{"go", "Go"}, true},
		{"[G|g]o", "Hugo is a static site generator written in Go.", -1, []string{"go", "Go"}, true},
		{"[G|g]o", "Hugo is a static site generator written in Go.", 1, []string{"go"}, true},
		{"[G|g]o", "Hugo is a static site generator written in Go.", "1", []string{"go"}, true},
		{"[G|g]o", "Hugo is a static site generator written in Go.", nil, []string(nil), true},
		{"[G|go", "Hugo is a static site generator written in Go.", nil, []string(nil), false},
		{"[G|g]o", t, nil, []string(nil), false},
	} {
		var (
			res []string
			err error
		)

		res, err = findRE(this.expr, this.content, this.limit)
		if err != nil && this.ok {
			t.Errorf("[%d] returned an unexpected error: %s", i, err)
		}

		assert.Equal(t, this.expect, res)
	}
}

func TestTrim(t *testing.T) {
	t.Parallel()

	for i, this := range []struct {
		v1     interface{}
		v2     string
		expect interface{}
	}{
		{"1234 my way 13", "123 ", "4 my way"},
		{"      my way  ", " ", "my way"},
		{1234, "14", "23"},
		{tstNoStringer{}, " ", false},
	} {
		result, err := trim(this.v1, this.v2)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] trim didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] got '%s' but expected %s", i, result, this.expect)
			}
		}
	}
}

func TestDateFormat(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		layout string
		value  interface{}
		expect interface{}
	}{
		{"Monday, Jan 2, 2006", "2015-01-21", "Wednesday, Jan 21, 2015"},
		{"Monday, Jan 2, 2006", time.Date(2015, time.January, 21, 0, 0, 0, 0, time.UTC), "Wednesday, Jan 21, 2015"},
		{"This isn't a date layout string", "2015-01-21", "This isn't a date layout string"},
		// The following test case gives either "Tuesday, Jan 20, 2015" or "Monday, Jan 19, 2015" depending on the local time zone
		{"Monday, Jan 2, 2006", 1421733600, time.Unix(1421733600, 0).Format("Monday, Jan 2, 2006")},
		{"Monday, Jan 2, 2006", 1421733600.123, false},
		{time.RFC3339, time.Date(2016, time.March, 3, 4, 5, 0, 0, time.UTC), "2016-03-03T04:05:00Z"},
		{time.RFC1123, time.Date(2016, time.March, 3, 4, 5, 0, 0, time.UTC), "Thu, 03 Mar 2016 04:05:00 UTC"},
		{time.RFC3339, "Thu, 03 Mar 2016 04:05:00 UTC", "2016-03-03T04:05:00Z"},
		{time.RFC1123, "2016-03-03T04:05:00Z", "Thu, 03 Mar 2016 04:05:00 UTC"},
	} {
		result, err := dateFormat(this.layout, this.value)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] DateFormat didn't return an expected error, got %v", i, result)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] DateFormat failed: %s", i, err)
				continue
			}
			if result != this.expect {
				t.Errorf("[%d] DateFormat got %v but expected %v", i, result, this.expect)
			}
		}
	}
}

func TestDefaultFunc(t *testing.T) {
	t.Parallel()
	then := time.Now()
	now := time.Now()

	for i, this := range []struct {
		dflt     interface{}
		given    interface{}
		expected interface{}
	}{
		{true, false, false},
		{"5", 0, "5"},

		{"test1", "set", "set"},
		{"test2", "", "test2"},
		{"test3", nil, "test3"},

		{[2]int{10, 20}, [2]int{1, 2}, [2]int{1, 2}},
		{[2]int{10, 20}, [0]int{}, [2]int{10, 20}},
		{[2]int{100, 200}, nil, [2]int{100, 200}},

		{[]string{"one"}, []string{"uno"}, []string{"uno"}},
		{[]string{"two"}, []string{}, []string{"two"}},
		{[]string{"three"}, nil, []string{"three"}},

		{map[string]int{"one": 1}, map[string]int{"uno": 1}, map[string]int{"uno": 1}},
		{map[string]int{"one": 1}, map[string]int{}, map[string]int{"one": 1}},
		{map[string]int{"two": 2}, nil, map[string]int{"two": 2}},

		{10, 1, 1},
		{10, 0, 10},
		{20, nil, 20},

		{float32(10), float32(1), float32(1)},
		{float32(10), 0, float32(10)},
		{float32(20), nil, float32(20)},

		{complex(2, -2), complex(1, -1), complex(1, -1)},
		{complex(2, -2), complex(0, 0), complex(2, -2)},
		{complex(3, -3), nil, complex(3, -3)},

		{struct{ f string }{f: "one"}, struct{ f string }{}, struct{ f string }{}},
		{struct{ f string }{f: "two"}, nil, struct{ f string }{f: "two"}},

		{then, now, now},
		{then, time.Time{}, then},
	} {
		res, err := dfault(this.dflt, this.given)
		if err != nil {
			t.Errorf("[%d] default returned an error: %s", i, err)
			continue
		}
		if !reflect.DeepEqual(this.expected, res) {
			t.Errorf("[%d] default returned %v, but expected %v", i, res, this.expected)
		}
	}
}

func TestDefault(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		input    interface{}
		tpl      string
		expected string
		ok       bool
	}{
		{map[string]string{"foo": "bar"}, `{{ index . "foo" | default "nope" }}`, `bar`, true},
		{map[string]string{"foo": "pop"}, `{{ index . "bar" | default "nada" }}`, `nada`, true},
		{map[string]string{"foo": "cat"}, `{{ default "nope" .foo }}`, `cat`, true},
		{map[string]string{"foo": "dog"}, `{{ default "nope" .foo "extra" }}`, ``, false},
		{map[string]interface{}{"images": []string{}}, `{{ default "default.jpg" (index .images 0) }}`, `default.jpg`, true},
	} {

		tmpl := newTestTemplate(t, "test", this.tpl)

		buf := new(bytes.Buffer)
		err := tmpl.Execute(buf, this.input)
		if (err == nil) != this.ok {
			t.Errorf("[%d] execute template returned unexpected error: %s", i, err)
			continue
		}

		if buf.String() != this.expected {
			t.Errorf("[%d] execute template got %v, but expected %v", i, buf.String(), this.expected)
		}
	}
}

func TestSafeHTML(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`<div></div>`, `{{ . }}`, `&lt;div&gt;&lt;/div&gt;`, `<div></div>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		v, err := safeHTML(this.str)
		if err != nil {
			t.Fatalf("[%d] unexpected error in safeHTML: %s", i, err)
		}

		err = tmpl.Execute(buf, v)
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by safeHTML returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by safeHTML, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}

func TestSafeHTMLAttr(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`href="irc://irc.freenode.net/#golang"`, `<a {{ . }}>irc</a>`, `<a ZgotmplZ>irc</a>`, `<a href="irc://irc.freenode.net/#golang">irc</a>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		v, err := safeHTMLAttr(this.str)
		if err != nil {
			t.Fatalf("[%d] unexpected error in safeHTMLAttr: %s", i, err)
		}

		err = tmpl.Execute(buf, v)
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by safeHTMLAttr returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by safeHTMLAttr, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}

func TestSafeCSS(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`width: 60px;`, `<div style="{{ . }}"></div>`, `<div style="ZgotmplZ"></div>`, `<div style="width: 60px;"></div>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		v, err := safeCSS(this.str)
		if err != nil {
			t.Fatalf("[%d] unexpected error in safeCSS: %s", i, err)
		}

		err = tmpl.Execute(buf, v)
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by safeCSS returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by safeCSS, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}

// TODO(bep) what is this? Also look above.
func TestSafeJS(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`619c16f`, `<script>var x{{ . }};</script>`, `<script>var x"619c16f";</script>`, `<script>var x619c16f;</script>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		v, err := safeJS(this.str)
		if err != nil {
			t.Fatalf("[%d] unexpected error in safeJS: %s", i, err)
		}

		err = tmpl.Execute(buf, v)
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by safeJS returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by safeJS, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}

// TODO(bep) what is this?
func TestSafeURL(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`irc://irc.freenode.net/#golang`, `<a href="{{ . }}">IRC</a>`, `<a href="#ZgotmplZ">IRC</a>`, `<a href="irc://irc.freenode.net/#golang">IRC</a>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		v, err := safeURL(this.str)
		if err != nil {
			t.Fatalf("[%d] unexpected error in safeURL: %s", i, err)
		}

		err = tmpl.Execute(buf, v)
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by safeURL returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by safeURL, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}

func TestBase64Decode(t *testing.T) {
	t.Parallel()
	testStr := "abc123!?$*&()'-=@~"
	enc := base64.StdEncoding.EncodeToString([]byte(testStr))
	result, err := base64Decode(enc)

	if err != nil {
		t.Error("base64Decode returned error:", err)
	}

	if result != testStr {
		t.Errorf("base64Decode: got '%s', expected '%s'", result, testStr)
	}

	_, err = base64Decode(t)
	if err == nil {
		t.Error("Expected error from base64Decode")
	}
}

func TestBase64Encode(t *testing.T) {
	t.Parallel()
	testStr := "YWJjMTIzIT8kKiYoKSctPUB+"
	dec, err := base64.StdEncoding.DecodeString(testStr)

	if err != nil {
		t.Error("base64Encode: the DecodeString function of the base64 package returned an error:", err)
	}

	result, err := base64Encode(string(dec))

	if err != nil {
		t.Errorf("base64Encode: Can't cast arg '%s' into a string:", testStr)
	}

	if result != testStr {
		t.Errorf("base64Encode: got '%s', expected '%s'", result, testStr)
	}

	_, err = base64Encode(t)
	if err == nil {
		t.Error("Expected error from base64Encode")
	}
}

func TestMD5(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		input        string
		expectedHash string
	}{
		{"Hello world, gophers!", "b3029f756f98f79e7f1b7f1d1f0dd53b"},
		{"Lorem ipsum dolor", "06ce65ac476fc656bea3fca5d02cfd81"},
	} {
		result, err := md5(this.input)
		if err != nil {
			t.Errorf("md5 returned error: %s", err)
		}

		if result != this.expectedHash {
			t.Errorf("[%d] md5: expected '%s', got '%s'", i, this.expectedHash, result)
		}
	}

	_, err := md5(t)
	if err == nil {
		t.Error("Expected error from md5")
	}
}

func TestSHA1(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		input        string
		expectedHash string
	}{
		{"Hello world, gophers!", "c8b5b0e33d408246e30f53e32b8f7627a7a649d4"},
		{"Lorem ipsum dolor", "45f75b844be4d17b3394c6701768daf39419c99b"},
	} {
		result, err := sha1(this.input)
		if err != nil {
			t.Errorf("sha1 returned error: %s", err)
		}

		if result != this.expectedHash {
			t.Errorf("[%d] sha1: expected '%s', got '%s'", i, this.expectedHash, result)
		}
	}

	_, err := sha1(t)
	if err == nil {
		t.Error("Expected error from sha1")
	}
}

func TestSHA256(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		input        string
		expectedHash string
	}{
		{"Hello world, gophers!", "6ec43b78da9669f50e4e422575c54bf87536954ccd58280219c393f2ce352b46"},
		{"Lorem ipsum dolor", "9b3e1beb7053e0f900a674dd1c99aca3355e1275e1b03d3cb1bc977f5154e196"},
	} {
		result, err := sha256(this.input)
		if err != nil {
			t.Errorf("sha256 returned error: %s", err)
		}

		if result != this.expectedHash {
			t.Errorf("[%d] sha256: expected '%s', got '%s'", i, this.expectedHash, result)
		}
	}

	_, err := sha256(t)
	if err == nil {
		t.Error("Expected error from sha256")
	}
}

func TestReadFile(t *testing.T) {
	t.Parallel()

	workingDir := "/home/hugo"

	v := viper.New()

	v.Set("workingDir", workingDir)

	f := newTestFuncsterWithViper(v)

	afero.WriteFile(f.Fs.Source, filepath.Join(workingDir, "/f/f1.txt"), []byte("f1-content"), 0755)
	afero.WriteFile(f.Fs.Source, filepath.Join("/home", "f2.txt"), []byte("f2-content"), 0755)

	for i, this := range []struct {
		filename string
		expect   interface{}
	}{
		{"", false},
		{"b", false},
		{filepath.FromSlash("/f/f1.txt"), "f1-content"},
		{filepath.FromSlash("f/f1.txt"), "f1-content"},
		{filepath.FromSlash("../f2.txt"), false},
	} {
		result, err := f.readFileFromWorkingDir(this.filename)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] readFile didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] readFile failed: %s", i, err)
				continue
			}
			if result != this.expect {
				t.Errorf("[%d] readFile got %q but expected %q", i, result, this.expect)
			}
		}
	}
}

func TestPartialHTMLAndText(t *testing.T) {
	t.Parallel()
	config := newDepsConfig(viper.New())

	data := struct {
		Name string
	}{
		Name: "a+b+c", // This should get encoded in HTML.
	}

	config.WithTemplate = func(templ tpl.TemplateHandler) error {
		if err := templ.AddTemplate("htmlTemplate.html", `HTML Test Partial: {{ partial "test.foo" . -}}`); err != nil {
			return err
		}

		if err := templ.AddTemplate("_text/textTemplate.txt", `Text Test Partial: {{ partial "test.foo" . -}}`); err != nil {
			return err
		}

		// Use "foo" here to say that the extension doesn't really matter in this scenario.
		// It will look for templates in "partials/test.foo" and "partials/test.foo.html".
		if err := templ.AddTemplate("partials/test.foo", "HTML Name: {{ .Name }}"); err != nil {
			return err
		}
		if err := templ.AddTemplate("_text/partials/test.foo", "Text Name: {{ .Name }}"); err != nil {
			return err
		}

		return nil
	}

	de, err := deps.New(config)
	require.NoError(t, err)
	require.NoError(t, de.LoadResources())

	templ := de.Tmpl.Lookup("htmlTemplate.html")
	require.NotNil(t, templ)
	resultHTML, err := templ.ExecuteToString(data)
	require.NoError(t, err)

	templ = de.Tmpl.Lookup("_text/textTemplate.txt")
	require.NotNil(t, templ)
	resultText, err := templ.ExecuteToString(data)
	require.NoError(t, err)

	require.Contains(t, resultHTML, "HTML Test Partial: HTML Name: a&#43;b&#43;c")
	require.Contains(t, resultText, "Text Test Partial: Text Name: a+b+c")

}

func TestPartialHTMLWithNoSuffix(t *testing.T) {
	t.Parallel()
	config := newDepsConfig(viper.New())

	data := struct {
		Name string
	}{
		Name: "a",
	}

	config.WithTemplate = func(templ tpl.TemplateHandler) error {
		if err := templ.AddTemplate("htmlTemplate.html", `HTML Test Partial: {{ partial "test" . -}}`); err != nil {
			return err
		}

		if err := templ.AddTemplate("partials/test.html", "HTML Name: {{ .Name }}"); err != nil {
			return err
		}
		return nil
	}

	de, err := deps.New(config)
	require.NoError(t, err)
	require.NoError(t, de.LoadResources())

	templ := de.Tmpl.Lookup("htmlTemplate.html")
	require.NotNil(t, templ)
	resultHTML, err := templ.ExecuteToString(data)
	require.NoError(t, err)

	require.Contains(t, resultHTML, "HTML Test Partial: HTML Name: a")

}

func TestPartialWithError(t *testing.T) {
	t.Parallel()
	config := newDepsConfig(viper.New())

	data := struct {
		Name string
	}{
		Name: "bep",
	}

	config.WithTemplate = func(templ tpl.TemplateHandler) error {
		if err := templ.AddTemplate("container.html", `HTML Test Partial: {{ partial "fail.foo" . -}}`); err != nil {
			return err
		}

		if err := templ.AddTemplate("partials/fail.foo", "Template: {{ .DoesNotExist }}"); err != nil {
			return err
		}

		return nil
	}

	de, err := deps.New(config)
	require.NoError(t, err)
	require.NoError(t, de.LoadResources())

	templ := de.Tmpl.Lookup("container.html")
	require.NotNil(t, templ)
	result, err := templ.ExecuteToString(data)
	require.Error(t, err)

	errStr := err.Error()

	require.Contains(t, errStr, `template: container.html:1:22: executing "container.html" at <partial "fail.foo" .>`)
	require.Contains(t, errStr, `can't evaluate field DoesNotExist`)

	require.Empty(t, result)

}

func TestPartialCached(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		partial string
		tmpl    string
		variant string
	}{
		// name and partial should match between test cases.
		{"test1", "{{ .Title }} seq: {{ shuffle (seq 1 20) }}", `{{ partialCached "test1" . }}`, ""},
		{"test1", "{{ .Title }} seq: {{ shuffle (seq 1 20) }}", `{{ partialCached "test1" . "%s" }}`, "header"},
		{"test1", "{{ .Title }} seq: {{ shuffle (seq 1 20) }}", `{{ partialCached "test1" . "%s" }}`, "footer"},
		{"test1", "{{ .Title }} seq: {{ shuffle (seq 1 20) }}", `{{ partialCached "test1" . "%s" }}`, "header"},
	}

	var data struct {
		Title   string
		Section string
		Params  map[string]interface{}
	}

	data.Title = "**BatMan**"
	data.Section = "blog"
	data.Params = map[string]interface{}{"langCode": "en"}

	for i, tc := range testCases {
		var tmp string
		if tc.variant != "" {
			tmp = fmt.Sprintf(tc.tmpl, tc.variant)
		} else {
			tmp = tc.tmpl
		}

		config := newDepsConfig(viper.New())

		config.WithTemplate = func(templ tpl.TemplateHandler) error {
			err := templ.AddTemplate("testroot", tmp)
			if err != nil {
				return err
			}
			err = templ.AddTemplate("partials/"+tc.name, tc.partial)
			if err != nil {
				return err
			}

			return nil
		}

		de, err := deps.New(config)
		require.NoError(t, err)
		require.NoError(t, de.LoadResources())

		buf := new(bytes.Buffer)
		templ := de.Tmpl.Lookup("testroot")
		err = templ.Execute(buf, &data)
		if err != nil {
			t.Fatalf("[%d] error executing template: %s", i, err)
		}

		for j := 0; j < 10; j++ {
			buf2 := new(bytes.Buffer)
			err := templ.Execute(buf2, nil)
			if err != nil {
				t.Fatalf("[%d] error executing template 2nd time: %s", i, err)
			}

			if !reflect.DeepEqual(buf, buf2) {
				t.Fatalf("[%d] cached results do not match:\nResult 1:\n%q\nResult 2:\n%q", i, buf, buf2)
			}
		}
	}
}

func BenchmarkPartial(b *testing.B) {
	config := newDepsConfig(viper.New())
	config.WithTemplate = func(templ tpl.TemplateHandler) error {
		err := templ.AddTemplate("testroot", `{{ partial "bench1" . }}`)
		if err != nil {
			return err
		}
		err = templ.AddTemplate("partials/bench1", `{{ shuffle (seq 1 10) }}`)
		if err != nil {
			return err
		}

		return nil
	}

	de, err := deps.New(config)
	require.NoError(b, err)
	require.NoError(b, de.LoadResources())

	buf := new(bytes.Buffer)
	tmpl := de.Tmpl.Lookup("testroot")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := tmpl.Execute(buf, nil); err != nil {
			b.Fatalf("error executing template: %s", err)
		}
		buf.Reset()
	}
}

func BenchmarkPartialCached(b *testing.B) {
	config := newDepsConfig(viper.New())
	config.WithTemplate = func(templ tpl.TemplateHandler) error {
		err := templ.AddTemplate("testroot", `{{ partialCached "bench1" . }}`)
		if err != nil {
			return err
		}
		err = templ.AddTemplate("partials/bench1", `{{ shuffle (seq 1 10) }}`)
		if err != nil {
			return err
		}

		return nil
	}

	de, err := deps.New(config)
	require.NoError(b, err)
	require.NoError(b, de.LoadResources())

	buf := new(bytes.Buffer)
	tmpl := de.Tmpl.Lookup("testroot")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := tmpl.Execute(buf, nil); err != nil {
			b.Fatalf("error executing template: %s", err)
		}
		buf.Reset()
	}
}

func newTestFuncster() *templateFuncster {
	return newTestFuncsterWithViper(viper.New())
}

func newTestFuncsterWithViper(v *viper.Viper) *templateFuncster {
	config := newDepsConfig(v)
	d, err := deps.New(config)
	if err != nil {
		panic(err)
	}

	if err := d.LoadResources(); err != nil {
		panic(err)
	}

	return d.Tmpl.(*templateHandler).html.funcster
}

func newTestTemplate(t *testing.T, name, template string) tpl.Template {
	config := newDepsConfig(viper.New())
	config.WithTemplate = func(templ tpl.TemplateHandler) error {
		err := templ.AddTemplate(name, template)
		if err != nil {
			return err
		}
		return nil
	}

	de, err := deps.New(config)
	require.NoError(t, err)
	require.NoError(t, de.LoadResources())

	return de.Tmpl.Lookup(name)
}
