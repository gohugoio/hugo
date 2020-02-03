// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"testing"

	qt "github.com/frankban/quicktest"
)

// Issue #5662
func TestHugoWithContentDirOverride(t *testing.T) {
	c := qt.New(t)

	hugoCmd := newCommandsBuilder().addAll().build()
	cmd := hugoCmd.getCommand()

	contentDir := "contentOverride"

	cfgStr := `

baseURL = "https://example.org"
title = "Hugo Commands"

contentDir = "thisdoesnotexist"

`
	dir, clean, err := createSimpleTestSite(t, testSiteConfig{configTOML: cfgStr, contentDir: contentDir})
	c.Assert(err, qt.IsNil)
	defer clean()

	cmd.SetArgs([]string{"-s=" + dir, "-c=" + contentDir})

	_, err = cmd.ExecuteC()
	c.Assert(err, qt.IsNil)

}
