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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
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
	gitHubGraphQLAPI  = "https://api.github.com/graphql"
	gitHubURLPrefix   = "https://github.com/"
	gitHubTokenEnv    = "HUGO_GITHUB_TOKEN"
	gitHubTokenEnvAlt = "GITHUB_TOKEN"

	// Maximum number of files to query in a single GraphQL request.
	graphQLBatchSize = 50
)

// gitHubClient fetches git info from GitHub API for Hugo modules.
type gitHubClient struct {
	cache *filecache.Cache
	token string

	// File commit cache keyed by "owner/repo/ref/path".
	fileCommits map[string]*gitHubCommit
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
		fileCommits: make(map[string]*gitHubCommit),
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
func isGitHubModule(m modules.Module) bool {
	origin := m.Origin()
	return origin.VCS == "git" && strings.HasPrefix(origin.URL, gitHubURLPrefix)
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

// fileCommitKey returns a cache key for a file commit.
func fileCommitKey(owner, repo, ref, path string) string {
	return owner + "/" + repo + "/" + ref + "/" + path
}

// fetchFileCommit fetches the commit info for a single file.
// It first checks the cache, then uses the batch fetcher if needed.
func (c *gitHubClient) fetchFileCommit(owner, repo, ref, path string) (*gitHubCommit, error) {
	key := fileCommitKey(owner, repo, ref, path)

	// Check in-memory cache.
	if commit, ok := c.fileCommits[key]; ok {
		return commit, nil
	}

	// Check disk cache.
	cacheKey := hashing.XxHashFromStringHexEncoded(key) + ".json"
	if data, err := c.cache.GetBytes(cacheKey); err == nil && len(data) > 0 {
		var commit gitHubCommit
		if err := json.Unmarshal(data, &commit); err == nil {
			c.fileCommits[key] = &commit
			return &commit, nil
		}
	}

	// Fetch from API (single file).
	commits, err := c.fetchFileCommitsGraphQL(owner, repo, ref, []string{path})
	if err != nil {
		return nil, err
	}

	commit := commits[path]
	if commit != nil {
		c.fileCommits[key] = commit
		c.persistFileCommit(owner, repo, ref, path, commit)
	}

	return commit, nil
}

// fetchFileCommitsGraphQL fetches commit info for multiple files using GraphQL.
func (c *gitHubClient) fetchFileCommitsGraphQL(owner, repo, ref string, paths []string) (map[string]*gitHubCommit, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	result := make(map[string]*gitHubCommit)

	// Process in batches.
	for i := 0; i < len(paths); i += graphQLBatchSize {
		batch := paths[i:min(i+graphQLBatchSize, len(paths))]

		commits, err := c.doGraphQLBatchRequest(owner, repo, ref, batch)
		if err != nil {
			return nil, err
		}

		maps.Copy(result, commits)
	}

	return result, nil
}

// buildGraphQLQuery builds a GraphQL query for fetching file history.
func buildGraphQLQuery(owner, repo, ref string, paths []string) string {
	var b strings.Builder
	b.WriteString("{\n  repository(owner: \"")
	b.WriteString(owner)
	b.WriteString("\", name: \"")
	b.WriteString(repo)
	b.WriteString("\") {\n    ref(qualifiedName: \"")
	b.WriteString(ref)
	b.WriteString("\") {\n      target {\n        ... on Commit {\n")

	for i, path := range paths {
		// Use f0, f1, f2... as aliases.
		fmt.Fprintf(&b, "          f%d: history(path: \"%s\", first: 1) {\n", i, escapeGraphQLString(path))
		b.WriteString("            nodes {\n")
		b.WriteString("              oid\n")
		b.WriteString("              message\n")
		b.WriteString("              author {\n")
		b.WriteString("                name\n")
		b.WriteString("                email\n")
		b.WriteString("                date\n")
		b.WriteString("              }\n")
		b.WriteString("              committer {\n")
		b.WriteString("                date\n")
		b.WriteString("              }\n")
		b.WriteString("            }\n")
		b.WriteString("          }\n")
	}

	b.WriteString("        }\n      }\n    }\n  }\n}")
	return b.String()
}

// escapeGraphQLString escapes special characters in a GraphQL string.
func escapeGraphQLString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// graphQLResponse represents the response from GitHub's GraphQL API.
type graphQLResponse struct {
	Data struct {
		Repository struct {
			Ref struct {
				Target map[string]json.RawMessage `json:"target"`
			} `json:"ref"`
		} `json:"repository"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// graphQLHistoryNode represents a commit node in the GraphQL response.
type graphQLHistoryNode struct {
	Nodes []struct {
		OID     string `json:"oid"`
		Message string `json:"message"`
		Author  struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Date time.Time `json:"date"`
		} `json:"committer"`
	} `json:"nodes"`
}

// doGraphQLBatchRequest executes a GraphQL batch request for file commits.
func (c *gitHubClient) doGraphQLBatchRequest(owner, repo, ref string, paths []string) (map[string]*gitHubCommit, error) {
	query := buildGraphQLQuery(owner, repo, ref, paths)

	reqBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", gitHubGraphQLAPI, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub GraphQL API returned %d: %s", resp.StatusCode, string(body))
	}

	var gqlResp graphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, fmt.Errorf("failed to parse GraphQL response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", gqlResp.Errors[0].Message)
	}

	result := make(map[string]*gitHubCommit)
	target := gqlResp.Data.Repository.Ref.Target

	for i, path := range paths {
		alias := fmt.Sprintf("f%d", i)
		historyData, ok := target[alias]
		if !ok {
			continue
		}

		var history graphQLHistoryNode
		if err := json.Unmarshal(historyData, &history); err != nil {
			continue
		}

		if len(history.Nodes) == 0 {
			continue
		}

		node := history.Nodes[0]
		commit := &gitHubCommit{
			SHA: node.OID,
		}
		commit.Commit.Message = node.Message
		commit.Commit.Author.Name = node.Author.Name
		commit.Commit.Author.Email = node.Author.Email
		commit.Commit.Author.Date = node.Author.Date
		commit.Commit.Committer.Date = node.Committer.Date

		result[path] = commit
	}

	return result, nil
}

// persistFileCommit saves a file commit to disk cache.
func (c *gitHubClient) persistFileCommit(owner, repo, ref, path string, commit *gitHubCommit) {
	key := fileCommitKey(owner, repo, ref, path)
	cacheKey := hashing.XxHashFromStringHexEncoded(key) + ".json"

	data, err := json.Marshal(commit)
	if err != nil {
		return
	}

	_ = c.cache.SetBytes(cacheKey, data)
}
