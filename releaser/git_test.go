// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package releaser

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGitInfos(t *testing.T) {
	c := qt.New(t)
	skipIfCI(t)
	infos, err := getGitInfos("v0.20", "hugo", "", false)

	c.Assert(err, qt.IsNil)
	c.Assert(len(infos) > 0, qt.Equals, true)
}

func TestIssuesRe(t *testing.T) {
	c := qt.New(t)

	body := `
This is a commit message.

Updates #123
Fix #345
closes #543
See #456
	`

	issues := extractIssues(body)

	c.Assert(len(issues), qt.Equals, 4)
	c.Assert(issues[0], qt.Equals, 123)
	c.Assert(issues[2], qt.Equals, 543)

}

func TestGitVersionTagBefore(t *testing.T) {
	skipIfCI(t)
	c := qt.New(t)
	v1, err := gitVersionTagBefore("v0.18")
	c.Assert(err, qt.IsNil)
	c.Assert(v1, qt.Equals, "v0.17")
}

func TestTagExists(t *testing.T) {
	skipIfCI(t)
	c := qt.New(t)
	b1, err := tagExists("v0.18")
	c.Assert(err, qt.IsNil)
	c.Assert(b1, qt.Equals, true)

	b2, err := tagExists("adfagdsfg")
	c.Assert(err, qt.IsNil)
	c.Assert(b2, qt.Equals, false)

}

func skipIfCI(t *testing.T) {
	if isCI() {
		// Travis has an ancient git with no --invert-grep: https://github.com/travis-ci/travis-ci/issues/6328
		// Also Travis clones very shallowly, making some of the tests above shaky.
		t.Skip("Skip git test on Linux to make Travis happy.")
	}
}
