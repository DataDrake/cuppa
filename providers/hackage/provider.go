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

package hackage

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	// TarballAPI is the format string for the Hackage Tarball API
	TarballAPI = "https://hackage.haskell.org/package/%s-%s/%s-%s.tar.gz"
	// UploadTimeAPI is the format string for the Hackage Upload Time API
	UploadTimeAPI = "https://hackage.haskell.org/package/%s-%s/upload-time"
	// VersionsAPI is the format string for the Hackage Versions API
	VersionsAPI = "https://hackage.haskell.org/package/%s/preferred"
)

// TarballRegex matches HAckage tarballs
var TarballRegex = regexp.MustCompile("https?://hackage.haskell.org/package/.*/(.*)-(.*?).tar.gz")

// Provider is the upstream provider interface for hackage
type Provider struct{}

// Latest finds the newest release for a hackage package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	rs, s := c.Releases(name)
	// Fail if not OK
	if s != results.OK {
		return
	}
	r = rs.First()
	return
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := TarballRegex.FindStringSubmatch(query)
	if len(sm) == 0 {
		return ""
	}
	return sm[1]
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "Hackage"
}

// Releases finds all matching releases for a hackage package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {

	// Query the API
	r, err := http.NewRequest("GET", fmt.Sprintf(VersionsAPI, name), nil)
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
	versions := &Versions{}
	err = dec.Decode(versions)
	if err != nil {
		panic(err.Error())
	}

	hrs := &Releases{}
	for _, v := range versions.Normal {
		hr := Release{
			name:    name,
			version: v,
		}
		r, err := http.Get(fmt.Sprintf(UploadTimeAPI, name, v))
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
