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

package cpan

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
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
var SearchRegex = regexp.MustCompile("https?://*(?:/.*cpan.org)(?:/CPAN)?/authors/id/(.*)")

// Provider is the upstream provider interface for CPAN
type Provider struct{}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := SearchRegex.FindStringSubmatch(query)
	if len(sm) == 0 {
		return ""
	}
	sms := strings.Split(sm[1], "/")
	filename := sms[len(sms)-1]
	pieces := strings.Split(filename, "-")
	pieces = pieces[0 : len(sms)-2]
    name := pieces[0]
    if len(pieces) > 2 {
        name = strings.Join(pieces, "-")
    }
	return name
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "CPAN"
}


type CPANRelease struct {
    Module string `json:"main_module"`
}

func nameToModule(name string) (module string, s results.Status) {
	// Query the Release API
	resp, err := http.Get(fmt.Sprintf(APIRelease, name))
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

	dec := json.NewDecoder(resp.Body)
	r := &CPANRelease{}
	err = dec.Decode(r)
	if err != nil {
		panic(err.Error())
	}
    module = r.Module
    return
}

// Latest finds the newest release for a CPAN package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	// Query the APIDownloadURL
	module, s := nameToModule(name)
    if s != results.OK {
        return
    }
	resp, err := http.Get(fmt.Sprintf(APIDownloadURL, module))
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

	dec := json.NewDecoder(resp.Body)
	rs := &Release{}
	err = dec.Decode(rs)
	if err != nil {
		panic(err.Error())
	}
	if len(rs.Error) > 0 {
		s = results.NotFound
		return
	}
	r = rs.Convert(name)
	return
}

// Releases finds all matching releases for a CPAN package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	r, s := c.Latest(name)
	if s != results.OK {
		return
	}
	rs = results.NewResultSet(name)
	rs.AddResult(r)
	return
}

/*
// Releases finds all matching releases for a CPAN package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	// Query the APIDownloadURL
	module, s:= nameToModule(name)
	if s != results.OK {
		return
	}
	resp, err := http.Get(fmt.Sprintf(APIDownloadURL, module))
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

	dec := json.NewDecoder(resp.Body)
	crs := &Releases{}
	err = dec.Decode(crs)
	if err != nil {
		panic(err.Error())
	}
	if len(crs.Releases) == 0 {
		s = results.NotFound
		return
	}
	rs = crs.Convert(name)
	return
}
*/
