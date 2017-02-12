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
	"time"
)

var pypiAPI = "https://pypi.python.org/pypi/%s/json"
var pypiRegex = regexp.MustCompile("https?://pypi.python.org/packages/\\w+/\\w+/\\w+/(.+)-.+.tar.gz")

const pythonTime = "2006-01-02T15:04:05"

type pypiInfo struct {
	Version string `json:"version"`
}

type pypiURL struct {
	UploadTime string `json:"upload_time"`
	URL        string `json:"url"`
}

type pypiLatest struct {
	Info pypiInfo  `json:"info"`
	URLs []pypiURL `json:"urls"`
}

func (cr *pypiLatest) Convert(name string) *results.Result {
	r := &results.Result{}
	r.Name = name
	r.Version = cr.Info.Version
	u := cr.URLs[len(cr.URLs)-1]
	r.Published, _ = time.Parse(pythonTime, u.UploadTime)
	r.Location = u.URL
	return r
}

func converURLs(cr []pypiURL, name, version string) *results.Result {
	r := &results.Result{}
	r.Name = name
	r.Version = version
	u := cr[len(cr)-1]
	r.Published, _ = time.Parse(pythonTime, u.UploadTime)
	r.Location = u.URL
	return r
}

type pypiResultSet struct {
	Releases map[string][]pypiURL
}

func (crs *pypiResultSet) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for ver, rel := range crs.Releases {
		r := converURLs(rel, name, ver)
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}

/*
PyPiProvider is the upstream provider interface for pypi
*/
type PyPiProvider struct{}

/*
Latest finds the newest release for a pypi package
*/
func (c PyPiProvider) Latest(name string) (r *results.Result, s results.Status) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(pypiAPI, name))
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
	cr := &pypiLatest{}
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
func (c PyPiProvider) Match(query string) string {
	sm := pypiRegex.FindStringSubmatch(query)
	if len(sm) != 2 {
		return ""
	}
	return sm[1]
}

/*
Name gives the name of this provider
*/
func (c PyPiProvider) Name() string {
	return "PyPi"
}

/*
Releases finds all matching releases for a pypi package
*/
func (c PyPiProvider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	//Query the API
	resp, err := http.Get(fmt.Sprintf(pypiAPI, name))
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
	crs := &pypiResultSet{}
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
