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
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"io/ioutil"
	"log"
	"os"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/i18n"
	"github.com/gohugoio/hugo/tpl"
	"github.com/gohugoio/hugo/tpl/internal"
	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
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

func TestTemplateFuncsExamples(t *testing.T) {
	t.Parallel()

	workingDir := "/home/hugo"

	v := viper.New()

	v.Set("workingDir", workingDir)
	v.Set("multilingual", true)
	v.Set("baseURL", "http://mysite.com/hugo/")
	v.Set("CurrentContentLanguage", helpers.NewLanguage("en", v))

	fs := hugofs.NewMem(v)

	afero.WriteFile(fs.Source, filepath.Join(workingDir, "README.txt"), []byte("Hugo Rocks!"), 0755)

	depsCfg := newDepsConfig(v)
	depsCfg.Fs = fs
	d, err := deps.New(depsCfg)
	require.NoError(t, err)

	var data struct {
		Title   string
		Section string
		Params  map[string]interface{}
	}

	data.Title = "**BatMan**"
	data.Section = "blog"
	data.Params = map[string]interface{}{"langCode": "en"}

	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns := nsf(d)
		for _, mm := range ns.MethodMappings {
			for i, example := range mm.Examples {
				in, expected := example[0], example[1]
				d.WithTemplate = func(templ tpl.TemplateHandler) error {
					require.NoError(t, templ.AddTemplate("test", in))
					require.NoError(t, templ.AddTemplate("partials/header.html", "<title>Hugo Rocks!</title>"))
					return nil
				}
				require.NoError(t, d.LoadResources())

				var b bytes.Buffer
				require.NoError(t, d.Tmpl.Lookup("test").Execute(&b, &data))
				if b.String() != expected {
					t.Fatalf("%s[%d]: got %q expected %q", ns.Name, i, b.String(), expected)
				}
			}
		}
	}

}

// TODO(bep) it would be dandy to put this one into the partials package, but
// we have some package cycle issues to solve first.
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
