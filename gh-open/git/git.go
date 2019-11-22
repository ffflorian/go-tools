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

var fullURLRegex = `(?i)^(?:.+?://(?:.+@)?|(?:.+@)?)(.+?)[:/](.+?)(?:.git)?/?$`
var rawURLRegex = `(?mi).*url = (.*)`
var gitBranchRegex = `(?mi)ref: refs/heads/(.*)$`

func readFile(fileName string) []byte {
	file, openError := os.Open(fileName)

	if openError != nil {
		fmt.Println(openError)
		os.Exit(1)
	}
	defer file.Close()

	content, readError := ioutil.ReadAll(file)

	if readError != nil {
		fmt.Println(readError)
		os.Exit(1)
	}

	return content
}

// ParseBranch takes a git directory and returns it's current branch.
func ParseBranch(gitDir string) []byte {
	gitHeadFile, absError := filepath.Abs(filepath.Join(gitDir, "HEAD"))

	if absError != nil {
		fmt.Println(absError)
		os.Exit(1)
	}

	if _, statError := os.Stat(gitHeadFile); os.IsNotExist(statError) {
		fmt.Println("Could not find git HEAD file in", gitDir)
		fmt.Println(statError)
		os.Exit(1)
	}

	gitHead := readFile(gitHeadFile)
	gitBranchRegExp := regexp.MustCompile(gitBranchRegex)
	branch := gitBranchRegExp.FindSubmatch(gitHead)

	if len(branch) != 2 {
		fmt.Println("No branch found in git HEAD file")
		os.Exit(1)
	}

	return branch[1]
}

// ParseRawURL takes a git directory and returns it's raw URL.
func ParseRawURL(gitDir string) []byte {
	gitConfigFile, absError := filepath.Abs(filepath.Join(gitDir, "config"))

	if absError != nil {
		fmt.Println(absError)
		os.Exit(1)
	}

	if _, statError := os.Stat(gitConfigFile); os.IsNotExist(statError) {
		fmt.Println("Could not find git config file in", gitDir)
		fmt.Println(statError)
		os.Exit(1)
	}

	gitHead := readFile(gitConfigFile)
	rawURLRegExp := regexp.MustCompile(rawURLRegex)
	branch := rawURLRegExp.FindSubmatch(gitHead)

	if len(branch) != 2 {
		fmt.Println("No branch found in git HEAD file")
		os.Exit(1)
	}

	return branch[1]
}

// FindGitDir takes a directory and returns it's next git directory.
func FindGitDir(mainDir string) string {
	foundDir, walkError := walk(mainDir, ".git")

	if walkError != nil {
		fmt.Println(walkError)
		os.Exit(1)
	}

	return foundDir
}

func walk(mainDir string, targetDir string) (string, error) {
	var targetPath = mainDir

	for {
		var joinedPath = filepath.Join(targetPath, targetDir)

		if _, statError := os.Stat(joinedPath); os.IsNotExist(statError) {
			var absError error
			targetPath, absError = filepath.Abs(filepath.Join(targetPath, "../"))

			if absError != nil {
				fmt.Println(absError)
				return "", absError
			}

			if filepath.Clean(targetPath) == "/" {
				return "", errors.New("Could not find a git repository in")
			}

			continue
		} else if _, statError := os.Stat(joinedPath); !os.IsNotExist(statError) {
			return joinedPath, nil
		}
	}
}

// GetFullURL takes a directory and (given it's inside a git repository) returns the repository's full URL.
func GetFullURL(mainDir string) string {
	gitDir := FindGitDir(mainDir)

	gitRawURL := ParseRawURL(gitDir)
	gitBranch := ParseBranch(gitDir)
	fullURLRegExp := regexp.MustCompile(fullURLRegex)
	fullURLMatch := fullURLRegExp.FindSubmatch(gitRawURL)

	if len(fullURLMatch) != 3 {
		fmt.Println("Could not convert raw URL")
		os.Exit(1)
	}

	parsedURL := fullURLRegExp.ReplaceAll(gitRawURL, []byte("https://$1/$2"))

	return fmt.Sprintf("%s/tree/%s", string(parsedURL), string(gitBranch))
}
