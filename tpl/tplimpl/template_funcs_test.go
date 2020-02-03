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

package tplimpl

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/gohugoio/hugo/modules"

	"github.com/gohugoio/hugo/resources/page"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/langs/i18n"
	"github.com/gohugoio/hugo/tpl"
	"github.com/gohugoio/hugo/tpl/internal"
	"github.com/gohugoio/hugo/tpl/partials"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var (
	logger = loggers.NewErrorLogger()
)

func newTestConfig() config.Provider {
	v := viper.New()
	v.Set("contentDir", "content")
	v.Set("dataDir", "data")
	v.Set("i18nDir", "i18n")
	v.Set("layoutDir", "layouts")
	v.Set("archetypeDir", "archetypes")
	v.Set("assetDir", "assets")
	v.Set("resourceDir", "resources")
	v.Set("publishDir", "public")

	langs.LoadLanguageSettings(v, nil)
	mod, err := modules.CreateProjectModule(v)
	if err != nil {
		panic(err)
	}
	v.Set("allModules", modules.Modules{mod})

	return v
}

func newDepsConfig(cfg config.Provider) deps.DepsCfg {
	l := langs.NewLanguage("en", cfg)
	return deps.DepsCfg{
		Language:            l,
		Site:                page.NewDummyHugoSite(cfg),
		Cfg:                 cfg,
		Fs:                  hugofs.NewMem(l),
		Logger:              logger,
		TemplateProvider:    DefaultTemplateProvider,
		TranslationProvider: i18n.NewTranslationProvider(),
	}
}

func TestTemplateFuncsExamples(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	workingDir := "/home/hugo"

	v := newTestConfig()

	v.Set("workingDir", workingDir)
	v.Set("multilingual", true)
	v.Set("contentDir", "content")
	v.Set("assetDir", "assets")
	v.Set("baseURL", "http://mysite.com/hugo/")
	v.Set("CurrentContentLanguage", langs.NewLanguage("en", v))

	fs := hugofs.NewMem(v)

	afero.WriteFile(fs.Source, filepath.Join(workingDir, "files", "README.txt"), []byte("Hugo Rocks!"), 0755)

	depsCfg := newDepsConfig(v)
	depsCfg.Fs = fs
	d, err := deps.New(depsCfg)
	c.Assert(err, qt.IsNil)

	var data struct {
		Title   string
		Section string
		Hugo    map[string]interface{}
		Params  map[string]interface{}
	}

	data.Title = "**BatMan**"
	data.Section = "blog"
	data.Params = map[string]interface{}{"langCode": "en"}
	data.Hugo = map[string]interface{}{"Version": hugo.MustParseVersion("0.36.1").Version()}

	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns := nsf(d)
		for _, mm := range ns.MethodMappings {
			for i, example := range mm.Examples {
				in, expected := example[0], example[1]
				d.WithTemplate = func(templ tpl.TemplateManager) error {
					c.Assert(templ.AddTemplate("test", in), qt.IsNil)
					c.Assert(templ.AddTemplate("partials/header.html", "<title>Hugo Rocks!</title>"), qt.IsNil)
					return nil
				}
				c.Assert(d.LoadResources(), qt.IsNil)

				var b bytes.Buffer
				templ, _ := d.Tmpl().Lookup("test")
				c.Assert(d.Tmpl().Execute(templ, &b, &data), qt.IsNil)
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

	c := qt.New(t)

	partial := `Now: {{ now.UnixNano }}`
	name := "testing"

	var data struct {
	}

	v := newTestConfig()

	config := newDepsConfig(v)

	config.WithTemplate = func(templ tpl.TemplateManager) error {
		err := templ.AddTemplate("partials/"+name, partial)
		if err != nil {
			return err
		}

		return nil
	}

	de, err := deps.New(config)
	c.Assert(err, qt.IsNil)
	c.Assert(de.LoadResources(), qt.IsNil)

	ns := partials.New(de)

	res1, err := ns.IncludeCached(name, &data)
	c.Assert(err, qt.IsNil)

	for j := 0; j < 10; j++ {
		time.Sleep(2 * time.Nanosecond)
		res2, err := ns.IncludeCached(name, &data)
		c.Assert(err, qt.IsNil)

		if !reflect.DeepEqual(res1, res2) {
			t.Fatalf("cache mismatch")
		}

		res3, err := ns.IncludeCached(name, &data, fmt.Sprintf("variant%d", j))
		c.Assert(err, qt.IsNil)

		if reflect.DeepEqual(res1, res3) {
			t.Fatalf("cache mismatch")
		}
	}

}

func BenchmarkPartial(b *testing.B) {
	doBenchmarkPartial(b, func(ns *partials.Namespace) error {
		_, err := ns.Include("bench1")
		return err
	})
}

func BenchmarkPartialCached(b *testing.B) {
	doBenchmarkPartial(b, func(ns *partials.Namespace) error {
		_, err := ns.IncludeCached("bench1", nil)
		return err
	})
}

func doBenchmarkPartial(b *testing.B, f func(ns *partials.Namespace) error) {
	c := qt.New(b)
	config := newDepsConfig(viper.New())
	config.WithTemplate = func(templ tpl.TemplateManager) error {
		err := templ.AddTemplate("partials/bench1", `{{ shuffle (seq 1 10) }}`)
		if err != nil {
			return err
		}

		return nil
	}

	de, err := deps.New(config)
	c.Assert(err, qt.IsNil)
	c.Assert(de.LoadResources(), qt.IsNil)

	ns := partials.New(de)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := f(ns); err != nil {
				b.Fatalf("error executing template: %s", err)
			}
		}
	})
}
