// Copyright 2026 The Hugo Authors. All rights reserved.
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

package hexec

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Node's permission model checks the lexical path only and follows symlinks
// even when they point outside the allowed set (documented as a known
// limitation). checkNodeSymlinks walks the allowed roots and fails on any
// symlink whose target resolves outside them. Each root is walked once per
// Exec; new symlinks created later are not seen.
func (e *Exec) checkNodeSymlinks(kind string, roots []string) error {
	if slices.Contains(roots, "*") {
		return nil
	}

	real := evalSymlinks(roots)

	for _, root := range topLevelDirs(real) {
		if _, err := e.nodeSymlinkChecks.GetOrCreate(kind+"|"+root, func() (bool, error) {
			return true, walkSymlinks(root, real, kind)
		}); err != nil {
			return err
		}
	}
	return nil
}

// CheckReadPath fails if filename resolves outside the allowRead roots.
// Used by in-process tools (e.g. ESBuild) that read files on their own,
// so they get the same boundary as the Node.js tools.
func (e *Exec) CheckReadPath(filename string) error {
	_, err := e.readPathChecks.GetOrCreate(filename, func() (bool, error) {
		roots := e.nodeReadRoots("")
		if slices.Contains(roots, "*") {
			return true, nil
		}
		target, err := filepath.EvalSymlinks(filename)
		if err != nil {
			return false, err
		}
		if !isBelowAny(target, evalSymlinks(roots)) {
			return false, fmt.Errorf("%q resolves to %q outside the project; only files below the working directory or the paths in security.node.permissions.allowRead can be imported", filename, target)
		}
		return true, nil
	})
	return err
}

// evalSymlinks resolves paths, dropping the ones that cannot be resolved.
func evalSymlinks(paths []string) []string {
	var real []string
	for _, p := range paths {
		if p, err := filepath.EvalSymlinks(p); err == nil {
			real = append(real, p)
		}
	}
	return real
}

func walkSymlinks(root string, allowed []string, kind string) error {
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if p == root {
				return nil
			}
			return err
		}
		if d.Type()&os.ModeSymlink == 0 {
			return nil
		}
		target, err := filepath.EvalSymlinks(p)
		if err != nil {
			// Dangling; nothing to read through.
			return nil
		}
		if !isBelowAny(target, allowed) {
			return fmt.Errorf("symlink %q resolves to %q outside the paths Node.js is allowed to access; remove it or add its target to security.node.permissions.%s", p, target, kind)
		}
		return nil
	})
}

// topLevelDirs drops dirs nested below another entry so they are not walked twice.
func topLevelDirs(dirs []string) []string {
	sorted := slices.Clone(dirs)
	slices.Sort(sorted)
	var top []string
	for _, d := range sorted {
		if len(top) > 0 && isBelow(d, top[len(top)-1]) {
			continue
		}
		top = append(top, d)
	}
	return top
}

func isBelowAny(p string, dirs []string) bool {
	return slices.ContainsFunc(dirs, func(d string) bool { return isBelow(p, d) })
}

func isBelow(p, dir string) bool {
	return p == dir || strings.HasPrefix(p, strings.TrimSuffix(dir, string(filepath.Separator))+string(filepath.Separator))
}
