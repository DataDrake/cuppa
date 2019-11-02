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
	"fmt"
	"strings"
	"time"

	"github.com/DataDrake/cuppa/results"
)

// Commit is a JSON representation of a GitLab commit
type Commit struct {
	AuthoredDate string `json:"authored_date"`
}

// Release is a JSON representation of a GitLab tag release
type Release struct {
	TagName string `json:"tag_name"`
}

// Tag is a JSON representation of a GitLab tag
type Tag struct {
	Name    string  `json:"name"`
	Commit  Commit  `json:"commit"`
	Release Release `json:"release"`
}

// Convert turns a GitLab tag into a Cuppa result
func (gl Tag) Convert(host, name string) *results.Result {
	published, _ := time.Parse(time.RFC3339, gl.Commit.AuthoredDate)
	file := fmt.Sprintf("%s-%s", strings.Split(name, "/")[1], gl.Name)
	loc := fmt.Sprintf(SourceFormat, host, name, gl.Name, file)
	version := gl.Release.TagName
	if len(version) == 0 {
		vs := strings.Split(gl.Name, "-")
		version = vs[len(vs)-1]
	}
	return results.NewResult(name, version, loc, published)
}
