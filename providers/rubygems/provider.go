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

package rubygems

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	// LatestAPI is the string format for the Rubygems latest API
	LatestAPI = "https://rubygems.org/api/v1/versions/%s/latest.json"
	// VersionsAPI is the string format for the Rubygems versions API
	VersionsAPI = "https://rubygems.org/api/v1/versions/%s.json"
	// SourceFormat is the string format for Gem sources
	SourceFormat = "https://rubygems.org/downloads/%s-%s.gem"
)

// GemRegex matches Rubygems sources
var GemRegex = regexp.MustCompile("https?://rubygems.org/downloads/(.+).gem")

// Provider is the upstream provider interface for rubygems
type Provider struct{}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "Rubygems"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	if sm := GemRegex.FindStringSubmatch(query); len(sm) > 1 {
		pieces := strings.Split(sm[1], "-")
		if len(pieces) > 2 {
			return strings.Join(pieces[0:len(pieces)-1], "-")
		}
		return pieces[0]
	}
	return ""
}

// Latest finds the newest release for a rubygems package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(LatestAPI, name))
	if err != nil {
		log.Debugf("Failed to get latest: %s\n", err)
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
	// Decode repsone
	dec := json.NewDecoder(resp.Body)
	var cr LatestVersion
	if err = dec.Decode(&cr); err != nil {
		log.Debugf("Failed to decode latest: %s\n", err)
		err = results.Unavailable
		return
	}
	r = cr.Convert(name)
	time.Sleep(time.Second)
	return
}

// Releases finds all matching releases for a rubygems package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
	// Query the API
	resp, err := http.Get(fmt.Sprintf(VersionsAPI, name))
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

	dec := json.NewDecoder(resp.Body)
	var crs Versions
	if err = dec.Decode(&crs); err != nil {
		log.Debugf("Failed to decode releases: %s\n", err)
		err = results.Unavailable
		return
	}
	if len(crs) == 0 {
		err = results.NotFound
		return
	}
	rs = crs.Convert(name)
	err = nil
	return
}
