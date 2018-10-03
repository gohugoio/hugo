// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/gohugoio/hugo/common/loggers"
	_errors "github.com/pkg/errors"

	"sort"
	"strings"
	"sync"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/helpers"

	"golang.org/x/sync/errgroup"

	"github.com/gohugoio/hugo/source"
)

var errSkipCyclicDir = errors.New("skip potential cyclic dir")

type capturer struct {
	// To prevent symbolic link cycles: Visit same folder only once.
	seen   map[string]bool
	seenMu sync.Mutex

	handler captureResultHandler

	sourceSpec *source.SourceSpec
	fs         afero.Fs
	logger     *loggers.Logger

	// Filenames limits the content to process to a list of filenames/directories.
	// This is used for partial building in server mode.
	filenames []string

	// Used to determine how to handle content changes in server mode.
	contentChanges *contentChangeMap

	// Semaphore used to throttle the concurrent sub directory handling.
	sem chan bool
}

func newCapturer(
	logger *loggers.Logger,
	sourceSpec *source.SourceSpec,
	handler captureResultHandler,
	contentChanges *contentChangeMap,
	filenames ...string) *capturer {

	numWorkers := 4
	if n := runtime.NumCPU(); n > numWorkers {
		numWorkers = n
	}

	// TODO(bep) the "index" vs "_index" check/strings should be moved in one place.
	isBundleHeader := func(filename string) bool {
		base := filepath.Base(filename)
		name := helpers.Filename(base)
		return IsContentFile(base) && (name == "index" || name == "_index")
	}

	// Make sure that any bundle header files are processed before the others. This makes
	// sure that any bundle head is processed before its resources.
	sort.Slice(filenames, func(i, j int) bool {
		a, b := filenames[i], filenames[j]
		ac, bc := isBundleHeader(a), isBundleHeader(b)

		if ac {
			return true
		}

		if bc {
			return false
		}

		return a < b
	})

	c := &capturer{
		sem:            make(chan bool, numWorkers),
		handler:        handler,
		sourceSpec:     sourceSpec,
		fs:             sourceSpec.SourceFs,
		logger:         logger,
		contentChanges: contentChanges,
		seen:           make(map[string]bool),
		filenames:      filenames}

	return c
}

// Captured files and bundles ready to be processed will be passed on to
// these channels.
type captureResultHandler interface {
	handleSingles(fis ...*fileInfo)
	handleCopyFiles(fis ...pathLangFile)
	captureBundlesHandler
}

type captureBundlesHandler interface {
	handleBundles(b *bundleDirs)
}

type captureResultHandlerChain struct {
	handlers []captureBundlesHandler
}

func (c *captureResultHandlerChain) handleSingles(fis ...*fileInfo) {
	for _, h := range c.handlers {
		if hh, ok := h.(captureResultHandler); ok {
			hh.handleSingles(fis...)
		}
	}
}
func (c *captureResultHandlerChain) handleBundles(b *bundleDirs) {
	for _, h := range c.handlers {
		h.handleBundles(b)
	}
}

func (c *captureResultHandlerChain) handleCopyFiles(files ...pathLangFile) {
	for _, h := range c.handlers {
		if hh, ok := h.(captureResultHandler); ok {
			hh.handleCopyFiles(files...)
		}
	}
}

func (c *capturer) capturePartial(filenames ...string) error {
	handled := make(map[string]bool)

	for _, filename := range filenames {
		dir, resolvedFilename, tp := c.contentChanges.resolveAndRemove(filename)
		if handled[resolvedFilename] {
			continue
		}

		handled[resolvedFilename] = true

		switch tp {
		case bundleLeaf:
			if err := c.handleDir(resolvedFilename); err != nil {
				// Directory may have been deleted.
				if !os.IsNotExist(err) {
					return err
				}
			}
		case bundleBranch:
			if err := c.handleBranchDir(resolvedFilename); err != nil {
				// Directory may have been deleted.
				if !os.IsNotExist(err) {
					return err
				}
			}
		default:
			fi, err := c.resolveRealPath(resolvedFilename)
			if os.IsNotExist(err) {
				// File has been deleted.
				continue
			}

			// Just in case the owning dir is a new symlink -- this will
			// create the proper mapping for it.
			c.resolveRealPath(dir)

			f, active := c.newFileInfo(fi, tp)
			if active {
				c.copyOrHandleSingle(f)
			}
		}
	}

	return nil
}

func (c *capturer) capture() error {
	if len(c.filenames) > 0 {
		return c.capturePartial(c.filenames...)
	}

	err := c.handleDir(helpers.FilePathSeparator)
	if err != nil {
		return err
	}

	return nil
}

func (c *capturer) handleNestedDir(dirname string) error {
	select {
	case c.sem <- true:
		var g errgroup.Group

		g.Go(func() error {
			defer func() {
				<-c.sem
			}()
			return c.handleDir(dirname)
		})
		return g.Wait()
	default:
		// For deeply nested file trees, waiting for a semaphore wil deadlock.
		return c.handleDir(dirname)
	}
}

// This handles a bundle branch and its resources only. This is used
// in server mode on changes. If this dir does not (anymore) represent a bundle
// branch, the handling is upgraded to the full handleDir method.
func (c *capturer) handleBranchDir(dirname string) error {
	files, err := c.readDir(dirname)
	if err != nil {

		return err
	}

	var (
		dirType bundleDirType
	)

	for _, fi := range files {
		if !fi.IsDir() {
			tp, _ := classifyBundledFile(fi.RealName())
			if dirType == bundleNot {
				dirType = tp
			}

			if dirType == bundleLeaf {
				return c.handleDir(dirname)
			}
		}
	}

	if dirType != bundleBranch {
		return c.handleDir(dirname)
	}

	dirs := newBundleDirs(bundleBranch, c)

	var secondPass []*fileInfo

	// Handle potential bundle headers first.
	for _, fi := range files {
		if fi.IsDir() {
			continue
		}

		tp, isContent := classifyBundledFile(fi.RealName())

		f, active := c.newFileInfo(fi, tp)

		if !active {
			continue
		}

		if !f.isOwner() {
			if !isContent {
				// This is a partial update -- we only care about the files that
				// is in this bundle.
				secondPass = append(secondPass, f)
			}
			continue
		}
		dirs.addBundleHeader(f)
	}

	for _, f := range secondPass {
		dirs.addBundleFiles(f)
	}

	c.handler.handleBundles(dirs)

	return nil

}

func (c *capturer) handleDir(dirname string) error {

	files, err := c.readDir(dirname)
	if err != nil {
		return err
	}

	type dirState int

	const (
		dirStateDefault dirState = iota

		dirStateAssetsOnly
		dirStateSinglesOnly
	)

	var (
		fileBundleTypes = make([]bundleDirType, len(files))

		// Start with the assumption that this dir contains only non-content assets (images etc.)
		// If that is still true after we had a first look at the list of files, we
		// can just copy the files to destination. We will still have to look at the
		// sub-folders for potential bundles.
		state = dirStateAssetsOnly

		// Start with the assumption that this dir is not a bundle.
		// A directory is a bundle if it contains a index content file,
		// e.g. index.md (a leaf bundle) or a _index.md (a branch bundle).
		bundleType = bundleNot
	)

	/* First check for any content files.
	- If there are none, then this is a assets folder only (images etc.)
	and we can just plainly copy them to
	destination.
	- If this is a section with no image etc. or similar, we can just handle it
	as it was a single content file.
	*/
	var hasNonContent, isBranch bool

	for i, fi := range files {
		if !fi.IsDir() {
			tp, isContent := classifyBundledFile(fi.RealName())

			fileBundleTypes[i] = tp
			if !isBranch {
				isBranch = tp == bundleBranch
			}

			if isContent {
				// This is not a assets-only folder.
				state = dirStateDefault
			} else {
				hasNonContent = true
			}
		}
	}

	if isBranch && !hasNonContent {
		// This is a section or similar with no need for any bundle handling.
		state = dirStateSinglesOnly
	}

	if state > dirStateDefault {
		return c.handleNonBundle(dirname, files, state == dirStateSinglesOnly)
	}

	var fileInfos = make([]*fileInfo, 0, len(files))

	for i, fi := range files {

		currentType := bundleNot

		if !fi.IsDir() {
			currentType = fileBundleTypes[i]
			if bundleType == bundleNot && currentType != bundleNot {
				bundleType = currentType
			}
		}

		if bundleType == bundleNot && currentType != bundleNot {
			bundleType = currentType
		}

		f, active := c.newFileInfo(fi, currentType)

		if !active {
			continue
		}

		fileInfos = append(fileInfos, f)
	}

	var todo []*fileInfo

	if bundleType != bundleLeaf {
		for _, fi := range fileInfos {
			if fi.FileInfo().IsDir() {
				// Handle potential nested bundles.
				if err := c.handleNestedDir(fi.Path()); err != nil {
					return err
				}
			} else if bundleType == bundleNot || (!fi.isOwner() && fi.isContentFile()) {
				// Not in a bundle.
				c.copyOrHandleSingle(fi)
			} else {
				// This is a section folder or similar with non-content files in it.
				todo = append(todo, fi)
			}
		}
	} else {
		todo = fileInfos
	}

	if len(todo) == 0 {
		return nil
	}

	dirs, err := c.createBundleDirs(todo, bundleType)
	if err != nil {
		return err
	}

	// Send the bundle to the next step in the processor chain.
	c.handler.handleBundles(dirs)

	return nil
}

func (c *capturer) handleNonBundle(
	dirname string,
	fileInfos pathLangFileFis,
	singlesOnly bool) error {

	for _, fi := range fileInfos {
		if fi.IsDir() {
			if err := c.handleNestedDir(fi.Filename()); err != nil {
				return err
			}
		} else {
			if singlesOnly {
				f, active := c.newFileInfo(fi, bundleNot)
				if !active {
					continue
				}
				c.handler.handleSingles(f)
			} else {
				c.handler.handleCopyFiles(fi)
			}
		}
	}

	return nil
}

func (c *capturer) copyOrHandleSingle(fi *fileInfo) {
	if fi.isContentFile() {
		c.handler.handleSingles(fi)
	} else {
		// These do not currently need any further processing.
		c.handler.handleCopyFiles(fi)
	}
}

func (c *capturer) createBundleDirs(fileInfos []*fileInfo, bundleType bundleDirType) (*bundleDirs, error) {
	dirs := newBundleDirs(bundleType, c)

	for _, fi := range fileInfos {
		if fi.FileInfo().IsDir() {
			var collector func(fis ...*fileInfo)

			if bundleType == bundleBranch {
				// All files in the current directory are part of this bundle.
				// Trying to include sub folders in these bundles are filled with ambiguity.
				collector = func(fis ...*fileInfo) {
					for _, fi := range fis {
						c.copyOrHandleSingle(fi)
					}
				}
			} else {
				// All nested files and directories are part of this bundle.
				collector = func(fis ...*fileInfo) {
					fileInfos = append(fileInfos, fis...)
				}
			}
			err := c.collectFiles(fi.Path(), collector)
			if err != nil {
				return nil, err
			}

		} else if fi.isOwner() {
			// There can be more than one language, so:
			// 1. Content files must be attached to its language's bundle.
			// 2. Other files must be attached to all languages.
			// 3. Every content file needs a bundle header.
			dirs.addBundleHeader(fi)
		}
	}

	for _, fi := range fileInfos {
		if fi.FileInfo().IsDir() || fi.isOwner() {
			continue
		}

		if fi.isContentFile() {
			if bundleType != bundleBranch {
				dirs.addBundleContentFile(fi)
			}
		} else {
			dirs.addBundleFiles(fi)
		}
	}

	return dirs, nil
}

func (c *capturer) collectFiles(dirname string, handleFiles func(fis ...*fileInfo)) error {

	filesInDir, err := c.readDir(dirname)
	if err != nil {
		return err
	}

	for _, fi := range filesInDir {
		if fi.IsDir() {
			err := c.collectFiles(fi.Filename(), handleFiles)
			if err != nil {
				return err
			}
		} else {
			f, active := c.newFileInfo(fi, bundleNot)
			if active {
				handleFiles(f)
			}
		}
	}

	return nil
}

func (c *capturer) readDir(dirname string) (pathLangFileFis, error) {
	if c.sourceSpec.IgnoreFile(dirname) {
		return nil, nil
	}

	dir, err := c.fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	fis, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	pfis := make(pathLangFileFis, 0, len(fis))

	for _, fi := range fis {
		fip := fi.(pathLangFileFi)

		if !c.sourceSpec.IgnoreFile(fip.Filename()) {

			err := c.resolveRealPathIn(fip)

			if err != nil {
				// It may have been deleted in the meantime.
				if err == errSkipCyclicDir || os.IsNotExist(err) {
					continue
				}
				return nil, err
			}

			pfis = append(pfis, fip)
		}
	}

	return pfis, nil
}

func (c *capturer) newFileInfo(fi pathLangFileFi, tp bundleDirType) (*fileInfo, bool) {
	f := newFileInfo(c.sourceSpec, "", "", fi, tp)
	return f, !f.disabled
}

type pathLangFile interface {
	hugofs.LanguageAnnouncer
	hugofs.FilePather
}

type pathLangFileFi interface {
	os.FileInfo
	pathLangFile
}

type pathLangFileFis []pathLangFileFi

type bundleDirs struct {
	tp bundleDirType
	// Maps languages to bundles.
	bundles map[string]*bundleDir

	// Keeps track of language overrides for non-content files, e.g. logo.en.png.
	langOverrides map[string]bool

	c *capturer
}

func newBundleDirs(tp bundleDirType, c *capturer) *bundleDirs {
	return &bundleDirs{tp: tp, bundles: make(map[string]*bundleDir), langOverrides: make(map[string]bool), c: c}
}

type bundleDir struct {
	tp bundleDirType
	fi *fileInfo

	resources map[string]*fileInfo
}

func (b bundleDir) clone() *bundleDir {
	b.resources = make(map[string]*fileInfo)
	fic := *b.fi
	b.fi = &fic
	return &b
}

func newBundleDir(fi *fileInfo, bundleType bundleDirType) *bundleDir {
	return &bundleDir{fi: fi, tp: bundleType, resources: make(map[string]*fileInfo)}
}

func (b *bundleDirs) addBundleContentFile(fi *fileInfo) {
	dir, found := b.bundles[fi.Lang()]
	if !found {
		// Every bundled content file needs a bundle header.
		// If one does not exist in its language, we pick the default
		// language version, or a random one if that doesn't exist, either.
		tl := b.c.sourceSpec.DefaultContentLanguage
		ldir, found := b.bundles[tl]
		if !found {
			// Just pick one.
			for _, v := range b.bundles {
				ldir = v
				break
			}
		}

		if ldir == nil {
			panic(fmt.Sprintf("bundle not found for file %q", fi.Filename()))
		}

		dir = ldir.clone()
		dir.fi.overriddenLang = fi.Lang()
		b.bundles[fi.Lang()] = dir
	}

	dir.resources[fi.Path()] = fi
}

func (b *bundleDirs) addBundleFiles(fi *fileInfo) {
	dir := filepath.ToSlash(fi.Dir())
	p := dir + fi.TranslationBaseName() + "." + fi.Ext()
	for lang, bdir := range b.bundles {
		key := path.Join(lang, p)

		// Given mypage.de.md (German translation) and mypage.md we pick the most
		// specific for that language.
		if fi.Lang() == lang || !b.langOverrides[key] {
			bdir.resources[key] = fi
		}
		b.langOverrides[key] = true
	}
}

func (b *bundleDirs) addBundleHeader(fi *fileInfo) {
	b.bundles[fi.Lang()] = newBundleDir(fi, b.tp)
}

func (c *capturer) isSeen(dirname string) bool {
	c.seenMu.Lock()
	defer c.seenMu.Unlock()
	seen := c.seen[dirname]
	c.seen[dirname] = true
	if seen {
		c.logger.WARN.Printf("Content dir %q already processed; skipped to avoid infinite recursion.", dirname)
		return true

	}
	return false
}

func (c *capturer) resolveRealPath(path string) (pathLangFileFi, error) {
	fileInfo, err := c.lstatIfPossible(path)
	if err != nil {
		return nil, err
	}
	return fileInfo, c.resolveRealPathIn(fileInfo)
}

func (c *capturer) resolveRealPathIn(fileInfo pathLangFileFi) error {

	basePath := fileInfo.BaseDir()
	path := fileInfo.Filename()

	realPath := path

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		link, err := filepath.EvalSymlinks(path)
		if err != nil {
			return _errors.Wrapf(err, "Cannot read symbolic link %q, error was:", path)
		}

		// This is a file on the outside of any base fs, so we have to use the os package.
		sfi, err := os.Stat(link)
		if err != nil {
			return _errors.Wrapf(err, "Cannot stat  %q, error was:", link)
		}

		// TODO(bep) improve all of this.
		if a, ok := fileInfo.(*hugofs.LanguageFileInfo); ok {
			a.FileInfo = sfi
		}

		realPath = link

		if realPath != path && sfi.IsDir() && c.isSeen(realPath) {
			// Avoid cyclic symlinks.
			// Note that this may prevent some uses that isn't cyclic and also
			// potential useful, but this implementation is both robust and simple:
			// We stop at the first directory that we have seen before, e.g.
			// /content/blog will only be processed once.
			return errSkipCyclicDir
		}

		if c.contentChanges != nil {
			// Keep track of symbolic links in watch mode.
			var from, to string
			if sfi.IsDir() {
				from = realPath
				to = path

				if !strings.HasSuffix(to, helpers.FilePathSeparator) {
					to = to + helpers.FilePathSeparator
				}
				if !strings.HasSuffix(from, helpers.FilePathSeparator) {
					from = from + helpers.FilePathSeparator
				}

				if !strings.HasSuffix(basePath, helpers.FilePathSeparator) {
					basePath = basePath + helpers.FilePathSeparator
				}

				if strings.HasPrefix(from, basePath) {
					// With symbolic links inside /content we need to keep
					// a reference to both. This may be confusing with --navigateToChanged
					// but the user has chosen this him or herself.
					c.contentChanges.addSymbolicLinkMapping(from, from)
				}

			} else {
				from = realPath
				to = path
			}

			c.contentChanges.addSymbolicLinkMapping(from, to)
		}
	}

	return nil
}

func (c *capturer) lstatIfPossible(path string) (pathLangFileFi, error) {
	fi, err := helpers.LstatIfPossible(c.fs, path)
	if err != nil {
		return nil, err
	}
	return fi.(pathLangFileFi), nil
}
