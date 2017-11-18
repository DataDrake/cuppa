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

package cpan

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"time"
)

// Release is a JSON representation of a CPAN release
type Release struct {
	Dist     string `json:"dist"`
	Archive  string `json:"archive"`
	Cpanid   string `json:"cpanid"`
	Version  string `json:"version"`
	Released string `json:"released"`
	Error    string `json:"error"`
	Status   string `json:"status"`
}

// Convert turns a CPAN release into a Cuppa result
func (cr *Release) Convert() *results.Result {
	if cr.Status != "stable" {
		return nil
	}
	r := &results.Result{}
	r.Name = cr.Dist
	r.Version = cr.Version
	r.Published, _ = time.Parse(time.RFC3339, cr.Released)
	r.Location = fmt.Sprintf(Source, cr.Cpanid[0:1], cr.Cpanid[0:2], cr.Cpanid, cr.Archive)
	return r
}
