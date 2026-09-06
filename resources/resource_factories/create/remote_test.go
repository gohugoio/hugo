// Copyright 2021 The Hugo Authors. All rights reserved.
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

package create

import (
	"net/http"
	"net/http/httptest"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/security"
)

func TestDecodeRemoteOptions(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	for _, test := range []struct {
		name    string
		args    map[string]any
		want    fromRemoteOptions
		wantErr bool
	}{
		{
			"POST",
			map[string]any{
				"meThod": "PoST",
				"headers": map[string]any{
					"foo": "bar",
				},
			},
			fromRemoteOptions{
				Method: "POST",
				Headers: map[string]any{
					"foo": "bar",
				},
			},
			false,
		},
		{
			"Body",
			map[string]any{
				"meThod": "POST",
				"body":   []byte("foo"),
			},
			fromRemoteOptions{
				Method: "POST",
				Body:   []byte("foo"),
			},
			false,
		},
		{
			"Body, string",
			map[string]any{
				"meThod": "POST",
				"body":   "foo",
			},
			fromRemoteOptions{
				Method: "POST",
				Body:   []byte("foo"),
			},
			false,
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			got, err := decodeRemoteOptions(test.args)
			isErr := qt.IsNil
			if test.wantErr {
				isErr = qt.IsNotNil
			}

			c.Assert(err, isErr)
			c.Assert(got, qt.DeepEquals, test.want)
		})
	}
}

func TestOptionsNewRequest(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	opts := fromRemoteOptions{
		Method: "GET",
		Body:   []byte("foo"),
	}

	req, err := opts.NewRequest("https://example.com/api")

	c.Assert(err, qt.IsNil)
	c.Assert(req.Method, qt.Equals, "GET")
	c.Assert(req.Header["User-Agent"], qt.DeepEquals, []string{"Hugo Static Site Generator"})

	opts = fromRemoteOptions{
		Method: "GET",
		Body:   []byte("foo"),
		Headers: map[string]any{
			"User-Agent": "foo",
		},
	}

	req, err = opts.NewRequest("https://example.com/api")

	c.Assert(err, qt.IsNil)
	c.Assert(req.Method, qt.Equals, "GET")
	c.Assert(req.Header["User-Agent"], qt.DeepEquals, []string{"foo"})
}

func TestRemoteResourceKeys(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	check := func(uri string, optionsm map[string]any, expect1, expect2 string) {
		c.Helper()
		got1, got2 := remoteResourceKeys(uri, optionsm)
		c.Assert(got1, qt.Equals, expect1)
		c.Assert(got2, qt.Equals, expect2)
	}

	check("foo", nil, "7763396052142361238", "7763396052142361238")
	check("foo", map[string]any{"bar": "baz"}, "5783339285578751849", "5783339285578751849")
	check("foo", map[string]any{"key": "1234", "bar": "baz"}, "15578353952571222948", "5783339285578751849")
	check("foo", map[string]any{"key": "12345", "bar": "baz"}, "14335752410685132726", "5783339285578751849")
	check("asdf", map[string]any{"key": "1234", "bar": "asdf"}, "15578353952571222948", "15615023578599429261")
	check("asdf", map[string]any{"key": "12345", "bar": "asdf"}, "14335752410685132726", "15615023578599429261")
}

// The transport used for remote fetches must refuse to connect to an internal
// (here loopback) address under the default security policy, even though the
// URL text check happens elsewhere. See CVE-2026-10582.
func TestSecureBaseTransportBlocksInternalAddress(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secret"))
	}))
	t.Cleanup(srv.Close)

	req, err := http.NewRequest("GET", srv.URL, nil)
	c.Assert(err, qt.IsNil)

	// Default policy: the loopback address the httptest server listens on is
	// refused at dial time.
	_, err = newSecureBaseTransport(security.DefaultConfig).RoundTrip(req)
	c.Assert(err, qt.IsNotNil)
	c.Assert(security.IsAccessDenied(err), qt.IsTrue)

	// Customized policy: the user has opted into their own hosts, so the dial
	// check stands down and the fetch succeeds.
	sec, err := security.DecodeConfig(config.FromTOMLConfigString(`
[security.http]
urls = ['.*']
`))
	c.Assert(err, qt.IsNil)
	resp, err := newSecureBaseTransport(sec).RoundTrip(req)
	c.Assert(err, qt.IsNil)
	resp.Body.Close()
}
