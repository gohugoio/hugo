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

package data

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/config/testconfig"

	"github.com/gohugoio/hugo/helpers"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
)

func TestScpGetLocal(t *testing.T) {
	t.Parallel()
	v := config.New()
	workingDir := "/my/working/dir"
	v.Set("workingDir", workingDir)
	v.Set("publishDir", "public")
	fs := hugofs.NewFromOld(afero.NewMemMapFs(), v)
	ps := helpers.FilePathSeparator

	tests := []struct {
		path    string
		content []byte
	}{
		{"testpath" + ps + "test.txt", []byte(`T€st Content 123 fOO,bar:foo%bAR`)},
		{"FOo" + ps + "BaR.html", []byte(`FOo/BaR.html T€st Content 123`)},
		{"трям" + ps + "трям", []byte(`T€st трям/трям Content 123`)},
		{"은행", []byte(`T€st C은행ontent 123`)},
		{"Банковский кассир", []byte(`Банковский кассир T€st Content 123`)},
	}

	for _, test := range tests {
		r := bytes.NewReader(test.content)
		err := helpers.WriteToDisk(filepath.Join(workingDir, test.path), r, fs.Source)
		if err != nil {
			t.Error(err)
		}

		c, err := getLocal(workingDir, test.path, fs.Source)
		if err != nil {
			t.Errorf("Error getting resource content: %s", err)
		}
		if !bytes.Equal(c, test.content) {
			t.Errorf("\nExpected: %s\nActual: %s\n", string(test.content), string(c))
		}
	}
}

func getTestServer(handler func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *http.Client) {
	testServer := httptest.NewServer(http.HandlerFunc(handler))
	client := &http.Client{
		Transport: &http.Transport{Proxy: func(r *http.Request) (*url.URL, error) {
			// Remove when https://github.com/golang/go/issues/13686 is fixed
			r.Host = "gohugo.io"
			return url.Parse(testServer.URL)
		}},
	}
	return testServer, client
}

func TestScpGetRemote(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	fs := new(afero.MemMapFs)
	cache := filecache.NewCache(fs, 100, "")

	tests := []struct {
		path    string
		content []byte
	}{
		{"http://Foo.Bar/foo_Bar-Foo", []byte(`T€st Content 123`)},
		{"http://Doppel.Gänger/foo_Bar-Foo", []byte(`T€st Cont€nt 123`)},
		{"http://Doppel.Gänger/Fizz_Bazz-Foo", []byte(`T€st Банковский кассир Cont€nt 123`)},
		{"http://Doppel.Gänger/Fizz_Bazz-Bar", []byte(`T€st Банковский кассир Cont€nt 456`)},
	}

	for _, test := range tests {
		msg := qt.Commentf("%v", test)

		req, err := http.NewRequest("GET", test.path, nil)
		c.Assert(err, qt.IsNil, msg)

		srv, cl := getTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.Write(test.content)
		})
		defer func() { srv.Close() }()

		ns := newTestNs()
		ns.client = cl

		var cb []byte
		f := func(b []byte) (bool, error) {
			cb = b
			return false, nil
		}

		err = ns.getRemote(cache, f, req)
		c.Assert(err, qt.IsNil, msg)
		c.Assert(string(cb), qt.Equals, string(test.content))

		c.Assert(string(cb), qt.Equals, string(test.content))

	}
}

func TestScpGetRemoteParallel(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	content := []byte(`T€st Content 123`)
	srv, cl := getTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	})

	defer func() { srv.Close() }()

	url := "http://Foo.Bar/foo_Bar-Foo"
	req, err := http.NewRequest("GET", url, nil)
	c.Assert(err, qt.IsNil)

	for _, ignoreCache := range []bool{false} {
		cfg := config.New()
		cfg.Set("ignoreCache", ignoreCache)

		ns := New(newDeps(cfg))
		ns.client = cl

		var wg sync.WaitGroup

		for i := 0; i < 1; i++ {
			wg.Add(1)
			go func(gor int) {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					var cb []byte
					f := func(b []byte) (bool, error) {
						cb = b
						return false, nil
					}
					err := ns.getRemote(ns.cacheGetJSON, f, req)

					c.Assert(err, qt.IsNil)
					if string(content) != string(cb) {
						t.Errorf("expected\n%q\ngot\n%q", content, cb)
					}

					time.Sleep(23 * time.Millisecond)
				}
			}(i)
		}

		wg.Wait()
	}
}

func newDeps(cfg config.Provider) *deps.Deps {
	conf := testconfig.GetTestConfig(nil, cfg)
	logger := loggers.NewDefault()
	fs := hugofs.NewFrom(afero.NewMemMapFs(), conf.BaseConfig())

	d := &deps.Deps{
		Fs:   fs,
		Log:  logger,
		Conf: conf,
	}
	if err := d.Init(); err != nil {
		panic(err)
	}
	return d
}

func newTestNs() *Namespace {
	return New(newDeps(config.New()))
}
