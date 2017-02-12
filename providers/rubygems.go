//
// Copyright © 2016 Bryan T. Meyers <bmeyers@datadrake.com>
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

package providers

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"regexp"
	"time"
)

var rubygemsAPILatest = "https://rubygems.org/api/v1/versions/%s/latest.json"
var rubygemsAPIVersions = "https://rubygems.org/api/v1/versions/%s.json"
var rubygemsSource = "https://rubygems.org/downloads/%s-%s.gem"
var rubygemsRegex = regexp.MustCompile("https?://rubygems.org/downloads/(.*)-.*\\.gem")

type rubygemsLatest struct {
	Version string `json:"version"`
}

func (cr *rubygemsLatest) Convert(name string) *results.Result {
	r := &results.Result{}
	r.Name = name
	r.Version = cr.Version
	r.Location = fmt.Sprintf(rubygemsSource, name, cr.Version)
	return r
}

type rubygemsVersion struct {
	CreatedAt  string `json:"created_at"`
	PreRelease bool   `json:"prerelease"`
	Number     string `json:"number"`
}

func (cr *rubygemsVersion) Convert(name string) *results.Result {
	if cr.PreRelease {
		return nil
	}
	r := &results.Result{}
	r.Name = name
	r.Version = cr.Number
	r.Published, _ = time.Parse(time.RFC3339, cr.CreatedAt)
	r.Location = fmt.Sprintf(rubygemsSource, name, cr.Number)
	return r
}

type rubygemsResultSet []rubygemsVersion

func (crs *rubygemsResultSet) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, rel := range *crs {
		r := rel.Convert(name)
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}

/*
RubygemsProvider is the upstream provider interface for rubygems
*/
type RubygemsProvider struct{}

/*
Latest finds the newest release for a rubygems package
*/
func (c RubygemsProvider) Latest(name string) (r *results.Result, s results.Status) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(rubygemsAPILatest, name))
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	//Translate Status Code
	switch resp.StatusCode {
	case 200:
		s = results.OK
	case 404:
		s = results.NotFound
	default:
		s = results.Unavailable
	}

	//Fail if not OK
	if s != results.OK {
		return
	}

	dec := json.NewDecoder(resp.Body)
	cr := &rubygemsLatest{}
	err = dec.Decode(cr)
	if err != nil {
		panic(err.Error())
	}
	r = cr.Convert(name)
	return
}

/*
Match checks to see if this provider can handle this kind of query
*/
func (c RubygemsProvider) Match(query string) string {
	sm := rubygemsRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}
	return sm[1]
}

/*
Name gives the name of this provider
*/
func (c RubygemsProvider) Name() string {
	return "Rubygems"
}

/*
Releases finds all matching releases for a rubygems package
*/
func (c RubygemsProvider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(rubygemsAPIVersions, name))
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	//Translate Status Code
	switch resp.StatusCode {
	case 200:
		s = results.OK
	case 404:
		s = results.NotFound
	default:
		s = results.Unavailable
	}

	//Fail if not OK
	if s != results.OK {
		return
	}

	dec := json.NewDecoder(resp.Body)
	crs := make(rubygemsResultSet, 0)
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
