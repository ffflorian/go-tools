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

package httpservice

import (
	"github.com/ffflorian/go-tools/simplelogger"
)

// HTTPService is a configuration struct for the HTTPService
type HTTPService struct {
	DebugMode bool
	Logger    *simplelogger.SimpleLogger
	Timeout   int
}

// New returns a new instance of the HTTPService
func New(timeout int, debugMode bool) *HTTPService {
	logger := simplelogger.New("electroninfo/httpservice", debugMode, true)

	return &HTTPService{
		DebugMode: debugMode,
		Logger:    logger,
		Timeout:   timeout,
	}
}

func (httpservice *HTTPService) Hello() {}
