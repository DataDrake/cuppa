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
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"strings"
)

// Tag is a JSON representation of a Github tag
type Tag struct {
	Name string `json:"name"`
}

// Convert turns a Github tag into a Cuppa result
func (t *Tag) Convert(name string) *results.Result {
	r := &results.Result{}
	pieces := strings.Split(name, "/")
	r.Name = pieces[len(pieces)-1]
	pieces = strings.Split(t.Name, "/")
	r.Version = pieces[len(pieces)-1]
	r.Location = fmt.Sprintf(SourceFormat, name, r.Version)
	return r
}

// Tags are a JSON representation of one or more tags
type Tags []Tag

// Convert a Github tagset to a Cuppa result set
func (ts *Tags) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, t := range *ts {
		r := t.Convert(name)
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}
