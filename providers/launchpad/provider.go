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

package launchpad

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
	"net/http"
	"regexp"
)

const (
	// FilesAPI is the format string for the Launchpad Files API
	FilesAPI = "https://api.launchpad.net/1.0/%s/%s/%s/files"
	// ReleasesAPI is the format string for the Launchpad Releases API
	ReleasesAPI = "https://api.launchpad.net/1.0/%s/%s/releases"
	// SeriesAPI is the format string for the Launchpad Series API
	SeriesAPI = "https://api.launchpad.net/1.0/%s/series"
	// SourceFormat is the format string for Launchpad tarballs
	SourceFormat = "https://launchpad.net/%s/%s/%s/+download/%s-%s.tar.gz"
)

// SourceRegex matches Launchpad source tarballs
var SourceRegex = regexp.MustCompile("https?://launchpad.net/(.*)/.*/.*/\\+download/.*.tar.gz")

// Provider is the upstream provider interface for launchpad
type Provider struct{}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "Launchpad"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	if sm := SourceRegex.FindStringSubmatch(query); len(sm) > 1 {
		return sm[1]
	}
	return ""
}

// Latest finds the newest release for a launchpad package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	rs, err := c.Releases(name)
	if err == nil {
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a launchpad package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
	// Query the API
	r, err := http.NewRequest("GET", fmt.Sprintf(SeriesAPI, name), nil)
	if err != nil {
		log.Debugf("Failed to build request: %s\n", err)
		err = results.Unavailable
		return
	}
	r.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Debugf("Failed to get series: %s\n", err)
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
	var seriesList SeriesList
	if err = dec.Decode(&seriesList); err != nil {
		log.Debugf("Failed to decode series: %s\n", err)
		err = results.Unavailable
		return
	}
	// Proccess Releases
	var lrs Releases
	for _, s := range seriesList.Entries {
		// Only Active Series
		if !s.Active {
			continue
		}
		// Only stable or supported
		switch s.Status {
		case "Active Development":
		case "Current Stable Release":
		case "Supported":
		default:
			continue
		}
		r, err := http.Get(fmt.Sprintf(ReleasesAPI, name, s.Name))
		if err != nil {
			log.Debugf("Failed to get releases: %s\n", err)
			continue
		}
		dec := json.NewDecoder(r.Body)
		var vl VersionList
		if err = dec.Decode(&vl); err != nil {
			log.Debugf("Failed to decode releases: %s\n", err)
			continue
		}
		for i := len(vl.Versions) - 1; i >= 0; i-- {
			r := vl.Versions[i]
			resp, err := http.Get(fmt.Sprintf(FilesAPI, name, s.Name, r.Number))
			if err != nil {
				log.Debugf("Failed to get files: %s\n", err)
				continue
			}
			dec := json.NewDecoder(resp.Body)
			var fl FileList
			if err = dec.Decode(&fl); err != nil {
				log.Debugf("Failed to decode files: %s\n", err)
				continue
			}
			var lr Release
			for _, f := range fl.Files {
				if f.Type != "Code Release Tarball" {
					continue
				}
				lr.name = name
				lr.series = s.Name
				lr.release = r.Number
				lr.uploaded = f.Uploaded
			}
			lrs = append(lrs, lr)
		}
	}
	if len(lrs) == 0 {
		err = results.NotFound
		return
	}
	rs = lrs.Convert(name)
	err = nil
	return
}
