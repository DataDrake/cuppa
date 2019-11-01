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

package gitlab

import (
	"github.com/DataDrake/cuppa/results"
)

// Tags is a set of one or more GitLab tags
type Tags struct {
	Tags []Tag
}

// Convert turns a GitLab result set into a Cuppa ResultSet
func (gls Tags) Convert(name string, isFreedesktop bool) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, rel := range gls.Tags {
		r := rel.Convert(name, isFreedesktop)
		if r != nil {
			rs.AddResult(r)
		}
	}

	return rs
}
