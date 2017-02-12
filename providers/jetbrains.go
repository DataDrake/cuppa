//
// Copyright © 2016 Bryan T. Meyers <bmeyers@datadrake.com>
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

var releaseCode = map[string]string{
	"appcode":              "AC",
	"clion":                "CL",
	"datagrip":             "DG",
	"ideaiu":               "IIU",
	"ideaic":               "IIC",
	"phpstorm":             "PS",
	"pycharm-professional": "PCP",
	"pycharm-ce":           "PCC",
	"rubymine":             "RM",
	"webstorm":             "WS",
}

var jetbrainsAPI = "https://data.services.jetbrains.com/products/releases?code=%s"
var jetbrainsAPILatest = "https://data.services.jetbrains.com/products/releases?code=%s&latest=true"
var jetbrainsRegex = regexp.MustCompile("https://download.jetbrains.com/.+?/(.+)-\\d.*")

type jetbrainsDownload struct {
	ChecksumLink string `json:"checksumLink"`
	Link         string `json:"link"`
	Size         uint64 `json:"size"`
}

type jetbrainsRelease struct {
	Build        string                       `json:"build"`
	Date         string                       `json:"date"`
	Downloads    map[string]jetbrainsDownload `json:"downloads"`
	Notes        string                       `json:"notesLink"`
	Type         string                       `json:"type"`
	MajorVersion string                       `json:"majorVersion"`
	Version      string                       `json:"version"`
}

func (jb jetbrainsRelease) Convert() *results.Result {
	r := &results.Result{}
	r.Version = jb.Version
	r.Published, _ = time.Parse("2006-01-02", jb.Date)
	if d, ok := jb.Downloads["linuxWithoutJDK"]; ok {
		r.Location = d.Link
		return r
	}
	if d, ok := jb.Downloads["linux"]; ok {
		r.Location = d.Link
		return r
	}
	return r
}

type jetbrainsResultSet map[string][]jetbrainsRelease

func (jbs jetbrainsResultSet) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	code := releaseCode[name]
	for _, rel := range jbs[code] {
		r := rel.Convert()
		if r != nil {
			r.Name = name
			rs.AddResult(r)
		}
	}
	return rs
}

/*
JetBrainsProvider is the upstream provider interface for JetBrains
*/
type JetBrainsProvider struct{}

/*
Latest finds the newest release for a JetBrains package
*/
func (c JetBrainsProvider) Latest(name string) (r *results.Result, s results.Status) {
	//Query the API
	code := releaseCode[name]
	resp, err := http.Get(fmt.Sprintf(jetbrainsAPI, code))
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
	jbs := make(jetbrainsResultSet)
	err = dec.Decode(&jbs)
	if err != nil {
		panic(err.Error())
	}
	if jbs[code] == nil || len(jbs[code]) == 0 {
		s = results.NotFound
		return
	}
	r = jbs.Convert(name).First()
	return
}

/*
Match checks to see if this provider can handle this kind of query
*/
func (c JetBrainsProvider) Match(query string) string {
	sm := jetbrainsRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}
	return strings.ToLower(sm[1])
}

/*
Name gives the name of this provider
*/
func (c JetBrainsProvider) Name() string {
	return "JetBrains"
}

/*
Releases finds all matching releases for a JetBrains package
*/
func (c JetBrainsProvider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	//Query the API
	code := releaseCode[name]
	resp, err := http.Get(fmt.Sprintf(jetbrainsAPI, code))
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
	jbs := make(jetbrainsResultSet)
	err = dec.Decode(&jbs)
	if err != nil {
		panic(err.Error())
	}
	if jbs[code] == nil || len(jbs[code]) == 0 {
		s = results.NotFound
		return
	}
	rs = jbs.Convert(name)
	return
}
