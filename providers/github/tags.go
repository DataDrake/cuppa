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

package github

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/config"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Ref is a representation of a Github Git Reference
type Ref struct {
	Object struct {
		URL string `json:"url"`
	} `json:"object"`
}

// Refs is a list of Refs
type Refs []Ref

// ToTags converts Refs to Tags
func (rs Refs) ToTags() Tags {
	tags := make(Tags, 0)
	for _, ref := range rs {
		if len(ref.Object.URL) == 0 {
			continue
		}
		req, _ := http.NewRequest("GET", ref.Object.URL, nil)
		if key := config.Global.Github.Key; len(key) > 0 {
			req.Header["Authorization"] = []string{"token " + key}
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err.Error())
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}
		dec := json.NewDecoder(resp.Body)
		tag := Tag{}
		err = dec.Decode(&tag)
		if err != nil {
			panic(err.Error())
			continue
		}
		resp.Body.Close()
		tags = append(tags, tag)
	}
	return tags
}

// Tag is a JSON representation of a Github tag
type Tag struct {
	Tag    string `json:"tag"`
	Tagger struct {
		Date string `json:"date"`
	} `json:"tagger"`
}

// Convert turns a Github tag into a Cuppa result
func (t *Tag) Convert(name string) *results.Result {
	r := &results.Result{}
	pieces := strings.Split(name, "/")
	r.Name = pieces[len(pieces)-1]
	pieces = strings.Split(t.Tag, "/")
	r.Version = pieces[len(pieces)-1]
	r.Location = fmt.Sprintf(SourceFormat, name, r.Version)
	return r
}

// Tags are a JSON representation of one or more tags
type Tags []Tag

// Len gets the length of Tags
func (ts Tags) Len() int {
	return len(ts)
}

// Swap is used by Sort
func (ts Tags) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

// Less is used by Sort
func (ts Tags) Less(i, j int) bool {
	d1, err := time.Parse(time.RFC3339, ts[i].Tagger.Date)
	if err != nil {
		return true
	}
	d2, err := time.Parse(time.RFC3339, ts[j].Tagger.Date)
	if err != nil {
		return false
	}
	return d1.Before(d2)
}

func getTags(name string) (rs *results.ResultSet, s results.Status) {
	// Query the API
	req, _ := http.NewRequest("GET", fmt.Sprintf(APITags, name), nil)
	if key := config.Global.Github.Key; len(key) > 0 {
		req.Header["Authorization"] = []string{"token " + key}
	}
	resp, err := http.DefaultClient.Do(req)
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
	refs := make(Refs, 0)
	err = dec.Decode(&refs)
	if err != nil {
		panic(err.Error())
	}
	if len(refs) == 0 {
		s = results.NotFound
		return
	}
	tags := refs.ToTags()
	if len(tags) == 0 {
		s = results.NotFound
		return
	}
	rs = tags.Convert(name)
	return
}

// Convert a Github tagset to a Cuppa result set
func (ts Tags) Convert(name string) *results.ResultSet {
	sort.Sort(ts)
	rs := results.NewResultSet(name)
	for _, t := range ts {
		r := t.Convert(name)
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}
