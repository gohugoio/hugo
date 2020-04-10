// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestHugoVersion(t *testing.T) {
	c := qt.New(t)

	c.Assert(version(0.15, 0, "-DEV"), qt.Equals, "0.15-DEV")
	c.Assert(version(0.15, 2, "-DEV"), qt.Equals, "0.15.2-DEV")

	v := Version{Number: 0.21, PatchLevel: 0, Suffix: "-DEV"}

	c.Assert(v.ReleaseVersion().String(), qt.Equals, "0.21")
	c.Assert(v.String(), qt.Equals, "0.21-DEV")
	c.Assert(v.Next().String(), qt.Equals, "0.22")
	nextVersionString := v.Next().Version()
	c.Assert(nextVersionString.String(), qt.Equals, "0.22")
	c.Assert(nextVersionString.Eq("0.22"), qt.Equals, true)
	c.Assert(nextVersionString.Eq("0.21"), qt.Equals, false)
	c.Assert(nextVersionString.Eq(nextVersionString), qt.Equals, true)
	c.Assert(v.NextPatchLevel(3).String(), qt.Equals, "0.20.3")

	// We started to use full semver versions even for main
	// releases in v0.54.0
	v = Version{Number: 0.53, PatchLevel: 0}
	c.Assert(v.String(), qt.Equals, "0.53")
	c.Assert(v.Next().String(), qt.Equals, "0.54.0")
	c.Assert(v.Next().Next().String(), qt.Equals, "0.55.0")
	v = Version{Number: 0.54, PatchLevel: 0, Suffix: "-DEV"}
	c.Assert(v.String(), qt.Equals, "0.54.0-DEV")
}

func TestCompareVersions(t *testing.T) {
	c := qt.New(t)

	c.Assert(compareVersions(0.20, 0, 0.20), qt.Equals, 0)
	c.Assert(compareVersions(0.20, 0, float32(0.20)), qt.Equals, 0)
	c.Assert(compareVersions(0.20, 0, float64(0.20)), qt.Equals, 0)
	c.Assert(compareVersions(0.19, 1, 0.20), qt.Equals, 1)
	c.Assert(compareVersions(0.19, 3, "0.20.2"), qt.Equals, 1)
	c.Assert(compareVersions(0.19, 1, 0.01), qt.Equals, -1)
	c.Assert(compareVersions(0, 1, 3), qt.Equals, 1)
	c.Assert(compareVersions(0, 1, int32(3)), qt.Equals, 1)
	c.Assert(compareVersions(0, 1, int64(3)), qt.Equals, 1)
	c.Assert(compareVersions(0.20, 0, "0.20"), qt.Equals, 0)
	c.Assert(compareVersions(0.20, 1, "0.20.1"), qt.Equals, 0)
	c.Assert(compareVersions(0.20, 1, "0.20"), qt.Equals, -1)
	c.Assert(compareVersions(0.20, 0, "0.20.1"), qt.Equals, 1)
	c.Assert(compareVersions(0.20, 1, "0.20.2"), qt.Equals, 1)
	c.Assert(compareVersions(0.21, 1, "0.22.1"), qt.Equals, 1)
	c.Assert(compareVersions(0.22, 0, "0.22-DEV"), qt.Equals, -1)
	c.Assert(compareVersions(0.22, 0, "0.22.1-DEV"), qt.Equals, 1)
	c.Assert(compareVersionsWithSuffix(0.22, 0, "-DEV", "0.22"), qt.Equals, 1)
	c.Assert(compareVersionsWithSuffix(0.22, 1, "-DEV", "0.22"), qt.Equals, -1)
	c.Assert(compareVersionsWithSuffix(0.22, 1, "-DEV", "0.22.1-DEV"), qt.Equals, 0)
}

func TestParseHugoVersion(t *testing.T) {
	c := qt.New(t)

	c.Assert(MustParseVersion("0.25").String(), qt.Equals, "0.25")
	c.Assert(MustParseVersion("0.25.2").String(), qt.Equals, "0.25.2")
	c.Assert(MustParseVersion("0.25-test").String(), qt.Equals, "0.25-test")
	c.Assert(MustParseVersion("0.25-DEV").String(), qt.Equals, "0.25-DEV")
}

func TestGoMinorVersion(t *testing.T) {
	c := qt.New(t)
	c.Assert(goMinorVersion("go1.12.5"), qt.Equals, 12)
	c.Assert(goMinorVersion("go1.14rc1"), qt.Equals, 14)
	c.Assert(GoMinorVersion() >= 11, qt.Equals, true)
}
