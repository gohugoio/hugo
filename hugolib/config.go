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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// config file items
type Config struct {
	SourceDir, PublishDir, BaseUrl, StaticDir string
	Path, CacheDir, LayoutDir, DefaultLayout  string
	Indexes                                   map[string]string // singular, plural
	ProcessFilters                            map[string][]string
	BuildDrafts                               bool
}

var c Config

// Read cfgfile or setup defaults.
func SetupConfig(cfgfile *string, path *string) *Config {
	c.setPath(*path)

	configPath, err := c.findConfigFile(*cfgfile)

	if err != nil {
		fmt.Printf("%v", err)
		fmt.Println(" using defaults instead")
	}

	// set defaults

	c.SourceDir = "content"
	c.LayoutDir = "layouts"
	c.PublishDir = "public"
	c.StaticDir = "static"
	c.DefaultLayout = "post"
	c.BuildDrafts = false

	file, err := ioutil.ReadFile(configPath)
	if err == nil {
		if err := json.Unmarshal(file, &c); err != nil {
			fmt.Printf("Error parsing config: %s", err)
			os.Exit(1)
		}
	}

	// set index defaults if none provided
	if len(c.Indexes) == 0 {
		c.Indexes = make(map[string]string)
		c.Indexes["tag"] = "tags"
		c.Indexes["category"] = "categories"
	}
	return &c
}

func (c *Config) setPath(p string) {
	if p == "" {
		path, err := FindPath()
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

func FindPath() (string, error) {
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

func (c *Config) GetAbsPath(name string) string {
	if path.IsAbs(name) {
		return name
	}

	p := filepath.Join(c.GetPath(), name)
	return p
}

func (c *Config) findConfigFile(configFileName string) (string, error) {
	// If the full path is given, just use that
	if path.IsAbs(configFileName) {
		return configFileName, nil
	}

	// Else check the local directory
	t := c.GetAbsPath(configFileName)
	if b, _ := exists(t); b {
		return t, nil
	} else {
		return "", fmt.Errorf("config file not found at: %s", t)
	}

	return "", nil // This line won't ever happen.. looking forward to go 1.1 when I don't need it
}
