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

package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ffflorian/go-tools/gh-open/github"
	"github.com/ffflorian/go-tools/gh-open/simplelogger"
)

// Client is a configuration struct for the git client
type Client struct {
	DebugMode bool
	Logger    *simplelogger.SimpleLogger
	Timeout   int
}

const (
	fullURLRegex     = `(?i)^(?:.+?://(?:.+@)?|(?:.+@)?)(.+?)[:/](.+?)(?:.git)?/?$`
	gitBranchRegex   = `(?mi)ref: refs/heads/(.*)$`
	pullRequestRegex = `(?i)github\.com/([^\/]+)/([^/]+)/tree/(.*)`
	rawURLRegex      = `(?mi).*url = (.*)`
)

// New returns a new instance of Client
func New(timeout int, debugMode bool) *Client {
	logger := simplelogger.New("gh-open/gitclient", debugMode, true)
	gitClient := &Client{
		DebugMode: debugMode,
		Logger:    logger,
		Timeout:   timeout,
	}

	return gitClient
}

func (gitClient *Client) readFile(fileName string) (*[]byte, error) {
	file, openError := os.Open(fileName)

	defer file.Close()

	if openError != nil {
		return nil, openError
	}

	content, readError := ioutil.ReadAll(file)

	if readError != nil {
		return nil, readError
	}

	return &content, nil
}

// ParseBranch takes a git directory and returns it's current branch.
func (gitClient *Client) ParseBranch(gitDir string) ([]byte, error) {
	gitHeadFile, absError := filepath.Abs(filepath.Join(gitDir, "HEAD"))

	if absError != nil {
		return nil, absError
	}

	if _, statError := os.Stat(gitHeadFile); os.IsNotExist(statError) {
		return nil, fmt.Errorf("Could not find git HEAD file in \"%s\"", gitDir)
	}

	gitHead, readFileError := gitClient.readFile(gitHeadFile)

	if readFileError != nil {
		return nil, readFileError
	}

	gitClient.Logger.Logf("Read git head file: \"%s\"", strings.TrimSpace(string(*gitHead)))

	gitBranchRegExp := regexp.MustCompile(gitBranchRegex)
	branchMatches := gitBranchRegExp.FindSubmatch(*gitHead)

	if len(branchMatches) != 2 {
		return nil, errors.New("No branch found in git HEAD file")
	}

	return branchMatches[1], nil
}

// ParseRawURL takes a git directory and returns it's raw URL.
func (gitClient *Client) ParseRawURL(gitDir string) ([]byte, error) {
	gitConfigFile, absError := filepath.Abs(filepath.Join(gitDir, "config"))

	if absError != nil {
		return nil, absError
	}

	gitClient.Logger.Logf("Found git config file \"%s\"", gitConfigFile)

	if _, statError := os.Stat(gitConfigFile); os.IsNotExist(statError) {
		return nil, fmt.Errorf("Could not find git config file in \"%s\"", gitDir)
	}

	gitConfig, readFileError := gitClient.readFile(gitConfigFile)

	if readFileError != nil {
		return nil, readFileError
	}

	rawURLRegExp := regexp.MustCompile(rawURLRegex)
	rawURLMatches := rawURLRegExp.FindSubmatch(*gitConfig)

	if len(rawURLMatches) != 2 {
		return nil, errors.New("No raw URL found in git config file")
	}

	return rawURLMatches[1], nil
}

// FindGitDir takes a directory and returns it's next git directory.
func (gitClient *Client) FindGitDir(mainDir string) (string, error) {
	foundDir, walkError := gitClient.findUp(mainDir, ".git")

	if walkError != nil {
		return "", walkError
	}

	return foundDir, nil
}

func (gitClient *Client) findUp(initialDir string, targetDir string) (string, error) {
	var mainDir = &initialDir

	if _, statError := os.Stat(initialDir); os.IsNotExist(statError) {
		return "", fmt.Errorf("Could not find the directory \"%s\"", initialDir)
	}

	for {
		var joinedPath = filepath.Join(*mainDir, targetDir)
		gitClient.Logger.Logf("Searching for git dir in \"%s\"", *mainDir)

		if _, statError := os.Stat(joinedPath); os.IsNotExist(statError) {
			absoluteDir, absError := filepath.Abs(filepath.Join(*mainDir, "../"))

			if absError != nil {
				fmt.Println(absError)
				return "", absError
			}

			if filepath.Clean(absoluteDir) == "/" {
				return "", fmt.Errorf("Could not find a git repository in \"%s\"", initialDir)
			}

			mainDir = &absoluteDir
			continue
		} else if _, statError := os.Stat(joinedPath); !os.IsNotExist(statError) {
			return joinedPath, nil
		}
	}
}

// GetFullURL takes a directory and (given it's inside a git repository) returns the repository's full URL.
func (gitClient *Client) GetFullURL(mainDir string) (string, error) {
	gitDir, findGitDirError := gitClient.FindGitDir(mainDir)

	if findGitDirError != nil {
		return "", findGitDirError
	}

	gitClient.Logger.Logf("Found git dir \"%s\"", string(gitDir))

	gitRawURL, gitRawURLError := gitClient.ParseRawURL(gitDir)

	if gitRawURLError != nil {
		return "", gitRawURLError
	}

	gitClient.Logger.Logf("Found raw URL \"%s\"", string(gitRawURL))

	gitBranch, gitBranchError := gitClient.ParseBranch(gitDir)

	if gitBranchError != nil {
		return "", gitBranchError
	}

	fullURLRegExp := regexp.MustCompile(fullURLRegex)
	fullURLMatches := fullURLRegExp.FindSubmatch(gitRawURL)

	if len(fullURLMatches) != 3 {
		return "", errors.New("Could not convert raw URL")
	}

	parsedURL := fullURLRegExp.ReplaceAll(gitRawURL, []byte("https://$1/$2"))
	gitClient.Logger.Logf("Found parsed URL \"%s\"", string(parsedURL))

	fullURL := fmt.Sprintf("%s/tree/%s", string(parsedURL), string(gitBranch))

	return fullURL, nil
}

// GetPullRequestURL gets the according pull request URL from GitHub
// if there is one
func (gitClient *Client) GetPullRequestURL(gitFullURL string) (string, error) {
	pullRequestRegExp := regexp.MustCompile(pullRequestRegex)
	fullURLMatches := pullRequestRegExp.FindStringSubmatch(gitFullURL)

	if len(fullURLMatches) != 4 {
		return "", errors.New("Could not convert GitHub URL to pull request")
	}

	repoUser := fullURLMatches[1]
	repoName := fullURLMatches[2]
	repoBranch := fullURLMatches[3]

	gitClient.Logger.Logf("Got user \"%s\", repo name \"%s\" and branch \"%s\"", repoUser, repoName, repoBranch)

	githubClient := github.New(gitClient.Timeout, gitClient.DebugMode)

	pullRequest, pullRequestError := githubClient.GetPullRequestByBranch(repoUser, repoName, repoBranch)
	if pullRequestError != nil {
		return "", pullRequestError
	}

	return pullRequest, nil
}
