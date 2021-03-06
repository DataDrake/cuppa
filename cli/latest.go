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

package cli

import (
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/DataDrake/cuppa/providers"
	log "github.com/DataDrake/waterlog"
)

func init() {
	cmd.Register(&Latest)
}

// Latest gets the most recent release for a given source
var Latest = cmd.Sub{
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
func LatestRun(r *cmd.Root, c *cmd.Sub) {
	args := c.Args.(*LatestArgs)
	found := false
	for _, p := range providers.All() {
		log.Infof("\033[1m%s\033[22m checking for match:\n", p)
		match := p.Match(args.URL)
		if len(match) == 0 {
			log.Warnf("\033[1m%s\033[22m does not match.\n", p)
			continue
		}
		r, err := p.Latest(match)
		if err != nil {
			log.Warnf("Could not get latest \033[1m%s\033[22m, reason: %s\n", match[0], err)
			continue
		}
		found = true
		r.Print()
		log.Goodf("\033[1m%s\033[22m match(es) found.\n", p)
	}
	if !found {
		log.Fatalln("No release found.")
	}
	log.Goodln("Done")
}
