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

package rubygems

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"regexp"
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
var GemRegex = regexp.MustCompile("https?://rubygems.org/downloads/(.*)-.*\\.gem")

// Provider is the upstream provider interface for rubygems
type Provider struct{}

// Latest finds the newest release for a rubygems package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(LatestAPI, name))
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
	cr := &LatestVersion{}
	err = dec.Decode(cr)
	if err != nil {
		panic(err.Error())
	}
	r = cr.Convert(name)
	return
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := GemRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}
	return sm[1]
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "Rubygems"
}

// Releases finds all matching releases for a rubygems package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	// Query the API
	resp, err := http.Get(fmt.Sprintf(VersionsAPI, name))
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
	crs := make(Versions, 0)
	err = dec.Decode(&crs)
	if err != nil {
		panic(err.Error())
	}
	if len(crs) == 0 {
		s = results.NotFound
		return
	}
	rs = crs.Convert(name)
	return
}
