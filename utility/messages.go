//
// Copyright © 2016 Bryan T. Meyers <bmeyers@datadrake.com>
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

package utility

import "fmt"

var msgStart = map[string]string{
	"error":   "\033[31m[ERROR] ",
	"normal":  "\033[39m",
	"status":  "\033[34m[STATUS] ",
	"warning": "\033[33m[WARNING] ",
}

var msgEnd = "\033[39m"

/*
Errorf prints an error message
*/
func Errorf(fmts string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Printf(msgStart["error"] + fmts + msgEnd + "\n\n")
	} else {
		fmt.Printf(msgStart["error"]+fmts+msgEnd+"\n\n", args...)
	}
}

/*
Statusf prints a status message
*/
func Statusf(fmts string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Printf(msgStart["status"] + fmts + msgEnd + "\n\n")
	} else {
		fmt.Printf(msgStart["status"]+fmts+msgEnd+"\n\n", args...)
	}
}

/*
Warningf prints a warning message
*/
func Warningf(fmts string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Printf(msgStart["warning"] + fmts + msgEnd + "\n\n")
	} else {
		fmt.Printf(msgStart["warning"]+fmts+msgEnd+"\n\n", args...)
	}
}
