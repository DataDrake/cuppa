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

package latest

import (
	"flag"
	"github.com/DataDrake/cuppa/providers"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"log"
	"os"
)

func latestUsage() {
	print("\t USAGE: cuppa latest <URL>")
}

func newLatestCMD() *flag.FlagSet {
	scmd := flag.NewFlagSet("latest", flag.ExitOnError)
	scmd.Usage = latestUsage
	return scmd
}

/*
Execute releases for all providers
*/
func Execute(ps []providers.Provider) {
	lcmd := newLatestCMD()
	lcmd.Parse(os.Args[2:])
	w := waterlog.New(os.Stdout, "", log.Ltime)
	w.SetLevel(level.Info)
	w.SetFormat(format.Min)
	found := false
	for _, p := range ps {
		w.Infof("\033[1m%s\033[21m checking for match:\n", p.Name())
		name := p.Match(lcmd.Arg(0))
		if name == "" {
			w.Warnf("\033[1m%s\033[21m does not match.\n", p.Name())
			continue
		}
		r, s := p.Latest(name)
		if s != results.OK {
			w.Warnf("Could not get latest \033[1m%s\033[21m, code: %d\n", name, s)
			continue
		}
		found = true
		r.Print()
		w.Goodf("\033[1m%s\033[21m match(es) found.\n", p.Name())
	}
	if found {
		w.Goodln("Done")
	} else {
		w.Fatalln("No release found.")
	}
}
