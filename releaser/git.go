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
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var issueRe = regexp.MustCompile(`(?i)[Updates?|Closes?|Fix.*|See] #(\d+)`)

const (
	notesChanges    = "notesChanges"
	templateChanges = "templateChanges"
	coreChanges     = "coreChanges"
	outChanges      = "outChanges"
	docsChanges     = "docsChanges"
	otherChanges    = "otherChanges"
)

type changeLog struct {
	Version      string
	Enhancements map[string]gitInfos
	Fixes        map[string]gitInfos
	Notes        gitInfos
	All          gitInfos

	// Overall stats
	Repo             *gitHubRepo
	ContributorCount int
	ThemeCount       int
}

func newChangeLog(infos gitInfos) *changeLog {
	return &changeLog{
		Enhancements: make(map[string]gitInfos),
		Fixes:        make(map[string]gitInfos),
		All:          infos,
	}
}

func (l *changeLog) addGitInfo(isFix bool, info gitInfo, category string) {
	var (
		infos   gitInfos
		found   bool
		segment map[string]gitInfos
	)

	if category == notesChanges {
		l.Notes = append(l.Notes, info)
		return
	} else if isFix {
		segment = l.Fixes
	} else {
		segment = l.Enhancements
	}

	infos, found = segment[category]
	if !found {
		infos = gitInfos{}
	}

	infos = append(infos, info)
	segment[category] = infos
}

func gitInfosToChangeLog(infos gitInfos) *changeLog {
	log := newChangeLog(infos)
	for _, info := range infos {
		los := strings.ToLower(info.Subject)
		isFix := strings.Contains(los, "fix")
		var category = otherChanges

		// TODO(bep) improve
		if regexp.MustCompile("(?i)deprecate").MatchString(los) {
			category = notesChanges
		} else if regexp.MustCompile("(?i)tpl|tplimpl:|layout").MatchString(los) {
			category = templateChanges
		} else if regexp.MustCompile("(?i)docs?:|documentation:").MatchString(los) {
			category = docsChanges
		} else if regexp.MustCompile("(?i)hugolib:").MatchString(los) {
			category = coreChanges
		} else if regexp.MustCompile("(?i)out(put)?:|media:|Output|Media").MatchString(los) {
			category = outChanges
		}

		// Trim package prefix.
		colonIdx := strings.Index(info.Subject, ":")
		if colonIdx != -1 && colonIdx < (len(info.Subject)/2) {
			info.Subject = info.Subject[colonIdx+1:]
		}

		info.Subject = strings.TrimSpace(info.Subject)

		log.addGitInfo(isFix, info, category)
	}

	return log
}

type gitInfo struct {
	Hash    string
	Author  string
	Subject string
	Body    string

	GitHubCommit *gitHubCommit
}

func (g gitInfo) Issues() []int {
	return extractIssues(g.Body)
}

func (g gitInfo) AuthorID() string {
	if g.GitHubCommit != nil {
		return g.GitHubCommit.Author.Login
	}
	return g.Author
}

func extractIssues(body string) []int {
	var i []int
	m := issueRe.FindAllStringSubmatch(body, -1)
	for _, mm := range m {
		issueID, err := strconv.Atoi(mm[1])
		if err != nil {
			continue
		}
		i = append(i, issueID)
	}
	return i
}

type gitInfos []gitInfo

func git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git failed: %q: %q", err, out)
	}
	return string(out), nil
}

func getGitInfos(tag string, remote bool) (gitInfos, error) {
	return getGitInfosBefore("HEAD", tag, remote)
}

type countribCount struct {
	Author       string
	GitHubAuthor gitHubAuthor
	Count        int
}

func (c countribCount) AuthorLink() string {
	if c.GitHubAuthor.HtmlURL != "" {
		return fmt.Sprintf("[@%s](%s)", c.GitHubAuthor.Login, c.GitHubAuthor.HtmlURL)
	}

	if !strings.Contains(c.Author, "@") {
		return c.Author
	}

	return c.Author[:strings.Index(c.Author, "@")]

}

type contribCounts []countribCount

func (c contribCounts) Less(i, j int) bool { return c[i].Count > c[j].Count }
func (c contribCounts) Len() int           { return len(c) }
func (c contribCounts) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (g gitInfos) ContribCountPerAuthor() contribCounts {
	var c contribCounts

	counters := make(map[string]countribCount)

	for _, gi := range g {
		authorID := gi.AuthorID()
		if count, ok := counters[authorID]; ok {
			count.Count = count.Count + 1
			counters[authorID] = count
		} else {
			var ghA gitHubAuthor
			if gi.GitHubCommit != nil {
				ghA = gi.GitHubCommit.Author
			}
			authorCount := countribCount{Count: 1, Author: gi.Author, GitHubAuthor: ghA}
			counters[authorID] = authorCount
		}
	}

	for _, v := range counters {
		c = append(c, v)
	}

	sort.Sort(c)
	return c
}

func getGitInfosBefore(ref, tag string, remote bool) (gitInfos, error) {

	var g gitInfos

	log, err := gitLogBefore(ref, tag)
	if err != nil {
		return g, err
	}

	log = strings.Trim(log, "\n\x1e'")
	entries := strings.Split(log, "\x1e")

	for _, entry := range entries {
		items := strings.Split(entry, "\x1f")
		gi := gitInfo{
			Hash:    items[0],
			Author:  items[1],
			Subject: items[2],
			Body:    items[3],
		}
		if remote {
			gc, err := fetchCommit(gi.Hash)
			if err == nil {
				gi.GitHubCommit = &gc
			}
		}
		g = append(g, gi)
	}

	return g, nil
}

// Ignore autogenerated commits etc. in change log. This is a regexp.
const ignoredCommits = "release:|vendor:|snapcraft:"

func gitLogBefore(ref, tag string) (string, error) {
	var prevTag string
	var err error
	if tag != "" {
		prevTag = tag
	} else {
		prevTag, err = gitVersionTagBefore(ref)
		if err != nil {
			return "", err
		}
	}
	log, err := git("log", "-E", fmt.Sprintf("--grep=%s", ignoredCommits), "--invert-grep", "--pretty=format:%x1e%h%x1f%aE%x1f%s%x1f%b", "--abbrev-commit", prevTag+".."+ref)
	if err != nil {
		return ",", err
	}

	return log, err
}

func gitVersionTagBefore(ref string) (string, error) {
	return gitShort("describe", "--tags", "--abbrev=0", "--always", "--match", "v[0-9]*", ref+"^")
}

func gitLog() (string, error) {
	return gitLogBefore("HEAD", "")
}

func gitShort(args ...string) (output string, err error) {
	output, err = git(args...)
	return strings.Replace(strings.Split(output, "\n")[0], "'", "", -1), err
}

func tagExists(tag string) (bool, error) {
	out, err := git("tag", "-l", tag)

	if err != nil {
		return false, err
	}

	if strings.Contains(out, tag) {
		return true, nil
	}

	return false, nil
}
