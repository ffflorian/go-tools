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

// ParseGitBranch takes a git directory and returns it's current branch.
func ParseGitBranch(gitDir string) string {
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

	file, openError := os.Open((gitHeadFile))

	if openError != nil {
		fmt.Println(openError)
		os.Exit(1)
	}
	defer file.Close()

	gitHead, readError := ioutil.ReadAll(file)

	if readError != nil {
		fmt.Println(readError)
		os.Exit(1)
	}

	var gitBranchRegEx = regexp.MustCompile(`(?mi)ref: refs/heads/(.*)$`)
	var branch = gitBranchRegEx.FindSubmatch(gitHead)

	if len(branch) != 2 {
		fmt.Println("No branch found in git HEAD file")
		os.Exit(1)
	}

	return string(branch[1])
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
