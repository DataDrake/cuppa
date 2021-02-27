//
// Copyright 2016-2021 Bryan T. Meyers <root@datadrake.com>
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
	"github.com/DataDrake/cuppa/version"
	"os"
	"text/tabwriter"
	"time"
)

// Result contains the information for a single query result
type Result struct {
	Name      string
	Version   version.Version
	Location  string
	Published time.Time
}

// NewResult creates a result with the specified values
func NewResult(name, v string, location string, published time.Time) *Result {
	r := &Result{name, version.NewVersion(v), location, published}
	if r.Published.IsZero() {
		r.Published = r.Version.FindDate()
	}
	return r
}

// Print pretty-prints a single Result
func (r *Result) Print() {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(tw, "%s\t: %s\n", "Name", r.Name)
	fmt.Fprintf(tw, "%s\t: %s\n", "Version", r.Version)
	if r.Location != "" {
		fmt.Fprintf(tw, "%s\t: %s\n", "Location", r.Location)
	}
	if !r.Published.IsZero() {
		fmt.Fprintf(tw, "%s\t: %s\n", "Published", r.Published.Format(time.RFC3339))
	}
	tw.Flush()
	fmt.Println()
}

// PrintSimple only prints the version and the location of the latest release
func (r *Result) PrintSimple() {
	fmt.Printf("%s %s\n", r.Version, r.Location)
}
