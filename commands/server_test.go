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
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"

	qt "github.com/frankban/quicktest"
)

func TestServer(t *testing.T) {
	c := qt.New(t)

	r := runServerTest(c, true, "")

	c.Assert(r.err, qt.IsNil)
	c.Assert(r.homeContent, qt.Contains, "List: Hugo Commands")
	c.Assert(r.homeContent, qt.Contains, "Environment: development")
}

// Issue 9518
func TestServerPanicOnConfigError(t *testing.T) {
	c := qt.New(t)

	config := `
[markup]
[markup.highlight]
linenos='table'
`

	r := runServerTest(c, false, config)

	c.Assert(r.err, qt.IsNotNil)
	c.Assert(r.err.Error(), qt.Contains, "cannot parse 'Highlight.LineNos' as bool:")
}

func TestServerFlags(t *testing.T) {
	c := qt.New(t)

	assertPublic := func(c *qt.C, r serverTestResult, renderStaticToDisk bool) {
		c.Assert(r.err, qt.IsNil)
		c.Assert(r.homeContent, qt.Contains, "Environment: development")
		c.Assert(r.publicDirnames["myfile.txt"], qt.Equals, renderStaticToDisk)

	}

	for _, test := range []struct {
		flag   string
		assert func(c *qt.C, r serverTestResult)
	}{
		{"", func(c *qt.C, r serverTestResult) {
			assertPublic(c, r, false)
		}},
		{"--renderToDisk", func(c *qt.C, r serverTestResult) {
			assertPublic(c, r, true)
		}},
	} {
		c.Run(test.flag, func(c *qt.C) {
			config := `
baseURL="https://example.org"
`

			var args []string
			if test.flag != "" {
				args = strings.Split(test.flag, "=")
			}

			r := runServerTest(c, true, config, args...)

			test.assert(c, r)

		})

	}

}

type serverTestResult struct {
	err            error
	homeContent    string
	publicDirnames map[string]bool
}

func runServerTest(c *qt.C, getHome bool, config string, args ...string) (result serverTestResult) {
	dir, clean, err := createSimpleTestSite(c, testSiteConfig{configTOML: config})
	defer clean()
	c.Assert(err, qt.IsNil)

	sp, err := helpers.FindAvailablePort()
	c.Assert(err, qt.IsNil)
	port := sp.Port

	defer func() {
		os.RemoveAll(dir)
	}()

	stop := make(chan bool)

	b := newCommandsBuilder()
	scmd := b.newServerCmdSignaled(stop)

	cmd := scmd.getCommand()
	args = append([]string{"-s=" + dir, fmt.Sprintf("-p=%d", port)}, args...)
	cmd.SetArgs(args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		_, err := cmd.ExecuteC()
		return err
	})

	if getHome {
		// Esp. on slow CI machines, we need to wait a little before the web
		// server is ready.
		time.Sleep(567 * time.Millisecond)
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
		c.Check(err, qt.IsNil)
		if err == nil {
			defer resp.Body.Close()
			result.homeContent = helpers.ReaderToString(resp.Body)
		}
	}

	select {
	case <-stop:
	case stop <- true:
	}

	pubFiles, err := os.ReadDir(filepath.Join(dir, "public"))
	c.Check(err, qt.IsNil)
	result.publicDirnames = make(map[string]bool)
	for _, f := range pubFiles {
		result.publicDirnames[f.Name()] = true
	}

	result.err = wg.Wait()

	return

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
			v := config.New()
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
