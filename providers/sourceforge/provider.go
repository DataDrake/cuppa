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

package sourceforge

import (
	"encoding/xml"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
	"net/http"
	"regexp"
	"time"
)

const (
	// API is the format string for a sourceforge RSS feed
	API = "https://sourceforge.net/projects/%s/rss?path=/%s"
)

var (
	// TarballRegex matches SourceForge sources
	TarballRegex = regexp.MustCompile("https?://.*sourceforge.net/projects?/(.+)/files/(.+/)?(.+?)-([\\d]+(?:.\\d+)*\\w*?)\\.(?:zip|tar\\..+z.*)(?:\\/download)?$")
	// ProjectRegex matches SourceForge sources
	ProjectRegex = regexp.MustCompile("https?://.*sourceforge.net/projects?/(.+)/(?:files/)?(.+?/)?(.+?)-([\\d]+(?:.\\d+)*\\w*?).+$")
)

// Item represents an entry in the RSS Feed
type Item struct {
	XMLName xml.Name `xml:"item"`
	Link    string   `xml:"link"`
	Date    string   `xml:"pubDate"`
}

// Feed represents the RSS feed itself
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Items   []Item   `xml:"channel>item"`
}

// toResults converts a Feed to a ResultSet
func (f *Feed) toResults(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, item := range f.Items {
		if sm := TarballRegex.FindStringSubmatch(item.Link); len(sm) > 4 {
			pub, _ := time.Parse(time.RFC1123, item.Date+"C")
			r := results.NewResult(name, sm[4], item.Link, pub)
			rs.AddResult(r)
		}
	}
	return rs
}

// Provider is the upstream provider interface for SourceForge
type Provider struct{}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "SourceForge"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := TarballRegex.FindStringSubmatch(query)
	if len(sm) != 5 {
		sm = ProjectRegex.FindStringSubmatch(query)
	}
	if len(sm) == 5 {
		return sm[0]
	}
	return ""
}

// Latest finds the newest release for a SourceForge package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	rs, err := c.Releases(name)
	if err == nil {
		r = rs.First()
	}
	return
}

// Releases finds all matching releases for a SourceForge package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
	sm := TarballRegex.FindStringSubmatch(name)
	if len(sm) != 5 {
		sm = ProjectRegex.FindStringSubmatch(name)
		sm[1], sm[3] = sm[3], sm[1]
	}
	// Query the API
	resp, err := http.Get(fmt.Sprintf(API, sm[1], sm[2]))
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
	// decode response
	dec := xml.NewDecoder(resp.Body)
	var feed Feed
	if err = dec.Decode(&feed); err != nil {
		log.Debugf("Failed to decode releases: %s\n", err)
		err = results.Unavailable
		return
	}
	rs = feed.toResults(sm[3])
	if rs.Len() == 0 {
		err = results.NotFound
	}
	return
}
