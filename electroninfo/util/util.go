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

package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/simonleung8/flags"
)

// Util is a configuration struct for the util
type Util struct {
	Description string
	FlagContext flags.FlagContext
	Name        string
	Version     string
}

// New returns a new instance of Util
func New(name string, version string, description string) *Util {
	flagContext := flags.New()

	return &Util{
		Description: description,
		FlagContext: flagContext,
		Name:        name,
		Version:     version,
	}
}

// CheckFlags checks which command line flags are set
func (util *Util) CheckFlags() {
	util.FlagContext.NewIntFlagWithDefault("timeout", "t", "Set a custom timeout for HTTP requests", 2000)
	util.FlagContext.NewBoolFlag("debug", "d", "Enable debug logging")
	util.FlagContext.NewBoolFlag("force", "f", "Force downloading the latest release file")
	util.FlagContext.NewBoolFlag("version", "v", "Output the version number")
	util.FlagContext.NewBoolFlag("raw", "r", "Output raw JSON")
	util.FlagContext.NewBoolFlag("help", "h", "Output usage information")

	parseError := util.FlagContext.Parse(os.Args...)
	util.CheckError(parseError, false)
}

// GetArgsDir returns the directory provided via arguments
func (util *Util) GetArgsDir() (string, error) {
	args := util.FlagContext.Args()

	switch len(args) {
	case 0:
	case 1:
		return ".", nil
	case 2:
		return args[1], nil
	default:
		return "", errors.New("Too many arguments")
	}

	return "", nil
}

// CheckError checks the error and if it exists, exits with exit code 1
func (util *Util) CheckError(err error, printUsage bool) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		if printUsage {
			fmt.Fprintln(os.Stderr, util.GetUsage())
		}
		os.Exit(1)
	}
}

// GetUsage returns the usage text
func (util *Util) GetUsage() string {
	return fmt.Sprintf(
		"%s\n\nUsage:\n  %s [options] [directory]\n\nOptions:\n%s",
		util.Description,
		util.Name,
		util.FlagContext.ShowUsage(2),
	)
}

// LogAndExit logs one or more messages and exits with exit code 0
func (util *Util) LogAndExit(messages ...interface{}) {
	fmt.Println(messages...)
	os.Exit(0)
}
