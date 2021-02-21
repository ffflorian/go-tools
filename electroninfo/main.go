/*
Copyright © 2021 Florian Imdahl <git@ffflorian.de>

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

var logger = simplelogger.New("electroninfo", false, true)

func main() {
	util := util.New(name, version, description)

	util.CheckFlags()

	debugMode := util.FlagContext.Bool("d")
	timeout := util.FlagContext.Int("t")

	logger.Enabled = true

	httpService := httpservice.New(timeout, debugMode)

	buildTable(httpService)
}

func buildTable(httpService *httpservice.HTTPService) {
	var releases, releasesErr = httpService.GetReleases()
	if releasesErr != nil {
		logger.Error(releasesErr)
		os.Exit(1)
	}

	firstRelease := (*releases)[0]
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk([][]string{
		{"Electron", firstRelease.Version},
		{"Published on", firstRelease.PublishedAt},
		{"Node.js", firstRelease.Deps.Node},
		{"Chrome", firstRelease.Deps.Chrome},
		{"OpenSSL", firstRelease.Deps.OpenSSL},
		{"Modules (Node ABI)", firstRelease.Deps.Modules},
		{"uv", firstRelease.Deps.Uv},
		{"V8", firstRelease.Deps.V8},
		{"zlib", firstRelease.Deps.Zlib},
	})
	table.Render()
}
