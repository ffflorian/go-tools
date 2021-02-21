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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ffflorian/go-tools/simplelogger"
)

// HTTPService is a configuration struct for the HTTPService
type HTTPService struct {
	DebugMode bool
	Logger    *simplelogger.SimpleLogger
	Timeout   int
}

// RawReleaseInfo describes the raw data of a release info
type RawReleaseInfo struct {
	Deps struct {
		Chrome  string `json:"chrome"`
		Modules string `json:"modules"`
		Node    string `json:"node"`
		OpenSSL string `json:"openssl"`
		Uv      string `json:"uv"`
		V8      string `json:"v8"`
		Zlib    string `json:"zlib"`
	} `json:"deps"`
	Name           string   `json:"name"`
	NodeID         string   `json:"node_id"`
	NpmDistTags    []string `json:"npm_dist_tags"`
	NpmPackageName string   `json:"npm_package_name"`
	Prerelease     bool     `json:"prerelease"`
	PublishedAt    string   `json:"published_at"`
	TagName        string   `json:"tag_name"`
	TotalDownloads int      `json:"total_downloads"`
	Version        string   `json:"version"`
}

// New returns a new instance of httpService
func New(timeout int, debugMode bool) *HTTPService {
	logger := simplelogger.New("electroninfo/httpService", debugMode, true)

	return &HTTPService{
		DebugMode: debugMode,
		Logger:    logger,
		Timeout:   timeout,
	}
}

func (httpservice *HTTPService) request(url string) (*[]byte, error) {
	var defaultTimeout = 2000
	timeout := time.Duration(defaultTimeout) * time.Millisecond
	httpClient := &http.Client{Timeout: timeout}

	httpservice.Logger.Logf("Downloading from \"%s\" with timeout \"%s\" ...", url, timeout)

	response, responseError := httpClient.Get(url)
	if responseError != nil {
		return nil, responseError
	}

	defer response.Body.Close()

	httpservice.Logger.Logf("Got response status code \"%d\"", response.StatusCode)

	if response.StatusCode != 200 {
		return nil, errors.New("Invalid response status code")
	}

	buffer, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return nil, readError
	}

	return &buffer, nil
}

// GetReleases downloads the releases file
func (httpservice *HTTPService) GetReleases() (*[]RawReleaseInfo, error) {
	var releases *[]RawReleaseInfo
	const downloadURL = "https://raw.githubusercontent.com/electron/releases/master/lite.json"

	requestBuffer, requestError := httpservice.request(downloadURL)
	if requestError != nil {
		return nil, requestError
	}

	unmarshalError := json.Unmarshal(*requestBuffer, &releases)
	if unmarshalError != nil {
		return nil, unmarshalError
	}

	httpservice.Logger.Logf("Got %d releases", len(*releases))

	return releases, nil
}
