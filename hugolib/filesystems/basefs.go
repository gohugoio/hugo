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

// Package filesystems provides the fine grained file systems used by Hugo. These
// are typically virtual filesystems that are composites of project and theme content.
package filesystems

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bep/overlayfs"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/types"

	"github.com/rogpeppe/go-internal/lockedfile"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/modules"

	hpaths "github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/paths"
	"github.com/spf13/afero"
)

const (
	// Used to control concurrency between multiple Hugo instances, e.g.
	// a running server and building new content with 'hugo new'.
	// It's placed in the project root.
	lockFileBuild = ".hugo_build.lock"
)

var filePathSeparator = string(filepath.Separator)

// BaseFs contains the core base filesystems used by Hugo. The name "base" is used
// to underline that even if they can be composites, they all have a base path set to a specific
// resource folder, e.g "/my-project/content". So, no absolute filenames needed.
type BaseFs struct {
	// SourceFilesystems contains the different source file systems.
	*SourceFilesystems

	// The source filesystem (needs absolute filenames).
	SourceFs afero.Fs

	// The project source.
	ProjectSourceFs afero.Fs

	// The filesystem used to publish the rendered site.
	// This usually maps to /my-project/public.
	PublishFs afero.Fs

	// The filesystem used for static files.
	PublishFsStatic afero.Fs

	// A read-only filesystem starting from the project workDir.
	WorkDir afero.Fs

	theBigFs *filesystemsCollector

	workingDir string

	// Locks.
	buildMu Lockable // <project>/.hugo_build.lock
}

type Lockable interface {
	Lock() (unlock func(), err error)
}

type fakeLockfileMutex struct {
	mu sync.Mutex
}

func (f *fakeLockfileMutex) Lock() (func(), error) {
	f.mu.Lock()
	return func() { f.mu.Unlock() }, nil
}

// Tries to acquire a build lock.
func (b *BaseFs) LockBuild() (unlock func(), err error) {
	return b.buildMu.Lock()
}

func (b *BaseFs) WatchFilenames() []string {
	var filenames []string
	sourceFs := b.SourceFs

	for _, rfs := range b.RootFss {
		for _, component := range files.ComponentFolders {
			fis, err := rfs.Mounts(component)
			if err != nil {
				continue
			}

			for _, fim := range fis {
				meta := fim.Meta()
				if !meta.Watch {
					continue
				}

				if !fim.IsDir() {
					filenames = append(filenames, meta.Filename)
					continue
				}

				w := hugofs.NewWalkway(hugofs.WalkwayConfig{
					Fs:   sourceFs,
					Root: meta.Filename,
					WalkFn: func(path string, fi hugofs.FileMetaInfo) error {
						if !fi.IsDir() {
							return nil
						}
						if fi.Name() == ".git" ||
							fi.Name() == "node_modules" || fi.Name() == "bower_components" {
							return filepath.SkipDir
						}
						filenames = append(filenames, fi.Meta().Filename)
						return nil
					},
				})

				w.Walk()
			}

		}
	}

	return filenames
}

func (b *BaseFs) mountsForComponent(component string) []hugofs.FileMetaInfo {
	var result []hugofs.FileMetaInfo
	for _, rfs := range b.RootFss {
		dirs, err := rfs.Mounts(component)
		if err == nil {
			result = append(result, dirs...)
		}
	}
	return result
}

// AbsProjectContentDir tries to construct a filename below the most
// relevant content directory.
func (b *BaseFs) AbsProjectContentDir(filename string) (string, string, error) {
	isAbs := filepath.IsAbs(filename)
	for _, fi := range b.mountsForComponent(files.ComponentFolderContent) {
		if !fi.IsDir() {
			continue
		}
		meta := fi.Meta()
		if !meta.IsProject {
			continue
		}

		if isAbs {
			if strings.HasPrefix(filename, meta.Filename) {
				return strings.TrimPrefix(filename, meta.Filename+filePathSeparator), filename, nil
			}
		} else {
			contentDir := strings.TrimPrefix(strings.TrimPrefix(meta.Filename, meta.BaseDir), filePathSeparator) + filePathSeparator

			if strings.HasPrefix(filename, contentDir) {
				relFilename := strings.TrimPrefix(filename, contentDir)
				absFilename := filepath.Join(meta.Filename, relFilename)
				return relFilename, absFilename, nil
			}
		}

	}

	if !isAbs {
		// A filename on the form "posts/mypage.md", put it inside
		// the first content folder, usually <workDir>/content.
		// Pick the first project dir (which is probably the most important one).
		for _, dir := range b.SourceFilesystems.Content.mounts() {
			if !dir.IsDir() {
				continue
			}
			meta := dir.Meta()
			if meta.IsProject {
				return filename, filepath.Join(meta.Filename, filename), nil
			}
		}
	}

	return "", "", fmt.Errorf("could not determine content directory for %q", filename)
}

// ResolveJSConfigFile resolves the JS-related config file to a absolute
// filename. One example of such would be postcss.config.js.
func (b *BaseFs) ResolveJSConfigFile(name string) string {
	// First look in assets/_jsconfig
	fi, err := b.Assets.Fs.Stat(filepath.Join(files.FolderJSConfig, name))
	if err == nil {
		return fi.(hugofs.FileMetaInfo).Meta().Filename
	}
	// Fall back to the work dir.
	fi, err = b.Work.Stat(name)
	if err == nil {
		return fi.(hugofs.FileMetaInfo).Meta().Filename
	}

	return ""
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

	AssetsWithDuplicatesPreserved *SourceFilesystem

	RootFss []*hugofs.RootMappingFs

	// Writable filesystem on top the project's resources directory,
	// with any sub module's resource fs layered below.
	ResourcesCache afero.Fs

	// The work folder (may be a composite of project and theme components).
	Work afero.Fs

	// When in multihost we have one static filesystem per language. The sync
	// static files is currently done outside of the Hugo build (where there is
	// a concept of a site per language).
	// When in non-multihost mode there will be one entry in this map with a blank key.
	Static map[string]*SourceFilesystem

	conf config.AllProvider
}

// A SourceFilesystem holds the filesystem for a given source type in Hugo (data,
// i18n, layouts, static) and additional metadata to be able to use that filesystem
// in server mode.
type SourceFilesystem struct {
	// Name matches one in files.ComponentFolders
	Name string

	// This is a virtual composite filesystem. It expects path relative to a context.
	Fs afero.Fs

	// The source filesystem (usually the OS filesystem).
	SourceFs afero.Fs

	// When syncing a source folder to the target (e.g. /public), this may
	// be set to publish into a subfolder. This is used for static syncing
	// in multihost mode.
	PublishFolder string
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
// Any non herrors.IsNotExist error will be returned.
// An herrors.IsNotExist error will be returned only if all filesystems return such an error.
// Note that if we only wanted to find the file, we could create a composite Afero fs,
// but we also need to know which filesystem root it lives in.
func (s SourceFilesystems) StatResource(lang, filename string) (fi os.FileInfo, fs afero.Fs, err error) {
	for _, fsToCheck := range []afero.Fs{s.StaticFs(lang), s.Assets.Fs, s.Content.Fs} {
		fs = fsToCheck
		fi, err = fs.Stat(filename)
		if err == nil || !herrors.IsNotExist(err) {
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

// ResolvePaths resolves the given filename to a list of paths in the filesystems.
func (s *SourceFilesystems) ResolvePaths(filename string) []hugofs.ComponentPath {
	var cpss []hugofs.ComponentPath
	for _, rfs := range s.RootFss {
		cps, err := rfs.ReverseLookup(filename)
		if err != nil {
			panic(err)
		}
		cpss = append(cpss, cps...)
	}
	return cpss
}

// MakeStaticPathRelative makes an absolute static filename into a relative one.
// It will return an empty string if the filename is not a member of a static filesystem.
func (s SourceFilesystems) MakeStaticPathRelative(filename string) string {
	for _, staticFs := range s.Static {
		rel, _ := staticFs.MakePathRelative(filename, true)
		if rel != "" {
			return rel
		}
	}
	return ""
}

// MakePathRelative creates a relative path from the given filename.
func (d *SourceFilesystem) MakePathRelative(filename string, checkExists bool) (string, bool) {
	cps, err := d.ReverseLookup(filename, checkExists)
	if err != nil {
		panic(err)
	}
	if len(cps) == 0 {
		return "", false
	}

	return filepath.FromSlash(cps[0].Path), true
}

// ReverseLookup returns the component paths for the given filename.
func (d *SourceFilesystem) ReverseLookup(filename string, checkExists bool) ([]hugofs.ComponentPath, error) {
	var cps []hugofs.ComponentPath
	hugofs.WalkFilesystems(d.Fs, func(fs afero.Fs) bool {
		if rfs, ok := fs.(hugofs.ReverseLookupProvder); ok {
			if c, err := rfs.ReverseLookupComponent(d.Name, filename); err == nil {
				if checkExists {
					n := 0
					for _, cp := range c {
						if _, err := d.Fs.Stat(filepath.FromSlash(cp.Path)); err == nil {
							c[n] = cp
							n++
						}
					}
					c = c[:n]
				}
				cps = append(cps, c...)
			}
		}
		return false
	})
	return cps, nil
}

func (d *SourceFilesystem) mounts() []hugofs.FileMetaInfo {
	var m []hugofs.FileMetaInfo
	hugofs.WalkFilesystems(d.Fs, func(fs afero.Fs) bool {
		if rfs, ok := fs.(*hugofs.RootMappingFs); ok {
			mounts, err := rfs.Mounts(d.Name)
			if err == nil {
				m = append(m, mounts...)
			}
		}
		return false
	})

	// Filter out any mounts not belonging to this filesystem.
	// TODO(bep) I think this is superflous.
	n := 0
	for _, mm := range m {
		if mm.Meta().Component == d.Name {
			m[n] = mm
			n++
		}
	}
	m = m[:n]

	return m
}

func (d *SourceFilesystem) RealFilename(rel string) string {
	fi, err := d.Fs.Stat(rel)
	if err != nil {
		return rel
	}
	if realfi, ok := fi.(hugofs.FileMetaInfo); ok {
		return realfi.Meta().Filename
	}

	return rel
}

// Contains returns whether the given filename is a member of the current filesystem.
func (d *SourceFilesystem) Contains(filename string) bool {
	for _, dir := range d.mounts() {
		if !dir.IsDir() {
			continue
		}
		if strings.HasPrefix(filename, dir.Meta().Filename) {
			return true
		}
	}
	return false
}

// RealDirs gets a list of absolute paths to directories starting from the given
// path.
func (d *SourceFilesystem) RealDirs(from string) []string {
	var dirnames []string
	for _, m := range d.mounts() {
		if !m.IsDir() {
			continue
		}
		dirname := filepath.Join(m.Meta().Filename, from)
		if _, err := d.SourceFs.Stat(dirname); err == nil {
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
func NewBase(p *paths.Paths, logger loggers.Logger, options ...func(*BaseFs) error) (*BaseFs, error) {
	fs := p.Fs
	if logger == nil {
		logger = loggers.NewDefault()
	}

	publishFs := hugofs.NewBaseFileDecorator(fs.PublishDir)
	projectSourceFs := hugofs.NewBaseFileDecorator(hugofs.NewBasePathFs(fs.Source, p.Cfg.BaseConfig().WorkingDir))
	sourceFs := hugofs.NewBaseFileDecorator(fs.Source)
	publishFsStatic := fs.PublishDirStatic

	var buildMu Lockable
	if p.Cfg.NoBuildLock() || htesting.IsTest {
		buildMu = &fakeLockfileMutex{}
	} else {
		buildMu = lockedfile.MutexAt(filepath.Join(p.Cfg.BaseConfig().WorkingDir, lockFileBuild))
	}

	b := &BaseFs{
		SourceFs:        sourceFs,
		ProjectSourceFs: projectSourceFs,
		WorkDir:         fs.WorkingDirReadOnly,
		PublishFs:       publishFs,
		PublishFsStatic: publishFsStatic,
		workingDir:      p.Cfg.BaseConfig().WorkingDir,
		buildMu:         buildMu,
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
		return nil, fmt.Errorf("build filesystems: %w", err)
	}

	b.SourceFilesystems = sourceFilesystems
	b.theBigFs = builder.theBigFs

	return b, nil
}

type sourceFilesystemsBuilder struct {
	logger   loggers.Logger
	p        *paths.Paths
	sourceFs afero.Fs
	result   *SourceFilesystems
	theBigFs *filesystemsCollector
}

func newSourceFilesystemsBuilder(p *paths.Paths, logger loggers.Logger, b *BaseFs) *sourceFilesystemsBuilder {
	sourceFs := hugofs.NewBaseFileDecorator(p.Fs.Source)
	return &sourceFilesystemsBuilder{
		p: p, logger: logger, sourceFs: sourceFs, theBigFs: b.theBigFs,
		result: &SourceFilesystems{
			conf: p.Cfg,
		},
	}
}

func (b *sourceFilesystemsBuilder) newSourceFilesystem(name string, fs afero.Fs) *SourceFilesystem {
	return &SourceFilesystem{
		Name:     name,
		Fs:       fs,
		SourceFs: b.sourceFs,
	}
}

func (b *sourceFilesystemsBuilder) Build() (*SourceFilesystems, error) {
	if b.theBigFs == nil {
		theBigFs, err := b.createMainOverlayFs(b.p)
		if err != nil {
			return nil, fmt.Errorf("create main fs: %w", err)
		}

		b.theBigFs = theBigFs
	}

	createView := func(componentID string, overlayFs *overlayfs.OverlayFs) *SourceFilesystem {
		if b.theBigFs == nil || b.theBigFs.overlayMounts == nil {
			return b.newSourceFilesystem(componentID, hugofs.NoOpFs)
		}

		fs := hugofs.NewComponentFs(
			hugofs.ComponentFsOptions{
				Fs:                     overlayFs,
				Component:              componentID,
				DefaultContentLanguage: b.p.Cfg.DefaultContentLanguage(),
				PathParser:             b.p.Cfg.PathParser(),
			},
		)

		return b.newSourceFilesystem(componentID, fs)
	}

	b.result.Archetypes = createView(files.ComponentFolderArchetypes, b.theBigFs.overlayMounts)
	b.result.Layouts = createView(files.ComponentFolderLayouts, b.theBigFs.overlayMounts)
	b.result.Assets = createView(files.ComponentFolderAssets, b.theBigFs.overlayMounts)
	b.result.ResourcesCache = b.theBigFs.overlayResources
	b.result.RootFss = b.theBigFs.rootFss

	// data and i18n  needs a different merge strategy.
	overlayMountsPreserveDupes := b.theBigFs.overlayMounts.WithDirsMerger(hugofs.AppendDirsMerger)
	b.result.Data = createView(files.ComponentFolderData, overlayMountsPreserveDupes)
	b.result.I18n = createView(files.ComponentFolderI18n, overlayMountsPreserveDupes)
	b.result.AssetsWithDuplicatesPreserved = createView(files.ComponentFolderAssets, overlayMountsPreserveDupes)

	contentFs := hugofs.NewComponentFs(
		hugofs.ComponentFsOptions{
			Fs:                     b.theBigFs.overlayMountsContent,
			Component:              files.ComponentFolderContent,
			DefaultContentLanguage: b.p.Cfg.DefaultContentLanguage(),
			PathParser:             b.p.Cfg.PathParser(),
		},
	)

	b.result.Content = b.newSourceFilesystem(files.ComponentFolderContent, contentFs)
	b.result.Work = hugofs.NewReadOnlyFs(b.theBigFs.overlayFull)

	// Create static filesystem(s)
	ms := make(map[string]*SourceFilesystem)
	b.result.Static = ms

	if b.theBigFs.staticPerLanguage != nil {
		// Multihost mode
		for k, v := range b.theBigFs.staticPerLanguage {
			sfs := b.newSourceFilesystem(files.ComponentFolderStatic, v)
			sfs.PublishFolder = k
			ms[k] = sfs
		}
	} else {
		bfs := hugofs.NewBasePathFs(b.theBigFs.overlayMountsStatic, files.ComponentFolderStatic)
		ms[""] = b.newSourceFilesystem(files.ComponentFolderStatic, bfs)
	}

	return b.result, nil
}

func (b *sourceFilesystemsBuilder) createMainOverlayFs(p *paths.Paths) (*filesystemsCollector, error) {
	var staticFsMap map[string]*overlayfs.OverlayFs
	if b.p.Cfg.IsMultihost() {
		languages := b.p.Cfg.Languages()
		staticFsMap = make(map[string]*overlayfs.OverlayFs)
		for _, l := range languages {
			staticFsMap[l.Lang] = overlayfs.New(overlayfs.Options{})
		}
	}

	collector := &filesystemsCollector{
		sourceProject:     b.sourceFs,
		sourceModules:     b.sourceFs,
		staticPerLanguage: staticFsMap,

		overlayMounts:        overlayfs.New(overlayfs.Options{}),
		overlayMountsContent: overlayfs.New(overlayfs.Options{DirsMerger: hugofs.LanguageDirsMerger}),
		overlayMountsStatic:  overlayfs.New(overlayfs.Options{DirsMerger: hugofs.LanguageDirsMerger}),
		overlayFull:          overlayfs.New(overlayfs.Options{}),
		overlayResources:     overlayfs.New(overlayfs.Options{FirstWritable: true}),
	}

	mods := p.AllModules()

	mounts := make([]mountsDescriptor, len(mods))

	for i := 0; i < len(mods); i++ {
		mod := mods[i]
		dir := mod.Dir()

		isMainProject := mod.Owner() == nil
		mounts[i] = mountsDescriptor{
			Module:        mod,
			dir:           dir,
			isMainProject: isMainProject,
			ordinal:       i,
		}

	}

	err := b.createOverlayFs(collector, mounts)

	return collector, err
}

func (b *sourceFilesystemsBuilder) isContentMount(mnt modules.Mount) bool {
	return strings.HasPrefix(mnt.Target, files.ComponentFolderContent)
}

func (b *sourceFilesystemsBuilder) isStaticMount(mnt modules.Mount) bool {
	return strings.HasPrefix(mnt.Target, files.ComponentFolderStatic)
}

func (b *sourceFilesystemsBuilder) createOverlayFs(
	collector *filesystemsCollector,
	mounts []mountsDescriptor,
) error {
	if len(mounts) == 0 {
		appendNopIfEmpty := func(ofs *overlayfs.OverlayFs) *overlayfs.OverlayFs {
			if ofs.NumFilesystems() > 0 {
				return ofs
			}
			return ofs.Append(hugofs.NoOpFs)
		}
		collector.overlayMounts = appendNopIfEmpty(collector.overlayMounts)
		collector.overlayMountsContent = appendNopIfEmpty(collector.overlayMountsContent)
		collector.overlayMountsStatic = appendNopIfEmpty(collector.overlayMountsStatic)
		collector.overlayMountsFull = appendNopIfEmpty(collector.overlayMountsFull)
		collector.overlayFull = appendNopIfEmpty(collector.overlayFull)
		collector.overlayResources = appendNopIfEmpty(collector.overlayResources)

		return nil
	}

	for _, md := range mounts {
		var (
			fromTo        []hugofs.RootMapping
			fromToContent []hugofs.RootMapping
			fromToStatic  []hugofs.RootMapping
		)

		absPathify := func(path string) (string, string) {
			if filepath.IsAbs(path) {
				return "", path
			}
			return md.dir, hpaths.AbsPathify(md.dir, path)
		}

		for i, mount := range md.Mounts() {
			// Add more weight to early mounts.
			// When two mounts contain the same filename,
			// the first entry wins.
			mountWeight := (10 + md.ordinal) * (len(md.Mounts()) - i)

			inclusionFilter, err := glob.NewFilenameFilter(
				types.ToStringSlicePreserveString(mount.IncludeFiles),
				types.ToStringSlicePreserveString(mount.ExcludeFiles),
			)
			if err != nil {
				return err
			}

			base, filename := absPathify(mount.Source)

			rm := hugofs.RootMapping{
				From:          mount.Target,
				To:            filename,
				ToBase:        base,
				Module:        md.Module.Path(),
				ModuleOrdinal: md.ordinal,
				IsProject:     md.isMainProject,
				Meta: &hugofs.FileMeta{
					Watch:           !mount.DisableWatch && md.Watch(),
					Weight:          mountWeight,
					InclusionFilter: inclusionFilter,
				},
			}

			isContentMount := b.isContentMount(mount)

			lang := mount.Lang
			if lang == "" && isContentMount {
				lang = b.p.Cfg.DefaultContentLanguage()
			}

			rm.Meta.Lang = lang

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

		sourceStatic := modBase

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

		// We need to keep the list of directories for watching.
		collector.addRootFs(rmfs)
		collector.addRootFs(rmfsContent)
		collector.addRootFs(rmfsStatic)

		if collector.staticPerLanguage != nil {
			for _, l := range b.p.Cfg.Languages() {
				lang := l.Lang

				lfs := rmfsStatic.Filter(func(rm hugofs.RootMapping) bool {
					rlang := rm.Meta.Lang
					return rlang == "" || rlang == lang
				})
				bfs := hugofs.NewBasePathFs(lfs, files.ComponentFolderStatic)
				collector.staticPerLanguage[lang] = collector.staticPerLanguage[lang].Append(bfs)
			}
		}

		getResourcesDir := func() string {
			if md.isMainProject {
				return b.p.AbsResourcesDir
			}
			_, filename := absPathify(files.FolderResources)
			return filename
		}

		collector.overlayMounts = collector.overlayMounts.Append(rmfs)
		collector.overlayMountsContent = collector.overlayMountsContent.Append(rmfsContent)
		collector.overlayMountsStatic = collector.overlayMountsStatic.Append(rmfsStatic)
		collector.overlayFull = collector.overlayFull.Append(hugofs.NewBasePathFs(modBase, md.dir))
		collector.overlayResources = collector.overlayResources.Append(hugofs.NewBasePathFs(modBase, getResourcesDir()))

	}

	return nil
}

//lint:ignore U1000 useful for debugging
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
			filename = fim.Meta().Filename
		}
		fmt.Fprintf(w, "    %q %q\n", path, filename)
		return nil
	})
}

type filesystemsCollector struct {
	sourceProject afero.Fs // Source for project folders
	sourceModules afero.Fs // Source for modules/themes

	overlayMounts        *overlayfs.OverlayFs
	overlayMountsContent *overlayfs.OverlayFs
	overlayMountsStatic  *overlayfs.OverlayFs
	overlayMountsFull    *overlayfs.OverlayFs
	overlayFull          *overlayfs.OverlayFs
	overlayResources     *overlayfs.OverlayFs

	rootFss []*hugofs.RootMappingFs

	// Set if in multihost mode
	staticPerLanguage map[string]*overlayfs.OverlayFs
}

func (c *filesystemsCollector) addRootFs(rfs *hugofs.RootMappingFs) {
	c.rootFss = append(c.rootFss, rfs)
}

type mountsDescriptor struct {
	modules.Module
	dir           string
	isMainProject bool
	ordinal       int // zero based starting from the project.
}
