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
	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/skratchdot/open-golang/open"
)

const (
	description = "Open a GitHub repository in your browser."
	name        = "gh-open"
	version     = "0.2.2"
)

func main() {
	var (
		logger = simplelogger.New("gh-open", false, true)
		utils  = util.New(name, version, description)
	)

	utils.CheckFlags()

	justPrint := utils.FlagContext.Bool("p")
	justBranch := utils.FlagContext.Bool("b")
	debugMode := utils.FlagContext.Bool("d")
	timeout := utils.FlagContext.Int("t")

	if debugMode == true {
		logger.Enabled = true
	}

	logger.Log("Got arguments:", utils.FlagContext.Args()[1:])

	if utils.FlagContext.IsSet("v") {
		utils.LogAndExit(version)
	}

	if utils.FlagContext.IsSet("h") {
		utils.LogAndExit(utils.GetUsage())
	}

	if utils.FlagContext.IsSet("t") {
		utils.FlagContext.Int("t")
	}

	gitClient := git.New(timeout, debugMode)

	argsDir, argsDirError := utils.GetArgsDir()
	utils.CheckError(argsDirError, true)

	mainDir, absError := filepath.Abs(argsDir)
	utils.CheckError(absError, false)

	fullURL, fullURLError := gitClient.GetFullURL(mainDir)
	utils.CheckError(fullURLError, false)

	if justBranch == false {
		pullRequest, pullRequestError := gitClient.GetPullRequestURL(fullURL)
		if pullRequestError != nil {
			logger.Error(pullRequestError)
		}
		if pullRequest != "" {
			fullURL = pullRequest
		}
	}

	if justPrint == true {
		utils.LogAndExit(fullURL)
	}

	open.Run(fullURL)
}
