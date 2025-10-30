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

package hugofs

import (
	"io/fs"

	"github.com/bep/overlayfs"
)

// AppendDirsMerger merges two directories keeping all regular files
// with the first slice as the base.
// Duplicate directories in the second slice will be ignored.
var AppendDirsMerger overlayfs.DirsMerger = func(lofi, bofi []fs.DirEntry) []fs.DirEntry {
	for _, fi1 := range bofi {
		var found bool
		// Remove duplicate directories.
		if fi1.IsDir() {
			for _, fi2 := range lofi {
				if fi2.IsDir() && fi2.Name() == fi1.Name() {
					found = true
					break
				}
			}
		}
		if !found {
			lofi = append(lofi, fi1)
		}
	}

	return lofi
}
