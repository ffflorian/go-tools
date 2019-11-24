/*
Copyright © 2019 Florian Keller <github@floriankeller.de>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package githubclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ffflorian/go-tools/gh-open/simplelogger"
)

// GitHubClient is a configuration struct for the client
type GitHubClient struct {
	Logger  simplelogger.SimpleLogger
	Timeout time.Duration
}

// PullRequest represents a pull request on GitHub
type PullRequest struct {
	Head struct {
		Ref string `json:"ref"` // The branch reference
	} `json:"head"`
	Links struct {
		HTML struct {
			Href string `json:"href"` // The pull request URL
		} `json:"html"`
	} `json:"_links"`
}

const baseURL = "https://api.github.com"

// New returns a new instance of GitHubClient
func New(logger simplelogger.SimpleLogger, timeout time.Duration) GitHubClient {
	client := GitHubClient{
		Logger:  logger,
		Timeout: timeout,
	}
	return client
}

func (githubClient GitHubClient) request(urlPath string) (*[]byte, error) {
	var httpClient = &http.Client{Timeout: githubClient.Timeout}
	var fullURL = fmt.Sprintf("%s/%s", baseURL, urlPath)

	githubClient.Logger.Logf("Sending GET request to \"%s\" with timeout \"%s\" ...", fullURL, githubClient.Timeout)

	response, responseError := httpClient.Get(fullURL)
	if responseError != nil {
		return nil, responseError
	}

	defer response.Body.Close()

	githubClient.Logger.Logf("Got response status code \"%d\"", response.StatusCode)

	buffer, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return nil, readError
	}

	return &buffer, nil
}

// GetPullRequests gets pull requests from GitHub,
// see https://developer.github.com/v3/pulls/#list-pull-requests
func (githubClient GitHubClient) GetPullRequests(repoUser string, repoName string) (*[]PullRequest, error) {
	var pullRequests *[]PullRequest

	fullURL := fmt.Sprintf("repos/%s/%s/pulls", repoUser, repoName)
	requestBuffer, requestError := githubClient.request(fullURL)
	if requestError != nil {
		return nil, requestError
	}

	unmarshalError := json.Unmarshal(*requestBuffer, &pullRequests)
	if unmarshalError != nil {
		return nil, unmarshalError
	}

	githubClient.Logger.Log("Got pull requests", *pullRequests)

	return pullRequests, nil
}

// GetPullRequestByBranch returns a pull request URL for the specified branch if it exists
func (githubClient GitHubClient) GetPullRequestByBranch(repoUser string, repoName string, branch string) (string, error) {
	pullRequests, pullRequestError := githubClient.GetPullRequests(repoUser, repoName)

	if pullRequestError != nil {
		return "", pullRequestError
	}

	for _, pullRequest := range *pullRequests {
		if pullRequest.Head.Ref == branch {
			pullRequestURL := pullRequest.Links.HTML.Href
			githubClient.Logger.Logf("Got pull request URL \"%s\"", pullRequestURL)
			return pullRequestURL, nil
		}
	}

	return "", nil
}
