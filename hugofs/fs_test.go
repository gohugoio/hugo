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
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitDefault(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	InitDefaultFs()

	assert.IsType(t, new(afero.OsFs), Source())
	assert.IsType(t, new(afero.OsFs), Destination())
	assert.IsType(t, new(afero.OsFs), Os())
	assert.Nil(t, WorkingDir())
}

func TestInitMemFs(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	InitMemFs()

	assert.IsType(t, new(afero.MemMapFs), Source())
	assert.IsType(t, new(afero.MemMapFs), Destination())
	assert.IsType(t, new(afero.OsFs), Os())
	assert.Nil(t, WorkingDir())
}

func TestSetSource(t *testing.T) {

	InitMemFs()

	SetSource(new(afero.OsFs))
	assert.IsType(t, new(afero.OsFs), Source())
}

func TestSetDestination(t *testing.T) {

	InitMemFs()

	SetDestination(new(afero.OsFs))
	assert.IsType(t, new(afero.OsFs), Destination())
}

func TestWorkingDir(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("WorkingDir", "/a/b/")

	InitMemFs()

	assert.IsType(t, new(afero.BasePathFs), WorkingDir())
}
