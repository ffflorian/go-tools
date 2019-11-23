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

package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/simonleung8/flags"
)

var flagContext flags.FlagContext

func init() {
	flagContext = flags.New()
}

// GetFlagContext just returns the flag context
func GetFlagContext() flags.FlagContext {
	return flagContext
}

// CheckFlags checks which command line flags are set
func CheckFlags(name string, version string, description string) {
	flagContext.NewBoolFlag("print", "p", "just print the URL")
	flagContext.NewBoolFlag("version", "v", "output the version number")
	flagContext.NewBoolFlag("help", "h", "output usage information")
	// fc.NewBoolFlag("branch", "b", "open the branch tree (and not the PR)")

	parseError := flagContext.Parse(os.Args...)
	CheckError(parseError)
}

// GetArgsDir returns the directory provided via arguments
func GetArgsDir() (string, error) {
	args := flagContext.Args()

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
func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// PrintUsageAndExit prints the usage text and exits with exit code 0
func PrintUsageAndExit(name string, description string) {
	fmt.Printf(
		"%s\n\nUsage:\n%s [flags] [directory]\n\nFlags:\n%s",
		description,
		name,
		flagContext.ShowUsage(2),
	)
	os.Exit(0)
}

// PrintAndExit prints one or more messages and exits with exit code 0
func PrintAndExit(messages ...interface{}) {
	fmt.Println(messages...)
	os.Exit(0)
}
