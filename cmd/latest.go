//
// Copyright 2017 Bryan T. Meyers <bmeyers@datadrake.com>
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

package cmd

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

// Latest fulfills the "latest" subcommand
type Latest struct{}

// Name provides the name of this command
func (l Latest) Name() string {
	return "latest"
}

// Short provides a quick description of this command
func (l Latest) Short() string {
	return "Get the latest stable release"
}

// Usage prints a general usage statement
func (l Latest) Usage() {
	print("USAGE: cuppa latest <URL>\n\n")
	l.Flags().PrintDefaults()
}

// Flags builds the flagset for this command
func (l Latest) Flags() *flag.FlagSet {
	scmd := flag.NewFlagSet("latest", flag.ExitOnError)
	scmd.Usage = l.Usage
	return scmd
}

/*
Execute releases for all providers
*/
func (l Latest) Execute() {
	ps := providers.All()
	flags := l.Flags()
	flags.Parse(os.Args[2:])
	w := waterlog.New(os.Stdout, "", log.Ltime)
	w.SetLevel(level.Info)
	w.SetFormat(format.Min)
	found := false
	for _, p := range ps {
		w.Infof("\033[1m%s\033[21m checking for match:\n", p.Name())
		name := p.Match(flags.Arg(0))
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
