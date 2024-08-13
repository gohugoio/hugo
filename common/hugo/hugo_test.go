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
	"fmt"
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
)

func TestHugoInfo(t *testing.T) {
	c := qt.New(t)

	conf := testConfig{environment: "production", workingDir: "/mywork", running: false}
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
	c.Assert(hugoInfo.IsDevelopment(), qt.Equals, false)
	c.Assert(hugoInfo.IsProduction(), qt.Equals, true)
	c.Assert(hugoInfo.IsExtended(), qt.Equals, IsExtended)
	c.Assert(hugoInfo.IsServer(), qt.Equals, false)

	devHugoInfo := NewInfo(testConfig{environment: "development", running: true}, nil)
	c.Assert(devHugoInfo.IsDevelopment(), qt.Equals, true)
	c.Assert(devHugoInfo.IsProduction(), qt.Equals, false)
	c.Assert(devHugoInfo.IsServer(), qt.Equals, true)
}

func TestDeprecationLogLevelFromVersion(t *testing.T) {
	c := qt.New(t)

	c.Assert(deprecationLogLevelFromVersion("0.55.0"), qt.Equals, logg.LevelError)
	ver := CurrentVersion
	c.Assert(deprecationLogLevelFromVersion(ver.String()), qt.Equals, logg.LevelInfo)
	ver.Minor -= 1
	c.Assert(deprecationLogLevelFromVersion(ver.String()), qt.Equals, logg.LevelInfo)
	ver.Minor -= 6
	c.Assert(deprecationLogLevelFromVersion(ver.String()), qt.Equals, logg.LevelWarn)
	ver.Minor -= 6
	c.Assert(deprecationLogLevelFromVersion(ver.String()), qt.Equals, logg.LevelError)
}

func TestMarkupScope(t *testing.T) {
	c := qt.New(t)

	conf := testConfig{environment: "production", workingDir: "/mywork", running: false}
	info := NewInfo(conf, nil)

	ctx := context.Background()

	ctx = SetMarkupScope(ctx, "foo")

	c.Assert(info.Context.MarkupScope(ctx), qt.Equals, "foo")
}

type testConfig struct {
	environment  string
	running      bool
	workingDir   string
	multihost    bool
	multilingual bool
}

func (c testConfig) Environment() string {
	return c.environment
}

func (c testConfig) Running() bool {
	return c.running
}

func (c testConfig) WorkingDir() string {
	return c.workingDir
}

func (c testConfig) IsMultihost() bool {
	return c.multihost
}

func (c testConfig) IsMultilingual() bool {
	return c.multilingual
}
