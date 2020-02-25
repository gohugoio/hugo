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
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/modules"

	"github.com/gohugoio/hugo/hugofs"

	"fmt"

	"github.com/gohugoio/hugo/hugolib/paths"
	"github.com/spf13/afero"
)

var filePathSeparator = string(filepath.Separator)

// BaseFs contains the core base filesystems used by Hugo. The name "base" is used
// to underline that even if they can be composites, they all have a base path set to a specific
// resource folder, e.g "/my-project/content". So, no absolute filenames needed.
type BaseFs struct {

	// SourceFilesystems contains the different source file systems.
	*SourceFilesystems

	// The filesystem used to publish the rendered site.
	// This usually maps to /my-project/public.
	PublishFs afero.Fs

	theBigFs *filesystemsCollector
}

// TODO(bep) we can get regular files in here and that is fine, but
// we need to clean up the naming.
func (fs *BaseFs) WatchDirs() []hugofs.FileMetaInfo {
	var dirs []hugofs.FileMetaInfo
	for _, dir := range fs.AllDirs() {
		if dir.Meta().Watch() {
			dirs = append(dirs, dir)
		}
	}
	return dirs
}

func (fs *BaseFs) AllDirs() []hugofs.FileMetaInfo {
	var dirs []hugofs.FileMetaInfo
	for _, dirSet := range [][]hugofs.FileMetaInfo{
		fs.Archetypes.Dirs,
		fs.I18n.Dirs,
		fs.Data.Dirs,
		fs.Content.Dirs,
		fs.Assets.Dirs,
		fs.Layouts.Dirs,
		//fs.Resources.Dirs,
		fs.StaticDirs,
	} {
		dirs = append(dirs, dirSet...)
	}

	return dirs
}

// RelContentDir tries to create a path relative to the content root from
// the given filename. The return value is the path and language code.
func (b *BaseFs) RelContentDir(filename string) string {
	for _, dir := range b.SourceFilesystems.Content.Dirs {
		dirname := dir.Meta().Filename()
		if strings.HasPrefix(filename, dirname) {
			rel := path.Join(dir.Meta().Path(), strings.TrimPrefix(filename, dirname))
			return strings.TrimPrefix(rel, filePathSeparator)
		}
	}
	// Either not a content dir or already relative.
	return filename
}

// SourceFilesystems contains the different source file systems. These can be
// composite file systems (theme and project etc.), and they have all root
// set to the source type the provides: data, i18n, static, layouts.
type SourceFilesystems struct {
	Content    *SourceFilesystem
	Data       *SourceFilesystem
	I18n       *SourceFilesystem
	Layouts    *SourceFilesystem
	Archetypes *SourceFilesystem
	Assets     *SourceFilesystem

	// Writable filesystem on top the project's resources directory,
	// with any sub module's resource fs layered below.
	ResourcesCache afero.Fs

	// The project folder.
	Work afero.Fs

	// When in multihost we have one static filesystem per language. The sync
	// static files is currently done outside of the Hugo build (where there is
	// a concept of a site per language).
	// When in non-multihost mode there will be one entry in this map with a blank key.
	Static map[string]*SourceFilesystem

	// All the /static dirs (including themes/modules).
	StaticDirs []hugofs.FileMetaInfo
}

// FileSystems returns the FileSystems relevant for the change detection
// in server mode.
// Note: This does currently not return any static fs.
func (s *SourceFilesystems) FileSystems() []*SourceFilesystem {
	return []*SourceFilesystem{
		s.Content,
		s.Data,
		s.I18n,
		s.Layouts,
		s.Archetypes,
		// TODO(bep) static
	}

}

// A SourceFilesystem holds the filesystem for a given source type in Hugo (data,
// i18n, layouts, static) and additional metadata to be able to use that filesystem
// in server mode.
type SourceFilesystem struct {
	// Name matches one in files.ComponentFolders
	Name string

	// This is a virtual composite filesystem. It expects path relative to a context.
	Fs afero.Fs

	// This filesystem as separate root directories, starting from project and down
	// to the themes/modules.
	Dirs []hugofs.FileMetaInfo

	// When syncing a source folder to the target (e.g. /public), this may
	// be set to publish into a subfolder. This is used for static syncing
	// in multihost mode.
	PublishFolder string
}

// ContentStaticAssetFs will create a new composite filesystem from the content,
// static, and asset filesystems. The site language is needed to pick the correct static filesystem.
// The order is content, static and then assets.
// TODO(bep) check usage
func (s SourceFilesystems) ContentStaticAssetFs(lang string) afero.Fs {
	staticFs := s.StaticFs(lang)

	base := afero.NewCopyOnWriteFs(s.Assets.Fs, staticFs)
	return afero.NewCopyOnWriteFs(base, s.Content.Fs)

}

// StaticFs returns the static filesystem for the given language.
// This can be a composite filesystem.
func (s SourceFilesystems) StaticFs(lang string) afero.Fs {
	var staticFs afero.Fs = hugofs.NoOpFs

	if fs, ok := s.Static[lang]; ok {
		staticFs = fs.Fs
	} else if fs, ok := s.Static[""]; ok {
		staticFs = fs.Fs
	}

	return staticFs
}

// StatResource looks for a resource in these filesystems in order: static, assets and finally content.
// If found in any of them, it returns FileInfo and the relevant filesystem.
// Any non os.IsNotExist error will be returned.
// An os.IsNotExist error wil be returned only if all filesystems return such an error.
// Note that if we only wanted to find the file, we could create a composite Afero fs,
// but we also need to know which filesystem root it lives in.
func (s SourceFilesystems) StatResource(lang, filename string) (fi os.FileInfo, fs afero.Fs, err error) {
	for _, fsToCheck := range []afero.Fs{s.StaticFs(lang), s.Assets.Fs, s.Content.Fs} {
		fs = fsToCheck
		fi, err = fs.Stat(filename)
		if err == nil || !os.IsNotExist(err) {
			return
		}
	}
	// Not found.
	return
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

// IsContent returns true if the given filename is a member of the content filesystem.
func (s SourceFilesystems) IsContent(filename string) bool {
	return s.Content.Contains(filename)
}

// IsLayout returns true if the given filename is a member of the layouts filesystem.
func (s SourceFilesystems) IsLayout(filename string) bool {
	return s.Layouts.Contains(filename)
}

// IsData returns true if the given filename is a member of the data filesystem.
func (s SourceFilesystems) IsData(filename string) bool {
	return s.Data.Contains(filename)
}

// IsAsset returns true if the given filename is a member of the asset filesystem.
func (s SourceFilesystems) IsAsset(filename string) bool {
	return s.Assets.Contains(filename)
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

	for _, dir := range d.Dirs {
		meta := dir.(hugofs.FileMetaInfo).Meta()
		currentPath := meta.Filename()

		if strings.HasPrefix(filename, currentPath) {
			rel := strings.TrimPrefix(filename, currentPath)
			if mp := meta.Path(); mp != "" {
				rel = filepath.Join(mp, rel)
			}
			return strings.TrimPrefix(rel, filePathSeparator)
		}
	}
	return ""
}

func (d *SourceFilesystem) RealFilename(rel string) string {
	fi, err := d.Fs.Stat(rel)
	if err != nil {
		return rel
	}
	if realfi, ok := fi.(hugofs.FileMetaInfo); ok {
		return realfi.Meta().Filename()
	}

	return rel
}

// Contains returns whether the given filename is a member of the current filesystem.
func (d *SourceFilesystem) Contains(filename string) bool {
	for _, dir := range d.Dirs {
		if strings.HasPrefix(filename, dir.Meta().Filename()) {
			return true
		}
	}
	return false
}

// Path returns the mount relative path to the given filename if it is a member of
// of the current filesystem, an empty string if not.
func (d *SourceFilesystem) Path(filename string) string {
	for _, dir := range d.Dirs {
		meta := dir.Meta()
		if strings.HasPrefix(filename, meta.Filename()) {
			p := strings.TrimPrefix(strings.TrimPrefix(filename, meta.Filename()), filePathSeparator)
			if mountRoot := meta.MountRoot(); mountRoot != "" {
				return filepath.Join(mountRoot, p)
			}
			return p
		}
	}
	return ""
}

// RealDirs gets a list of absolute paths to directories starting from the given
// path.
func (d *SourceFilesystem) RealDirs(from string) []string {
	var dirnames []string
	for _, dir := range d.Dirs {
		meta := dir.Meta()
		dirname := filepath.Join(meta.Filename(), from)
		_, err := meta.Fs().Stat(from)

		if err == nil {
			dirnames = append(dirnames, dirname)
		}
	}
	return dirnames
}

// WithBaseFs allows reuse of some potentially expensive to create parts that remain
// the same across sites/languages.
func WithBaseFs(b *BaseFs) func(*BaseFs) error {
	return func(bb *BaseFs) error {
		bb.theBigFs = b.theBigFs
		bb.SourceFilesystems = b.SourceFilesystems
		return nil
	}
}

// NewBase builds the filesystems used by Hugo given the paths and options provided.NewBase
func NewBase(p *paths.Paths, logger *loggers.Logger, options ...func(*BaseFs) error) (*BaseFs, error) {
	fs := p.Fs
	if logger == nil {
		logger = loggers.NewWarningLogger()
	}

	publishFs := hugofs.NewBaseFileDecorator(afero.NewBasePathFs(fs.Destination, p.AbsPublishDir))

	b := &BaseFs{
		PublishFs: publishFs,
	}

	for _, opt := range options {
		if err := opt(b); err != nil {
			return nil, err
		}
	}

	if b.theBigFs != nil && b.SourceFilesystems != nil {
		return b, nil
	}

	builder := newSourceFilesystemsBuilder(p, logger, b)
	sourceFilesystems, err := builder.Build()
	if err != nil {
		return nil, errors.Wrap(err, "build filesystems")
	}

	b.SourceFilesystems = sourceFilesystems
	b.theBigFs = builder.theBigFs

	return b, nil
}

type sourceFilesystemsBuilder struct {
	logger   *loggers.Logger
	p        *paths.Paths
	sourceFs afero.Fs
	result   *SourceFilesystems
	theBigFs *filesystemsCollector
}

func newSourceFilesystemsBuilder(p *paths.Paths, logger *loggers.Logger, b *BaseFs) *sourceFilesystemsBuilder {
	sourceFs := hugofs.NewBaseFileDecorator(p.Fs.Source)
	return &sourceFilesystemsBuilder{p: p, logger: logger, sourceFs: sourceFs, theBigFs: b.theBigFs, result: &SourceFilesystems{}}
}

func (b *sourceFilesystemsBuilder) newSourceFilesystem(name string, fs afero.Fs, dirs []hugofs.FileMetaInfo) *SourceFilesystem {
	return &SourceFilesystem{
		Name: name,
		Fs:   fs,
		Dirs: dirs,
	}
}

func (b *sourceFilesystemsBuilder) Build() (*SourceFilesystems, error) {

	if b.theBigFs == nil {

		theBigFs, err := b.createMainOverlayFs(b.p)
		if err != nil {
			return nil, errors.Wrap(err, "create main fs")
		}

		b.theBigFs = theBigFs
	}

	createView := func(componentID string) *SourceFilesystem {
		if b.theBigFs == nil || b.theBigFs.overlayMounts == nil {
			return b.newSourceFilesystem(componentID, hugofs.NoOpFs, nil)
		}

		dirs := b.theBigFs.overlayDirs[componentID]

		return b.newSourceFilesystem(componentID, afero.NewBasePathFs(b.theBigFs.overlayMounts, componentID), dirs)

	}

	b.theBigFs.finalizeDirs()

	b.result.Archetypes = createView(files.ComponentFolderArchetypes)
	b.result.Layouts = createView(files.ComponentFolderLayouts)
	b.result.Assets = createView(files.ComponentFolderAssets)
	b.result.ResourcesCache = b.theBigFs.overlayResources

	// Data, i18n and content cannot use the overlay fs
	dataDirs := b.theBigFs.overlayDirs[files.ComponentFolderData]
	dataFs, err := hugofs.NewSliceFs(dataDirs...)
	if err != nil {
		return nil, err
	}

	b.result.Data = b.newSourceFilesystem(files.ComponentFolderData, dataFs, dataDirs)

	i18nDirs := b.theBigFs.overlayDirs[files.ComponentFolderI18n]
	i18nFs, err := hugofs.NewSliceFs(i18nDirs...)
	if err != nil {
		return nil, err
	}
	b.result.I18n = b.newSourceFilesystem(files.ComponentFolderI18n, i18nFs, i18nDirs)

	contentDirs := b.theBigFs.overlayDirs[files.ComponentFolderContent]
	contentBfs := afero.NewBasePathFs(b.theBigFs.overlayMountsContent, files.ComponentFolderContent)

	contentFs, err := hugofs.NewLanguageFs(b.p.LanguagesDefaultFirst.AsOrdinalSet(), contentBfs)
	if err != nil {
		return nil, errors.Wrap(err, "create content filesystem")
	}

	b.result.Content = b.newSourceFilesystem(files.ComponentFolderContent, contentFs, contentDirs)

	b.result.Work = afero.NewReadOnlyFs(b.theBigFs.overlayFull)

	// Create static filesystem(s)
	ms := make(map[string]*SourceFilesystem)
	b.result.Static = ms
	b.result.StaticDirs = b.theBigFs.overlayDirs[files.ComponentFolderStatic]

	if b.theBigFs.staticPerLanguage != nil {
		// Multihost mode
		for k, v := range b.theBigFs.staticPerLanguage {
			sfs := b.newSourceFilesystem(files.ComponentFolderStatic, v, b.result.StaticDirs)
			sfs.PublishFolder = k
			ms[k] = sfs
		}
	} else {
		bfs := afero.NewBasePathFs(b.theBigFs.overlayMountsStatic, files.ComponentFolderStatic)
		ms[""] = b.newSourceFilesystem(files.ComponentFolderStatic, bfs, b.result.StaticDirs)
	}

	return b.result, nil

}

func (b *sourceFilesystemsBuilder) createMainOverlayFs(p *paths.Paths) (*filesystemsCollector, error) {

	var staticFsMap map[string]afero.Fs
	if b.p.Cfg.GetBool("multihost") {
		staticFsMap = make(map[string]afero.Fs)
	}

	collector := &filesystemsCollector{
		sourceProject:     b.sourceFs,
		sourceModules:     hugofs.NewNoSymlinkFs(b.sourceFs, b.logger, false),
		overlayDirs:       make(map[string][]hugofs.FileMetaInfo),
		staticPerLanguage: staticFsMap,
	}

	mods := p.AllModules

	if len(mods) == 0 {
		return collector, nil
	}

	modsReversed := make([]mountsDescriptor, len(mods))

	// The theme components are ordered from left to right.
	// We need to revert it to get the
	// overlay logic below working as expected, with the project on top.
	j := 0
	for i := len(mods) - 1; i >= 0; i-- {
		mod := mods[i]
		dir := mod.Dir()

		isMainProject := mod.Owner() == nil
		modsReversed[j] = mountsDescriptor{
			Module:        mod,
			dir:           dir,
			isMainProject: isMainProject,
		}
		j++
	}

	err := b.createOverlayFs(collector, modsReversed)

	return collector, err

}

func (b *sourceFilesystemsBuilder) isContentMount(mnt modules.Mount) bool {
	return strings.HasPrefix(mnt.Target, files.ComponentFolderContent)
}

func (b *sourceFilesystemsBuilder) isStaticMount(mnt modules.Mount) bool {
	return strings.HasPrefix(mnt.Target, files.ComponentFolderStatic)
}

func (b *sourceFilesystemsBuilder) createModFs(
	collector *filesystemsCollector,
	md mountsDescriptor) error {

	var (
		fromTo        []hugofs.RootMapping
		fromToContent []hugofs.RootMapping
		fromToStatic  []hugofs.RootMapping
	)

	absPathify := func(path string) (string, string) {
		if filepath.IsAbs(path) {
			return "", path
		}
		return md.dir, paths.AbsPathify(md.dir, path)
	}

	for _, mount := range md.Mounts() {

		mountWeight := 1
		if md.isMainProject {
			mountWeight++
		}

		base, filename := absPathify(mount.Source)

		rm := hugofs.RootMapping{
			From:      mount.Target,
			To:        filename,
			ToBasedir: base,
			Module:    md.Module.Path(),
			Meta: hugofs.FileMeta{
				"watch":       md.Watch(),
				"mountWeight": mountWeight,
			},
		}

		isContentMount := b.isContentMount(mount)

		lang := mount.Lang
		if lang == "" && isContentMount {
			lang = b.p.DefaultContentLanguage
		}

		rm.Meta["lang"] = lang

		if isContentMount {
			fromToContent = append(fromToContent, rm)
		} else if b.isStaticMount(mount) {
			fromToStatic = append(fromToStatic, rm)
		} else {
			fromTo = append(fromTo, rm)
		}
	}

	modBase := collector.sourceProject
	if !md.isMainProject {
		modBase = collector.sourceModules
	}
	sourceStatic := hugofs.NewNoSymlinkFs(modBase, b.logger, true)

	rmfs, err := hugofs.NewRootMappingFs(modBase, fromTo...)
	if err != nil {
		return err
	}
	rmfsContent, err := hugofs.NewRootMappingFs(modBase, fromToContent...)
	if err != nil {
		return err
	}
	rmfsStatic, err := hugofs.NewRootMappingFs(sourceStatic, fromToStatic...)
	if err != nil {
		return err
	}

	// We need to keep the ordered list of directories for watching and
	// some special merge operations (data, i18n).
	collector.addDirs(rmfs)
	collector.addDirs(rmfsContent)
	collector.addDirs(rmfsStatic)

	if collector.staticPerLanguage != nil {
		for _, l := range b.p.Languages {
			lang := l.Lang

			lfs := rmfsStatic.Filter(func(rm hugofs.RootMapping) bool {
				rlang := rm.Meta.Lang()
				return rlang == "" || rlang == lang
			})

			bfs := afero.NewBasePathFs(lfs, files.ComponentFolderStatic)

			sfs, found := collector.staticPerLanguage[lang]
			if found {
				collector.staticPerLanguage[lang] = afero.NewCopyOnWriteFs(sfs, bfs)

			} else {
				collector.staticPerLanguage[lang] = bfs
			}
		}
	}

	getResourcesDir := func() string {
		if md.isMainProject {
			return b.p.AbsResourcesDir
		}
		_, filename := absPathify(files.FolderResources)
		return filename
	}

	if collector.overlayMounts == nil {
		collector.overlayMounts = rmfs
		collector.overlayMountsContent = rmfsContent
		collector.overlayMountsStatic = rmfsStatic
		collector.overlayFull = afero.NewBasePathFs(modBase, md.dir)
		collector.overlayResources = afero.NewBasePathFs(modBase, getResourcesDir())
	} else {

		collector.overlayMounts = afero.NewCopyOnWriteFs(collector.overlayMounts, rmfs)
		collector.overlayMountsContent = hugofs.NewLanguageCompositeFs(collector.overlayMountsContent, rmfsContent)
		collector.overlayMountsStatic = hugofs.NewLanguageCompositeFs(collector.overlayMountsStatic, rmfsStatic)
		collector.overlayFull = afero.NewCopyOnWriteFs(collector.overlayFull, afero.NewBasePathFs(modBase, md.dir))
		collector.overlayResources = afero.NewCopyOnWriteFs(collector.overlayResources, afero.NewBasePathFs(modBase, getResourcesDir()))
	}

	return nil

}

func printFs(fs afero.Fs, path string, w io.Writer) {
	if fs == nil {
		return
	}
	afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		var filename string
		if fim, ok := info.(hugofs.FileMetaInfo); ok {
			filename = fim.Meta().Filename()
		}
		fmt.Fprintf(w, "    %q %q\n", path, filename)
		return nil
	})
}

type filesystemsCollector struct {
	sourceProject afero.Fs // Source for project folders
	sourceModules afero.Fs // Source for modules/themes

	overlayMounts        afero.Fs
	overlayMountsContent afero.Fs
	overlayMountsStatic  afero.Fs
	overlayFull          afero.Fs
	overlayResources     afero.Fs

	// Maps component type (layouts, static, content etc.) an ordered list of
	// directories representing the overlay filesystems above.
	overlayDirs map[string][]hugofs.FileMetaInfo

	// Set if in multihost mode
	staticPerLanguage map[string]afero.Fs

	finalizerInit sync.Once
}

func (c *filesystemsCollector) addDirs(rfs *hugofs.RootMappingFs) {
	for _, componentFolder := range files.ComponentFolders {
		dirs, err := rfs.Dirs(componentFolder)

		if err == nil {
			c.overlayDirs[componentFolder] = append(c.overlayDirs[componentFolder], dirs...)
		}
	}
}

func (c *filesystemsCollector) finalizeDirs() {
	c.finalizerInit.Do(func() {
		// Order the directories from top to bottom (project, theme a, theme ...).
		for _, dirs := range c.overlayDirs {
			c.reverseFis(dirs)
		}
	})

}

func (c *filesystemsCollector) reverseFis(fis []hugofs.FileMetaInfo) {
	for i := len(fis)/2 - 1; i >= 0; i-- {
		opp := len(fis) - 1 - i
		fis[i], fis[opp] = fis[opp], fis[i]
	}
}

type mountsDescriptor struct {
	modules.Module
	dir           string
	isMainProject bool
}

func (b *sourceFilesystemsBuilder) createOverlayFs(collector *filesystemsCollector, mounts []mountsDescriptor) error {
	if len(mounts) == 0 {
		return nil
	}

	err := b.createModFs(collector, mounts[0])
	if err != nil {
		return err
	}

	if len(mounts) == 1 {
		return nil
	}

	return b.createOverlayFs(collector, mounts[1:])
}
