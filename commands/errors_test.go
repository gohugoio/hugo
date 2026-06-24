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

package commands

import (
	"errors"
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestWrapStaticSyncError(t *testing.T) {
	c := qt.New(t)

	err := &os.PathError{Op: "chtimes", Path: "public", Err: os.ErrPermission}
	got := wrapStaticSyncError(err)

	c.Assert(got, qt.ErrorMatches, `chtimes public: permission denied; you have insufficient permissions to update timestamps in publishDir; try --noTimes or set noTimes = true, and check publishDir ownership and permissions`)
	c.Assert(errors.Is(got, os.ErrPermission), qt.Equals, true)

	err = &os.PathError{Op: "stat", Path: "public/favicon.ico", Err: os.ErrPermission}
	got = wrapStaticSyncError(err)

	c.Assert(got, qt.ErrorMatches, `stat public/favicon\.ico: permission denied`)
	c.Assert(errors.Is(got, os.ErrPermission), qt.Equals, true)

	err2 := errors.New("boom")
	c.Assert(wrapStaticSyncError(err2), qt.Equals, err2)
}
