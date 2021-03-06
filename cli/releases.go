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
	cmd.Register(&Releases)
}

// Releases gets all releases for a given source
var Releases = cmd.Sub{
	Name:  "releases",
	Alias: "r",
	Short: "Get all stable releases",
	Args:  &ReleasesArgs{},
	Run:   ReleasesRun,
}

// ReleasesArgs contains the arguments for the "releases" subcommand
type ReleasesArgs struct {
	URL string `desc:"Location of a previous source archive"`
}

// ReleasesRun carries out finding all releases
func ReleasesRun(r *cmd.Root, c *cmd.Sub) {
	args := c.Args.(*ReleasesArgs)
	found := false
	for _, p := range providers.All() {
		log.Infof("\033[1m%s\033[22m checking for match:\n", p)
		match := p.Match(args.URL)
		if len(match) == 0 {
			log.Warnf("\033[1m%s\033[22m does not match.\n", p)
			continue
		}
		rs, err := p.Releases(match)
		if err != nil {
			log.Warnf("Could not get latest \033[1m%s\033[22m, reason: %s\n", match[0], err)
			continue
		}
		found = true
		rs.PrintAll()
		log.Goodf("\033[1m%s\033[22m match(es) found.\n", p)
	}
	if !found {
		log.Fatalln("No release found.")
	}
	log.Goodln("Done")
}
