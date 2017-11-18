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

package launchpad

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
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

// Latest finds the newest release for a launchpad package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	rs, s := c.Releases(name)
	// Fail if not OK
	if s != results.OK {
		return
	}
	r = rs.Last()
	return
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := SourceRegex.FindStringSubmatch(query)
	if len(sm) == 0 {
		return ""
	}
	return sm[1]
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "Launchpad"
}

// Releases finds all matching releases for a launchpad package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {

	// Query the API
	r, err := http.NewRequest("GET", fmt.Sprintf(SeriesAPI, name), nil)
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
	seriesList := &SeriesList{}
	err = dec.Decode(seriesList)
	if err != nil {
		panic(err.Error())
	}

	lrs := make(Releases, 0)
	for _, s := range seriesList.Entries {
		// Only Active Series
		if !s.Active {
			continue
		}

		// Only stable or supported
		switch s.Status {
		case "Current Stable Release":
		case "Supported":
		default:
			continue
		}
		r, err := http.Get(fmt.Sprintf(ReleasesAPI, name, s.Name))
		if err != nil {
			panic(err.Error())
		}
		dec := json.NewDecoder(r.Body)
		vl := &VersionList{}
		err = dec.Decode(vl)
		if err != nil {
			panic(err.Error())
		}
		for _, r := range vl.Versions {
			resp, err := http.Get(fmt.Sprintf(FilesAPI, name, s.Name, r.Number))
			if err != nil {
				panic(err.Error())
			}
			dec := json.NewDecoder(resp.Body)
			fl := &FileList{}
			err = dec.Decode(fl)
			if err != nil {
				panic(err.Error())
			}
			lr := Release{}
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
		s = results.NotFound
		return
	}
	rs = lrs.Convert(name)
	return
}
