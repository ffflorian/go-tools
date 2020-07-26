/*
Copyright © 2019 Florian Keller <git@ffflorian.de>

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

package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ffflorian/go-tools/simplelogger"
)

// Client is a configuration struct for the GitHub client
type Client struct {
	DebugMode bool
	Logger    *simplelogger.SimpleLogger
	Timeout   int
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

// New returns a new instance of Client
func New(timeout int, debugMode bool) *Client {
	logger := simplelogger.New("gh-open/githubclient", debugMode, true)

	return &Client{
		DebugMode: debugMode,
		Logger:    logger,
		Timeout:   timeout,
	}
}

func (githubClient *Client) request(urlPath string) (*[]byte, error) {
	timeout := time.Duration(githubClient.Timeout) * time.Millisecond
	httpClient := &http.Client{Timeout: timeout}
	fullURL := fmt.Sprintf("%s/%s", baseURL, urlPath)

	githubClient.Logger.Logf("Sending GET request to \"%s\" with timeout \"%s\" ...", fullURL, timeout)

	response, responseError := httpClient.Get(fullURL)
	if responseError != nil {
		return nil, responseError
	}

	defer response.Body.Close()

	githubClient.Logger.Logf("Got response status code \"%d\"", response.StatusCode)

	if response.StatusCode != 200 {
		return nil, errors.New("Invalid response status code")
	}

	buffer, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return nil, readError
	}

	return &buffer, nil
}

// GetPullRequests gets pull requests from GitHub,
// see https://developer.github.com/v3/pulls/#list-pull-requests
func (githubClient *Client) GetPullRequests(repoUser string, repoName string) (*[]PullRequest, error) {
	var pullRequests *[]PullRequest

	urlPath := fmt.Sprintf("repos/%s/%s/pulls", repoUser, repoName)
	requestBuffer, requestError := githubClient.request(urlPath)
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
func (githubClient *Client) GetPullRequestByBranch(repoUser string, repoName string, branch string) (string, error) {
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
