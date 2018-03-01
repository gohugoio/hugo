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

// Package filesystems provides the fine grained file systems used by Hugo. These
// are typically virtual filesystems that are composites of project and theme content.
package filesystems

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/hugofs"

	"fmt"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugolib/paths"
	"github.com/gohugoio/hugo/langs"
	"github.com/spf13/afero"
)

// When we create a virtual filesystem with data and i18n bundles for the project and the themes,
// this is the name of the project's virtual root. It got it's funky name to make sure
// (or very unlikely) that it collides with a theme name.
const projectVirtualFolder = "__h__project"

var filePathSeparator = string(filepath.Separator)

// BaseFs contains the core base filesystems used by Hugo. The name "base" is used
// to underline that even if they can be composites, they all have a base path set to a specific
// resource folder, e.g "/my-project/content". So, no absolute filenames needed.
type BaseFs struct {
	// TODO(bep) make this go away
	AbsContentDirs []types.KeyValueStr

	// The filesystem used to capture content. This can be a composite and
	// language aware file system.
	ContentFs afero.Fs

	// SourceFilesystems contains the different source file systems.
	*SourceFilesystems

	// The filesystem used to store resources (processed images etc.).
	// This usually maps to /my-project/resources.
	ResourcesFs afero.Fs

	// The filesystem used to publish the rendered site.
	// This usually maps to /my-project/public.
	PublishFs afero.Fs

	themeFs afero.Fs

	// TODO(bep) improve the "theme interaction"
	AbsThemeDirs []string
}

// RelContentDir tries to create a path relative to the content root from
// the given filename. The return value is the path and language code.
func (b *BaseFs) RelContentDir(filename string) (string, string) {
	for _, dir := range b.AbsContentDirs {
		if strings.HasPrefix(filename, dir.Value) {
			rel := strings.TrimPrefix(filename, dir.Value)
			return strings.TrimPrefix(rel, filePathSeparator), dir.Key
		}
	}
	// Either not a content dir or already relative.
	return filename, ""
}

// IsContent returns whether the given filename is in the content filesystem.
func (b *BaseFs) IsContent(filename string) bool {
	for _, dir := range b.AbsContentDirs {
		if strings.HasPrefix(filename, dir.Value) {
			return true
		}
	}
	return false
}

// SourceFilesystems contains the different source file systems. These can be
// composite file systems (theme and project etc.), and they have all root
// set to the source type the provides: data, i18n, static, layouts.
type SourceFilesystems struct {
	Data       *SourceFilesystem
	I18n       *SourceFilesystem
	Layouts    *SourceFilesystem
	Archetypes *SourceFilesystem

	// When in multihost we have one static filesystem per language. The sync
	// static files is currently done outside of the Hugo build (where there is
	// a concept of a site per language).
	// When in non-multihost mode there will be one entry in this map with a blank key.
	Static map[string]*SourceFilesystem
}

// A SourceFilesystem holds the filesystem for a given source type in Hugo (data,
// i18n, layouts, static) and additional metadata to be able to use that filesystem
// in server mode.
type SourceFilesystem struct {
	Fs afero.Fs

	Dirnames []string

	// When syncing a source folder to the target (e.g. /public), this may
	// be set to publish into a subfolder. This is used for static syncing
	// in multihost mode.
	PublishFolder string
}

// IsStatic returns true if the given filename is a member of one of the static
// filesystems.
func (s SourceFilesystems) IsStatic(filename string) bool {
	for _, staticFs := range s.Static {
		if staticFs.Contains(filename) {
			return true
		}
	}
	return false
}

// IsLayout returns true if the given filename is a member of the layouts filesystem.
func (s SourceFilesystems) IsLayout(filename string) bool {
	return s.Layouts.Contains(filename)
}

// IsData returns true if the given filename is a member of the data filesystem.
func (s SourceFilesystems) IsData(filename string) bool {
	return s.Data.Contains(filename)
}

// IsI18n returns true if the given filename is a member of the i18n filesystem.
func (s SourceFilesystems) IsI18n(filename string) bool {
	return s.I18n.Contains(filename)
}

// MakeStaticPathRelative makes an absolute static filename into a relative one.
// It will return an empty string if the filename is not a member of a static filesystem.
func (s SourceFilesystems) MakeStaticPathRelative(filename string) string {
	for _, staticFs := range s.Static {
		rel := staticFs.MakePathRelative(filename)
		if rel != "" {
			return rel
		}
	}
	return ""
}

// MakePathRelative creates a relative path from the given filename.
// It will return an empty string if the filename is not a member of this filesystem.
func (d *SourceFilesystem) MakePathRelative(filename string) string {
	for _, currentPath := range d.Dirnames {
		if strings.HasPrefix(filename, currentPath) {
			return strings.TrimPrefix(filename, currentPath)
		}
	}
	return ""
}

// Contains returns whether the given filename is a member of the current filesystem.
func (d *SourceFilesystem) Contains(filename string) bool {
	for _, dir := range d.Dirnames {
		if strings.HasPrefix(filename, dir) {
			return true
		}
	}
	return false
}

// WithBaseFs allows reuse of some potentially expensive to create parts that remain
// the same across sites/languages.
func WithBaseFs(b *BaseFs) func(*BaseFs) error {
	return func(bb *BaseFs) error {
		bb.themeFs = b.themeFs
		bb.AbsThemeDirs = b.AbsThemeDirs
		return nil
	}
}

// NewBase builds the filesystems used by Hugo given the paths and options provided.NewBase
func NewBase(p *paths.Paths, options ...func(*BaseFs) error) (*BaseFs, error) {
	fs := p.Fs

	resourcesFs := afero.NewBasePathFs(fs.Source, p.AbsResourcesDir)
	publishFs := afero.NewBasePathFs(fs.Destination, p.AbsPublishDir)

	contentFs, absContentDirs, err := createContentFs(fs.Source, p.WorkingDir, p.DefaultContentLanguage, p.Languages)
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

	b := &BaseFs{
		AbsContentDirs: absContentDirs,
		ContentFs:      contentFs,
		ResourcesFs:    resourcesFs,
		PublishFs:      publishFs,
	}

	for _, opt := range options {
		if err := opt(b); err != nil {
			return nil, err
		}
	}

	builder := newSourceFilesystemsBuilder(p, b)
	sourceFilesystems, err := builder.Build()
	if err != nil {
		return nil, err
	}

	b.SourceFilesystems = sourceFilesystems
	b.themeFs = builder.themeFs
	b.AbsThemeDirs = builder.absThemeDirs

	return b, nil
}

type sourceFilesystemsBuilder struct {
	p            *paths.Paths
	result       *SourceFilesystems
	themeFs      afero.Fs
	hasTheme     bool
	absThemeDirs []string
}

func newSourceFilesystemsBuilder(p *paths.Paths, b *BaseFs) *sourceFilesystemsBuilder {
	return &sourceFilesystemsBuilder{p: p, themeFs: b.themeFs, absThemeDirs: b.AbsThemeDirs, result: &SourceFilesystems{}}
}

func (b *sourceFilesystemsBuilder) Build() (*SourceFilesystems, error) {
	if b.themeFs == nil && b.p.ThemeSet() {
		themeFs, absThemeDirs, err := createThemesOverlayFs(b.p)
		if err != nil {
			return nil, err
		}
		if themeFs == nil {
			panic("createThemesFs returned nil")
		}
		b.themeFs = themeFs
		b.absThemeDirs = absThemeDirs

	}

	b.hasTheme = len(b.absThemeDirs) > 0

	sfs, err := b.createRootMappingFs("dataDir", "data")
	if err != nil {
		return nil, err
	}
	b.result.Data = sfs

	sfs, err = b.createRootMappingFs("i18nDir", "i18n")
	if err != nil {
		return nil, err
	}
	b.result.I18n = sfs

	sfs, err = b.createFs("layoutDir", "layouts")
	if err != nil {
		return nil, err
	}
	b.result.Layouts = sfs

	sfs, err = b.createFs("archetypeDir", "archetypes")
	if err != nil {
		return nil, err
	}
	b.result.Archetypes = sfs

	err = b.createStaticFs()
	if err != nil {
		return nil, err
	}

	return b.result, nil
}

func (b *sourceFilesystemsBuilder) createFs(dirKey, themeFolder string) (*SourceFilesystem, error) {
	s := &SourceFilesystem{}
	dir := b.p.Cfg.GetString(dirKey)
	if dir == "" {
		return s, fmt.Errorf("config %q not set", dirKey)
	}

	var fs afero.Fs

	absDir := b.p.AbsPathify(dir)
	if b.existsInSource(absDir) {
		fs = afero.NewBasePathFs(b.p.Fs.Source, absDir)
		s.Dirnames = []string{absDir}
	}

	if b.hasTheme {
		themeFolderFs := afero.NewBasePathFs(b.themeFs, themeFolder)
		if fs == nil {
			fs = themeFolderFs
		} else {
			fs = afero.NewCopyOnWriteFs(themeFolderFs, fs)
		}

		for _, absThemeDir := range b.absThemeDirs {
			absThemeFolderDir := filepath.Join(absThemeDir, themeFolder)
			if b.existsInSource(absThemeFolderDir) {
				s.Dirnames = append(s.Dirnames, absThemeFolderDir)
			}
		}
	}

	if fs == nil {
		s.Fs = hugofs.NoOpFs
	} else {
		s.Fs = afero.NewReadOnlyFs(fs)
	}

	return s, nil
}

// Used for data, i18n -- we cannot use overlay filsesystems for those, but we need
// to keep a strict order.
func (b *sourceFilesystemsBuilder) createRootMappingFs(dirKey, themeFolder string) (*SourceFilesystem, error) {
	s := &SourceFilesystem{}

	projectDir := b.p.Cfg.GetString(dirKey)
	if projectDir == "" {
		return nil, fmt.Errorf("config %q not set", dirKey)
	}

	var fromTo []string
	to := b.p.AbsPathify(projectDir)

	if b.existsInSource(to) {
		s.Dirnames = []string{to}
		fromTo = []string{projectVirtualFolder, to}
	}

	for _, theme := range b.p.AllThemes {
		to := b.p.AbsPathify(filepath.Join(b.p.ThemesDir, theme.Name, themeFolder))
		if b.existsInSource(to) {
			s.Dirnames = append(s.Dirnames, to)
			from := theme
			fromTo = append(fromTo, from.Name, to)
		}
	}

	if len(fromTo) == 0 {
		s.Fs = hugofs.NoOpFs
		return s, nil
	}

	fs, err := hugofs.NewRootMappingFs(b.p.Fs.Source, fromTo...)
	if err != nil {
		return nil, err
	}

	s.Fs = afero.NewReadOnlyFs(fs)

	return s, nil

}

func (b *sourceFilesystemsBuilder) existsInSource(abspath string) bool {
	exists, _ := afero.Exists(b.p.Fs.Source, abspath)
	return exists
}

func (b *sourceFilesystemsBuilder) createStaticFs() error {
	isMultihost := b.p.Cfg.GetBool("multihost")
	ms := make(map[string]*SourceFilesystem)
	b.result.Static = ms

	if isMultihost {
		for _, l := range b.p.Languages {
			s := &SourceFilesystem{PublishFolder: l.Lang}
			staticDirs := removeDuplicatesKeepRight(getStaticDirs(l))
			if len(staticDirs) == 0 {
				continue
			}

			for _, dir := range staticDirs {
				absDir := b.p.AbsPathify(dir)
				if !b.existsInSource(absDir) {
					continue
				}

				s.Dirnames = append(s.Dirnames, absDir)
			}

			fs, err := createOverlayFs(b.p.Fs.Source, s.Dirnames)
			if err != nil {
				return err
			}

			s.Fs = fs
			ms[l.Lang] = s

		}

		return nil
	}

	s := &SourceFilesystem{}
	var staticDirs []string

	for _, l := range b.p.Languages {
		staticDirs = append(staticDirs, getStaticDirs(l)...)
	}

	staticDirs = removeDuplicatesKeepRight(staticDirs)
	if len(staticDirs) == 0 {
		return nil
	}

	for _, dir := range staticDirs {
		absDir := b.p.AbsPathify(dir)
		if !b.existsInSource(absDir) {
			continue
		}
		s.Dirnames = append(s.Dirnames, absDir)
	}

	fs, err := createOverlayFs(b.p.Fs.Source, s.Dirnames)
	if err != nil {
		return err
	}

	if b.hasTheme {
		themeFolder := "static"
		fs = afero.NewCopyOnWriteFs(afero.NewBasePathFs(b.themeFs, themeFolder), fs)
		for _, absThemeDir := range b.absThemeDirs {
			s.Dirnames = append(s.Dirnames, filepath.Join(absThemeDir, themeFolder))
		}
	}

	s.Fs = fs
	ms[""] = s

	return nil
}

func getStaticDirs(cfg config.Provider) []string {
	var staticDirs []string
	for i := -1; i <= 10; i++ {
		staticDirs = append(staticDirs, getStringOrStringSlice(cfg, "staticDir", i)...)
	}
	return staticDirs
}

func getStringOrStringSlice(cfg config.Provider, key string, id int) []string {

	if id >= 0 {
		key = fmt.Sprintf("%s%d", key, id)
	}

	return config.GetStringSlicePreserveString(cfg, key)

}

func createContentFs(fs afero.Fs,
	workingDir,
	defaultContentLanguage string,
	languages langs.Languages) (afero.Fs, []types.KeyValueStr, error) {

	var contentLanguages langs.Languages
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
	languages langs.Languages,
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

	absContentDir := paths.AbsPathify(workingDir, language.ContentDir)
	if !strings.HasSuffix(absContentDir, paths.FilePathSeparator) {
		absContentDir += paths.FilePathSeparator
	}

	// If root, remove the second '/'
	if absContentDir == "//" {
		absContentDir = paths.FilePathSeparator
	}

	if len(absContentDir) < 6 {
		return nil, fmt.Errorf("invalid content dir %q: Path is too short", absContentDir)
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

func createThemesOverlayFs(p *paths.Paths) (afero.Fs, []string, error) {

	themes := p.AllThemes

	if len(themes) == 0 {
		panic("AllThemes not set")
	}

	themesDir := p.AbsPathify(p.ThemesDir)
	if themesDir == "" {
		return nil, nil, errors.New("no themes dir set")
	}

	absPaths := make([]string, len(themes))

	// The themes are ordered from left to right. We need to revert it to get the
	// overlay logic below working as expected.
	for i := 0; i < len(themes); i++ {
		absPaths[i] = filepath.Join(themesDir, themes[len(themes)-1-i].Name)
	}

	fs, err := createOverlayFs(p.Fs.Source, absPaths)

	return fs, absPaths, err

}

func createOverlayFs(source afero.Fs, absPaths []string) (afero.Fs, error) {
	if len(absPaths) == 0 {
		return hugofs.NoOpFs, nil
	}

	if len(absPaths) == 1 {
		return afero.NewReadOnlyFs(afero.NewBasePathFs(source, absPaths[0])), nil
	}

	base := afero.NewReadOnlyFs(afero.NewBasePathFs(source, absPaths[0]))
	overlay, err := createOverlayFs(source, absPaths[1:])
	if err != nil {
		return nil, err
	}

	return afero.NewCopyOnWriteFs(base, overlay), nil
}

func removeDuplicatesKeepRight(in []string) []string {
	seen := make(map[string]bool)
	var out []string
	for i := len(in) - 1; i >= 0; i-- {
		v := in[i]
		if seen[v] {
			continue
		}
		out = append([]string{v}, out...)
		seen[v] = true
	}

	return out
}

func printFs(fs afero.Fs, path string, w io.Writer) {
	if fs == nil {
		return
	}
	afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			s := path
			if lang, ok := info.(hugofs.LanguageAnnouncer); ok {
				s = s + "\tLANG: " + lang.Lang()
			}
			if fp, ok := info.(hugofs.FilePather); ok {
				s = s + "\tRF: " + fp.Filename() + "\tBP: " + fp.BaseDir()
			}
			fmt.Fprintln(w, "    ", s)
		}
		return nil
	})
}
