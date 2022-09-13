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

// Issue 9518
func TestServerPanicOnConfigError(t *testing.T) {
	c := qt.New(t)

	config := `
[markup]
[markup.highlight]
linenos='table'
`

	r := runServerTest(c,
		serverTestOptions{
			config: config,
		},
	)

	c.Assert(r.err, qt.IsNotNil)
	c.Assert(r.err.Error(), qt.Contains, "cannot parse 'Highlight.LineNos' as bool:")
}

func TestServer404(t *testing.T) {
	c := qt.New(t)

	r := runServerTest(c,
		serverTestOptions{
			test404:     true,
			getNumHomes: 1,
		},
	)

	c.Assert(r.err, qt.IsNil)
	c.Assert(r.content404, qt.Contains, "404: 404 Page not found|Not Found.")
}

func TestServerFlags(t *testing.T) {
	c := qt.New(t)

	assertPublic := func(c *qt.C, r serverTestResult, renderStaticToDisk bool) {
		c.Assert(r.err, qt.IsNil)
		c.Assert(r.homesContent[0], qt.Contains, "Environment: development")
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
		{"--renderStaticToDisk", func(c *qt.C, r serverTestResult) {
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

			opts := serverTestOptions{
				config:      config,
				args:        args,
				getNumHomes: 1,
			}

			r := runServerTest(c, opts)

			test.assert(c, r)

		})

	}

}

func TestServerBugs(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		name       string
		config     string
		flag       string
		numservers int
		assert     func(c *qt.C, r serverTestResult)
	}{
		// Issue 9788
		{"PostProcess, memory", "", "", 1, func(c *qt.C, r serverTestResult) {
			c.Assert(r.err, qt.IsNil)
			c.Assert(r.homesContent[0], qt.Contains, "PostProcess: /foo.min.css")
		}},
		{"PostProcess, disk", "", "--renderToDisk", 1, func(c *qt.C, r serverTestResult) {
			c.Assert(r.err, qt.IsNil)
			c.Assert(r.homesContent[0], qt.Contains, "PostProcess: /foo.min.css")
		}},
		// Isue 9901
		{"Multihost", `
defaultContentLanguage = 'en'
[languages]
[languages.en]
baseURL = 'https://example.com'
title = 'My blog'
weight = 1
[languages.fr]
baseURL = 'https://example.fr'
title = 'Mon blogue'
weight = 2
`, "", 2, func(c *qt.C, r serverTestResult) {
			c.Assert(r.err, qt.IsNil)
			for i, s := range []string{"My blog", "Mon blogue"} {
				c.Assert(r.homesContent[i], qt.Contains, s)
			}
		}},
	} {
		c.Run(test.name, func(c *qt.C) {
			if test.config == "" {
				test.config = `
baseURL="https://example.org"
`
			}

			var args []string
			if test.flag != "" {
				args = strings.Split(test.flag, "=")
			}

			opts := serverTestOptions{
				config:      test.config,
				getNumHomes: test.numservers,
				test404:     true,
				args:        args,
			}

			r := runServerTest(c, opts)
			c.Assert(r.content404, qt.Contains, "404: 404 Page not found|Not Found.")
			test.assert(c, r)

		})

	}

}

type serverTestResult struct {
	err            error
	homesContent   []string
	content404     string
	publicDirnames map[string]bool
}

type serverTestOptions struct {
	getNumHomes int
	test404     bool
	config      string
	args        []string
}

func runServerTest(c *qt.C, opts serverTestOptions) (result serverTestResult) {
	dir := createSimpleTestSite(c, testSiteConfig{configTOML: opts.config})

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
	args := append([]string{"-s=" + dir, fmt.Sprintf("-p=%d", port)}, opts.args...)
	cmd.SetArgs(args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		_, err := cmd.ExecuteC()
		return err
	})

	if opts.getNumHomes > 0 {
		// Esp. on slow CI machines, we need to wait a little before the web
		// server is ready.
		time.Sleep(567 * time.Millisecond)
		result.homesContent = make([]string, opts.getNumHomes)
		for i := 0; i < opts.getNumHomes; i++ {
			func() {
				resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port+i))
				c.Check(err, qt.IsNil)
				c.Check(resp.StatusCode, qt.Equals, http.StatusOK)
				if err == nil {
					defer resp.Body.Close()
					result.homesContent[i] = helpers.ReaderToString(resp.Body)
				}
			}()
		}
	}

	if opts.test404 {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/this-page-does-not-exist", port))
		c.Check(err, qt.IsNil)
		c.Check(resp.StatusCode, qt.Equals, http.StatusNotFound)
		if err == nil {
			defer resp.Body.Close()
			result.content404 = helpers.ReaderToString(resp.Body)
		}
	}

	time.Sleep(1 * time.Second)

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
			v := config.NewWithTestDefaults()
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
