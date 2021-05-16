// Copyright 2015 The Hugo Authors. All rights reserved.
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

package commands

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/helpers"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/viper"
)

func TestServer(t *testing.T) {
	if isWindowsCI() {
		// TODO(bep) not sure why server tests have started to fail on the Windows CI server.
		t.Skip("Skip server test on appveyor")
	}
	c := qt.New(t)
	dir, clean, err := createSimpleTestSite(t, testSiteConfig{})
	defer clean()
	c.Assert(err, qt.IsNil)

	// Let us hope that this port is available on all systems ...
	port := 1313

	defer func() {
		os.RemoveAll(dir)
	}()

	stop := make(chan bool)

	b := newCommandsBuilder()
	scmd := b.newServerCmdSignaled(stop)

	cmd := scmd.getCommand()
	cmd.SetArgs([]string{"-s=" + dir, fmt.Sprintf("-p=%d", port)})

	go func() {
		_, err = cmd.ExecuteC()
		c.Assert(err, qt.IsNil)
	}()

	// There is no way to know exactly when the server is ready for connections.
	// We could improve by something like https://golang.org/pkg/net/http/httptest/#Server
	// But for now, let us sleep and pray!
	time.Sleep(2 * time.Second)

	checkHttpResp(c, t, "http://localhost:1313/", "List: Hugo Commands")
	checkHttpResp(c, t, "http://localhost:1313/", "Environment: development")

	// Stop the server.
	stop <- true
}

func TestServerRenderStaticToDisk(t *testing.T) {
	// Basically same as TestServer
	if isWindowsCI() {
		t.Skip("Skip server test on appveyor")
	}
	c := qt.New(t)
	dir, clean, err := createSimpleTestSite(t, testSiteConfig{})

	contentDir := "content"
	staticDir := "static"
	writeFile(t, filepath.Join(dir, contentDir, "content.csv"), "static data in contentDir")
	writeFile(t, filepath.Join(dir, staticDir, "static.csv"), "static data in staticDir")
	defer clean()
	c.Assert(err, qt.IsNil)

	// If server tests run in pararrel, port 1313 is blocked.
	// Expect that this port is available.
	port := 1333

	defer func() {
		os.RemoveAll(dir)
	}()

	stop := make(chan bool)

	b := newCommandsBuilder()
	scmd := b.newServerCmdSignaled(stop)

	cmd := scmd.getCommand()
	cmd.SetArgs([]string{"-s=" + dir, fmt.Sprintf("-p=%d", port), "--renderStaticToDisk"})

	go func() {
		_, err = cmd.ExecuteC()
		c.Assert(err, qt.IsNil)
	}()

	time.Sleep(2 * time.Second)

	checkHttpResp(c, t, "http://localhost:1333/", "List: Hugo Commands")
	checkHttpResp(c, t, "http://localhost:1333/", "Environment: development")

	// After inital build, in public,
	// "static.csv" and "content.csv" exist
	// "p1.md"-related html does NOT exist
	files, err := ioutil.ReadDir(dir + "/public")
	c.Assert(err, qt.IsNil)
	for _, f := range files {
		c.Assert([]string{"static.csv", "content.csv"}, qt.Any(qt.Equals), f.Name())
	}

	checkHttpResp(c, t, "http://localhost:1333/content.csv", "static data in contentDir")
	checkHttpResp(c, t, "http://localhost:1333/static.csv", "static data in staticDir")

	writeFile(t, filepath.Join(dir, contentDir, "content.csv"), "static data in contentDir can be updated")
	writeFile(t, filepath.Join(dir, staticDir, "static.csv"), "static data in staticDir can be updated")

	writeFile(t, filepath.Join(dir, contentDir, "added/content.csv"), "static data in contentDir can be added")
	writeFile(t, filepath.Join(dir, staticDir, "added/static.csv"), "static data in staticDir can be added")

	writeFile(t, filepath.Join(dir, contentDir, "conflicted.csv"), "if file name is conflicted, data in contentDir is served")
	writeFile(t, filepath.Join(dir, staticDir, "conflicted.csv"), "if file name is conflicted, data in staticDir is served")

	time.Sleep(2 * time.Second)

	checkHttpResp(c, t, "http://localhost:1333/content.csv", "static data in contentDir can be updated")
	checkHttpResp(c, t, "http://localhost:1333/static.csv", "static data in staticDir can be updated")

	checkHttpResp(c, t, "http://localhost:1333/added/content.csv", "static data in contentDir can be added")
	checkHttpResp(c, t, "http://localhost:1333/added/static.csv", "static data in staticDir can be added")

	checkHttpResp(c, t, "http://localhost:1333/conflicted.csv", "if file name is conflicted, data in contentDir is served")

	files, err = ioutil.ReadDir(dir + "/public")
	c.Assert(err, qt.IsNil)
	for _, f := range files {
		c.Assert([]string{"static.csv", "content.csv", "conflicted.csv", "added"}, qt.Any(qt.Equals), f.Name())
	}

	// can be removed
	os.RemoveAll(filepath.Join(dir, staticDir, "added"))
	os.RemoveAll(filepath.Join(dir, contentDir, "added"))

	time.Sleep(1 * time.Second)

	checkHttpResp(c, t, "http://localhost:1333/added/content.csv", "404 page not found")
	checkHttpResp(c, t, "http://localhost:1333/added/static.csv", "404 page not found")

	files, err = ioutil.ReadDir(dir + "/public")
	c.Assert(err, qt.IsNil)
	for _, f := range files {
		c.Assert([]string{"static.csv", "content.csv", "conflicted.csv"}, qt.Any(qt.Equals), f.Name())
	}

	stop <- true
}

func checkHttpResp(c *qt.C, t *testing.T, url string, part string) {
	resp, err := http.Get(url)
	c.Assert(err, qt.IsNil)
	defer resp.Body.Close()
	content := helpers.ReaderToString(resp.Body)
	c.Assert(content, qt.Contains, part)
}

func TestFixURL(t *testing.T) {
	type data struct {
		TestName   string
		CLIBaseURL string
		CfgBaseURL string
		AppendPort bool
		Port       int
		Result     string
	}
	tests := []data{
		{"Basic http localhost", "", "http://foo.com", true, 1313, "http://localhost:1313/"},
		{"Basic https production, http localhost", "", "https://foo.com", true, 1313, "http://localhost:1313/"},
		{"Basic subdir", "", "http://foo.com/bar", true, 1313, "http://localhost:1313/bar/"},
		{"Basic production", "http://foo.com", "http://foo.com", false, 80, "http://foo.com/"},
		{"Production subdir", "http://foo.com/bar", "http://foo.com/bar", false, 80, "http://foo.com/bar/"},
		{"No http", "", "foo.com", true, 1313, "//localhost:1313/"},
		{"Override configured port", "", "foo.com:2020", true, 1313, "//localhost:1313/"},
		{"No http production", "foo.com", "foo.com", false, 80, "//foo.com/"},
		{"No http production with port", "foo.com", "foo.com", true, 2020, "//foo.com:2020/"},
		{"No config", "", "", true, 1313, "//localhost:1313/"},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			b := newCommandsBuilder()
			s := b.newServerCmd()
			v := viper.New()
			baseURL := test.CLIBaseURL
			v.Set("baseURL", test.CfgBaseURL)
			s.serverAppend = test.AppendPort
			s.serverPort = test.Port
			result, err := s.fixURL(v, baseURL, s.serverPort)
			if err != nil {
				t.Errorf("Unexpected error %s", err)
			}
			if result != test.Result {
				t.Errorf("Expected %q, got %q", test.Result, result)
			}
		})
	}
}

func TestRemoveErrorPrefixFromLog(t *testing.T) {
	c := qt.New(t)
	content := `ERROR 2018/10/07 13:11:12 Error while rendering "home": template: _default/baseof.html:4:3: executing "main" at <partial "logo" .>: error calling partial: template: partials/logo.html:5:84: executing "partials/logo.html" at <$resized.AHeight>: can't evaluate field AHeight in type *resource.Image
ERROR 2018/10/07 13:11:12 Rebuild failed: logged 1 error(s)
`

	withoutError := removeErrorPrefixFromLog(content)

	c.Assert(strings.Contains(withoutError, "ERROR"), qt.Equals, false)
}

func isWindowsCI() bool {
	return runtime.GOOS == "windows" && os.Getenv("CI") != ""
}
