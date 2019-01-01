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
	"time"
)

// Download is a JSON representation of a JetBrains downloadable source
type Download struct {
	ChecksumLink string `json:"checksumLink"`
	Link         string `json:"link"`
	Size         uint64 `json:"size"`
}

// Release is a JSON representation of a JetBrains release
type Release struct {
	Build        string              `json:"build"`
	Date         string              `json:"date"`
	Downloads    map[string]Download `json:"downloads"`
	Notes        string              `json:"notesLink"`
	Type         string              `json:"type"`
	MajorVersion string              `json:"majorVersion"`
	Version      string              `json:"version"`
}

// Convert turns a JetBrains release into a Cuppa result
func (jb Release) Convert() *results.Result {
	published, _ := time.Parse("2006-01-02", jb.Date)
	if d, ok := jb.Downloads["linuxWithoutJDK"]; ok {
		return results.NewResult("", jb.Version, d.Link, published)
	}
	if d, ok := jb.Downloads["linux"]; ok {
		return results.NewResult("", jb.Version, d.Link, published)
	}
	return nil
}
