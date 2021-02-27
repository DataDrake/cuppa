//
// Copyright 2016-2021 Bryan T. Meyers <root@datadrake.com>
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
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/util"
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

// String gives the name of this provider
func (c Provider) String() string {
	return "PyPi"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) (params []string) {
	if sm := TarballRegex.FindStringSubmatch(query); len(sm) > 1 {
		pieces := strings.Split(sm[1], "-")
		if len(pieces) > 2 {
			params = append(params, strings.Join(pieces[0:len(pieces)-1], "-"))
			return
		}
		params = append(params, pieces[0])
	}
	return
}

// Latest finds the newest release for a pypi package
func (c Provider) Latest(params []string) (r *results.Result, err error) {
	name := params[0]
	url := fmt.Sprintf(SourceAPI, name)
	var cr LatestSource
	if err = util.FetchJSON(url, "latest", &cr); err == nil {
		r = cr.Convert(name)
	}
	return
}

// Releases finds all matching releases for a pypi package
func (c Provider) Releases(params []string) (rs *results.ResultSet, err error) {
	name := params[0]
	url := fmt.Sprintf(SourceAPI, name)
	var crs Releases
	if err = util.FetchJSON(url, "releases", &crs); err != nil {
		return
	}
	if len(crs.Releases) == 0 {
		err = results.NotFound
		return
	}
	rs = crs.Convert(name)
	return
}
