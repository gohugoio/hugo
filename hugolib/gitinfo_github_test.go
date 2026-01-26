// Copyright 2025 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"os"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/modules"
)

func TestParseGitHubURL(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	tests := []struct {
		url           string
		expectedOwner string
		expectedRepo  string
	}{
		{"https://github.com/bep/hugo-testing-git-versions", "bep", "hugo-testing-git-versions"},
		{"https://github.com/gohugoio/hugo", "gohugoio", "hugo"},
		{"https://github.com/owner/repo.git", "owner", "repo"},
		{"https://gitlab.com/owner/repo", "", ""},
		{"https://bitbucket.org/owner/repo", "", ""},
		{"", "", ""},
	}

	for _, test := range tests {
		owner, repo := parseGitHubURL(test.url)
		c.Assert(owner, qt.Equals, test.expectedOwner, qt.Commentf("URL: %s", test.url))
		c.Assert(repo, qt.Equals, test.expectedRepo, qt.Commentf("URL: %s", test.url))
	}
}

func TestIsGitHubModule(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	tests := []struct {
		origin   modules.ModuleOrigin
		expected bool
	}{
		{modules.ModuleOrigin{VCS: "git", URL: "https://github.com/bep/hugo-testing-git-versions"}, true},
		{modules.ModuleOrigin{VCS: "git", URL: "https://gitlab.com/owner/repo"}, false},
		{modules.ModuleOrigin{VCS: "hg", URL: "https://github.com/owner/repo"}, false},
		{modules.ModuleOrigin{}, false},
	}

	for _, test := range tests {
		c.Assert(isGitHubModule(test.origin), qt.Equals, test.expected, qt.Commentf("Origin: %+v", test.origin))
	}
}

func TestSplitCommitMessage(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	tests := []struct {
		message         string
		expectedSubject string
		expectedBody    string
	}{
		{"Simple subject", "Simple subject", ""},
		{"Subject\n\nBody text", "Subject", "Body text"},
		{"Subject\n\nMulti\nline\nbody", "Subject", "Multi\nline\nbody"},
		{"  Trimmed subject  \n\n  Trimmed body  ", "Trimmed subject", "Trimmed body"},
		{"", "", ""},
	}

	for _, test := range tests {
		subject, body := splitCommitMessage(test.message)
		c.Assert(subject, qt.Equals, test.expectedSubject)
		c.Assert(body, qt.Equals, test.expectedBody)
	}
}

func TestAbbreviateHash(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	tests := []struct {
		hash     string
		expected string
	}{
		{"3e0f3930f1ec9a29a7442da5f1bfc0b7e58f167a", "3e0f393"},
		{"abcdefg", "abcdefg"},
		{"abc", "abc"},
		{"", ""},
	}

	for _, test := range tests {
		c.Assert(abbreviateHash(test.hash), qt.Equals, test.expected)
	}
}

func TestGitHubCommitToGitInfo(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	authorDate := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	commitDate := time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC)

	commit := &gitHubCommit{
		SHA: "3e0f3930f1ec9a29a7442da5f1bfc0b7e58f167a",
	}
	commit.Commit.Message = "Test commit message\n\nThis is the body."
	commit.Commit.Author.Name = "Test Author"
	commit.Commit.Author.Email = "test@example.com"
	commit.Commit.Author.Date = authorDate
	commit.Commit.Committer.Date = commitDate

	gi := commit.toGitInfo()

	c.Assert(gi.Hash, qt.Equals, "3e0f3930f1ec9a29a7442da5f1bfc0b7e58f167a")
	c.Assert(gi.AbbreviatedHash, qt.Equals, "3e0f393")
	c.Assert(gi.Subject, qt.Equals, "Test commit message")
	c.Assert(gi.Body, qt.Equals, "This is the body.")
	c.Assert(gi.AuthorName, qt.Equals, "Test Author")
	c.Assert(gi.AuthorEmail, qt.Equals, "test@example.com")
	c.Assert(gi.AuthorDate, qt.Equals, authorDate)
	c.Assert(gi.CommitDate, qt.Equals, commitDate)
}

func TestGitInfoFromGitHubModule(t *testing.T) {
	if !isGitHubAPIAvailable() {
		t.Skip("Skipping test: GitHub API requires HUGO_GITHUB_TOKEN or GITHUB_TOKEN")
	}

	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
enableGitInfo = true

[module]
[[module.imports]]
path = "github.com/bep/hugo-testing-git-versions"
version = "v3.0.1"
[[module.imports.mounts]]
source = "content"
target = "content"
-- layouts/_default/single.html --
Title: {{ .Title }}|
GitInfo: {{ with .GitInfo }}Hash: {{ .Hash }}|Subject: {{ .Subject }}|AuthorName: {{ .AuthorName }}{{ end }}|
-- layouts/_default/list.html --
List: {{ .Title }}
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"Hash: 3e0f3930f1ec9a29a7442da5f1bfc0b7e58f167a|",
		"AuthorName: Bjørn Erik Pedersen|",
	)
}

func isGitHubAPIAvailable() bool {
	return getGitHubToken() != ""
}

func getGitHubToken() string {
	token := os.Getenv(gitHubTokenEnv)
	if token == "" {
		token = os.Getenv(gitHubTokenEnvAlt)
	}
	return token
}
