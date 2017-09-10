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

	"github.com/stretchr/testify/require"
)

func TestGitInfos(t *testing.T) {
	skipIfCI(t)
	infos, err := getGitInfos("v0.20", "hugo", "", false)

	require.NoError(t, err)
	require.True(t, len(infos) > 0)

}

func TestIssuesRe(t *testing.T) {

	body := `
This is a commit message.

Updates #123
Fix #345
closes #543
See #456
	`

	issues := extractIssues(body)

	require.Len(t, issues, 4)
	require.Equal(t, 123, issues[0])
	require.Equal(t, 543, issues[2])

}

func TestGitVersionTagBefore(t *testing.T) {
	skipIfCI(t)
	v1, err := gitVersionTagBefore("v0.18")
	require.NoError(t, err)
	require.Equal(t, "v0.17", v1)
}

func TestTagExists(t *testing.T) {
	skipIfCI(t)
	b1, err := tagExists("v0.18")
	require.NoError(t, err)
	require.True(t, b1)

	b2, err := tagExists("adfagdsfg")
	require.NoError(t, err)
	require.False(t, b2)

}

func skipIfCI(t *testing.T) {
	if isCI() {
		// Travis has an ancient git with no --invert-grep: https://github.com/travis-ci/travis-ci/issues/6328
		// Also Travis clones very shallowly, making some of the tests above shaky.
		t.Skip("Skip git test on Linux to make Travis happy.")
	}
}
