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

package main

import (
	"github.com/DataDrake/cuppa/cmd"
	"os"
)

func main() {

	// Setup the Sub-Commands
	cmd.RegisterCMD("help", cmd.Help{})
	cmd.RegisterCMD("latest", cmd.Latest{})
	cmd.RegisterCMD("releases", cmd.Releases{})
	cmd.RegisterCMD("quick", cmd.Quick{})

	// Run the program
	cmd.Run()

	// Terminate Gracefully
	os.Exit(0)
}
