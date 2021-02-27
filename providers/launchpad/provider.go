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

package launchpad

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/util"
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

// String gives the name of this provider
func (c Provider) String() string {
	return "Launchpad"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) (params []string) {
	if sm := SourceRegex.FindStringSubmatch(query); len(sm) > 1 {
		params = sm[1:]
	}
	return
}

// Latest finds the newest release for a launchpad package
func (c Provider) Latest(params []string) (r *results.Result, err error) {
	rs, err := c.Releases(params)
	if err == nil {
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a launchpad package
func (c Provider) Releases(params []string) (rs *results.ResultSet, err error) {
	name := params[0]
	// Query the API
	url := fmt.Sprintf(SeriesAPI, name)
	var seriesList SeriesList
	if err = util.FetchJSON(url, "series", &seriesList); err != nil {
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
		url := fmt.Sprintf(ReleasesAPI, name, s.Name)
		var vl VersionList
		if err = util.FetchJSON(url, "releases", &vl); err != nil {
			continue
		}
		for i := len(vl.Versions) - 1; i >= 0; i-- {
			r := vl.Versions[i]
			url := fmt.Sprintf(FilesAPI, name, s.Name, r.Number)
			var fl FileList
			if err = util.FetchJSON(url, "files", &fl); err != nil {
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
