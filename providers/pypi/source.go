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

package pypi

import (
	"github.com/DataDrake/cuppa/results"
	"time"
)

// Info contains a PyPi Version number
type Info struct {
	Version string `json:"version"`
}

// URL contains a JSON representation of a PyPi tarball URL
type URL struct {
	UploadTime string `json:"upload_time"`
	URL        string `json:"url"`
}

// LatestSource contains a JSON representation of a PyPi Source
type LatestSource struct {
	Info Info  `json:"info"`
	URLs []URL `json:"urls"`
}

// Convert turns a PyPi latest into a Cuppa Result
func (cr *LatestSource) Convert(name string) *results.Result {
	u := cr.URLs[len(cr.URLs)-1]
	published, _ := time.Parse(DateFormat, u.UploadTime)
	return results.NewResult(name, cr.Info.Version, u.URL, published)
}
