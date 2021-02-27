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

package jetbrains

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/util"
	"regexp"
	"strings"
)

// ReleaseCodes provides a mapping between JetBrains products and their API codename
var ReleaseCodes = map[string]string{
	"appcode":              "AC",
	"clion":                "CL",
	"datagrip":             "DG",
	"goglang":              "GO",
	"ideaiu":               "IIU",
	"ideaic":               "IIC",
	"phpstorm":             "PS",
	"pycharm-professional": "PCP",
	"pycharm-ce":           "PCC",
	"pycharm-community":    "PCC",
	"pycharm-edu":          "PCE",
	"rider":                "RD",
	"rubymine":             "RM",
	"upsource":             "US",
	"webstorm":             "WS",
}

const (
	// ReleasesAPI is the format string for the JetBrains Releases API
	ReleasesAPI = "https://data.services.jetbrains.com/products/releases?code=%s"
	// LatestAPI is the format string for the JetBrains Releases API when asking for latest
	LatestAPI = "https://data.services.jetbrains.com/products/releases?code=%s&latest=true"
)

// SourceRegex matches JetBrains sources
var SourceRegex = regexp.MustCompile("https?://download.jetbrains.com/.+?/(.+?)-\\d.*")

// Provider is the upstream provider interface for JetBrains
type Provider struct{}

// String gives the name of this provider
func (c Provider) String() string {
	return "JetBrains"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) (params []string) {
	if sm := SourceRegex.FindStringSubmatch(query); len(sm) > 1 {
		params = append(params, strings.ToLower(sm[1]))
	}
	return
}

// Latest finds the newest release for a JetBrains package
func (c Provider) Latest(params []string) (r *results.Result, err error) {
	rs, err := c.fetchReleases(LatestAPI, "latest", params)
	if err != nil {
		return
	}
	r = rs.First()
	return
}

// Releases finds all matching releases for a JetBrains package
func (c Provider) Releases(params []string) (rs *results.ResultSet, err error) {
	return c.fetchReleases(ReleasesAPI, "releases", params)
}

func (c Provider) fetchReleases(api, kind string, params []string) (rs *results.ResultSet, err error) {
	name := params[0]
	// Query the API
	code := ReleaseCodes[name]
	var jbs Releases
	url := fmt.Sprintf(api, code)
	if err = util.FetchJSON(url, kind, &jbs); err != nil {
		return
	}
	if jbs[code] == nil || len(jbs[code]) == 0 {
		err = results.NotFound
		return
	}
	rs = jbs.Convert(name)
	return
}
