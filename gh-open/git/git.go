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
)

const (
	fullURLRegex   = `(?i)^(?:.+?://(?:.+@)?|(?:.+@)?)(.+?)[:/](.+?)(?:.git)?/?$`
	gitBranchRegex = `(?mi)ref: refs/heads/(.*)$`
	rawURLRegex    = `(?mi).*url = (.*)`
)

func readFile(fileName string) ([]byte, error) {
	file, openError := os.Open(fileName)

	defer file.Close()

	if openError != nil {
		return []byte(""), openError
	}

	content, readError := ioutil.ReadAll(file)

	if readError != nil {
		return []byte(""), readError
	}

	return content, nil
}

// ParseBranch takes a git directory and returns it's current branch.
func ParseBranch(gitDir string) ([]byte, error) {
	gitHeadFile, absError := filepath.Abs(filepath.Join(gitDir, "HEAD"))

	if absError != nil {
		fmt.Println(absError)
		os.Exit(1)
	}

	if _, statError := os.Stat(gitHeadFile); os.IsNotExist(statError) {
		return []byte(""), fmt.Errorf("Could not find git HEAD file in \"%s\"", gitDir)
	}

	gitHead, readFileError := readFile(gitHeadFile)

	if readFileError != nil {
		return []byte(""), readFileError
	}

	gitBranchRegExp := regexp.MustCompile(gitBranchRegex)
	branch := gitBranchRegExp.FindSubmatch(gitHead)

	if len(branch) != 2 {
		return []byte(""), errors.New("No branch found in git HEAD file")
	}

	return branch[1], nil
}

// ParseRawURL takes a git directory and returns it's raw URL.
func ParseRawURL(gitDir string) ([]byte, error) {
	gitConfigFile, absError := filepath.Abs(filepath.Join(gitDir, "config"))

	if absError != nil {
		return []byte(""), absError
	}

	if _, statError := os.Stat(gitConfigFile); os.IsNotExist(statError) {
		return []byte(""), fmt.Errorf("Could not find git config file in \"%s\"", gitDir)
	}

	gitConfig, readFileError := readFile(gitConfigFile)

	if readFileError != nil {
		return []byte(""), readFileError
	}

	rawURLRegExp := regexp.MustCompile(rawURLRegex)
	branch := rawURLRegExp.FindSubmatch(gitConfig)

	if len(branch) != 2 {
		return []byte(""), errors.New("No branch found in git config file")
	}

	return branch[1], nil
}

// FindGitDir takes a directory and returns it's next git directory.
func FindGitDir(mainDir string) (string, error) {
	foundDir, walkError := walk(mainDir, ".git")

	if walkError != nil {
		return "", walkError
	}

	return foundDir, nil
}

func walk(mainDir string, targetDir string) (string, error) {
	for {
		var joinedPath = filepath.Join(mainDir, targetDir)

		if _, statError := os.Stat(joinedPath); os.IsNotExist(statError) {
			var absError error
			mainDir, absError = filepath.Abs(filepath.Join(mainDir, "../"))

			if absError != nil {
				fmt.Println(absError)
				return "", absError
			}

			if filepath.Clean(mainDir) == "/" {
				return "", fmt.Errorf("Could not find a git repository in \"%s\"", mainDir)
			}

			continue
		} else if _, statError := os.Stat(joinedPath); !os.IsNotExist(statError) {
			return joinedPath, nil
		}
	}
}

// GetFullURL takes a directory and (given it's inside a git repository) returns the repository's full URL.
func GetFullURL(mainDir string) (string, error) {
	gitDir, findGitDirError := FindGitDir(mainDir)

	if findGitDirError != nil {
		return "", findGitDirError
	}

	gitRawURL, gitRawURLError := ParseRawURL(gitDir)

	if gitRawURLError != nil {
		return "", gitRawURLError
	}

	gitBranch, gitBranchError := ParseBranch(gitDir)

	if gitBranchError != nil {
		return "", gitBranchError
	}

	fullURLRegExp := regexp.MustCompile(fullURLRegex)
	fullURLMatch := fullURLRegExp.FindSubmatch(gitRawURL)

	if len(fullURLMatch) != 3 {
		fmt.Println("Could not convert raw URL")
		os.Exit(1)
	}

	parsedURL := fullURLRegExp.ReplaceAll(gitRawURL, []byte("https://$1/$2"))
	fullURL := fmt.Sprintf("%s/tree/%s", string(parsedURL), string(gitBranch))

	return fullURL, nil
}
