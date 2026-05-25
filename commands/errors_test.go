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

package commands

import (
	"io/fs"
	"os"
	"syscall"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestWrapStaticSyncError(t *testing.T) {
	c := qt.New(t)

	err := wrapStaticSyncError(&os.PathError{Op: "chtimes", Path: "public", Err: syscall.EPERM})
	c.Assert(err, qt.ErrorMatches, `.*--noTimes.*noTimes = true`)

	err = wrapStaticSyncError(&os.PathError{Op: "chmod", Path: "public", Err: syscall.EPERM})
	c.Assert(err.Error(), qt.Equals, "chmod public: operation not permitted")

	err = wrapStaticSyncError(&fs.PathError{Op: "chtimes", Path: "public", Err: syscall.ENOENT})
	c.Assert(err.Error(), qt.Equals, "chtimes public: no such file or directory")
}
