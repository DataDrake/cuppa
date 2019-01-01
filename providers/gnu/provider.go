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

package gnu

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"github.com/jlaffaye/ftp"
	"regexp"
	"sort"
)

const (
	// MirrorsFTP is the host to use as a GNU mirror
	MirrorsFTP = "mirrors.rit.edu:21"
	// GNUFormat is the format string for GNU sources
	GNUFormat = "https://mirrors.rit.edu/gnu/%s/%s"
)

// MirrorsRegex is a regex for a GNU mirror source
var MirrorsRegex = regexp.MustCompile("(?:https?|ftp)://[^\\/]+/gnu/(.+)/[^\\/]+$")

// TarballRegex is a regex for finding tarball files
var TarballRegex = regexp.MustCompile("^(.+)-(.+)\\.tar\\..+z$")

// Provider is the upstream provider interface for GNU
type Provider struct{}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := MirrorsRegex.FindStringSubmatch(query)
	if len(sm) == 0 {
		return ""
	}
	return sm[1]
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "GNU"
}

// Latest finds the newest release for a GNU package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	rs, s := c.Releases(name)
	if s == results.OK {
		sort.Sort(rs)
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a GNU package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	client, err := ftp.Dial(MirrorsFTP)
	if err != nil {
		s = results.Unavailable
		return
	}
	err = client.Login("anonymous", "anonymous")
	if err != nil {
		s = results.Unavailable
		return
	}
	entries, err := client.List("gnu" + "/" + name)
	if err != nil {
		fmt.Printf("FTP Error: %s\n", err.Error())
		s = results.NotFound
		return
	}
	rs = results.NewResultSet(name)
	for _, entry := range entries {
		if entry.Type != ftp.EntryTypeFile {
			continue
		}
		sm := TarballRegex.FindStringSubmatch(entry.Name)
		if len(sm) == 0 {
			continue
		}
		r := results.NewResult(sm[1], sm[2], fmt.Sprintf(GNUFormat, name, entry.Name), entry.Time )
		rs.AddResult(r)
		s = results.OK
	}
	return
}
