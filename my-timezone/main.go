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
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/beevik/ntp"
	"github.com/ffflorian/go-tools/my-timezone/nominatim"
	"github.com/ffflorian/go-tools/my-timezone/util"
	"github.com/ffflorian/go-tools/simplelogger"
)

const (
	description = "Calculate the exact time based on your location."
	name        = "my-timezone"
	version     = "0.0.1"
)

var (
	debugMode   = false
	offlineMode = false
	timeout     int
)

func main() {
	var (
		logger = simplelogger.New("my-timezone", false, true)
		utils  = util.New(name, version, description)
	)

	utils.CheckFlags()

	ntpServer := utils.FlagContext.String("s")
	offlineMode = utils.FlagContext.Bool("o")
	debugMode = utils.FlagContext.Bool("d")
	timeout = utils.FlagContext.Int("t")

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

	argsLocation, argsLocationError := utils.GetArgsLocation()
	utils.CheckError(argsLocationError, true)

	parsedLongitude, parseLongitudeError := locationToLongitude(argsLocation)
	utils.CheckError(parseLongitudeError, true)

	myTime, getTimeError := getTimeByLocation(ntpServer, parsedLongitude)
	utils.CheckError(getTimeError, false)

	fmt.Printf("Your personal time in \"%s\": %s\n", argsLocation, myTime.Format("15:04:05"))
}

func getUTCDate(ntpServer string) (time.Time, error) {
	options := ntp.QueryOptions{
		Timeout: time.Duration(time.Duration(timeout) * time.Millisecond),
	}

	queryResult, queryError := ntp.QueryWithOptions(ntpServer, options)

	if queryError != nil {
		return time.Now(), queryError
	}

	return queryResult.Time, nil
}

func calculateDistance(from float64, to float64) float64 {
	return math.Abs(from - to)
}

func getTimeByLocation(ntpServer string, longitude float64) (time.Time, error) {
	var now time.Time

	if offlineMode == true {
		now = time.Now()
	} else {
		utcTime, getUTCDateError := getUTCDate(ntpServer)
		if getUTCDateError != nil {
			return time.Now().UTC(), getUTCDateError
		}
		now = utcTime
	}

	distance := calculateDistance(0, longitude)
	distanceSeconds := distance / 0.004167

	if longitude < 0 {
		distanceSeconds = distanceSeconds * -1
	}

	return now.UTC().Add(time.Duration(distanceSeconds) * time.Second), nil
}

func locationToLongitude(location string) (float64, error) {
	longitudeRegExp := regexp.MustCompile(`(?m)[-?\W\d\.]+,([-?\W\d\.]+)`)
	longitudeMatch := longitudeRegExp.FindStringSubmatch(location)

	if len(longitudeMatch) == 0 {
		nominatimClient := nominatim.New(10000, debugMode)
		longitude, longitudeError := nominatimClient.GetLongitudeByName(location)

		if longitudeError != nil {
			return 0, longitudeError
		}

		return longitude, nil
	}

	parsedLongitude, parseFloatError := strconv.ParseFloat(longitudeMatch[1], 64)

	if parseFloatError != nil {
		return 0, parseFloatError
	}

	return parsedLongitude, nil
}
