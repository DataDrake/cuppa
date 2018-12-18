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

package kde

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	// ListingURL is the location of the KDE FTP file listing
	ListingURL = "https://download.kde.org/ls-lR.bz2"
	// ListingPrefix is the prefix of all paths in the KDE listing that is hidden by HTTP
	ListingPrefix = "/srv/archives/ftp/"
	// SourceFormat4 is the string format for KDE sources with 4 pieces
	SourceFormat4 = "https://download.kde.org/%s/%s/%s/%s-%s.tar.bz2"
	// SourceFormat5 is the string format for KDE sources with 5 pieces
	SourceFormat5 = "https://download.kde.org/%s/%s/%s/%s/%s-%s.tar.bz2"
)

// TarballRegex matches KDE sources
var TarballRegex = regexp.MustCompile("https?://download.kde.org/(.+)")

// Provider is the upstream provider interface for KDE
type Provider struct{}

// Latest finds the newest release for a KDE package
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
	sm := TarballRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}
	pieces := strings.Split(sm[1], "/")
	if len(pieces) < 4 || len(pieces) > 5 {
		return ""
	}
	return sm[1]
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "KDE"
}

// Releases finds all matching releases for a KDE package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	// Query the API
	resp, err := http.Get(ListingURL)
	if err != nil {
		panic(err.Error())
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
	pieces := strings.Split(name, "/")
	name = strings.Split(pieces[len(pieces)-1], "-")[0]
	var searchPrefix string
	switch len(pieces) {
	case 4:
		searchPrefix = ListingPrefix + strings.Join(pieces[0:len(pieces)-2], "/") + ":\n"
	case 5:
		searchPrefix = ListingPrefix + strings.Join(pieces[0:len(pieces)-3], "/") + ":\n"
	}
	body := bzip2.NewReader(resp.Body)
	listing := bufio.NewReader(body)
	rs = results.NewResultSet(name)
	line := ""
	for {
		line, err = listing.ReadString('\n')
		if err != nil {
			break
		}
		if line != searchPrefix {
			continue
		}
		for line != "\n" {
			line, err = listing.ReadString('\n')
			if err != nil {
				break
			}
			if line == "\n" {
				break
			}
			fields := strings.Fields(line)
			version := fields[len(fields)-1]
			updated, _ := time.Parse("2006-01-02 15:04", strings.Join(fields[len(fields)-3:len(fields)-2], " "))
			var location string
			switch len(pieces) {
			case 4:
				location = fmt.Sprintf(SourceFormat4, pieces[0], pieces[1], version, name, version)
			case 5:
				location = fmt.Sprintf(SourceFormat5, pieces[0], pieces[1], version, pieces[3], name, version)
			}
			r := &results.Result{
				Name:      name,
				Version:   version,
				Location:  location,
				Published: updated,
			}
			rs.AddResult(r)
		}
		break
	}
	return
}
