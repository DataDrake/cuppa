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

package search

import (
	"flag"
	"fmt"
	"github.com/DataDrake/cuppa/providers"
	"github.com/DataDrake/cuppa/results"
	"os"
)

func searchUsage() {
	print("\t USAGE: cuppa search <URL> or cuppa search <NAME>")
}

func newSearchCMD() *flag.FlagSet {
	scmd := flag.NewFlagSet("search", flag.ExitOnError)
	scmd.Usage = searchUsage
	return scmd
}

/*
Execute search for all providers
*/
func Execute(ps []providers.Provider) {
	scmd := newSearchCMD()
	scmd.Parse(os.Args[2:])
	for _, p := range ps {
		name := p.Match(scmd.Arg(0))
		if name == "" {
			continue
		}
		rs, s := p.Search(name)
		if s != results.OK {
			fmt.Fprintf(os.Stderr, "Failed to perform search, code: %d\n", s)
			os.Exit(1)
		}
		rs.PrintAll()
	}
}
