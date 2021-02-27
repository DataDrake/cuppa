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
	log "github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	log2 "log"
)

// GlobalFlags contains the flags for all commands
type GlobalFlags struct{}

// Root is the main command for this application
var Root *cmd.RootCMD

func init() {
	// Build Application
	Root = &cmd.RootCMD{
		Name:  "cuppa",
		Short: "Comprehensive Upstream Provider Polling Assistant",
		Flags: &GlobalFlags{},
	}
	// Setup the Sub-Commands
	Root.RegisterCMD(&cmd.Help)
	Root.RegisterCMD(&Latest)
	Root.RegisterCMD(&Quick)
	Root.RegisterCMD(&Releases)

	//Set up logging
	log.SetFlags(log2.Ltime)
	log.SetLevel(level.Info)
	log.SetFormat(format.Min)
}
