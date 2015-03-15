// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpl

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
)

func TestScpCache(t *testing.T) {

	tests := []struct {
		path    string
		content []byte
		ignore  bool
	}{
		{"http://Foo.Bar/foo_Bar-Foo", []byte(`T€st Content 123`), false},
		{"fOO,bar:foo%bAR", []byte(`T€st Content 123 fOO,bar:foo%bAR`), false},
		{"FOo/BaR.html", []byte(`FOo/BaR.html T€st Content 123`), false},
		{"трям/трям", []byte(`T€st трям/трям Content 123`), false},
		{"은행", []byte(`T€st C은행ontent 123`), false},
		{"Банковский кассир", []byte(`Банковский кассир T€st Content 123`), false},
		{"Банковский кассир", []byte(`Банковский кассир T€st Content 456`), true},
	}

	fs := new(afero.MemMapFs)

	for _, test := range tests {
		c, err := resGetCache(test.path, fs, test.ignore)
		if err != nil {
			t.Errorf("Error getting cache: %s", err)
		}
		if c != nil {
			t.Errorf("There is content where there should not be anything: %s", string(c))
		}

		err = resWriteCache(test.path, test.content, fs)
		if err != nil {
			t.Errorf("Error writing cache: %s", err)
		}

		c, err = resGetCache(test.path, fs, test.ignore)
		if err != nil {
			t.Errorf("Error getting cache after writing: %s", err)
		}
		if test.ignore {
			if c != nil {
				t.Errorf("Cache ignored but content is not nil: %s", string(c))
			}
		} else {
			if bytes.Compare(c, test.content) != 0 {
				t.Errorf("\nExpected: %s\nActual: %s\n", string(test.content), string(c))
			}
		}
	}
}

func TestScpGetLocal(t *testing.T) {
	fs := new(afero.MemMapFs)
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
		err := helpers.WriteToDisk(test.path, r, fs)
		if err != nil {
			t.Error(err)
		}

		c, err := resGetLocal(test.path, fs)
		if err != nil {
			t.Errorf("Error getting resource content: %s", err)
		}
		if bytes.Compare(c, test.content) != 0 {
			t.Errorf("\nExpected: %s\nActual: %s\n", string(test.content), string(c))
		}
	}

}

func getTestServer(handler func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *http.Client) {
	testServer := httptest.NewServer(http.HandlerFunc(handler))
	client := &http.Client{
		Transport: &http.Transport{Proxy: func(*http.Request) (*url.URL, error) { return url.Parse(testServer.URL) }},
	}
	return testServer, client
}

func TestScpGetRemote(t *testing.T) {
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

		srv, cl := getTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.Write(test.content)
		})
		defer func() { srv.Close() }()

		c, err := resGetRemote(test.path, fs, cl)
		if err != nil {
			t.Errorf("Error getting resource content: %s", err)
		}
		if bytes.Compare(c, test.content) != 0 {
			t.Errorf("\nNet Expected: %s\nNet Actual: %s\n", string(test.content), string(c))
		}
		cc, cErr := resGetCache(test.path, fs, test.ignore)
		if cErr != nil {
			t.Error(cErr)
		}
		if test.ignore {
			if cc != nil {
				t.Errorf("Cache ignored but content is not nil: %s", string(cc))
			}
		} else {
			if bytes.Compare(cc, test.content) != 0 {
				t.Errorf("\nCache Expected: %s\nCache Actual: %s\n", string(test.content), string(cc))
			}
		}
	}
}

func TestParseCsv(t *testing.T) {

	tests := []struct {
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
	}
	for _, test := range tests {
		csv, err := parseCsv(test.csv, test.sep)
		if test.err && err == nil {
			t.Error("Expecting an error")
		}
		if test.err {
			continue
		}
		if !test.err && err != nil {
			t.Error(err)
		}

		act := ""
		for _, v := range csv {
			act = act + strings.Join(v, "")
		}

		if act != test.exp {
			t.Errorf("\nExpected: %s\nActual: %s\n%#v\n", test.exp, act, csv)
		}

	}
}
