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
	"context"
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/gohugoio/hugo/helpers"
	"golang.org/x/sync/errgroup"
)

type bundleState int

const (
	bundleNone bundleState = iota
	bundleLeaf
	bundleBranch
)

type FileInfo struct {
	os.FileInfo
	Owner bool
}

func newFileInfo(fi os.FileInfo, owner bool) FileInfo {
	return FileInfo{FileInfo: fi, Owner: owner}
}

type dirs map[string][]FileInfo

// For testing/debugging.
func (d dirs) String() string {
	s := "\n"
	for k, v := range d {
		ss := ""
		for i, fi := range v {
			ss += fi.Name()
			if i+1 < len(v) {
				ss += "|"
			}
		}
		s += fmt.Sprintf("%s: %s\n", k, ss)
	}
	return s
}

// Basics:
// * A file in `/content` can be either a single file or part of a bundle.
// * A bundle is always a directory.
// * A bundle does not contain other bundles.
// * A bundle can contain (non-bundle) directories.
//
// So in the Hugo file and page conversion chain, the smallest _unit of work_ is a directory.
// Any smaller unit must be resolved to its parent directory.
type bundler struct {
	ps *helpers.PathSpec

	// Maps path to a list of files directly below that folder.
	dirs dirs

	// Contains paths to the page bundles.
	bundles map[string]bool
}

func (b *bundler) addBundle(dir string) {
	parts := strings.Split(dir, helpers.FilePathSeparator)
	for i := len(parts); i >= 0; i-- {
		b.bundles[filepath.Join(parts[:i]...)] = true
	}
}

func newBundler(ps *helpers.PathSpec) *bundler {
	return &bundler{ps: ps, dirs: make(map[string][]FileInfo), bundles: make(map[string]bool)}
}

func (b *bundler) processFiles() error {
	g, ctx := errgroup.WithContext(context.TODO())

	filesChan := make(chan []FileInfo)
	numWorkers := 1 // getGoMaxProcs() * 4

	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for {
				select {
				case files, more := <-filesChan:
					if !more {
						return nil
					}
					fmt.Println("Received Files:")
					for _, f := range files {
						fmt.Println(">>", f.Name(), f.Owner)
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	for _, fis := range b.dirs {
		filesChan <- fis
	}

	close(filesChan)

	if err := g.Wait(); err != nil {
		// Cancel?
		return err

	}

	return nil
}

func newBool(b bool) *bool {
	return &b
}

func (b *bundler) captureFiles() error {

	// TODO(bep) bundler check ignore logic etc.
	walker := func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		if fi.IsDir() {
			return nil
		}

		if filename == b.ps.ContentDir() {
			return nil
		}

		// Find the ancestor node to attach to.

		bState := evaluateBundle(fi)
		isBundle := bState > bundleNone

		dir := filepath.Dir(filename)

		if isBundle {
			b.addBundle(dir)
		}

		//fmt.Println(">>>", filename, dir, ">", &parent, ">", parent.isBundle)

		//bundle-site
		//├── 0-home-illustration.jpg
		//├── _index.md - OK
		//└── first-section
		//    ├── _index.md - OK
		//    ├── images - ?
		//    │   └── hugo.jpg
		//    ├── my-bundle-post
		//    │   ├── images
		//    │   │   └── hugo.jpg
		//    │   └── index.md
		//    ├── my-non-bundle-post
		//    │   ├── post-image.jpg
		//    │   └── post.md
		//    └── sub-section
		//        ├── sub-section-image.jpg
		//        └── sub-sub-section
		//            ├── _index.md
		//            └── myimage.jpg
		// From the above:
		//
		// If the lowermost path is a bundle (contains _index or index) then:
		// * All parents up to the root are also bundles.
		// * All child folders excluding bundles are part of the bundle.
		// Examples:
		// sub-sub-section (myimage.jpg) > sub-section, first-section, home

		b.dirs[dir] = append(b.dirs[dir], newFileInfo(fi, isBundle))

		return nil
	}

	return helpers.SymbolicWalk(b.ps.Fs.Source, b.ps.ContentDir(), walker)

}

func evaluateBundle(fi os.FileInfo) bundleState {
	if fi.IsDir() {
		return bundleNone
	}

	filename := fi.Name()

	if strings.HasPrefix(filename, "_index.") {
		return bundleBranch
	}

	if strings.HasPrefix(filename, "index.") {
		return bundleLeaf
	}

	return bundleNone
}

// Outline of some new processing interfaces

type PagesProcessor interface {
	ProcessPages(files ...FileInfo) (Pages, error)
}
