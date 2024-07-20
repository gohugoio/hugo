// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"html/template"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/maps"

	qt "github.com/frankban/quicktest"
)

func TestGetCSV(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for i, test := range []struct {
		sep     string
		url     string
		content string
		expect  any
	}{
		// Remotes
		{
			",",
			`http://success/`,
			"gomeetup,city\nyes,Sydney\nyes,San Francisco\nyes,Stockholm\n",
			[][]string{{"gomeetup", "city"}, {"yes", "Sydney"}, {"yes", "San Francisco"}, {"yes", "Stockholm"}},
		},
		{
			",",
			`http://error.extra.field/`,
			"gomeetup,city\nyes,Sydney\nyes,San Francisco\nyes,Stockholm,EXTRA\n",
			false,
		},
		{
			",",
			`http://nofound/404`,
			``,
			false,
		},

		// Locals
		{
			";",
			"pass/semi",
			"gomeetup;city\nyes;Sydney\nyes;San Francisco\nyes;Stockholm\n",
			[][]string{{"gomeetup", "city"}, {"yes", "Sydney"}, {"yes", "San Francisco"}, {"yes", "Stockholm"}},
		},
		{
			";",
			"fail/no-file",
			"",
			false,
		},
	} {
		c.Run(test.url, func(c *qt.C) {
			msg := qt.Commentf("Test %d", i)

			ns := newTestNs()

			// Setup HTTP test server
			var srv *httptest.Server
			srv, ns.client = getTestServer(func(w http.ResponseWriter, r *http.Request) {
				if !hasHeaderValue(r.Header, "Accept", "text/csv") && !hasHeaderValue(r.Header, "Accept", "text/plain") {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}

				if r.URL.Path == "/404" {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				w.Header().Add("Content-type", "text/csv")

				w.Write([]byte(test.content))
			})
			defer func() { srv.Close() }()

			// Setup local test file for schema-less URLs
			if !strings.Contains(test.url, ":") && !strings.HasPrefix(test.url, "fail/") {
				f, err := ns.deps.Fs.Source.Create(filepath.Join(ns.deps.Conf.BaseConfig().WorkingDir, test.url))
				c.Assert(err, qt.IsNil, msg)
				f.WriteString(test.content)
				f.Close()
			}

			// Get on with it
			got, err := ns.GetCSV(test.sep, test.url)

			if _, ok := test.expect.(bool); ok {
				c.Assert(int(ns.deps.Log.LoggCount(logg.LevelError)), qt.Equals, 1)
				c.Assert(got, qt.IsNil)
				return
			}

			c.Assert(err, qt.IsNil, msg)
			c.Assert(int(ns.deps.Log.LoggCount(logg.LevelError)), qt.Equals, 0)
			c.Assert(got, qt.Not(qt.IsNil), msg)
			c.Assert(got, qt.DeepEquals, test.expect, msg)
		})
	}
}

func TestGetJSON(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for i, test := range []struct {
		url     string
		content string
		expect  any
	}{
		{
			`http://success/`,
			`{"gomeetup":["Sydney","San Francisco","Stockholm"]}`,
			map[string]any{"gomeetup": []any{"Sydney", "San Francisco", "Stockholm"}},
		},
		{
			`http://malformed/`,
			`{gomeetup:["Sydney","San Francisco","Stockholm"]}`,
			false,
		},
		{
			`http://nofound/404`,
			``,
			false,
		},
		// Locals
		{
			"pass/semi",
			`{"gomeetup":["Sydney","San Francisco","Stockholm"]}`,
			map[string]any{"gomeetup": []any{"Sydney", "San Francisco", "Stockholm"}},
		},
		{
			"fail/no-file",
			"",
			false,
		},
		{
			`pass/üńīçøðê-url.json`,
			`{"gomeetup":["Sydney","San Francisco","Stockholm"]}`,
			map[string]any{"gomeetup": []any{"Sydney", "San Francisco", "Stockholm"}},
		},
	} {
		c.Run(test.url, func(c *qt.C) {
			msg := qt.Commentf("Test %d", i)
			ns := newTestNs()

			// Setup HTTP test server
			var srv *httptest.Server
			srv, ns.client = getTestServer(func(w http.ResponseWriter, r *http.Request) {
				if !hasHeaderValue(r.Header, "Accept", "application/json") {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}

				if r.URL.Path == "/404" {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				w.Header().Add("Content-type", "application/json")

				w.Write([]byte(test.content))
			})
			defer func() { srv.Close() }()

			// Setup local test file for schema-less URLs
			if !strings.Contains(test.url, ":") && !strings.HasPrefix(test.url, "fail/") {
				f, err := ns.deps.Fs.Source.Create(filepath.Join(ns.deps.Conf.BaseConfig().WorkingDir, test.url))
				c.Assert(err, qt.IsNil, msg)
				f.WriteString(test.content)
				f.Close()
			}

			// Get on with it
			got, _ := ns.GetJSON(test.url)

			if _, ok := test.expect.(bool); ok {
				c.Assert(int(ns.deps.Log.LoggCount(logg.LevelError)), qt.Equals, 1)
				return
			}

			c.Assert(int(ns.deps.Log.LoggCount(logg.LevelError)), qt.Equals, 0, msg)
			c.Assert(got, qt.Not(qt.IsNil), msg)
			c.Assert(got, qt.DeepEquals, test.expect)
		})
	}
}

func TestHeaders(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		name    string
		headers any
		assert  func(c *qt.C, headers string)
	}{
		{
			`Misc header variants`,
			map[string]any{
				"Accept-Charset": "utf-8",
				"Max-forwards":   "10",
				"X-Int":          32,
				"X-Templ":        template.HTML("a"),
				"X-Multiple":     []string{"a", "b"},
				"X-MultipleInt":  []int{3, 4},
			},
			func(c *qt.C, headers string) {
				c.Assert(headers, qt.Contains, "Accept-Charset: utf-8")
				c.Assert(headers, qt.Contains, "Max-Forwards: 10")
				c.Assert(headers, qt.Contains, "X-Int: 32")
				c.Assert(headers, qt.Contains, "X-Templ: a")
				c.Assert(headers, qt.Contains, "X-Multiple: a")
				c.Assert(headers, qt.Contains, "X-Multiple: b")
				c.Assert(headers, qt.Contains, "X-Multipleint: 3")
				c.Assert(headers, qt.Contains, "X-Multipleint: 4")
				c.Assert(headers, qt.Contains, "User-Agent: Hugo Static Site Generator")
			},
		},
		{
			`Params`,
			maps.Params{
				"Accept-Charset": "utf-8",
			},
			func(c *qt.C, headers string) {
				c.Assert(headers, qt.Contains, "Accept-Charset: utf-8")
			},
		},
		{
			`Override User-Agent`,
			map[string]any{
				"User-Agent": "007",
			},
			func(c *qt.C, headers string) {
				c.Assert(headers, qt.Contains, "User-Agent: 007")
			},
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			ns := newTestNs()

			// Setup HTTP test server
			var srv *httptest.Server
			var headers bytes.Buffer
			srv, ns.client = getTestServer(func(w http.ResponseWriter, r *http.Request) {
				c.Assert(r.URL.String(), qt.Equals, "http://gohugo.io/api?foo")
				w.Write([]byte("{}"))
				r.Header.Write(&headers)
			})
			defer func() { srv.Close() }()

			testFunc := func(fn func(args ...any) error) {
				defer headers.Reset()
				err := fn("http://example.org/api", "?foo", test.headers)

				c.Assert(err, qt.IsNil)
				c.Assert(int(ns.deps.Log.LoggCount(logg.LevelError)), qt.Equals, 0)
				test.assert(c, headers.String())
			}

			testFunc(func(args ...any) error {
				_, err := ns.GetJSON(args...)
				return err
			})
			testFunc(func(args ...any) error {
				_, err := ns.GetCSV(",", args...)
				return err
			})
		})
	}
}

func TestToURLAndHeaders(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	url, headers := toURLAndHeaders([]any{"https://foo?id=", 32})
	c.Assert(url, qt.Equals, "https://foo?id=32")
	c.Assert(headers, qt.IsNil)

	url, headers = toURLAndHeaders([]any{"https://foo?id=", 32, map[string]any{"a": "b"}})
	c.Assert(url, qt.Equals, "https://foo?id=32")
	c.Assert(headers, qt.DeepEquals, map[string]any{"a": "b"})
}

func TestParseCSV(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for i, test := range []struct {
		csv []byte
		sep string
		exp string
		err bool
	}{
		{[]byte("a,b,c\nd,e,f\n"), "", "", true},
		{[]byte("a,b,c\nd,e,f\n"), "~/", "", true},
		{[]byte("a,b,c\nd,e,f"), "|", "a,b,cd,e,f", false},
		{[]byte("q,w,e\nd,e,f"), ",", "qwedef", false},
		{[]byte("a|b|c\nd|e|f|g"), "|", "abcdefg", true},
		{[]byte("z|y|c\nd|e|f"), "|", "zycdef", false},
	} {
		msg := qt.Commentf("Test %d: %v", i, test)

		csv, err := parseCSV(test.csv, test.sep)
		if test.err {
			c.Assert(err, qt.Not(qt.IsNil), msg)
			continue
		}
		c.Assert(err, qt.IsNil, msg)

		act := ""
		for _, v := range csv {
			act = act + strings.Join(v, "")
		}

		c.Assert(act, qt.Equals, test.exp, msg)
	}
}
