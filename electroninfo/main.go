/*
Copyright © 2019 Florian Imdahl <git@ffflorian.de>

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
	"os"

	"github.com/ffflorian/go-tools/electroninfo/fileservice"
	"github.com/ffflorian/go-tools/electroninfo/httpservice"
	"github.com/ffflorian/go-tools/electroninfo/util"
	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/olekukonko/tablewriter"
)

const (
	description = "Get information about electron releases."
	name        = "electron-info"
	version     = "0.0.1"
)

func main() {
	var (
		logger      = simplelogger.New("electroninfo", false, true)
		tablewriter = tablewriter.NewWriter(os.Stdout)
		util        = util.New(name, version, description)
	)

	util.CheckFlags()

	timeout := util.FlagContext.Int("t")
	debug := util.FlagContext.Bool("d")
	force := util.FlagContext.Bool("f")
	version := util.FlagContext.Bool("v")
	raw := util.FlagContext.Bool("r")
	help := util.FlagContext.Bool("h")

	if debug == true {
		logger.Enabled = true
	}

	fileservice := fileservice.New(debug)
	httpservice := httpservice.New(timeout, debug)

	logger.Log("Got arguments:", util.FlagContext.Args()[1:])

	if version == true {
		util.LogAndExit(version)
	}

	if help == true {
		util.LogAndExit(util.GetUsage())
	}

	if util.FlagContext.IsSet("t") {
		util.FlagContext.Int("t")
	}

	fileservice.Hello()
	httpservice.Hello()

	// gitClient := git.New(timeout, debugMode)

	// argsDir, argsDirError := util.GetArgsDir()
	// util.CheckError(argsDirError, true)

	// mainDir, absError := filepath.Abs(argsDir)
	// util.CheckError(absError, false)

	// fullURL, fullURLError := gitClient.GetFullURL(mainDir)
	// util.CheckError(fullURLError, false)

	// if justBranch == false {
	// 	pullRequest, pullRequestError := gitClient.GetPullRequestURL(fullURL)
	// 	if pullRequestError != nil {
	// 		logger.Error(pullRequestError)
	// 	}
	// 	if pullRequest != "" {
	// 		fullURL = pullRequest
	// 	}
	// }

	if raw == true {
		util.LogAndExit("hey there")
	}
}
