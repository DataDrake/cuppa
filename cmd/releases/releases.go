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
	"fmt"
	"github.com/DataDrake/cuppa/providers"
	"github.com/DataDrake/cuppa/results"
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
		name := p.Match(rcmd.Arg(0))
		if name == "" {
			continue
		}
		rs, s := p.Releases(name)
		if s != results.OK {
			fmt.Fprintf(os.Stderr, "Failed to perform releases, code: %d\n", s)
			os.Exit(1)
		}
		rs.PrintAll()
	}
}
