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
	"os"
	"path/filepath"

	"github.com/ffflorian/go-tools/gh-open/git"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var (
	justPrint    bool
	openBranch   bool
	printVersion bool
)

const version = "0.0.1"

var rootCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("too many arguments")
		}
		return nil
	},
	Use:   "gh-open",
	Short: "Open a GitHub repository in your browser.",
	Run: func(cmd *cobra.Command, args []string) {
		if printVersion == true {
			fmt.Println(version)
			os.Exit(0)
		}

		var mainDir string
		var absError error

		if len(args) > 0 {
			mainDir = args[0]
		}

		mainDir, absError = filepath.Abs(mainDir)

		if absError != nil {
			fmt.Println(absError)
			os.Exit(1)
		}

		fullURL, fullURLError := git.GetFullURL(mainDir)

		if fullURLError != nil {
			fmt.Println(fullURLError)
			os.Exit(1)
		}

		if justPrint == true {
			fmt.Println(fullURL)
			os.Exit(0)
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
