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

package quick

import (
	"flag"
	"github.com/DataDrake/cuppa/providers"
	"github.com/DataDrake/cuppa/results"
	"os"
)

func quickUsage() {
	print("\t USAGE: cuppa quick <URL>")
}

func newQuickCMD() *flag.FlagSet {
	scmd := flag.NewFlagSet("quick", flag.ExitOnError)
	scmd.Usage = quickUsage
	return scmd
}

/*
Execute quick for all providers
*/
func Execute(ps []providers.Provider) {
	lcmd := newQuickCMD()
	lcmd.Parse(os.Args[2:])
	found := false
	for _, p := range ps {
		name := p.Match(lcmd.Arg(0))
		if name == "" {
			continue
		}
		r, s := p.Latest(name)
		if s != results.OK {
			continue
		}
		found = true
		r.PrintSimple()
	}
	if !found {
		println("No release found.")
		os.Exit(1)
	}
}
