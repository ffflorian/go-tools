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

package simplelogger

import (
	"fmt"
	"os"
	"strings"
)

// Logger is a configuration struct for the logger
type Logger struct {
	Enabled bool
}

// New returns a new instance of Logger
func New(enabled bool, checkEnvironment bool) Logger {
	if checkEnvironment == true {
		DEBUG := os.Getenv("DEBUG")
		if strings.Contains(DEBUG, "gh-open") {
			enabled = true
		}
	}

	logger := Logger{Enabled: enabled}
	return logger
}

// Log logs one or more messages if the logger is enabled
func (logger Logger) Log(messages ...interface{}) {
	if logger.Enabled == true {
		fmt.Print("debug: ", fmt.Sprintln(messages...))
	}
}
