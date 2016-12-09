//
// Copyright © 2016 Bryan T. Meyers <bmeyers@datadrake.com>
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

package results

import (
	"fmt"
	"net/url"
	"time"
)

/*
Result contains the information for a single query result
*/
type Result struct {
	Name      string
	Version   string
	Location  *url.URL
	Published time.Time
}

/*
NewResult creates a result with the specified values
*/
func NewResult(name, version string, location *url.URL, published time.Time) *Result {
	return &Result{name, version, location, published}
}

/*
Print pretty-prints a single Result
*/
func (r *Result) Print() {
	fmt.Printf("%-10s: %s\n", "Name", r.Name)
	fmt.Printf("%-10s: %s\n", "Version", r.Version)
	if r.Location != nil {
		fmt.Printf("%-10s: %s\n", "Location", r.Location.String())
	}
	fmt.Printf("%-10s: %s\n", "Published", r.Published.Format(time.RFC3339))
	fmt.Println()
}
