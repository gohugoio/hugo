// Copyright 2016-present The Hugo Authors. All rights reserved.
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

package helpers

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/cast"
)

// PathSpec holds methods that decides how paths in URLs and files in Hugo should look like.
type PathSpec struct {
	BaseURL

	// If the baseURL contains a base path, e.g. https://example.com/docs, then "/docs" will be the BasePath.
	// This will not be set if canonifyURLs is enabled.
	BasePath string

	disablePathToLower bool
	removePathAccents  bool
	uglyURLs           bool
	canonifyURLs       bool

	Language  *Language
	Languages Languages

	// pagination path handling
	paginatePath string

	theme string

	// Directories
	contentDir     string
	themesDir      string
	layoutDir      string
	workingDir     string
	staticDirs     []string
	absContentDirs []types.KeyValueStr

	PublishDir string

	// The PathSpec looks up its config settings in both the current language
	// and then in the global Viper config.
	// Some settings, the settings listed below, does not make sense to be set
	// on per-language-basis. We have no good way of protecting against this
	// other than a "white-list". See language.go.
	defaultContentLanguageInSubdir bool
	defaultContentLanguage         string
	multilingual                   bool

	ProcessingStats *ProcessingStats

	// The file systems to use
	Fs *hugofs.Fs

	// The fine grained filesystems in play (resources, content etc.).
	BaseFs *hugofs.BaseFs

	// The config provider to use
	Cfg config.Provider
}

func (p PathSpec) String() string {
	return fmt.Sprintf("PathSpec, language %q, prefix %q, multilingual: %T", p.Language.Lang, p.getLanguagePrefix(), p.multilingual)
}

// NewPathSpec creats a new PathSpec from the given filesystems and Language.
func NewPathSpec(fs *hugofs.Fs, cfg config.Provider) (*PathSpec, error) {

	baseURLstr := cfg.GetString("baseURL")
	baseURL, err := newBaseURLFromString(baseURLstr)

	if err != nil {
		return nil, fmt.Errorf("Failed to create baseURL from %q: %s", baseURLstr, err)
	}

	var staticDirs []string

	for i := -1; i <= 10; i++ {
		staticDirs = append(staticDirs, getStringOrStringSlice(cfg, "staticDir", i)...)
	}

	var (
		lang      string
		language  *Language
		languages Languages
	)

	if l, ok := cfg.(*Language); ok {
		language = l
		lang = l.Lang

	}

	if l, ok := cfg.Get("languagesSorted").(Languages); ok {
		languages = l
	}

	defaultContentLanguage := cfg.GetString("defaultContentLanguage")

	// We will eventually pull out this badly placed path logic.
	contentDir := cfg.GetString("contentDir")
	workingDir := cfg.GetString("workingDir")
	resourceDir := cfg.GetString("resourceDir")
	publishDir := cfg.GetString("publishDir")

	if len(languages) == 0 {
		// We have some old tests that does not test the entire chain, hence
		// they have no languages. So create one so we get the proper filesystem.
		languages = Languages{&Language{Lang: "en", ContentDir: contentDir}}
	}

	absPuslishDir := AbsPathify(workingDir, publishDir)
	if !strings.HasSuffix(absPuslishDir, FilePathSeparator) {
		absPuslishDir += FilePathSeparator
	}
	// If root, remove the second '/'
	if absPuslishDir == "//" {
		absPuslishDir = FilePathSeparator
	}
	absResourcesDir := AbsPathify(workingDir, resourceDir)
	if !strings.HasSuffix(absResourcesDir, FilePathSeparator) {
		absResourcesDir += FilePathSeparator
	}
	if absResourcesDir == "//" {
		absResourcesDir = FilePathSeparator
	}

	contentFs, absContentDirs, err := createContentFs(fs.Source, workingDir, defaultContentLanguage, languages)
	if err != nil {
		return nil, err
	}

	// Make sure we don't have any overlapping content dirs. That will never work.
	for i, d1 := range absContentDirs {
		for j, d2 := range absContentDirs {
			if i == j {
				continue
			}
			if strings.HasPrefix(d1.Value, d2.Value) || strings.HasPrefix(d2.Value, d1.Value) {
				return nil, fmt.Errorf("found overlapping content dirs (%q and %q)", d1, d2)
			}
		}
	}

	resourcesFs := afero.NewBasePathFs(fs.Source, absResourcesDir)
	publishFs := afero.NewBasePathFs(fs.Destination, absPuslishDir)

	baseFs := &hugofs.BaseFs{
		ContentFs:   contentFs,
		ResourcesFs: resourcesFs,
		PublishFs:   publishFs,
	}

	ps := &PathSpec{
		Fs:                             fs,
		BaseFs:                         baseFs,
		Cfg:                            cfg,
		disablePathToLower:             cfg.GetBool("disablePathToLower"),
		removePathAccents:              cfg.GetBool("removePathAccents"),
		uglyURLs:                       cfg.GetBool("uglyURLs"),
		canonifyURLs:                   cfg.GetBool("canonifyURLs"),
		multilingual:                   cfg.GetBool("multilingual"),
		Language:                       language,
		Languages:                      languages,
		defaultContentLanguageInSubdir: cfg.GetBool("defaultContentLanguageInSubdir"),
		defaultContentLanguage:         defaultContentLanguage,
		paginatePath:                   cfg.GetString("paginatePath"),
		BaseURL:                        baseURL,
		contentDir:                     contentDir,
		themesDir:                      cfg.GetString("themesDir"),
		layoutDir:                      cfg.GetString("layoutDir"),
		workingDir:                     workingDir,
		staticDirs:                     staticDirs,
		absContentDirs:                 absContentDirs,
		theme:                          cfg.GetString("theme"),
		ProcessingStats:                NewProcessingStats(lang),
	}

	if !ps.canonifyURLs {
		basePath := ps.BaseURL.url.Path
		if basePath != "" && basePath != "/" {
			ps.BasePath = basePath
		}
	}

	// TODO(bep) remove this, eventually
	ps.PublishDir = absPuslishDir

	return ps, nil
}

func getStringOrStringSlice(cfg config.Provider, key string, id int) []string {

	if id >= 0 {
		key = fmt.Sprintf("%s%d", key, id)
	}

	var out []string

	sd := cfg.Get(key)

	if sds, ok := sd.(string); ok {
		out = []string{sds}
	} else if sd != nil {
		out = cast.ToStringSlice(sd)
	}

	return out
}

func createContentFs(fs afero.Fs,
	workingDir,
	defaultContentLanguage string,
	languages Languages) (afero.Fs, []types.KeyValueStr, error) {

	var contentLanguages Languages
	var contentDirSeen = make(map[string]bool)
	languageSet := make(map[string]bool)

	// The default content language needs to be first.
	for _, language := range languages {
		if language.Lang == defaultContentLanguage {
			contentLanguages = append(contentLanguages, language)
			contentDirSeen[language.ContentDir] = true
		}
		languageSet[language.Lang] = true
	}

	for _, language := range languages {
		if contentDirSeen[language.ContentDir] {
			continue
		}
		if language.ContentDir == "" {
			language.ContentDir = defaultContentLanguage
		}
		contentDirSeen[language.ContentDir] = true
		contentLanguages = append(contentLanguages, language)

	}

	var absContentDirs []types.KeyValueStr

	fs, err := createContentOverlayFs(fs, workingDir, contentLanguages, languageSet, &absContentDirs)
	return fs, absContentDirs, err

}

func createContentOverlayFs(source afero.Fs,
	workingDir string,
	languages Languages,
	languageSet map[string]bool,
	absContentDirs *[]types.KeyValueStr) (afero.Fs, error) {
	if len(languages) == 0 {
		return source, nil
	}

	language := languages[0]

	contentDir := language.ContentDir
	if contentDir == "" {
		panic("missing contentDir")
	}

	absContentDir := AbsPathify(workingDir, language.ContentDir)
	if !strings.HasSuffix(absContentDir, FilePathSeparator) {
		absContentDir += FilePathSeparator
	}

	// If root, remove the second '/'
	if absContentDir == "//" {
		absContentDir = FilePathSeparator
	}

	if len(absContentDir) < 6 {
		return nil, fmt.Errorf("invalid content dir %q: %s", absContentDir, ErrPathTooShort)
	}

	*absContentDirs = append(*absContentDirs, types.KeyValueStr{Key: language.Lang, Value: absContentDir})

	overlay := hugofs.NewLanguageFs(language.Lang, languageSet, afero.NewBasePathFs(source, absContentDir))
	if len(languages) == 1 {
		return overlay, nil
	}

	base, err := createContentOverlayFs(source, workingDir, languages[1:], languageSet, absContentDirs)
	if err != nil {
		return nil, err
	}

	return hugofs.NewLanguageCompositeFs(base, overlay), nil

}

// RelContentDir tries to create a path relative to the content root from
// the given filename. The return value is the path and language code.
func (p *PathSpec) RelContentDir(filename string) (string, string) {
	for _, dir := range p.absContentDirs {
		if strings.HasPrefix(filename, dir.Value) {
			rel := strings.TrimPrefix(filename, dir.Value)
			return strings.TrimPrefix(rel, FilePathSeparator), dir.Key
		}
	}
	// Either not a content dir or already relative.
	return filename, ""
}

// ContentDirs returns all the content dirs (absolute paths).
func (p *PathSpec) ContentDirs() []types.KeyValueStr {
	return p.absContentDirs
}

// PaginatePath returns the configured root path used for paginator pages.
func (p *PathSpec) PaginatePath() string {
	return p.paginatePath
}

// ContentDir returns the configured workingDir.
func (p *PathSpec) ContentDir() string {
	return p.contentDir
}

// WorkingDir returns the configured workingDir.
func (p *PathSpec) WorkingDir() string {
	return p.workingDir
}

// StaticDirs returns the relative static dirs for the current configuration.
func (p *PathSpec) StaticDirs() []string {
	return p.staticDirs
}

// LayoutDir returns the relative layout dir in the current configuration.
func (p *PathSpec) LayoutDir() string {
	return p.layoutDir
}

// Theme returns the theme name if set.
func (p *PathSpec) Theme() string {
	return p.theme
}

// Theme returns the theme relative theme dir.
func (p *PathSpec) ThemesDir() string {
	return p.themesDir
}

// PermalinkForBaseURL creates a permalink from the given link and baseURL.
func (p *PathSpec) PermalinkForBaseURL(link, baseURL string) string {
	link = strings.TrimPrefix(link, "/")
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return baseURL + link

}
