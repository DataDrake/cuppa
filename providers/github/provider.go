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

package github

import (
	"github.com/DataDrake/cuppa/results"
	"regexp"
)

const (
	// SourceFormat is the format string for Github release tarballs
	SourceFormat = "https://github.com/%s/archive/%s.tar.gz"
)

var (
	// SourceRegex is the regex for Github sources
	SourceRegex = regexp.MustCompile("github.com/([^/]+/[^/.]+)")
	// VersionRegex is used to parse Github version numbers
	VersionRegex = regexp.MustCompile("(?:\\d+\\.)*\\d+\\w*")
)

// Provider is the upstream provider interface for github
type Provider struct{}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "GitHub"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	if sm := SourceRegex.FindStringSubmatch(query); len(sm) > 1 {
		return sm[1]
	}
	return ""
}

// Latest finds the newest release for a github package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	rs, err := c.GetReleases(name, 100)
	if err == nil {
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a github package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
	rs, err = c.GetReleases(name, 100)
	return
}
