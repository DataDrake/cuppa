//
// Copyright 2016-2017 Bryan T. Meyers <bmeyers@datadrake.com>
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

package rubygems

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
)

// LatestVersion is a JSON representation of the latest Version of a Gem
type LatestVersion struct {
	Version string `json:"version"`
}

// Convert turns a Rubygems latest release into a Cuppa result
func (cr *LatestVersion) Convert(name string) *results.Result {
	r := &results.Result{}
	r.Name = name
	r.Version = cr.Version
	r.Location = fmt.Sprintf(SourceFormat, name, cr.Version)
	return r
}
