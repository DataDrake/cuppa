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

package github

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"regexp"
)

const (
	// APILatest is a format string for the Latest releases API
	APILatest = "https://api.github.com/repos/%s/releases/latest"
	// APIReleases is a format string for the Releases API
	APIReleases = "https://api.github.com/repos/%s/releases"
	// APITags is a format string for the Tags API
	APITags = "https://api.github.com/repos/%s/git/refs/tags"
	// SourceFormat is the format string for Github release tarballs
	SourceFormat = "https://github.com/%s/archive/%s.tar.gz"
)

// SourceRegex is the regex for Github sources
var SourceRegex = regexp.MustCompile("https?://github.com/([^/]*/[^/]*)/.*/[^/]*.tar.gz")

// VersionRegex is used to parse Github version numbers
var VersionRegex = regexp.MustCompile("(?:\\d+\\.)*\\d+\\w*")

// Provider is the upstream provider interface for github
type Provider struct{}

// Latest finds the newest release for a github package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	// Query the API
	resp, err := http.Get(fmt.Sprintf(APILatest, name))
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	// Translate Status Code
	switch resp.StatusCode {
	case 200:
		s = results.OK
	case 404:
		rs, st := getTags(name)
		s = st
		if s != results.OK || rs.Empty() {
			return
		}
		r = rs.Last()
		return
	default:
		s = results.Unavailable
	}

	// Fail if not OK
	if s != results.OK {
		return
	}

	dec := json.NewDecoder(resp.Body)
	cr := &Release{}
	err = dec.Decode(cr)
	if err != nil {
		panic(err.Error())
	}
	r = cr.Convert(name)
	return
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := SourceRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}
	return sm[1]
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "GitHub"
}

func getTags(name string) (rs *results.ResultSet, s results.Status) {
	// Query the API
	resp, err := http.Get(fmt.Sprintf(APITags, name))
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
	tags := make(Tags, 0)
	err = dec.Decode(&tags)
	if err != nil {
		panic(err.Error())
	}
	if len(tags) == 0 {
		s = results.NotFound
		return
	}
	rs = tags.Convert(name)
	return
}

// Releases finds all matching releases for a github package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	// Query the API
	resp, err := http.Get(fmt.Sprintf(APIReleases, name))
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
	crs := make(Releases, 0)
	err = dec.Decode(&crs)
	if err != nil {
		panic(err.Error())
	}
	if len(crs) == 0 {
		rs, s = getTags(name)
		return
	}
	rs = crs.Convert(name)
	return
}
