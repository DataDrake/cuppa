//
// Copyright 2016-2020 Bryan T. Meyers <root@datadrake.com>
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

package cpan

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
	// APIRelease is the format string for the metacpan release API
	APIRelease = "https://fastapi.metacpan.org/v1/release/%s"
	// APIDownloadURL is the format string for the metacpan download_url API
	APIDownloadURL = "https://fastapi.metacpan.org/v1/download_url/%s"
)

// SearchRegex is the regexp for "search.cpan.org"
var SearchRegex = regexp.MustCompile("https?://*(?:/.*cpan.org)(?:/CPAN)?/authors/id/(.+)$")

// Provider is the upstream provider interface for CPAN
type Provider struct{}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "CPAN"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	if sm := SearchRegex.FindStringSubmatch(query); len(sm) > 0 {
		sms := strings.Split(sm[1], "/")
		filename := sms[len(sms)-1]
		pieces := strings.Split(filename, "-")
		if len(pieces) > 2 {
			return strings.Join(pieces[0:len(pieces)-1], "-")
		}
		return pieces[0]
	}
	return ""
}

// Latest finds the newest release for a CPAN package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	// Query the APIDownloadURL
	module, err := nameToModule(name)
	if err != nil {
		return
	}
	resp, err := http.Get(fmt.Sprintf(APIDownloadURL, module))
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
	var rel Release
	if err = dec.Decode(&rel); err != nil {
		log.Debugf("Failed to decode response: %s\n", err)
		err = results.Unavailable
		return
	}
	if len(rel.Error) > 0 {
		err = results.NotFound
		return
	}
	if r = rel.Convert(name); r == nil {
		err = results.NotFound
	}
	return
}

// Releases finds all matching releases for a CPAN package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
	if r, err := c.Latest(name); err == nil {
		rs = results.NewResultSet(name)
		rs.AddResult(r)
	}
	return
}
