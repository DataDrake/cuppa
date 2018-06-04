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

package launchpad

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"time"
)

// Release is an internal representation of a Launchpad Release
type Release struct {
	name     string
	release  string
	series   string
	uploaded string
}

// Convert turns a Launchpad release into a Cuppa result
func (lr *Release) Convert() *results.Result {
	r := &results.Result{}
	r.Name = lr.name
	r.Version = lr.release
	r.Published, _ = time.Parse(time.RFC3339, lr.uploaded)
	r.Location = fmt.Sprintf(SourceFormat, lr.name, lr.series, lr.release, lr.name, lr.release)
	return r
}

// Releases holds one or more Launchpad Releases
type Releases []Release

// Convert turns a Launchpad result set to a Cuppa result set
func (lrs Releases) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, rel := range lrs {
		r := rel.Convert()
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}
