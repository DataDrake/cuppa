//
// Copyright © 2017 Bryan T. Meyers <bmeyers@datadrake.com>
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
	"os"
)

// Quick fulfills the "Quick" subcommand
type Quick struct {}

// Name provides the name of this command
func (q Quick) Name() string {
    return "quick"
}

// Short prints a quick description of this command
func (q Quick) Short() string {
    return "Get the version and location of the most recent release"
}

// Usage prints a simple description of how to use this command
func (q Quick) Usage() {
	print("\t USAGE: cuppa quick <URL>")
}

// Flags builds the flagset for this command
func (q Quick) Flags() *flag.FlagSet {
	qcmd := flag.NewFlagSet("quick", flag.ExitOnError)
	qcmd.Usage = q.Usage
	return qcmd
}

/*
Execute quick for all providers
*/
func (q Quick) Execute() {
    ps := providers.All()
	flags := q.Flags()
	flags.Parse(os.Args[2:])
	found := false
	for _, p := range ps {
		name := p.Match(flags.Arg(0))
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
