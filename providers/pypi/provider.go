//
// Copyright 2016-2018 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pypi

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
	"net/http"
	"regexp"
	"strings"
)

const (
	// SourceAPI is the format string for the PyPi API
	SourceAPI = "https://pypi.python.org/pypi/%s/json"
	// DateFormat is the Time format used by PyPi
	DateFormat = "2006-01-02T15:04:05"
)

// TarballRegex matches PyPi source tarballs
var TarballRegex = regexp.MustCompile("https?://[^/]*py[^/]*/packages/(?:[^/]+/)+(.+)$")

// Provider is the upstream provider interface for pypi
type Provider struct{}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "PyPi"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	if sm := TarballRegex.FindStringSubmatch(query); len(sm) > 1 {
		pieces := strings.Split(sm[1], "-")
		if len(pieces) > 2 {
			return strings.Join(pieces[0:len(pieces)-1], "-")
		}
		return pieces[0]
	}
	return ""
}

// Latest finds the newest release for a pypi package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	// Query the API
	resp, err := http.Get(fmt.Sprintf(SourceAPI, name))
	if err != nil {
		log.Debugf("Failed to get latest: %s\n", err)
		err = results.Unavailable
		return
	}
	defer resp.Body.Close()
	// Translate Status Code
	switch resp.StatusCode {
	case 200:
		break
	case 404:
		err = results.NotFound
		return
	default:
		err = results.Unavailable
		return
	}
	// Decode response
	dec := json.NewDecoder(resp.Body)
	var cr LatestSource
	if err = dec.Decode(&cr); err != nil {
		log.Debugf("Failed to decode latest: %s\n", err)
		err = results.Unavailable
		return
	}
	r = cr.Convert(name)
	return
}

// Releases finds all matching releases for a pypi package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
	// Query the API
	resp, err := http.Get(fmt.Sprintf(SourceAPI, name))
	if err != nil {
		log.Debugf("Failed to get releases: %s\n", err)
		err = results.Unavailable
		return
	}
	defer resp.Body.Close()
	// Translate Status Code
	switch resp.StatusCode {
	case 200:
		break
	case 404:
		err = results.NotFound
		return
	default:
		err = results.Unavailable
		return
	}
	// Decode response
	dec := json.NewDecoder(resp.Body)
	var crs Releases
	if err = dec.Decode(&crs); err != nil {
		log.Debugf("Failed to decode releases: %s\n", err)
		err = results.Unavailable
		return
	}
	if len(crs.Releases) == 0 {
		err = results.NotFound
		return
	}
	rs = crs.Convert(name)
	err = nil
	return
}
