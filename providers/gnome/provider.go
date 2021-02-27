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

package gnome

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// CacheAPI is the string format for GNOME cache.json files
	CacheAPI = "https://download.gnome.org/sources/%s/cache.json"
	// SourceFormat is the string format for GNOME sources
	SourceFormat = "https://download.gnome.org/sources/%s/%s"
)

// TarballRegex matches GNOME sources
var TarballRegex = regexp.MustCompile("https?://(?:ftp.gnome.org/pub/gnome|download.gnome.org)/sources/(.+?)/.*")

// Provider is the upstream provider interface for GNOME
type Provider struct{}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "GNOME"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	if sm := TarballRegex.FindStringSubmatch(query); len(sm) > 1 {
		return sm[1]
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
	// Query the API
	resp, err := http.Get(fmt.Sprintf(CacheAPI, name))
	if err != nil {
		log.Debugf("Failed to fetch releases: %s\n", err)
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
	var raw []interface{}
	if err = dec.Decode(&raw); err != nil {
		log.Debugf("Failed to decode response: %s\n", err)
		err = results.Unavailable
		return
	}
	if len(raw) < 3 {
		err = results.Unavailable
		return
	}
	rs = Merge(name, raw[1].(map[string]interface{}), raw[2].(map[string]interface{}))
	return
}

// Merge combines Source and Versions into a ResultSet
func Merge(name string, srcs, vs map[string]interface{}) (rs *results.ResultSet) {
	rs = results.NewResultSet(name)
	if srcs[name] == nil || vs[name] == nil {
		return
	}
	for _, v := range vs[name].([]interface{}) {
		pieces := strings.Split(v.(string), ".")
		if len(pieces) < 2 {
			continue
		}
		minor, err := strconv.Atoi(pieces[1])
		if err != nil {
			continue
		}
		// Filter out unstable releases
		if minor%2 == 1 {
			continue
		}
		files := srcs[name].(map[string]interface{})[v.(string)].(map[string]interface{})
		if len(files) == 0 {
			continue
		}
		// get location of tarball
		var location string
		switch {
		case files["tar.xz"] != nil:
			location = fmt.Sprintf(SourceFormat, name, files["tar.xz"].(string))
		case files["tar.gz"] != nil:
			location = fmt.Sprintf(SourceFormat, name, files["tar.gz"].(string))
		case files["tar.bz2"] != nil:
			location = fmt.Sprintf(SourceFormat, name, files["tar.bz2"].(string))
		default:
			continue
		}
		r := results.NewResult(name, v.(string), location, time.Time{})
		rs.AddResult(r)
	}
	return
}
