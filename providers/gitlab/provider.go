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

package gitlab

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/util"
	"regexp"
	"strings"
)

const (
	// SourceFormat is the format string for GitLab release tarballs
	SourceFormat = "https://%s/%s/-/archive/%s/%s.tar.gz"

	// TagsEndpoint is the API endpoint URL for GitLab project tags
	TagsEndpoint = "https://%s/api/v4/projects/%s/repository/tags"
)

var (
	// SourceRegex is the regex for GitLab sources
	SourceRegex = regexp.MustCompile("https?://(gitlab[^/]+)/(.+/[^/.]+)/\\-/")
	// VersionRegex is used to parse GitLab version numbers
	VersionRegex = regexp.MustCompile("(?:\\d+\\.)*\\d+\\w*")
)

// Provider is the upstream provider interface for GitLab
type Provider struct{}

// String gives the name of this provider
func (c Provider) String() string {
	return "GitLab"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) (params []string) {
	if sm := SourceRegex.FindStringSubmatch(query); len(sm) > 2 {
		params = sm[1:]
	}
	return
}

// Latest finds the newest release for a GitLab package
func (c Provider) Latest(params []string) (r *results.Result, err error) {
	rs, err := c.Releases(params)
	if err == nil {
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a GitLab package
func (c Provider) Releases(params []string) (rs *results.ResultSet, err error) {
	// Query the API
	id := strings.Join(strings.Split(params[1], "/"), "%2f")
	url := fmt.Sprintf(TagsEndpoint, params[0], id)
	var tags Tags
	if err = util.FetchJSON(url, "releases", &tags); err != nil {
		return
	}
	rs = tags.Convert(params[0], params[1])
	return
}
