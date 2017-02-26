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
	"os"
)

// Help fulfills the "help" subcommand
type Help struct{}

// Name provides the name of this command
func (h Help) Name() string {
	return "help"
}

// Short provides a quick description of this command
func (h Help) Short() string {
	return "Get help with a specific subcomand"
}

// Usage prints a general usage statement
func (h Help) Usage() {
	print("USAGE: cuppa help <subcomand>\n\n")
	h.Flags().PrintDefaults()
}

// Flags builds the flagset for this command
func (h Help) Flags() *flag.FlagSet {
	scmd := flag.NewFlagSet("help", flag.ExitOnError)
	scmd.Usage = h.Usage
	return scmd
}

/*
Execute releases for all providers
*/
func (h Help) Execute() {

	flags := h.Flags()
	flags.Parse(os.Args[2:])

	if flags.NArg() != 1 {
		h.Usage()
		os.Exit(1)
	}

	found := false
	for _, c := range All {
		if flags.Arg(0) == c.Name() {
			c.Usage()
			found = true
			break
		}
	}
	if !found {
		h.Usage()
		os.Exit(1)
	}
}
