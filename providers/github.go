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

var githubAPILatest = "https://api.github.com/repos/%s/releases/latest"
var githubAPIReleases = "https://api.github.com/repos/%s/releases"
var githubAPITags = "https://api.github.com/repos/%s/git/refs/tags"
var githubSource = "https://github.com/%s/archive/%s.tar.gz"
var githubRegex = regexp.MustCompile("https?://github.com/([^/]*/[^/]*)/.*/[^/]*.tar.gz")
var githubVersionRegex = regexp.MustCompile("(?:\\d+\\.)*\\d+\\w*")

type githubRelease struct {
	CreatedAt  string `json:"created_at"`
	Name       string `json:"name"`
	PreRelease bool   `json:"prerelease"`
	Tag        string `json:"tag_name"`
}

func (cr *githubRelease) Convert(name string) *results.Result {
	if cr.PreRelease {
		return nil
	}
	r := &results.Result{}
	r.Name = cr.Name
	r.Version = githubVersionRegex.FindString(cr.Tag)
	r.Published, _ = time.Parse(time.RFC3339, cr.CreatedAt)
	r.Location = fmt.Sprintf(githubSource, name, cr.Tag)
	return r
}

type githubResultSet []githubRelease

func (crs *githubResultSet) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, rel := range *crs {
		r := rel.Convert(name)
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}

type githubTag struct {
	Ref string `json:"ref"`
}

func (ts *githubTags) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, t := range *ts {
		r := t.Convert(name)
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}

type githubTags []githubTag

func (t *githubTag) Convert(name string) *results.Result {
	r := &results.Result{}
	pieces := strings.Split(name, "/")
	r.Name = pieces[len(pieces)-1]
	pieces = strings.Split(t.Ref, "/")
	r.Version = pieces[len(pieces)-1]
	r.Location = fmt.Sprintf(githubSource, name, r.Version)
	return r
}

/*
GitHubProvider is the upstream provider interface for github
*/
type GitHubProvider struct{}

/*
Latest finds the newest release for a github package
*/
func (c GitHubProvider) Latest(name string) (r *results.Result, s results.Status) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(githubAPILatest, name))
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	//Translate Status Code
	switch resp.StatusCode {
	case 200:
		s = results.OK
	case 404:
		rs, st := getTags(name)
		s = st
		if s != results.OK || rs.Empty() {
			return
		}
		r = rs.Last()
		return
	default:
		s = results.Unavailable
	}

	//Fail if not OK
	if s != results.OK {
		return
	}

	dec := json.NewDecoder(resp.Body)
	cr := &githubRelease{}
	err = dec.Decode(cr)
	if err != nil {
		panic(err.Error())
	}
	r = cr.Convert(name)
	return
}

/*
Match checks to see if this provider can handle this kind of query
*/
func (c GitHubProvider) Match(query string) string {
	sm := githubRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}
	return sm[1]
}

/*
Name gives the name of this provider
*/
func (c GitHubProvider) Name() string {
	return "GitHub"
}

func getTags(name string) (rs *results.ResultSet, s results.Status) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(githubAPITags, name))
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
	tags := make(githubTags, 0)
	err = dec.Decode(&tags)
	if err != nil {
		panic(err.Error())
	}
	if len(tags) == 0 {
		s = results.NotFound
		return
	}
	rs = tags.Convert(name)
	return
}

/*
Releases finds all matching releases for a github package
*/
func (c GitHubProvider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(githubAPIReleases, name))
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
	println(s)

	//Fail if not OK
	if s != results.OK {
		return
	}

	dec := json.NewDecoder(resp.Body)
	crs := make(githubResultSet, 0)
	err = dec.Decode(&crs)
	if err != nil {
		panic(err.Error())
	}
	if len(crs) == 0 {
		rs, s = getTags(name)
		return
	}
	rs = crs.Convert(name)
	return
}
