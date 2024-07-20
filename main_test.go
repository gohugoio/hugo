// Copyright 2024 The Hugo Authors. All rights reserved.
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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bep/helpers/envhelpers"
	"github.com/gohugoio/hugo/commands"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestCommands(t *testing.T) {
	p := commonTestScriptsParam
	p.Dir = "testscripts/commands"
	testscript.Run(t, p)
}

// Tests in development can be put in "testscripts/unfinished".
// Also see the watch_testscripts.sh script.
func TestUnfinished(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skip unfinished tests on CI")
	}

	p := commonTestScriptsParam
	p.Dir = "testscripts/unfinished"
	// p.UpdateScripts = true

	testscript.Run(t, p)
}

func TestMain(m *testing.M) {
	os.Exit(
		testscript.RunMain(m, map[string]func() int{
			// The main program.
			"hugo": func() int {
				err := commands.Execute(os.Args[1:])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return 1
				}
				return 0
			},
		}),
	)
}

var commonTestScriptsParam = testscript.Params{
	Setup: func(env *testscript.Env) error {
		return testSetupFunc()(env)
	},
	Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
		// log prints to stderr.
		"log": func(ts *testscript.TestScript, neg bool, args []string) {
			log.Println(args)
		},
		// dostounix converts \r\n to \n.
		"dostounix": func(ts *testscript.TestScript, neg bool, args []string) {
			filename := ts.MkAbs(args[0])
			b, err := os.ReadFile(filename)
			if err != nil {
				ts.Fatalf("%v", err)
			}
			b = bytes.Replace(b, []byte("\r\n"), []byte{'\n'}, -1)
			if err := os.WriteFile(filename, b, 0o666); err != nil {
				ts.Fatalf("%v", err)
			}
		},
		// cat prints a file to stdout.
		"cat": func(ts *testscript.TestScript, neg bool, args []string) {
			filename := ts.MkAbs(args[0])
			b, err := os.ReadFile(filename)
			if err != nil {
				ts.Fatalf("%v", err)
			}
			fmt.Print(string(b))
		},
		// sleep sleeps for a second.
		"sleep": func(ts *testscript.TestScript, neg bool, args []string) {
			i := 1
			if len(args) > 0 {
				var err error
				i, err = strconv.Atoi(args[0])
				if err != nil {
					i = 1
				}
			}
			time.Sleep(time.Duration(i) * time.Second)
		},
		// ls lists a directory to stdout.
		"ls": func(ts *testscript.TestScript, neg bool, args []string) {
			dirname := ts.MkAbs(args[0])

			dir, err := os.Open(dirname)
			if err != nil {
				ts.Fatalf("%v", err)
			}
			fis, err := dir.Readdir(-1)
			if err != nil {
				ts.Fatalf("%v", err)
			}
			if len(fis) == 0 {
				// To simplify empty dir checks.
				fmt.Fprintln(ts.Stdout(), "Empty dir")
				return
			}
			for _, fi := range fis {
				fmt.Fprintf(ts.Stdout(), "%s %04o %s %s\n", fi.Mode(), fi.Mode().Perm(), fi.ModTime().Format(time.RFC3339Nano), fi.Name())
			}
		},
		// append appends to a file with a leading newline.
		"append": func(ts *testscript.TestScript, neg bool, args []string) {
			if len(args) < 2 {
				ts.Fatalf("usage: append FILE TEXT")
			}

			filename := ts.MkAbs(args[0])
			words := args[1:]
			for i, word := range words {
				words[i] = strings.Trim(word, "\"")
			}
			text := strings.Join(words, " ")

			_, err := os.Stat(filename)
			if err != nil {
				if os.IsNotExist(err) {
					ts.Fatalf("file does not exist: %s", filename)
				}
				ts.Fatalf("failed to stat file: %v", err)
			}

			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				ts.Fatalf("failed to open file: %v", err)
			}
			defer f.Close()

			_, err = f.WriteString("\n" + text)
			if err != nil {
				ts.Fatalf("failed to write to file: %v", err)
			}
		},
		// replace replaces a string in a file.
		"replace": func(ts *testscript.TestScript, neg bool, args []string) {
			if len(args) < 3 {
				ts.Fatalf("usage: replace FILE OLD NEW")
			}
			filename := ts.MkAbs(args[0])
			oldContent, err := os.ReadFile(filename)
			if err != nil {
				ts.Fatalf("failed to read file %v", err)
			}
			newContent := bytes.Replace(oldContent, []byte(args[1]), []byte(args[2]), -1)
			err = os.WriteFile(filename, newContent, 0o644)
			if err != nil {
				ts.Fatalf("failed to write file: %v", err)
			}
		},

		// httpget checks that a HTTP resource's body matches (if it compiles as a regexp) or contains all of the strings given as arguments.
		"httpget": func(ts *testscript.TestScript, neg bool, args []string) {
			if len(args) < 2 {
				ts.Fatalf("usage: httpgrep URL STRING...")
			}

			tryget := func() error {
				resp, err := http.Get(args[0])
				if err != nil {
					return fmt.Errorf("failed to get URL %q: %v", args[0], err)
				}

				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("failed to read response body: %v", err)
				}
				for _, s := range args[1:] {
					re, err := regexp.Compile(s)
					if err == nil {
						ok := re.Match(body)
						if ok != !neg {
							return fmt.Errorf("response body %q for URL %q does not match %q", body, args[0], s)
						}
					} else {
						ok := bytes.Contains(body, []byte(s))
						if ok != !neg {
							return fmt.Errorf("response body %q for URL %q does not contain %q", body, args[0], s)
						}
					}
				}
				return nil
			}

			// The timing on server rebuilds can be a little tricky to get right,
			// so we try again a few times until the server is ready.
			// There may be smarter ways to do this, but this works.
			start := time.Now()
			for {
				time.Sleep(200 * time.Millisecond)
				err := tryget()
				if err == nil {
					return
				}
				if time.Since(start) > 6*time.Second {
					ts.Fatalf("timeout waiting for %q: %v", args[0], err)
				}
			}
		},
		// checkfile checks that a file exists and is not empty.
		"checkfile": func(ts *testscript.TestScript, neg bool, args []string) {
			var readonly, exec bool
		loop:
			for len(args) > 0 {
				switch args[0] {
				case "-readonly":
					readonly = true
					args = args[1:]
				case "-exec":
					exec = true
					args = args[1:]
				default:
					break loop
				}
			}
			if len(args) == 0 {
				ts.Fatalf("usage: checkfile [-readonly] [-exec] file...")
			}

			for _, filename := range args {
				filename = ts.MkAbs(filename)
				fi, err := os.Stat(filename)
				ok := err == nil != neg
				if !ok {
					ts.Fatalf("stat %s: %v", filename, err)
				}
				if fi.Size() == 0 {
					ts.Fatalf("%s is empty", filename)
				}
				if readonly && fi.Mode()&0o222 != 0 {
					ts.Fatalf("%s is writable", filename)
				}
				if exec && runtime.GOOS != "windows" && fi.Mode()&0o111 == 0 {
					ts.Fatalf("%s is not executable", filename)
				}
			}
		},

		// checkfilecount checks that the number of files in a directory is equal to the given count.
		"checkfilecount": func(ts *testscript.TestScript, neg bool, args []string) {
			if len(args) != 2 {
				ts.Fatalf("usage: checkfilecount count dir")
			}
			count, err := strconv.Atoi(args[0])
			if err != nil {
				ts.Fatalf("invalid count: %v", err)
			}
			if count < 0 {
				ts.Fatalf("count must be non-negative")
			}
			dir := args[1]
			dir = ts.MkAbs(dir)

			found := 0

			filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				found++
				return nil
			})

			ok := found == count != neg
			if !ok {
				ts.Fatalf("found %d files, want %d", found, count)
			}
		},
		// waitServer waits for the .ready file to be created by the server.
		"waitServer": func(ts *testscript.TestScript, neg bool, args []string) {
			type testInfo struct {
				BaseURLs []string
			}

			// The server will write a .ready file when ready.
			// We wait for that.
			readyFilename := ts.MkAbs(".ready")
			limit := time.Now().Add(5 * time.Second)
			for {
				_, err := os.Stat(readyFilename)
				if err != nil {
					time.Sleep(500 * time.Millisecond)
					if time.Now().After(limit) {
						ts.Fatalf("timeout waiting for .ready file")
					}
					continue
				}
				var info testInfo
				// Read the .ready file's JSON into info.
				f, err := os.Open(readyFilename)
				if err != nil {
					ts.Fatalf("failed to open .ready file: %v", err)
				}
				err = json.NewDecoder(f).Decode(&info)
				if err != nil {
					ts.Fatalf("error decoding json: %v", err)
				}
				f.Close()

				for i, s := range info.BaseURLs {
					ts.Setenv(fmt.Sprintf("HUGOTEST_BASEURL_%d", i), s)
				}

				return
			}
		},
		"stopServer": func(ts *testscript.TestScript, neg bool, args []string) {
			baseURL := ts.Getenv("HUGOTEST_BASEURL_0")
			if baseURL == "" {
				ts.Fatalf("HUGOTEST_BASEURL_0 not set")
			}
			if !strings.HasSuffix(baseURL, "/") {
				baseURL += "/"
			}
			resp, err := http.Head(baseURL + "__stop")
			if err != nil {
				ts.Fatalf("failed to shutdown server: %v", err)
			}
			resp.Body.Close()
			// Allow some time for the server to shut down.
			time.Sleep(2 * time.Second)
		},
	},
}

func testSetupFunc() func(env *testscript.Env) error {
	sourceDir, _ := os.Getwd()
	return func(env *testscript.Env) error {
		var keyVals []string
		keyVals = append(keyVals, "HUGO_TESTRUN", "true")
		keyVals = append(keyVals, "HUGO_CACHEDIR", filepath.Join(env.WorkDir, "hugocache"))
		xdghome := filepath.Join(env.WorkDir, "xdgcachehome")
		keyVals = append(keyVals, "XDG_CACHE_HOME", xdghome)
		home := filepath.Join(env.WorkDir, "home")
		keyVals = append(keyVals, "HOME", home)

		if runtime.GOOS == "darwin" {
			if err := os.MkdirAll(filepath.Join(home, "Library", "Caches"), 0o777); err != nil {
				return err
			}
		}

		if runtime.GOOS == "linux" {
			if err := os.MkdirAll(xdghome, 0o777); err != nil {
				return err
			}
		}

		keyVals = append(keyVals, "SOURCE", sourceDir)

		goVersion := runtime.Version()

		goVersion = strings.TrimPrefix(goVersion, "go")
		if strings.HasPrefix(goVersion, "1.20") {
			// Strip patch version.
			goVersion = goVersion[:strings.LastIndex(goVersion, ".")]
		}

		keyVals = append(keyVals, "GOVERSION", goVersion)
		envhelpers.SetEnvVars(&env.Vars, keyVals...)

		return nil
	}
}
