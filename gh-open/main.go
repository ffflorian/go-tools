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
	"path/filepath"

	"github.com/ffflorian/go-tools/gh-open/git"
	"github.com/ffflorian/go-tools/gh-open/util"
	"github.com/skratchdot/open-golang/open"
)

const name = "gh-open"
const version = "0.0.2"
const description = "Open a GitHub repository in your browser."

func init() {
	util.CheckFlags(name, version, description)
}

func main() {
	argsDir, argsDirError := util.GetArgsDir()
	util.CheckError(argsDirError)

	mainDir, absError := filepath.Abs(argsDir)
	util.CheckError(absError)

	fullURL, fullURLError := git.GetFullURL(mainDir)
	util.CheckError(fullURLError)

	if util.GetFlagContext().IsSet("p") {
		util.PrintAndExit(fullURL)
	}

	open.Run(fullURL)
}
