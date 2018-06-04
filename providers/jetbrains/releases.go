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

package jetbrains

import (
	"github.com/DataDrake/cuppa/results"
)

// Releases is a collection of JetBrains releases
type Releases map[string][]Release

// Convert turns JetBrains releases into a Cuppa result set
func (jbs Releases) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	code := ReleaseCodes[name]
	for _, rel := range jbs[code] {
		r := rel.Convert()
		if r != nil {
			r.Name = name
			rs.AddResult(r)
		}
	}
	return rs
}
