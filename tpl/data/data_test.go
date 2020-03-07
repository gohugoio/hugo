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
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGetCSV(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for i, test := range []struct {
		sep     string
		url     string
		content string
		expect  interface{}
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
			`http://error.no.sep/`,
			"gomeetup;city\nyes;Sydney\nyes;San Francisco\nyes;Stockholm\n",
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
		msg := qt.Commentf("Test %d", i)

		ns := newTestNs()

		// Setup HTTP test server
		var srv *httptest.Server
		srv, ns.client = getTestServer(func(w http.ResponseWriter, r *http.Request) {
			if !haveHeader(r.Header, "Accept", "text/csv") && !haveHeader(r.Header, "Accept", "text/plain") {
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
			f, err := ns.deps.Fs.Source.Create(filepath.Join(ns.deps.Cfg.GetString("workingDir"), test.url))
			c.Assert(err, qt.IsNil, msg)
			f.WriteString(test.content)
			f.Close()
		}

		// Get on with it
		got, err := ns.GetCSV(test.sep, test.url)

		if _, ok := test.expect.(bool); ok {
			c.Assert(int(ns.deps.Log.ErrorCounter.Count()), qt.Equals, 1)
			//c.Assert(err, msg, qt.Not(qt.IsNil))
			c.Assert(got, qt.IsNil)
			continue
		}

		c.Assert(err, qt.IsNil, msg)
		c.Assert(int(ns.deps.Log.ErrorCounter.Count()), qt.Equals, 0)
		c.Assert(got, qt.Not(qt.IsNil), msg)
		c.Assert(got, qt.DeepEquals, test.expect, msg)

	}
}

func TestGetJSON(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for i, test := range []struct {
		url     string
		content string
		expect  interface{}
	}{
		{
			`http://success/`,
			`{"gomeetup":["Sydney","San Francisco","Stockholm"]}`,
			map[string]interface{}{"gomeetup": []interface{}{"Sydney", "San Francisco", "Stockholm"}},
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
			map[string]interface{}{"gomeetup": []interface{}{"Sydney", "San Francisco", "Stockholm"}},
		},
		{
			"fail/no-file",
			"",
			false,
		},
		{
			`pass/üńīçøðê-url.json`,
			`{"gomeetup":["Sydney","San Francisco","Stockholm"]}`,
			map[string]interface{}{"gomeetup": []interface{}{"Sydney", "San Francisco", "Stockholm"}},
		},
	} {

		msg := qt.Commentf("Test %d", i)
		ns := newTestNs()

		// Setup HTTP test server
		var srv *httptest.Server
		srv, ns.client = getTestServer(func(w http.ResponseWriter, r *http.Request) {
			if !haveHeader(r.Header, "Accept", "application/json") {
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
			f, err := ns.deps.Fs.Source.Create(filepath.Join(ns.deps.Cfg.GetString("workingDir"), test.url))
			c.Assert(err, qt.IsNil, msg)
			f.WriteString(test.content)
			f.Close()
		}

		// Get on with it
		got, _ := ns.GetJSON(test.url)

		if _, ok := test.expect.(bool); ok {
			c.Assert(int(ns.deps.Log.ErrorCounter.Count()), qt.Equals, 1)
			//c.Assert(err, msg, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(int(ns.deps.Log.ErrorCounter.Count()), qt.Equals, 0, msg)
		c.Assert(got, qt.Not(qt.IsNil), msg)
		c.Assert(got, qt.DeepEquals, test.expect)
	}
}

func TestJoinURL(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	c.Assert(joinURL([]interface{}{"https://foo?id=", 32}), qt.Equals, "https://foo?id=32")
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

func haveHeader(m http.Header, key, needle string) bool {
	var s []string
	var ok bool

	if s, ok = m[key]; !ok {
		return false
	}

	for _, v := range s {
		if v == needle {
			return true
		}
	}
	return false
}
