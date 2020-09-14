// Copyright 2018 The Hugo Authors. All rights reserved.
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

package hugo

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
)

const (
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

var (
	// commitHash contains the current Git revision. Use make to build to make
	// sure this gets set.
	commitHash string

	// buildDate contains the date of the current build.
	buildDate string
)

// Info contains information about the current Hugo environment
type Info struct {
	CommitHash string
	BuildDate  string

	// The build environment.
	// Defaults are "production" (hugo) and "development" (hugo server).
	// This can also be set by the user.
	// It can be any string, but it will be all lower case.
	Environment string
}

// Version returns the current version as a comparable version string.
func (i Info) Version() VersionString {
	return CurrentVersion.Version()
}

// Generator a Hugo meta generator HTML tag.
func (i Info) Generator() template.HTML {
	return template.HTML(fmt.Sprintf(`<meta name="generator" content="Hugo %s" />`, CurrentVersion.String()))
}

func (i Info) IsProduction() bool {
	return i.Environment == EnvironmentProduction
}

// NewInfo creates a new Hugo Info object.
func NewInfo(environment string) Info {
	if environment == "" {
		environment = EnvironmentProduction
	}
	return Info{
		CommitHash:  commitHash,
		BuildDate:   buildDate,
		Environment: environment,
	}
}

func GetExecEnviron(workDir string, cfg config.Provider, fs afero.Fs) []string {
	env := os.Environ()
	nodepath := filepath.Join(workDir, "node_modules")
	if np := os.Getenv("NODE_PATH"); np != "" {
		nodepath = workDir + string(os.PathListSeparator) + np
	}
	config.SetEnvVars(&env, "NODE_PATH", nodepath)
	config.SetEnvVars(&env, "PWD", workDir)
	config.SetEnvVars(&env, "HUGO_ENVIRONMENT", cfg.GetString("environment"))
	fis, err := afero.ReadDir(fs, files.FolderJSConfig)
	if err == nil {
		for _, fi := range fis {
			key := fmt.Sprintf("HUGO_FILE_%s", strings.ReplaceAll(strings.ToUpper(fi.Name()), ".", "_"))
			value := fi.(hugofs.FileMetaInfo).Meta().Filename()
			config.SetEnvVars(&env, key, value)
		}
	}

	return env
}
