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

package providers

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var hackageAPITarball = "https://hackage.haskell.org/package/%s-%s/%s-%s.tar.gz"
var hackageAPIUploadTime = "https://hackage.haskell.org/package/%s-%s/upload-time"
var hackageAPIVersions = "https://hackage.haskell.org/package/%s/preferred"
var hackageRegex = regexp.MustCompile("https?://hackage.haskell.org/package/(.*)-(.*)/.*.tar.gz")

type hackageVersions struct {
	Normal []string `json:"normal-version"`
}

type hackageRelease struct {
	name     string
	released string
	version  string
}

// Convert turns a Hackage release into a Cuppa result
func (hr *hackageRelease) Convert() *results.Result {
	r := &results.Result{}
	r.Name = hr.name
	r.Version = hr.version
	r.Published, _ = time.Parse(time.UnixDate, hr.released)
	r.Location = fmt.Sprintf(hackageAPITarball, hr.name, hr.version, hr.name, hr.version)
	return r
}

type hackageResultSet struct {
	Releases []hackageRelease
}

// Convert turns a Hackage result set into a Cuppa result set
func (hrs *hackageResultSet) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, rel := range hrs.Releases {
		r := rel.Convert()
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}

// HackageProvider is the upstream provider interface for hackage
type HackageProvider struct{}

// Latest finds the newest release for a hackage package
func (c HackageProvider) Latest(name string) (r *results.Result, s results.Status) {
	rs, s := c.Releases(name)
	// Fail if not OK
	if s != results.OK {
		return
	}
	r = rs.First()
	return
}

// Match checks to see if this provider can handle this kind of query
func (c HackageProvider) Match(query string) string {
	sm := hackageRegex.FindStringSubmatch(query)
	if len(sm) == 0 {
		return ""
	}
	return sm[1]
}

// Name gives the name of this provider
func (c HackageProvider) Name() string {
	return "Hackage"
}

// Releases finds all matching releases for a hackage package
func (c HackageProvider) Releases(name string) (rs *results.ResultSet, s results.Status) {

	// Query the API
	r, err := http.NewRequest("GET", fmt.Sprintf(hackageAPIVersions, name), nil)
	if err != nil {
		panic(err.Error())
	}
	r.Header.Set("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(r)
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
	versions := &hackageVersions{}
	err = dec.Decode(versions)
	if err != nil {
		panic(err.Error())
	}

	hrs := &hackageResultSet{}
	for _, v := range versions.Normal {
		hr := hackageRelease{
			name:    name,
			version: v,
		}
		r, err := http.Get(fmt.Sprintf(hackageAPIUploadTime, name, v))
		if err != nil {
			panic(err.Error())
		}
		defer r.Body.Close()
		dateRaw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}
		hr.released = string(dateRaw)
		hrs.Releases = append(hrs.Releases, hr)
	}
	if len(hrs.Releases) == 0 {
		s = results.NotFound
		return
	}
	rs = hrs.Convert(name)
	return
}
