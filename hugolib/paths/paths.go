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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/modules"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/hugofs"
)

var FilePathSeparator = string(filepath.Separator)

type Paths struct {
	Fs  *hugofs.Fs
	Cfg config.Provider

	BaseURL

	// If the baseURL contains a base path, e.g. https://example.com/docs, then "/docs" will be the BasePath.
	BasePath string

	// Directories
	// TODO(bep) when we have trimmed down mos of the dirs usage outside of this package, make
	// these into an interface.
	ThemesDir  string
	WorkingDir string

	// Directories to store Resource related artifacts.
	AbsResourcesDir string

	AbsPublishDir string

	// pagination path handling
	PaginatePath string

	PublishDir string

	// When in multihost mode, this returns a list of base paths below PublishDir
	// for each language.
	MultihostTargetBasePaths []string

	DisablePathToLower bool
	RemovePathAccents  bool
	UglyURLs           bool
	CanonifyURLs       bool

	Language              *langs.Language
	Languages             langs.Languages
	LanguagesDefaultFirst langs.Languages

	// The PathSpec looks up its config settings in both the current language
	// and then in the global Viper config.
	// Some settings, the settings listed below, does not make sense to be set
	// on per-language-basis. We have no good way of protecting against this
	// other than a "white-list". See language.go.
	defaultContentLanguageInSubdir bool
	DefaultContentLanguage         string
	multilingual                   bool

	AllModules    modules.Modules
	ModulesClient *modules.Client
}

func New(fs *hugofs.Fs, cfg config.Provider) (*Paths, error) {
	baseURLstr := cfg.GetString("baseURL")
	baseURL, err := newBaseURLFromString(baseURLstr)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create baseURL from %q:", baseURLstr)
	}

	contentDir := filepath.Clean(cfg.GetString("contentDir"))
	workingDir := filepath.Clean(cfg.GetString("workingDir"))
	resourceDir := filepath.Clean(cfg.GetString("resourceDir"))
	publishDir := filepath.Clean(cfg.GetString("publishDir"))

	if publishDir == "" {
		return nil, fmt.Errorf("publishDir not set")
	}

	defaultContentLanguage := cfg.GetString("defaultContentLanguage")

	var (
		language              *langs.Language
		languages             langs.Languages
		languagesDefaultFirst langs.Languages
	)

	if l, ok := cfg.(*langs.Language); ok {
		language = l

	}

	if l, ok := cfg.Get("languagesSorted").(langs.Languages); ok {
		languages = l
	}

	if l, ok := cfg.Get("languagesSortedDefaultFirst").(langs.Languages); ok {
		languagesDefaultFirst = l
	}

	//

	if len(languages) == 0 {
		// We have some old tests that does not test the entire chain, hence
		// they have no languages. So create one so we get the proper filesystem.
		languages = langs.Languages{&langs.Language{Lang: "en", Cfg: cfg, ContentDir: contentDir}}
	}

	absPublishDir := AbsPathify(workingDir, publishDir)
	if !strings.HasSuffix(absPublishDir, FilePathSeparator) {
		absPublishDir += FilePathSeparator
	}
	// If root, remove the second '/'
	if absPublishDir == "//" {
		absPublishDir = FilePathSeparator
	}
	absResourcesDir := AbsPathify(workingDir, resourceDir)
	if !strings.HasSuffix(absResourcesDir, FilePathSeparator) {
		absResourcesDir += FilePathSeparator
	}
	if absResourcesDir == "//" {
		absResourcesDir = FilePathSeparator
	}

	var multihostTargetBasePaths []string
	if languages.IsMultihost() {
		for _, l := range languages {
			multihostTargetBasePaths = append(multihostTargetBasePaths, l.Lang)
		}
	}

	p := &Paths{
		Fs:      fs,
		Cfg:     cfg,
		BaseURL: baseURL,

		DisablePathToLower: cfg.GetBool("disablePathToLower"),
		RemovePathAccents:  cfg.GetBool("removePathAccents"),
		UglyURLs:           cfg.GetBool("uglyURLs"),
		CanonifyURLs:       cfg.GetBool("canonifyURLs"),

		ThemesDir:  cfg.GetString("themesDir"),
		WorkingDir: workingDir,

		AbsResourcesDir: absResourcesDir,
		AbsPublishDir:   absPublishDir,

		multilingual:                   cfg.GetBool("multilingual"),
		defaultContentLanguageInSubdir: cfg.GetBool("defaultContentLanguageInSubdir"),
		DefaultContentLanguage:         defaultContentLanguage,

		Language:                 language,
		Languages:                languages,
		LanguagesDefaultFirst:    languagesDefaultFirst,
		MultihostTargetBasePaths: multihostTargetBasePaths,

		PaginatePath: cfg.GetString("paginatePath"),
	}

	if cfg.IsSet("allModules") {
		p.AllModules = cfg.Get("allModules").(modules.Modules)
	}

	if cfg.IsSet("modulesClient") {
		p.ModulesClient = cfg.Get("modulesClient").(*modules.Client)
	}

	// TODO(bep) remove this, eventually
	p.PublishDir = absPublishDir

	return p, nil
}

// GetBasePath returns any path element in baseURL if needed.
func (p *Paths) GetBasePath(isRelativeURL bool) string {
	if isRelativeURL && p.CanonifyURLs {
		// The baseURL will be prepended later.
		return ""
	}
	return p.BasePath
}

func (p *Paths) Lang() string {
	if p == nil || p.Language == nil {
		return ""
	}
	return p.Language.Lang
}

func (p *Paths) GetTargetLanguageBasePath() string {
	if p.Languages.IsMultihost() {
		// In a multihost configuration all assets will be published below the language code.
		return p.Lang()
	}
	return p.GetLanguagePrefix()
}

func (p *Paths) GetURLLanguageBasePath() string {
	if p.Languages.IsMultihost() {
		return ""
	}
	return p.GetLanguagePrefix()
}

func (p *Paths) GetLanguagePrefix() string {
	if !p.multilingual {
		return ""
	}

	defaultLang := p.DefaultContentLanguage
	defaultInSubDir := p.defaultContentLanguageInSubdir

	currentLang := p.Language.Lang
	if currentLang == "" || (currentLang == defaultLang && !defaultInSubDir) {
		return ""
	}
	return currentLang
}

// GetLangSubDir returns the given language's subdir if needed.
func (p *Paths) GetLangSubDir(lang string) string {
	if !p.multilingual {
		return ""
	}

	if p.Languages.IsMultihost() {
		return ""
	}

	if lang == "" || (lang == p.DefaultContentLanguage && !p.defaultContentLanguageInSubdir) {
		return ""
	}

	return lang
}

// AbsPathify creates an absolute path if given a relative path. If already
// absolute, the path is just cleaned.
func (p *Paths) AbsPathify(inPath string) string {
	return AbsPathify(p.WorkingDir, inPath)
}

// RelPathify trims any WorkingDir prefix from the given filename. If
// the filename is not considered to be absolute, the path is just cleaned.
func (p *Paths) RelPathify(filename string) string {
	filename = filepath.Clean(filename)
	if !filepath.IsAbs(filename) {
		return filename
	}

	return strings.TrimPrefix(strings.TrimPrefix(filename, p.WorkingDir), FilePathSeparator)

}

// AbsPathify creates an absolute path if given a working dir and arelative path.
// If already absolute, the path is just cleaned.
func AbsPathify(workingDir, inPath string) string {
	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}
	return filepath.Join(workingDir, inPath)
}
