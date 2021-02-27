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
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/util"
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

// String gives the name of this provider
func (c Provider) String() string {
	return "Hackage"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) (params []string) {
	if sm := TarballRegex.FindStringSubmatch(query); len(sm) > 1 {
		params = sm[1:]
	}
	return
}

// Latest finds the newest release for a hackage package
func (c Provider) Latest(params []string) (r *results.Result, err error) {
	rs, err := c.Releases(params)
	if err == nil {
		r = rs.First()
	}
	return
}

// Releases finds all matching releases for a hackage package
func (c Provider) Releases(params []string) (rs *results.ResultSet, err error) {
	name := params[0]
	url := fmt.Sprintf(VersionsAPI, name)
	var versions Versions
	if err = util.FetchJSON(url, "versions", &versions); err != nil {
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
