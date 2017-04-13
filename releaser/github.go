package releaser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	gitHubCommitsApi      = "https://api.github.com/repos/spf13/hugo/commits/%s"
	gitHubRepoApi         = "https://api.github.com/repos/spf13/hugo"
	gitHubContributorsApi = "https://api.github.com/repos/spf13/hugo/contributors"
)

type gitHubCommit struct {
	Author  gitHubAuthor `json:"author"`
	HtmlURL string       `json:"html_url"`
}

type gitHubAuthor struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	HtmlURL   string `json:"html_url"`
	AvatarURL string `json:"avatar_url"`
}

type gitHubRepo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	HtmlURL      string `json:"html_url"`
	Stars        int    `json:"stargazers_count"`
	Contributors []gitHubContributor
}

type gitHubContributor struct {
	ID            int    `json:"id"`
	Login         string `json:"login"`
	HtmlURL       string `json:"html_url"`
	Contributions int    `json:"contributions"`
}

func fetchCommit(ref string) (gitHubCommit, error) {
	var commit gitHubCommit

	u := fmt.Sprintf(gitHubCommitsApi, ref)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return commit, err
	}

	err = doGitHubRequest(req, &commit)

	return commit, err
}

func fetchRepo() (gitHubRepo, error) {
	var repo gitHubRepo

	req, err := http.NewRequest("GET", gitHubRepoApi, nil)
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
		url := fmt.Sprintf(gitHubContributorsApi+"?page=%d", page)

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
