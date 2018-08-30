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
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitHubLookupCommit(t *testing.T) {
	skipIfNoToken(t)
	client := newGitHubAPI("hugo")
	commit, err := client.fetchCommit("793554108763c0984f1a1b1a6ee5744b560d78d0")
	require.NoError(t, err)
	fmt.Println(commit)
}

func TestFetchRepo(t *testing.T) {
	skipIfNoToken(t)
	client := newGitHubAPI("hugo")
	repo, err := client.fetchRepo()
	require.NoError(t, err)
	fmt.Println(">>", len(repo.Contributors))
}

func skipIfNoToken(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("Skip test against GitHub as no GITHUB_TOKEN set.")
	}
}
