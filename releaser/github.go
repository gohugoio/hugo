package releaser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	gitHubCommitsAPI      = "https://api.github.com/repos/gohugoio/REPO/commits/%s"
	gitHubRepoAPI         = "https://api.github.com/repos/gohugoio/REPO"
	gitHubContributorsAPI = "https://api.github.com/repos/gohugoio/REPO/contributors"
)

type gitHubAPI struct {
	commitsAPITemplate      string
	repoAPI                 string
	contributorsAPITemplate string
}

func newGitHubAPI(repo string) *gitHubAPI {
	return &gitHubAPI{
		commitsAPITemplate:      strings.Replace(gitHubCommitsAPI, "REPO", repo, -1),
		repoAPI:                 strings.Replace(gitHubRepoAPI, "REPO", repo, -1),
		contributorsAPITemplate: strings.Replace(gitHubContributorsAPI, "REPO", repo, -1),
	}
}

type gitHubCommit struct {
	Author  gitHubAuthor `json:"author"`
	HTMLURL string       `json:"html_url"`
}

type gitHubAuthor struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	HTMLURL   string `json:"html_url"`
	AvatarURL string `json:"avatar_url"`
}

type gitHubRepo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	HTMLURL      string `json:"html_url"`
	Stars        int    `json:"stargazers_count"`
	Contributors []gitHubContributor
}

type gitHubContributor struct {
	ID            int    `json:"id"`
	Login         string `json:"login"`
	HTMLURL       string `json:"html_url"`
	Contributions int    `json:"contributions"`
}

func (g *gitHubAPI) fetchCommit(ref string) (gitHubCommit, error) {
	var commit gitHubCommit

	u := fmt.Sprintf(g.commitsAPITemplate, ref)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return commit, err
	}

	err = doGitHubRequest(req, &commit)

	return commit, err
}

func (g *gitHubAPI) fetchRepo() (gitHubRepo, error) {
	var repo gitHubRepo

	req, err := http.NewRequest("GET", g.repoAPI, nil)
	if err != nil {
		return repo, err
	}

	err = doGitHubRequest(req, &repo)
	if err != nil {
		return repo, err
	}

	var contributors []gitHubContributor
	page := 0
	for {
		page++
		var currPage []gitHubContributor
		url := fmt.Sprintf(g.contributorsAPITemplate+"?page=%d", page)

		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return repo, err
		}

		err = doGitHubRequest(req, &currPage)
		if err != nil {
			return repo, err
		}
		if len(currPage) == 0 {
			break
		}

		contributors = append(contributors, currPage...)

	}

	repo.Contributors = contributors

	return repo, err

}

func doGitHubRequest(req *http.Request, v interface{}) error {
	addGitHubToken(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if isError(resp) {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("GitHub lookup failed: %s", string(b))
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func isError(resp *http.Response) bool {
	return resp.StatusCode < 200 || resp.StatusCode > 299
}

func addGitHubToken(req *http.Request) {
	gitHubToken := os.Getenv("GITHUB_TOKEN")
	if gitHubToken != "" {
		req.Header.Add("Authorization", "token "+gitHubToken)
	}
}
