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

package gnu

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
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

var (
	// MirrorsRegex is a regex for a GNU mirror source
	MirrorsRegex = regexp.MustCompile("(?:https?|ftp)://[^\\/]+/gnu/(.+)/[^\\/]+$")
	// TarballRegex is a regex for finding tarball files
	TarballRegex = regexp.MustCompile("^(.+)-(.+)\\.tar\\..+z$")
)

// Provider is the upstream provider interface for GNU
type Provider struct{}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "GNU"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	if sm := MirrorsRegex.FindStringSubmatch(query); len(sm) > 1 {
		return sm[1]
	}
	return ""
}

// Latest finds the newest release for a GNU package
func (c Provider) Latest(name string) (r *results.Result, err error) {
	rs, err := c.Releases(name)
	if err == nil {
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a GNU package
func (c Provider) Releases(name string) (rs *results.ResultSet, err error) {
	client, err := ftp.Dial(MirrorsFTP)
	if err != nil {
		log.Debugf("Failed to connect to FTP server: %s\n", err)
		err = results.Unavailable
		return
	}
	if err = client.Login("anonymous", "anonymous"); err != nil {
		log.Debugf("Failed to login to FTP server: %s\n", err)
		err = results.Unavailable
		return
	}
	entries, err := client.List("gnu" + "/" + name)
	if err != nil {
		log.Debugf("FTP Error: %s\n", err.Error())
		err = results.NotFound
		return
	}
	rs = results.NewResultSet(name)
	for _, entry := range entries {
		if entry.Type != ftp.EntryTypeFile {
			continue
		}
		if sm := TarballRegex.FindStringSubmatch(entry.Name); len(sm) > 2 {
			r := results.NewResult(sm[1], sm[2], fmt.Sprintf(GNUFormat, name, entry.Name), entry.Time)
			rs.AddResult(r)
		}
	}
	if rs.Len() == 0 {
		err = results.NotFound
	}
	sort.Sort(rs)
	return
}
