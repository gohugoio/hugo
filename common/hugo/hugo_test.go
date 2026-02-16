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
	"context"
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
)

func TestDeprecationLogLevelFromVersion(t *testing.T) {
	c := qt.New(t)

	c.Assert(deprecationLogLevelFromVersion("0.55.0"), qt.Equals, logg.LevelError)
	ver := CurrentVersion
	c.Assert(deprecationLogLevelFromVersion(ver.String()), qt.Equals, logg.LevelInfo)
	ver.Minor -= 3
	c.Assert(deprecationLogLevelFromVersion(ver.String()), qt.Equals, logg.LevelWarn)
	ver.Minor -= 4
	c.Assert(deprecationLogLevelFromVersion(ver.String()), qt.Equals, logg.LevelWarn)
	ver.Minor -= 13
	c.Assert(deprecationLogLevelFromVersion(ver.String()), qt.Equals, logg.LevelError)

	// Added just to find the threshold for where we can remove deprecated items.
	// Subtract 5 from the minor version of the first ERRORed version => 0.136.0.
	c.Assert(deprecationLogLevelFromVersion("0.141.0"), qt.Equals, logg.LevelError)
}

func TestMarkupScope(t *testing.T) {
	c := qt.New(t)

	ctx := context.Background()
	ctx = SetMarkupScope(ctx, "foo")

	var hugoCtx Context
	c.Assert(hugoCtx.MarkupScope(ctx), qt.Equals, "foo")
	c.Assert(GetMarkupScope(ctx), qt.Equals, "foo")
}

func TestGetBuildInfo(t *testing.T) {
	c := qt.New(t)

	bi := GetBuildInfo()
	// In test mode, build info may or may not be available.
	if bi != nil {
		c.Assert(bi.GoVersion, qt.Not(qt.Equals), "")
	}
}
