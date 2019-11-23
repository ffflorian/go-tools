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

	"github.com/ffflorian/go-tools/gh-open/gitclient"
	"github.com/ffflorian/go-tools/gh-open/simplelogger"
	"github.com/ffflorian/go-tools/gh-open/util"
	"github.com/skratchdot/open-golang/open"
)

const (
	description = "Open a GitHub repository in your browser."
	name        = "gh-open"
	version     = "0.0.3"
)

var (
	gitClient gitclient.Git
	justPrint = false
	logger    simplelogger.Logger
)

func init() {
	logger = simplelogger.New(false, true)

	util.CheckFlags(name, version, description)

	if util.GetFlagContext().IsSet("d") {
		logger.Enabled = true
	}

	if util.GetFlagContext().IsSet("v") {
		util.LogAndExit(version)
	}

	if util.GetFlagContext().IsSet("h") {
		util.PrintUsageAndExit(name, description)
	}

	if util.GetFlagContext().IsSet("p") {
		justPrint = true
	}

	gitClient = gitclient.New(logger)
}

func main() {
	argsDir, argsDirError := util.GetArgsDir()
	util.CheckError(argsDirError)

	mainDir, absError := filepath.Abs(argsDir)
	util.CheckError(absError)

	fullURL, fullURLError := gitClient.GetFullURL(mainDir)
	util.CheckError(fullURLError)

	if justPrint == true {
		util.LogAndExit(fullURL)
	}

	open.Run(fullURL)
}
