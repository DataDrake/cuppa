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

package jetbrains

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"os"
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

// Latest finds the newest release for a JetBrains package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	// Query the API
	code := ReleaseCodes[name]
	resp, err := http.Get(fmt.Sprintf(LatestAPI, code))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
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
	jbs := make(Releases)
	err = dec.Decode(&jbs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
	}
	if jbs[code] == nil || len(jbs[code]) == 0 {
		s = results.NotFound
		return
	}
	r = jbs.Convert(name).First()
	return
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := SourceRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}
	return strings.ToLower(sm[1])
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "JetBrains"
}

// Releases finds all matching releases for a JetBrains package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	// Query the API
	code := ReleaseCodes[name]
	resp, err := http.Get(fmt.Sprintf(ReleasesAPI, code))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
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
	jbs := make(Releases)
	err = dec.Decode(&jbs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
	}
	if jbs[code] == nil || len(jbs[code]) == 0 {
		s = results.NotFound
		return
	}
	rs = jbs.Convert(name)
	return
}
