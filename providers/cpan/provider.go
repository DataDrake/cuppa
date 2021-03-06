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

package cpan

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/util"
	"regexp"
	"strings"
)

// APIDownloadURL is the format string for the metacpan download_url API
const APIDownloadURL = "https://fastapi.metacpan.org/v1/download_url/%s"

// SearchRegex is the regexp for "search.cpan.org"
var SearchRegex = regexp.MustCompile("https?://*(?:/.*cpan.org)(?:/CPAN)?/authors/id/(.+)$")

// Provider is the upstream provider interface for CPAN
type Provider struct{}

// String gives the name of this provider
func (c Provider) String() string {
	return "CPAN"
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) (params []string) {
	if sm := SearchRegex.FindStringSubmatch(query); len(sm) > 0 {
		sms := strings.Split(sm[1], "/")
		filename := sms[len(sms)-1]
		pieces := strings.Split(filename, "-")
		if len(pieces) > 2 {
			params = append(params, strings.Join(pieces[0:len(pieces)-1], "-"))
			return
		}
		params = append(params, pieces[0])
	}
	return
}

// Latest finds the newest release for a CPAN package
func (c Provider) Latest(params []string) (r *results.Result, err error) {
	name := params[0]
	module, err := nameToModule(name)
	if err != nil {
		return
	}
	url := fmt.Sprintf(APIDownloadURL, module)
	var rel Release
	if err = util.FetchJSON(url, "latest", &rel); err != nil {
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
func (c Provider) Releases(params []string) (rs *results.ResultSet, err error) {
	if r, err := c.Latest(params); err == nil {
		rs = results.NewResultSet(params[0])
		rs.AddResult(r)
	}
	return
}
