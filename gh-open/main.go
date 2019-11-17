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

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var cfgFile string
var justPrint bool
var openBranch bool
var printVersion bool
var version string = "0.0.1"

var rootCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("too many arguments")
		}
		return nil
	},
	Use:   "gh-open",
	Short: "Open a GitHub repository in your browser.",
	Long:  "Open a GitHub repository in your browser.",
	Run: func(cmd *cobra.Command, args []string) {
		if printVersion == true {
			fmt.Println(version)
			os.Exit(0)
		}

		var mainDir string
		var gitDir string = "."
		var absError error

		if len(args) > 0 {
			mainDir = args[0]
		}

		mainDir, absError = filepath.Abs(mainDir)

		if absError != nil {
			fmt.Println(absError)
			os.Exit(1)
		}

		fmt.Println("mainDir", mainDir)

		gitDir = findGitDir(mainDir)
		fmt.Println("gitDir", gitDir)

		fullURL := getFullURL(gitDir)

		fmt.Println("full URL", fullURL)

		if justPrint == true {
			fmt.Println(fullURL)
			os.Exit((0))
		}

		open.Run(fullURL)
	},
}

func main() {
	if cmdError := rootCmd.Execute(); cmdError != nil {
		fmt.Println(cmdError)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&justPrint, "print", "p", false, "just print the URL")
	rootCmd.PersistentFlags().BoolVarP(&openBranch, "branch", "b", false, "open the branch tree (and not the PR)")
	rootCmd.PersistentFlags().BoolVarP(&printVersion, "version", "v", false, "output the version number")
}

func getFullURL(gitDir string) string {
	branch := parseGitBranch(gitDir)

	fmt.Println("branch", branch)

	return ""
}

func parseGitBranch(gitDir string) string {
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

func findGitDir(mainDir string) string {
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
