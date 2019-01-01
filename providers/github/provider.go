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
	"github.com/DataDrake/cuppa/results"
	"regexp"
)

const (
	// SourceFormat is the format string for Github release tarballs
	SourceFormat = "https://github.com/%s/archive/%s.tar.gz"
)

// SourceRegex is the regex for Github sources
var SourceRegex = regexp.MustCompile("github.com/([^/]+/[^/.]+)")

// VersionRegex is used to parse Github version numbers
var VersionRegex = regexp.MustCompile("(?:\\d+\\.)*\\d+\\w*")

// Provider is the upstream provider interface for github
type Provider struct{}

// Latest finds the newest release for a github package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	rs, s := c.GetReleases(name, 100)
	if s != results.OK {
		return
	}
	r = rs.Last()
	if r == nil {
		s = results.NotFound
	}
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

// Releases finds all matching releases for a github package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	rs, s = c.GetReleases(name, 100)
	return
}
