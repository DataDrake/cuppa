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

package html

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"os"
)

// Provider is the upstream provider interface for HTML
type Provider struct{}

// Latest finds the newest release for a GNOME package
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
	for _, upstream := range upstreams {
		name := upstream.Match(query)
		if len(name) != 0 {
			return name
		}
	}
	return ""
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "HTML"
}

// Releases finds all matching releases for a rubygems package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	var upstream Upstream
	for i := range upstreams {
		if len(upstreams[i].Match(name)) != 0 {
			upstream = upstreams[i]
			break
		}
	}
	sm := upstream.HostPattern.FindStringSubmatch(name)
	resp, err := http.Get(sm[1])
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
	rs, err = upstream.Parse(name, resp.Body)
	if err != nil {
		s = results.NotFound
	}
	return
}
