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

package releases

import (
	"flag"
	"github.com/DataDrake/cuppa/providers"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/utility"
	"os"
)

func releasesUsage() {
	print("\t USAGE: cuppa releases <URL> or cuppa releases <NAME>")
}

func newReleaserCMD() *flag.FlagSet {
	rcmd := flag.NewFlagSet("releases", flag.ExitOnError)
	rcmd.Usage = releasesUsage
	return rcmd
}

/*
Execute releases for all providers
*/
func Execute(ps []providers.Provider) {
	rcmd := newReleaserCMD()
	rcmd.Parse(os.Args[2:])
	for _, p := range ps {
		utility.Statusf("Checking provider '%s':", p.Name())
		name := p.Match(rcmd.Arg(0))
		if name == "" {
			utility.Warningf("Query does not supported by this provider.")
			continue
		}
		rs, s := p.Releases(name)
		if s != results.OK {
			utility.Errorf("Failed to fetch releases, code: %d", s)
			os.Exit(1)
		}
		rs.PrintAll()
	}
}
