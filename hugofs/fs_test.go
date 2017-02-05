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

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewDefault(t *testing.T) {
	v := viper.New()
	f := NewDefault(v)

	assert.NotNil(t, f.Source)
	assert.IsType(t, new(afero.OsFs), f.Source)
	assert.NotNil(t, f.Destination)
	assert.IsType(t, new(afero.OsFs), f.Destination)
	assert.NotNil(t, f.Os)
	assert.IsType(t, new(afero.OsFs), f.Os)
	assert.Nil(t, f.WorkingDir)

	assert.IsType(t, new(afero.OsFs), Os)
}

func TestNewMem(t *testing.T) {
	v := viper.New()
	f := NewMem(v)

	assert.NotNil(t, f.Source)
	assert.IsType(t, new(afero.MemMapFs), f.Source)
	assert.NotNil(t, f.Destination)
	assert.IsType(t, new(afero.MemMapFs), f.Destination)
	assert.IsType(t, new(afero.OsFs), f.Os)
	assert.Nil(t, f.WorkingDir)
}

func TestWorkingDir(t *testing.T) {
	v := viper.New()

	v.Set("workingDir", "/a/b/")

	f := NewMem(v)

	assert.NotNil(t, f.WorkingDir)
	assert.IsType(t, new(afero.BasePathFs), f.WorkingDir)
}
