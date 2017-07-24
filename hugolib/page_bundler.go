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
	"fmt"
	"os"

	"github.com/gohugoio/hugo/helpers"
)

// Basics:
// * A file in `/content` can be either a single file or part of a bundle.
// * A bundle is always a directory.
//
// So in the Hugo file and page conversion chain, the smallest _unit of work_ is a directory.
// Any smaller unit must be resolved to its parent directory.
type bundler struct {
	ps *helpers.PathSpec
}

type parent struct {
	path string
	b    *bundler
}

func (b *bundler) newParent(path string) *parent {
	return &parent{path: path, b: b}
}

func (p *parent) String() string {
	return p.path
}

func (p *parent) walk() error {

	walker := func(filePath string, fi os.FileInfo, err error) error {
		fmt.Println(">>>", p.String(), ">", fi.IsDir(), filePath)
		return nil
	}

	return helpers.SymbolicWalk(p.b.ps.Fs.Source, p.path, walker)

}

func newBundler(ps *helpers.PathSpec) *bundler {
	return &bundler{ps: ps}
}

func (b *bundler) captureFiles() error {

	// TODO(bep) bundler check ignore logic etc.
	rootDirWalker := func(filePath string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !fi.IsDir() {
			return nil
		}

		parent := b.newParent(filePath)

		return parent.walk()
	}

	return helpers.SymbolicWalk(b.ps.Fs.Source, b.ps.ContentDir(), rootDirWalker)

}
