// Copyright 2018 The Hugo Authors. All rights reserved.
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

package hugo

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestHugoInfo(t *testing.T) {
	c := qt.New(t)

	conf := testConfig{environment: "production", workingDir: "/mywork"}
	hugoInfo := NewInfo(conf, nil)

	c.Assert(hugoInfo.Version(), qt.Equals, CurrentVersion.Version())
	c.Assert(fmt.Sprintf("%T", VersionString("")), qt.Equals, fmt.Sprintf("%T", hugoInfo.Version()))
	c.Assert(hugoInfo.WorkingDir(), qt.Equals, "/mywork")

	bi := getBuildInfo()
	if bi != nil {
		c.Assert(hugoInfo.CommitHash, qt.Equals, bi.Revision)
		c.Assert(hugoInfo.BuildDate, qt.Equals, bi.RevisionTime)
		c.Assert(hugoInfo.GoVersion, qt.Equals, bi.GoVersion)
	}
	c.Assert(hugoInfo.Environment, qt.Equals, "production")
	c.Assert(string(hugoInfo.Generator()), qt.Contains, fmt.Sprintf("Hugo %s", hugoInfo.Version()))
	c.Assert(hugoInfo.IsProduction(), qt.IsTrue)
	c.Assert(hugoInfo.IsExtended(), qt.Equals, IsExtended)

	devHugoInfo := NewInfo(testConfig{environment: "development"}, nil)
	c.Assert(devHugoInfo.IsProduction(), qt.IsFalse)
}

type testConfig struct {
	environment string
	workingDir  string
}

func (c testConfig) Environment() string {
	return c.environment
}

func (c testConfig) WorkingDir() string {
	return c.workingDir
}
