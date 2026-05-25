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

package modules

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

func TestPathKey(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		in     string
		expect string
	}{
		{"github.com/foo", "github.com/foo"},
		{"github.com/foo/v2", "github.com/foo"},
		{"github.com/foo/v12", "github.com/foo"},
		{"github.com/foo/v3d", "github.com/foo/v3d"},
		{"MyTheme", "mytheme"},
	} {
		c.Assert(pathBase(test.in), qt.Equals, test.expect)
	}
}

func TestResolveWorkspacePattern(t *testing.T) {
	c := qt.New(t)

	fs := afero.NewMemMapFs()
	root := filepath.FromSlash("/project/mod")

	// Create workspace package.json files.
	for _, dir := range []string{
		"packages/aws1",
		"packages/aws2",
		"packages/hugoautogen",
		"other/foo",
	} {
		p := filepath.Join(root, dir, "package.json")
		c.Assert(afero.WriteFile(fs, p, []byte(`{}`), 0o644), qt.IsNil)
	}

	// Literal path, no glob.
	c.Assert(ResolveWorkspacePattern(fs, root, "packages/aws1"), qt.DeepEquals, []string{"packages/aws1"})

	got := ResolveWorkspacePattern(fs, root, "packages/*")
	c.Assert(got, qt.DeepEquals, []string{
		filepath.FromSlash("packages/aws1"),
		filepath.FromSlash("packages/aws2"),
		filepath.FromSlash("packages/hugoautogen"),
	})

	// We currently support only one level of globbing, so this should give the same result as above.
	got = ResolveWorkspacePattern(fs, root, "packages/**")
	c.Assert(got, qt.DeepEquals, []string{
		filepath.FromSlash("packages/aws1"),
		filepath.FromSlash("packages/aws2"),
		filepath.FromSlash("packages/hugoautogen"),
	})

	got = ResolveWorkspacePattern(fs, root, "packages/{aws1,aws2}")
	c.Assert(got, qt.DeepEquals, []string{
		filepath.FromSlash("packages/aws1"),
		filepath.FromSlash("packages/aws2"),
	})

	// ** matches recursively.
	nestedDir := filepath.Join(root, "packages", "aws1", "sub", "package.json")
	c.Assert(afero.WriteFile(fs, nestedDir, []byte(`{}`), 0o644), qt.IsNil)
	got = ResolveWorkspacePattern(fs, root, "packages/**")
	c.Assert(got, qt.DeepEquals, []string{
		filepath.FromSlash("packages/aws1"),
		filepath.FromSlash("packages/aws1/sub"),
		filepath.FromSlash("packages/aws2"),
		filepath.FromSlash("packages/hugoautogen"),
	})

	// No matches.
	got = ResolveWorkspacePattern(fs, root, "nope/*")
	c.Assert(got, qt.HasLen, 0)
}

func TestFilterUnwantedMounts(t *testing.T) {
	mounts := []Mount{
		{Source: "a", Target: "b", Lang: "en"},
		{Source: "a", Target: "b", Lang: "en"},
		{Source: "b", Target: "c", Lang: "en"},
	}

	filtered := filterDuplicateMounts(mounts)

	c := qt.New(t)
	c.Assert(len(filtered), qt.Equals, 2)
	c.Assert(filtered, qt.DeepEquals, []Mount{{Source: "a", Target: "b", Lang: "en"}, {Source: "b", Target: "c", Lang: "en"}})
}
