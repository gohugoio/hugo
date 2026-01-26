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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/modules"
	"github.com/gohugoio/hugo/source"
)

const (
	gitHubAPIBase     = "https://api.github.com"
	gitHubURLPrefix   = "https://github.com/"
	gitHubTokenEnv    = "HUGO_GITHUB_TOKEN"
	gitHubTokenEnvAlt = "GITHUB_TOKEN"
)

// gitHubClient fetches git info from GitHub API for Hugo modules.
type gitHubClient struct {
	cache *filecache.Cache
	token string

	// Commit maps keyed by "owner/repo".
	repoCommits map[string]map[string]*gitHubCommit
}

// newGitHubClient creates a new GitHub client.
func newGitHubClient(cache *filecache.Cache) *gitHubClient {
	token := os.Getenv(gitHubTokenEnv)
	if token == "" {
		token = os.Getenv(gitHubTokenEnvAlt)
	}
	return &gitHubClient{
		cache:       cache,
		token:       token,
		repoCommits: make(map[string]map[string]*gitHubCommit),
	}
}

// gitHubCommit represents a commit from the GitHub API.
type gitHubCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
}

// parseGitHubURL parses a GitHub URL and returns owner and repo.
// Returns empty strings if the URL is not a valid GitHub URL.
func parseGitHubURL(url string) (owner, repo string) {
	if !strings.HasPrefix(url, gitHubURLPrefix) {
		return "", ""
	}
	path := strings.TrimPrefix(url, gitHubURLPrefix)
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// isGitHubModule returns true if the module origin is from GitHub.
func isGitHubModule(origin modules.ModuleOrigin) bool {
	return origin.VCS == "git" && strings.HasPrefix(origin.URL, gitHubURLPrefix)
}

// repoKey returns the cache key for a repo's commit map.
func repoKey(owner, repo string) string {
	return owner + "/" + repo
}

// fetchRepoCommits fetches the commit history for a repo starting from sha,
// populating the in-memory map. Results are cached to disk per repo.
// It fetches in pages of 100, stopping when all needed hashes are found.
func (c *gitHubClient) fetchRepoCommits(owner, repo, sha string, needed map[string]bool) error {
	rk := repoKey(owner, repo)

	if c.repoCommits[rk] == nil {
		c.repoCommits[rk] = make(map[string]*gitHubCommit)

		// Try to load from disk cache.
		cacheKey := hashing.XxHashFromStringHexEncoded(owner, repo) + ".json"
		if data, err := c.cache.GetBytes(cacheKey); err == nil && len(data) > 0 {
			var commits []*gitHubCommit
			if err := json.Unmarshal(data, &commits); err == nil {
				for _, commit := range commits {
					c.repoCommits[rk][commit.SHA] = commit
				}
			}
		}
	}

	m := c.repoCommits[rk]

	// Check if all needed hashes are already available.
	if allFound(m, needed) {
		return nil
	}

	// Fetch pages of commits from the API.
	for page := 1; ; page++ {
		url := fmt.Sprintf("%s/repos/%s/%s/commits?sha=%s&per_page=100&page=%d", gitHubAPIBase, owner, repo, sha, page)
		commits, err := c.doGitHubRequest(url)
		if err != nil {
			return err
		}
		if len(commits) == 0 {
			break
		}

		for _, commit := range commits {
			m[commit.SHA] = commit
		}

		if allFound(m, needed) {
			break
		}
	}

	// Persist the full map to disk cache.
	c.persistRepoCache(owner, repo)

	return nil
}

func allFound(m map[string]*gitHubCommit, needed map[string]bool) bool {
	for h := range needed {
		if _, ok := m[h]; !ok {
			return false
		}
	}
	return true
}

func (c *gitHubClient) persistRepoCache(owner, repo string) {
	rk := repoKey(owner, repo)
	m := c.repoCommits[rk]
	if len(m) == 0 {
		return
	}

	commits := make([]*gitHubCommit, 0, len(m))
	for _, commit := range m {
		commits = append(commits, commit)
	}

	data, err := json.Marshal(commits)
	if err != nil {
		return
	}

	cacheKey := hashing.XxHashFromStringHexEncoded(owner, repo) + ".json"
	_ = c.cache.SetBytes(cacheKey, data)
}

func (c *gitHubClient) doGitHubRequest(url string) ([]*gitHubCommit, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var commits []*gitHubCommit
	if err := json.Unmarshal(data, &commits); err != nil {
		return nil, err
	}

	return commits, nil
}

// getCommit returns the commit for a given hash, using the cached repo history.
func (c *gitHubClient) getCommit(owner, repo, hash string) *gitHubCommit {
	rk := repoKey(owner, repo)
	if m, ok := c.repoCommits[rk]; ok {
		return m[hash]
	}
	return nil
}

// toGitInfo converts a GitHub commit to a GitInfo struct.
func (c *gitHubCommit) toGitInfo() *source.GitInfo {
	subject, body := splitCommitMessage(c.Commit.Message)
	return &source.GitInfo{
		Hash:            c.SHA,
		AbbreviatedHash: abbreviateHash(c.SHA),
		Subject:         subject,
		AuthorName:      c.Commit.Author.Name,
		AuthorEmail:     c.Commit.Author.Email,
		AuthorDate:      c.Commit.Author.Date,
		CommitDate:      c.Commit.Committer.Date,
		Body:            body,
	}
}

// splitCommitMessage splits a commit message into subject and body.
func splitCommitMessage(message string) (subject, body string) {
	message = strings.TrimSpace(message)
	parts := strings.SplitN(message, "\n", 2)
	subject = strings.TrimSpace(parts[0])
	if len(parts) > 1 {
		body = strings.TrimSpace(parts[1])
	}
	return
}

// abbreviateHash returns the first 7 characters of a hash.
func abbreviateHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}
