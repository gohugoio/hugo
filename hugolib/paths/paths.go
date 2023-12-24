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

package paths

import (
	"path/filepath"
	"strings"

	hpaths "github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/modules"

	"github.com/gohugoio/hugo/hugofs"
)

var FilePathSeparator = string(filepath.Separator)

type Paths struct {
	Fs  *hugofs.Fs
	Cfg config.AllProvider

	// Directories to store Resource related artifacts.
	AbsResourcesDir string

	AbsPublishDir string

	// When in multihost mode, this returns a list of base paths below PublishDir
	// for each language.
	MultihostTargetBasePaths []string
}

func New(fs *hugofs.Fs, cfg config.AllProvider) (*Paths, error) {
	bcfg := cfg.BaseConfig()
	publishDir := bcfg.PublishDir
	if publishDir == "" {
		panic("publishDir not set")
	}

	absPublishDir := hpaths.AbsPathify(bcfg.WorkingDir, publishDir)
	if !strings.HasSuffix(absPublishDir, FilePathSeparator) {
		absPublishDir += FilePathSeparator
	}
	// If root, remove the second '/'
	if absPublishDir == "//" {
		absPublishDir = FilePathSeparator
	}
	absResourcesDir := hpaths.AbsPathify(bcfg.WorkingDir, cfg.Dirs().ResourceDir)
	if !strings.HasSuffix(absResourcesDir, FilePathSeparator) {
		absResourcesDir += FilePathSeparator
	}
	if absResourcesDir == "//" {
		absResourcesDir = FilePathSeparator
	}

	var multihostTargetBasePaths []string
	if cfg.IsMultihost() && len(cfg.Languages()) > 1 {
		for _, l := range cfg.Languages() {
			multihostTargetBasePaths = append(multihostTargetBasePaths, l.Lang)
		}
	}

	p := &Paths{
		Fs:                       fs,
		Cfg:                      cfg,
		AbsResourcesDir:          absResourcesDir,
		AbsPublishDir:            absPublishDir,
		MultihostTargetBasePaths: multihostTargetBasePaths,
	}

	return p, nil
}

func (p *Paths) AllModules() modules.Modules {
	return p.Cfg.GetConfigSection("allModules").(modules.Modules)
}

// GetBasePath returns any path element in baseURL if needed.
// The path returned will have a leading, but no trailing slash.
func (p *Paths) GetBasePath(isRelativeURL bool) string {
	if isRelativeURL && p.Cfg.CanonifyURLs() {
		// The baseURL will be prepended later.
		return ""
	}
	return p.Cfg.BaseURL().BasePathNoTrailingSlash
}

func (p *Paths) Lang() string {
	if p == nil || p.Cfg.Language() == nil {
		return ""
	}
	return p.Cfg.Language().Lang
}

func (p *Paths) GetTargetLanguageBasePath() string {
	if p.Cfg.IsMultihost() {
		// In a multihost configuration all assets will be published below the language code.
		return p.Lang()
	}
	return p.GetLanguagePrefix()
}

func (p *Paths) GetLanguagePrefix() string {
	return p.Cfg.LanguagePrefix()
}

// AbsPathify creates an absolute path if given a relative path. If already
// absolute, the path is just cleaned.
func (p *Paths) AbsPathify(inPath string) string {
	return hpaths.AbsPathify(p.Cfg.BaseConfig().WorkingDir, inPath)
}

// RelPathify trims any WorkingDir prefix from the given filename. If
// the filename is not considered to be absolute, the path is just cleaned.
func (p *Paths) RelPathify(filename string) string {
	filename = filepath.Clean(filename)
	if !filepath.IsAbs(filename) {
		return filename
	}

	return strings.TrimPrefix(strings.TrimPrefix(filename, p.Cfg.BaseConfig().WorkingDir), FilePathSeparator)
}
