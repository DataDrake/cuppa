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

	if len(os.Args) < 2 {
		cmd.Usage()
		os.Exit(1)
	}

	found := false
	for _, c := range cmd.All {
		if os.Args[1] == c.Name() {
			c.Execute()
			found = true
			break
		}
	}
	if !found {
		cmd.Usage()
		os.Exit(1)
	}

	os.Exit(0)
}
