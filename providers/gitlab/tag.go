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
func (gl Tag) Convert(name string, isFreedesktop bool) *results.Result {
	published, _ := time.Parse(time.RFC3339, gl.Commit.AuthoredDate)
	file := fmt.Sprintf("%s-%s", strings.Split(name, "/")[1], gl.Name)

	// For some projects, GitLab does not return a release object in the
	// JSON response.
	var version string
	if gl.Release.TagName != "" {
		version = gl.Release.TagName
	} else {
		version = gl.Name
	}

	var loc string
	if isFreedesktop {
		loc = fmt.Sprintf(SourceFormatFreedesktop, name, gl.Name, file)
	} else {
		loc = fmt.Sprintf(SourceFormat, name, gl.Name, file)
	}

	return results.NewResult(name, version, loc, published)
}
