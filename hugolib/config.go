// Copyright © 2013 Steve Francia <spf@spf13.com>.
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

package hugolib

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/spf13/hugo/helpers"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// config file items
type Config struct {
	ContentDir, PublishDir, BaseUrl, StaticDir string
	Path, CacheDir, LayoutDir, DefaultLayout   string
	ConfigFile                                 string
	Title                                      string
	Indexes                                    map[string]string // singular, plural
	ProcessFilters                             map[string][]string
	Params                                     map[string]interface{}
	Permalinks                                 PermalinkOverrides
	BuildDrafts, UglyUrls, Verbose             bool
	CanonifyUrls                               bool
}

var c Config

// Read cfgfile or setup defaults.
func SetupConfig(cfgfile *string, path *string) *Config {
	c.setPath(*path)

	cfg, err := c.findConfigFile(*cfgfile)
	c.ConfigFile = cfg

	if err != nil {
		fmt.Printf("%v", err)
		fmt.Println(" using defaults instead")
	}

	// set defaults
	c.ContentDir = "content"
	c.LayoutDir = "layouts"
	c.PublishDir = "public"
	c.StaticDir = "static"
	c.DefaultLayout = "post"
	c.BuildDrafts = false
	c.UglyUrls = false
	c.Verbose = false
	c.CanonifyUrls = true

	c.readInConfig()

	// set index defaults if none provided
	if c.Indexes == nil {
		c.Indexes = make(map[string]string)
		c.Indexes["tag"] = "tags"
		c.Indexes["category"] = "categories"
	}

	// ensure map exists, albeit empty
	if c.Permalinks == nil {
		c.Permalinks = make(PermalinkOverrides, 0)
	}

	if !strings.HasSuffix(c.BaseUrl, "/") {
		c.BaseUrl = c.BaseUrl + "/"
	}

	return &c
}

func (c *Config) readInConfig() {
	file, err := ioutil.ReadFile(c.ConfigFile)
	if err == nil {
		switch path.Ext(c.ConfigFile) {
		case ".yaml":
			if err := goyaml.Unmarshal(file, &c); err != nil {
				fmt.Printf("Error parsing config: %s", err)
				os.Exit(1)
			}

		case ".json":
			if err := json.Unmarshal(file, &c); err != nil {
				fmt.Printf("Error parsing config: %s", err)
				os.Exit(1)
			}

		case ".toml":
			if _, err := toml.Decode(string(file), &c); err != nil {
				fmt.Printf("Error parsing config: %s", err)
				os.Exit(1)
			}
		}
	}
}

func (c *Config) setPath(p string) {
	if p == "" {
		path, err := findPath()
		if err != nil {
			fmt.Printf("Error finding path: %s", err)
		}
		c.Path = path
	} else {
		path, err := filepath.Abs(p)
		if err != nil {
			fmt.Printf("Error finding path: %s", err)
		}
		c.Path = path
	}
}

func (c *Config) GetPath() string {
	if c.Path == "" {
		c.setPath("")
	}
	return c.Path
}

func findPath() (string, error) {
	serverFile, err := filepath.Abs(os.Args[0])

	if err != nil {
		return "", fmt.Errorf("Can't get absolute path for executable: %v", err)
	}

	path := filepath.Dir(serverFile)
	realFile, err := filepath.EvalSymlinks(serverFile)

	if err != nil {
		if _, err = os.Stat(serverFile + ".exe"); err == nil {
			realFile = filepath.Clean(serverFile + ".exe")
		}
	}

	if err == nil && realFile != serverFile {
		path = filepath.Dir(realFile)
	}

	return path, nil
}

// GetAbsPath return the absolute path for a given path with the internal slashes
// properly converted.
func (c *Config) GetAbsPath(name string) string {
	if filepath.IsAbs(name) {
		return name
	}

	return filepath.Join(c.GetPath(), name)
}

func (c *Config) findConfigFile(configFileName string) (string, error) {

	if configFileName == "" { // config not specified, let's search
		if b, _ := helpers.Exists(c.GetAbsPath("config.json")); b {
			return c.GetAbsPath("config.json"), nil
		}

		if b, _ := helpers.Exists(c.GetAbsPath("config.toml")); b {
			return c.GetAbsPath("config.toml"), nil
		}

		if b, _ := helpers.Exists(c.GetAbsPath("config.yaml")); b {
			return c.GetAbsPath("config.yaml"), nil
		}

		return "", fmt.Errorf("config file not found in: %s", c.GetPath())

	} else {
		// If the full path is given, just use that
		if path.IsAbs(configFileName) {
			return configFileName, nil
		}

		// Else check the local directory
		t := c.GetAbsPath(configFileName)
		if b, _ := helpers.Exists(t); b {
			return t, nil
		} else {
			return "", fmt.Errorf("config file not found at: %s", t)
		}
	}
}
