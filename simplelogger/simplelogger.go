/*
Copyright © 2019 Florian Keller <git@ffflorian.de>

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

const version = "0.0.3"

// SimpleLogger is a configuration struct for the logger
type SimpleLogger struct {
	Enabled bool
	Prefix  string
}

// New returns a new instance of Logger
func New(prefix string, enabled bool, checkEnvironment bool) *SimpleLogger {
	if checkEnvironment == true {
		DEBUG := os.Getenv("DEBUG")
		if strings.Contains(DEBUG, prefix) {
			enabled = true
		}
	}

	return &SimpleLogger{
		Enabled: enabled,
		Prefix:  prefix,
	}
}

func bold(message string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", message)
}

func red(message string) string {
	return fmt.Sprintf("\033[91m%s\033[0m", message)
}

// Log logs one or more unformatted messages if the logger is enabled
func (logger *SimpleLogger) Log(messages ...interface{}) {
	if logger.Enabled == true {
		fmt.Printf("%s %s", bold(logger.Prefix), fmt.Sprintln(messages...))
	}
}

// Logf logs one or more formatted messages if the logger is enabled
func (logger *SimpleLogger) Logf(format string, messages ...interface{}) {
	if logger.Enabled == true {
		fmt.Printf("%s %s\n", bold(logger.Prefix), fmt.Sprintf(format, messages...))
	}
}

// Error logs one or more unformatted messages to stderr if the logger is enabled
func (logger *SimpleLogger) Error(messages ...interface{}) {
	if logger.Enabled == true {
		fmt.Fprintf(os.Stderr, "%s %s %s", bold(red(logger.Prefix)), red("Error:"), red(fmt.Sprintln(messages...)))
	}
}

// Errorf logs one or more formatted messages to stderr if the logger is enabled
func (logger *SimpleLogger) Errorf(format string, messages ...interface{}) {
	if logger.Enabled == true {
		fmt.Fprintf(os.Stderr, "%s %s %s\n", bold(red(logger.Prefix)), red("Error:"), red(fmt.Sprintf(format, messages...)))
	}
}
