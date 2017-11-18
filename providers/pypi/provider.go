//
// Copyright 2016-2017 Bryan T. Meyers <bmeyers@datadrake.com>
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
	"net/http"
	"regexp"
)

const (
	// SourceAPI is the format string for the PyPi API
	SourceAPI = "https://pypi.python.org/pypi/%s/json"
	// DateFormat is the Time format used by PyPi
	DateFormat = "2006-01-02T15:04:05"
)

// TarballRegex matches PyPi source tarballs
var TarballRegex = regexp.MustCompile("https?://pypi.python.org/packages/\\w+/\\w+/\\w+/(.+)-.+.tar.gz")

// Provider is the upstream provider interface for pypi
type Provider struct{}

// Latest finds the newest release for a pypi package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	// Query the API
	resp, err := http.Get(fmt.Sprintf(SourceAPI, name))
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	// Translate Status Code
	switch resp.StatusCode {
	case 200:
		s = results.OK
	case 404:
		s = results.NotFound
	default:
		s = results.Unavailable
	}

	// Fail if not OK
	if s != results.OK {
		return
	}

	dec := json.NewDecoder(resp.Body)
	cr := &LatestSource{}
	err = dec.Decode(cr)
	if err != nil {
		panic(err.Error())
	}
	r = cr.Convert(name)
	return
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := TarballRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}
	return sm[1]
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "PyPi"
}

// Releases finds all matching releases for a pypi package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	// Query the API
	resp, err := http.Get(fmt.Sprintf(SourceAPI, name))
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	// Translate Status Code
	switch resp.StatusCode {
	case 200:
		s = results.OK
	case 404:
		s = results.NotFound
	default:
		s = results.Unavailable
	}

	// Fail if not OK
	if s != results.OK {
		return
	}

	dec := json.NewDecoder(resp.Body)
	crs := &Releases{}
	err = dec.Decode(crs)
	if err != nil {
		panic(err.Error())
	}
	if len(crs.Releases) == 0 {
		s = results.NotFound
		return
	}
	rs = crs.Convert(name)
	return
}
