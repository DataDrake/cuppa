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

package html

import (
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
	"net/http"
)

// Provider is the upstream provider interface for HTML
type Provider struct{}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "HTML"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	for _, upstream := range upstreams {
		if name := upstream.Match(query); len(name) > 0 {
			return name
		}
	}
	return ""
}

// Latest finds the newest release for a GNOME package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	rs, err := c.Releases(name)
	if err == nil {
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a rubygems package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
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
	if rs, err = upstream.Parse(name, resp.Body); err != nil {
		err = results.NotFound
	}
	return
}
