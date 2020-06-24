//
// Copyright 2016-2020 Bryan T. Meyers <root@datadrake.com>
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

package hackage

import (
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"time"
)

// Versions is a JSON representation of Hackage release version numbers
type Versions struct {
	Normal []string `json:"normal-version"`
}

// Release is a local representation of a Hackage release
type Release struct {
	name     string
	released string
	version  string
}

// Convert turns a Hackage release into a Cuppa result
func (hr *Release) Convert() *results.Result {
	pub, _ := time.Parse(time.UnixDate, hr.released)
	return results.NewResult(hr.name, hr.version, fmt.Sprintf(TarballAPI, hr.name, hr.version, hr.name, hr.version), pub)
}
