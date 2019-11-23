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
	u "github.com/ffflorian/go-tools/gh-open/util"
	"github.com/skratchdot/open-golang/open"
)

const (
	description = "Open a GitHub repository in your browser."
	name        = "gh-open"
	version     = "0.0.3"
)

func main() {
	justPrint := false
	justBranch := false
	logger := simplelogger.New(false, true)
	util := u.New(name, version, description)

	util.CheckFlags()

	if util.FlagContext.IsSet("d") {
		logger.Enabled = true
	}

	logger.Log("Got arguments:", util.FlagContext.Args()[1:])

	if util.FlagContext.IsSet("v") {
		util.LogAndExit(version)
	}

	if util.FlagContext.IsSet("h") {
		util.LogAndExit(util.GetUsage())
	}

	if util.FlagContext.IsSet("p") {
		justPrint = true
	}

	if util.FlagContext.IsSet("b") {
		justBranch = true
	}

	gitClient := gitclient.New(logger)

	argsDir, argsDirError := util.GetArgsDir()
	util.CheckError(argsDirError, true)

	mainDir, absError := filepath.Abs(argsDir)
	util.CheckError(absError, false)

	fullURL, fullURLError := gitClient.GetFullURL(mainDir)
	util.CheckError(fullURLError, false)

	if justBranch == false {
		pullRequest, pullRequestError := gitClient.GetPullRequestURL(fullURL)
		if pullRequestError != nil {
			logger.Log(pullRequestError)
		}
		if pullRequest != "" {
			fullURL = pullRequest
		}
	}

	if justPrint == true {
		util.LogAndExit(fullURL)
	}

	open.Run(fullURL)
}
