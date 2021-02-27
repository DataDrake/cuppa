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

package hackage

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
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

// Name gives the name of this provider
func (c Provider) Name() string {
	return "Hackage"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	if sm := TarballRegex.FindStringSubmatch(query); len(sm) > 1 {
		return sm[1]
	}
	return ""
}

// Latest finds the newest release for a hackage package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	rs, err := c.Releases(name)
	if err == nil {
		r = rs.First()
	}
	return
}

// Releases finds all matching releases for a hackage package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
	// Query the API
	r, err := http.NewRequest("GET", fmt.Sprintf(VersionsAPI, name), nil)
	if err != nil {
		log.Debugf("Failed to create request: %s\n", err)
		err = results.Unavailable
		return
	}
	r.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Debugf("Failed to get versions: %s\n", err)
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
	var versions Versions
	if err = dec.Decode(&versions); err != nil {
		log.Debugf("Failed to decode versions: %s\n", err)
		err = results.Unavailable
		return
	}
	// Process releases
	var hrs Releases
	for _, v := range versions.Normal {
		hr := Release{
			name:    name,
			version: v,
		}
		r, err := http.Get(fmt.Sprintf(UploadTimeAPI, name, v))
		if err != nil {
			log.Debugf("Failed to get upload time: %s\n", err)
			continue
		}
		defer r.Body.Close()
		dateRaw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Debugf("Failed to read response: %s\n", err)
			continue
		}
		hr.released = string(dateRaw)
		hrs.Releases = append(hrs.Releases, hr)
	}
	if len(hrs.Releases) == 0 {
		err = results.NotFound
		return
	}
	rs = hrs.Convert(name)
	err = nil
	return
}
