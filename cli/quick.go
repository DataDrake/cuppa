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
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/DataDrake/cuppa/providers"
	log "github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
)

// Quick gets the most recent release for a given source, without pretty printing
var Quick = cmd.CMD{
	Name:  "quick",
	Alias: "q",
	Short: "Get the version and location of the most recent release",
	Args:  &QuickArgs{},
	Run:   QuickRun,
}

// QuickArgs contains the arguments for the "quick" subcommand
type QuickArgs struct {
	URL string `desc:"Location of a previous source archive"`
}

// QuickRun carries out finding the latest release
func QuickRun(r *cmd.RootCMD, c *cmd.CMD) {
	args := c.Args.(*QuickArgs)
	found := false
	log.SetFormat(format.Un)
	for _, p := range providers.All() {
		match := p.Match(args.URL)
		if len(match) == 0 {
			continue
		}
		r, err := p.Latest(match)
		if err != nil {
			continue
		}
		found = true
		r.PrintSimple()
		break
	}
	if !found {
		log.Fatalln("No release found.")
	}
}
