//
// Copyright 2017 Bryan T. Meyers <bmeyers@datadrake.com>
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

package providers

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var cpanAPIDist = "http://search.cpan.org/api/dist/%s"
var cpanSource = "http://search.cpan.org/CPAN/authors/id/%s/%s/%s/%s"
var cpanRegex = regexp.MustCompile("http://search.cpan.org/CPAN/authors/id/(.*)")

type cpanRelease struct {
	Dist     string `json:"dist"`
	Archive  string `json:"archive"`
	Cpanid   string `json:"cpanid"`
	Version  string `json:"version"`
	Released string `json:"released"`
	Error    string `json:"error"`
	Status   string `json:"status"`
}

func (cr *cpanRelease) Convert() *results.Result {
	if cr.Status != "stable" {
		return nil
	}
	r := &results.Result{}
	r.Name = cr.Dist
	r.Version = cr.Version
	r.Published, _ = time.Parse(time.RFC3339, cr.Released)
	r.Location = fmt.Sprintf(cpanSource, cr.Cpanid[0:1], cr.Cpanid[0:2], cr.Cpanid, cr.Archive)
	return r
}

type cpanResultSet struct {
	Releases []cpanRelease `json:"releases"`
}

func (crs *cpanResultSet) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, rel := range crs.Releases {
		r := rel.Convert()
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}

/*
CPANProvider is the upstream provider interface for CPAN
*/
type CPANProvider struct{}

/*
Latest finds the newest release for a CPAN package
*/
func (c CPANProvider) Latest(name string) (r *results.Result, s results.Status) {
	rs, s := c.Releases(name)
	//Fail if not OK
	if s != results.OK {
		return
	}
	r = rs.First()
	return
}

/*
Match checks to see if this provider can handle this kind of query
*/
func (c CPANProvider) Match(query string) string {
	sm := cpanRegex.FindStringSubmatch(query)
	if len(sm) == 0 {
		return ""
	}
	sms := strings.Split(sm[1], "/")
	filename := sms[len(sms)-1]
	pieces := strings.Split(filename, "-")
	pieces = pieces[0 : len(sms)-2]
	return strings.Join(pieces, "-")
}

/*
Name gives the name of this provider
*/
func (c CPANProvider) Name() string {
	return "CPAN"
}

/*
Releases finds all matching releases for a CPAN package
*/
func (c CPANProvider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(cpanAPIDist, name))
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	//Translate Status Code
	switch resp.StatusCode {
	case 200:
		s = results.OK
	case 404:
		s = results.NotFound
	default:
		s = results.Unavailable
	}

	//Fail if not OK
	if s != results.OK {
		return
	}

	dec := json.NewDecoder(resp.Body)
	crs := &cpanResultSet{}
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
