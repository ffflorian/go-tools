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

package nominatimclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/ffflorian/go-tools/gh-open/simplelogger"
)

// Client is a configuration struct for nominatim
type Client struct {
	DebugMode bool
	Logger    *simplelogger.SimpleLogger
	Timeout   int
}

// Location represents a location result from Nominatim
type Location struct {
	Boundingbox []string `json:"boundingbox"`
	Class       string   `json:"class"`
	DisplayName string   `json:"display_name"`
	Icon        string   `json:"icon,omitempty"`
	Importance  float64  `json:"importance"`
	Lat         string   `json:"lat"`
	Licence     string   `json:"licence"`
	Lon         string   `json:"lon"`
	OsmID       int      `json:"osm_id"`
	OsmType     string   `json:"osm_type"`
	PlaceID     int      `json:"place_id"`
	Type        string   `json:"type"`
}

const baseURL = "https://nominatim.openstreetmap.org"

// New returns a new instance of Nominatim
func New(timeout int, debugMode bool) *Client {
	logger := simplelogger.New("my-timezone/nominatim", debugMode, true)
	return &Client{
		DebugMode: debugMode,
		Logger:    logger,
		Timeout:   timeout,
	}
}

func (nominatimClient *Client) request(urlPath string) (*[]byte, error) {
	timeout := time.Duration(nominatimClient.Timeout) * time.Millisecond
	httpClient := &http.Client{Timeout: timeout}
	fullURL := fmt.Sprintf("%s/%s", baseURL, urlPath)

	nominatimClient.Logger.Logf("Sending GET request to \"%s\" with timeout \"%s\" ...", fullURL, timeout)

	response, responseError := httpClient.Get(fullURL)
	if responseError != nil {
		return nil, responseError
	}

	defer response.Body.Close()

	nominatimClient.Logger.Logf("Got response status code \"%d\"", response.StatusCode)

	buffer, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return nil, readError
	}

	return &buffer, nil
}

// GetLongitudeByName takes a city name and returns it's longitude by using
// the Nominatim API (https://nominatim.org/release-docs/develop/api/Overview/).
func (nominatimClient *Client) GetLongitudeByName(locationName string) (float64, error) {
	locationsPtr, locationsError := nominatimClient.getLocationsByName(locationName)

	if locationsError != nil {
		return 0, locationsError
	}

	locations := *locationsPtr

	if len(locations) == 0 {
		return 0, errors.New("Could not find any place with that name")
	}

	parsedLongitude, parseFloatError := strconv.ParseFloat(locations[0].Lon, 64)

	if parseFloatError != nil {
		return 0, parseFloatError
	}

	return parsedLongitude, nil
}

func (nominatimClient *Client) getLocationsByName(locationName string) (*[]Location, error) {
	var locations *[]Location

	urlPath := fmt.Sprintf("search?q=%s&limit=9&format=json", locationName)
	requestBuffer, requestError := nominatimClient.request(urlPath)
	if requestError != nil {
		return nil, requestError
	}

	unmarshalError := json.Unmarshal(*requestBuffer, &locations)
	if unmarshalError != nil {
		return nil, unmarshalError
	}

	nominatimClient.Logger.Logf("Got %d locations", len(*locations))

	return locations, nil
}
