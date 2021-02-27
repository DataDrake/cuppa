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

package rubygems

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/util"
	"regexp"
	"strings"
	"time"
)

const (
	// LatestAPI is the string format for the Rubygems latest API
	LatestAPI = "https://rubygems.org/api/v1/versions/%s/latest.json"
	// VersionsAPI is the string format for the Rubygems versions API
	VersionsAPI = "https://rubygems.org/api/v1/versions/%s.json"
	// SourceFormat is the string format for Gem sources
	SourceFormat = "https://rubygems.org/downloads/%s-%s.gem"
)

// GemRegex matches Rubygems sources
var GemRegex = regexp.MustCompile("https?://rubygems.org/downloads/(.+).gem")

// Provider is the upstream provider interface for rubygems
type Provider struct{}

// String gives the name of this provider
func (c Provider) String() string {
	return "Rubygems"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) (params []string) {
	if sm := GemRegex.FindStringSubmatch(query); len(sm) > 1 {
		pieces := strings.Split(sm[1], "-")
		if len(pieces) > 2 {
			params = append(params, strings.Join(pieces[0:len(pieces)-1], "-"))
			return
		}
		params = pieces[:1]
	}
	return
}

// Latest finds the newest release for a rubygems package
func (c Provider) Latest(params []string) (r *results.Result, err error) {
	name := params[0]
	url := fmt.Sprintf(LatestAPI, name)
	var cr LatestVersion
	if err = util.FetchJSON(url, "latest", &cr); err == nil {
		r = cr.Convert(name)
	}
	time.Sleep(time.Second)
	return
}

// Releases finds all matching releases for a rubygems package
func (c Provider) Releases(params []string) (rs *results.ResultSet, err error) {
	name := params[0]
	url := fmt.Sprintf(VersionsAPI, name)
	var crs Versions
	if err = util.FetchJSON(url, "releases", &crs); err != nil {
		return
	}
	if len(crs) == 0 {
		err = results.NotFound
		return
	}
	rs = crs.Convert(name)
	return
}
