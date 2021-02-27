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
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
	"net/http"
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

// Name gives the name of this provider
func (c Provider) Name() string {
	return "GitLab"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	if sm := SourceRegex.FindStringSubmatch(query); len(sm) > 2 {
		return sm[0]
	}
	return ""
}

// Latest finds the newest release for a GitLab package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	rs, err := c.Releases(name)
	if err == nil {
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a GitLab package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
	// Query the API
	sm := SourceRegex.FindStringSubmatch(name)
	id := strings.Join(strings.Split(sm[2], "/"), "%2f")
	resp, err := http.Get(fmt.Sprintf(TagsEndpoint, sm[1], id))
	if err != nil {
		log.Debugf("Failed to get releases: %s\n", err)
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
	var tags Tags
	if err = dec.Decode(&tags); err != nil {
		log.Debugf("Failed to decode response: %s\n", err)
		err = results.Unavailable
		return
	}
	rs = tags.Convert(sm[1], sm[2])
	return
}
