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

package cpan

import (
	"github.com/DataDrake/cuppa/results"
	"time"
)

// Release is a JSON representation of a CPAN release
type Release struct {
	Version  string `json:"version"`
	Status   string `json:"status"`
	Date     string `json:"date"`
	Location string `json:"download_url"`
	Error    string `json:"error"`
}

// Convert turns a CPAN release into a Cuppa result
func (cr *Release) Convert(name string) *results.Result {
	if cr.Status != "latest" {
		return nil
	}
	r := &results.Result{}
	r.Name = name
	r.Version = cr.Version
	r.Published, _ = time.Parse(time.RFC3339, cr.Date)
	r.Location = cr.Location
	return r
}
