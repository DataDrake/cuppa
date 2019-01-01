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

package rubygems

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"time"
)

// Version is a JSON representation of a version of a Gem
type Version struct {
	CreatedAt  string `json:"created_at"`
	PreRelease bool   `json:"prerelease"`
	Number     string `json:"number"`
}

// Convert turns a Rubygems version to a Cuppa result
func (cr *Version) Convert(name string) *results.Result {
	if cr.PreRelease {
		return nil
	}
	published, _ := time.Parse(time.RFC3339, cr.CreatedAt)
	location := fmt.Sprintf(SourceFormat, name, cr.Number)
	return results.NewResult(name, cr.Number, location, published)
}
