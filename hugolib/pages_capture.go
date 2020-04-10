// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"context"
	"fmt"
	"os"
	pth "path"
	"path/filepath"
	"reflect"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/parser/pageparser"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
)

const (
	walkIsRootFileMetaKey = "walkIsRootFileMetaKey"
)

func newPagesCollector(
	sp *source.SourceSpec,
	contentMap *pageMaps,
	logger *loggers.Logger,
	contentTracker *contentChangeMap,
	proc pagesCollectorProcessorProvider, filenames ...string) *pagesCollector {

	return &pagesCollector{
		fs:         sp.SourceFs,
		contentMap: contentMap,
		proc:       proc,
		sp:         sp,
		logger:     logger,
		filenames:  filenames,
		tracker:    contentTracker,
	}
}

type contentDirKey struct {
	dirname  string
	filename string
	tp       bundleDirType
}

type fileinfoBundle struct {
	header    hugofs.FileMetaInfo
	resources []hugofs.FileMetaInfo
}

func (b *fileinfoBundle) containsResource(name string) bool {
	for _, r := range b.resources {
		if r.Name() == name {
			return true
		}
	}

	return false

}

type pageBundles map[string]*fileinfoBundle

type pagesCollector struct {
	sp     *source.SourceSpec
	fs     afero.Fs
	logger *loggers.Logger

	contentMap *pageMaps

	// Ordered list (bundle headers first) used in partial builds.
	filenames []string

	// Content files tracker used in partial builds.
	tracker *contentChangeMap

	proc pagesCollectorProcessorProvider
}

// isCascadingEdit returns whether the dir represents a cascading edit.
// That is, if a front matter cascade section is removed, added or edited.
// If this is the case we must re-evaluate its descendants.
func (c *pagesCollector) isCascadingEdit(dir contentDirKey) (bool, string) {
	// This is eiter a section or a taxonomy node. Find it.
	prefix := cleanTreeKey(dir.dirname)

	section := "/"
	var isCascade bool

	c.contentMap.walkBranchesPrefix(prefix, func(s string, n *contentNode) bool {
		if n.fi == nil || dir.filename != n.fi.Meta().Filename() {
			return false
		}

		f, err := n.fi.Meta().Open()
		if err != nil {
			// File may have been removed, assume a cascading edit.
			// Some false positives is not too bad.
			isCascade = true
			return true
		}

		pf, err := pageparser.ParseFrontMatterAndContent(f)
		f.Close()
		if err != nil {
			isCascade = true
			return true
		}

		if n.p == nil || n.p.bucket == nil {
			return true
		}

		section = s

		maps.ToLower(pf.FrontMatter)
		cascade1, ok := pf.FrontMatter["cascade"]
		hasCascade := n.p.bucket.cascade != nil && len(n.p.bucket.cascade) > 0
		if !ok {
			isCascade = hasCascade
			return true
		}

		if !hasCascade {
			isCascade = true
			return true
		}

		isCascade = !reflect.DeepEqual(cascade1, n.p.bucket.cascade)

		return true

	})

	return isCascade, section
}

// Collect.
func (c *pagesCollector) Collect() (collectErr error) {
	c.proc.Start(context.Background())
	defer func() {
		err := c.proc.Wait()
		if collectErr == nil {
			collectErr = err
		}
	}()

	if len(c.filenames) == 0 {
		// Collect everything.
		collectErr = c.collectDir("", false, nil)
	} else {
		for _, pm := range c.contentMap.pmaps {
			pm.cfg.isRebuild = true
		}
		dirs := make(map[contentDirKey]bool)
		for _, filename := range c.filenames {
			dir, btype := c.tracker.resolveAndRemove(filename)
			dirs[contentDirKey{dir, filename, btype}] = true
		}

		for dir := range dirs {
			for _, pm := range c.contentMap.pmaps {
				pm.s.ResourceSpec.DeleteBySubstring(dir.dirname)
			}

			switch dir.tp {
			case bundleLeaf:
				collectErr = c.collectDir(dir.dirname, true, nil)
			case bundleBranch:
				isCascading, section := c.isCascadingEdit(dir)
				if isCascading {
					c.contentMap.deleteSection(section)
				}
				collectErr = c.collectDir(dir.dirname, !isCascading, nil)
			default:
				// We always start from a directory.
				collectErr = c.collectDir(dir.dirname, true, func(fim hugofs.FileMetaInfo) bool {
					return dir.filename == fim.Meta().Filename()
				})
			}

			if collectErr != nil {
				break
			}
		}

	}

	return

}

func (c *pagesCollector) isBundleHeader(fi hugofs.FileMetaInfo) bool {
	class := fi.Meta().Classifier()
	return class == files.ContentClassLeaf || class == files.ContentClassBranch
}

func (c *pagesCollector) getLang(fi hugofs.FileMetaInfo) string {
	lang := fi.Meta().Lang()
	if lang != "" {
		return lang
	}

	return c.sp.DefaultContentLanguage
}

func (c *pagesCollector) addToBundle(info hugofs.FileMetaInfo, btyp bundleDirType, bundles pageBundles) error {
	getBundle := func(lang string) *fileinfoBundle {
		return bundles[lang]
	}

	cloneBundle := func(lang string) *fileinfoBundle {
		// Every bundled content file needs a content file header.
		// Use the default content language if found, else just
		// pick one.
		var (
			source *fileinfoBundle
			found  bool
		)

		source, found = bundles[c.sp.DefaultContentLanguage]
		if !found {
			for _, b := range bundles {
				source = b
				break
			}
		}

		if source == nil {
			panic(fmt.Sprintf("no source found, %d", len(bundles)))
		}

		clone := c.cloneFileInfo(source.header)
		clone.Meta()["lang"] = lang

		return &fileinfoBundle{
			header: clone,
		}
	}

	lang := c.getLang(info)
	bundle := getBundle(lang)
	isBundleHeader := c.isBundleHeader(info)
	if bundle != nil && isBundleHeader {
		// index.md file inside a bundle, see issue 6208.
		info.Meta()["classifier"] = files.ContentClassContent
		isBundleHeader = false
	}
	classifier := info.Meta().Classifier()
	isContent := classifier == files.ContentClassContent
	if bundle == nil {
		if isBundleHeader {
			bundle = &fileinfoBundle{header: info}
			bundles[lang] = bundle
		} else {
			if btyp == bundleBranch {
				// No special logic for branch bundles.
				// Every language needs its own _index.md file.
				// Also, we only clone bundle headers for lonsesome, bundled,
				// content files.
				return c.handleFiles(info)
			}

			if isContent {
				bundle = cloneBundle(lang)
				bundles[lang] = bundle
			}
		}
	}

	if !isBundleHeader && bundle != nil {
		bundle.resources = append(bundle.resources, info)
	}

	if classifier == files.ContentClassFile {
		translations := info.Meta().Translations()

		for lang, b := range bundles {
			if !stringSliceContains(lang, translations...) && !b.containsResource(info.Name()) {

				// Clone and add it to the bundle.
				clone := c.cloneFileInfo(info)
				clone.Meta()["lang"] = lang
				b.resources = append(b.resources, clone)
			}
		}
	}

	return nil
}

func (c *pagesCollector) cloneFileInfo(fi hugofs.FileMetaInfo) hugofs.FileMetaInfo {
	cm := hugofs.FileMeta{}
	meta := fi.Meta()
	if meta == nil {
		panic(fmt.Sprintf("not meta: %v", fi.Name()))
	}
	for k, v := range meta {
		cm[k] = v
	}

	return hugofs.NewFileMetaInfo(fi, cm)
}

func (c *pagesCollector) collectDir(dirname string, partial bool, inFilter func(fim hugofs.FileMetaInfo) bool) error {
	fi, err := c.fs.Stat(dirname)
	if err != nil {
		if os.IsNotExist(err) {
			// May have been deleted.
			return nil
		}
		return err
	}

	handleDir := func(
		btype bundleDirType,
		dir hugofs.FileMetaInfo,
		path string,
		readdir []hugofs.FileMetaInfo) error {

		if btype > bundleNot && c.tracker != nil {
			c.tracker.add(path, btype)
		}

		if btype == bundleBranch {
			if err := c.handleBundleBranch(readdir); err != nil {
				return err
			}
			// A branch bundle is only this directory level, so keep walking.
			return nil
		} else if btype == bundleLeaf {
			if err := c.handleBundleLeaf(dir, path, readdir); err != nil {
				return err
			}

			return nil
		}

		if err := c.handleFiles(readdir...); err != nil {
			return err
		}

		return nil

	}

	filter := func(fim hugofs.FileMetaInfo) bool {
		if fim.Meta().SkipDir() {
			return false
		}

		if c.sp.IgnoreFile(fim.Meta().Filename()) {
			return false
		}

		if inFilter != nil {
			return inFilter(fim)
		}
		return true
	}

	preHook := func(dir hugofs.FileMetaInfo, path string, readdir []hugofs.FileMetaInfo) ([]hugofs.FileMetaInfo, error) {
		var btype bundleDirType

		filtered := readdir[:0]
		for _, fi := range readdir {
			if filter(fi) {
				filtered = append(filtered, fi)

				if c.tracker != nil {
					// Track symlinks.
					c.tracker.addSymbolicLinkMapping(fi)
				}
			}
		}
		walkRoot := dir.Meta().GetBool(walkIsRootFileMetaKey)
		readdir = filtered

		// We merge language directories, so there can be duplicates, but they
		// will be ordered, most important first.
		var duplicates []int
		seen := make(map[string]bool)

		for i, fi := range readdir {

			if fi.IsDir() {
				continue
			}

			meta := fi.Meta()
			if walkRoot {
				meta[walkIsRootFileMetaKey] = true
			}
			class := meta.Classifier()
			translationBase := meta.TranslationBaseNameWithExt()
			key := pth.Join(meta.Lang(), translationBase)

			if seen[key] {
				duplicates = append(duplicates, i)
				continue
			}
			seen[key] = true

			var thisBtype bundleDirType

			switch class {
			case files.ContentClassLeaf:
				thisBtype = bundleLeaf
			case files.ContentClassBranch:
				thisBtype = bundleBranch
			}

			// Folders with both index.md and _index.md type of files have
			// undefined behaviour and can never work.
			// The branch variant will win because of sort order, but log
			// a warning about it.
			if thisBtype > bundleNot && btype > bundleNot && thisBtype != btype {
				c.logger.WARN.Printf("Content directory %q have both index.* and _index.* files, pick one.", dir.Meta().Filename())
				// Reclassify it so it will be handled as a content file inside the
				// section, which is in line with the <= 0.55 behaviour.
				meta["classifier"] = files.ContentClassContent
			} else if thisBtype > bundleNot {
				btype = thisBtype
			}

		}

		if len(duplicates) > 0 {
			for i := len(duplicates) - 1; i >= 0; i-- {
				idx := duplicates[i]
				readdir = append(readdir[:idx], readdir[idx+1:]...)
			}
		}

		err := handleDir(btype, dir, path, readdir)
		if err != nil {
			return nil, err
		}

		if btype == bundleLeaf || partial {
			return nil, filepath.SkipDir
		}

		// Keep walking.
		return readdir, nil

	}

	var postHook hugofs.WalkHook
	if c.tracker != nil {
		postHook = func(dir hugofs.FileMetaInfo, path string, readdir []hugofs.FileMetaInfo) ([]hugofs.FileMetaInfo, error) {
			if c.tracker == nil {
				// Nothing to do.
				return readdir, nil
			}

			return readdir, nil
		}
	}

	wfn := func(path string, info hugofs.FileMetaInfo, err error) error {
		if err != nil {
			return err
		}

		return nil
	}

	fim := fi.(hugofs.FileMetaInfo)
	// Make sure the pages in this directory gets re-rendered,
	// even in fast render mode.
	fim.Meta()[walkIsRootFileMetaKey] = true

	w := hugofs.NewWalkway(hugofs.WalkwayConfig{
		Fs:       c.fs,
		Logger:   c.logger,
		Root:     dirname,
		Info:     fim,
		HookPre:  preHook,
		HookPost: postHook,
		WalkFn:   wfn})

	return w.Walk()

}

func (c *pagesCollector) handleBundleBranch(readdir []hugofs.FileMetaInfo) error {

	// Maps bundles to its language.
	bundles := pageBundles{}

	var contentFiles []hugofs.FileMetaInfo

	for _, fim := range readdir {

		if fim.IsDir() {
			continue
		}

		meta := fim.Meta()

		switch meta.Classifier() {
		case files.ContentClassContent:
			contentFiles = append(contentFiles, fim)
		default:
			if err := c.addToBundle(fim, bundleBranch, bundles); err != nil {
				return err
			}
		}

	}

	// Make sure the section is created before its pages.
	if err := c.proc.Process(bundles); err != nil {
		return err
	}

	return c.handleFiles(contentFiles...)

}

func (c *pagesCollector) handleBundleLeaf(dir hugofs.FileMetaInfo, path string, readdir []hugofs.FileMetaInfo) error {
	// Maps bundles to its language.
	bundles := pageBundles{}

	walk := func(path string, info hugofs.FileMetaInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		return c.addToBundle(info, bundleLeaf, bundles)

	}

	// Start a new walker from the given path.
	w := hugofs.NewWalkway(hugofs.WalkwayConfig{
		Root:       path,
		Fs:         c.fs,
		Logger:     c.logger,
		Info:       dir,
		DirEntries: readdir,
		WalkFn:     walk})

	if err := w.Walk(); err != nil {
		return err
	}

	return c.proc.Process(bundles)

}

func (c *pagesCollector) handleFiles(fis ...hugofs.FileMetaInfo) error {
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		if err := c.proc.Process(fi); err != nil {
			return err
		}
	}
	return nil
}

func stringSliceContains(k string, values ...string) bool {
	for _, v := range values {
		if k == v {
			return true
		}
	}
	return false
}
