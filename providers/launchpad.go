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
	"net/http"
	"regexp"
	"time"
)

var launchpadAPIFiles = "https://api.launchpad.net/1.0/%s/%s/%s/files"
var launchpadAPIRelease = "https://api.launchpad.net/1.0/%s/%s/releases"
var launchpadAPISeries = "https://api.launchpad.net/1.0/%s/series"
var launchpadRegex = regexp.MustCompile("https://launchpad.net/(.*)/.*/.*/\\+download/.*.tar.gz")
var launchpadSource = "https://launchpad.net/%s/%s/%s/+download/%s-%s.tar.gz"

type launchpadSeriesList struct {
	Entries []launchpadSeries `json:"entries"`
}

type launchpadSeries struct {
	Active bool   `json:"active"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type launchpadReleaseList struct {
	Versions []launchpadVersion `json:"entries"`
}

type launchpadVersion struct {
	Number string `json:"version"`
}

type launchpadFileList struct {
	Files []launchpadFile `json:"entries"`
}

type launchpadFile struct {
	Link     string `json:"file_link"`
	Type     string `json:"file_type"`
	Uploaded string `json:"date_uploaded"`
}

type launchpadRelease struct {
	name     string
	release  string
	series   string
	uploaded string
}

// Convert turns a Launchpad release into a Cuppa result
func (lr *launchpadRelease) Convert() *results.Result {
	r := &results.Result{}
	r.Name = lr.name
	r.Version = lr.release
	r.Published, _ = time.Parse(time.RFC3339, lr.uploaded)
	r.Location = fmt.Sprintf(launchpadSource, lr.name, lr.series, lr.release, lr.name, lr.release)
	return r
}

type launchpadResultSet struct {
	Releases []launchpadRelease
}

// Convert turns a Launchpad result set to a Cuppa result set
func (lrs *launchpadResultSet) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, rel := range lrs.Releases {
		r := rel.Convert()
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}

// LaunchpadProvider is the upstream provider interface for launchpad
type LaunchpadProvider struct{}

// Latest finds the newest release for a launchpad package
func (c LaunchpadProvider) Latest(name string) (r *results.Result, s results.Status) {
	rs, s := c.Releases(name)
	// Fail if not OK
	if s != results.OK {
		return
	}
	r = rs.Last()
	return
}

// Match checks to see if this provider can handle this kind of query
func (c LaunchpadProvider) Match(query string) string {
	sm := launchpadRegex.FindStringSubmatch(query)
	if len(sm) == 0 {
		return ""
	}
	return sm[1]
}

// Name gives the name of this provider
func (c LaunchpadProvider) Name() string {
	return "Launchpad"
}

// Releases finds all matching releases for a launchpad package
func (c LaunchpadProvider) Releases(name string) (rs *results.ResultSet, s results.Status) {

	// Query the API
	r, err := http.NewRequest("GET", fmt.Sprintf(launchpadAPISeries, name), nil)
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
	seriesList := &launchpadSeriesList{}
	err = dec.Decode(seriesList)
	if err != nil {
		panic(err.Error())
	}

	lrs := &launchpadResultSet{}
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
		r, err := http.Get(fmt.Sprintf(launchpadAPIRelease, name, s.Name))
		if err != nil {
			panic(err.Error())
		}
		dec := json.NewDecoder(r.Body)
		releaseList := &launchpadReleaseList{}
		err = dec.Decode(releaseList)
		if err != nil {
			panic(err.Error())
		}
		for _, r := range releaseList.Versions {
			resp, err := http.Get(fmt.Sprintf(launchpadAPIFiles, name, s.Name, r.Number))
			if err != nil {
				panic(err.Error())
			}
			dec := json.NewDecoder(resp.Body)
			fileList := &launchpadFileList{}
			err = dec.Decode(fileList)
			if err != nil {
				panic(err.Error())
			}
			lr := launchpadRelease{}
			for _, f := range fileList.Files {
				if f.Type != "Code Release Tarball" {
					continue
				}
				lr.name = name
				lr.series = s.Name
				lr.release = r.Number
				lr.uploaded = f.Uploaded
			}
			lrs.Releases = append(lrs.Releases, lr)
		}
	}
	if len(lrs.Releases) == 0 {
		s = results.NotFound
		return
	}
	rs = lrs.Convert(name)
	return
}
