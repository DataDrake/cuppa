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

package cli

import (
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/DataDrake/cuppa/providers"
	"github.com/DataDrake/cuppa/results"
	"os"
)

// Latest gets the most recent release for a given source
var Latest = cmd.CMD{
	Name:  "latest",
	Alias: "l",
	Short: "Get the latest stable release",
	Args:  &LatestArgs{},
	Run:   LatestRun,
}

// LatestArgs contains the arguments for the "latest" subcommand
type LatestArgs struct {
	URL string `desc:"Location of a previous source archive"`
}

// LatestRun carries out finding the latest release
func LatestRun(r *cmd.RootCMD, c *cmd.CMD) {
	args := c.Args.(*LatestArgs)
	found := false
	for _, p := range providers.All() {
		log.Infof("\033[1m%s\033[21m checking for match:\n", p.Name())
		name := p.Match(args.URL)
		if name == "" {
			log.Warnf("\033[1m%s\033[21m does not match.\n", p.Name())
			continue
		}
		r, s := p.Latest(name)
		if s != results.OK {
			log.Warnf("Could not get latest \033[1m%s\033[21m, code: %d\n", name, s)
			continue
		}
		found = true
		r.Print()
		log.Goodf("\033[1m%s\033[21m match(es) found.\n", p.Name())
	}
	if found {
		log.Goodln("Done")
	} else {
		log.Fatalln("No release found.")
	}
	os.Exit(0)
}
