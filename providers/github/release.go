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
	"time"
)

// Release is a JSON representation fo a Github release
type Release struct {
	CreatedAt  string `json:"created_at"`
	Name       string `json:"name"`
	PreRelease bool   `json:"prerelease"`
	Tag        string `json:"tag_name"`
}

// Convert turns a Github release into a Cuppa release
func (cr *Release) Convert(name string) *results.Result {
	if cr.PreRelease {
		return nil
	}
	r := &results.Result{}
	r.Name = cr.Name
	r.Version = VersionRegex.FindString(cr.Tag)
	r.Published, _ = time.Parse(time.RFC3339, cr.CreatedAt)
	r.Location = fmt.Sprintf(SourceFormat, name, cr.Tag)
	return r
}
