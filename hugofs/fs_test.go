// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"testing"

	"github.com/gohugoio/hugo/config"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting/hqt"
	"github.com/spf13/afero"
)

func TestIsOsFs(t *testing.T) {
	c := qt.New(t)

	c.Assert(IsOsFs(Os), qt.Equals, true)
	c.Assert(IsOsFs(&afero.MemMapFs{}), qt.Equals, false)
	c.Assert(IsOsFs(NewBasePathFs(&afero.MemMapFs{}, "/public")), qt.Equals, false)
	c.Assert(IsOsFs(NewBasePathFs(Os, t.TempDir())), qt.Equals, true)
}

func TestNewDefault(t *testing.T) {
	c := qt.New(t)
	v := config.New()
	v.Set("workingDir", t.TempDir())
	v.Set("publishDir", "public")
	f := NewDefault(v)

	c.Assert(f.Source, qt.IsNotNil)
	c.Assert(f.Source, hqt.IsSameType, new(afero.OsFs))
	c.Assert(f.Os, qt.IsNotNil)
	c.Assert(f.WorkingDirReadOnly, qt.IsNotNil)
	c.Assert(IsOsFs(f.WorkingDirReadOnly), qt.IsTrue)
	c.Assert(IsOsFs(f.Source), qt.IsTrue)
	c.Assert(IsOsFs(f.PublishDir), qt.IsTrue)
	c.Assert(IsOsFs(f.Os), qt.IsTrue)
}
