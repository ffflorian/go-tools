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
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
	"os"
)

var cfgFile string
var justPrint bool
var openBranch bool
var printVersion bool
var version string = "0.0.1"

var rootCmd = &cobra.Command{
	Use:   "gh-open",
	Short: "Open a GitHub repository in your browser.",
	Long:  "Open a GitHub repository in your browser.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args", args)

		if printVersion {
			fmt.Println(version)
			os.Exit(0)
		}

		fmt.Println(getFullURL())

		open.Run("https://google.com/")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&justPrint, "print", "p", false, "just print the URL")
	rootCmd.PersistentFlags().BoolVarP(&openBranch, "branch", "b", false, "open the branch tree (and not the PR)")
	rootCmd.PersistentFlags().BoolVarP(&printVersion, "version", "v", false, "output the version number")
}

func getFullURL() string {
	url, err := gitconfig.OriginURL()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return url
}
