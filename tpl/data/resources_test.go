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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScpGetLocal(t *testing.T) {
	t.Parallel()
	v := viper.New()
	fs := hugofs.NewMem(v)
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
		err := helpers.WriteToDisk(test.path, r, fs.Source)
		if err != nil {
			t.Error(err)
		}

		c, err := getLocal(test.path, fs.Source, v)
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
	fs := new(afero.MemMapFs)

	tests := []struct {
		path    string
		content []byte
		ignore  bool
	}{
		{"http://Foo.Bar/foo_Bar-Foo", []byte(`T€st Content 123`), false},
		{"http://Doppel.Gänger/foo_Bar-Foo", []byte(`T€st Cont€nt 123`), false},
		{"http://Doppel.Gänger/Fizz_Bazz-Foo", []byte(`T€st Банковский кассир Cont€nt 123`), false},
		{"http://Doppel.Gänger/Fizz_Bazz-Bar", []byte(`T€st Банковский кассир Cont€nt 456`), true},
	}

	for _, test := range tests {
		msg := fmt.Sprintf("%v", test)

		req, err := http.NewRequest("GET", test.path, nil)
		require.NoError(t, err, msg)

		srv, cl := getTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.Write(test.content)
		})
		defer func() { srv.Close() }()

		cfg := viper.New()

		c, err := getRemote(req, fs, cfg, cl)
		require.NoError(t, err, msg)
		assert.Equal(t, string(test.content), string(c))

		c, err = getCache(req.URL.String(), fs, cfg, test.ignore)
		require.NoError(t, err, msg)

		if test.ignore {
			assert.Empty(t, c, msg)
		} else {
			assert.Equal(t, string(test.content), string(c))

		}
	}
}

func TestScpGetRemoteParallel(t *testing.T) {
	t.Parallel()

	ns := New(newDeps(viper.New()))

	content := []byte(`T€st Content 123`)
	srv, cl := getTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	})
	defer func() { srv.Close() }()

	url := "http://Foo.Bar/foo_Bar-Foo"
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	for _, ignoreCache := range []bool{false, true} {
		cfg := viper.New()
		cfg.Set("ignoreCache", ignoreCache)

		var wg sync.WaitGroup

		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(gor int) {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					c, err := getRemote(req, ns.deps.Fs.Source, ns.deps.Cfg, cl)
					assert.NoError(t, err)
					assert.Equal(t, string(content), string(c))

					time.Sleep(23 * time.Millisecond)
				}
			}(i)
		}

		wg.Wait()
	}
}

func newDeps(cfg config.Provider) *deps.Deps {
	l := helpers.NewLanguage("en", cfg)
	l.Set("i18nDir", "i18n")
	cs, err := helpers.NewContentSpec(l)
	if err != nil {
		panic(err)
	}
	return &deps.Deps{
		Cfg:         cfg,
		Fs:          hugofs.NewMem(l),
		ContentSpec: cs,
	}
}
