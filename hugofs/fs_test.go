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

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting/hqt"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func TestNewDefault(t *testing.T) {
	c := qt.New(t)
	v := viper.New()
	f := NewDefault(v)

	c.Assert(f.Source, qt.Not(qt.IsNil))
	c.Assert(f.Source, hqt.IsSameType, new(afero.OsFs))
	c.Assert(f.Os, qt.Not(qt.IsNil))
	c.Assert(f.WorkingDir, qt.IsNil)

}

func TestNewMem(t *testing.T) {
	c := qt.New(t)
	v := viper.New()
	f := NewMem(v)

	c.Assert(f.Source, qt.Not(qt.IsNil))
	c.Assert(f.Source, hqt.IsSameType, new(afero.MemMapFs))
	c.Assert(f.Destination, qt.Not(qt.IsNil))
	c.Assert(f.Destination, hqt.IsSameType, new(afero.MemMapFs))
	c.Assert(f.Os, hqt.IsSameType, new(afero.OsFs))
	c.Assert(f.WorkingDir, qt.IsNil)
}

func TestWorkingDir(t *testing.T) {
	c := qt.New(t)
	v := viper.New()

	v.Set("workingDir", "/a/b/")

	f := NewMem(v)

	c.Assert(f.WorkingDir, qt.Not(qt.IsNil))
	c.Assert(f.WorkingDir, hqt.IsSameType, new(afero.BasePathFs))

}
