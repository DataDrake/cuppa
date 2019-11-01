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

package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/DataDrake/cuppa/results"
)

const (
	// SourceFormat is the format string for GitLab release tarballs
	SourceFormat = "https://gitlab.com/%s/-/archive/%s/%s.tar.bz2"

	// SourceFormatFreedesktop is the format string for release tarballs from Freedesktop's GitLab
	SourceFormatFreedesktop = "https://gitlab.freedesktop.org/%s/-/archive/%s/%s.tar.bz2"

	// GitlabEndpoint is the API endpoint URL for GitLab project tags
	GitlabEndpoint = "https://gitlab.com/api/v4/projects/%s/repository/tags"

	// FreedesktopEndpoint is the API endpoint URL for Freedesktop project tags
	FreedesktopEndpoint = "https://gitlab.freedesktop.org/api/v4/projects/%s/repository/tags"
)

// SourceRegex is the regex for GitLab sources. It also matches Freedesktop's GitLab format
var SourceRegex = regexp.MustCompile("https?://.*(?:gitlab.com|gitlab.freedesktop.org)/([^/]+/[^/.]+)")

// VersionRegex is used to parse GitLab version numbers
var VersionRegex = regexp.MustCompile("(?:\\d+\\.)*\\d+\\w*")

// Provider is the upstream provider interface for GitLab
type Provider struct{}

// Latest finds the newest release for a GitLab package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	rs, s := c.Releases(name)
	if s != results.OK {
		return
	}
	r = rs.Last()
	return
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := SourceRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}

	return strings.TrimPrefix(sm[0], "https://")
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "GitLab"
}

// Releases finds all matching releases for a GitLab package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	// Figure out which endpoint we need to hit
	parts := strings.SplitN(name, "/", 2)
	hostname := parts[0]
	project := parts[1]
	encoded := strings.Replace(project, "/", "%2f", 1)
	isFreedesktop := strings.HasSuffix(hostname, "freedesktop.org")

	var endpoint string
	if isFreedesktop {
		endpoint = fmt.Sprintf(FreedesktopEndpoint, encoded)
	} else {
		endpoint = fmt.Sprintf(GitlabEndpoint, encoded)
	}

	// Query the API
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
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
	keys := make([]Tag, 0)
	err = dec.Decode(&keys)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
	}

	tags := &Tags{keys}
	rs = tags.Convert(project, isFreedesktop)
	return
}
