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

package pypi

import (
	"github.com/DataDrake/cuppa/results"
	"time"
)

// ConvertURLS translates PyPi URLs to Cuppa results
func ConvertURLS(cr []URL, name, version string) *results.Result {
	u := cr[len(cr)-1]
	published, _ := time.Parse(DateFormat, u.UploadTime)
	return results.NewResult(name, version, u.URL, published)
}

// Releases holds one or more Source URLs
type Releases struct {
	Releases map[string][]URL `json:"releases"`
}

// Convert turns PyPi releases into a Cuppa results set
func (crs *Releases) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for ver, rel := range crs.Releases {
		if r := ConvertURLS(rel, name, ver); r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}
